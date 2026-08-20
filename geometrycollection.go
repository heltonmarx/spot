package spot

import (
	"encoding/json"
)

// GeometryCollection is a collection of other geometry objects, held as
// already-decoded concrete geometries in Geometries.
type GeometryCollection struct {
	Type       string `json:"type"`
	Geometries []any  `json:"geometries"`
}

// NewGeometryCollection returns a GeometryCollection of the given geometries.
func NewGeometryCollection(geometries []any) *GeometryCollection {
	return &GeometryCollection{
		Type:       TypeGeometryCollection,
		Geometries: geometries,
	}
}

// UnmarshalJSON decodes a GeoJSON geometry collection into m, decoding each
// member geometry according to its own "type".
func (m *GeometryCollection) UnmarshalJSON(data []byte) error {
	raw := &rawGeometry{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	return m.decode(raw)
}

func (m *GeometryCollection) decode(raw *rawGeometry) error {
	m.Type = raw.Type
	m.Geometries = make([]any, 0, len(raw.Geometries))
	for i := range raw.Geometries {
		geom, err := decodeGeometry(&raw.Geometries[i])
		if err != nil {
			return err
		}
		m.Geometries = append(m.Geometries, geom)
	}
	return nil
}
