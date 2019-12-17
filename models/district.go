package models

import (
	"github.com/go-spatial/geom"
	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/geo"
)

type District struct {
	gorm.Model
	Name     string  `gorm:"size:20;unique_index"`
	Geometry string  `gorm:"type:geometry(MULTIPOLYGON, 4326)"`
	Points   []Point `gorm:"many2many:district_points"`
}

func NewDistrictFromJson(name string, value string) (*District, error) {
	var district District

	geojson, err := geo.NewGeojsonFromValue(value)

	if err != nil {
		return nil, err
	}

	district.Name = name
	district.Geometry = geojson.String

	return &district, nil
}

func (d *District) GeoJson() (*geo.Geojson, error) {
	geojson, err := geo.NewGeojsonFromRawMultiPolygon(d.Geometry)

	return geojson, err
}

func (d *District) Multipolygon() (*geom.MultiPolygon, error) {
	geojson, err := d.GeoJson()

	if err != nil {
		return nil, err
	}

	return geojson.MultiPolygon, nil
}

func (d *District) String() string {
	return d.Name
}
