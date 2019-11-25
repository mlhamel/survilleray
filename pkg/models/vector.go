package models

import (
	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/runtime"
)

type Vector struct {
	gorm.Model
	Icao24   string  `gorm:"not null"`
	CallSign string  `gorm:"not null"`
	Country  string  `gorm:"not null"`
	Closed   bool    `gorm:"default:false"`
	Points   []Point `gorm:"many2many:vector_points"`
}

type VectorRepository interface {
	Find() ([]Vector, error)
	FindByPoint(*Point) ([]Vector, error)
	FindByCallSign(string) (*Vector, error)
	Insert(*Vector) error
	AppendPoints(*Vector, []Point) error
	Update(*Vector, ...interface{}) error
}

func NewVectorFromPoint(point *Point) *Vector {
	return &Vector{
		Icao24:   point.Icao24,
		CallSign: point.CallSign,
		Country:  point.OriginCountry,
		Points:   []Point{},
	}
}

func NewVectorRepository(context *runtime.Context) VectorRepository {
	return &vectoryRepository{context}
}

type vectoryRepository struct {
	context *runtime.Context
}

func (repository *vectoryRepository) Find() ([]Vector, error) {
	var vectors []Vector

	err := repository.context.Database().Find(&vectors).Error

	if err != nil {
		return nil, err
	}

	return vectors, nil
}

func (repository *vectoryRepository) FindByPoint(point *Point) ([]Vector, error) {
	var vectors []Vector

	err := repository.context.
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

	err := repository.context.
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

func (repository *vectoryRepository) Insert(vector *Vector) error {
	return repository.context.
		Database().
		Create(vector).Error
}

func (repository *vectoryRepository) AppendPoints(vector *Vector, points []Point) error {
	return repository.context.
		Database().
		Model(vector).
		Association("Points").
		Append(points).Error
}

func (repository *vectoryRepository) Update(vector *Vector, attrs ...interface{}) error {
	return repository.context.
		Database().
		Model(vector).
		Update(attrs...).Error
}
