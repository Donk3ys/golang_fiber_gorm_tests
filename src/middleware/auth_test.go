package middleware_test

import (
	"api/src/constants"
	"api/src/middleware"
	"api/src/models"
	"api/src/views"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestParseBearerToken(t *testing.T) {
	parsed := middleware.ParseBearerToken(" bearer 123456 ")
	assert.Equal(t, "123456", parsed)
}

func createLoggedInUser() (models.User, string) {
	newUser := models.User{
		FirstName: fake.Person().FirstName(),
		LastName:  fake.Person().LastName(),
		Email:     fake.Internet().Email(),
		Password:  "hash",
		Status:    1,
	}
	db.Create(&newUser)

	bearer, _ := mware.Repo.UpdateUserSession("", newUser.ID)
	return newUser, bearer
}

func TestAuthTokenValidParseUserIdSuccess(t *testing.T) {
	setupTest()
	defer tearDownTest()

	user, bearer := createLoggedInUser()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(constants.AUTH_HEADER, "bearer "+bearer)

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	// t.Logf("uuid %v | %s", user.ID, string(body))

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, user.ID.String(), string(body))
}

func TestAuthTokenInvalidSessionValid(t *testing.T) {
	setupTest()
	defer tearDownTest()

	user, _ := createLoggedInUser()
	claims := views.SessionClaims{
		UserID: user.ID,
		// UserRole: userRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-constants.AUTH_TOKEN_DURATION).Unix(),
		},
	}
	bearer, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv(constants.AUTH_SECRET)))
	// Update session token to exipred one
	var exSession models.Session
	db.First(&exSession)
	exSession.Token = bearer
	db.Save(&exSession)

	// ACT
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(constants.AUTH_HEADER, "bearer "+bearer)

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	// t.Logf("uuid %v | %s", user.ID, string(body))

	// ASSERT
	// t.Logf("bearer %s", resp.Header.Get(constants.AUTH_HEADER))
	newBearer := resp.Header.Get(constants.AUTH_HEADER)
	assert.NotEmpty(t, newBearer)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, user.ID.String(), string(body))

	// Check session updated in db
	var updatedSession models.Session
	db.First(&updatedSession)
	assert.Equal(t, middleware.ParseBearerToken(newBearer), updatedSession.Token)
	assert.Equal(t, bearer, updatedSession.FromToken)
}

func TestAuthTokenInvalidSessionValidAndLongRequest(t *testing.T) {
	setupTest()
	defer tearDownTest()

	user, _ := createLoggedInUser()

	claims := views.SessionClaims{
		UserID: user.ID,
		// UserRole: userRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-constants.AUTH_TOKEN_DURATION).Unix(),
		},
	}
	bearer, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv(constants.AUTH_SECRET)))
	// Update session token to exipred one
	var exSession models.Session
	db.First(&exSession)
	exSession.Token = bearer
	db.Save(&exSession)

	// Set duration so terst doesnt need to wait long
	constants.AUTH_TOKEN_DURATION = time.Second * 2

	// ACT
	req := httptest.NewRequest(http.MethodGet, "/long", nil)
	req.Header.Set(constants.AUTH_HEADER, "bearer "+bearer)
	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)

	// ASSERT
	// t.Logf("bearer %s", resp.Header.Get(constants.AUTH_HEADER))
	newBearer := resp.Header.Get(constants.AUTH_HEADER)
	assert.NotEmpty(t, newBearer)
	assert.NotEmpty(t, resp.Header.Get("duration"))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, user.ID.String(), string(body))

	// Check session updated in db
	var updatedSession models.Session
	db.First(&updatedSession)
	assert.Equal(t, middleware.ParseBearerToken(newBearer), updatedSession.Token)
	assert.NotEqual(t, bearer, updatedSession.FromToken) // A new token should be created during the request -> bearer -> lost bearer due to duration -> new bearer

	constants.AUTH_TOKEN_DURATION = time.Second * 20
}
