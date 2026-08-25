package spot

import "encoding/json"

// Envelope represents a bounding rectangle by the coordinates of its upper
// left and lower right corners, each a [lon, lat] position.
type Envelope struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

// NewEnvelope returns an Envelope from its two corner coordinates.
// Panics if exactly two corner positions are not provided.
func NewEnvelope(coordinates [][]float64) *Envelope {
	if len(coordinates) != 2 {
		panic("spot: Envelope requires exactly 2 corner positions [[lon,lat],[lon,lat]]")
	}
	return &Envelope{
		Type:        TypeEnvelope,
		Coordinates: coordinates,
	}
}

// GeoType returns the geometry type constant for Envelope.
func (m *Envelope) GeoType() string { return TypeEnvelope }
func (m *Envelope) isGeometry()     {}

// UnmarshalJSON decodes an envelope geometry into m.
func (m *Envelope) UnmarshalJSON(data []byte) error { return unmarshalGeometry(data, m) }

func (m *Envelope) decode(raw *rawGeometry) error {
	m.Type = TypeEnvelope
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
