package vectorization

import (
	"context"
	"fmt"
	"time"

	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/config"
)

type Job struct {
	cfg *config.Config
}

func NewJob(cfg *config.Config) *Job {
	return &Job{cfg}
}

func (job *Job) Run(ctx context.Context) error {
	pointRepos := models.NewPointRepository(job.cfg)
	vectorRepos := models.NewVectorRepository(job.cfg)

	points, err := pointRepos.FindByVectorizedAt(nil)

	if err != nil {
		return fmt.Errorf("Cannot find points to vectorize: %w", err)
	}

	operation := NewOperation(job.cfg)

	for i := 0; i < len(points); i++ {
		tx := job.cfg.Database().Begin()
		if tx.Error != nil {
			return err
		}

		point := points[i]

		vector, err := operation.GetOrCreateVectorFromPoint(ctx, &point)

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
