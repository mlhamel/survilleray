package app

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/migrations"
)

type MigrateApp struct {
	cfg *config.Config
}

type migration struct {
	cfg *config.Config
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

func (m *migration) migrate(desc string, migrator func(*config.Config) error) {
	if m.err == nil {
		if err := migrator(m.cfg); err != nil {
			m.err = errors.Wrapf(err, "Failed migrating: %s", desc)
		}
	}
}

func (m *MigrateApp) Migrate() error {
	migrator := migration{cfg: m.cfg}

	migrator.migrate("creating point", migrations.CreatePoint)
	migrator.migrate("enabling postgis", migrations.EnablePostgis)
	migrator.migrate("creating district", migrations.CreateDistrict)
	migrator.migrate("creating villeray", migrations.CreateVilleray)

	return migrator.err
}
