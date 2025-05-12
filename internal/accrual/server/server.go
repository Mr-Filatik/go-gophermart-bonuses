package server

import (
	"errors"
	"net/http"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/server/models"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/service"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/server"
	"github.com/go-chi/chi/v5"
	"golang.org/x/time/rate"
)

const (
	ServerLimitPerSecond rate.Limit = 2
	ServerLimitBuffer    int        = 4
)

type Server struct {
	router  *chi.Mux
	service *service.Service
	limiter *rate.Limiter
	log     logger.Logger
}

func New(srvc *service.Service, log logger.Logger) *Server {
	srv := Server{
		router:  chi.NewRouter(),
		service: srvc,
		limiter: rate.NewLimiter(ServerLimitPerSecond, ServerLimitBuffer),
		log:     log,
	}

	srv.registerHandlers()

	log.Info("Server created")

	return &srv
}

func (s *Server) registerHandlers() {
	s.router.Handle("/api/orders/{number}", http.HandlerFunc(s.OrderNumber))
	s.router.Handle("/api/orders", http.HandlerFunc(s.Orders))
	s.router.Handle("/api/goods", http.HandlerFunc(s.Goods))
}

func (s *Server) Start(addr string) {
	s.log.Info(
		"Server starting",
		"run address", addr,
	)
	err := http.ListenAndServe(addr, s.router)
	if err != nil {
		s.log.Error("Server work error", err)
	}
}

func (s *Server) OrderNumber(w http.ResponseWriter, r *http.Request) {
	ok := server.ValidateRequestMethod(w, s.log, r.Method, http.MethodGet)
	if !ok {
		return
	}

	if !s.limiter.Allow() {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	number := r.PathValue("number")

	order, err := s.service.OrderGet(number)
	if err != nil {
		if errors.Is(err, service.ErrEntityNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	serr := server.SetDataToBodyInJSON(w, order)
	if serr != nil {
		http.Error(w, serr.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) Orders(w http.ResponseWriter, r *http.Request) {
	ok := server.ValidateRequestMethod(w, s.log, r.Method, http.MethodPost)
	if !ok {
		return
	}

	data, err := server.GetDataFromBodyInJSON[models.OrderRequest](r)
	if err != nil {
		s.log.Error("Error get data from body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if data != nil && (data.Order == "" || len(data.Goods) == 0) {
		http.Error(w, "data is empty", http.StatusBadRequest)
		return
	}

	cerr := s.service.OrderCreate(*data)
	if cerr != nil {
		if errors.Is(cerr, service.ErrEntityAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) Goods(w http.ResponseWriter, r *http.Request) {
	ok := server.ValidateRequestMethod(w, s.log, r.Method, http.MethodPost)
	if !ok {
		return
	}

	data, err := server.GetDataFromBodyInJSON[models.GoodRequest](r)
	if err != nil {
		s.log.Error("Error get data from body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if data != nil && (data.Match == "" ||
		(data.RewardType != models.RewardTypePercent && data.RewardType != models.RewardTypeCount) || data.Reward <= 0.0) {
		http.Error(w, "data is empty", http.StatusBadRequest)
		return
	}

	cerr := s.service.GoodCreate(*data)
	if cerr != nil {
		if errors.Is(cerr, service.ErrEntityAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
