package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-spatial/geom"
	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/rs/zerolog"
)

var ErrorPointalreadyExisted = errors.New("Point already existed")

// Point represent a flight point from Opensky
type Point struct {
	gorm.Model
	Icao24         string `gorm:"not null;unique_index:idx_icao24_callsign_lastcontact"`
	CallSign       string `gorm:"not null;unique_index:idx_icao24_callsign_lastcontact"`
	OriginCountry  string
	TimePosition   float64
	LastContact    float64 `gorm:"not null;unique_index:idx_icao24_callsign_lastcontact"`
	Longitude      float64
	Latitude       float64
	GeoAltitude    float64
	OnGround       bool
	Velocity       float64
	Heading        float64
	VerticalRate   float64
	Sensors        string
	BaroAltitude   float64
	Squawk         string
	Spi            bool
	PositionSource float64
	VectorizedAt   time.Time `gorm:"default:null"`
}

type PointRepository interface {
	Find() ([]Point, error)
	FindByIcao24(string) ([]Point, error)
	FindByVectorizedAt(*time.Time) ([]Point, error)
	Create(*Point) error
	Insert(*Point) error
	Update(*Point, ...interface{}) error
}

// BeforeSave is adding addional validations
func (p *Point) BeforeSave() error {
	p.CallSign = strings.TrimSpace(p.CallSign)

	if p.CallSign == "" {
		return fmt.Errorf("CallSign cannot be empty")
	}

	return nil
}

// String return the string representation of the point
func (p *Point) String() string {
	return fmt.Sprintf("(%s, %s, %f)", p.Icao24, p.CallSign, p.LastContact)
}

func (p *Point) Geography() *geom.Point {
	return &geom.Point{p.Longitude, p.Latitude}
}

func (p *Point) CheckOverlaps(district *District) (bool, error) {
	geojson, err := district.GeoJson()
	if err != nil {
		return false, err
	}

	extent, err := geom.NewExtentFromGeometry(geojson.MultiPolygon)
	if err != nil {
		return false, err
	}

	return extent.ContainsGeom(p.Geography())
}

func (p *Point) MarshalZerologObject(e *zerolog.Event) {
	e.Str("Icao24", p.Icao24).
		Str("CallSign", p.CallSign).
		Str("OriginCountry", p.OriginCountry).
		Float64("TimePosition", p.TimePosition).
		Float64("LastContact", p.LastContact).
		Float64("Longitude", p.Longitude).
		Float64("Latitude", p.Latitude).
		Float64("GeoAltitude", p.GeoAltitude).
		Bool("OnGround", p.OnGround).
		Float64("Velocity", p.Velocity).
		Float64("Heading", p.Heading).
		Float64("VerticalRate", p.VerticalRate).
		Str("Sensors", p.Sensors).
		Float64("BaroAltitude", p.BaroAltitude).
		Str("Squawk", p.Squawk).
		Bool("Spi", p.Spi).
		Float64("PositionSource", p.PositionSource).
		Time("VectorizedAt", p.VectorizedAt)
}

func NewPointRepository(cfg *config.Config) PointRepository {
	return &pointRepository{cfg}
}

type pointRepository struct {
	cfg *config.Config
}

func (repository *pointRepository) Find() ([]Point, error) {
	points := []Point{}

	err := repository.cfg.Orm().Find(&points).Error

	if err != nil {
		return nil, err
	}

	return points, nil
}

func (repository *pointRepository) FindByIcao24(icao24 string) ([]Point, error) {
	points := []Point{}

	err := repository.cfg.Orm().
		Debug().
		Where("icao24 = ?", icao24).
		Find(&points).Error

	if err != nil {
		return nil, err
	}

	return points, nil
}

func (repository *pointRepository) FindByVectorizedAt(vectorizedAt *time.Time) ([]Point, error) {
	points := []Point{}

	query := repository.cfg.Orm().Debug()

	if vectorizedAt == nil {
		query = query.Where("vectorized_at IS NULL")
	} else {
		query = query.Where("vectorized_at = ?", vectorizedAt)
	}

	err := query.Find(&points).Error

	if err != nil {
		return nil, err
	}

	return points, nil
}

func (repository *pointRepository) Create(point *Point) error {
	if !repository.cfg.Orm().NewRecord(point) {
		return ErrorPointalreadyExisted
	}
	return repository.cfg.Orm().Create(&point).Error
}

func (repository *pointRepository) Insert(point *Point) error {
	return repository.cfg.Orm().Create(point).Error
}

func (repository *pointRepository) Update(point *Point, attrs ...interface{}) error {
	return repository.cfg.
		Orm().
		Model(point).
		Update(attrs...).Error
}
