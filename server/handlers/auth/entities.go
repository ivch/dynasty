package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Session struct {
	ID           string
	UserID       uint
	RefreshToken uuid.UUID
	ExpiresIn    int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Token struct {
	ID   uint
	Name string
	Role uint
	jwt.StandardClaims
}
