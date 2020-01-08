package app

import (
	"context"
	"time"

	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/config"
)

const TIMEOUT = time.Hour * -2

type collectionApp struct {
	cfg   *config.Config
	repos models.VectorRepository
}

func NewCollectionApp(cfg *config.Config) *collectionApp {
	repos := models.NewVectorRepository(cfg)
	return &collectionApp{cfg, repos}
}

func (a *collectionApp) Run(ctx context.Context) error {
	a.cfg.Logger().Info().Msg("Running collection")

	yesterday := time.Now().Add(TIMEOUT)

	vectors, err := a.repos.FindByClosed(false)

	if err != nil {
		a.cfg.Logger().Error().Err(err).Msg("Cannot find for vectors")
		return err
	}

	a.cfg.Logger().Info().Msgf("Found %d open vectors", len(vectors))

	for i := range vectors {
		vector := vectors[i]

		if vector.CreatedAt.UTC().After(yesterday) {
			a.cfg.Logger().Debug().Object("vector", &vector).Msg("Vector is still active")
			continue
		}

		a.repos.Update(&vector, map[string]interface{}{"Closed": true})

		a.cfg.Logger().Info().Object("vector", &vector).Msg("Closing vector")
	}

	a.cfg.Logger().Info().Msg("Done: Running collection")
	return nil
}
