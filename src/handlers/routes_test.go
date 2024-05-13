package handlers_test

import (
	"api/src/handlers"
	"api/src/middleware"
	"api/src/mocks"
	"api/src/notification"
	"api/src/repo"
	"api/src/storage"
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jaswdr/faker"
	"github.com/joho/godotenv"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
)

var (
	app       *fiber.App
	db        *gorm.DB
	dbTc      testcontainers.Container
	fake      faker.Faker
	mockEmail *mocks_test.EmailCLient
	mockSms   *mocks_test.SmsCLient
	mockFs    *mocks_test.FileSysetmClient
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env-test"); err != nil {
		log.Panic("Test environment variables not set or error parsing!")
	}

	fake = faker.New()
	setupTestApp()
	code := m.Run()
	tearDownTestApp()
	os.Exit(code)
}

func setupTestApp() {
	// createFolders()

	db, dbTc = storage.TestConnectPostgres(context.Background())
	storage.AutoMigratePostgres(db)

	// database.Seed(db)

	mockEmail = &mocks_test.EmailCLient{}
	mockSms = &mocks_test.SmsCLient{}
	mockFs = &mocks_test.FileSysetmClient{}

	repo := repos.Instance{
		Db: db,
		Fs: storage.FileSystem{
			Client: mockFs,
		},
	}

	notification := notification.Instance{
		Email: notification.Email{
			Client: mockEmail,
		},
		SMS: notification.SMS{
			Client: mockSms,
		},
	}

	appHandlers := handlers.Instance{
		Middleware:   middleware.Instance{Repo: repo},
		Notification: notification,
		Repo:         repo,
	}

	app = fiber.New(fiber.Config{
		Immutable: true,
	})
	// app.Use(cors.New())
	appHandlers.Setup(app)
}

func tearDownTestApp() {
	dbTc.Terminate(context.Background())
}

func setupTest() {
	storage.AutoMigratePostgres(db)
}

func tearDownTest() {
	storage.TestDropTablesPostgres(db)
}

func responseMap(b []byte) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	return m
}
