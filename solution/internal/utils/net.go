package utils

import (
	"solution/internal/contract"

	"github.com/gofiber/fiber/v2"
)

func SendError(c *fiber.Ctx, err error, status int) error {
	return c.Status(status).JSON(contract.NewErrorResp(err))
}
