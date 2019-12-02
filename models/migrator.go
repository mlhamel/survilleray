package models

import (
	"context"
	"fmt"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/pkg/errors"
)

type migration struct {
	cfg *config.Config
	err error
}

func (m *migration) migrate(ctx context.Context, description string, migrator func(context.Context, *config.Config) error) {
	if m.err == nil {
		fmt.Printf("=== Running: %s ===\n", description)
		if err := migrator(ctx, m.cfg); err != nil {
			m.err = errors.Wrapf(err, "Failed migrating: %s", description)
		}
	}
}

func Migrate(ctx context.Context, cfg *config.Config) error {
	migrator := migration{cfg: cfg}

	migrator.migrate(ctx, "creating point", CreatePoint)
	migrator.migrate(ctx, "enabling postgis", EnablePostgis)
	migrator.migrate(ctx, "creating district", CreateDistrict)
	migrator.migrate(ctx, "creating vector", CreateVector)
	migrator.migrate(ctx, "creating villeray", CreateVilleray)

	return migrator.err
}
