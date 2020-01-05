package acquisition

import (
	"context"
	"errors"

	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/opensky"
	"github.com/rs/zerolog"
)

var ErrPointNotOverlaps = errors.New("point does not overlaps with district")

type Operation interface {
	GetLatestPoint(context.Context, string) ([]models.Point, error)
	InsertPoint(context.Context, *models.Point) error
	UpdateOverlaps(context.Context, *models.District, *models.Point) error
}

type operationImpl struct {
	logger        *zerolog.Logger
	pointRepos    models.PointRepository
	districtRepos models.DistrictRepository
}

func NewOperation(logger *zerolog.Logger, pointRepos models.PointRepository, districtRepos models.DistrictRepository) Operation {
	return &operationImpl{logger, pointRepos, districtRepos}
}

func (o *operationImpl) GetLatestPoint(ctx context.Context, url string) ([]models.Point, error) {
	return opensky.NewRequestWithLogger(url, o.logger).GetPlanes(ctx)
}

func (o *operationImpl) InsertPoint(ctx context.Context, point *models.Point) error {
	err := o.pointRepos.Create(point)

	if err == nil {
		o.logger.Info().Str("point", point.Icao24).Msg("Inserting point")
		return nil
	}

	if err.Error() == models.ErrorPointalreadyExisted.Error() {
		o.logger.Info().Str("point", point.Icao24).Msg("Point already existed")
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
		o.logger.Info().Str("point", point.Icao24).Str("district", district.Name).Msg("Point does not overlaps with district")
		return ErrPointNotOverlaps
	}

	o.logger.Info().Str("point", point.Icao24).Str("district", district.Name).Msg("Point overlaps with district")

	return o.districtRepos.AppendPoint(district, point)
}
