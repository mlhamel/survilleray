package vectorization

import (
	"context"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/config"
)

type Operation interface {
	GetOrCreateVectorFromPoint(context.Context, *models.Point) (*models.Vector, error)
}

type OperationImpl struct {
	cfg *config.Config
}

func NewOperation(cfg *config.Config) Operation {
	return &OperationImpl{cfg}
}

func (operation *OperationImpl) GetOrCreateVectorFromPoint(ctx context.Context, point *models.Point) (*models.Vector, error) {
	repository := models.NewVectorRepository(operation.cfg)

	vector, err := repository.FindByCallSign(point.CallSign)

	if gorm.IsRecordNotFoundError(err) {
		log.Printf("Creating vector for point %s", point.String())
		vector = models.NewVectorFromPoint(point)

		if err = operation.cfg.Database().Create(&vector).Error; err != nil {
			return nil, fmt.Errorf("Cannot create vector: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("Cannot find vector: %w", err)
	}

	return vector, nil
}
