package acquisition

import (
	"log"

	"github.com/mlhamel/survilleray/pkg/runtime"
)

type AcquisitionJob struct {
	context *runtime.Context
}

func NewJob(context *runtime.Context) *AcquisitionJob {
	return &AcquisitionJob{context}
}

func (job *AcquisitionJob) Run() error {
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
