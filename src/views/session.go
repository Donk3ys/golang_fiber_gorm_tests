package views

import (
	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
)

type SessionClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}
