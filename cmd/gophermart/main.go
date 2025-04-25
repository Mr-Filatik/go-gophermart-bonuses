package main

import (
	"errors"

	logger "github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger/zap/sugar"
)

func main() {
	log := logger.New(logger.LevelDebug)
	defer log.Close()

	log.Debug("Test")
	log.Info("Test")
	log.Warning("Test")
	log.Error("Test", errors.New("test error"))
}
