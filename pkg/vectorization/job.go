package vectorization

import (
	"fmt"
	"time"

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
	pointRepos := models.NewPointRepository(job.context)
	vectorRepos := models.NewVectorRepository(job.context)

	points, err := pointRepos.FindByVectorizedAt(nil)

	if err != nil {
		return fmt.Errorf("Cannot find points to vectorize: %w", err)
	}

	operation := NewOperation()

	for i := 0; i < len(points); i++ {
		tx := job.context.Database().Begin()
		if tx.Error != nil {
			return err
		}

		point := points[i]

		vector, err := operation.GetOrCreateVectorFromPoint(job.context, &point)

		if err != nil {
			return fmt.Errorf("Cannot find or create vector: %w", err)
		}

		if err = vectorRepos.AppendPoints(vector, []models.Point{point}); err != nil {
			return fmt.Errorf("Cannot add point to the matching vector: %w", err)
		}

		if err = pointRepos.Update(&point, map[string]interface{}{"VectorizedAt": time.Now()}); err != nil {
			return fmt.Errorf("Cannot update VectorizedAt for a point: %w", err)
		}

		if point.OnGround {
			if err = vectorRepos.Update(vector, map[string]interface{}{"Closed": true}); err != nil {
				return fmt.Errorf("Cannot update Closed for a vector: %w", err)
			}
		}

		if err = tx.Commit().Error; err != nil {
			return fmt.Errorf("Cannot commit transaction: %w", err)
		}
	}

	return nil
}
