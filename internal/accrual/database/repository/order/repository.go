package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/database/models"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/database/repository"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/database/connector"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
	"gorm.io/gorm"
)

type OrderRepository struct {
	connector *connector.Connector
	log       logger.Logger
}

func New(conn *connector.Connector, log logger.Logger) *OrderRepository {
	return &OrderRepository{
		connector: conn,
		log:       log,
	}
}

func (r *OrderRepository) Create(entity *models.Order) error {
	r.log.Debug("OrderRepository.Create() was called.", "order_number", entity.Number)

	entity.UploadedAt = time.Now()
	entity.Accrual = 0
	entity.Status = models.OrderStatusRegistered

	conn := r.connector.GetDB()

	result := conn.Create(entity)
	if result.Error != nil {
		r.log.Error("OrderRepository.Create error.", result.Error, "order_number", entity.Number)
		if strings.Contains(result.Error.Error(), "ERROR: duplicate key value violates unique constraint") {
			return repository.ErrEntityAlreadyExists
		}
		return errors.New(result.Error.Error())
	}

	return nil
}

func (r *OrderRepository) GetByNumber(number uint64) (*models.Order, error) {
	r.log.Debug("OrderRepository.GetByNumber() was called.", "number", number)

	conn := r.connector.GetDB()

	var user models.Order
	result := conn.First(&user, "number = ?", number)
	if result.Error != nil {
		r.log.Error("OrderRepository.GetByNumber error.", result.Error, "number", number)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repository.ErrEntityNotFound
		}
		return nil, errors.New(result.Error.Error())
	}

	return &user, nil
}

func (r *OrderRepository) Update(item *models.Order) error {
	r.log.Debug(
		"OrderRepository.Update() was called.",
		"order_number", item.Number,
		"accrual", item.Accrual,
		"status", item.Status,
	)

	conn := r.connector.GetDB()

	result := conn.Save(&item)
	if result.Error != nil {
		r.log.Error("OrderRepository.Update error.", result.Error, "number", item.Number)
		return errors.New(result.Error.Error())
	}

	return nil
}
