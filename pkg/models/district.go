package models

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/config"
	geom "github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkt"
)

type District struct {
	gorm.Model
	Name     string          `gorm:"size:20;unique_index"`
	Geometry json.RawMessage `gorm:"type:geometry(MULTIPOLYGON, 4326)"`
}

type DistrictRepository interface {
	Find() ([]*District, error)
	FindByName(name string) (*District, error)
}

func NewDistrictRepository(cfg *config.Config) DistrictRepository {
	return &districtRepository{cfg}
}

type districtRepository struct {
	cfg *config.Config
}

func (d *districtRepository) Find() ([]*District, error) {
	var districts []*District

	errors := d.cfg.DB().Find(&districts).GetErrors()

	if len(errors) > 0 {
		return nil, errors[0]
	}

	return districts, nil
}

func (d *districtRepository) FindByName(name string) (*District, error) {
	var district District

	d.cfg.DB().Where("name = ?", name).First(&district)

	errors := d.cfg.DB().GetErrors()

	if len(errors) > 0 {
		return nil, errors[0]
	}

	return &district, nil
}

func NewDistrictFromJson(name string, path string) (*District, error) {
	var district District

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(file).Decode(&district); err != nil {
		return nil, err
	}

	district.Name = name

	return &district, nil
}

func (d *District) Multipolygon() (*geom.MultiPolygon, error) {
	var geometry geom.T

	if err := geojson.Unmarshal(d.Geometry, &geometry); err != nil {
		return nil, err
	}

	multipolygon, ok := geometry.(*geom.MultiPolygon)
	if !ok {
		return nil, errors.New("geometry is not a multipolygon")
	}

	return multipolygon, nil
}

func (d *District) Serialize() (string, error) {
	multipolygon, err := d.Multipolygon()

	if err != nil {
		return "", err
	}

	multipolygonAsString, err := wkt.Marshal(multipolygon)

	if err != nil {
		return "", err
	}

	return multipolygonAsString, nil
}
