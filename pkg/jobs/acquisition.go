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
	request := opensky.NewRequest(openskyURL)
	points, err := request.GetPlanes()

	if err != nil {
		return err
	}

	for i := 0; i < len(points); i++ {
		point := points[i]

		if !job.context.Database().NewRecord(point) {
			continue
		}

		log.Printf("Inserting point with `%s`", point.String())

		err := job.context.Database().Create(&point).Error

		if err != nil {
			log.Printf("Cannot insert point for %s, error is %s", point.Icao24, err)
		}
	}

	return nil
}
