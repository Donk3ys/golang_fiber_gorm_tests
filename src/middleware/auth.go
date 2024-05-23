package middleware

import (
	"api/src/constants"
	"api/src/token"
	"api/src/views"

	"github.com/gofiber/fiber/v2"
)

func (i *Instance) AuthenticateAuthTokenAndCreateNewIfExpired(c *fiber.Ctx) error {
	authHeader := c.Get(constants.ACCESS_TOKEN_HEADER)
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No token!"})
	}

	token, err := token.DecodeAuthToken(authHeader)
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token error!"})
	}

	claims := token.Claims.(*views.SessionClaims)
	c.Locals(constants.REQ_USER_ID, claims.UserID)
	return c.Next()
}
