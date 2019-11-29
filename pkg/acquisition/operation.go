package acquisition

import (
	"log"

	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/opensky"
	"github.com/mlhamel/survilleray/pkg/runtime"
)

type AcquisitionOperation interface {
	GetLatestPoint(string) ([]models.Point, error)
	InsertPoint(*runtime.Context, *models.Point) error
}

type AcquisitionOperationImpl struct {
}

func NewAcquisitionOperation() AcquisitionOperation {
	return &AcquisitionOperationImpl{}
}

func (operation *AcquisitionOperationImpl) GetLatestPoint(url string) ([]models.Point, error) {
	var r = opensky.NewRequest(url)

	return r.GetPlanes()
}

func (operation *AcquisitionOperationImpl) InsertPoint(context *runtime.Context, point *models.Point) error {
	if !context.Database().NewRecord(point) {
		log.Printf("Point `%s` already existed", point.String())

		return nil
	}

	log.Printf("Inserting point with `%s`", point.String())

	return context.Database().Create(&point).Error

}
