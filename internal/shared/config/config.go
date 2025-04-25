package config

import (
	"os"
	"strconv"
)

func GetStringFromEnvironment(key string) (string, bool) {
	val, ok := os.LookupEnv(key)
	if ok && val != "" {
		return val, true
	}
	return "", false
}

func GetInt64FromEnvironment(key string) (int64, bool) {
	val, ok := os.LookupEnv(key)
	if ok && val != "" {
		if num, err := strconv.ParseInt(val, 10, 64); err == nil {
			return num, true
		}
	}
	return 0, false
}
