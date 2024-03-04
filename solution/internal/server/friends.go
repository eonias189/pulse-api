package server

import (
	"fmt"
	"solution/internal/auth"
	"solution/internal/contract"
	"solution/internal/service"
	"solution/internal/utils"
	"solution/internal/validation"

	"github.com/gofiber/fiber/v2"
)

func handleFriends(r fiber.Router, s *service.Service) {
	r.Use(auth.AuthRequired(s))

	r.Post("/add", func(c *fiber.Ctx) error {

		payload, err := auth.GetJWTPayload(c)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusUnauthorized)
		}

		body := contract.LoginBody{}
		err = c.BodyParser(&body)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusBadRequest)
		}

		err = validation.ValidateLoginBody(body)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusBadRequest)
		}

		_, accepterNotExists := s.GetUserByLogin(body.Login)
		if accepterNotExists != nil {
			return utils.SendError(c, contract.NOT_FOUND("user", body.Login), fiber.StatusBadRequest)
		}

		_, relationNotExists := s.FindRelation(payload.Login, body.Login)
		if relationNotExists == nil {
			return c.JSON(contract.StatusResponse{Status: "ok"})
		}

		err = s.AddToFriends(payload.Login, body.Login)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusInternalServerError)
		}
		return c.JSON(contract.StatusResponse{Status: "ok"})
	})

	r.Post("/remove", func(c *fiber.Ctx) error {

		payload, err := auth.GetJWTPayload(c)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusUnauthorized)
		}

		body := contract.LoginBody{}
		err = c.BodyParser(&body)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusBadRequest)
		}

		err = validation.ValidateLoginBody(body)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusBadRequest)
		}

		relation, err := s.FindRelation(payload.Login, body.Login)
		if err != nil {
			return utils.SendError(c, contract.NOT_FOUND("friend", body.Login), fiber.StatusBadRequest)
		}

		err = s.DeleteRelation(relation)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusInternalServerError)
		}

		return c.JSON(contract.StatusResponse{Status: "ok"})
	})

	r.Get("/", func(c *fiber.Ctx) error {

		payload, err := auth.GetJWTPayload(c)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusUnauthorized)
		}

		limit := c.QueryInt("limit", 5)
		offset := c.QueryInt("offset", 0)

		fmt.Println(limit, offset)

		friends, err := s.GetFriends(payload.Login, limit, offset)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusInternalServerError)
		}

		return c.JSON(utils.Map(friends, func(a service.AccepterRelation) contract.Friend {
			return a.ToFriend()
		}))
	})
}
