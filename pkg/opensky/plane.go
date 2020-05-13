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

const planeBasePath = "/api/flights/aircraft?icao=%s&begin=%d&end=%d"

type planeRequest struct {
	URL            string
	icao           string
	begin          time.Time
	end            time.Time
	logger         *zerolog.Logger
	results        *parsedPlaneRequest
	defaultTimeout time.Duration
}

type parsedPlaneRequest []map[string]interface{}

func NewPlaneRequest(url string, icao string, begin time.Time, end time.Time, logger *zerolog.Logger) *planeRequest {
	return &planeRequest{
		URL:            url,
		icao:           icao,
		begin:          begin,
		end:            end,
		logger:         logger,
		results:        &parsedPlaneRequest{},
		defaultTimeout: time.Second * 5,
	}
}

func (r *planeRequest) Run(ctx context.Context) error {
	var url = r.planePathURL()

	ctx, cancelFunc := context.WithTimeout(context.Background(), r.defaultTimeout)
	defer cancelFunc()

	r.logger.Info().Str("url", url).Msg("Getting data")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("Cannot reach opensky url at %s: %w", url, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	r.logger.Info().Str("status", resp.Status).Msg("Getting response")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	r.logger.Info().Str("body", string(body)).Msg("Parsing body")

	if err = json.Unmarshal([]byte(body), r.results); err != nil {
		return fmt.Errorf("Cannot unmarshall %s to json: %w", string(body), err)
	}

	return nil
}

func (r *planeRequest) Result() []models.Plane {
	var planes []models.Plane

	for i := 0; i < len(*r.results); i++ {
		var v = (*r.results)[i]

		plane := models.Plane{
			Icao24:                           v["icao24"].(string),
			FirstSeen:                        v["firstSeen"].(float64),
			EstDepartureAirport:              v["estDepartureAirport"].(string),
			LastSeen:                         v["lastSeen"].(float64),
			EstArrivalAirport:                v["estArrivalAirport"].(string),
			Callsign:                         v["callsign"].(string),
			EstDepartureAirportHorizDistance: v["estDepartureAirportHorizDistance"].(float64),
			EstDepartureAirportVertDistance:  v["estDepartureAirportVertDistance"].(float64),
			EstArrivalAirportHorizDistance:   v["estArrivalAirportHorizDistance"].(float64),
			EstArrivalAirportVertDistance:    v["estArrivalAirportVertDistance"].(float64),
			DepartureAirportCandidatesCount:  v["departureAirportCandidatesCount"].(float64),
			ArrivalAirportCandidatesCount:    v["arrivalAirportCandidatesCount"].(float64),
		}

		planes = append(planes, plane)
	}

	return planes
}

func (r *planeRequest) planePathURL() string {
	path := fmt.Sprintf(planeBasePath, r.icao, r.begin.Unix(), r.end.Unix())

	return fmt.Sprintf("%s/%s", r.URL, path)
}
