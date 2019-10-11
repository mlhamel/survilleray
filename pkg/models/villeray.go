package models

import (
	"fmt"

	"github.com/mlhamel/survilleray/pkg/config"
)

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
