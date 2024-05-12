package handlers_test

import (
	"api/src/constants"
	otp "api/src/enums/otps"
	"api/src/handlers"
	"api/src/models"
	"api/src/util"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// WARNING: Test Healpers
func signupTestUser() (models.User, *http.Response) {
	sendMailCall := mockEmail.On("SendMail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	defer sendMailCall.Unset()

	pw := fake.Internet().Password()
	newUser := models.User{
		FirstName: fake.Person().FirstName(),
		LastName:  fake.Person().LastName(),
		Email:     fake.Internet().Email(),
		Password:  pw,
	}
	reqBody := fmt.Sprintf(`{"email":"%s","password":"%s","first_name":"%s","last_name":"%s"}`, newUser.Email, newUser.Password, newUser.FirstName, newUser.LastName)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sign-up", strings.NewReader(reqBody))
	req.Header.Set("Content-type", "application/json")

	resp, _ := app.Test(req, -1)

	db.First(&newUser, "email=? AND status<2", newUser.Email)
	newUser.Password = pw
	return newUser, resp
}

func newLoggedInTestUser() (models.User, string) {
	user, _ := signupTestUser()
	db.Model(&user).Update("status", 1)

	reqBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, user.Email, user.Password)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(reqBody))
	req.Header.Set("Content-type", "application/json")

	resp, _ := app.Test(req, -1)
	bearer := resp.Header.Get(constants.AUTH_HEADER)

	return user, bearer
}

func createUserPasswordResetRequest(user *models.User) (uint16, *http.Response) {
	sendMailCall := mockEmail.On("SendMail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	defer sendMailCall.Unset()
	url := fmt.Sprintf("/api/v1/password/reset/request?email=%s", user.Email)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	resp, _ := app.Test(req, -1)

	var userOTP models.UserOTPCode
	db.First(&userOTP, "email=? AND type=?", user.Email, otp.PASSWORD_RESET)

	return userOTP.Code, resp
}

// INFO: Tests Success
func TestSignupUserSuccess(t *testing.T) {
	setupTest()
	defer tearDownTest()

	newUser, resp := signupTestUser()

	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)
	t.Log(json)

	if assert.Equal(t, http.StatusCreated, resp.StatusCode) {
		assert.Equal(t, "User successfully created. Please verify your email address.", json["message"])

		var user models.User
		db.First(&user, "email=? AND status=0", newUser.Email)
		assert.NotEmpty(t, user.ID)

		// Get user code from db
		var count int64
		db.Model(&models.UserOTPCode{}).Count(&count)
		t.Log(count)
		assert.Equal(t, int64(1), count)

		var otpCode models.UserOTPCode
		db.First(&otpCode, "email=?", newUser.Email)
		assert.Equal(t, newUser.Email, otpCode.Email)
		assert.NotEmpty(t, otpCode.ExpiresAt)
	}
}

func TestVerifyUserSuccess(t *testing.T) {
	setupTest()
	defer tearDownTest()

	// Signup user - Tested: TestSignupUserSuccess
	user, _ := signupTestUser()

	// Get user code from db
	var otpCode models.UserOTPCode
	db.First(&otpCode, "email=?", user.Email)
	assert.Equal(t, user.Email, otpCode.Email)

	// TEST
	// Verify email
	url := fmt.Sprintf("/api/v1/sign-up/verify?email=%s&code=%d", user.Email, otpCode.Code)
	t.Log(url)
	req := httptest.NewRequest(http.MethodGet, url, nil)

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)
	t.Log(json)

	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		assert.Equal(t, "Email successfully verified!", json["message"])
		// assert.NotEmpty(t, json["token"])

		var userDb models.User
		db.First(&userDb, "email=? AND status=1", user.Email)
		assert.NotEmpty(t, userDb.ID)
	}
}

func TestLoginSuccess(t *testing.T) {
	setupTest()
	defer tearDownTest()

	user, _ := newLoggedInTestUser()

	reqBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, user.Email, user.Password)
	t.Log(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(reqBody))
	req.Header.Set("Content-type", "application/json")

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)
	bearer := resp.Header.Get(constants.AUTH_HEADER)

	if resp.StatusCode != 200 {
		t.Log(json["error"])
	}

	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		assert.Equal(t, user.Email, json["email"])
		assert.Nil(t, json["password"])
		assert.Equal(t, float64(1), json["status"])

		assert.NotEmpty(t, bearer)
	}
}

func TestGetUserSuccess(t *testing.T) {
	setupTest()
	defer tearDownTest()

	user, bearer := newLoggedInTestUser()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user", nil)
	req.Header.Set(constants.AUTH_HEADER, bearer)

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)

	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		assert.Equal(t, user.ID.String(), json["id"])
		assert.Equal(t, user.Email, json["email"])
		assert.Nil(t, json["password"])
		assert.Equal(t, float64(1), json["status"])
	}
}

