package storage

import (
	"api/src/models"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var tables = []interface{}{
	models.Session{},
	models.User{},
	models.UserOTPCode{},
}

func ConnectPostgres() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOSTNAME"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Postgres connection error", err)
	}

	return db
}

func AutoMigratePostgres(db *gorm.DB) {
	db.AutoMigrate(
		tables...,
	)
}
