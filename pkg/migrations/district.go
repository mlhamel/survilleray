package migrations

import (
	"fmt"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/models"
)

func CreateDistrict(cfg *config.Config) error {
	fmt.Println("... Creating district table")

	db := cfg.DB()

	if db.HasTable(&models.District{}) {
		return nil
	}

	db.Debug().CreateTable(&models.District{})

	return db.Error
}
