package spot

import "encoding/json"

// Point is a single GeoJSON position.
type Point struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// NewPoint returns a Point positioned at the given [lon, lat] coordinates.
// Panics if fewer than 2 coordinates are supplied.
func NewPoint(coordinates []float64) *Point {
	if len(coordinates) < 2 {
		panic("spot: Point requires at least 2 coordinates [lon, lat]")
	}
	return &Point{
		Type:        TypePoint,
		Coordinates: coordinates,
	}
}

// GeoType returns the geometry type constant for Point.
func (m *Point) GeoType() string { return TypePoint }
func (m *Point) isGeometry()     {}

// UnmarshalJSON decodes a GeoJSON point into m.
func (m *Point) UnmarshalJSON(data []byte) error { return unmarshalGeometry(data, m) }

func (m *Point) decode(raw *rawGeometry) error {
	m.Type = TypePoint
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
