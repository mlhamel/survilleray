package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/config"
)

func CreatePoint(cfg *config.Config) error {
	fmt.Println("... Creating point table")

	db := cfg.DB()

	if db.HasTable(&Point{}) {
		return nil
	}

	db.CreateTable(&Point{})

	return db.Error
}

func CreateDistrict(cfg *config.Config) error {
	fmt.Println("... Creating district table")

	db := cfg.DB()

	if db.HasTable(&District{}) {
		return nil
	}

	db.CreateTable(&District{})

	return db.Error
}

func EnablePostgis(cfg *config.Config) error {
	fmt.Println("... Enabling postgis extension")

	db := cfg.DB()

	db.Exec("CREATE EXTENSION IF NOT EXISTS postgis")

	return db.Error
}

func CreateVector(cfg *config.Config) error {
	fmt.Println("... Creating vector table")

	db := cfg.DB()

	if db.HasTable(&Vector{}) {
		fmt.Println("	Vector already exists")
		return nil
	}

	return db.AutoMigrate(&Vector{}).Error
}

func CreateVilleray(cfg *config.Config) error {
	db := cfg.DB()

	fmt.Println("... Creating villeray district")

	repository := NewDistrictRepository(cfg)
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

	return db.Exec(query, "villeray", district.Geometry).Error
}
