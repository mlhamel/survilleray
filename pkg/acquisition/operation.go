package acquisition

import (
	"log"

	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/opensky"
)

type Operation interface {
	GetLatestPoint(string) ([]models.Point, error)
	InsertPoint(*config.Config, *models.Point) error
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

func (operation *OperationImpl) InsertPoint(cfg *config.Config, point *models.Point) error {
	if !cfg.Database().NewRecord(point) {
		log.Printf("Point `%s` already existed", point.String())

		return nil
	}

	log.Printf("Inserting point with `%s`", point.String())

	return cfg.Database().Create(&point).Error

}
