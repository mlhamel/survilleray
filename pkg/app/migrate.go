package app

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/runtime"
)

type MigrateApp struct {
	context *runtime.Context
}

func NewMigrateApp(context *runtime.Context) *MigrateApp {
	return &MigrateApp{context}
}

func (m *MigrateApp) Run() error {
	fmt.Printf("Migrating %s\n", m.context.Config().DSN())

	return models.Migrate(m.context)
}
