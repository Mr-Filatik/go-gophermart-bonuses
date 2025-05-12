package worker

import (
	"strings"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/database/models"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/database/repository"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/helper"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
)

const (
	workerChannelLimit        = 10
	workerFullPercentForPrice = 100
	workerFailedChanceLimit   = 5
)

var workerFailedChance = 0

type Worker struct {
	ordRep  repository.IOrderRepository
	goodRep repository.IGoodRepository
	rulRep  repository.IRuleRepository
	log     logger.Logger
	jobs    chan uint64
}

func New(
	ordRep repository.IOrderRepository,
	goodRep repository.IGoodRepository,
	rulRep repository.IRuleRepository,
	log logger.Logger) *Worker {
	srv := Worker{
		ordRep:  ordRep,
		goodRep: goodRep,
		rulRep:  rulRep,
		log:     log,
		jobs:    nil,
	}

	log.Info("Worker created")

	return &srv
}

func (w *Worker) Run() {
	w.jobs = make(chan uint64, workerChannelLimit)

	for i := 1; i <= workerChannelLimit; i++ {
		go w.Processing(i)
	}
}

func (w *Worker) AddTask(number uint64) {
	w.jobs <- number
}

func (w *Worker) Processing(number int) {
	w.log.Info("Run worker", "worker_number", number)
	for value := range w.jobs {
		w.log.Info("Run worker processed", "number", number)

		order, oErr := w.ordRep.GetByNumber(value)
		if oErr != nil {
			continue
		}

		order.Status = models.OrderStatusProcessing
		uErr := w.ordRep.Update(order)
		if uErr != nil {
			continue
		}

		// time.Sleep(time.Second)

		goods, gErr := w.goodRep.GetAllByOrderID(order.ID)
		if gErr != nil {
			continue
		}
		if len(goods) == 0 {
			continue
		}

		rules, rErr := w.rulRep.GetAll()
		if rErr != nil {
			continue
		}
		if len(rules) == 0 {
			continue
		}

		// Failed processing
		workerFailedChance++
		if workerFailedChance >= workerFailedChanceLimit {
			workerFailedChance = 0
			order.Status = models.OrderStatusInvalid
			uErr = w.ordRep.Update(order)
			if uErr != nil {
				continue
			}
		}

		// time.Sleep(time.Second)

		accrual := uint64(0)
		for gi := range goods {
			w.log.Info("Good", "description", goods[gi].Description)
			for ri := range rules {
				w.log.Info("Rule", "match", rules[ri].Match)
				if strings.Contains(goods[gi].Description, rules[ri].Match) {
					switch rules[ri].RewardType {
					case models.RuleRewardTypeCount:
						accrual += rules[ri].Reward
					case models.RuleRewardTypePercent:
						percent := helper.ConvertPriceToFloat64(rules[ri].Reward)
						price := helper.ConvertPriceToFloat64(goods[gi].Price)
						value := price * percent / workerFullPercentForPrice
						accrual += helper.ConvertPriceToUint64(value)
					default:
						w.log.Info("Accrual", "add", "error")
						// error
					}
				}
			}
		}

		order.Accrual = accrual
		order.Status = models.OrderStatusProcessed
		uErr = w.ordRep.Update(order)
		if uErr != nil {
			continue
		}
	}
}

func (w *Worker) Close() {
	if w.jobs != nil {
		close(w.jobs)
	}
}
