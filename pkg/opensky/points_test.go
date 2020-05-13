package opensky

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"net/http/httptest"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewPointRequest(t *testing.T) {
	var baseURL = "https://opensky-network.org/api/"

	log := zerolog.Nop().With().Logger()

	r := NewPointsRequest(baseURL, &log)
	assert.Equal(t, baseURL, r.URL)
}

func TestPointsRun(t *testing.T) {
	server := makeOpenskyPointServer()

	defer server.Close()

	log := zerolog.Nop().With().Logger()
	r := NewPointsRequest(server.URL, &log)

	err := r.Run(context.Background())
	assert.NoError(t, err)

	ps := r.Result()
	assert.Equal(t, 1, len(ps))

	p := ps[0]
	assert.Equal(t, "e8027e", p.Icao24)
}

func makeOpenskyPointServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"time": 1580522020,
			"states": []interface{}{
				[]interface{}{
					"e8027e",
					"LPE2225 ",
					"Chile",
					1580521951,
					1580521965,
					-77.138,
					-11.9739,
					213.36,
					false,
					66.95,
					154.03,
					-3.58,
					nil,
					281.94,
					"2675",
					false,
					0,
				},
			},
		})
	}))
}
