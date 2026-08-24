package spot

import "encoding/json"

// Polygon is an object consisting of one or more linear rings: the first ring
// is the exterior boundary and each subsequent ring is a hole. Rings are closed
// [lon, lat] position sequences (first and last position equal).
type Polygon struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

// NewPolygon returns a Polygon from its array of linear rings.
func NewPolygon(coordinates [][][]float64) *Polygon {
	return &Polygon{
		Type:        TypePolygon,
		Coordinates: coordinates,
	}
}

// GeoType returns the geometry type constant for Polygon.
func (m *Polygon) GeoType() string { return TypePolygon }
func (m *Polygon) isGeometry()     {}

// UnmarshalJSON decodes a GeoJSON polygon into m.
func (m *Polygon) UnmarshalJSON(data []byte) error {
	raw := &rawGeometry{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	return m.decode(raw)
}

func (m *Polygon) decode(raw *rawGeometry) error {
	m.Type = raw.Type
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
