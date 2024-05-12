package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Session struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Token     string    `gorm:"index"`
	FromToken string    `gorm:"index"`
	UserID    uuid.UUID `gorm:"index"`
	CreatedAt time.Time
	ExpiresAt time.Time
}
