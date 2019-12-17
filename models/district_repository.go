package models

import "github.com/mlhamel/survilleray/pkg/config"

import "log"

type DistrictRepository interface {
	Find() ([]*District, error)
	FindByName(name string) (*District, error)
	Insert(*District) error
	AppendPoint(*District, *Point) error
}

func NewDistrictRepository(cfg *config.Config) DistrictRepository {
	return &districtRepository{cfg}
}

type districtRepository struct {
	cfg *config.Config
}

func (d *districtRepository) Find() ([]*District, error) {
	var districts []*District

	err := d.cfg.Database().
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

	err := d.cfg.Database().
		Table("districts").
		Select("name, ST_AsText(geometry) as geometry").
		Where("name = ?", name).
		First(&district).Error

	if err != nil {
		return nil, err
	}

	return &district, nil
}

func (d *districtRepository) Insert(district *District) error {
	query := "INSERT INTO districts(name, geometry) VALUES ($1, ST_GeomFromText($2, 4326));"

	return d.cfg.Database().
		Exec(query, "villeray", district.Geometry).
		Error
}

func (d *districtRepository) AppendPoint(district *District, point *Point) error {
	log.Printf("Inserting point %s in district %s", point.Icao24, district.Name)
	return d.cfg.
		Database().
		Model(district).
		Association("Points").
		Append(point).Error
}
