package opensky

// ParsedVector hold the data about flight
type ParsedVector struct {
	Icao24         string
	CallSign       string
	OriginCountry  string
	TimePosition   float64
	LastContact    float64
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
