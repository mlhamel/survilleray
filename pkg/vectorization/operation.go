package vectorization

import (
	"context"
	"fmt"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/models"
	"github.com/rs/zerolog"
)

type Operation interface {
	GetVectorFromPoint(context.Context, *models.Point) (*models.Vector, error)
	CreateVectorFromPoint(context.Context, *models.Point) (*models.Vector, error)
	AddPointToVector(context.Context, *models.Point, *models.Vector) error
	MarkPointAsVectorized(context.Context, *models.Point) error
}

type operationImpl struct {
	statsd           *statsd.Client
	logger           *zerolog.Logger
	pointRepository  models.PointRepository
	vectorRepository models.VectorRepository
}

func NewOperation(statsd *statsd.Client, logger *zerolog.Logger, pointRepository models.PointRepository, vectorRepository models.VectorRepository) Operation {
	return &operationImpl{statsd, logger, pointRepository, vectorRepository}
}

func (o *operationImpl) GetVectorFromPoint(ctx context.Context, point *models.Point) (*models.Vector, error) {
	vector, err := o.vectorRepository.FindByCallSign(point.CallSign)

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("Cannot get vector: %w", err)
	}

	o.logger.Info().Object("point", point).Msg("Point already vectorized")
	o.statsd.Incr("vectorization.retrieve_vector_from_point.update", makeTags(point), 1)

	o.statsd.Gauge("GeoAltitude", point.GeoAltitude, makeTags(point), 1)
	o.statsd.Gauge("BaroAltitude", point.BaroAltitude, makeTags(point), 1)
	o.statsd.Gauge("Velocity", point.Velocity, makeTags(point), 1)

	return vector, err
}

func (o *operationImpl) CreateVectorFromPoint(ctx context.Context, point *models.Point) (*models.Vector, error) {
	vector := models.NewVectorFromPoint(point)
	err := o.vectorRepository.Create(vector)

	if err != nil {
		return nil, fmt.Errorf("Cannot create vector: %w", err)
	}

	o.logger.Info().Str("point", point.String()).Msg("Vector created from point")
	o.statsd.Incr("vectorization.retrieve_vector_from_point.new", makeTags(point), 1)

	o.statsd.Gauge("GeoAltitude", point.GeoAltitude, makeTags(point), 1)
	o.statsd.Gauge("BaroAltitude", point.BaroAltitude, makeTags(point), 1)
	o.statsd.Gauge("Velocity", point.Velocity, makeTags(point), 1)

	return vector, err
}

func (o *operationImpl) AddPointToVector(ctx context.Context, point *models.Point, vector *models.Vector) error {
	err := o.vectorRepository.AppendPoints(vector, []*models.Point{point})

	if err != nil {
		return fmt.Errorf("Cannot add point to the matching vector: %w", err)
	}

	return nil
}

func (o *operationImpl) MarkPointAsVectorized(ctx context.Context, point *models.Point) error {
	return o.pointRepository.Update(point, map[string]interface{}{"VectorizedAt": time.Now()})
}

func makeTags(point *models.Point) []string {
	return []string{fmt.Sprintf("OriginCountry:%s", point.OriginCountry)}
}
