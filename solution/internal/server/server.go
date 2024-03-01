package server

import (
	"log/slog"
	"net/http"
	serv "solution/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type Server struct {
	service *serv.Service
	logger  *slog.Logger
}

func NewServer(service *serv.Service, logger *slog.Logger) *Server {
	return &Server{
		logger:  logger,
		service: service,
	}
}

func (s *Server) Run(address string) error {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	api := app.Group("/api")

	api.Get("/ping", s.handlePing)

	countries := api.Group("/countries")
	countries.Get("/", s.handleCountriesIndex)
	countries.Get("/:alpha2", s.handleCountriesAlpha2)

	s.logger.Info("server has been started", "address", address)

	err := app.Listen(address)
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}
