package spot

// GeometryCollection is a collection of other geometry objects.
type GeometryCollection struct {
	Type       string     `json:"type"`
	Geometries []Geometry `json:"geometries"`
}

// NewGeometryCollection returns a GeometryCollection of the given geometries.
func NewGeometryCollection(geometries []Geometry) *GeometryCollection {
	return &GeometryCollection{
		Type:       TypeGeometryCollection,
		Geometries: geometries,
	}
}

// GeoType returns the geometry type constant for GeometryCollection.
func (m *GeometryCollection) GeoType() string { return TypeGeometryCollection }
func (m *GeometryCollection) isGeometry()     {}

// UnmarshalJSON decodes a GeoJSON geometry collection into m, decoding each
// member geometry according to its own "type".
func (m *GeometryCollection) UnmarshalJSON(data []byte) error { return unmarshalGeometry(data, m) }

func (m *GeometryCollection) decode(raw *rawGeometry) error {
	m.Type = TypeGeometryCollection
	m.Geometries = make([]Geometry, 0, len(raw.Geometries))
	for i := range raw.Geometries {
		geom, err := decodeGeometry(&raw.Geometries[i])
		if err != nil {
			return err
		}
		m.Geometries = append(m.Geometries, geom)
	}
	return nil
}
