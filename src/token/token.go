package token

import (
	"api/src/constants"
	"api/src/views"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
)

func ParseBearerToken(bearer string) string {
	bearer = strings.TrimSpace(bearer)
	return strings.Replace(bearer, "bearer ", "", 1)
}

func DecodeAuthToken(authToken string) (*jwt.Token, error) {
	bearer := ParseBearerToken(authToken)
	token, err := jwt.ParseWithClaims(bearer, &views.SessionClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(constants.ACCESS_TOKEN_SECRET)), nil
	})

	return token, err
}

func DecodeRefreshToken(refreshToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &views.SessionClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(constants.REFRESH_TOKEN_HEADER)), nil
	})

	return token, err
}

func CreateAuthToken(userID uuid.UUID) (string, *time.Time, error) {
	bearerExpiry := time.Now().Add(constants.ACCESS_TOKEN_DURATION)
	bearerClaims := views.SessionClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: bearerExpiry.Unix(),
		},
	}
	bearer, err := jwt.NewWithClaims(jwt.SigningMethodHS256, bearerClaims).SignedString([]byte(os.Getenv(constants.ACCESS_TOKEN_SECRET)))
	if err != nil {
		log.Error("jwt sign bearer error. ", err)
		return "", nil, errors.New("Could not create authentication token for user[001]")
	}
	return bearer, &bearerExpiry, nil
}

func CreateRefreshToken(userID uuid.UUID) (string, *time.Time, error) {
	refreshExpiry := time.Now().Add(constants.REFRESH_TOKEN_DURATION)
	refreshClaims := views.SessionClaims{
		UserID: userID,
		// UserRole: userModel.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpiry.Unix(),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv(constants.REFRESH_TOKEN_SECRET)))
	if err != nil {
		log.Error("jwt sign refresh error. ", err)
		return "", nil, errors.New("Could not crete authentication for user")
	}

	return refreshToken, &refreshExpiry, nil
}
