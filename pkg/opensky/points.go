package opensky

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mlhamel/survilleray/models"
	"github.com/rs/zerolog"
)

const pointBasePath = "/api/states/all?lamin=%d&lamax=%d&lomin=%d&lomax=%d"

type pointsRequest struct {
	URL            string
	logger         *zerolog.Logger
	results        *parsedPointRequest
	defaultTimeout time.Duration
}

type parsedPointRequest struct {
	Time   interface{}   `json:"time"`
	States []interface{} `json:"states"`
}

func NewPointsRequest(url string, logger *zerolog.Logger) *pointsRequest {
	return &pointsRequest{
		URL:            url,
		logger:         logger,
		results:        &parsedPointRequest{},
		defaultTimeout: time.Second * 5,
	}
}

func (r *pointsRequest) Run(ctx context.Context) error {
	var url = r.statePathURL()

	ctx, cancelFunc := context.WithTimeout(context.Background(), r.defaultTimeout)
	defer cancelFunc()

	r.logger.Info().Str("url", url).Msg("Getting data")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("Cannot reach opensky url at %s: %w", url, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Cannot reach opensky url at %s: %w", url, err)
	}
	defer resp.Body.Close()

	r.logger.Info().Str("status", resp.Status).Msg("Getting response")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Cannot read body: %w", err)
	}

	r.logger.Info().Str("body", string(body)).Msg("Parsing body")

	if err = json.Unmarshal([]byte(body), r.results); err != nil {
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
