package middleware_test

import (
	"api/src/constants"
	"api/src/middleware"
	mocks_test "api/src/mocks"
	repos "api/src/repo"
	"api/src/storage"
	"context"
	"log"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jaswdr/faker"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
)

var (
	app   *fiber.App
	db    *gorm.DB
	dbTc  testcontainers.Container
	fake  faker.Faker
	mware middleware.Instance
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
	db, dbTc = storage.TestConnectPostgres(context.Background())
	storage.AutoMigratePostgres(db)

	repo := repos.Instance{
		Db: db,
		Fs: storage.FileSystem{Client: &mocks_test.FileSysetmClient{}},
	}

	mware = middleware.Instance{
		Repo: repo,
	}

	app = fiber.New(fiber.Config{
		Immutable: true,
	})
	app.Use(mware.AuthenticateAuthTokenAndCreateNewIfExpired)
	app.Get("/", func(c *fiber.Ctx) error {
		uId := (c.Locals(constants.REQ_USER_ID)).(uuid.UUID)
		return c.SendString(uId.String())
	})
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
