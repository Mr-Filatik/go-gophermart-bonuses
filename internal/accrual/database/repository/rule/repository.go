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

type RuleRepository struct {
	connector *connector.Connector
	log       logger.Logger
}

func New(conn *connector.Connector, log logger.Logger) *RuleRepository {
	return &RuleRepository{
		connector: conn,
		log:       log,
	}
}

func (r *RuleRepository) Create(entity *models.Rule) error {
	r.log.Debug("RuleRepository.Create() was called.", "match", entity.Match)

	conn := r.connector.GetDB()

	result := conn.Create(entity)
	if result.Error != nil {
		r.log.Error("OrderRepository.Create error.", result.Error, "match", entity.Match)
		if strings.Contains(result.Error.Error(), "ERROR: duplicate key value violates unique constraint") {
			return repository.ErrEntityAlreadyExists
		}
		return errors.New(result.Error.Error())
	}

	return nil
}

func (r *RuleRepository) GetAll() ([]models.Rule, error) {
	r.log.Debug("RuleRepository.GetAll() was called.")

	conn := r.connector.GetDB()

	var users []models.Rule
	result := conn.Find(&users)
	if result.Error != nil {
		r.log.Error("OrderRepository.GetAllByUserID error.", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repository.ErrEntityNotFound
		}
		return nil, errors.New(result.Error.Error())
	}

	return users, nil
}
