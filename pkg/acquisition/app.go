package acquisition

import (
	"github.com/mlhamel/survilleray/pkg/runtime"
)

type AcquisitionApp struct {
	context *runtime.Context
}

func NewApp(context *runtime.Context) *AcquisitionApp {
	return &AcquisitionApp{context}
}

func (a *AcquisitionApp) Run() error {
	job := NewJob(a.context)

	return job.Run()
}
