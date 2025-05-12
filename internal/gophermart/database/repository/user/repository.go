package repository

import (
	"errors"
	"strings"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/database/models"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/database/repository"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/database/connector"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
	"gorm.io/gorm"
)

type UserRepository struct {
	connector *connector.Connector
	log       logger.Logger
}

func New(conn *connector.Connector, log logger.Logger) *UserRepository {
	return &UserRepository{
		connector: conn,
		log:       log,
	}
}

func (r *UserRepository) Create(entity *models.User) error {
	r.log.Debug("UserRepository.Create() was called.", "user_login", entity.Login)

	entity.Current = 0
	entity.Withdrawn = 0

	conn := r.connector.GetDB()

	result := conn.Create(entity)
	if result.Error != nil {
		r.log.Error("UserRepository.Create error.", result.Error, "user_login", entity.Login)
		if strings.Contains(result.Error.Error(), "ERROR: duplicate key value violates unique constraint") {
			return repository.ErrEntityAlreadyExists
		}
		return errors.New(result.Error.Error())
	}

	return nil
}

func (r *UserRepository) GetByID(userID uint64) (*models.User, error) {
	r.log.Debug("UserRepository.GetByLogin() was called.", "userID", userID)

	return r.getByField("id = ?", userID)
}

func (r *UserRepository) GetByLogin(login string) (*models.User, error) {
	r.log.Debug("UserRepository.GetByLogin() was called.", "login", login)

	return r.getByField("login = ?", login)
}

func (r *UserRepository) getByField(condition string, value any) (*models.User, error) {
	conn := r.connector.GetDB()

	var user models.User
	result := conn.First(&user, condition, value)
	if result.Error != nil {
		r.log.Error("UserRepository.GetByLogin error.", result.Error, condition, value)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repository.ErrEntityNotFound
		}
		return nil, errors.New(result.Error.Error())
	}

	return &user, nil
}

func (r *UserRepository) Update(item *models.User) error {
	r.log.Debug(
		"UserRepository.Update() was called.",
		"user_login", item.Login,
		"current", item.Current,
		"withdrawn", item.Withdrawn,
	)

	conn := r.connector.GetDB()

	result := conn.Save(&item)
	if result.Error != nil {
		r.log.Error("UserRepository.Update error.", result.Error, "login", item.Login)
		return errors.New(result.Error.Error())
	}

	return nil
}
