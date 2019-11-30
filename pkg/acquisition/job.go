package acquisition

import (
	"log"

	"github.com/mlhamel/survilleray/pkg/runtime"
)

type Job struct {
	context *runtime.Context
}

func NewJob(context *runtime.Context) *Job {
	return &Job{context}
}

func (job *Job) Run() error {
	cfg := job.context.Config()
	operation := NewOperation()

	points, err := operation.GetLatestPoint(cfg.OpenSkyURL())

	if err != nil {
		return err
	}

	for i := 0; i < len(points); i++ {
		point := points[i]

		if err = operation.InsertPoint(job.context, &point); err != nil {
			log.Printf("Cannot insert point for %s, error is %s", point.Icao24, err)
		}
	}

	return nil
}
