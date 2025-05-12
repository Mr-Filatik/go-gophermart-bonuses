package service

import (
	"errors"
	"strconv"

	dbModels "github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/database/models"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/database/repository"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/server/models"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/worker"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/helper"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
)

var (
	ErrLoginAlreadyTaken      = errors.New("login already taken")
	ErrInvalidLoginOrPassword = errors.New("invalid login/password pair")
	ErrInvalidOrderNumber     = errors.New("invalid order number")
	ErrInsufficientFunds      = errors.New("insufficient funds")
	ErrAlreadyUploadThisUser  = errors.New("already upload this user")
	ErrAlreadyUploadOtherUser = errors.New("already upload other user")
)

type Service struct {
	userRep   repository.IUserRepository
	orderRep  repository.IOrderRepository
	wthrwlRep repository.IWithdrawalRepository
	log       logger.Logger
	worker    *worker.Worker
}

func New(
	userRep repository.IUserRepository,
	orderRep repository.IOrderRepository,
	wthrwlRep repository.IWithdrawalRepository,
	worker *worker.Worker,
	log logger.Logger) *Service {
	srv := Service{
		userRep:   userRep,
		orderRep:  orderRep,
		wthrwlRep: wthrwlRep,
		worker:    worker,
		log:       log,
	}

	log.Info("Service created")

	return &srv
}

func (s *Service) UserRegister(data models.UserRegisterRequest) error {
	s.log.Debug("Service.UserRegister() was called.", "register_login", data.Login)

	hashedPassword, passErr := helper.GeneratePasswordHash(data.Password)
	if passErr != nil {
		return errors.New(passErr.Error())
	}

	user := dbModels.User{
		Login:        data.Login,
		PasswordHash: hashedPassword,
	}

	err := s.userRep.Create(&user)
	if err != nil {
		if errors.Is(err, repository.ErrEntityAlreadyExists) {
			return ErrLoginAlreadyTaken
		}
		return errors.New(err.Error())
	}

	return nil
}

func (s *Service) UserLogin(data models.UserLoginRequest) error {
	s.log.Debug("Service.UserLogin() was called.", "login_login", data.Login)

	user, err := s.userRep.GetByLogin(data.Login)
	if err != nil {
		if errors.Is(err, repository.ErrEntityNotFound) {
			return ErrInvalidLoginOrPassword
		}
		return errors.New(err.Error())
	}

	ok := helper.ComparePasswordHashes(user.PasswordHash, data.Password)
	if !ok {
		return ErrInvalidLoginOrPassword
	}

	return nil
}

func (s *Service) UserOrdersGet(login string) ([]models.UserOrder, error) {
	s.log.Debug("Service.UserOrdersGet() was called.", "order_login", login)

	user, err := s.userRep.GetByLogin(login)
	if err != nil {
		return make([]models.UserOrder, 0), errors.New(err.Error())
	}

	orders, oerr := s.orderRep.GetAllByUserID(user.ID)
	if oerr != nil {
		return make([]models.UserOrder, 0), errors.New(oerr.Error())
	}

	userOrders := make([]models.UserOrder, len(orders))
	for i := range orders {
		userOrders[i] = models.UserOrder{
			Number:     strconv.FormatUint(orders[i].Number, 10),
			Status:     models.UserOrderStatus(orders[i].Status),
			Accrual:    helper.ConvertPriceToFloat64(orders[i].Accrual),
			UploadedAt: orders[i].UploadedAt, // .Format(time.RFC3339)
		}
	}

	return userOrders, nil
}

func (s *Service) UserOrdersCreate(login string, number string) error {
	s.log.Debug(
		"Service.UserOrdersCreate() was called.",
		"order_login", login,
		"number", number,
	)

	num, err := strconv.ParseUint(number, 10, 64)
	if err != nil {
		return errors.New(err.Error())
	}

	if ok := helper.ValidateAlgorithmLuhn(num); !ok {
		return ErrInvalidOrderNumber
	}

	order, err := s.orderRep.GetByNumber(num)
	if err != nil {
		if errors.Is(err, repository.ErrEntityNotFound) {
			usr, uerr := s.userRep.GetByLogin(login)
			if uerr != nil {
				return errors.New(uerr.Error())
			}
			createdOrder := dbModels.UserOrder{
				Number: num,
				UserID: usr.ID,
				User:   *usr,
			}
			cerr := s.orderRep.Create(&createdOrder)
			if cerr != nil {
				return errors.New(cerr.Error())
			}

			s.worker.AddTask(num)

			return nil
		}
		return errors.New(err.Error())
	}

	// autoload user entity
	user, usErr := s.userRep.GetByID(order.UserID)
	if usErr != nil {
		return errors.New(usErr.Error())
	}

	if user.Login == login {
		return ErrAlreadyUploadThisUser
	} else {
		return ErrAlreadyUploadOtherUser
	}
}

func (s *Service) UserBalanceGet(login string) (models.UserBalanceResponse, error) {
	s.log.Debug("Service.UserBalanceGet() was called.", "login", login)

	user, uerr := s.userRep.GetByLogin(login)
	if uerr != nil {
		return models.UserBalanceResponse{}, errors.New(uerr.Error())
	}

	return models.UserBalanceResponse{
		Current:   helper.ConvertPriceToFloat64(user.Current),
		Withdrawn: helper.ConvertPriceToFloat64(user.Withdrawn),
	}, nil
}

func (s *Service) UserBalanceWithdraw(data models.UserBalanceWithdrawRequest, login string) error {
	s.log.Debug(
		"Service.UserBalanceWithdraw() was called.",
		"order", data.Order,
		"sum", data.Sum,
	)

	sum := helper.ConvertPriceToUint64(data.Sum)

	num, perr := strconv.ParseUint(data.Order, 10, 64)
	if perr != nil {
		return errors.New(perr.Error())
	}
	if ok := helper.ValidateAlgorithmLuhn(num); !ok {
		return ErrInvalidOrderNumber
	}

	user, uerr := s.userRep.GetByLogin(login)
	if uerr != nil {
		return errors.New(uerr.Error())
	}

	if sum > user.Current {
		return ErrInsufficientFunds
	}

	// BEGIN TRANSACTION

	createdWithdrawal := dbModels.UserWithdrawal{
		Order:  data.Order,
		Sum:    sum,
		UserID: user.ID,
		User:   *user,
	}
	err := s.wthrwlRep.Create(&createdWithdrawal)
	if err != nil {
		return errors.New(err.Error())
	}
	user.Current -= sum
	user.Withdrawn += sum
	updErr := s.userRep.Update(user)
	if updErr != nil {
		return errors.New(updErr.Error())
	}

	// END TRANSACTION

	return nil
}

func (s *Service) UserWithdrawalsGet(login string) ([]models.UserWithdraw, error) {
	s.log.Debug("Service.UserWithdrawalsGet() was called.", "login", login)

	user, uerr := s.userRep.GetByLogin(login)
	if uerr != nil {
		return make([]models.UserWithdraw, 0), errors.New(uerr.Error())
	}

	withdrawals, err := s.wthrwlRep.GetAllByUserID(user.ID)
	if err != nil {
		return make([]models.UserWithdraw, 0), errors.New(err.Error())
	}

	results := make([]models.UserWithdraw, len(withdrawals))
	for i, item := range withdrawals {
		results[i] = models.UserWithdraw{
			Order:       item.Order,
			Sum:         helper.ConvertPriceToFloat64(item.Sum),
			ProcessedAt: item.ProcessedAt, // .Format(time.RFC3339)
		}
	}

	return results, nil
}
