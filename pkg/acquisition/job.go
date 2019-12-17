package acquisition

import (
	"context"
	"log"

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
	operation := NewOperation(j.pointRepos, j.districtRepos)

	villeray, err := j.districtRepos.FindByName("villeray")

	if err != nil {
		return err
	}

	points, err := operation.GetLatestPoint(ctx, j.cfg.OpenSkyURL())

	if err != nil {
		return err
	}

	for i := 0; i < len(points); i++ {
		point := points[i]

		log.Println("")

		log.Printf("Trying to insert point for %s", point.Icao24)

		if err = operation.InsertPoint(ctx, &point); err != nil {
			log.Printf("Cannot insert point for %s, error is %s", point.Icao24, err)
			continue
		}

		log.Printf("Figuring out if %s overlaps with %s", point.Icao24, villeray.Name)

		err = operation.UpdateOverlaps(ctx, villeray, &point)

		if err == ErrPointNotOverlaps {
			err = nil
		}

		if err != nil {
			return err
		}
	}

	return nil
}
