package config

import (
	"flag"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/config"
)

type EnvName = config.EnvName

const (
	// Env names.
	EnvNameRunAddress  = config.EnvNameRunAddress
	EnvNameDatabaseURI = config.EnvNameDatabaseURI

	// Default values.
	defaultValueRunAddress      string = "localhost:8081"
	defaultValueDatabaseURI     string = "postgres://user:password@host:port/database"
	defaultValueWorkerPoolCount int    = 10
)

type Config struct {
	RunAddress      string
	DatabaseURI     string
	WorkerPoolCount int
}

func Initialize() *Config {
	conf := Config{
		RunAddress:      defaultValueRunAddress,
		DatabaseURI:     defaultValueDatabaseURI,
		WorkerPoolCount: defaultValueWorkerPoolCount,
	}

	conf.getFlags()
	conf.getEnvironments()

	return &conf
}

func (c *Config) getFlags() {
	argRunAddress := flag.String("a", defaultValueRunAddress, "Run address")
	argDatabaseURI := flag.String("d", defaultValueDatabaseURI, "Database uri")

	flag.Parse()

	if argRunAddress != nil && *argRunAddress != "" {
		c.RunAddress = *argRunAddress
	}
	if argDatabaseURI != nil && *argDatabaseURI != "" {
		c.DatabaseURI = *argDatabaseURI
	}
}

func (c *Config) getEnvironments() {
	runAdrress, ok := config.GetStringFromEnvironment(EnvNameRunAddress)
	if ok {
		c.RunAddress = runAdrress
	}

	databaseURI, ok := config.GetStringFromEnvironment(EnvNameDatabaseURI)
	if ok {
		c.DatabaseURI = databaseURI
	}
}
