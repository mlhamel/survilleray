package models

type Plane struct {
	Icao24                           string
	FirstSeen                        float64
	EstDepartureAirport              string
	LastSeen                         float64
	EstArrivalAirport                string
	Callsign                         string
	EstDepartureAirportHorizDistance float64
	EstDepartureAirportVertDistance  float64
	EstArrivalAirportHorizDistance   float64
	EstArrivalAirportVertDistance    float64
	DepartureAirportCandidatesCount  float64
	ArrivalAirportCandidatesCount    float64
}
