package migrations

import (
	"errors"
	"fmt"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/models"
	geom "github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkt"
)

func CreateVilleray(cfg *config.Config) error {
	const PATH = "data/districts/villeray.geojson"

	db := cfg.DB()

	fmt.Println("... Creating villeray district")

	repository := models.NewDistrictRepository(cfg)
	villeray, err := repository.FindByName("villeray")

	if err != nil {
		return err
	}

	if villeray != nil {
		fmt.Printf("    Villeray already exists\n")
		return nil
	}

	district, err := models.NewDistrictFromJson("villeray", PATH)

	if err != nil {
		return err
	}

	var geometry geom.T
	if err := geojson.Unmarshal(district.Geometry, &geometry); err != nil {
		return err
	}

	multipolygon, ok := geometry.(*geom.MultiPolygon)
	if !ok {
		return errors.New("geometry is not a multipolygon")
	}

	geomStr, err := wkt.Marshal(multipolygon.SetSRID(4326))

	if err != nil {
		return err
	}

	db.Debug().Exec("INSERT INTO districts(name, geometry) VALUES ($1, ST_GeomFromText($2, 4326));", "villeray", geomStr)

	return db.Error
}
