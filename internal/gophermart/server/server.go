package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/models"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/service"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/server"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	router  *chi.Mux
	service *service.Service
	log     logger.Logger
}

func New(srvc *service.Service, log logger.Logger) *Server {
	srv := Server{
		router:  chi.NewRouter(),
		service: srvc,
		log:     log,
	}

	srv.registerHandlers()

	log.Info("Server created")

	return &srv
}

func (s *Server) registerHandlers() {
	s.router.Handle("/api/user/register", http.HandlerFunc(s.UserRegister))
	s.router.Handle("/api/user/login", http.HandlerFunc(s.UserLogin))
	s.router.Handle("/api/user/orders", http.HandlerFunc(s.UserOrders))
	s.router.Handle("/api/user/balance", http.HandlerFunc(s.UserBalance))
	s.router.Handle("/api/user/balance/withdraw", http.HandlerFunc(s.UserBalanceWithdraw))
	s.router.Handle("/api/user/withdrawals", http.HandlerFunc(s.UserWithdrawals))
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

func (s *Server) UserRegister(w http.ResponseWriter, r *http.Request) {
	ok := server.ValidateRequestMethod(w, s.log, r.Method, http.MethodPost)
	if !ok {
		return
	}

	data, err := server.GetDataFromBody[models.UserRegisterRequest](r)
	if err != nil {
		s.log.Error("Error get data from body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if data != nil && (data.Login == "" || data.Password == "") {
		http.Error(w, "data is empty", http.StatusBadRequest)
		return
	}

	err = s.service.UserRegister(*data)
	if err != nil {
		s.log.Error("Error register user", err)
		if errors.Is(err, service.ErrLoginAlreadyTaken) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	token, terr := server.CreateToken(data.Login)
	if terr != nil {
		http.Error(w, terr.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/api/user",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	// 200 — пользователь успешно зарегистрирован и аутентифицирован;
	// 400 — неверный формат запроса;
	// 409 — логин уже занят;
	// 500 — внутренняя ошибка сервера.
}

func (s *Server) UserLogin(w http.ResponseWriter, r *http.Request) {
	ok := server.ValidateRequestMethod(w, s.log, r.Method, http.MethodPost)
	if !ok {
		return
	}

	data, err := server.GetDataFromBody[models.UserLoginRequest](r)
	if err != nil {
		s.log.Error("Error get data from body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if data != nil && (data.Login == "" || data.Password == "") {
		http.Error(w, "data is empty", http.StatusBadRequest)
		return
	}

	lerr := s.service.UserLogin(*data)
	if lerr != nil {
		s.log.Error("Error login user", lerr)
		if errors.Is(lerr, service.ErrInvalidLoginOrPassword) {
			http.Error(w, lerr.Error(), http.StatusUnauthorized)
			return
		} else {
			http.Error(w, lerr.Error(), http.StatusInternalServerError)
			return
		}
	}

	token, terr := server.CreateToken(data.Login)
	if terr != nil {
		http.Error(w, terr.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/api/user",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	// 200 — пользователь успешно аутентифицирован;
	// 400 — неверный формат запроса;
	// 401 — неверная пара логин/пароль;
	// 500 — внутренняя ошибка сервера.
}

func (s *Server) UserOrders(w http.ResponseWriter, r *http.Request) {
	ok := server.ValidateRequestMethods(w, s.log, r.Method, http.MethodPost, http.MethodGet)
	if !ok {
		return
	}

	s.log.Info("tun tun tun tun tun saur")

	// POST
	// 200 — номер заказа уже был загружен этим пользователем;
	// 202 — новый номер заказа принят в обработку;
	// 400 — неверный формат запроса;
	// 401 — пользователь не аутентифицирован;
	// 409 — номер заказа уже был загружен другим пользователем;
	// 422 — неверный формат номера заказа;
	// 500 — внутренняя ошибка сервера.

	// GET
	// 200 — успешная обработка запроса.
	// 204 — нет данных для ответа.
	// 401 — пользователь не авторизован.
	// 500 — внутренняя ошибка сервера.
}

func (s *Server) UserBalance(w http.ResponseWriter, r *http.Request) {
	ok := server.ValidateRequestMethod(w, s.log, r.Method, http.MethodGet)
	if !ok {
		return
	}

	s.log.Info("tun tun tun tun tun saur")

	// 200 — успешная обработка запроса.
	// 401 — пользователь не авторизован.
	// 500 — внутренняя ошибка сервера.
}

func (s *Server) UserBalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	ok := server.ValidateRequestMethod(w, s.log, r.Method, http.MethodPost)
	if !ok {
		return
	}

	s.log.Info("tun tun tun tun tun saur")

	// 200 — успешная обработка запроса;
	// 401 — пользователь не авторизован;
	// 402 — на счету недостаточно средств;
	// 422 — неверный номер заказа;
	// 500 — внутренняя ошибка сервера.
}

func (s *Server) UserWithdrawals(w http.ResponseWriter, r *http.Request) {
	ok := server.ValidateRequestMethod(w, s.log, r.Method, http.MethodGet)
	if !ok {
		return
	}

	s.log.Info("tun tun tun tun tun saur")

	// 200 — успешная обработка запроса;
	// 204 — нет ни одного списания.
	// 401 — пользователь не авторизован.
	// 500 — внутренняя ошибка сервера.
}
