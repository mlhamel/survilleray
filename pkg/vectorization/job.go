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
		point := points[i]
		vector, err := operation.GetVectorFromPoint(ctx, &point)
		if err != nil {
			return err
		}

		if vector == nil {
			vector, err = operation.CreateVectorFromPoint(ctx, &point)

			if err != nil {
				return err
			}
		}

		if err = operation.MarkPointAsVectorized(ctx, &point); err != nil {
			return err
		}

		if err = operation.AddPointToVector(ctx, &point, vector); err != nil {
			return err
		}
	}

	return nil
}
