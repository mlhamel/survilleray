package vectorization

import (
	"fmt"

	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/runtime"
)

type Job struct {
	context *runtime.Context
}

func NewJob(context *runtime.Context) *Job {
	return &Job{context}
}

func (job *Job) Run() error {
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
