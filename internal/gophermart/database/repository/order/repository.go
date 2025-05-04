package repository

import (
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/database/connector"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
)

type OrderRepository struct {
	connector *connector.Connector
	log       logger.Logger
}

func (r *OrderRepository) Create() {

}
