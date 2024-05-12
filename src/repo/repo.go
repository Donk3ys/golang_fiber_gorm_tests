package repos

import (
	"api/src/storage"

	"gorm.io/gorm"
)

type Instance struct {
	Db *gorm.DB
	Fs storage.FileSystem
}
