package middleware_test

import (
	"api/src/constants"
	"api/src/models"
	"api/src/token"
	"api/src/util"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseBearerToken(t *testing.T) {
	parsed := token.ParseBearerToken(" bearer 123456 ")
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

	bearer, _, _, _ := mware.Repo.CreateUserSession(newUser.ID)
	return newUser, bearer
}

func TestAuthTokenValidParseUserIdSuccess(t *testing.T) {
	setupTest()
	defer tearDownTest()

	user, bearer := createLoggedInUser()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(constants.ACCESS_TOKEN_HEADER, "bearer "+bearer)

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	// t.Logf("uuid %v | %s", user.ID, string(body))

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, user.ID.String(), string(body))
}

func TestAuthTokenNotSentFail(t *testing.T) {
	// ACT
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(constants.ACCESS_TOKEN_HEADER, "")
	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := util.JsonMapFromBytes(body)

	// ASSERT
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Equal(t, "No token!", json["error"])
}

func TestAuthTokenTimedoutFail(t *testing.T) {
	setupTest()
	defer tearDownTest()
	// Reset token duration after test runs
	defer func() {
		dur, _ := util.ParseDuration(os.Getenv("AUTH_TOKEN_DURATION"))
		constants.ACCESS_TOKEN_DURATION = dur
	}()

	constants.ACCESS_TOKEN_DURATION = time.Millisecond
	_, bearer := createLoggedInUser()

	time.Sleep(time.Second)

	// ACT
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(constants.ACCESS_TOKEN_HEADER, "bearer "+bearer)
	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := util.JsonMapFromBytes(body)
	// t.Logf("uuid %v | %s", user.ID, string(body))

	// ASSERT
	// t.Logf("bearer %s", resp.Header.Get(constants.AUTH_HEADER))
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Equal(t, "Token invalid!", json["error"])
}
