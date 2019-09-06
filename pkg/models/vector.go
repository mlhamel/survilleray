package main

import (
	"github.com/jinzhu/gorm"
)

type Vector struct {
	gorm.Model
	Code  string
	Price uint
}
