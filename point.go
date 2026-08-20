package spot

import "encoding/json"

// Point is a single GeoJSON position.
type Point struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// NewPoint returns a Point positioned at the given [lon, lat] coordinates.
func NewPoint(coordinates []float64) *Point {
	return &Point{
		Type:        TypePoint,
		Coordinates: coordinates,
	}
}

// UnmarshalJSON decodes a GeoJSON point into m.
func (m *Point) UnmarshalJSON(data []byte) error {
	raw := &rawGeometry{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	return m.decode(raw)
}

func (m *Point) decode(raw *rawGeometry) error {
	m.Type = raw.Type
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
