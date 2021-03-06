package vectorization

import (
	"context"

	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/config"
)

type App struct {
	cfg         *config.Config
	pointRepos  models.PointRepository
	vectorRepos models.VectorRepository
}

func NewApp(cfg *config.Config) *App {
	pointRepos := models.NewPointRepository(cfg)
	vectorRepos := models.NewVectorRepository(cfg)

	return &App{cfg, pointRepos, vectorRepos}
}

/*
Run is running the vectorization of point of planes.

The vectorization is the process of associating open points of the same route.
We are defining "same route" as sharing the same `Icao24` and `CallSign`.

A `Closed` point is closing the corresponding vector.
*/
func (a *App) Run(ctx context.Context) error {
	a.cfg.Logger().Info().Msg("Running vectorization")
	job := NewJob(a.cfg, a.pointRepos, a.vectorRepos)

	err := job.Run(ctx)

	a.cfg.Logger().Info().Msg("Done: Running vectorization")

	return err
}
