package acquisition

import (
	"context"
	"log"

	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/opensky"
)

type Operation interface {
	GetLatestPoint(context.Context, string) ([]models.Point, error)
	InsertPoint(context.Context, *models.Point) error
}

type OperationImpl struct {
	cfg *config.Config
}

func NewOperation(cfg *config.Config) Operation {
	return &OperationImpl{cfg}
}

func (operation *OperationImpl) GetLatestPoint(ctx context.Context, url string) ([]models.Point, error) {
	var r = opensky.NewRequest(url)

	return r.GetPlanes()
}

func (operation *OperationImpl) InsertPoint(ctx context.Context, point *models.Point) error {
	if !operation.cfg.Database().NewRecord(point) {
		log.Printf("Point `%s` already existed", point.String())

		return nil
	}

	log.Printf("Inserting point with `%s`", point.String())

	return operation.cfg.Database().Create(&point).Error

}
