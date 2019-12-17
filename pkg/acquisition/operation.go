package acquisition

import (
	"context"
	"errors"
	"log"

	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/opensky"
)

var ErrPointNotOverlaps = errors.New("point does not overlaps with district")

type Operation interface {
	GetLatestPoint(context.Context, string) ([]models.Point, error)
	InsertPoint(context.Context, *models.Point) error
	UpdateOverlaps(context.Context, *models.District, *models.Point) error
}

type operationImpl struct {
	pointRepos    models.PointRepository
	districtRepos models.DistrictRepository
}

func NewOperation(pointRepos models.PointRepository, districtRepos models.DistrictRepository) Operation {
	return &operationImpl{pointRepos, districtRepos}
}

func (o *operationImpl) GetLatestPoint(ctx context.Context, url string) ([]models.Point, error) {
	return opensky.NewRequest(url).GetPlanes(ctx)
}

func (o *operationImpl) InsertPoint(ctx context.Context, point *models.Point) error {
	err := o.pointRepos.Create(point)

	if err == nil {
		log.Printf("Inserting point with `%s`", point.String())
		return nil
	}

	if err.Error() == models.ErrorPointalreadyExisted.Error() {
		log.Printf("Point `%s` already existed", point.String())
		return nil
	}

	return err
}

func (o *operationImpl) UpdateOverlaps(ctx context.Context, district *models.District, point *models.Point) error {
	overlaps, err := point.CheckOverlaps(district)

	if err != nil {
		return err
	}

	if !overlaps {
		log.Printf("Point `%s` does not overlaps with `%s`", point.Icao24, district.Name)
		return ErrPointNotOverlaps
	}

	log.Printf("Point `%s` overlaps with `%s`", point.Icao24, district.Name)

	return o.districtRepos.AppendPoint(district, point)
}
