package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-spatial/geom"
	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/runtime"
)

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

func (p *Point) FindOverlaps(district *District) (bool, error) {
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

func NewPointRepository(context *runtime.Context) PointRepository {
	return &pointRepository{context}
}

type pointRepository struct {
	context *runtime.Context
}

func (repository *pointRepository) Find() ([]Point, error) {
	points := []Point{}

	err := repository.context.Database().Find(&points).Error

	if err != nil {
		return nil, err
	}

	return points, nil
}

func (repository *pointRepository) FindByIcao24(icao24 string) ([]Point, error) {
	points := []Point{}

	err := repository.context.
		Database().
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

	query := repository.context.Database().Debug()

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

func (repository *pointRepository) Insert(point *Point) error {
	return repository.context.Database().Create(point).Error
}

func (repository *pointRepository) Update(point *Point, attrs ...interface{}) error {
	return repository.context.
		Database().
		Model(point).
		Update(attrs...).Error
}
