package migrations

import (
	"fmt"

	"github.com/mlhamel/survilleray/pkg/config"
)

func EnablePostgis(cfg *config.Config) error {
	fmt.Println("... Enabling postgis extension")

	db := cfg.DB()

	db.Exec("CREATE EXTENSION IF NOT EXISTS postgis")

	return db.Error
}
