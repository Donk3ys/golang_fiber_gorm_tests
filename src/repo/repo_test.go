package repos_test

import (
	"api/src/constants"
	"api/src/mocks"
	repos "api/src/repo"
	"api/src/storage"
	"context"
	"os"
	"testing"

	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/jaswdr/faker"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
)

var (
	cache  *marshaler.Marshaler
	db     *gorm.DB
	dbTc   testcontainers.Container
	fake   faker.Faker
	repo   repos.Instance
	mockFs *mocks_test.FileSysetmClient
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

	cache = storage.ConnectRistrettoCache()

	mockFs = &mocks_test.FileSysetmClient{}

	// database.Seed(db)
	repo = repos.Instance{
		Cache: cache,
		Db:    db,
		Fs:    &storage.FileSystem{Client: mockFs},
	}
}

func tearDownTestApp() {
	cache.Clear(context.Background())
	dbTc.Terminate(context.Background())
}

func setupTest() {
	storage.AutoMigratePostgres(db)
}

func tearDownTest() {
	storage.TestDropTablesPostgres(db)
}
