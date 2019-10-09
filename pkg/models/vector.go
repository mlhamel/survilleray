package models

import (
	"github.com/jinzhu/gorm"
)

type Vector struct {
	gorm.Model
	points []Point `gorm:"many2many:vector_points"`
}

