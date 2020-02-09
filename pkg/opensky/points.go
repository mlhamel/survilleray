package opensky

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mlhamel/survilleray/models"
	"github.com/rs/zerolog"
)

const pointBasePath = "/api/states/all?lamin=%d&lamax=%d&lomin=%d&lomax=%d"

type pointsRequest struct {
	URL     string
	logger  *zerolog.Logger
	results *parsedPointRequest
}

type parsedPointRequest struct {
	Time   interface{}   `json:"time"`
	States []interface{} `json:"states"`
}

func NewPointsRequest(url string, logger *zerolog.Logger) *pointsRequest {
	return &pointsRequest{URL: url, logger: logger, results: &parsedPointRequest{}}
}

func (r *pointsRequest) Run(ctx context.Context) error {
	r.logger.Info().Str("url", r.statePathURL()).Msg("Getting data")
	resp, err := http.Get(r.statePathURL())
	if err != nil {
		return fmt.Errorf("Cannot reach opensky url at %s: %w", r.statePathURL(), err)
	}
	defer resp.Body.Close()

	r.logger.Info().Str("status", resp.Status).Msg("Getting response")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Cannot read body: %w", err)
	}

	r.logger.Info().Str("body", string(body)).Msg("Parsing body")
	err = json.Unmarshal([]byte(body), r.results)
	if err != nil {
		return fmt.Errorf("Cannot unmarshall %s to json: %w", string(body), err)
	}

	return nil
}

func (r *pointsRequest) Result() []models.Point {
	var points []models.Point

	for i := 0; i < len(r.results.States); i++ {
		var v = r.results.States[i].([]interface{})

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
	return points
}

func (r *pointsRequest) statePathURL() string {
	path := fmt.Sprintf(pointBasePath, 44, 47, -74, -72)

	return fmt.Sprintf("%s/%s", r.URL, path)
}
