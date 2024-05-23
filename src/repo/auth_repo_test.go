package repos_test

import (
	"api/src/constants"
	otp "api/src/enums/otps"
	"api/src/models"
	"context"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2/log"
	uuid "github.com/satori/go.uuid"
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

func TestCacheTest(t *testing.T) {
	setupTest()
	defer tearDownTest()

	session := models.Session{
		ID:        uuid.NewV4(),
		Token:     fake.RandomStringWithLength(18),
		FromToken: fake.RandomStringWithLength(18),
		UserID:    uuid.NewV4(),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(constants.REFRESH_TOKEN_DURATION),
	}

	ctx := context.Background()
	err := cache.Set(ctx, "session_1", session)
	if err != nil {
		log.Error("Err Set ", err)
	}

	time.Sleep(time.Millisecond * 10)

	session2 := models.Session{}
	_, err = cache.Get(ctx, "session_1", &session2)
	if err != nil {
		log.Error("Err Get", err)
	}

	log.Debug("User Session ", session2)
}
