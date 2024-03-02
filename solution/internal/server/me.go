package server

import (
	"solution/internal/contract"
	"solution/internal/service"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func handleMe(r fiber.Router, s *service.Service) {
	r.Use(AuthRequired(s))

	r.Get("/profile", func(c *fiber.Ctx) error {

		payload, err := GetJWTPayload(c)
		if err != nil {
			return sendError(c, jwt.ErrTokenMalformed, fiber.StatusUnauthorized)
		}

		user, err := s.GetUserByLogin(payload.Login)
		if err != nil {
			return sendError(c, err, fiber.StatusInternalServerError)
		}

		return c.JSON(user.ToUserProfile())
	})

	r.Patch("/profile", func(c *fiber.Ctx) error {
		payload, err := GetJWTPayload(c)
		if err != nil {
			return sendError(c, jwt.ErrTokenMalformed, fiber.StatusUnauthorized)
		}

		body := contract.MeProfileBody{}
		if err = c.BodyParser(&body); err != nil {
			return sendError(c, err, fiber.StatusBadRequest)
		}

		user, err := s.GetUserByLogin(payload.Login)
		if err != nil {
			return sendError(c, err, fiber.StatusInternalServerError)
		}

		if !strings.Contains(string(c.Body()), `"isPublic"`) {
			body.IsPublic = user.IsPublic
		}

		if body.Image == "" {
			body.Image = user.Image
		}

		if body.Phone == "" {
			body.Phone = user.Phone
		}

		if body.CountryCode == "" {
			body.CountryCode = user.CountryCode
		} else {
			_, err = s.GetCountryByAlpha2(body.CountryCode)
			if err != nil {
				return sendError(c, contract.NOT_FOUND("country with alpha2", body.CountryCode), fiber.StatusBadRequest)
			}
		}

		if len(body.Image) > 200 {
			return sendError(c, contract.BAD_BODY_PARAM("image must be not loneger than 200 symbols"), fiber.StatusBadRequest)
		}

		if s.PhoneExists(body.ToUser()) {
			return sendError(c, contract.USER_ALREADY_EXISTS, fiber.StatusConflict)
		}

		err = s.UpdateUser(user, body.ToUser())
		if err != nil {
			return sendError(c, err, fiber.StatusInternalServerError)
		}

		newUser, err := s.GetUserByLogin(user.Login)
		if err != nil {
			return sendError(c, err, fiber.StatusInternalServerError)
		}

		return c.JSON(newUser.ToUserProfile())
	})
}
