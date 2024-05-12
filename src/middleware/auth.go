package middleware

import (
	"api/src/constants"
	"api/src/models"
	"api/src/views"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
)

func (i *Instance) AuthenticateAuthTokenAndCreateNewIfExpired(c *fiber.Ctx) error {
	token, bearer, err := getAuthToken(c)

	// Valid
	if err == nil && token.Valid {
		addClaimsToCtx(c, token)
		return c.Next()
	}

	// Check jwt only expired
	ve, ok := err.(*jwt.ValidationError)
	if !ok || ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) == 0 {
		// log.Error("JWT INVAILD")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token unauthorized"})
	}
	// log.Error("JWT EXPIRED")

	// Check session not expired
	var exSession models.Session
	i.Repo.Db.First(&exSession, "(token=? AND expires_at>?) OR from_token=?", bearer, time.Now(), bearer)
	if exSession.Token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Session expired"})
	}

	// Add claims to request
	addClaimsToCtx(c, token)

	// Check if new session needs to be created
	if bearer == exSession.Token {
		// log.Error("JWT CREATE SESSION")
		// Create new session
		newBearer, err := i.Repo.UpdateUserSession(bearer, exSession.UserID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Error creating new session!"})
		}
		c.Set("Access-Control-Expose-Headers", "*") // Needed to show headers in web app
		c.Set(constants.AUTH_HEADER, "bearer "+newBearer)

		// If request takes longer then the auth token expiry duratation then create new sesssion
		reqTimerStart := time.Now()
		defer func() {
			reqTimerEnd := time.Now()
			reqTimerDur := reqTimerEnd.Sub(reqTimerStart)
			if reqTimerDur > constants.AUTH_TOKEN_DURATION {
				newBearer, _ := i.Repo.UpdateUserSession(bearer, exSession.UserID)
				c.Set(constants.AUTH_HEADER, "bearer "+newBearer)
			}
		}()
		return c.Next()
	}

	// log.Error("JWT USE OLD TOKEN FOR SESSION")
	return c.Next()
}

func ParseBearerToken(bearer string) string {
	return strings.Replace(bearer, constants.AUTH_HEADER+" ", "", 1)
}

func getAuthToken(c *fiber.Ctx) (*jwt.Token, string, error) {
	authHeader := c.Get(constants.AUTH_HEADER)
	if authHeader == "" {
		authHeader = c.Query(constants.AUTH_HEADER)
	}
	bearer := ParseBearerToken(authHeader)
	token, err := jwt.ParseWithClaims(bearer, &views.SessionClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(constants.AUTH_SECRET)), nil
	})

	return token, bearer, err
}

func addClaimsToCtx(c *fiber.Ctx, token *jwt.Token) {
	claims := token.Claims.(*views.SessionClaims)
	c.Locals(constants.REQ_USER_ID, claims.UserID)
}

func UserIDFromAuthToken(c *fiber.Ctx) *uuid.UUID {
	at, _, _ := getAuthToken(c)
	if at == nil {
		return nil
	}

	return &at.Claims.(*views.SessionClaims).UserID
}
