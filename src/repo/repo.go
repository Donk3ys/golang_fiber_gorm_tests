package repos

import (
	"api/src/storage"

	"github.com/valkey-io/valkey-go"
	"gorm.io/gorm"
)

type Instance struct {
	Cache valkey.Client
	Db    *gorm.DB
	Fs    *storage.FileSystem
}
