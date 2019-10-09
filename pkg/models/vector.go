package models

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/config"
	geom "github.com/twpayne/go-geom"
)

// Vector represent a flight vector from Opensky
type Vector struct {
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

type VectorRepository interface {
	Find() ([]*Vector, error)
}

// BeforeSave is adding addional validations
func (v *Vector) BeforeSave() error {
	v.CallSign = strings.TrimSpace(v.CallSign)

	if v.CallSign == "" {
		return fmt.Errorf("CallSign cannot be empty")
	}

	return nil
}

// String return the string representation of the vector
func (v *Vector) String() string {
	return fmt.Sprintf("(%s, %s, %f)", v.Icao24, v.CallSign, v.LastContact)
}

func (v *Vector) Point() *geom.Point {
	return geom.NewPoint(geom.XY).MustSetCoords([]float64{v.Longitude, v.Latitude}).SetSRID(4326)
}

func (v *Vector) FindOverlaps(district *District) (bool, error) {
	polygons, err := district.Multipolygon()

	if err != nil {
		return false, err
	}

	return polygons.Bounds().
		OverlapsPoint(polygons.Layout(), v.Point().Coords()), nil
}

func NewVectorRepository(cfg *config.Config) VectorRepository {
	return &vectorRepository{cfg}
}

type vectorRepository struct {
	cfg *config.Config
}

func (v *vectorRepository) Find() ([]*Vector, error) {
	var vectors []*Vector

	errors := v.cfg.DB().Find(&vectors).GetErrors()

	if len(errors) > 0 {
		return nil, errors[0]
	}

	return vectors, nil
}
