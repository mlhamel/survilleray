package migrations

import (
	"fmt"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/models"
)

func CreateVector(cfg *config.Config) error {
	fmt.Println("... Creating vector table")

	db := cfg.DB()

	if db.HasTable(&models.Vector{}) {
		return nil
	}

	db.Debug().CreateTable(&models.Vector{})

	return db.Error
}
