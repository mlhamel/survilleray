package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/rs/zerolog"
)

var ErrorVectorAlreadyExisted = errors.New("Point already existed")

type Vector struct {
	gorm.Model
	Icao24   string  `gorm:"not null;index:idx_icao24_callsign_country_closed"`
	CallSign string  `gorm:"not null;index:idx_callsign,idx_icao24_callsign_country_closed"`
	Country  string  `gorm:"not null"`
	Closed   bool    `gorm:"default:false;index:idx_closed,idx_icao24_callsign_country_closed"`
	Points   []Point `gorm:"many2many:vector_points"`
}

func NewVectorFromPoint(point *Point) *Vector {
	return &Vector{
		Icao24:   point.Icao24,
		CallSign: point.CallSign,
		Country:  point.OriginCountry,
		Points:   []Point{},
		Closed:   false,
	}
}

func (v *Vector) MarshalZerologObject(e *zerolog.Event) {
	e.Str("Icao24", v.Icao24).
		Str("CallSign", v.CallSign).
		Str("Country", v.Country).
		Bool("Closed", v.Closed)
}

type VectorRepository interface {
	Find() ([]Vector, error)
	FindByPoint(*Point) ([]Vector, error)
	FindByCallSign(string) (*Vector, error)
	FindByClosed(bool) ([]Vector, error)
	Create(*Vector) error
	Insert(*Vector) error
	AppendPoints(*Vector, []*Point) error
	Update(*Vector, ...interface{}) error
}

func NewVectorRepository(cfg *config.Config) VectorRepository {
	return &vectorRepository{cfg}
}

type vectorRepository struct {
	cfg *config.Config
}

func (r *vectorRepository) Find() ([]Vector, error) {
	var vectors []Vector

	err := r.cfg.Orm().Find(&vectors).Error

	if err != nil {
		return nil, err
	}

	return vectors, nil
}

func (r *vectorRepository) FindByPoint(point *Point) ([]Vector, error) {
	var vectors []Vector

	err := r.cfg.Orm().
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

func (r *vectorRepository) FindByCallSign(callsign string) (*Vector, error) {
	var vector Vector

	err := r.cfg.Orm().
		Where(map[string]interface{}{"call_sign": callsign, "closed": false}).
		Last(&vector).
		Error

	if err != nil {
		return nil, err
	}

	return &vector, err
}

func (r *vectorRepository) FindByClosed(closed bool) ([]Vector, error) {
	var vectors []Vector

	err := r.cfg.Orm().
		Where(map[string]interface{}{"closed": closed}).
		Find(&vectors).Error

	if err != nil {
		return nil, err
	}

	return vectors, nil
}

func (r *vectorRepository) Create(vector *Vector) error {
	if !r.cfg.Orm().NewRecord(vector) {
		return ErrorVectorAlreadyExisted
	}
	return r.cfg.Orm().Create(&vector).Error
}

func (r *vectorRepository) Insert(vector *Vector) error {
	return r.cfg.Orm().Create(vector).Error
}

func (r *vectorRepository) AppendPoints(vector *Vector, points []*Point) error {
	return r.cfg.Orm().
		Model(vector).
		Association("Points").
		Append(points).Error
}

func (r *vectorRepository) Update(vector *Vector, attrs ...interface{}) error {
	return r.cfg.Orm().
		Model(vector).
		Update(attrs...).Error
}
