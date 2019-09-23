package app

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/migrations"
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

	return m.Up()
}

func (m *MigrateApp) Up() error {
	db := m.cfg.DB()

	e := migrations.CreateVector(db)

	if e != nil {
		return e
	}

	return nil
}
