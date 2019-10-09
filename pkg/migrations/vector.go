package migrations

import (
	"fmt"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/models"
)

func CreateVector(cfg *config.Config) error {
	fmt.Println("... Creating point table")

	db := cfg.DB()

	if db.HasTable(&models.Vector{}) {
		return nil
	}

	db.Debug().AutoMigrate(&models.Vector{})

	return db.Error
}
