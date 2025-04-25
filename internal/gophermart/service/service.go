package service

import (
	"errors"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/models"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
)

var (
	ErrLoginAlreadyTaken = errors.New("login already taken")
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
		return errors.New("aaa bb ccc")
	}

	s.log.Info(
		"User register",
	)
	return nil
}

func (s *Service) UserLogin(data models.UserLoginRequest) (string, error) {
	s.log.Info(
		"User login",
	)
	return "token", nil
}
