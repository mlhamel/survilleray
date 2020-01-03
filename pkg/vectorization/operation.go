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
	RetrieveVectorFromPoint(context.Context, *models.Point) (*models.Vector, error)
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

func (o *operationImpl) RetrieveVectorFromPoint(ctx context.Context, point *models.Point) (*models.Vector, error) {
	vector, err := o.vectorRepository.FindByCallSign(point.CallSign)

	if gorm.IsRecordNotFoundError(err) && point.OnGround != true {
		vector = models.NewVectorFromPoint(point)
		err = o.vectorRepository.Create(vector)

		if err != nil {
			return nil, fmt.Errorf("Cannot create vector: %w", err)
		}

		o.logger.Info().Str("point", point.String()).Msg("Vector created from point")
		o.statsd.Incr("vectorization.retrieve_vector_from_point.new", makeTags(point), 1)
	} else if gorm.IsRecordNotFoundError(err) && point.OnGround == true {
		vector = nil
		err = nil
	} else {
		o.statsd.Incr("vectorization.retrieve_vector_from_point.update", makeTags(point), 1)
	}

	if vector != nil {
		o.statsd.Gauge("GeoAltitude", point.GeoAltitude, makeTags(point), 1)
		o.statsd.Gauge("BaroAltitude", point.BaroAltitude, makeTags(point), 1)
		o.statsd.Gauge("Velocity", point.Velocity, makeTags(point), 1)
	}

	if err != nil {
		return nil, fmt.Errorf("Cannot find vector: %w", err)
	}

	return vector, nil
}

func makeTags(point *models.Point) []string {
	return []string{fmt.Sprintf("OriginCountry:%s", point.OriginCountry)}
}

func (o *operationImpl) AddPointToVector(ctx context.Context, point *models.Point, vector *models.Vector) error {
	if err := o.vectorRepository.AppendPoints(vector, []*models.Point{point}); err != nil {
		return fmt.Errorf("Cannot add point to the matching vector: %w", err)
	}

	if point.OnGround {
		if err := o.vectorRepository.Update(vector, map[string]interface{}{"Closed": true}); err != nil {
			return fmt.Errorf("Cannot update Closed for a vector: %w", err)
		}
	}

	return nil
}

func (o *operationImpl) MarkPointAsVectorized(ctx context.Context, point *models.Point) error {
	return o.pointRepository.Update(point, map[string]interface{}{"VectorizedAt": time.Now()})
}
