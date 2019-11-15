package app

import (
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/jobs"
)

type VectorizeApp struct {
	cfg *config.Config
}

func NewVectorizeApp(cfg *config.Config) *VectorizeApp {
	return &VectorizeApp{cfg: cfg}
}

func (a *VectorizeApp) Run() error {
	job := jobs.NewVectorizeJob(a.cfg)

	return job.Run()
}
