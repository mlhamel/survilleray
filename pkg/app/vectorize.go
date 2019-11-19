package app

import (
	"github.com/mlhamel/survilleray/pkg/jobs"
	"github.com/mlhamel/survilleray/pkg/runtime"
)

type VectorizeApp struct {
	context *runtime.Context
}

func NewVectorizeApp(context *runtime.Context) *VectorizeApp {
	return &VectorizeApp{context: context}
}

func (a *VectorizeApp) Run() error {
	job := jobs.NewVectorizeJob(a.context)

	return job.Run()
}
