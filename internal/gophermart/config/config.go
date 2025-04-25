package config

import (
	"flag"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/config"
)

type EnvName = config.EnvName

const (
	// Env Names
	EnvNameRunAddress             = config.EnvNameRunAddress
	EnvNameDatabaseUri            = config.EnvNameDatabaseUri
	EnvNameAccuralAddress EnvName = "ACCRUAL_SYSTEM_ADDRESS"

	// Default Values
	defaultValueRunAddress     string = "localhost:8080"
	defaultValueDatabaseUri    string = "..."
	defaultValueAccuralAddress string = "localhost:8081"
)

type Config struct {
	RunAddress           string
	DatabaseUri          string
	AccuralSystemAddress string
}

func Initialize() *Config {
	conf := Config{
		RunAddress:           defaultValueRunAddress,
		DatabaseUri:          defaultValueDatabaseUri,
		AccuralSystemAddress: defaultValueAccuralAddress,
	}

	conf.getFlags()
	conf.getEnvironments()

	return &conf
}

func (c *Config) getFlags() {
	argRunAddress := flag.String("a", defaultValueRunAddress, "Run address")
	argDatabaseUri := flag.String("d", defaultValueDatabaseUri, "Database uri")
	argAccuralAddress := flag.String("r", defaultValueAccuralAddress, "Accural system address")

	flag.Parse()

	if argRunAddress != nil && *argRunAddress != "" {
		c.RunAddress = *argRunAddress
	}
	if argDatabaseUri != nil && *argDatabaseUri != "" {
		c.DatabaseUri = *argDatabaseUri
	}
	if argAccuralAddress != nil && *argAccuralAddress != "" {
		c.AccuralSystemAddress = *argAccuralAddress
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

	accuralAddress, ok := config.GetStringFromEnvironment(EnvNameAccuralAddress)
	if ok {
		c.AccuralSystemAddress = accuralAddress
	}
}
