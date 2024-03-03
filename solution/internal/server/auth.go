package server

import (
	"solution/internal/auth"
	"solution/internal/contract"
	"solution/internal/service"
	"solution/internal/utils"
	"solution/internal/validation"
	"time"

	"github.com/gofiber/fiber/v2"
)

func handleAuth(r fiber.Router, s *service.Service) {
	r.Post("/register", func(c *fiber.Ctx) error {

		body := contract.RegisterBody{}
		err := c.BodyParser(&body)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusBadRequest)
		}

		err = validation.ValidateRegister(body, s)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusBadRequest)
		}

		user := contract.User(body)
		if s.UserExists(user) {
			return utils.SendError(c, contract.USER_ALREADY_EXISTS, fiber.StatusConflict)
		}

		user.Password, err = user.HashPassword()
		if err != nil {
			return utils.SendError(c, err, fiber.StatusInternalServerError)
		}

		user.PasswordChanged = time.Now().Unix()
		err = s.AddUser(user)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusCreated).JSON(contract.RegisterResp{Profile: user.ToUserProfile()})

	})

	r.Post("/sign-in", func(c *fiber.Ctx) error {

		body := contract.SignInBody{}
		err := c.BodyParser(&body)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusBadRequest)
		}

		err = validation.ValidateSignIn(body, s)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusBadRequest)
		}

		user, err := s.GetUserByLogin(body.Login)
		if err != nil {
			return utils.SendError(c, contract.BAD_CRENDIALS, fiber.StatusUnauthorized)
		}

		if err = user.CheckPassword(body.Password); err != nil {
			return utils.SendError(c, contract.BAD_CRENDIALS, fiber.StatusUnauthorized)
		}

		token, err := auth.GenerateJWT(contract.JWTPayload{
			Login:      user.Login,
			CreateTime: time.Now().Unix(),
		})
		if err != nil {
			return utils.SendError(c, contract.BAD_CRENDIALS, fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusOK).JSON(contract.SignInResp{Token: token})

	})
}
