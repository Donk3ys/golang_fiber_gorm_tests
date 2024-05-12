package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestConnectPostgres(ctx context.Context) (*gorm.DB, testcontainers.Container) {
	pw := os.Getenv("POSTGRES_PASSWORD")
	usr := os.Getenv("POSTGRES_USER")
	dbName := os.Getenv("POSTGRES_DB")

	var env = map[string]string{
		"POSTGRES_PASSWORD": pw,
		"POSTGRES_USER":     usr,
		"POSTGRES_DB":       dbName,
	}
	var port = "5432/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:14-alpine",
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	}
	tc, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		log.Panicf("failed to start container: %v", err)
	}

	p, err := tc.MappedPort(ctx, "5432")
	if err != nil {
		log.Panicf("failed to get container external port: %v", err)
	}

	log.Println("postgres container ready and running at port: ", p.Port())

	time.Sleep(time.Millisecond * 400)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOSTNAME"),
		p.Port(),
		usr,
		pw,
		dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic("Test postgres connection error")
	}

	return db, tc
}

func TestDropTablesPostgres(db *gorm.DB) {
	db.Migrator().DropTable(tables...)
}
