package repository

import (
	"errors"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/database/models"
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

type IUserRepository interface {
	IRepository[models.User]
	GetByLogin(login string) (*models.User, error)
}

type IOrderRepository interface {
	IRepository[models.UserOrder]
	GetByNumber(number uint64) (*models.UserOrder, error)
	GetAllByUserID(userID uint64) ([]models.UserOrder, error)
}
