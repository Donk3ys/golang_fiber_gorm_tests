package handlers

import (
	"api/src/constants"
	otp "api/src/enums/otps"
	"api/src/models"
	"api/src/util"
	"api/src/views"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	uuid "github.com/satori/go.uuid"
)

func (i *Instance) signupUser(c *fiber.Ctx) error {
	var userView views.User
	var json struct {
		Password string `json:"password"`
	}
	errUser := c.BodyParser(&userView)
	errJson := c.BodyParser(&json)
	if errUser != nil || errJson != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ERROR: "Invalid data!"})
	}
	json.Password = strings.Trim(json.Password, " ")
	userView.Email = strings.Trim(strings.ToLower(userView.Email), " ")

	if userView.Email == "" || json.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ERROR: "Email or password not set"})
	}

	// Check user doesnt exist and has been veritfied
	var userFromDb models.User
	userFromDbRes := i.Repo.Db.First(&userFromDb, "email=?", userView.Email)

	if userFromDb.Status == 1 {
		// User already exists and is verified
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{ERROR: "Email address has already been registered"})
	} else if userFromDb.ID != uuid.Nil && userFromDb.Status == 2 {
		// User already exists and is removed
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{ERROR: "Previous account has been removed. Please contact our support team to assist."})
	}

	hashedPw, err := util.GenPasswordHash(json.Password)
	if err != nil {
		log.Error(err)

		if err.Error() == "Password to short" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ERROR: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ERROR: "Unexpected error hashing user password!"})
	}

	// User to be inserted into db
	userModel := views.UserViewToUserModel(&userView)
	userModel.Password = hashedPw
	userModel.Status = 0 // unverified
	userModel.CreatedAt = time.Now().UTC()

	// Check if user already exists then update db or create in db
	if userFromDbRes.RowsAffected != 0 {
		userModel.ID = userFromDb.ID
		res := i.Repo.Db.Save(&userModel)
		if res.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ERROR: "Unexpected error updating existing user!"})
		}

		// Remove existing sessions
		i.Repo.Db.Where("user_id=?", userModel.ID).Delete(&models.Session{})
	} else {
		res := i.Repo.Db.Create(&userModel)
		if res.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ERROR: "Unexpected error creating user!"})
		}
	}

	code, err := i.Repo.CreateUserOTPCode(userModel.Email, "", otp.SIGNUP_EMAIL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ERROR: err.Error()})
	}

	// Send mail
	err = i.Notification.Email.SendUserVerificationCode(userView.Email, code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ERROR: "Error sending signup verification email!"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{MESSAGE: "User successfully created. Please verify your email address."})
}

func (i *Instance) loginUser(c *fiber.Ctx) error {
	var req views.LoginReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ERROR: "Invalid data!"})
	}

	var userModel models.User
	i.Repo.Db.First(&userModel, "email=? AND status<2", req.Email)
	if userModel.ID == uuid.Nil {
		// log.Printf("No user found for email: %s\n", req.Email)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ERROR: "Invalid login credentials! (no user)"})
		// return jsonError(c, fiber.StatusBadRequest, "Invalid login credentials!")
	}

	err := util.CheckPasswordMatch(userModel.Password, req.Password)
	if err != nil {
		// log.Printf("Password missmatch for email: %s\n%s\n", req.Email, err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{ERROR: "Invalid login credentials! (password)"})
		// return jsonError(c, fiber.StatusBadRequest, "Invalid login credentials!")
	}

	// Check user is verified
	if userModel.Status == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{ERROR: "User not verified. Please check your mail for verification code!"})
	}

	// code, err := i.Repo.CreateUserOTPCode(user.Email, user.MobileNumber, otp.LOGIN)
	// if err != nil {
	// 	return jsonError(c, fiber.StatusInternalServerError, err.Error())
	// }
	// err = i.Notification.SMS.SendUserLoginVerificationCode(user.MobileNumber, user.Email, code)
	// if err != nil {
	// 	return jsonError(c, fiber.StatusInternalServerError, "Error sending code!")
	// }
	// return jsonMessage(c, "Mobile OTP sent.")

	bearerToken, bearerExpiry, refreshToken, err := i.Repo.CreateUserSession(userModel.ID)
	if err != nil {
		log.Error("CreateUserSession: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ERROR: "Invalid login credentials! (token)"})
		// return jsonError(c, fiber.StatusBadRequest, "Invalid login credentials!")
	}

	c.Set("Access-Control-Expose-Headers", "*") // Needed to show headers in web app
	c.Set(constants.ACCESS_TOKEN_HEADER, "bearer "+bearerToken)
	c.Set(constants.ACCESS_TOKEN_EXPIRY_HEADER, fmt.Sprint(bearerExpiry))
	c.Set(constants.REFRESH_TOKEN_HEADER, refreshToken)
	return c.JSON(views.UserModelToUserView(&userModel))
}

