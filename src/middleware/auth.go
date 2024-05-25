package middleware

import (
	"api/src/constants"
	"api/src/token"
	"api/src/views"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (i *Instance) AuthenticateAuthToken(c *fiber.Ctx) error {
	authHeader := c.Get(constants.ACCESS_TOKEN_HEADER)
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No token!"})
	}

	token, err := token.DecodeAuthToken(authHeader)
	if err != nil || !token.Valid {
		log.Error("Bearer error: ", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token error!"})
	}

	claims := token.Claims.(*views.SessionClaims)
	c.Locals(constants.REQ_USER_ID, claims.UserID)
	return c.Next()
}
