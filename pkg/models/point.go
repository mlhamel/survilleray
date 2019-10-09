package models

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/config"
	geom "github.com/twpayne/go-geom"
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
	Find() ([]*Point, error)
}

// BeforeSave is adding addional validations
func (v *Point) BeforeSave() error {
	v.CallSign = strings.TrimSpace(v.CallSign)

	if v.CallSign == "" {
		return fmt.Errorf("CallSign cannot be empty")
	}

	return nil
}

// String return the string representation of the point
func (v *Point) String() string {
	return fmt.Sprintf("(%s, %s, %f)", v.Icao24, v.CallSign, v.LastContact)
}

func (v *Point) Point() *geom.Point {
	return geom.NewPoint(geom.XY).MustSetCoords([]float64{v.Longitude, v.Latitude}).SetSRID(4326)
}

func (v *Point) FindOverlaps(district *District) (bool, error) {
	polygons, err := district.Multipolygon()

	if err != nil {
		return false, err
	}

	return polygons.Bounds().
		OverlapsPoint(polygons.Layout(), v.Point().Coords()), nil
}

func NewPointRepository(cfg *config.Config) PointRepository {
	return &pointRepository{cfg}
}

type pointRepository struct {
	cfg *config.Config
}

func (v *pointRepository) Find() ([]*Point, error) {
	var points []*Point

	errors := v.cfg.DB().Find(&points).GetErrors()

	if len(errors) > 0 {
		return nil, errors[0]
	}

	return points, nil
}
