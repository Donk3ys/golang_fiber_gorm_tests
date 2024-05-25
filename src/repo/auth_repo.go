package repos

import (
	"api/src/constants"
	otp "api/src/enums/otps"
	"api/src/models"
	"api/src/util"
	"api/src/views"

	"api/src/token"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	uuid "github.com/satori/go.uuid"
)

func (i *Instance) CreateUserOTPCode(email string, mobile string, otpType string) (uint16, error) {
	if !otp.IsValid(otpType) {
		return 0, fmt.Errorf("Type %s does not exist!", otpType)
	}
	duration := otp.GetDuration(otpType)

	// Create verifiaction userOTP
	userOTP := models.UserOTPCode{
		Code:      util.GenerateCode(),
		Email:     email,
		Mobile:    mobile,
		Type:      otpType,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(duration),
	}
	res := i.Db.Create(&userOTP)
	if res.Error != nil {
		return 0, errors.New("Unexpected error generating your verification code!")
	}

	return userOTP.Code, nil
}

func (i *Instance) DeleteUserAuthOtpCodes(email string) error {
	res := i.Db.Where("email=? AND type IN (?, ?, ?, ?)", email, otp.LOGIN, otp.REGISTER_MOBILE, otp.SIGNUP_EMAIL, otp.PASSWORD_RESET).Delete(&models.UserOTPCode{})
	return res.Error
}

func (i *Instance) CheckPasswordResetCode(email string, code string) error {
	var resetOTP models.UserOTPCode
	i.Db.Order("expires_at DESC").First(&resetOTP, "email=? AND code=? AND type=? AND expires_at>?", email, code, otp.PASSWORD_RESET, time.Now().UTC())
	if resetOTP.Email == "" {
		return errors.New("Code not found or expired for reset!")
	}
	return nil
}

func (i *Instance) fetchExisingUserSession(token string, userID uuid.UUID) *models.Session {
	// Check for existing session in cache

	// Check for existing session in db/cache
	var exSession models.Session
	i.Db.First(&exSession, "(token=? AND expires_at>?) OR from_token=?", token, time.Now(), token)
	if exSession.UserID == uuid.Nil {
		return nil
	}

	return &exSession
}

func (i *Instance) CreateUserSession(userID uuid.UUID) (string, int64, string, error) {
	// Create new bearerToken
	bearer, bearerExpiry, err := token.CreateAuthToken(userID)
	if err != nil {
		return "", 0, "", err
	}

	// Check for existing session in db
	var exSession models.Session
	i.Db.First(&exSession).Where("user_id", userID)
	if exSession.UserID != uuid.Nil {
		exRefreshToken, err := token.DecodeRefreshToken(exSession.Token)
		// Check refresh isnt about to expire
		dur := exSession.ExpiresAt.Sub(time.Now().UTC())
		notExpiring := dur > constants.ACCESS_TOKEN_DURATION

		// If refresh token still valid and not expiring return exisitng one
		if err == nil && exRefreshToken.Valid && notExpiring {
			return bearer, bearerExpiry.Unix(), exSession.Token, nil
		}
	}

	// Create new session
	// Create new refreshToken
	refreshToken, refreshExpiry, err := token.CreateRefreshToken(userID)
	if err != nil {
		return "", 0, "", err
	}

	// Save new session to db/cache
	session := models.Session{
		UserID:    userID,
		Token:     refreshToken,
		FromToken: exSession.Token,
		CreatedAt: time.Now(),
		ExpiresAt: refreshExpiry.UTC(),
	}
	res := i.Db.Save(&session)
	if res.RowsAffected == 0 {
		return "", 0, "", errors.New("Could not create auth token!")
	}

	return bearer, bearerExpiry.Unix(), refreshToken, nil
}

func (i *Instance) UpdateUserSession(existingBearerToken string, existingRefreshToken string) (string, int64, string, error) {
	exBearerToken, _ := token.DecodeAuthToken(existingBearerToken)
	exRefreshToken, errRT := token.DecodeRefreshToken(existingRefreshToken)
	userID := exBearerToken.Claims.(*views.SessionClaims).UserID
	userIDRT := exRefreshToken.Claims.(*views.SessionClaims).UserID
	if userID != userIDRT {
		return "", 0, "", errors.New("User missmatch from tokens")
	}

	log.Debug("")
	// Check for existing session in db/cache
	var exSession models.Session
	i.Db.First(&exSession, "(token=? AND expires_at>?) OR from_token=?", existingRefreshToken, time.Now().UTC(), existingRefreshToken)
	if exSession.UserID == uuid.Nil {
		return "", 0, "", errors.New("Could not find exising session!")
	}

	// Create new bearerToken
	bearer, bearerExpiry, err := token.CreateAuthToken(userID)
	if err != nil {
		return "", 0, "", err
	}

	// Check refresh isnt about to expire
	dur := exSession.ExpiresAt.Sub(time.Now().UTC())
	validDur := dur > constants.ACCESS_TOKEN_DURATION
	if errRT == nil && exRefreshToken.Valid && validDur {
		// refreshToken vaild send new bearer
		return bearer, bearerExpiry.Unix(), "", nil
	}

	// // Check if refresh only expired
	// ve, ok := errRT.(*jwt.ValidationError)
	// if !ok || ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) == 0 {
	// 	log.Debug("Reresh error not expired")
	// 	return "", 0, "", errors.New("Token invaild [001]")
	// }

	// Create new refreshToken
	refreshToken, refreshExpiry, err := token.CreateRefreshToken(userID)
	if err != nil {
		log.Error("Refresh token create ", err)
		return "", 0, "", err
	}

	// Create session
	session := models.Session{
		UserID:    userID,
		Token:     refreshToken,
		FromToken: existingRefreshToken,
		CreatedAt: time.Now(),
		ExpiresAt: refreshExpiry.UTC(),
	}
	res := i.Db.Save(&session)
	if res.RowsAffected == 0 {
		return "", 0, "", errors.New("Could not save new session")
	}

	return bearer, bearerExpiry.Unix(), refreshToken, nil
}
