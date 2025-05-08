package main

import (
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/config"
	goodRepository "github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/database/repository/good"
	orderRepository "github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/database/repository/order"
	ruleRepository "github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/database/repository/rule"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/server"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/accrual/service"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/database/connector"
	logger "github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger/zap/sugar"
)

func main() {
	log := logger.New(logger.LevelDebug)
	defer log.Close()

	conf := config.Initialize()

	conn := connector.New(log)
	err := conn.Connect(conf.DatabaseURI)
	if err != nil {
		log.Error("Database not allowed", err)
		return
	}
	if migrErr := conn.ApplyMigrations(); migrErr != nil {
		log.Error("Migrations not applied", migrErr)
		return
	}
	orderRep := orderRepository.New(conn, log)
	goodRep := goodRepository.New(conn, log)
	ruleRep := ruleRepository.New(conn, log)

	srvc := service.New(orderRep, goodRep, ruleRep, log)

	serv := server.New(srvc, log)
	serv.Start(conf.RunAddress)
}
