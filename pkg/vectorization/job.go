package vectorization

import (
	"context"
	"fmt"
	"time"

	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/config"
)

type Job struct {
	cfg         *config.Config
	pointRepos  models.PointRepository
	vectorRepos models.VectorRepository
}

func NewJob(cfg *config.Config, pointRepos models.PointRepository, vectorRepos models.VectorRepository) *Job {
	return &Job{cfg, pointRepos, vectorRepos}
}

func (job *Job) Run(ctx context.Context) error {
	points, err := job.pointRepos.FindByVectorizedAt(nil)

	if err != nil {
		return fmt.Errorf("Cannot find points to vectorize: %w", err)
	}

	operation := NewOperation(job.pointRepos, job.vectorRepos)

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

		if err = job.vectorRepos.AppendPoints(vector, []models.Point{point}); err != nil {
			return fmt.Errorf("Cannot add point to the matching vector: %w", err)
		}

		if err = job.pointRepos.Update(&point, map[string]interface{}{"VectorizedAt": time.Now()}); err != nil {
			return fmt.Errorf("Cannot update VectorizedAt for a point: %w", err)
		}

		if point.OnGround {
			if err = job.vectorRepos.Update(vector, map[string]interface{}{"Closed": true}); err != nil {
				return fmt.Errorf("Cannot update Closed for a vector: %w", err)
			}
		}

		if err = tx.Commit().Error; err != nil {
			return fmt.Errorf("Cannot commit transaction: %w", err)
		}
	}

	return nil
}
