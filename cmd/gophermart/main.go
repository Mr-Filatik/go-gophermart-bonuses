package main

import (
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/config"
	orderRepository "github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/database/repository/order"
	userRepository "github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/database/repository/user"
	withdrawalRepository "github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/database/repository/withdrawal"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/server"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/service"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/worker"
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
	userRep := userRepository.New(conn, log)
	orderRep := orderRepository.New(conn, log)
	wdrwlsRep := withdrawalRepository.New(conn, log)

	wrkr := worker.New(orderRep, userRep, conf.AccuralSystemAddress, log)
	wrkr.Run()
	defer wrkr.Close()

	srvc := service.New(userRep, orderRep, wdrwlsRep, wrkr, log)

	serv := server.New(srvc, log)
	serv.Start(conf.RunAddress)
}
