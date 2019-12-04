package acquisition

import (
	"context"
	"log"

	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/opensky"
)

type Operation interface {
	GetLatestPoint(context.Context, string) ([]models.Point, error)
	InsertPoint(context.Context, *models.Point) error
}

type OperationImpl struct {
	repository models.PointRepository
}

func NewOperation(repository models.PointRepository) Operation {
	return &OperationImpl{repository}
}

func (operation *OperationImpl) GetLatestPoint(ctx context.Context, url string) ([]models.Point, error) {
	return opensky.NewRequest(url).GetPlanes(ctx)
}

func (operation *OperationImpl) InsertPoint(ctx context.Context, point *models.Point) error {
	err := operation.repository.Create(point)

	if err == nil {
		log.Printf("Inserting point with `%s`", point.String())
		return nil
	}

	if err.Error() == models.ErrorPointalreadyExisted.Error() {
		log.Printf("Point `%s` already existed", point.String())
		return nil
	}

	return err
}
