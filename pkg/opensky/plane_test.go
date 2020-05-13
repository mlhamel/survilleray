package opensky

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewPlaneRequest(t *testing.T) {
	var baseURL = "https://opensky-network.org/api/"
	var icao = "c07c71"
	var begin, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	var end, _ = time.Parse(time.RFC3339, "2006-01-02T15:05:05Z")

	log := zerolog.Nop().With().Logger()

	r := NewPlaneRequest(baseURL, icao, begin, end, &log)
	assert.Equal(t, baseURL, r.URL)
}

func TestPlaneRun(t *testing.T) {
	var icao = "c07c71"
	var begin, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	var end, _ = time.Parse(time.RFC3339, "2006-01-02T15:05:05Z")
	server := makeOpenskyPlaneServer()

	defer server.Close()

	log := zerolog.Nop().With().Logger()

	r := NewPlaneRequest(server.URL, icao, begin, end, &log)

	err := r.Run(context.Background())
	assert.NoError(t, err)
}

func makeOpenskyPlaneServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{
			map[string]interface{}{
				"arrivalAirportCandidatesCount":    3,
				"callsign":                         "DLH2VC  ",
				"departureAirportCandidatesCount":  1,
				"estArrivalAirport":                "ESSA",
				"estArrivalAirportHorizDistance":   7194,
				"estArrivalAirportVertDistance":    423,
				"estDepartureAirport":              "EDDF",
				"estDepartureAirportHorizDistance": 1462,
				"estDepartureAirportVertDistance":  49,
				"firstSeen":                        1517258040,
				"icao24":                           "3c675a",
				"lastSeen":                         1517263900,
			},
		})
	}))
}
