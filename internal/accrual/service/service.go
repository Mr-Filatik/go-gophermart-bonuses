package service

import (
	"errors"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/models"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
)

var (
	ErrEntityNotFound      = errors.New("entity not found")
	ErrEntityAlreadyExists = errors.New("entity already exists")
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

func (s *Service) OrderGet(number string) (models.OrderResponse, error) {
	if number == "000000" {
		return models.OrderResponse{}, ErrEntityNotFound
	}
	return models.OrderResponse{
		Order:   number,
		Status:  models.OrderStatusProcessed,
		Accrual: 12.8,
	}, nil
}

func (s *Service) OrderCreate(data models.OrderRequest) error {
	if data.Order == "000000" {
		return ErrEntityAlreadyExists
	}
	return nil
}

func (s *Service) GoodCreate(data models.GoodRequest) error {
	if data.Match == "exists" {
		return ErrEntityAlreadyExists
	}
	return nil
}