// func (i *Instance) verifyMobileLogin(c *fiber.Ctx) error {
// 	strCode := c.Query("code")
// 	email := c.Query("email")
// 	code, err := strconv.ParseUint(strCode, 10, 16)
// 	if err != nil {
// 		return jsonError(c, fiber.StatusBadRequest, "Invalid code!")
// 	}
//
// 	var otpCode models.UserOTPCode
// 	res := i.Repo.Db.First(&otpCode, "email=? AND code=? AND expires_at>? AND type=?", email, uint16(code), time.Now().UTC(), otp.LOGIN)
// 	if res.RowsAffected == 0 {
// 		return jsonError(c, fiber.StatusUnauthorized, "Code does not exist or has expired!")
// 	}
// 	// i.Repo.Db.Where("email=? AND type=?", email, otp.LOGIN).Delete(&models.UserOTPCode{})
//
// 	var user models.User
// 	i.Repo.Db.First(&user, "email=? AND status=1", email)
// 	if user.ID == uuid.Nil {
// 		return jsonError(c, fiber.StatusUnauthorized, "User does not exist or has not verified their email or mobile!")
// 	}
//
// 	token, err := i.Repo.UpdateUserSession("", user.ID, user.Role)
// 	if err != nil {
// 		return jsonError(c, fiber.StatusInternalServerError, "Invalid login credentials!")
// 	}
//
// 	c.Set("Access-Control-Expose-Headers", "*") // Needed to show headers in web app
// 	c.Set(constants.AUTH_HEADER, "bearer "+token)
// 	return c.JSON(&user)
// }

func (i *Instance) signupVerifyEmail(c *fiber.Ctx) error {
	strCode := c.Query("code")
	email := c.Query("email")
	code, err := strconv.ParseUint(strCode, 10, 16)
	if err != nil {
		return c.JSON(fiber.Map{MESSAGE: "Email successfully verified!"})
	}

	var userModel models.User
	i.Repo.Db.First(&userModel, "email=? AND status=1", email)
	if userModel.ID != uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{ERROR: "Code does not exist!"})
	}

	var otpCode models.UserOTPCode
	res := i.Repo.Db.First(&otpCode, "email=? AND code=?", email, uint16(code))

	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{ERROR: "Code does not exist!"})
	}

	// Update user status to verified
	i.Repo.Db.Model(models.User{}).Where("email=?", email).Update("status", 1)

	// Delete all codes shere email matches
	i.Repo.DeleteUserAuthOtpCodes(email)
	return c.JSON(fiber.Map{MESSAGE: "Email successfully verified!"})
}

// func (i *Instance) signUpLinkMobileNumber(c *fiber.Ctx) error {
// 	email := c.Query("email")
// 	mobile := c.Query("mobile")
//
// 	var user models.User
// 	i.Repo.Db.First(&user, "email=? AND mobile_number='' AND status=1", email)
// 	if user.ID == uuid.Nil {
// 		return jsonError(c, fiber.StatusBadRequest, "User does not exsit")
// 	}
//
// 	code, err := i.Repo.CreateUserOTPCode(email, mobile, otp.REGISTER_MOBILE)
// 	if err != nil {
// 		return jsonError(c, fiber.StatusInternalServerError, err.Error())
// 	}
//
// 	err = i.Notification.SMS.SendUserMobileVerificationCode(mobile, email, code)
// 	if err != nil {
// 		return jsonError(c, fiber.StatusInternalServerError, "Error sending code!")
// 	}
//
// 	return jsonMessage(c, "Code sent to mobile.")
// }
//
// func (i *Instance) updateUserMobileNumber(c *fiber.Ctx) error {
// 	uId := (c.Locals(constants.REQ_USER_ID)).(string)
// 	newMobile := c.Query("mobile")
//
// 	var user models.User
// 	i.Repo.Db.First(&user, "ID=?", uId)
// 	if user.ID == uuid.Nil {
// 		return jsonError(c, fiber.StatusUnauthorized, "No user found!")
// 	}
// 	i.Repo.DeleteUserAuthOtpCodes(user.Email)
//
// 	code, err := i.Repo.CreateUserOTPCode(user.Email, newMobile, otp.UPDATE_MOBILE)
// 	if err != nil {
// 		return jsonError(c, fiber.StatusInternalServerError, err.Error())
// 	}
//
// 	err = i.Notification.SMS.SendUserMobileVerificationCode(newMobile, user.Email, code)
// 	if err != nil {
// 		return jsonError(c, fiber.StatusInternalServerError, "Error sending code!")
// 	}
//
// 	return jsonMessage(c, "Code sent to mobile.")
// }
//
// func (i *Instance) verifyUserMobileLink(c *fiber.Ctx) error {
// 	strCode := c.Query("code")
// 	email := c.Query("email")
// 	mobile := c.Query("mobile")
// 	code, err := strconv.ParseUint(strCode, 10, 16)
// 	if err != nil {
// 		return jsonError(c, fiber.StatusBadRequest, "Invalid data!")
// 	}
//
// 	var otpCode models.UserOTPCode
// 	res := i.Repo.Db.First(&otpCode, "email=? AND mobile=? AND code=? AND expires_at>? AND (type=? OR type=?)", email, mobile, uint16(code), time.Now().UTC(), otp.REGISTER_MOBILE, otp.UPDATE_MOBILE)
// 	if res.RowsAffected == 0 {
// 		return jsonError(c, fiber.StatusUnauthorized, "Code does not exist or has expired!")
// 	}
//
// 	// Update user mobile number
// 	i.Repo.Db.Model(models.User{}).Where("email=?", email).Update("mobile_number", mobile)
//
// 	// Delete all codes where mobile number matches
// 	i.Repo.DeleteUserAuthOtpCodes(email)
// 	return jsonMessage(c, "Mobile successfully verified!")
// }

