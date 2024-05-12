package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID              uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email           string    `gorm:"unique"`
	FirstName       string
	LastName        string
	MobileNumber    string
	ProfileImageUrl string
	Password        string
	Status          uint8 // 0: deactivated, 1: activated, 2: removed
}

type UserOTPCode struct {
	Code      uint16    `json:"code" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"primaryKey"`
	Mobile    string    `json:"mobile" gorm:"primaryKey"`
	Type      string    `json:"type" gorm:"primaryKey"`
	CreatedAt time.Time `json:"-"`
	ExpiresAt time.Time
}
