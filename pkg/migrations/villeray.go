package migrations

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/models"
	geom "github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkt"
)

func CreateVilleray(db *gorm.DB) error {
	const PATH = "data/districts/villeray.geojson"
	const NAME = "villeray"

	fmt.Println("... Creating villeray district")

	villeray, err := models.GetVilleray(db)

	if err != nil {
		return err
	}

	if villeray != nil {
		fmt.Printf("    Villeray already exists\n")
		return nil
	}

	district, err := models.NewDistrictFromJson(NAME, PATH)

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

	db.Debug().Exec("INSERT INTO districts(name, geometry) VALUES ($1, ST_GeomFromText($2, 4326));", NAME, geomStr)

	return db.Error
}
