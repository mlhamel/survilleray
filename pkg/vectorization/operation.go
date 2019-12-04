package vectorization

import (
	"context"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/models"
)

type Operation interface {
	GetOrCreateVectorFromPoint(context.Context, *models.Point) (*models.Vector, error)
}

type OperationImpl struct {
	pointRepository  models.PointRepository
	vectorRepository models.VectorRepository
}

func NewOperation(pointRepository models.PointRepository, vectorRepository models.VectorRepository) Operation {
	return &OperationImpl{pointRepository, vectorRepository}
}

func (operation *OperationImpl) GetOrCreateVectorFromPoint(ctx context.Context, point *models.Point) (*models.Vector, error) {
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
