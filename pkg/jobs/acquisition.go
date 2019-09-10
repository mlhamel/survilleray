package jobs

import (
	"fmt"

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

func (job *AcquisitionJob) Run() {
	db := job.Config.DB()
	var r = opensky.NewRequest(openskyURL)

	vectors, err := r.GetPlanes()

	if err != nil {
		panic(err)
	}

	for i := 0; i < len(vectors); i++ {
		v := vectors[i]
		if db.NewRecord(v) {
			db.Create(&v)

			errors := db.GetErrors()

			if len(errors) == 0 {
				fmt.Printf("Inserted: %s\n", v.String())
			}
		}
	}

	defer db.Close()
}
