package vectorization

import (
	"github.com/mlhamel/survilleray/pkg/runtime"
)

type App struct {
	context *runtime.Context
}

func NewApp(context *runtime.Context) *App {
	return &App{context: context}
}

func (a *App) Run() error {
	job := NewJob(a.context)

	return job.Run()
}
