package spot

import "encoding/json"

// Envelope represents a bounding rectangle by the coordinates of its upper
// left and lower right corners, each a [lon, lat] position.
type Envelope struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

// NewEnvelope returns an Envelope from its two corner coordinates.
func NewEnvelope(coordinates [][]float64) *Envelope {
	return &Envelope{
		Type:        TypeEnvelope,
		Coordinates: coordinates,
	}
}

// UnmarshalJSON decodes an envelope geometry into m.
func (m *Envelope) UnmarshalJSON(data []byte) error {
	raw := &rawGeometry{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	return m.decode(raw)
}

func (m *Envelope) decode(raw *rawGeometry) error {
	m.Type = raw.Type
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
