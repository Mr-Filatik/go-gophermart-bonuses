package repository

import (
	"errors"
	"strings"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/database/models"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/database/repository"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/database/connector"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
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
