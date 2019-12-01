package vectorization

import (
	"github.com/mlhamel/survilleray/pkg/config"
)

type App struct {
	cfg *config.Config
}

func NewApp(cfg *config.Config) *App {
	return &App{cfg}
}

func (app *App) Run() error {
	job := NewJob(app.cfg)

	return job.Run()
}