func TestGetPasswordResetCodeSuccess(t *testing.T) {
	// ARRANGE
	setupTest()
	defer tearDownTest()

	user, _ := newLoggedInTestUser()

	// ACT
	code, resp := createUserPasswordResetRequest(&user)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)

	// ASSERT
	// Assert repsonse
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Email sent", json[handlers.MESSAGE])

	// Asset OTP created in db
	assert.NotEmpty(t, code)

	// Assert email was sent
	mockEmail.AssertCalled(t, "SendMail", mock.Anything, []string{user.Email}, "Password Reset Code", mock.Anything, mock.Anything)
	// mockMailer.AssertExpectations(t)
}
func TestVerifyPasswordResetCodeSuccess(t *testing.T) {
	// ARRANGE
	setupTest()
	defer tearDownTest()

	user, _ := newLoggedInTestUser()
	code, _ := createUserPasswordResetRequest(&user)

	// ACT
	url := fmt.Sprintf("/api/v1/password/reset/verify?email=%s&code=%v", user.Email, code)
	req := httptest.NewRequest(http.MethodGet, url, nil)

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)

	// ASSERT
	// Assert repsonse
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Reset code valid", json[handlers.MESSAGE])
}

func TestResetPasswordSuccess(t *testing.T) {
	// ARRANGE
	setupTest()
	defer tearDownTest()

	user, _ := newLoggedInTestUser()
	code, _ := createUserPasswordResetRequest(&user)

	sendMailCall := mockEmail.On("SendMail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	defer sendMailCall.Unset()

	// ACT
	newPassword := "my_new_password"
	reqBody := fmt.Sprintf(`{"email":"%s","password":"%s","code":%d}`, user.Email, newPassword, code)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/password/reset", strings.NewReader(reqBody))
	req.Header.Set("Content-type", "application/json")

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)

	// ASSERT
	// Assert repsonse
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Password successfully updated", json[handlers.MESSAGE])

	// New password is set and hashed
	db.First(&user)
	err := util.CheckPasswordMatch(user.Password, newPassword)
	assert.NoError(t, err)

	// Check otps deleted
	var count int64
	db.Debug().Model(&models.UserOTPCode{}).Where("email=? AND type=?", user.Email, otp.PASSWORD_RESET).Count(&count)
	assert.Equal(t, int64(0), count)

	// Check user sessions removed
	db.Model(&models.Session{}).Count(&count)
	assert.Equal(t, int64(0), count)

	// Check email sent
	mockEmail.AssertCalled(t, "SendMail", mock.Anything, []string{user.Email}, "Password Reset Success", mock.Anything, mock.Anything)
	mockEmail.AssertExpectations(t)
}

// INFO: Tests Fails
func TestSignupUserFailInvaildData(t *testing.T) {
	setupTest()
	defer tearDownTest()

	reqBody := fmt.Sprintf(``)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sign-up", strings.NewReader(reqBody))
	req.Header.Set("Content-type", "application/json")

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)

	if assert.Equal(t, http.StatusBadRequest, resp.StatusCode) {
		assert.Equal(t, "Invalid data!", json[handlers.ERROR])
	}
}

func TestSignupUserFailEmailOrPasswordNotSet(t *testing.T) {
	setupTest()
	defer tearDownTest()

	// Email
	reqBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, "", "pw")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sign-up", strings.NewReader(reqBody))
	req.Header.Set("Content-type", "application/json")

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)

	if assert.Equal(t, http.StatusBadRequest, resp.StatusCode) {
		assert.Equal(t, "Email or password not set", json[handlers.ERROR])
	}

	// Password
	reqBody = fmt.Sprintf(`{"email":"%s","password":"%s"}`, "test@email.com", "")
	req = httptest.NewRequest(http.MethodPost, "/api/v1/sign-up", strings.NewReader(reqBody))
	req.Header.Set("Content-type", "application/json")

	resp, _ = app.Test(req, -1)
	body, _ = ioutil.ReadAll(resp.Body)
	json = responseMap(body)

	if assert.Equal(t, http.StatusBadRequest, resp.StatusCode) {
		assert.Equal(t, "Email or password not set", json[handlers.ERROR])
	}
}

func TestSignupUserFailEmailAlreadyRegistered(t *testing.T) {
	setupTest()
	defer tearDownTest()

	user, _ := newLoggedInTestUser()

	reqBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, user.Email, user.Password)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sign-up", strings.NewReader(reqBody))
	req.Header.Set("Content-type", "application/json")

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)

	if assert.Equal(t, http.StatusConflict, resp.StatusCode) {
		assert.Equal(t, "Email address has already been registered", json[handlers.ERROR])
	}
}

func TestSignupUserFailAccountRemoved(t *testing.T) {
	setupTest()
	defer tearDownTest()

	user, _ := newLoggedInTestUser()
	db.Model(&user).Update("status", 2)

	reqBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, user.Email, user.Password)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sign-up", strings.NewReader(reqBody))
	req.Header.Set("Content-type", "application/json")

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)

	if assert.Equal(t, http.StatusConflict, resp.StatusCode) {
		assert.Equal(t, "Previous account has been removed. Please contact our support team to assist.", json[handlers.ERROR])
	}
}

