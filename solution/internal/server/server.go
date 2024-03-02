package server

import (
	"log/slog"
	"net/http"
	"solution/internal/contract"
	serv "solution/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

func sendError(c *fiber.Ctx, err error, status int) error {
	return c.Status(status).JSON(contract.NewErrorResp(err))
}

func handlePing(c *fiber.Ctx) error {
	return c.SendString("ok")
}

func (s *Server) Run(address string) error {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New(recover.Config{}))
	app.Use(cors.New(cors.Config{}))

	api := app.Group("/api")

	api.Get("/ping", handlePing)
	handleCountries(api.Group("/countries"), s.service)
	handleAuth(api.Group("/auth"), s.service)

	r := app.Group("/nr")
	r.Use(AuthRequired())
	r.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("access")
	})

	s.logger.Info("server has been started", "address", address)

	err := app.Listen(address)
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}
