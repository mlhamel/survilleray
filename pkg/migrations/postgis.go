package migrations

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

func EnablePostgis(db *gorm.DB) error {
	fmt.Println("... Enabling postgis extension")

	db.Exec("CREATE EXTENSION IF NOT EXISTS postgis")

	return db.Error
}
