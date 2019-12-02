package models

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/config"
)

func CreatePoint(ctx context.Context, cfg *config.Config) error {
	return cfg.Database().Debug().AutoMigrate(&Point{}).Error
}

func CreateDistrict(ctx context.Context, cfg *config.Config) error {
	return cfg.Database().Debug().AutoMigrate(&District{}).Error
}

func EnablePostgis(ctx context.Context, cfg *config.Config) error {
	return cfg.Database().Debug().Exec("CREATE EXTENSION IF NOT EXISTS postgis").Error
}

func CreateVector(ctx context.Context, cfg *config.Config) error {
	return cfg.Database().Debug().AutoMigrate(&Vector{}).Error
}

func CreateVilleray(ctx context.Context, cfg *config.Config) error {
	repository := NewDistrictRepository(cfg)
	villeray, err := repository.FindByName("villeray")

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}

	if villeray != nil {
		return nil
	}

	district, err := NewDistrictFromJson("villeray", VILLERAY)

	if err != nil {
		return err
	}

	query := "INSERT INTO districts(name, geometry) VALUES ($1, ST_GeomFromText($2, 4326));"

	return cfg.Database().Debug().Exec(query, "villeray", district.Geometry).Error
}
