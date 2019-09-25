package migrations

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/models"
)

func CreateDistrict(db *gorm.DB) error {
	fmt.Println("... Creating district table")

	if db.HasTable(&models.District{}) {
		return nil
	}

	db.Debug().CreateTable(&models.District{})

	return db.Error
}
