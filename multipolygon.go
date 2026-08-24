package spot

import "encoding/json"

// MultiPolygon represents a GeoJSON object of multiple Polygons.
type MultiPolygon struct {
	Type        string          `json:"type"`
	Coordinates [][][][]float64 `json:"coordinates"`
}

// NewMultiPolygon returns a MultiPolygon, one Polygon per top-level element.
func NewMultiPolygon(coordinates [][][][]float64) *MultiPolygon {
	return &MultiPolygon{
		Type:        TypeMultiPolygon,
		Coordinates: coordinates,
	}
}

// GeoType returns the geometry type constant for MultiPolygon.
func (m *MultiPolygon) GeoType() string { return TypeMultiPolygon }
func (m *MultiPolygon) isGeometry()     {}

// UnmarshalJSON decodes a GeoJSON multipolygon into m.
func (m *MultiPolygon) UnmarshalJSON(data []byte) error {
	raw := &rawGeometry{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	return m.decode(raw)
}

func (m *MultiPolygon) decode(raw *rawGeometry) error {
	m.Type = raw.Type
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
