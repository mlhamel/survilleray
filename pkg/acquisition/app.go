package acquisition

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

func (a *App) Run(ctx context.Context) error {
	job := NewJob(a.cfg)

	return job.Run(ctx)
}
