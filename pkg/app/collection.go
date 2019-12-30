package app

import (
	"context"
	"time"

	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/rs/zerolog/log"
)

type collectionApp struct {
	cfg   *config.Config
	repos models.VectorRepository
}

func NewCollectionApp(cfg *config.Config) *collectionApp {
	repos := models.NewVectorRepository(cfg)
	return &collectionApp{cfg, repos}
}

func (a *collectionApp) Run(ctx context.Context) error {
	log.Info().Msg("Running collection")

	yesterday := time.Now().Add(time.Hour * -24)

	vectors, err := a.repos.FindByClosed(false)

	if err != nil {
		log.Error().Err(err).Msg("Cannot find for vectors")
		return err
	}

	log.Info().Msgf("Found %d open vectors", len(vectors))

	for i := range vectors {
		vector := vectors[i]

		if vector.CreatedAt.UTC().After(yesterday) {
			log.Debug().Object("vector", &vector).Msg("Vector is still active")
			continue
		}

		a.repos.Update(&vector, map[string]interface{}{"Closed": true})

		log.Info().Object("vector", &vector).Msg("Closing vector")
	}

	return nil
}
