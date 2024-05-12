package repos_test

import (
	"api/src/mocks"
	repos "api/src/repo"
	"api/src/storage"
	"context"
	"log"
	"os"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/joho/godotenv"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	dbTc   testcontainers.Container
	fake   faker.Faker
	repo   repos.Instance
	mockFs *mocks_test.FileSysetmClient
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

	mockFs = &mocks_test.FileSysetmClient{}

	// database.Seed(db)
	repo = repos.Instance{
		Db: db,
		Fs: storage.FileSystem{Client: mockFs},
	}
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
