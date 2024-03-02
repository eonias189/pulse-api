package server

import (
	"solution/internal/contract"
	"solution/internal/service"
	"strings"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKey    = []byte("very very secret")
	contextKey   = "user"
	tokenTimeout = time.Hour * 4
	jwtCfg       = jwtware.Config{
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
		SigningKey: jwtware.SigningKey{
			Key:    secretKey,
			JWTAlg: jwtware.HS256,
		},
		ContextKey: contextKey,
	}
)

func AuthRequired(s *service.Service) func(*fiber.Ctx) error {
	newCfg := jwtCfg
	newCfg.ErrorHandler = func(c *fiber.Ctx, err error) error {
		return sendError(c, err, fiber.StatusUnauthorized)
	}
	newCfg.SuccessHandler = func(c *fiber.Ctx) error {

		payload, err := GetJWTPayload(c)
		if err != nil {
			return sendError(c, err, fiber.StatusUnauthorized)
		}

		if time.Now().Unix() > payload.CreateTime+int64(tokenTimeout.Seconds()) {
			return sendError(c, jwt.ErrTokenExpired, fiber.StatusUnauthorized)
		}

		user, err := s.GetUserByLogin(payload.Login)
		if err != nil {
			return sendError(c, err, fiber.StatusInternalServerError)
		}

		if user.PasswordChanged > payload.CreateTime {
			return sendError(c, contract.PASSWORD_CHANGED, fiber.StatusUnauthorized)
		}

		return c.Next()
	}
	return jwtware.New(newCfg)
}

func GetJWTPayload(c *fiber.Ctx) (contract.JWTPayload, error) {
	payload := contract.JWTPayload{}
	jwtToken, ok := c.Context().Value(contextKey).(*jwt.Token)
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

func validatePassword(password string) error {
	if len(password) < 6 {
		return contract.INVALID_PASSWORD
	}

	latin := "abcdefghijklmnopqrstuvwxyz"
	digits := "0123456789"

	if !strings.ContainsAny(password, latin) {
		return contract.INVALID_PASSWORD
	}

	if !strings.ContainsAny(password, strings.ToUpper(latin)) {
		return contract.INVALID_PASSWORD
	}

	if !strings.ContainsAny(password, digits) {
		return contract.INVALID_PASSWORD
	}

	return nil

}

func validateRegister(user contract.RegisterBody, s *service.Service) error {
	switch "" {
	case user.Login:
		return contract.MISSING_FIELD("login")
	case user.Email:
		return contract.MISSING_FIELD("email")
	case user.Password:
		return contract.MISSING_FIELD("password")
	case user.CountryCode:
		return contract.MISSING_FIELD("countryCode")
	}

	if len(user.Image) > 200 {
		return contract.BAD_BODY_PARAM("image len must be not longer than 200 symbols")
	}

	_, err := s.GetCountryByAlpha2(user.CountryCode)
	if err != nil {
		return contract.UNKNOWN_COUNTRY_CODE(user.CountryCode)
	}

	return validatePassword(user.Password)
}

func validateSignIn(user contract.SignInBody) error {
	switch "" {
	case user.Login:
		return contract.MISSING_FIELD("login")
	case user.Password:
		return contract.MISSING_FIELD("password")
	}
	return nil
}

func generateJWT(payload contract.JWTPayload) (string, error) {
	claims := payload.ToClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}

func handleAuth(r fiber.Router, s *service.Service) {
	r.Post("/register", func(c *fiber.Ctx) error {

		body := contract.RegisterBody{}
		err := c.BodyParser(&body)
		if err != nil {
			return sendError(c, err, fiber.StatusBadRequest)
		}

		err = validateRegister(body, s)
		if err != nil {
			return sendError(c, err, fiber.StatusBadRequest)
		}

		user := contract.User(body)
		if s.UserExists(user) {
			return sendError(c, contract.USER_ALREADY_EXISTS, fiber.StatusConflict)
		}

		user.Id = contract.GenerateUUID()
		hashedPassword, err := user.HashPassword()
		if err != nil {
			return sendError(c, err, fiber.StatusInternalServerError)
		}

		user.Password = string(hashedPassword)
		user.PasswordChanged = time.Now().Unix()
		err = s.AddUser(user)
		if err != nil {
			return sendError(c, err, fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusCreated).JSON(contract.RegisterResp{Profile: user.ToUserProfile()})

	})

	r.Post("/sign-in", func(c *fiber.Ctx) error {

		body := contract.SignInBody{}
		err := c.BodyParser(&body)
		if err != nil {
			return sendError(c, err, fiber.StatusBadRequest)
		}

		err = validateSignIn(body)
		if err != nil {
			return sendError(c, err, fiber.StatusBadRequest)
		}

		user, err := s.GetUserByLogin(body.Login)
		if err != nil {
			return sendError(c, contract.BAD_CRENDIALS, fiber.StatusUnauthorized)
		}

		if err = user.CheckPassword(body.Password); err != nil {
			return sendError(c, contract.BAD_CRENDIALS, fiber.StatusUnauthorized)
		}

		token, err := generateJWT(contract.JWTPayload{
			Login:      user.Login,
			CreateTime: time.Now().Add(tokenTimeout).Unix(),
		})
		if err != nil {
			return sendError(c, contract.BAD_CRENDIALS, fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusOK).JSON(contract.SignInResp{Token: token})

	})
}
