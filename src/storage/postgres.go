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
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Fatalf("Postgres connection error db %s at %s:%s\n", dbname, host, port, err)
	}
	log.Infof("Postgres connected to db %s at %s:%s", dbname, host, port)

	return db
}

func AutoMigratePostgres(db *gorm.DB) {
	db.AutoMigrate(
		tables...,
	)
}
