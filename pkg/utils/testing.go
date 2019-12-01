package utils

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Suite interface {
	T() *testing.T
	SetupSuite()
}
