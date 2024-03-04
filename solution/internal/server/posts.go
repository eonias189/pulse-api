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

func handlePosts(r fiber.Router, s *service.Service) {
	r.Use(auth.AuthRequired(s))

	r.Post("/new", func(c *fiber.Ctx) error {

		payload, err := auth.GetJWTPayload(c)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusUnauthorized)
		}

		body := contract.PostNewBody{}
		err = c.BodyParser(&body)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusBadRequest)
		}

		err = validation.ValidateNewPostBody(string(c.BodyRaw()))
		if err != nil {
			return utils.SendError(c, err, fiber.StatusBadRequest)
		}

		post := contract.Post{Id: contract.GenerateUUID(), Content: body.Content, Tags: body.Tags, Author: payload.Login, CreatedAt: time.Now().Unix()}
		err = s.AddPost(post)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusInternalServerError)
		}

		post, err = s.GetPostById(post.Id)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusInternalServerError)
		}
		return c.JSON(post.ToPostPreview())
	})

	r.Get("/:id", func(c *fiber.Ctx) error {

		payload, err := auth.GetJWTPayload(c)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusUnauthorized)
		}

		id := c.Params("id")
		if id == "" {
			return utils.SendError(c, contract.BAD_PATH_PARAM("postId"), fiber.StatusBadRequest)
		}

		post, err := s.GetPostById(id)
		if err != nil {
			return utils.SendError(c, contract.NOT_FOUND("post with id", id), fiber.StatusNotFound)
		}
		author, err := s.GetUserByLogin(post.Author)
		if err != nil {
			return utils.SendError(c, contract.NOT_FOUND("author", author.Login), fiber.StatusNotFound)
		}

		_, err = s.FindRelation(author.Login, payload.Login)
		if author.Login != payload.Login && !author.IsPublic && err != nil {
			return utils.SendError(c, contract.ACCESS_DENIED, fiber.StatusNotFound)
		}

		return c.JSON(post.ToPostPreview())
	})

	r.Get("/feed/my", func(c *fiber.Ctx) error {

		limit := c.QueryInt("limit", 5)
		offset := c.QueryInt("offset", 0)

		payload, err := auth.GetJWTPayload(c)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusUnauthorized)
		}
		posts, err := s.GetPostsOf(payload.Login, limit, offset)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusInternalServerError)
		}

		return c.JSON(utils.Map(posts, func(p contract.Post) contract.PostPreview {
			return p.ToPostPreview()
		}))
	})

	r.Get("/feed/:login", func(c *fiber.Ctx) error {

		limit := c.QueryInt("limit", 5)
		offset := c.QueryInt("offset", 0)
		login := c.Params("login")

		payload, err := auth.GetJWTPayload(c)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusUnauthorized)
		}

		author, err := s.GetUserByLogin(login)
		if err != nil {
			return utils.SendError(c, contract.NOT_FOUND("user", login), fiber.StatusNotFound)
		}

		_, err = s.FindRelation(author.Login, payload.Login)
		if author.Login != payload.Login && !author.IsPublic && err != nil {
			return utils.SendError(c, contract.ACCESS_DENIED, fiber.StatusNotFound)
		}

		posts, err := s.GetPostsOf(author.Login, limit, offset)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusInternalServerError)
		}
		return c.JSON(utils.Map(posts, func(p contract.Post) contract.PostPreview {
			return p.ToPostPreview()
		}))

	})

	r.Post("/:id/like", func(c *fiber.Ctx) error {

		postId := c.Params("id")

		payload, err := auth.GetJWTPayload(c)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusUnauthorized)
		}

		post, err := s.GetPostById(postId)
		if err != nil {
			return utils.SendError(c, contract.NOT_FOUND("post with id", postId), fiber.StatusForbidden)
		}

		author, err := s.GetUserByLogin(post.Author)
		if err != nil {
			return utils.SendError(c, contract.NOT_FOUND("author of post", post.Author), fiber.StatusNotFound)
		}

		_, err = s.FindRelation(author.Login, payload.Login)
		if author.Login != payload.Login && !author.IsPublic && err != nil {
			return utils.SendError(c, contract.ACCESS_DENIED, fiber.StatusNotFound)
		}

		err = s.SetReaction(payload.Login, postId, contract.Like)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusInternalServerError)
		}

		post, err = s.GetPostById(postId)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusInternalServerError)
		}

		return c.JSON(post.ToPostPreview())
	})

	r.Post("/:id/dislike", func(c *fiber.Ctx) error {

		postId := c.Params("id")

		payload, err := auth.GetJWTPayload(c)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusUnauthorized)
		}

		post, err := s.GetPostById(postId)
		if err != nil {
			return utils.SendError(c, contract.NOT_FOUND("post with id", postId), fiber.StatusForbidden)
		}

		author, err := s.GetUserByLogin(post.Author)
		if err != nil {
			return utils.SendError(c, contract.NOT_FOUND("author of post", post.Author), fiber.StatusNotFound)
		}

		_, err = s.FindRelation(author.Login, payload.Login)
		if author.Login != payload.Login && !author.IsPublic && err != nil {
			return utils.SendError(c, contract.ACCESS_DENIED, fiber.StatusNotFound)
		}

		err = s.SetReaction(payload.Login, postId, contract.Dislike)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusInternalServerError)
		}

		post, err = s.GetPostById(postId)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusInternalServerError)
		}

		return c.JSON(post.ToPostPreview())
	})
}