func (i *Instance) getUser(c *fiber.Ctx) error {
	uId := (c.Locals(constants.REQ_USER_ID)).(uuid.UUID)

	var userModel models.User
	i.Repo.Db.First(&userModel, "id=? AND status=?", uId, 1)

	if userModel.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{ERROR: "User not found or disabled. Please conteact support!"})
	}

	return c.JSON(views.UserModelToUserView(&userModel))
}

func (i *Instance) updateUser(c *fiber.Ctx) error {
	uId := (c.Locals(constants.REQ_USER_ID)).(string)

	var userView views.User
	if err := c.BodyParser(&userView); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ERROR: "Invalid data!"})
	}

	// Check user exists in db
	var userFromDb models.User
	exUserRes := i.Repo.Db.First(&userFromDb, "id=?", uId)
	if exUserRes.RowsAffected == 0 || userFromDb.Status != 1 || userView.Email != userFromDb.Email {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{ERROR: "Unauthorized to update this user"})
	}

	// var json map[string]any
	// if err := c.BodyParser(&json); err != nil {
	// 	return jsonError(c, fiber.StatusBadRequest, "Invalid update request!")
	// }
	// imgBytes := json["image_bytes"]
	// imgExt := json["image_extension"]
	// removeImg := json["remove_image"]
	//
	// // Save image to file system
	// if imgBytes != nil && imgExt != nil {
	// 	newUrl, err := i.Repo.Fs.UpdateOrCreateFile(
	// 		uId,
	// 		user.ProfileImageUrl,
	// 		imgBytes.(string),
	// 		imgExt.(string),
	// 		constants.USER_PROFILE_PICS_PATH,
	// 		util.GenerateCode(),
	// 	)
	// 	if err != nil {
	// 		return jsonError(c, fiber.StatusBadRequest, "Error storing user image")
	// 	}
	//
	// 	// 192.168.1.210:8000/public/profile/bc931248-e087-4368-9603-e74546583680--2222.png
	// 	user.ProfileImageUrl = *newUrl
	// } else if removeImg.(bool) && user.ProfileImageUrl != "" {
	// 	i.Repo.Fs.RemoveFile(user.ProfileImageUrl, constants.USER_PROFILE_PICS_PATH)
	// 	user.ProfileImageUrl = ""
	// }

	// updateUser := map[string]interface{}{
	// 	"first_name":        user.FirstName,
	// 	"last_name":         user.LastName,
	// 	"phone_number":      user.MobileNumber,
	// 	"profile_image_url": user.ProfileImageUrl,
	// }
	// res := i.Repo.Db.Model(&user).Clauses(clause.Returning{}).Where("id=?", uId).Updates(updateUser)
	// if res.Error != nil {
	// 	return jsonError(c, fiber.StatusInternalServerError, "Unexpected error trying to update user")
	// }

	return c.JSON(&userView)
}

