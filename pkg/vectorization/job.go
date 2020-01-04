package vectorization

import (
	"context"
	"fmt"

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

	operation := NewOperation(job.cfg.Statsd(), job.cfg.Logger(), job.pointRepos, job.vectorRepos)

	for i := 0; i < len(points); i++ {
		tx := job.cfg.Orm().Begin()
		if tx.Error != nil {
			return err
		}

		point := points[i]

		vector, err := operation.RetrieveVectorFromPoint(ctx, &point)

		if err != nil {
			return fmt.Errorf("Cannot find or create vector: %w", err)
		}

		if vector == nil {
			continue
		}

		if err = operation.MarkPointAsVectorized(ctx, &point); err != nil {
			return fmt.Errorf("Cannot update VectorizedAt for a point: %w", err)
		}

		if err = operation.AddPointToVector(ctx, &point, vector); err != nil {
			return fmt.Errorf("Cannot add point to vector: %w", err)
		}

		if err = tx.Commit().Error; err != nil {
			return fmt.Errorf("Cannot commit transaction: %w", err)
		}
	}

	return nil
}