func TestGetUserFailNotFound(t *testing.T) {
	// Arrange
	setupTest()
	defer tearDownTest()

	user, bearer := newLoggedInTestUser()
	// Set user to disabled
	user.Status = 0
	db.Save(&user)

	// Action
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user", nil)
	req.Header.Set(constants.AUTH_HEADER, bearer)

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)

	// Assert
	if assert.Equal(t, http.StatusNotFound, resp.StatusCode) {
		assert.Equal(t, "User not found or disabled. Please conteact support!", json[handlers.ERROR])
	}
}

func TestGetPasswordCodeFailNoEmailSet(t *testing.T) {
	// Arrange
	// setupTest()
	// defer tearDownTest()

	// Action
	url := fmt.Sprintf("/api/v1/password/reset/request?email=%s", "")
	req := httptest.NewRequest(http.MethodGet, url, nil)

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)

	// Assert
	if assert.Equal(t, http.StatusBadRequest, resp.StatusCode) {
		assert.Equal(t, "No email address provided!", json[handlers.ERROR])
	}
}

func TestGetPasswordCodeFailNoActiveUser(t *testing.T) {
	// Arrange
	setupTest()
	defer tearDownTest()

	user, _ := signupTestUser()

	// Action
	url := fmt.Sprintf("/api/v1/password/reset/request?email=%s", user.Email)
	req := httptest.NewRequest(http.MethodGet, url, nil)

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)

	// Assert
	if assert.Equal(t, http.StatusBadRequest, resp.StatusCode) {
		assert.Equal(t, "Email or password not set", json[handlers.ERROR])
	}
}

func TestGetPasswordCodeFailSendMail(t *testing.T) {
	// Arrange
	setupTest()
	defer tearDownTest()

	user, _ := newLoggedInTestUser()

	sendMailCall := mockEmail.On("SendMail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Mock Send Mail Error"))
	defer sendMailCall.Unset()

	// Action
	url := fmt.Sprintf("/api/v1/password/reset/request?email=%s", user.Email)
	req := httptest.NewRequest(http.MethodGet, url, nil)

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)

	// Assert
	if assert.Equal(t, http.StatusInternalServerError, resp.StatusCode) {
		assert.Equal(t, "Error sending OTP email! Please contact support.", json[handlers.ERROR])
	}

	// Check mail was attempted to send
	mockEmail.AssertExpectations(t)
}

func TestVerifyPasswordResetCodeFailNoEmailOrCode(t *testing.T) {
	// Arrange
	// setupTest()
	// defer tearDownTest()

	// Action
	req := httptest.NewRequest(http.MethodGet, "/api/v1/password/reset/verify?email=mail@email.com=&code=", nil)
	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "No email address or code provided!", json[handlers.ERROR])

	// Action
	req = httptest.NewRequest(http.MethodGet, "/api/v1/password/reset/verify?email=&code=00000", nil)
	resp, _ = app.Test(req, -1)
	body, _ = ioutil.ReadAll(resp.Body)
	json = responseMap(body)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "No email address or code provided!", json[handlers.ERROR])
}

func TestVerifyPasswordResetCodeFailNoCodeFoundOrExpired(t *testing.T) {
	// Arrange
	setupTest()
	defer tearDownTest()

	// Action
	req := httptest.NewRequest(http.MethodGet, "/api/v1/password/reset/verify?email=mail@email.com=&code=00000000", nil)
	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)

	// Assert
	if assert.Equal(t, http.StatusNotFound, resp.StatusCode) {
		assert.Equal(t, "Code not found or expired for reset!", json[handlers.ERROR])
	}
}

func TestResetPasswordFailNoBadData(t *testing.T) {
	// Arrange
	// setupTest()
	// defer tearDownTest()

	// Action
	req := httptest.NewRequest(http.MethodPost, "/api/v1/password/reset", nil)

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)

	// Assert
	if assert.Equal(t, http.StatusBadRequest, resp.StatusCode) {
		assert.Equal(t, "Invalid data!", json[handlers.ERROR])
	}
}

func TestResetPasswordFailPasswordNotLongEnough(t *testing.T) {
	// Arrange
	// setupTest()
	// defer tearDownTest()

	// Action
	reqBody := fmt.Sprintf(`{"email":"%s","password":"%s","code":%d}`, "mail@email.com", "12345", 1000)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/password/reset", strings.NewReader(reqBody))
	req.Header.Set("Content-type", "application/json")

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)

	// Assert
	if assert.Equal(t, http.StatusBadRequest, resp.StatusCode) {
		assert.Equal(t, "Password not 6 characters long!", json[handlers.ERROR])
	}
}

func TestResetPasswordFailUserNotFound(t *testing.T) {
	// Arrange
	setupTest()
	defer tearDownTest()

	// Action
	reqBody := fmt.Sprintf(`{"email":"%s","password":"%s","code":%d}`, "mail@email.com", "123456", 1000)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/password/reset", strings.NewReader(reqBody))
	req.Header.Set("Content-type", "application/json")

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := responseMap(body)

	// Assert
	if assert.Equal(t, http.StatusNotFound, resp.StatusCode) {
		assert.Equal(t, "User not found or disabled!", json[handlers.ERROR])
	}
}
