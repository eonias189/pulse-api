package auth

import (
	"context"
	"solution/internal/contract"
	"solution/internal/service"
	"solution/internal/utils"
	"solution/internal/validation"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKey    = []byte("very very secret")
	tokenKey     = "user"
	payloadKey   = "payload"
	tokenTimeout = time.Hour * 4
)

func AuthRequired(s *service.Service) func(*fiber.Ctx) error {
	cfg := jwtware.Config{
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
		SigningKey: jwtware.SigningKey{
			Key:    secretKey,
			JWTAlg: jwtware.HS256,
		},
		ContextKey: tokenKey,
	}
	cfg.ErrorHandler = func(c *fiber.Ctx, err error) error {
		return utils.SendError(c, err, fiber.StatusUnauthorized)
	}
	cfg.SuccessHandler = func(c *fiber.Ctx) error {
		payload, err := getTokenData(c)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusUnauthorized)
		}

		user, err := validation.ValidateJWTPayload(payload, s, tokenTimeout)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusUnauthorized)
		}

		payload.User = user
		c.SetUserContext(context.WithValue(c.Context(), payloadKey, payload))

		return c.Next()
	}
	return jwtware.New(cfg)
}

func getTokenData(c *fiber.Ctx) (contract.JWTPayload, error) {
	payload := contract.JWTPayload{}
	jwtToken, ok := c.Context().Value(tokenKey).(*jwt.Token)
	if !ok {
		return payload, jwt.ErrTokenMalformed
	}

	claimMap, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return payload, jwt.ErrTokenMalformed
	}

	payload.Login, ok = claimMap["login"].(string)
	if !ok {
		return payload, jwt.ErrTokenInvalidClaims
	}

	timeout, ok := claimMap["createTime"].(float64)
	if !ok {
		return payload, jwt.ErrTokenInvalidClaims
	}

	payload.CreateTime = int64(timeout)
	return payload, nil
}

func GetJWTPayload(c *fiber.Ctx) (contract.JWTPayload, error) {
	payload, ok := c.UserContext().Value(payloadKey).(contract.JWTPayload)
	if !ok {
		return payload, jwt.ErrTokenMalformed
	}
	return payload, nil

}

func GenerateJWT(payload contract.JWTPayload) (string, error) {
	claims := payload.ToClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}
