package config

import (
	"flag"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/config"
)

type EnvName = config.EnvName

const (
	// Env names.
	EnvNameRunAddress             = config.EnvNameRunAddress
	EnvNameDatabaseURI            = config.EnvNameDatabaseURI
	EnvNameAccuralAddress EnvName = "ACCRUAL_SYSTEM_ADDRESS"
	EnvNameSecretKey              = "SECRET_KEY"

	// Default values.
	defaultValueRunAddress      string = "localhost:8080"
	defaultValueDatabaseURI     string = "postgres://user:password@host:port/database"
	defaultValueAccuralAddress  string = "localhost:8081"
	defaultValueWorkerPoolCount int    = 10
	defaultValueSecretKey       string = "FILATIK_SECRET_KEY_FOR_TOKEN"
)

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccuralSystemAddress string
	SecretKey            string
	WorkerPoolCount      int
}

func Initialize() *Config {
	conf := Config{
		RunAddress:           defaultValueRunAddress,
		DatabaseURI:          defaultValueDatabaseURI,
		AccuralSystemAddress: defaultValueAccuralAddress,
		WorkerPoolCount:      defaultValueWorkerPoolCount,
		SecretKey:            defaultValueSecretKey,
	}

	conf.getFlags()
	conf.getEnvironments()

	return &conf
}

func (c *Config) getFlags() {
	argRunAddress := flag.String("a", defaultValueRunAddress, "Run address")
	argDatabaseURI := flag.String("d", defaultValueDatabaseURI, "Database uri")
	argAccuralAddress := flag.String("r", defaultValueAccuralAddress, "Accural system address")
	argSecretKey := flag.String("k", defaultValueSecretKey, "Secret key for token")

	flag.Parse()

	if argRunAddress != nil && *argRunAddress != "" {
		c.RunAddress = *argRunAddress
	}
	if argDatabaseURI != nil && *argDatabaseURI != "" {
		c.DatabaseURI = *argDatabaseURI
	}
	if argAccuralAddress != nil && *argAccuralAddress != "" {
		c.AccuralSystemAddress = *argAccuralAddress
	}
	if argSecretKey != nil && *argSecretKey != "" {
		c.SecretKey = *argSecretKey
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

	accuralAddress, ok := config.GetStringFromEnvironment(EnvNameAccuralAddress)
	if ok {
		c.AccuralSystemAddress = accuralAddress
	}

	secretKey, ok := config.GetStringFromEnvironment(EnvNameSecretKey)
	if ok {
		c.SecretKey = secretKey
	}
}
