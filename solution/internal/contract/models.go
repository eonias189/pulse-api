package contract

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	Like    ReactionType = "like"
	Dislike ReactionType = "dislike"
)

type Country struct {
	Id     int    `json:"-"`
	Name   string `json:"name"`
	Alpha2 string `json:"alpha2"`
	Alpha3 string `json:"alpha3"`
	Region string `json:"region"`
}

type User struct {
	Login           string `json:"login"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	CountryCode     string `json:"countryCode"`
	IsPublic        bool   `json:"isPublic"`
	Phone           string `json:"phone"`
	Image           string `json:"image"`
	PasswordChanged int64  `json:"-"`
}

type Relation struct {
	Id            int
	SenderLogin   string
	AccepterLogin string
	CreateTime    int64
}

type Post struct {
	Id            string
	Content       string
	Author        string
	Tags          []string
	CreatedAt     int64
	LikesCount    int
	DislikesCount int
}

type ReactionType string

type Reaction struct {
	Login  string
	PostId string
	Type   ReactionType
}

type UserProfile struct {
	Login       string `json:"login"`
	Email       string `json:"email"`
	CountryCode string `json:"countryCode"`
	IsPublic    bool   `json:"isPublic"`
	Phone       string `json:"phone,omitempty"`
	Image       string `json:"image,omitempty"`
}

type PostPreview struct {
	Id            string   `json:"id"`
	Content       string   `json:"content"`
	Author        string   `json:"author"`
	Tags          []string `json:"tags"`
	CreatedAt     string   `json:"createdAt"`
	LikesCount    int      `json:"likesCount"`
	DislikesCount int      `json:"dislikesCount"`
}

type JWTPayload struct {
	User       User   `json:"-"`
	Login      string `json:"login"`
	CreateTime int64  `json:"createTime"`
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

func (p Post) ToPostPreview() PostPreview {
	return PostPreview{
		Id:            p.Id,
		Content:       p.Content,
		Author:        p.Author,
		Tags:          p.Tags,
		CreatedAt:     time.Unix(p.CreatedAt, 0).Format(time.RFC3339),
		LikesCount:    p.LikesCount,
		DislikesCount: p.DislikesCount,
	}
}

func (j JWTPayload) ToClaims() jwt.Claims {
	return jwt.MapClaims{
		"login":      j.Login,
		"createTime": j.CreateTime,
	}
}

func (u User) HashPassword() (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(password), nil

}

func (u User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func GenerateUUID() string {
	id := uuid.New()
	return fmt.Sprint(id)
}
