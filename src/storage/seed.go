package storage

import (
	"api/src/models"
	"api/src/util"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	pw, _ := util.GenPasswordHash("123456")
	newUser := models.User{
		FirstName: "David",
		LastName:  "Gericke",
		Email:     "d@e.com",
		Password:  pw,
		Status:    1,
	}
	db.Save(&newUser)
}
