package service

import (
	"errors"
	"strconv"
	"time"

	dbModels "github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/database/models"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/database/repository"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/models"
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
	userRep  repository.IUserRepository
	orderRep repository.IOrderRepository
	log      logger.Logger
}

func New(
	userRep repository.IUserRepository,
	orderRep repository.IOrderRepository,
	log logger.Logger) *Service {
	srv := Service{
		userRep:  userRep,
		orderRep: orderRep,
		log:      log,
	}

	log.Info("Service created")

	return &srv
}

func (s *Service) UserRegister(data models.UserRegisterRequest) error {
	s.log.Debug("Service.UserRegister() was called.", "login", data.Login)

	hashedPassword, passErr := helper.GeneratePasswordHash(data.Password)
	if passErr != nil {
		return errors.New(passErr.Error())
	}

	user := dbModels.User{
		Login:        data.Login,
		PasswordHash: string(hashedPassword),
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
	s.log.Debug("Service.UserLogin() was called.", "login", data.Login)

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
	s.log.Debug("Service.UserOrdersGet() was called.", "login", login)

	user, err := s.userRep.GetByLogin(login)
	if err != nil {
		return make([]models.UserOrder, 0), errors.New(err.Error())
	}

	orders, oerr := s.orderRep.GetAllByUserID(user.ID)
	if oerr != nil {
		return make([]models.UserOrder, 0), errors.New(oerr.Error())
	}

	userOrders := make([]models.UserOrder, len(orders))
	for i, order := range orders {
		userOrders[i] = models.UserOrder{
			Number:     strconv.FormatUint(order.Number, 10),
			Status:     models.UserOrderStatus(order.Status),
			Accrual:    0,
			UploadedAt: order.UploadedAt, // .Format(time.RFC3339)
		}
	}

	return userOrders, nil
}

func (s *Service) UserOrdersCreate(login string, number string) error {
	s.log.Debug(
		"Service.UserOrdersCreate() was called.",
		"login", login,
		"number", number,
	)

	num, err := strconv.ParseUint(number, 10, 64)
	if err != nil {
		return errors.New(err.Error())
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
			return nil
		}
		return errors.New(err.Error())
	}

	if order.User.Login == login {
		return ErrAlreadyUploadThisUser
	} else {
		return ErrAlreadyUploadOtherUser
	}
}

func (s *Service) UserBalanceGet(login string) (models.UserBalanceResponse, error) {
	s.log.Info(
		"User balance",
	)
	return models.UserBalanceResponse{
		Current:   12.5,
		Withdrawn: 6.5,
	}, nil
}

func (s *Service) UserBalanceWithdraw(data models.UserBalanceWithdrawRequest) error {
	if data.Order == "000000" {
		return ErrInvalidOrderNumber
	}
	if data.Sum > 50.0 {
		return ErrInsufficientFunds
	}

	s.log.Info(
		"Balance withdraw",
	)
	return nil
}

func (s *Service) UserWithdrawalsGet(login string) ([]models.UserWithdraw, error) {
	if login == "empty" {
		return make([]models.UserWithdraw, 0), nil
	}

	s.log.Info(
		"Get withdrawals",
	)

	parsedTime, err := time.Parse(time.RFC3339, "2023-10-01T12:00:00Z")
	if err != nil {
		return make([]models.UserWithdraw, 0), errors.New(err.Error())
	}

	withdrawals := []models.UserWithdraw{
		{Number: "12345", Sum: 100.50, ProcessedAt: parsedTime},
		{Number: "67890", Sum: 50.25, ProcessedAt: parsedTime},
	}
	return withdrawals, nil
}
