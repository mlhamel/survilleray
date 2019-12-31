package acquisition

import (
	"context"

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
	logger := j.cfg.Logger()
	operation := NewOperation(logger, j.pointRepos, j.districtRepos)

	villeray, err := j.districtRepos.FindByName("villeray")

	if err != nil {
		return err
	}

	points, err := operation.GetLatestPoint(ctx, j.cfg.OpenSkyURL())

	if err != nil {
		return err
	}

	j.cfg.Statsd().
		Gauge("acquistion.job.found", float64(len(points)), []string{}, 1)

	for i := 0; i < len(points); i++ {
		point := points[i]

		logger.Info().Str("point", point.Icao24).Msg("Trying to insert point")

		if err = operation.InsertPoint(ctx, &point); err != nil {
			j.cfg.Statsd().Incr("acquistion.job.invalid", []string{}, 1)
			logger.Warn().Err(err).Str("point", point.Icao24).Msg("Cannot insert point")
			continue
		}

		logger.Info().
			Str("point", point.Icao24).
			Str("district", villeray.Name).
			Msg("Figuring out if point overlaps with district")

		err = operation.UpdateOverlaps(ctx, villeray, &point)

		if err == ErrPointNotOverlaps {
			j.cfg.Statsd().Incr("acquistion.job.nooverlaps", []string{}, 1)
			err = nil
		} else {
			j.cfg.Statsd().Incr("acquistion.job.overlaps", []string{}, 1)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
