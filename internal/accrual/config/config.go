package config

import (
	"flag"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/config"
)

type EnvName = config.EnvName

const (
	// Env Names
	EnvNameRunAddress  = config.EnvNameRunAddress
	EnvNameDatabaseUri = config.EnvNameDatabaseUri

	// Default Values
	defaultValueRunAddress  string = "localhost:8081"
	defaultValueDatabaseUri string = "postgres://user:password@host:port/database"
)

type Config struct {
	RunAddress  string
	DatabaseUri string
}

func Initialize() *Config {
	conf := Config{
		RunAddress:  defaultValueRunAddress,
		DatabaseUri: defaultValueDatabaseUri,
	}

	conf.getFlags()
	conf.getEnvironments()

	return &conf
}

func (c *Config) getFlags() {
	argRunAddress := flag.String("a", defaultValueRunAddress, "Run address")
	argDatabaseUri := flag.String("d", defaultValueDatabaseUri, "Database uri")

	flag.Parse()

	if argRunAddress != nil && *argRunAddress != "" {
		c.RunAddress = *argRunAddress
	}
	if argDatabaseUri != nil && *argDatabaseUri != "" {
		c.DatabaseUri = *argDatabaseUri
	}
}

func (c *Config) getEnvironments() {
	runAdrress, ok := config.GetStringFromEnvironment(EnvNameRunAddress)
	if ok {
		c.RunAddress = runAdrress
	}

	databaseUri, ok := config.GetStringFromEnvironment(EnvNameDatabaseUri)
	if ok {
		c.DatabaseUri = databaseUri
	}
}
