package server

import (
	"solution/internal/contract"
	"solution/internal/service"
	"strings"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	secretKey    = []byte("very very secret")
	contextKey   = "user"
	tokenTimeout = time.Second * 30
	jwtCfg       = jwtware.Config{
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
		SigningKey: jwtware.SigningKey{
			Key:    secretKey,
			JWTAlg: jwtware.HS256,
		},
		ContextKey: contextKey,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return sendError(c, err, fiber.StatusUnauthorized)
		},
	}
)

func AuthRequired() func(*fiber.Ctx) error {
	return jwtware.New(jwtCfg)
}

// func GetJWTPayload(c *fiber.Ctx) (contract.JWTPayload, error)

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
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return sendError(c, err, fiber.StatusInternalServerError)
		}

		user.Password = string(hashedPassword)
		err = s.AddUser(user)
		if err != nil {
			return sendError(c, err, fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusCreated).JSON(user.ToUserProfile())

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

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
			return sendError(c, contract.BAD_CRENDIALS, fiber.StatusUnauthorized)
		}

		token, err := generateJWT(contract.JWTPayload{Login: user.Login, Timeout: int(time.Now().Add(tokenTimeout).Unix())})
		if err != nil {
			return sendError(c, contract.BAD_CRENDIALS, fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusOK).JSON(contract.SignInResp{Token: token})

	})
}
