package handlers_test

import (
	"api/src/constants"
	"api/src/handlers"
	"api/src/middleware"
	"api/src/mocks"
	"api/src/notification"
	"api/src/repo"
	"api/src/storage"
	"context"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jaswdr/faker"
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
	constants.SetConstantsFromEnvs("../../.env-test")

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

	cache := storage.ConnectRistrettoCache()

	mockEmail = &mocks_test.EmailCLient{}
	mockSms = &mocks_test.SmsCLient{}
	mockFs = &mocks_test.FileSysetmClient{}

	repo := repos.Instance{
		Cache: cache,
		Db:    db,
		Fs: &storage.FileSystem{
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
