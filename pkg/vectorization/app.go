package vectorization

import (
	"context"

	"github.com/mlhamel/survilleray/pkg/config"
)

type App struct {
	cfg *config.Config
}

func NewApp(cfg *config.Config) *App {
	return &App{cfg}
}

func (app *App) Run(ctx context.Context) error {
	job := NewJob(app.cfg)

	return job.Run(ctx)
}
