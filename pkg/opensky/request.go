package opensky

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mlhamel/survilleray/models"
	"github.com/rs/zerolog/log"
)

// Request is used for requesting new data to OpenSky
type Request struct {
	url string
}

type parsedRequest struct {
	Time   int
	States []interface{}
}

// NewRequest is creating a new OpenSky request
func NewRequest(url string) *Request {
	return &Request{url: url}
}

// GetPlanes a request to OpenSky
func (r *Request) GetPlanes(ctx context.Context) (points []models.Point, e error) {
	parsedURL := fmt.Sprintf(r.url, 44, 47, -74, -72)

	log.Info().Str("url", parsedURL).Msg("Getting data")

	resp, err := http.Get(parsedURL)

	if err != nil {
		return points, err
	}

	defer resp.Body.Close()

	log.Info().Str("status", resp.Status).Msg("Getting response")

	body, err := ioutil.ReadAll(resp.Body)

	log.Info().Str("body", string(body)).Msg("Parsing body")

	if err != nil {
		return points, err
	}

	var results parsedRequest

	err = json.Unmarshal([]byte(body), &results)

	if err != nil {
		return points, err
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

		point := models.Point{
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

		points = append(points, point)
	}
	return points, nil
}
