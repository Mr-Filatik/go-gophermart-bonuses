package server

import (
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/middleware"
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
	s.router.Handle(
		"/api/user/orders",
		middleware.AuthMiddleware(http.HandlerFunc(s.UserOrders)))
	s.router.Handle(
		"/api/user/balance",
		middleware.AuthMiddleware(http.HandlerFunc(s.UserBalance)))
	s.router.Handle(
		"/api/user/balance/withdraw",
		middleware.AuthMiddleware(http.HandlerFunc(s.UserBalanceWithdraw)))
	s.router.Handle(
		"/api/user/withdrawals",
		middleware.AuthMiddleware(http.HandlerFunc(s.UserWithdrawals)))
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

	data, err := server.GetDataFromBodyInJSON[models.UserRegisterRequest](r)
	if err != nil {
		s.log.Error("Error get data from body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if data != nil && (data.Login == "" || data.Password == "") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.service.UserRegister(*data)
	if err != nil {
		s.log.Error("Error register user", err)
		if errors.Is(err, service.ErrLoginAlreadyTaken) {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	token, terr := server.CreateToken(data.Login)
	if terr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/api/user",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}

func (s *Server) UserLogin(w http.ResponseWriter, r *http.Request) {
	ok := server.ValidateRequestMethod(w, s.log, r.Method, http.MethodPost)
	if !ok {
		return
	}

	data, err := server.GetDataFromBodyInJSON[models.UserLoginRequest](r)
	if err != nil {
		s.log.Error("Error get data from body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if data != nil && (data.Login == "" || data.Password == "") {
		s.log.Error("Data from body is empty", errors.New("body is empty"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lerr := s.service.UserLogin(*data)
	if lerr != nil {
		s.log.Error("Error login user", lerr)
		if errors.Is(lerr, service.ErrInvalidLoginOrPassword) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	token, terr := server.CreateToken(data.Login)
	if terr != nil {
		s.log.Error("Create token error", terr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/api/user",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}

func (s *Server) UserOrders(w http.ResponseWriter, r *http.Request) {
	ok := server.ValidateRequestMethod(w, s.log, r.Method, http.MethodPost, http.MethodGet)
	if !ok {
		return
	}

	ctx := r.Context()
	login, ok := ctx.Value("login").(string)
	if !ok || login == "" {
		http.Error(w, "Unauthorized: Login not found in context", http.StatusUnauthorized)
		s.log.Error("Login not found in context", errors.New("not login in token"))
		return
	}

	if r.Method == http.MethodGet {
		orders, err := s.service.UserOrdersGet(login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(orders) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		serr := server.SetDataToBodyInJSON(w, orders)
		if serr != nil {
			http.Error(w, serr.Error(), http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {
		data, err := server.GetStringFromBody(r)
		if err != nil {
			s.log.Error("Error get data from body", err)
			http.Error(w, err.Error(), http.StatusBadRequest) // dont used
			return
		}

		s.log.Info("Data", "data", data)
		ok, _ := regexp.MatchString(`^\d+$`, data)
		if !ok {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		cerr := s.service.UserOrdersCreate(login, data)
		if cerr != nil {
			s.log.Error("Error create order", cerr)
			if errors.Is(cerr, service.ErrAlreadyUploadThisUser) {
				w.WriteHeader(http.StatusOK)
				return
			}
			if errors.Is(cerr, service.ErrAlreadyUploadOtherUser) {
				http.Error(w, cerr.Error(), http.StatusConflict)
				return
			}
			http.Error(w, cerr.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *Server) UserBalance(w http.ResponseWriter, r *http.Request) {
	ok := server.ValidateRequestMethod(w, s.log, r.Method, http.MethodGet)
	if !ok {
		return
	}

	ctx := r.Context()
	login, ok := ctx.Value("login").(string)
	if !ok || login == "" {
		http.Error(w, "Unauthorized: Login not found in context", http.StatusUnauthorized)
		s.log.Error("Login not found in context", errors.New("not login in token"))
		return
	}

	balance, err := s.service.UserBalanceGet(login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serr := server.SetDataToBodyInJSON(w, balance)
	if serr != nil {
		http.Error(w, serr.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) UserBalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	ok := server.ValidateRequestMethod(w, s.log, r.Method, http.MethodPost)
	if !ok {
		return
	}

	data, err := server.GetDataFromBodyInJSON[models.UserBalanceWithdrawRequest](r)
	if err != nil {
		s.log.Error("Error get data from body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if data != nil && (data.Order == "" || data.Sum <= 0.0) {
		http.Error(w, "data is empty", http.StatusBadRequest)
		return
	}

	berr := s.service.UserBalanceWithdraw(*data)
	if berr != nil {
		s.log.Error("Error balance withdraw", berr)
		if errors.Is(berr, service.ErrInvalidOrderNumber) {
			http.Error(w, berr.Error(), http.StatusUnprocessableEntity)
			return
		} else if errors.Is(berr, service.ErrInsufficientFunds) {
			http.Error(w, berr.Error(), http.StatusPaymentRequired)
			return
		} else {
			http.Error(w, berr.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) UserWithdrawals(w http.ResponseWriter, r *http.Request) {
	ok := server.ValidateRequestMethod(w, s.log, r.Method, http.MethodGet)
	if !ok {
		return
	}

	ctx := r.Context()
	login, ok := ctx.Value("login").(string)
	if !ok || login == "" {
		http.Error(w, "Unauthorized: Login not found in context", http.StatusUnauthorized)
		s.log.Error("Login not found in context", errors.New("not login in token"))
		return
	}

	withdrawals, err := s.service.UserWithdrawalsGet(login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	serr := server.SetDataToBodyInJSON(w, withdrawals)
	if serr != nil {
		http.Error(w, serr.Error(), http.StatusInternalServerError)
		return
	}
}
