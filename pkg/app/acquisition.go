package app

import (
	"github.com/mlhamel/survilleray/pkg/acquisition"
	"github.com/mlhamel/survilleray/pkg/runtime"
)

type AcquisitionApp struct {
	context *runtime.Context
}

func NewAcquisitionApp(context *runtime.Context) *AcquisitionApp {
	return &AcquisitionApp{context}
}

func (a *AcquisitionApp) Run() error {
	job := acquisition.NewAcquisition(a.context)

	return job.Run()
}
