package jobs

import (
	"log"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/opensky"
)

const openskyURL = "https://opensky-network.org/api/states/all?lamin=%d&lamax=%d&lomin=%d&lomax=%d"

type AcquisitionJob struct {
	Config *config.Config
}

func NewAcquisition(c *config.Config) *AcquisitionJob {
	return &AcquisitionJob{Config: c}
}

func (job *AcquisitionJob) Run() error {
	db := job.Config.DB()

	var r = opensky.NewRequest(openskyURL)

	points, err := r.GetPlanes()

	if err != nil {
		return err
	}

	for i := 0; i < len(points); i++ {
		v := points[i]

		if !db.NewRecord(v) {
			continue
		}

		log.Printf("Inserting point with `%s`", v.String())

		err := db.Create(&v).Error

		if err != nil {
			log.Printf("Cannot insert point for %s, error is %s", v.Icao24, err)
		}
	}

	defer db.Close()

	return nil
}
