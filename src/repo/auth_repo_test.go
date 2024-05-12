package repos_test

import (
	otp "api/src/enums/otps"
	"api/src/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkTestOtpCodeValid(t *testing.T, code uint16) {
	if code >= 10000 {
		t.Fatalf("Should not generate number larger than 10000 %v", code)
	}
	if code < 1000 {
		t.Fatalf("Should not generate number smaller than 1000 %v", code)
	}
}

func TestCreateUserOTPCodeSuccessSignup(t *testing.T) {
	setupTest()
	defer tearDownTest()

	email := fake.Person().Contact().Email

	code, err := repo.CreateUserOTPCode(email, "", otp.SIGNUP_EMAIL)
	if err != nil {
		t.Fatalf("Error creating OTP code %v", err)
	}

	checkTestOtpCodeValid(t, code)

	var otpCode models.UserOTPCode
	db.First(&otpCode, "email=?", email)
	assert.Equal(t, code, otpCode.Code)
	// TODO: check create and expired duration
}