// TODO: test
func (i *Instance) removeUser(c *fiber.Ctx) error {
	uId := (c.Locals(constants.REQ_USER_ID)).(string)

	var json map[string]string
	if err := c.BodyParser(&json); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ERROR: "Invalid data!"})
	}
	password := json["password"]

	var userModel models.User
	i.Repo.Db.First(&userModel, "id=? AND status>0", uId)

	err := util.CheckPasswordMatch(userModel.Password, password)
	if err != nil {
		// log.Printf("Password missmatch for email: %s\n", req.Email)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{ERROR: "Invalid password or user already removed!"})
	}

	res := i.Repo.Db.Model(&models.User{}).Where("id=?", uId).Update("status", 0)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ERROR: "Unexpected error trying to remove user"})
	}

	i.Repo.Db.Where("user_id=?", uId).Delete(&models.Session{})

	return c.JSON(fiber.Map{MESSAGE: "User succesfully removed."})
}

func (i *Instance) refreshSession(c *fiber.Ctx) error {
	existingBearerToken := c.Get(constants.ACCESS_TOKEN_HEADER)
	existingRefreshToken := c.Get(constants.REFRESH_TOKEN_HEADER)
	if existingBearerToken == "" || existingRefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ERROR: "Bad Request - No token missing"})
	}

	bearerToken, bearerExpiry, refreshToken, err := i.Repo.UpdateUserSession(existingBearerToken, existingRefreshToken)
	if err != nil {
		log.Error("UpdateUserSession: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ERROR: "Could not update auth token"})
		// return jsonError(c, fiber.StatusBadRequest, "Invalid login credentials!")
	}

	c.Set("Access-Control-Expose-Headers", "*") // Needed to show headers in web app
	c.Set(constants.ACCESS_TOKEN_HEADER, "bearer "+bearerToken)
	c.Set(constants.ACCESS_TOKEN_EXPIRY_HEADER, fmt.Sprint(bearerExpiry))
	if refreshToken != "" {
		c.Set(constants.REFRESH_TOKEN_HEADER, refreshToken)
	}
	return c.JSON(fiber.Map{MESSAGE: "Session updated"})
}

func (i *Instance) getPasswordResetCode(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ERROR: "No email address provided!"})
	}

	var userModel models.User
	i.Repo.Db.First(&userModel, "email=? AND status=1", email)
	if userModel.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ERROR: "Email or password not set"})
	}

	code, err := i.Repo.CreateUserOTPCode(userModel.Email, "", otp.PASSWORD_RESET)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ERROR: err.Error()})
	}

	err = i.Notification.Email.SendPasswordResetCode(userModel.Email, code)
	if err != nil {
		log.Errorf("getPasswordResetCode Send Mail error\n%s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ERROR: "Error sending OTP email! Please contact support."})
	}
	return c.JSON(fiber.Map{MESSAGE: "Email sent"})
}

func (i *Instance) verifyPasswordResetCode(c *fiber.Ctx) error {
	email := c.Query("email")
	code := c.Query("code")
	if email == "" || code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ERROR: "No email address or code provided!"})
	}

	err := i.Repo.CheckPasswordResetCode(email, code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{ERROR: err.Error()})
	}

	return c.JSON(fiber.Map{MESSAGE: "Reset code valid"})
}

func (i *Instance) resetPassword(c *fiber.Ctx) error {
	var resetReq views.PasswordResetReq
	if err := c.BodyParser(&resetReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ERROR: "Invalid data!"})
	}

	if len(resetReq.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ERROR: "Password not 6 characters long!"})
	}

	var userModel models.User
	i.Repo.Db.First(&userModel, "email=? AND status=1", resetReq.Email)
	if userModel.Email == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{ERROR: "User not found or disabled!"})
	}

	err := i.Repo.CheckPasswordResetCode(resetReq.Email, fmt.Sprintf("%d", resetReq.Code))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{ERROR: err.Error()})
	}

	passwordHash, err := util.GenPasswordHash(resetReq.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ERROR: "Unexpected server error [001]"})
	}

	userModel.Password = passwordHash
	res := i.Repo.Db.Save(&userModel)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ERROR: "Unexpected server error [002]"})
	}

	i.Repo.DeleteUserAuthOtpCodes(userModel.Email)
	i.Repo.Db.Where("user_id=?", userModel.ID).Delete(&models.Session{})

	i.Notification.Email.SendPasswordResetSuccess(userModel.Email)
	return c.JSON(fiber.Map{MESSAGE: "Password successfully updated"})
}

// Healpers (Maybe move to util/services package)
func (i *Instance) RemoveExpiredUserOTPCodes() {
	log.Info("Removing expired User OTP codes.")
	i.Repo.Db.Where("expires_at<?", time.Now().UTC()).Delete(&models.UserOTPCode{})
}

func (i *Instance) RemoveExpiredSessions() {
	log.Info("Removing expired sessions.")
	i.Repo.Db.Where("expires_at<?", time.Now().UTC()).Delete(&models.Session{})
}
