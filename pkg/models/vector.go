package models

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
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

func (v *Vector) BeforeSave() error {
	v.CallSign = strings.TrimSpace(v.CallSign)

	if v.CallSign == "" {
		return fmt.Errorf("CallSign cannot be empty")
	}

	return nil
}

func (v *Vector) String() string {
	return fmt.Sprintf("(%s, %s, %f)", v.Icao24, v.CallSign, v.LastContact)
}
