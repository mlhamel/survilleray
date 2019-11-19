package app

import (
	"github.com/mlhamel/survilleray/pkg/jobs"
	"github.com/mlhamel/survilleray/pkg/runtime"
)

type AcquisitionApp struct {
	context *runtime.Context
}

func NewAcquisitionApp(context *runtime.Context) *AcquisitionApp {
	return &AcquisitionApp{context}
}

func (a *AcquisitionApp) Run() error {
	job := jobs.NewAcquisition(a.context)

	return job.Run()
}
