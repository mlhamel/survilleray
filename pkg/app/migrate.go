package app

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/migrations"
)

type MigrateApp struct {
	cfg *config.Config
}

type migration struct {
	db  *gorm.DB
	err error
}

func NewMigrateApp(cfg *config.Config) *MigrateApp {
	return &MigrateApp{
		cfg: cfg,
	}
}

func (m *MigrateApp) Run() error {
	fmt.Printf("Migrating %s\n", m.cfg.DSN())

	return m.Migrate()
}

func (m *migration) migrate(desc string, migrator func(*gorm.DB) error) {
	if m.err == nil {
		if err := migrator(m.db); err != nil {
			m.err = errors.Wrapf(err, "Failed migrating: %s", desc)
		}
	}
}

func (m *MigrateApp) Migrate() error {
	migrator := migration{db: m.cfg.DB()}

	migrator.migrate("creating vector", migrations.CreateVector)
	migrator.migrate("enabling postgis", migrations.EnablePostgis)
	migrator.migrate("creating district", migrations.CreateDistrict)
	migrator.migrate("creating villeray", migrations.CreateVilleray)

	return migrator.err
}
