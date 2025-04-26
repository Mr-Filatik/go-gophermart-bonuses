package main

import (
	"errors"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/config"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/server"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/service"
	logger "github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger/zap/sugar"
)

func main() {
	log := logger.New(logger.LevelDebug)
	defer log.Close()

	log.Debug("Test")
	log.Info("Test")
	log.Warning("Test")
	log.Error("Test", errors.New("test error"))

	conf := config.Initialize()

	srvc := service.New(log)

	serv := server.New(srvc, log)
	serv.Start(conf.RunAddress)
}
