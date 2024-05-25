package repos_test

import (
	"api/src/constants"
	"api/src/mocks"
	repos "api/src/repo"
	"api/src/storage"
	"context"
	"os"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/testcontainers/testcontainers-go"
	"github.com/valkey-io/valkey-go"
	"gorm.io/gorm"
)

var (
	cache   valkey.Client
	db      *gorm.DB
	dbTc    testcontainers.Container
	cacheTc testcontainers.Container
	fake    faker.Faker
	repo    repos.Instance
	mockFs  *mocks_test.FileSysetmClient
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
	ctx := context.Background()
	// createFolders()
	db, dbTc = storage.TestConnectPostgres(ctx)
	storage.AutoMigratePostgres(db)

	cache, cacheTc = storage.TestConnectValkey(ctx)

	mockFs = &mocks_test.FileSysetmClient{}

	// database.Seed(db)
	repo = repos.Instance{
		Cache: cache,
		Db:    db,
		Fs:    &storage.FileSystem{Client: mockFs},
	}
}

func tearDownTestApp() {
	ctx := context.Background()
	// cache.Do(ctx, cache.B().Reset().Build()).Error()
	dbTc.Terminate(ctx)
	cacheTc.Terminate(ctx)
}

func setupTest() {
	storage.AutoMigratePostgres(db)
}

func tearDownTest() {
	storage.TestDropTablesPostgres(db)
}
