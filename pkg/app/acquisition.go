package app

import (
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/jobs"
)

type AcquisitionApp struct {
	cfg *config.Config
}

func NewAcquisitionApp(cfg *config.Config) *AcquisitionApp {
	return &AcquisitionApp{
		cfg: cfg,
	}
}

func (a *AcquisitionApp) Run() error {
	job := jobs.NewAcquisition(a.cfg)

	if errors := job.Run(); len(errors) > 0 {
		return errors[0]
	}

	return nil
}
