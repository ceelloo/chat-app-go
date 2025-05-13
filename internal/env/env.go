package env

import (
	"os"
	"strconv"
)
// Use direnv and .envrc to set environment variables

func GetString(key, fallback string) string {
	val, err := os.LookupEnv(key)
	if !err {
		return fallback
	}

	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	toInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	
	return toInt
}