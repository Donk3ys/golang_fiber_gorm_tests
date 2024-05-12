package views

import (
	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
)

type Session struct {
	ID uuid.UUID `json:"id"`
}

type SessionClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}
