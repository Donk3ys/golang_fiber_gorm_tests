package repos_test

import (
	"api/src/constants"
	"api/src/mocks"
	repos "api/src/repo"
	"api/src/storage"
	"api/src/util"
	"context"
	"os"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/valkey-io/valkey-go"
	"gorm.io/gorm"
)

var (
	cache valkey.Client
	db    *gorm.DB
	// dbTc    testcontainers.Container
	// cacheTc testcontainers.Container
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
	ctx := context.Background()
	// createFolders()
	db, _ = storage.TestConnectPostgres(ctx)
	storage.AutoMigratePostgres(db)

	cache, _ = storage.TestConnectValkey(ctx)
	mockFs = &mocks_test.FileSysetmClient{}

	// database.Seed(db)
	repo = repos.Instance{
		Cache: cache,
		Db:    db,
		Fs:    &storage.FileSystem{Client: mockFs},
	}
}

func tearDownTestApp() {
	// ctx := context.Background()
	// dbTc.Terminate(ctx)
	// cacheTc.Terminate(ctx)
}

func setupTest() {
	storage.AutoMigratePostgres(db)
	// Restore durations as these are sometines changed for tests
	durAt, _ := util.ParseDuration(os.Getenv("AUTH_TOKEN_DURATION"))
	constants.ACCESS_TOKEN_DURATION = durAt
	durRt, _ := util.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION"))
	constants.REFRESH_TOKEN_DURATION = durRt
}

func tearDownTest() {
	ctx := context.Background()
	storage.TestDropTablesPostgres(db)
	cache.Do(ctx, cache.B().Reset().Build()).Error()
}
