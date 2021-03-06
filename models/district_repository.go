package models

import "github.com/mlhamel/survilleray/pkg/config"

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

	err := d.cfg.Orm().
		Table("districts").
		Select("ID, name, created_at, updated_at, deleted_at, ST_AsText(geometry) as geometry").
		Find(&districts).Error

	if err != nil {
		return nil, err
	}

	return districts, nil
}

func (d *districtRepository) FindByName(name string) (*District, error) {
	var district District

	err := d.cfg.Orm().
		Table("districts").
		Select("ID, name, created_at, updated_at, deleted_at, ST_AsText(geometry) as geometry").
		Where("name = ?", name).
		First(&district).Error

	if err != nil {
		return nil, err
	}

	return &district, nil
}

func (d *districtRepository) Insert(district *District) error {
	query := "INSERT INTO districts(name, geometry) VALUES ($1, ST_GeomFromText($2, 4326));"

	return d.cfg.Orm().
		Exec(query, "villeray", district.Geometry).
		Error
}

func (d *districtRepository) AppendPoint(district *District, point *Point) error {
	d.cfg.Logger().Info().
		Str("point", point.Icao24).
		Str("district", district.Name).
		Msg("Inserting point in district")
	return d.cfg.
		Orm().
		Model(district).
		Association("Points").
		Append(point).Error
}
