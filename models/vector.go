package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/config"
)

var ErrorVectorAlreadyExisted = errors.New("Point already existed")

type Vector struct {
	gorm.Model
	Icao24   string  `gorm:"not null"`
	CallSign string  `gorm:"not null"`
	Country  string  `gorm:"not null"`
	Closed   bool    `gorm:"default:false"`
	Points   []Point `gorm:"many2many:vector_points"`
}

func NewVectorFromPoint(point *Point) *Vector {
	return &Vector{
		Icao24:   point.Icao24,
		CallSign: point.CallSign,
		Country:  point.OriginCountry,
		Points:   []Point{},
	}
}

type VectorRepository interface {
	Find() ([]Vector, error)
	FindByPoint(*Point) ([]Vector, error)
	FindByCallSign(string) (*Vector, error)
	Create(*Vector) error
	Insert(*Vector) error
	AppendPoints(*Vector, []Point) error
	Update(*Vector, ...interface{}) error
}

func NewVectorRepository(cfg *config.Config) VectorRepository {
	return &vectoryRepository{cfg}
}

type vectoryRepository struct {
	cfg *config.Config
}

func (repository *vectoryRepository) Find() ([]Vector, error) {
	var vectors []Vector

	err := repository.cfg.Database().Find(&vectors).Error

	if err != nil {
		return nil, err
	}

	return vectors, nil
}

func (repository *vectoryRepository) FindByPoint(point *Point) ([]Vector, error) {
	var vectors []Vector

	err := repository.cfg.
		Database().
		Where(map[string]interface{}{
			"icao24":    point.Icao24,
			"call_sign": point.CallSign,
			"country":   point.OriginCountry,
			"closed":    false}).
		Find(&vectors).Error

	if err != nil {
		return nil, err
	}

	return vectors, nil
}

func (repository *vectoryRepository) FindByCallSign(callsign string) (*Vector, error) {
	var vector Vector

	err := repository.cfg.
		Database().
		Where(map[string]interface{}{
			"call_sign": callsign,
			"closed":    false}).
		Last(&vector).
		Error

	if err != nil {
		return nil, err
	}

	return &vector, err
}

func (repository *vectoryRepository) Create(vector *Vector) error {
	if !repository.cfg.Database().NewRecord(vector) {
		return ErrorVectorAlreadyExisted
	}
	return repository.cfg.Database().Create(&vector).Error
}

func (repository *vectoryRepository) Insert(vector *Vector) error {
	return repository.cfg.
		Database().
		Create(vector).Error
}

func (repository *vectoryRepository) AppendPoints(vector *Vector, points []Point) error {
	return repository.cfg.
		Database().
		Model(vector).
		Association("Points").
		Append(points).Error
}

func (repository *vectoryRepository) Update(vector *Vector, attrs ...interface{}) error {
	return repository.cfg.
		Database().
		Model(vector).
		Update(attrs...).Error
}
