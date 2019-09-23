package migrations

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/models"
)

func CreateVector(db *gorm.DB) error {
	fmt.Println("... Creating vector table")
	if db.HasTable(&models.Vector{}) {
		return nil
	}

	db.Debug().CreateTable(&models.Vector{})

	return db.Error
}
