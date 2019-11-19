package utils

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/runtime"
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

	orm, err = config.NewTestDatabase(config.GetEnv("DATABASE_URL", ""))

	if err != nil {
		panic(err)
	}

	cfg := config.NewConfig()
	runtime.NewContext(cfg, orm)
}
