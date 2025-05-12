package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateRequestMethod(w http.ResponseWriter, log logger.Logger, current string, allowed ...string) bool {
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

	errorMessage := fmt.Sprintf("Invalid request method. Expected one of: %v.", allowed)
	http.Error(w, errorMessage, http.StatusMethodNotAllowed)
	return false
}

var (
	ErrInvalidRequestFormat = errors.New("invalid request format")
	ErrInternalServerError  = errors.New("internal server error")
)

func GetStringFromBody(r *http.Request) (string, error) {
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		return "", ErrInvalidRequestFormat
	}

	return buf.String(), nil
}

func GetDataFromBodyInJSON[Tout any](r *http.Request) (*Tout, error) {
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

func SetDataToBodyInJSON(w http.ResponseWriter, data any) error {
	jdata, err := json.Marshal(data)
	if err != nil {
		return errors.New(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jdata)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

const TokenExpiredHours = 24

func CreateToken(login string) (string, error) {
	secretKey := []byte("FILATIK_SECRET_KEY_FOR_TOKEN")

	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * TokenExpiredHours).Unix(),
		"iss": "gophermart-bonuces.gophermart",
		"sub": login,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	strToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", errors.New(err.Error())
	}
	return strToken, nil
}

type ContextKey string

const ContextKeyUserLogin ContextKey = ContextKey("login")

func GetStringFromContext(ctx context.Context, key ContextKey) (string, bool) {
	value, ok := ctx.Value(key).(string)
	return value, ok
}

func SetStringToContext(ctx context.Context, key ContextKey, value string) context.Context {
	return context.WithValue(ctx, key, value)
}
