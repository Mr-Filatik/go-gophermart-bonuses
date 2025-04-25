package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
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

func GetDataFromBody[Tout any](r *http.Request) (Tout, error) {
	var metr Tout
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		return metr, ErrInvalidRequestFormat
	}

	if err := json.Unmarshal(buf.Bytes(), &metr); err != nil {
		return metr, ErrInvalidRequestFormat
	}

	return metr, nil
}
