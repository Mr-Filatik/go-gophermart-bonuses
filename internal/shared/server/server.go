package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
	"github.com/dgrijalva/jwt-go"
)

func ValidateRequestMethod(w http.ResponseWriter, log logger.Logger, current string, allowed string) bool {
	if current == allowed {
		return true
	}

	log.Error(
		"Invalid request method",
		errors.New("invalid request method"),
		"actual", current,
		"allowed", allowed,
	)

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	return false
}

func ValidateRequestMethods(w http.ResponseWriter, log logger.Logger, current string, allowed ...string) bool {
	for _, method := range allowed {
		if current == method {
			return true
		}
	}

	log.Error(
		"Invalid request method",
		errors.New("invalid request method"),
		"actual", current,
		"allowed", allowed,
	)

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	return false
}

var (
	ErrInvalidRequestFormat = errors.New("invalid request format")
	ErrInternalServerError  = errors.New("internal server error")
)

func GetDataFromBody[Tout any](r *http.Request) (*Tout, error) {
	data := new(Tout)
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		return data, ErrInvalidRequestFormat
	}

	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		return data, ErrInvalidRequestFormat
	}

	return data, nil
}

func CreateToken(login string) (string, error) {
	secretKey := []byte("FILATIK_SECRET_KEY_FOR_TOKEN")

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Issuer:    "gophermart-bonuces.gophermart",
		Subject:   login,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}
