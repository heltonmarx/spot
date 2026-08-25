package spot

import "encoding/json"

// LineString is an array of two or more positions.
type LineString struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

// NewLineString returns a LineString through the given [lon, lat] positions.
// Panics if fewer than 2 positions are supplied.
func NewLineString(coordinates [][]float64) *LineString {
	if len(coordinates) < 2 {
		panic("spot: LineString requires at least 2 positions")
	}
	return &LineString{
		Type:        TypeLineString,
		Coordinates: coordinates,
	}
}

// GeoType returns the geometry type constant for LineString.
func (m *LineString) GeoType() string { return TypeLineString }
func (m *LineString) isGeometry()     {}

// UnmarshalJSON decodes a GeoJSON linestring into m.
func (m *LineString) UnmarshalJSON(data []byte) error { return unmarshalGeometry(data, m) }

func (m *LineString) decode(raw *rawGeometry) error {
	m.Type = TypeLineString
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
