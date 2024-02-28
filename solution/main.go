package main

import (
	"log/slog"
	"os"
	"solution/internal/config"
	"solution/internal/db"
	"solution/internal/server"
)

func main() {
	logger := slog.Default()

	cfg, err := config.Get()
	if err != nil {
		logger.Error(err.Error())
	}

	dataBase := db.NewDB(cfg.PostressConn)
	logger.Info("start connecting")
	err = dataBase.Connect()
	if err != nil {
		logger.Error("unable to connect to database: " + err.Error())
		os.Exit(1)
	}
	logger.Info("successfully connected")

	s := server.NewServer(cfg.ServerAddress, dataBase, logger)
	if err = s.Start(); err != nil {
		logger.Error("server has been stopped", "error", err)
	}

}
