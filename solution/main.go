package main

import (
	"log/slog"
	"os"
	"solution/internal/config"
	"solution/internal/server"
	db "solution/internal/service"
)

func main() {
	logger := slog.Default()

	cfg, err := config.Get()
	if err != nil {
		logger.Error(err.Error())
	}

	service, err := db.New(cfg.PostressConn)
	if err != nil {
		logger.Error("Service creating error: " + err.Error())
		os.Exit(1)
	}

	s := server.NewServer(service, logger)
	if err = s.Run(cfg.ServerAddress); err != nil {
		logger.Error("server has been stopped", "error", err)
	}

}
