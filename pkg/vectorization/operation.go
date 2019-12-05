package vectorization

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/models"
)

type Operation interface {
	RetrieveVectorFromPoint(context.Context, *models.Point) (*models.Vector, error)
	AddPointToVector(context.Context, *models.Point, *models.Vector) error
	MarkPointAsVectorized(context.Context, *models.Point) error
}

type OperationImpl struct {
	pointRepository  models.PointRepository
	vectorRepository models.VectorRepository
}

func NewOperation(pointRepository models.PointRepository, vectorRepository models.VectorRepository) Operation {
	return &OperationImpl{pointRepository, vectorRepository}
}

func (operation *OperationImpl) RetrieveVectorFromPoint(ctx context.Context, point *models.Point) (*models.Vector, error) {
	vector, err := operation.vectorRepository.FindByCallSign(point.CallSign)

	if gorm.IsRecordNotFoundError(err) {
		log.Printf("Creating vector for point %s", point.String())
		vector = models.NewVectorFromPoint(point)

		if err = operation.vectorRepository.Create(vector); err != nil {
			return nil, fmt.Errorf("Cannot create vector: %w", err)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("Cannot find vector: %w", err)
	}

	return vector, nil
}

func (operation *OperationImpl) AddPointToVector(ctx context.Context, point *models.Point, vector *models.Vector) error {
	if err := operation.vectorRepository.AppendPoints(vector, []*models.Point{point}); err != nil {
		return fmt.Errorf("Cannot add point to the matching vector: %w", err)
	}

	if point.OnGround {
		if err := operation.vectorRepository.Update(vector, map[string]interface{}{"Closed": true}); err != nil {
			return fmt.Errorf("Cannot update Closed for a vector: %w", err)
		}
	}

	return nil
}

func (operation *OperationImpl) MarkPointAsVectorized(ctx context.Context, point *models.Point) error {
	return operation.pointRepository.Update(point, map[string]interface{}{"VectorizedAt": time.Now()})
}
