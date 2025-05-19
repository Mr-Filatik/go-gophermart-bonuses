package config

import (
	"os"
	"strconv"
)

type EnvName string

const (
	EnvNameRunAddress  EnvName = "RUN_ADDRESS"
	EnvNameDatabaseURI EnvName = "DATABASE_URI"
)

func GetStringFromEnvironment(key EnvName) (string, bool) {
	val, ok := os.LookupEnv(string(key))
	if ok && val != "" {
		return val, true
	}
	return "", false
}

func GetInt64FromEnvironment(key EnvName) (int64, bool) {
	val, ok := os.LookupEnv(string(key))
	if ok && val != "" {
		if num, err := strconv.ParseInt(val, 10, 64); err == nil {
			return num, true
		}
	}
	return 0, false
}
