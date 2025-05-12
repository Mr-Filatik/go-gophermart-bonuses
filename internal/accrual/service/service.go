package service

import (
	"errors"
	"strconv"

	dbModels "github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/database/models"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/database/repository"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/server/models"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/helper"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/worker"
)

var (
	ErrEntityNotFound      = errors.New("entity not found")
	ErrEntityAlreadyExists = errors.New("entity already exists")
)

type Service struct {
	ordRep  repository.IOrderRepository
	goodRep repository.IGoodRepository
	rulRep  repository.IRuleRepository
	log     logger.Logger
	worker  worker.IWorker[uint64]
}

func New(
	ordRep repository.IOrderRepository,
	goodRep repository.IGoodRepository,
	rulRep repository.IRuleRepository,
	worker worker.IWorker[uint64],
	log logger.Logger) *Service {
	srv := Service{
		ordRep:  ordRep,
		goodRep: goodRep,
		rulRep:  rulRep,
		log:     log,
		worker:  worker,
	}

	log.Info("Service created")

	return &srv
}

func (s *Service) OrderGet(number string) (models.OrderResponse, error) {
	s.log.Debug("Service.OrderCreate() was called.", "number", number)

	num, err := strconv.ParseUint(number, 10, 64)
	if err != nil {
		return models.OrderResponse{}, errors.New(err.Error())
	}

	order, gerr := s.ordRep.GetByNumber(num)
	if gerr != nil {
		if errors.Is(gerr, repository.ErrEntityNotFound) {
			return models.OrderResponse{}, ErrEntityNotFound
		}
		return models.OrderResponse{}, errors.New(gerr.Error())
	}

	return models.OrderResponse{
		Order:   strconv.FormatUint(order.Number, 10),
		Status:  models.OrderStatus(order.Status),
		Accrual: helper.ConvertPriceToFloat64(order.Accrual),
	}, nil
}

func (s *Service) OrderCreate(data models.OrderRequest) error {
	s.log.Debug(
		"Service.OrderCreate() was called.",
		"order", data.Order,
		"goods_count", len(data.Goods),
	)

	num, err := strconv.ParseUint(data.Order, 10, 64)
	if err != nil {
		return errors.New(err.Error())
	}

	// BEGIN TRANSACTION

	createdOrder := dbModels.Order{
		Number: num,
	}
	cerr := s.ordRep.Create(&createdOrder)
	if cerr != nil {
		if errors.Is(cerr, repository.ErrEntityAlreadyExists) {
			return ErrEntityAlreadyExists
		}
		return errors.New(cerr.Error())
	}

	// ADD RETURN ID FROM CREATE
	order, gerr := s.ordRep.GetByNumber(num)
	if gerr != nil {
		return errors.New(gerr.Error())
	}

	for _, item := range data.Goods {
		createdGood := dbModels.Good{
			Description: item.Description,
			Price:       helper.ConvertPriceToUint64(item.Price),
			OrderID:     order.ID,
			Order:       *order,
		}
		err := s.goodRep.Create(&createdGood)
		if err != nil {
			return errors.New(err.Error())
		}
	}

	// END TRANSACTION

	s.worker.AddTask(createdOrder.Number)

	return nil
}

func (s *Service) GoodCreate(data models.GoodRequest) error {
	s.log.Debug(
		"Service.GoodCreate() was called.",
		"match", data.Match,
		"reward", data.Reward,
		"reward_type", data.RewardType,
	)

	createdRule := dbModels.Rule{
		Match:      data.Match,
		Reward:     helper.ConvertPriceToUint64(data.Reward),
		RewardType: dbModels.RuleRewardType(data.RewardType),
	}
	cerr := s.rulRep.Create(&createdRule)
	if cerr != nil {
		if errors.Is(cerr, repository.ErrEntityAlreadyExists) {
			return ErrEntityAlreadyExists
		}
		return errors.New(cerr.Error())
	}

	return nil
}
