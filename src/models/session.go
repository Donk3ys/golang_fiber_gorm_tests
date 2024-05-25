package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Session struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Token     string    `gorm:"index"`
	FromToken string    `gorm:"index"`
	CreatedAt time.Time
	ExpiresAt time.Time
}
