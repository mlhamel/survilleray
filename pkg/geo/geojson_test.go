package geo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGeojsonFromPath(t *testing.T) {
	value, err := NewGeojsonFromPath("villeray.geojson")
	assert.NoError(t, err)
	assert.NotNil(t, value)
	assert.Equal(t, 4326, value.Geometry.SRID())
	assert.Equal(t, float64(-0.001898821957849052), value.MultiPolygon.Area())
	assert.Contains(t, value.String, "MULTIPOLYGON")
}
