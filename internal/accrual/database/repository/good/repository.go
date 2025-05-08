package repository

import (
	"errors"
	"strings"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/database/models"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/database/repository"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/database/connector"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
	"gorm.io/gorm"
)

type GoodRepository struct {
	connector *connector.Connector
	log       logger.Logger
}

func New(conn *connector.Connector, log logger.Logger) *GoodRepository {
	return &GoodRepository{
		connector: conn,
		log:       log,
	}
}

func (r *GoodRepository) Create(entity *models.Good) error {
	r.log.Debug("GoodRepository.Create() was called.")

	conn := r.connector.GetDB()

	result := conn.Create(entity)
	if result.Error != nil {
		r.log.Error("GoodRepository.Create error.", result.Error)
		if strings.Contains(result.Error.Error(), "ERROR: duplicate key value violates unique constraint") {
			return repository.ErrEntityAlreadyExists
		}
		return errors.New(result.Error.Error())
	}

	return nil
}

func (r *GoodRepository) GetAllByOrderID(orderID uint64) ([]models.Good, error) {
	r.log.Debug("OrderRepository.GetByNumber() was called.", "order_id", orderID)

	conn := r.connector.GetDB()

	var users []models.Good
	result := conn.Where("order_id = ?", orderID).Find(&users)
	if result.Error != nil {
		r.log.Error("OrderRepository.GetAllByUserID error.", result.Error, "order_id", orderID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repository.ErrEntityNotFound
		}
		return nil, errors.New(result.Error.Error())
	}

	return users, nil
}
