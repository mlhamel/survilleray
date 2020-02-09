package acquisition

import (
	"context"
	"errors"
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/opensky"
	"github.com/rs/zerolog"
)

var ErrPointNotOverlaps = errors.New("point does not overlaps with district")

type Operation interface {
	GetLatestPoint(context.Context, string) ([]models.Point, error)
	InsertPoint(context.Context, *models.Point) error
	UpdateOverlaps(context.Context, *models.District, *models.Point) (bool, error)
}

type operationImpl struct {
	statsd        *statsd.Client
	logger        *zerolog.Logger
	pointRepos    models.PointRepository
	districtRepos models.DistrictRepository
}

func NewOperation(statsd *statsd.Client, logger *zerolog.Logger, pointRepos models.PointRepository, districtRepos models.DistrictRepository) Operation {
	return &operationImpl{statsd, logger, pointRepos, districtRepos}
}

func (o *operationImpl) GetLatestPoint(ctx context.Context, url string) ([]models.Point, error) {
	request := opensky.NewPointsRequest(url, o.logger)
	err := request.Run(ctx)
	if err != nil {
		return []models.Point{}, err
	}
	return request.Result(), nil
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

func (o *operationImpl) UpdateOverlaps(ctx context.Context, district *models.District, point *models.Point) (bool, error) {
	overlaps, err := point.CheckOverlaps(district)

	if err != nil {
		return false, fmt.Errorf("Cannot check overlaps: %w", err)
	}

	if !overlaps {
		o.logger.Info().Str("point", point.Icao24).Str("district", district.Name).Msg("Point does not overlaps with district")
		o.statsd.Incr("acquisition.notmatch", []string{fmt.Sprintf("district:%d", district.ID)}, 1)
		return false, nil
	}

	o.statsd.Incr("acquisition.match", []string{fmt.Sprintf("district:%d", district.ID)}, 1)

	o.logger.Info().
		Str("point", point.Icao24).
		Uint("point-id", point.ID).
		Str("district", district.Name).
		Uint("district-id", district.ID).
		Msg("Point overlaps with district")

	if err := o.districtRepos.AppendPoint(district, point); err != nil {
		return false, fmt.Errorf("Cannot append point: %w", err)
	}

	return true, nil
}
