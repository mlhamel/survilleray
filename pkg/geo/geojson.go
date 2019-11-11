package geo

import (
	"encoding/json"
	"errors"
	"os"

	geom "github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkt"
)

type Geojson struct {
	Geometry     geom.T
	MultiPolygon *geom.MultiPolygon
	String       string
}

func NewGeojsonFromPath(path string) (*Geojson, error) {
	var raw map[string]json.RawMessage
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(file).Decode(&raw); err != nil {
		return nil, err
	}

	return NewGeojson(raw["geometry"])
}

func NewGeojson(raw json.RawMessage) (*Geojson, error) {
	var err error
	var value Geojson

	if err := geojson.Unmarshal(raw, &value.Geometry); err != nil {
		return nil, err
	}

	multipolygon, ok := value.Geometry.(*geom.MultiPolygon)
	if !ok {
		return nil, errors.New("geometry is not a multipolygon")
	}

	value.MultiPolygon = multipolygon
	value.String, err = wkt.Marshal(value.MultiPolygon.SetSRID(4326))

	if err != nil {
		return nil, err
	}

	return &value, nil
}
