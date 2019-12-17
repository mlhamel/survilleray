package models

import (
	"context"
	"fmt"

	"github.com/mlhamel/survilleray/pkg/config"
)

type Migration = struct {
	description string
	migration   func(context.Context, *config.Config) error
}

type Migrator struct {
	cfg *config.Config
	err error
}

var migrations = []Migration{
	Migration{"enabling postgis", EnablePostgis},
	Migration{"creating point", CreatePoint},
	Migration{"creating district", CreateDistrict},
	Migration{"creating vector", CreateVector},
	Migration{"creating villeray", CreateVilleray},
}

func NewMigrator(cfg *config.Config) *Migrator {
	return &Migrator{cfg: cfg}
}

func (migrator *Migrator) Migrations() []Migration {
	return migrations
}

func (migrator *Migrator) Execute(ctx context.Context) error {
	for i := range migrator.Migrations() {
		fmt.Printf("=== Running: %s ===", migrations[i].description)
		if err := migrations[i].migration(ctx, migrator.cfg); err != nil {
			return fmt.Errorf("Failed migrating %s: %w", migrations[i].description, err)
		}
	}
	return nil
}
