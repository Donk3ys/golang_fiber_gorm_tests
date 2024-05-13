package middleware_test

import (
	"api/src/constants"
	"api/src/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestAuthenticateAuthTokenValidateParseUserIdSuccess(t *testing.T) {
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
