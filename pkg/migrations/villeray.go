package migrations

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/models"
	geom "github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
	"github.com/twpayne/go-geom/encoding/geojson"
)

func CreateVilleray(db *gorm.DB) error {
	const NAME = "villeray"
	const PATH = "data/districts/villeray.geojson"

	var count int
	var district models.District

	fmt.Println("... Creating villeray district")

	db.Where("name = ?", NAME).First(&district).Count(&count)

	if count > 0 {
		fmt.Printf("    Already exists (%d)\n", count)
		return nil
	}

	file, err := os.Open(PATH)

	if err != nil {
		return err
	}

	if err := json.NewDecoder(file).Decode(&district); err != nil {
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

	ewkbhexGeom, err := ewkbhex.Encode(multipolygon.SetSRID(4326), ewkbhex.XDR)

	if err != nil {
		return err
	}

	db.Debug().Exec("INSERT INTO districts(name, geometry) VALUES ($1, $2);", NAME, ewkbhexGeom)

	return db.Error
}
