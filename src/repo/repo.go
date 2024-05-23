package repos

import (
	"api/src/storage"

	"github.com/eko/gocache/lib/v4/marshaler"
	"gorm.io/gorm"
)

type Instance struct {
	Cache *marshaler.Marshaler
	Db    *gorm.DB
	Fs    *storage.FileSystem
}
