package acquisition

import (
	"context"

	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/config"
)

type App struct {
	cfg        *config.Config
	repository models.PointRepository
}

func NewApp(cfg *config.Config) *App {
	repository := models.NewPointRepository(cfg)

	return &App{cfg: cfg, repository: repository}
}

func (app *App) Run(ctx context.Context) error {
	job := NewJob(app.cfg)

	return job.Run(ctx, app.repository)
}
