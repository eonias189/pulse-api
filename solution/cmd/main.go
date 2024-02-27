package main

import (
	"log/slog"
	"os"
	"solution/internal/env"
	"solution/internal/server"
)

func main() {
	logger := slog.Default()

	pgUrl, err := env.GetPGURL()
	if err != nil {
		logger.Error("missed POSTGRES_CONN env")
		os.Exit(1)
	}
	logger.Info(pgUrl)

	serverAddress, err := env.GetServerAddress()
	if err != nil {
		logger.Error("missed SERVER_ADDRESS env (export smth like ':8080')")
		os.Exit(1)
	}

	s := server.NewServer(serverAddress, logger)
	if err = s.Start(); err != nil {
		logger.Error("server has been stopped", "error", err)
	}
}
