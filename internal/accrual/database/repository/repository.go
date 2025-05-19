package repository

import (
	"errors"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/database/models"
)

var (
	ErrEntityNotFound      = errors.New("entity not found")
	ErrEntityAlreadyExists = errors.New("entity already exists")
)

type IRepository[T any] interface {
	Create(item *T) error
	// GetByID(entity_id uint64) (*T, error)
	// Update()
	// Delete()
}

type IOrderRepository interface {
	IRepository[models.Order]
	GetByNumber(number uint64) (*models.Order, error)
	Update(entity *models.Order) error
}

type IGoodRepository interface {
	IRepository[models.Good]
	GetAllByOrderID(orderID uint64) ([]models.Good, error)
}

type IRuleRepository interface {
	IRepository[models.Rule]
	GetAll() ([]models.Rule, error)
}
