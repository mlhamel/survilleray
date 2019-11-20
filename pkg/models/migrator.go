package models

import (
	"fmt"

	"github.com/mlhamel/survilleray/pkg/runtime"
	"github.com/pkg/errors"
)

type migration struct {
	context *runtime.Context
	err     error
}

func (m *migration) migrate(desc string, migrator func(*runtime.Context) error) {
	if m.err == nil {
		fmt.Printf("=== Running: %s ===\n", desc)
		if err := migrator(m.context); err != nil {
			m.err = errors.Wrapf(err, "Failed migrating: %s", desc)
		}
	}
}

func Migrate(context *runtime.Context) error {
	migrator := migration{context: context}

	migrator.migrate("creating point", CreatePoint)
	migrator.migrate("enabling postgis", EnablePostgis)
	migrator.migrate("creating district", CreateDistrict)
	migrator.migrate("creating vector", CreateVector)
	migrator.migrate("creating villeray", CreateVilleray)

	return migrator.err
}
