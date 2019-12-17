package acquisition

import (
	"context"

	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/config"
)

type app struct {
	cfg           *config.Config
	pointRepos    models.PointRepository
	districtRepos models.DistrictRepository
}

func NewApp(cfg *config.Config) *app {
	pointRepos := models.NewPointRepository(cfg)
	districtRepos := models.NewDistrictRepository(cfg)

	return &app{cfg, pointRepos, districtRepos}
}

func (a *app) Run(ctx context.Context) error {
	job := NewJob(a.cfg, a.pointRepos, a.districtRepos)

	return job.Run(ctx)
}
