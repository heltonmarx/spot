package spot

import "encoding/json"

// LineString is an array of two or more positions.
type LineString struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

// NewLineString returns a LineString through the given [lon, lat] positions.
func NewLineString(coordinates [][]float64) *LineString {
	return &LineString{
		Type:        TypeLineString,
		Coordinates: coordinates,
	}
}

// UnmarshalJSON decodes a GeoJSON linestring into m.
func (m *LineString) UnmarshalJSON(data []byte) error {
	raw := &rawGeometry{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	return m.decode(raw)
}

func (m *LineString) decode(raw *rawGeometry) error {
	m.Type = raw.Type
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
