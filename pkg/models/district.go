package models

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
)

type District struct {
	gorm.Model
	Name     string          `gorm:"size:20;unique_index"`
	Geometry json.RawMessage `gorm:"type:geometry(MULTIPOLYGON, 4326)"`
}
