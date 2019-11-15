package models

import (
	"github.com/go-spatial/geom"
	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/geo"
)

type District struct {
	gorm.Model
	Name     string `gorm:"size:20;unique_index"`
	Geometry string `gorm:"type:geometry(MULTIPOLYGON, 4326)"`
}

type DistrictRepository interface {
	Find() ([]*District, error)
	FindByName(name string) (*District, error)
	Insert(*District) error
}

func NewDistrictRepository(cfg *config.Config) DistrictRepository {
	return &districtRepository{cfg}
}

type districtRepository struct {
	cfg *config.Config
}

func (d *districtRepository) Find() ([]*District, error) {
	var districts []*District

	err := d.cfg.DB().
		Table("districts").
		Select("name, ST_AsText(geometry) as geometry").
		Find(&districts).Error

	if err != nil {
		return nil, err
	}

	return districts, nil
}

func (d *districtRepository) FindByName(name string) (*District, error) {
	var district District

	err := d.cfg.DB().Where("name = ?", name).First(&district).Error

	return &district, err
}

func (d *districtRepository) Insert(district *District) error {
	query := "INSERT INTO districts(name, geometry) VALUES ($1, ST_GeomFromText($2, 4326));"
	return d.cfg.DB().Exec(query, "villeray", district.Geometry).Error
}

func NewDistrictFromJson(name string, path string) (*District, error) {
	var district District

	value, err := geo.NewGeojsonFromPath(path)

	if err != nil {
		return nil, err
	}

	district.Name = name
	district.Geometry = value.String

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
