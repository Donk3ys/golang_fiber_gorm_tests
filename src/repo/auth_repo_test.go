package repos_test

import (
	"api/src/constants"
	otp "api/src/enums/otps"
	"api/src/models"
	"testing"
	"time"

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

// TODO: TEST FAILED
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

func TestCreateUserSessionSuccess(t *testing.T) {
	setupTest()
	defer tearDownTest()

	now := time.Now().Unix()
	userID := uuid.NewV4()
	bearer, bearerExp, refreshToken, err := repo.CreateUserSession(userID)
	assert.NotEmpty(t, bearer)
	assert.Less(t, now, bearerExp)
	assert.NotEmpty(t, refreshToken)
	assert.NoError(t, err)

	// Check session in db
	var exSession models.Session
	db.First(&exSession).Where("user_id", userID)
	assert.Equal(t, userID, exSession.UserID)
	assert.Equal(t, refreshToken, exSession.Token)
	assert.Empty(t, exSession.FromToken)
	assert.Less(t, bearerExp, exSession.ExpiresAt.Unix())

	// Check 1 session in db
	var count int64
	db.Model(&models.Session{}).Count(&count)
	assert.Equal(t, int64(1), count)

	// ADD another user

	userID2 := uuid.NewV4()
	bearer2, bearerExp2, refreshToken2, err := repo.CreateUserSession(userID2)

	// Check session in db
	var exSession2 models.Session
	db.First(&exSession2).Where("user_id", userID2)
	assert.Equal(t, userID2, exSession2.UserID)
	assert.Equal(t, refreshToken2, exSession2.Token)
	assert.Empty(t, exSession2.FromToken)
	assert.Less(t, bearerExp2, exSession2.ExpiresAt.Unix())

	assert.NotEqual(t, bearer, bearer2)

	db.Model(&models.Session{}).Count(&count)
	assert.Equal(t, int64(2), count)
}

func TestCreateUserSessionUsingExistingRefreshTokenSuccess(t *testing.T) {
	setupTest()
	defer tearDownTest()

	userID := uuid.NewV4()
	bearer, bearerExp, refreshToken, err := repo.CreateUserSession(userID)

	time.Sleep(time.Second)

	bearer2, bearerExp2, refreshToken2, err := repo.CreateUserSession(userID)
	assert.NotEqual(t, bearer, bearer2)
	assert.Less(t, bearerExp, bearerExp2)
	assert.Equal(t, refreshToken, refreshToken2)
	assert.NoError(t, err)

	// Check only session in db as session sould be update and not a new on created
	var count int64
	db.Model(&models.Session{}).Count(&count)
	assert.Equal(t, int64(1), count)

	// Make sure session hasnt changed in db
	var exSession models.Session
	db.First(&exSession).Where("user_id", userID)
	assert.Equal(t, userID, exSession.UserID)
	assert.Equal(t, refreshToken, exSession.Token)
	assert.Empty(t, exSession.FromToken)
	assert.Less(t, bearerExp, exSession.ExpiresAt.Unix())
}

func TestCreateUserSessionUpdatingRefreshTokenSuccess(t *testing.T) {
	setupTest()
	defer tearDownTest()

	constants.ACCESS_TOKEN_DURATION = time.Second * 2
	constants.REFRESH_TOKEN_DURATION = time.Second * 3

	userID := uuid.NewV4()

	start1 := time.Now()
	bearer, bearerExp, refreshToken, err := repo.CreateUserSession(userID)

	var exSession models.Session
	db.First(&exSession).Where("user_id", userID)
	assert.Equal(t, start1.Add(constants.ACCESS_TOKEN_DURATION).Unix(), bearerExp)
	assert.Equal(t, start1.Add(constants.REFRESH_TOKEN_DURATION).Unix(), exSession.ExpiresAt.Unix())
	// log.Debug("bearerExp:\t", bearerExp)
	// log.Debug("refreshExp1:\t", exSession.ExpiresAt.Unix())

	time.Sleep(time.Second * 2)

	start2 := time.Now()
	bearer2, bearerExp2, refreshToken2, err := repo.CreateUserSession(userID)
	assert.NotEqual(t, bearer, bearer2)
	assert.Less(t, bearerExp, bearerExp2)
	assert.NotEqual(t, refreshToken, refreshToken2)
	assert.NoError(t, err)

	// Check only session in db as session sould be update and not a new on created
	var count int64
	db.Model(&models.Session{}).Count(&count)
	assert.Equal(t, int64(1), count)

	// Make sure session has changed in db
	var exSession2 models.Session
	db.First(&exSession2).Where("user_id", userID)
	assert.Equal(t, userID, exSession2.UserID)
	assert.Equal(t, refreshToken2, exSession2.Token)
	assert.Equal(t, refreshToken, exSession2.FromToken)
	assert.NotEmpty(t, exSession2.FromToken)
	assert.Less(t, bearerExp2, exSession2.ExpiresAt.Unix())

	assert.Equal(t, start2.Add(constants.ACCESS_TOKEN_DURATION).Unix(), bearerExp2)
	assert.Equal(t, start2.Add(constants.REFRESH_TOKEN_DURATION).Unix(), exSession2.ExpiresAt.Unix())
	// log.Debug("bearerExp2:\t", bearerExp2)
	// log.Debug("refreshExp2:\t", exSession2.ExpiresAt.Unix())
}

// func TestCache(t *testing.T) {
// 	setupTest()
// 	defer tearDownTest()
//
// 	// session := models.Session{
// 	// 	ID:        uuid.NewV4(),
// 	// 	Token:     fake.RandomStringWithLength(18),
// 	// 	FromToken: fake.RandomStringWithLength(18),
// 	// 	UserID:    uuid.NewV4(),
// 	// 	CreatedAt: time.Now(),
// 	// 	ExpiresAt: time.Now().Add(constants.REFRESH_TOKEN_DURATION),
// 	// }
//
// 	ctx := context.Background()
//
// 	err := cache.Do(ctx, cache.B().Set().Key("key").Value("bob").Build()).Error()
// 	if err != nil {
// 		log.Error("Err Set: ", err)
// 	}
// 	val, err := cache.Do(ctx, cache.B().Get().Key("key").Build()).ToString()
// 	if err != nil {
// 		log.Error("Err Get: ", err)
// 	}
// 	log.Debugf("GET %s", val)
//
// 	// err := cache.Set(ctx, "session_1", session)
// 	// if err != nil {
// 	// 	log.Error("Err Set ", err)
// 	// }
// 	//
// 	// time.Sleep(time.Millisecond * 10)
// 	//
// 	// session2 := models.Session{}
// 	// _, err = cache.Get(ctx, "session_1", &session2)
// 	// if err != nil {
// 	// 	log.Error("Err Get", err)
// 	// }
//
// 	// log.Debug("User Session ", session2)
// }
