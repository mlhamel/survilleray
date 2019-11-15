package models

import (
	"fmt"
	"strings"

	"github.com/go-spatial/geom"
	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/config"
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
}

type PointRepository interface {
	Find() ([]Point, error)
	Insert(*Point) error
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

func NewPointRepository(cfg *config.Config) PointRepository {
	return &pointRepository{cfg}
}

type pointRepository struct {
	cfg *config.Config
}

func (p *pointRepository) Find() ([]Point, error) {
	points := []Point{}

	err := p.cfg.DB().Find(&points).Error

	if err != nil {
		return nil, err
	}

	return points, nil
}

func (p *pointRepository) Insert(point *Point) error {
	return p.cfg.DB().Create(point).Error
}
