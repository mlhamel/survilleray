package jobs

import (
	"fmt"

	"github.com/mlhamel/survilleray/pkg/models"
	"github.com/mlhamel/survilleray/pkg/runtime"
)

type VectorizeJob struct {
	context *runtime.Context
}

func NewVectorizeJob(context *runtime.Context) *VectorizeJob {
	return &VectorizeJob{context}
}

func (job *VectorizeJob) Run() error {
	repository := models.NewPointRepository(job.context)
	points, err := repository.Find()

	if err != nil {
		return err
	}
	for i := 0; i < len(points); i++ {
		fmt.Printf("%d %s\n", points[i].ID, points[i].CallSign)
	}

	return nil
}
