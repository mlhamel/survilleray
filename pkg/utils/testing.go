package utils

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/mlhamel/survilleray/pkg/config"
)

type Suite interface {
	T() *testing.T
	SetupSuite()
}

func SetupSuite(s Suite) {
	var (
		orm *gorm.DB
		err error
	)

	orm, err = config.NewTestDatabase("postgres://postgres:docker@survilleray.railgun:5432/survilleray_test")

	if err != nil {
		panic(err)
	}

	config.NewConfigWithDB(orm)
}
