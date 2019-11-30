package acquisition

import (
	"log"

	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/opensky"
	"github.com/mlhamel/survilleray/pkg/runtime"
)

type Operation interface {
	GetLatestPoint(string) ([]models.Point, error)
	InsertPoint(*runtime.Context, *models.Point) error
	CloseVectors(*runtime.Context, *models.Point) error
}

type OperationImpl struct {
}

func NewOperation() Operation {
	return &OperationImpl{}
}

func (operation *OperationImpl) GetLatestPoint(url string) ([]models.Point, error) {
	var r = opensky.NewRequest(url)

	return r.GetPlanes()
}

func (operation *OperationImpl) InsertPoint(context *runtime.Context, point *models.Point) error {
	if !context.Database().NewRecord(point) {
		log.Printf("Point `%s` already existed", point.String())

		return nil
	}

	log.Printf("Inserting point with `%s`", point.String())

	return context.Database().Create(&point).Error

}

func (operation *OperationImpl) CloseVectors(context *runtime.Context, point *models.Point) error {
	if !point.OnGround {
		return nil
	}

	repository := models.NewVectorRepository(context)

	vectors, err := repository.FindByPoint(point)

	if err != nil {
		return err
	}

	for i := 0; i < len(vectors); i++ {
		if !vectors[i].Closed {
			repository.Update(&vectors[i], map[string]interface{}{"Closed": true})
		}
	}

	return nil
}
