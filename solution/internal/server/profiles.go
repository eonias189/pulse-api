package server

import (
	"solution/internal/auth"
	"solution/internal/contract"
	"solution/internal/service"
	"solution/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func handleProfiles(r fiber.Router, s *service.Service) {
	r.Use(auth.AuthRequired(s))

	r.Get("/:login", func(c *fiber.Ctx) error {
		login := c.Params("login")
		if login == "" {
			return utils.SendError(c, contract.BAD_PATH_PARAM("login"), fiber.StatusBadRequest)
		}

		payload, err := auth.GetJWTPayload(c)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusUnauthorized)
		}

		user, err := s.GetUserByLogin(login)
		if err != nil {
			return utils.SendError(c, contract.NOT_FOUND("user", login), fiber.StatusForbidden)
		}

		_, err = s.FindRelation(user.Login, payload.Login)

		if user.Login != payload.Login && !user.IsPublic && err != nil {
			return utils.SendError(c, contract.ACCESS_DENIED, fiber.StatusForbidden)
		}

		return c.JSON(user.ToUserProfile())
	})
}
