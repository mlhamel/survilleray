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
	Insert(*Vector) error
	AppendPoints(*Vector, []Point) error
}

func NewVectorRepository(context *runtime.Context) VectorRepository {
	return &vectoryRepository{context}
}

type vectoryRepository struct {
	context *runtime.Context
}

func (v *vectoryRepository) Find() ([]Vector, error) {
	var vectors []Vector

	err := v.context.Database().Find(&vectors).Error

	if err != nil {
		return nil, err
	}

	return vectors, nil
}

func (v *vectoryRepository) FindByPoint(point *Point) ([]Vector, error) {
	var vectors []Vector

	err := v.context.Database().Where(map[string]interface{}{
		"icao24":    point.Icao24,
		"call_sign": point.CallSign,
		"country":   point.OriginCountry,
		"closed":    false,
	}).Find(&vectors).Error

	if err != nil {
		return nil, err
	}

	return vectors, nil
}

func (v *vectoryRepository) Insert(vector *Vector) error {
	return v.context.Database().Create(vector).Error
}

func (v *vectoryRepository) AppendPoints(vector *Vector, points []Point) error {
	return v.context.Database().Model(vector).Association("Points").Append(points).Error
}
