package models

import (
	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/runtime"
)

func CreatePoint(context *runtime.Context) error {
	return context.Database().Debug().AutoMigrate(&Point{}).Error
}

func CreateDistrict(context *runtime.Context) error {
	return context.Database().Debug().AutoMigrate(&District{}).Error
}

func EnablePostgis(context *runtime.Context) error {
	return context.Database().Debug().Exec("CREATE EXTENSION IF NOT EXISTS postgis").Error
}

func CreateVector(context *runtime.Context) error {
	return context.Database().Debug().AutoMigrate(&Vector{}).Error
}

func CreateVilleray(context *runtime.Context) error {
	repository := NewDistrictRepository(context)
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

	return context.Database().Debug().Exec(query, "villeray", district.Geometry).Error
}
