package geo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGeojsonFromPath(t *testing.T) {
	value, err := NewGeojsonFromPath("villeray.geojson")
	assert.NoError(t, err)
	assert.NotNil(t, value)
	assert.Contains(t, value.String, "MULTIPOLYGON")
}
