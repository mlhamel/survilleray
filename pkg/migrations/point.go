package migrations

import (
	"fmt"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/models"
)

func CreatePoint(cfg *config.Config) error {
	fmt.Println("... Creating point table")

	db := cfg.DB()

	if db.HasTable(&models.Point{}) {
		return nil
	}

	db.Debug().CreateTable(&models.Point{})

	return db.Error
}
