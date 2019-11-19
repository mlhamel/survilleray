package jobs

import (
	"log"

	"github.com/mlhamel/survilleray/pkg/opensky"
	"github.com/mlhamel/survilleray/pkg/runtime"
)

const openskyURL = "https://opensky-network.org/api/states/all?lamin=%d&lamax=%d&lomin=%d&lomax=%d"

type AcquisitionJob struct {
	context *runtime.Context
}

func NewAcquisition(context *runtime.Context) *AcquisitionJob {
	return &AcquisitionJob{context}
}

func (job *AcquisitionJob) Run() error {
	var r = opensky.NewRequest(openskyURL)

	points, err := r.GetPlanes()

	if err != nil {
		return err
	}

	for i := 0; i < len(points); i++ {
		v := points[i]

		if !job.context.Database().NewRecord(v) {
			continue
		}

		log.Printf("Inserting point with `%s`", v.String())

		err := job.context.Database().Create(&v).Error

		if err != nil {
			log.Printf("Cannot insert point for %s, error is %s", v.Icao24, err)
		}
	}

	return nil
}
