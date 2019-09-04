package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const base_url = "https://opensky-network.org/api/states/all?lamin=%d&lamax=%d&lomin=%d&lomax=%d"

type Input struct {
	Time   int
	States []interface{}
}

type StateResponse struct {
	Time   int           `json:"time"`
	States []StateVector `json:"states"`
}

type StateVector struct {
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

func main() {
	parsed_url := fmt.Sprintf(base_url, 44, 47, -74, -72)
	resp, err := http.Get(parsed_url)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var results Input

	err = json.Unmarshal([]byte(body), &results)

	if err != nil {
		panic(err)
	}

	for i := 0; i < len(results.States); i++ {
		var v = results.States[i].([]interface{})

		var longitude float64
		var latitude float64
		var geoAltitude float64
		var velocity float64
		var verticalRate float64
		var sensors string
		var baroAltitude float64
		var squawk string

		if v[5] != nil {
			longitude = v[5].(float64)
		}

		if v[6] != nil {
			latitude = v[6].(float64)
		}

		if v[7] != nil {
			geoAltitude = v[7].(float64)
		}

		if v[9] != nil {
			velocity = v[9].(float64)
		}

		if v[11] != nil {
			verticalRate = v[11].(float64)
		}

		if v[12] != nil {
			sensors = v[12].(string)
		}

		if v[13] != nil {
			baroAltitude = v[13].(float64)
		}

		if v[14] != nil {
			squawk = v[14].(string)
		}

		vector := StateVector{
			Icao24:         v[0].(string),
			CallSign:       v[1].(string),
			OriginCountry:  v[2].(string),
			TimePosition:   v[3].(float64),
			LastContact:    v[4].(float64),
			Longitude:      longitude,
			Latitude:       latitude,
			GeoAltitude:    geoAltitude,
			OnGround:       v[8].(bool),
			Velocity:       velocity,
			Heading:        v[10].(float64),
			VerticalRate:   verticalRate,
			Sensors:        sensors,
			BaroAltitude:   baroAltitude,
			Squawk:         squawk,
			Spi:            v[15].(bool),
			PositionSource: v[16].(float64),
		}

		fmt.Println(vector)
	}
}
