package server

import (
	"solution/internal/auth"
	"solution/internal/contract"
	"solution/internal/service"
	"solution/internal/utils"
	"solution/internal/validation"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func handleMe(r fiber.Router, s *service.Service) {
	r.Use(auth.AuthRequired(s))
	r.Get("/profile", func(c *fiber.Ctx) error {

		payload, err := auth.GetJWTPayload(c)
		if err != nil {
			return utils.SendError(c, jwt.ErrTokenMalformed, fiber.StatusUnauthorized)
		}

		return c.JSON(payload.User.ToUserProfile())
	})

	r.Patch("/profile", func(c *fiber.Ctx) error {
		payload, err := auth.GetJWTPayload(c)
		if err != nil {
			return utils.SendError(c, jwt.ErrTokenMalformed, fiber.StatusUnauthorized)
		}

		body := contract.MeProfileBody{}
		if err = c.BodyParser(&body); err != nil {
			return utils.SendError(c, err, fiber.StatusBadRequest)
		}

		if !strings.Contains(string(c.Body()), `"isPublic"`) {
			body.IsPublic = payload.User.IsPublic
		}

		if body.Image == "" {
			body.Image = payload.User.Image
		}

		if body.Phone == "" {
			body.Phone = payload.User.Phone
		}

		if body.CountryCode == "" {
			body.CountryCode = payload.User.CountryCode
		} else {
			_, err = s.GetCountryByAlpha2(body.CountryCode)
			if err != nil {
				return utils.SendError(c, contract.NOT_FOUND("country with alpha2", body.CountryCode), fiber.StatusBadRequest)
			}
		}

		err = validation.ValidateImage(body.Image)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusBadRequest)
		}

		if s.UserDataExists(body.ToUser()) {
			return utils.SendError(c, contract.USER_ALREADY_EXISTS, fiber.StatusConflict)
		}

		newUser := body.ToUser()
		newUser.Password = payload.User.Password
		newUser.PasswordChanged = payload.User.PasswordChanged
		err = s.UpdateUser(payload.User, newUser)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusInternalServerError)
		}

		newUser, err = s.GetUserByLogin(payload.Login)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusInternalServerError)
		}

		return c.JSON(newUser.ToUserProfile())
	})

	r.Post("/updatePassword", func(c *fiber.Ctx) error {

		payload, err := auth.GetJWTPayload(c)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusUnauthorized)
		}

		body := contract.MeUpdatePasswordBody{}
		err = c.BodyParser(&body)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusBadRequest)
		}

		err = validation.ValidateUpdatePassword(body)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusBadRequest)
		}

		err = payload.User.CheckPassword(body.OldPassword)
		if err != nil {
			return utils.SendError(c, contract.INCORRECT_PASSWORD, fiber.StatusForbidden)
		}

		err = validation.ValidatePassword(body.NewPassword)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusBadRequest)
		}

		newUser := payload.User
		newUser.Password = body.NewPassword
		newUser.PasswordChanged = time.Now().Unix()
		newUser.Password, err = newUser.HashPassword()
		if err != nil {
			return utils.SendError(c, err, fiber.StatusInternalServerError)
		}

		err = s.UpdateUser(payload.User, newUser)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusInternalServerError)
		}

		return c.JSON(contract.StatusResponse{Status: "ok"})

	})
}
