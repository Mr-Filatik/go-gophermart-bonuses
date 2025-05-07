package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/database/models"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/database/repository"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/database/connector"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
	"gorm.io/gorm"
)

type WithdrawalRepository struct {
	connector *connector.Connector
	log       logger.Logger
}

func New(conn *connector.Connector, log logger.Logger) *WithdrawalRepository {
	return &WithdrawalRepository{
		connector: conn,
		log:       log,
	}
}

func (r *WithdrawalRepository) Create(entity *models.UserWithdrawal) error {
	r.log.Debug(
		"WithdrawalRepository.Create() was called.",
		"order", entity.Order,
		"sum", entity.Sum,
	)

	entity.ProcessedAt = time.Now()

	conn := r.connector.GetDB()

	result := conn.Create(entity)
	if result.Error != nil {
		r.log.Error("WithdrawalRepository.Create error.", result.Error, "order", entity.Order)
		if strings.Contains(result.Error.Error(), "ERROR: duplicate key value violates unique constraint") {
			return repository.ErrEntityAlreadyExists
		}
		return errors.New(result.Error.Error())
	}

	return nil
}

func (r *WithdrawalRepository) GetAllByUserID(userID uint64) ([]models.UserWithdrawal, error) {
	r.log.Debug("WithdrawalRepository.GetAllByUserID() was called.", "user_id", userID)

	conn := r.connector.GetDB()

	var datas []models.UserWithdrawal
	result := conn.Where("user_id = ?", userID).Find(&datas)
	if result.Error != nil {
		r.log.Error("WithdrawalRepository.GetAllByUserLogin error.", result.Error, "user_id", userID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repository.ErrEntityNotFound
		}
		return nil, errors.New(result.Error.Error())
	}

	return datas, nil
}
