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

func (app *App) Run(ctx context.Context) error {
	job := NewJob(app.cfg, app.pointRepos, app.vectorRepos)

	return job.Run(ctx)
}
