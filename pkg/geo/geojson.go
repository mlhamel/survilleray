package geo

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/go-spatial/geom"
	"github.com/go-spatial/geom/encoding/geojson"
	"github.com/go-spatial/geom/encoding/wkt"
)

type Geojson struct {
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

func NewGeojsonFromValue(value string) (*Geojson, error) {
	var raw map[string]json.RawMessage

	if err := json.NewDecoder(strings.NewReader(value)).Decode(&raw); err != nil {
		return nil, err
	}

	return NewGeojson(raw["geometry"])
}

func NewGeojson(raw json.RawMessage) (*Geojson, error) {
	var value Geojson
	var geometry geojson.Geometry
	var buf bytes.Buffer

	if err := geometry.UnmarshalJSON(raw); err != nil {
		return nil, err
	}

	multipolygon, ok := geometry.Geometry.(geom.MultiPolygon)
	if !ok {
		return nil, errors.New("geometry is not a multipolygon")
	}

	encoder := wkt.NewDefaultEncoder(&buf)
	err := encoder.Encode(multipolygon)

	if err != nil {
		return nil, err
	}

	value.MultiPolygon = &multipolygon
	value.String = buf.String()

	return &value, nil
}

func NewGeojsonFromRawMultiPolygon(raw string) (*Geojson, error) {
	var value Geojson

	decoder := wkt.NewDecoder(strings.NewReader(raw))
	geometry, err := decoder.Decode()

	if err != nil {
		return nil, err
	}

	multipolygon, ok := geometry.(geom.MultiPolygon)
	if !ok {
		return nil, errors.New("geometry is not a multipolygon")
	}

	value.MultiPolygon = &multipolygon
	value.String = raw

	return &value, nil
}
