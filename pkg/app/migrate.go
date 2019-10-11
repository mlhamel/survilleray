package app

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/models"
)

type MigrateApp struct {
	cfg *config.Config
}

func NewMigrateApp(cfg *config.Config) *MigrateApp {
	return &MigrateApp{
		cfg: cfg,
	}
}

func (m *MigrateApp) Run() error {
	fmt.Printf("Migrating %s\n", m.cfg.DSN())

	return models.Migrate(m.cfg)
}



