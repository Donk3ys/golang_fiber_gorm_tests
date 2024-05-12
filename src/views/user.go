package views

import (
	"api/src/models"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type UserSlim struct {
	ID              uuid.UUID `json:"id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	MobileNumber    string    `json:"mobile_number"`
	ProfileImageUrl string    `json:"profile_image_url"`
	Status          uint8     `json:"status"` // 0: deactivated, 1: activated, 2: removed
}

type User struct {
	UserSlim
	Email     string    `json:"email" `
	CreatedAt time.Time `json:"created_at"`
}

type UserOTPCode struct {
	Code      uint16    `json:"code"`
	Email     string    `json:"email"`
	Mobile    string    `json:"mobile"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"-"`
	ExpiresAt time.Time
}

// Requests
type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PasswordResetReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     uint16 `json:"code"`
}

// Projections
func UserViewToUserModel(user *User) *models.User {
	return &models.User{
		ID:              user.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		MobileNumber:    user.MobileNumber,
		ProfileImageUrl: user.ProfileImageUrl,
		Status:          user.Status,
		Model: gorm.Model{
			CreatedAt: user.CreatedAt,
		},
	}
}
func UserModelToUserView(user *models.User) *User {
	return &User{
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UserSlim: UserSlim{
			ID:              user.ID,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			MobileNumber:    user.MobileNumber,
			ProfileImageUrl: user.ProfileImageUrl,
			Status:          user.Status,
		},
	}
}
