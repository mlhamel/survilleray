package vectorization

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/runtime"
)

type Operation interface {
	GetOrCreateVectorFromPoint(*runtime.Context, *models.Point) (*models.Vector, error)
}

type OperationImpl struct {
}

func NewOperation() Operation {
	return &OperationImpl{}
}

func (operation *OperationImpl) GetOrCreateVectorFromPoint(context *runtime.Context, point *models.Point) (*models.Vector, error) {
	repository := models.NewVectorRepository(context)

	vector, err := repository.FindByCallSign(point.CallSign)

	if gorm.IsRecordNotFoundError(err) {
		log.Printf("Creating vector for point %s", point.String())
		vector = models.NewVectorFromPoint(point)

		if err = context.Database().Create(&vector).Error; err != nil {
			return nil, fmt.Errorf("Cannot create vector: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("Cannot find vector: %w", err)
	}

	return vector, nil
}
