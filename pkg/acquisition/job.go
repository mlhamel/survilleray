package acquisition

import (
	"log"

	"github.com/mlhamel/survilleray/pkg/config"
)

type Job struct {
	cfg *config.Config
}

func NewJob(cfg *config.Config) *Job {
	return &Job{cfg}
}

func (job *Job) Run() error {
	operation := NewOperation()

	points, err := operation.GetLatestPoint(job.cfg.OpenSkyURL())

	if err != nil {
		return err
	}

	for i := 0; i < len(points); i++ {
		point := points[i]

		if err = operation.InsertPoint(job.cfg, &point); err != nil {
			log.Printf("Cannot insert point for %s, error is %s", point.Icao24, err)
		}
	}

	return nil
}
