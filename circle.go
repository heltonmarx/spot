package spot

import (
	"encoding/json"
)

// Circle is specified by a center point and a radius with units, defaulting to
// meters when no unit suffix is present.
type Circle struct {
	Type        string    `json:"type"`
	Radius      string    `json:"radius"`
	Coordinates []float64 `json:"coordinates"`
}

// NewCircle returns a Circle centered at the [lon, lat] coordinates with the
// given radius (e.g. "25m", "1km"). Panics if fewer than 2 coordinates are
// supplied or if radius is empty.
func NewCircle(radius string, coordinates []float64) *Circle {
	switch {
	case radius == "":
		panic("spot: Circle radius must not be empty")
	case len(coordinates) < 2:
		panic("spot: Circle requires at least 2 coordinates [lon, lat]")
	}
	return &Circle{
		Type:        TypeCircle,
		Radius:      radius,
		Coordinates: coordinates,
	}
}

// GeoType returns the geometry type constant for Circle.
func (m *Circle) GeoType() string { return TypeCircle }
func (m *Circle) isGeometry()     {}

// UnmarshalJSON decodes a circle geometry into m.
func (m *Circle) UnmarshalJSON(data []byte) error {
	raw := &rawGeometry{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	return m.decode(raw)
}

func (m *Circle) decode(raw *rawGeometry) error {
	m.Type = raw.Type
	m.Radius = raw.Radius
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
