package vectorization

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/config"
)

type Operation interface {
	GetOrCreateVectorFromPoint(*config.Config, *models.Point) (*models.Vector, error)
}

type OperationImpl struct {
}

func NewOperation() Operation {
	return &OperationImpl{}
}

func (operation *OperationImpl) GetOrCreateVectorFromPoint(cfg *config.Config, point *models.Point) (*models.Vector, error) {
	repository := models.NewVectorRepository(cfg)

	vector, err := repository.FindByCallSign(point.CallSign)

	if gorm.IsRecordNotFoundError(err) {
		log.Printf("Creating vector for point %s", point.String())
		vector = models.NewVectorFromPoint(point)

		if err = cfg.Database().Create(&vector).Error; err != nil {
			return nil, fmt.Errorf("Cannot create vector: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("Cannot find vector: %w", err)
	}

	return vector, nil
}
