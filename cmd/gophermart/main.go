package main

import (
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/config"
	userRepository "github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/database/repository/user"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/server"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/gophermart/service"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/database/connector"
	logger "github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger/zap/sugar"
)

func main() {
	log := logger.New(logger.LevelDebug)
	defer log.Close()

	conf := config.Initialize()

	conn := connector.New(log)
	err := conn.Connect(conf.DatabaseUri)
	if err != nil {
		log.Error("Database not allowed", err)
		return
	}
	if migrErr := conn.ApplyMigrations(); migrErr != nil {
		log.Error("Migrations not applied", migrErr)
		return
	}
	userRep := userRepository.New(conn, log)

	srvc := service.New(userRep, log)

	serv := server.New(srvc, log)
	serv.Start(conf.RunAddress)
}
