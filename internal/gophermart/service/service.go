package service

import (
	"errors"
	"time"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/models"
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
	log logger.Logger
}

func New(log logger.Logger) *Service {
	srv := Service{
		log: log,
	}

	log.Info("Service created")

	return &srv
}

func (s *Service) UserRegister(data models.UserRegisterRequest) error {
	if data.Login == "login" {
		return ErrLoginAlreadyTaken
	}
	if data.Password != "password" {
		return errors.New("tun tun tun tun tun saur")
	}

	s.log.Info(
		"User register",
	)
	return nil
}

func (s *Service) UserLogin(data models.UserLoginRequest) error {
	if data.Login == "error" && data.Password == "error" {
		return errors.New("tun tun tun tun tun saur")
	}

	if data.Password != "password" {
		return ErrInvalidLoginOrPassword
	}

	s.log.Info(
		"User login",
	)
	return nil
}

func (s *Service) UserOrdersGet(login string) ([]models.UserOrder, error) {
	if login == "empty" {
		return make([]models.UserOrder, 0), nil
	}

	s.log.Info(
		"Get orders",
	)

	parsedTime, err := time.Parse(time.RFC3339, "2023-10-01T12:00:00Z")
	if err != nil {
		return make([]models.UserOrder, 0), errors.New(err.Error())
	}

	orders := []models.UserOrder{
		{Number: "12345", Status: models.UserOrderStatusNew, UploadedAt: parsedTime},
		{Number: "67890", Status: models.UserOrderStatusProcessed, Accrual: 50.25, UploadedAt: parsedTime},
	}
	return orders, nil
}

func (s *Service) UserOrdersCreate(login string, number string) error {
	if login == "login" && number == "000001" {
		return ErrAlreadyUploadThisUser
	}
	if login == "login" && number == "000002" {
		return ErrAlreadyUploadOtherUser
	}
	if number == "000000" {
		return nil
	}
	return errors.New("tun tun tun tun saur")
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
