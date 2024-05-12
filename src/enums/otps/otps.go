package otp

import (
	"api/src/util"
	"time"
)

const LOGIN = "login"
const PASSWORD_RESET = "password-reset"
const REGISTER_MOBILE = "register-mobile"
const SIGNUP_EMAIL = "signup-email"
const UPDATE_MOBILE = "update-mobile"

var ALL_OTP_TYPES = []string{LOGIN, PASSWORD_RESET, REGISTER_MOBILE, SIGNUP_EMAIL, UPDATE_MOBILE}

const SIGNUP_EMAIL_OTP_DURATION = time.Hour * 24 * 5 // 5 days
const REGISTER_MOBILE_OTP_DURATION = time.Hour * 12  // 12 hours
const DEFAULT_OTP_DURATION = time.Minute * 15        // 15 min

func IsValid(otp string) bool {
	return util.IndexOfItemInSlice(otp, ALL_OTP_TYPES) != -1
}

func GetDuration(otpType string) time.Duration {
	var duration time.Duration
	switch otpType {
	case REGISTER_MOBILE, UPDATE_MOBILE:
		duration = REGISTER_MOBILE_OTP_DURATION
		break
	case SIGNUP_EMAIL:
		duration = SIGNUP_EMAIL_OTP_DURATION
		break
	default:
		duration = DEFAULT_OTP_DURATION
		break
	}
	return duration
}
