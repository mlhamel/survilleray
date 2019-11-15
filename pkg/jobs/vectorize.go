package jobs

import (
	"fmt"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/models"
)

type VectorizeJob struct {
	cfg *config.Config
}

func NewVectorizeJob(cfg *config.Config) *VectorizeJob {
	return &VectorizeJob{cfg: cfg}
}

func (job *VectorizeJob) Run() error {
	repository := models.NewPointRepository(job.cfg)
	points, err := repository.Find()

	if err != nil {
		return err
	}
	for i := 0; i < len(points); i++ {
		fmt.Printf("%d %s\n", points[i].ID, points[i].CallSign)
	}

	return nil
}
