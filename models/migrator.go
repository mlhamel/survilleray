package models

import (
	"fmt"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/pkg/errors"
)

type migration struct {
	cfg *config.Config
	err error
}

func (m *migration) migrate(desc string, migrator func(*config.Config) error) {
	if m.err == nil {
		fmt.Printf("=== Running: %s ===\n", desc)
		if err := migrator(m.cfg); err != nil {
			m.err = errors.Wrapf(err, "Failed migrating: %s", desc)
		}
	}
}

func Migrate(cfg *config.Config) error {
	migrator := migration{cfg: cfg}

	migrator.migrate("creating point", CreatePoint)
	migrator.migrate("enabling postgis", EnablePostgis)
	migrator.migrate("creating district", CreateDistrict)
	migrator.migrate("creating vector", CreateVector)
	migrator.migrate("creating villeray", CreateVilleray)

	return migrator.err
}
