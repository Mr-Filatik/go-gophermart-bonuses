package worker

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	dbModels "github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/database/models"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/database/repository"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/helper"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/server/models"
	"github.com/go-resty/resty/v2"
)

const (
	workerChannelLimit = 10
)

type Worker struct {
	ordRep   repository.IOrderRepository
	useRep   repository.IUserRepository
	log      logger.Logger
	jobs     chan uint64
	endpoint string
}

func New(
	ordRep repository.IOrderRepository,
	useRep repository.IUserRepository,
	endpoint string,
	log logger.Logger) *Worker {
	srv := Worker{
		ordRep:   ordRep,
		useRep:   useRep,
		endpoint: endpoint,
		log:      log,
		jobs:     nil,
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

		// request to accrual
		var address string
		if strings.Contains(w.endpoint, "http") {
			address = w.endpoint + "/api/orders/" + strconv.FormatUint(value, 10)
		} else {
			address = "http://" + w.endpoint + "/api/orders/" + strconv.FormatUint(value, 10)
		}

		client := resty.New()
		resp, rerr := client.R().
			// SetHeader("Content-Type", "application/json").
			// SetHeader(ContentEncodingHeader, EncodingType).
			// SetHeader(AcceptEncodingHeader, EncodingType).
			// SetBody(dat).
			Get(address)

		if rerr != nil {
			// log.Error("Response error", rerr)
			w.log.Error("Responce to accrual error.", rerr)
			continue
		}

		if resp.StatusCode() != http.StatusOK {
			w.log.Info("Request error.", errors.New("status not ok"))
			continue
		}

		data := resp.Body()
		w.log.Info("Result", "body", data)

		var resporder models.OrderResponse
		// Декодируем JSON в структуру
		err := json.Unmarshal(data, &resporder)
		if err != nil {
			w.log.Error("Unmarshal error.", err)
			continue
		}

		// if resporder.Status != OrderStatusRegistred {
		// 	continue
		// }

		order, oErr := w.ordRep.GetByNumber(value)
		if oErr != nil {
			w.log.Error("Get by number error.", oErr)
			continue
		}

		switch resporder.Status {
		case models.OrderStatusProcessing:
			order.Status = dbModels.UserOrderStatusProcessing
		case models.OrderStatusProcessed:
			order.Status = dbModels.UserOrderStatusProcessed
		case models.OrderStatusInvalidProcessing:
			order.Status = dbModels.UserOrderStatusInvalid
		default:
			continue
		}
		order.Accrual = helper.ConvertPriceToUint64(resporder.Accrual)

		uErr := w.ordRep.Update(order)
		if uErr != nil {
			w.log.Error("Update number error.", uErr)
			continue
		}

		if order.Status == dbModels.UserOrderStatusProcessed {
			user, usErr := w.useRep.GetByID(order.UserID)
			if usErr != nil {
				w.log.Error("Get user info error.", usErr)
				continue
			}
			user.Current += order.Accrual
			eeerr := w.useRep.Update(user)
			if eeerr != nil {
				w.log.Error("Get user info error.", eeerr)
				continue
			}
		}
	}
}

func (w *Worker) Close() {
	if w.jobs != nil {
		close(w.jobs)
	}
}
