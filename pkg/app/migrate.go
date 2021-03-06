package app

import (
	"context"
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/config"
)

type MigrateApp struct {
	cfg *config.Config
}

func NewMigrateApp(cfg *config.Config) *MigrateApp {
	return &MigrateApp{cfg}
}

func (app *MigrateApp) Run(ctx context.Context) error {
	fmt.Printf("Migrating %s\n", app.cfg.DSN())

	migrator := models.NewMigrator(app.cfg)

	if err := migrator.Execute(ctx); err != nil {
		return fmt.Errorf("Cannot migrate database: %w", err)
	}

	return nil
}
