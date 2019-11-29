package geo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGeojsonFromValue(t *testing.T) {
	const multipolygon = `{
		"geometry": {
			"coordinates": [
		  		[[[102.0, 2.0], [103.0, 2.0], [103.0, 3.0], [102.0, 3.0], [102.0, 2.0]]],
		  		[[[100.0, 0.0], [101.0, 0.0], [101.0, 1.0], [100.0, 1.0], [100.0, 0.0]],
				[[100.2, 0.2], [100.8, 0.2], [100.8, 0.8], [100.2, 0.8], [100.2, 0.2]]]
			],
			"type": "MultiPolygon"
	  	}
	}`
	value, err := NewGeojsonFromValue(multipolygon)
	assert.NoError(t, err)
	assert.NotNil(t, value)
	assert.Contains(t, value.String, "MULTIPOLYGON")
}
