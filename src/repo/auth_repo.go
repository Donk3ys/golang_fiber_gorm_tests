package repos

import (
	"api/src/constants"
	otp "api/src/enums/otps"
	"api/src/models"
	"api/src/util"
	"api/src/views"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
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

func (i *Instance) UpdateUserSession(existingToken string, userID uuid.UUID) (string, error) {
	claims := views.SessionClaims{
		UserID: userID,
		// UserRole: userRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(constants.AUTH_TOKEN_DURATION).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv(constants.AUTH_SECRET)))
	if err != nil {
		return token, err
	}

	// Create session
	session := models.Session{
		Token:     token,
		FromToken: existingToken,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(constants.AUTH_SESSION_DURATION),
	}
	if existingToken != "" {
		i.Db.Where("token=? AND user_id=?", existingToken, userID).Delete(models.Session{})
	}
	i.Db.Create(&session)

	return token, err
}
