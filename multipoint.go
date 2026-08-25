package spot

import "encoding/json"

// MultiPoint is an array of positions.
type MultiPoint struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

// NewMultiPoint returns a MultiPoint holding each [lon, lat] position.
func NewMultiPoint(coordinates [][]float64) *MultiPoint {
	return &MultiPoint{
		Type:        TypeMultiPoint,
		Coordinates: coordinates,
	}
}

// GeoType returns the geometry type constant for MultiPoint.
func (m *MultiPoint) GeoType() string { return TypeMultiPoint }
func (m *MultiPoint) isGeometry()     {}

// UnmarshalJSON decodes a GeoJSON multipoint into m.
func (m *MultiPoint) UnmarshalJSON(data []byte) error { return unmarshalGeometry(data, m) }

func (m *MultiPoint) decode(raw *rawGeometry) error {
	m.Type = TypeMultiPoint
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
