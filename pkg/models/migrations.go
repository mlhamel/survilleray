package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/runtime"
)

func CreatePoint(context *runtime.Context) error {
	fmt.Println("... Creating point table")

	if context.Database().HasTable(&Point{}) {
		return nil
	}

	return context.Database().CreateTable(&Point{}).Error
}

func CreateDistrict(context *runtime.Context) error {
	fmt.Println("... Creating district table")

	if context.Database().HasTable(&District{}) {
		return nil
	}

	return context.Database().CreateTable(&District{}).Error
}

func EnablePostgis(context *runtime.Context) error {
	fmt.Println("... Enabling postgis extension")

	return context.Database().Exec("CREATE EXTENSION IF NOT EXISTS postgis").Error
}

func CreateVector(context *runtime.Context) error {
	fmt.Println("... Creating vector table")

	if context.Database().HasTable(&Vector{}) {
		fmt.Println("	Vector already exists")
		return nil
	}

	return context.Database().AutoMigrate(&Vector{}).Error
}

func CreateVilleray(context *runtime.Context) error {
	fmt.Println("... Creating villeray district")

	repository := NewDistrictRepository(context)
	villeray, err := repository.FindByName("villeray")

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}

	if villeray != nil {
		fmt.Printf("    Villeray already exists with ID: %d\n", villeray.ID)
		return nil
	}

	district, err := NewDistrictFromJson("villeray", VILLERAY)

	if err != nil {
		return err
	}

	query := "INSERT INTO districts(name, geometry) VALUES ($1, ST_GeomFromText($2, 4326));"

	return context.Database().Exec(query, "villeray", district.Geometry).Error
}
