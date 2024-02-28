package config

import (
	"os"
	"solution/internal/contract"
)

type Config struct {
	PostressConn  string
	ServerAddress string
}

func Get() (*Config, error) {
	cfg := &Config{}

	pgConn := os.Getenv("POSTGRES_CONN")
	if pgConn == "" {
		return cfg, contract.ENV_ERROR("POSTHRES_CONN")
	}

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		return cfg, contract.ENV_ERROR("SERVER_ADDRESS")
	}

	cfg.PostressConn = pgConn
	cfg.ServerAddress = serverAddress
	return cfg, nil
}
