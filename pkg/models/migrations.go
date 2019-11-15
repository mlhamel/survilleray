package models

import (
	"fmt"

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

	return db.CreateTable(&Vector{}).Error
}

func CreateVilleray(cfg *config.Config) error {
	const PATH = "data/districts/villeray.geojson"

	db := cfg.DB()

	fmt.Println("... Creating villeray district")

	repository := NewDistrictRepository(cfg)
	villeray, err := repository.FindByName("villeray")

	if err != nil {
		return err
	}

	if villeray != nil {
		fmt.Printf("    Villeray already exists\n")
		return nil
	}

	district, err := NewDistrictFromJson("villeray", PATH)

	query := "INSERT INTO districts(name, geometry) VALUES ($1, ST_GeomFromText($2, 4326));"

	return db.Exec(query, "villeray", district.Geometry).Error
}
