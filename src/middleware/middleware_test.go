package middleware_test

import (
	"api/src/constants"
	"api/src/middleware"
	mocks_test "api/src/mocks"
	repos "api/src/repo"
	"api/src/storage"
	"context"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jaswdr/faker"
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
	constants.SetConstantsFromEnvs("../../.env-test")

	fake = faker.New()
	setupTestApp()
	code := m.Run()
	tearDownTestApp()
	os.Exit(code)
}

func setupTestApp() {
	db, dbTc = storage.TestConnectPostgres(context.Background())
	storage.AutoMigratePostgres(db)
	cache := storage.ConnectRistrettoCache()

	repo := repos.Instance{
		Cache: cache,
		Db:    db,
		Fs:    &storage.FileSystem{Client: &mocks_test.FileSysetmClient{}},
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
	app.Get("/long", func(c *fiber.Ctx) error {
		uId := (c.Locals(constants.REQ_USER_ID)).(uuid.UUID)
		time.Sleep(time.Second * 3)
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
