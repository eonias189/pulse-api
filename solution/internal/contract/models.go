package contract

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/golang-jwt/jwt/v5"
)

type Country struct {
	Id     int    `json:"-"`
	Name   string `json:"name"`
	Alpha2 string `json:"alpha2"`
	Alpha3 string `json:"alpha3"`
	Region string `json:"region"`
}

type User struct {
	Id          string `json:"-"`
	Login       string `json:"login"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	CountryCode string `json:"countryCode"`
	IsPublic    bool   `json:"isPublic"`
	Phone       string `json:"phone"`
	Image       string `json:"image"`
}

type UserProfile struct {
	Login       string `json:"login"`
	Email       string `json:"email"`
	CountryCode string `json:"countryCode"`
	IsPublic    bool   `json:"isPublic"`
	Phone       string `json:"phone,omitempty"`
	Image       string `json:"image,omitempty"`
}

type JWTPayload struct {
	Login   string `json:"login"`
	Timeout int    `json:"timeout"`
}

func (u User) ToUserProfile() UserProfile {
	return UserProfile{
		Login:       u.Login,
		Email:       u.Email,
		CountryCode: u.CountryCode,
		IsPublic:    u.IsPublic,
		Phone:       u.Phone,
		Image:       u.Image,
	}
}

func (j JWTPayload) ToClaims() jwt.Claims {
	return jwt.MapClaims{
		"login":   j.Login,
		"timeout": j.Timeout,
	}
}

func GenerateUUID() string {
	u := make([]byte, 16)
	rand.Read(u)

	u[8] = (u[8] | 0x80) & 0xBF // what does this do?
	u[6] = (u[6] | 0x40) & 0x4F // what does this do?

	return hex.EncodeToString(u)
}
