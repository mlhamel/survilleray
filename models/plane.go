package models

import "github.com/jinzhu/gorm"

type Plane struct {
	gorm.Model
	Icao24                           string `gorm:"not null;unique_index:idx_icao24_callsign_lastseen"`
	FirstSeen                        float64
	EstDepartureAirport              string
	LastSeen                         float64 `gorm:"not null;unique_index:idx_icao24_callsign_lastseen"`
	EstArrivalAirport                string
	Callsign                         string `gorm:"not null;unique_index:idx_icao24_callsign_lastseen"`
	EstDepartureAirportHorizDistance float64
	EstDepartureAirportVertDistance  float64
	EstArrivalAirportHorizDistance   float64
	EstArrivalAirportVertDistance    float64
	DepartureAirportCandidatesCount  float64
	ArrivalAirportCandidatesCount    float64
}
