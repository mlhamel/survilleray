package acquisition

import (
	"context"
	"fmt"

	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/config"
)

type job struct {
	cfg           *config.Config
	pointRepos    models.PointRepository
	districtRepos models.DistrictRepository
}

func NewJob(cfg *config.Config, pointRepos models.PointRepository, districtRepos models.DistrictRepository) *job {
	return &job{cfg, pointRepos, districtRepos}
}

func (j *job) Run(ctx context.Context) error {
	statsd := j.cfg.Statsd()
	logger := j.cfg.Logger()
	operation := NewOperation(statsd, logger, j.pointRepos, j.districtRepos)

	villeray, err := j.districtRepos.FindByName("villeray")

	if err != nil {
		return err
	}

	points, err := operation.GetLatestPoint(ctx, j.cfg.OpenSkyURL())
	logger.Debug().Msg("Request done")

	if err != nil {
		return err
	}

	j.cfg.Statsd().Gauge("acquistion.job.found", float64(len(points)), []string{}, 1)

	for i := 0; i < len(points); i++ {
		point := points[i]

		logger.Info().Str("point", point.Icao24).Msg("Trying to insert point")

		if err = operation.InsertPoint(ctx, &point); err != nil {
			j.cfg.Statsd().Incr("acquistion.job.invalid", []string{}, 1)
			logger.Debug().Err(err).Object("point", &point).Msg("Cannot insert point")
			continue
		} else if !point.OnGround {
			j.cfg.Statsd().Incr("acquistion.job.insert", []string{"onground:0"}, 1)
		} else {
			j.cfg.Statsd().Incr("acquistion.job.insert", []string{"onground:1"}, 1)
		}

		logger.Debug().Object("point", &point).Msg("Point was inserted")

		isOverlaps, err := operation.UpdateOverlaps(ctx, villeray, &point)

		if err != nil {
			return fmt.Errorf("Cannot update overlaps: %w", err)
		}

		if isOverlaps {
			j.cfg.Statsd().Incr("acquistion.job.overlaps", []string{}, 1)
		} else {
			j.cfg.Statsd().Incr("acquistion.job.nooverlaps", []string{}, 1)
		}
	}

	return nil
}
