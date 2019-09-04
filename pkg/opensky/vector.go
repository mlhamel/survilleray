package opensky

// Vector hold the data about flight
type Vector struct {
	Icao24         string  `json:"icao24"`
	CallSign       string  `json:"callsign"`
	OriginCountry  string  `json:"origin_country"`
	TimePosition   float64 `json:"time_position"`
	LastContact    float64 `json:"last_contact"`
	Longitude      float64 `json:"longitude"`
	Latitude       float64 `json:"latitude"`
	GeoAltitude    float64 `json:"geo_altitude"`
	OnGround       bool    `json:"on_ground"`
	Velocity       float64 `json:"velocity"`
	Heading        float64 `json:"heading"`
	VerticalRate   float64 `json:"vertical_rate"`
	Sensors        string  `json:"sensors"`
	BaroAltitude   float64 `json:"baro_altitude"`
	Squawk         string  `json:"squawk"`
	Spi            bool    `json:"spi"`
	PositionSource float64 `json:"position_source"`
}
