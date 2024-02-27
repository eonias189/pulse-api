package env

import (
	"os"
	"solution/internal/contract"
)

func GetPGURL() (string, error) {
	pgURL := os.Getenv("POSTGRES_CONN")
	if pgURL == "" {
		return "", contract.ENV_ERROR
	}
	return pgURL, nil
}

func GetServerAddress() (string, error) {
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		return "", contract.ENV_ERROR
	}
	return serverAddress, nil
}
