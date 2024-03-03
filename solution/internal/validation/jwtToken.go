package validation

import (
	"solution/internal/contract"
	"solution/internal/service"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWTPayload(payload contract.JWTPayload, s *service.Service, timeout time.Duration) error {
	if time.Now().Unix() > payload.CreateTime+int64(timeout.Seconds()) {
		return jwt.ErrTokenExpired
	}

	user, err := s.GetUserByLogin(payload.Login)
	if err != nil {
		return contract.NOT_FOUND("user", payload.Login)
	}

	if user.PasswordChanged > payload.CreateTime {
		return contract.PASSWORD_CHANGED
	}
	return nil
}
