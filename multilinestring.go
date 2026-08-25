package spot

import "encoding/json"

// MultiLineString is an array of LineString coordinate arrays.
type MultiLineString struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

// NewMultiLineString returns a MultiLineString, one [lon, lat] position array
// per line.
func NewMultiLineString(coordinates [][][]float64) *MultiLineString {
	return &MultiLineString{
		Type:        TypeMultiLineString,
		Coordinates: coordinates,
	}
}

// GeoType returns the geometry type constant for MultiLineString.
func (m *MultiLineString) GeoType() string { return TypeMultiLineString }
func (m *MultiLineString) isGeometry()     {}

// UnmarshalJSON decodes a GeoJSON multilinestring into m.
func (m *MultiLineString) UnmarshalJSON(data []byte) error { return unmarshalGeometry(data, m) }

func (m *MultiLineString) decode(raw *rawGeometry) error {
	m.Type = TypeMultiLineString
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
