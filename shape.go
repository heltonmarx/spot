package spot

import (
	"encoding/json"
	"fmt"
)

// Shape is a discriminated container that can hold any single geometry: its
// Type field names the geometry kind, and exactly one concrete geometry field
// (Point, LineString, Polygon, etc.) is set to carry the value. Use it when the
// geometry type is not known at compile time — e.g. when unmarshaling arbitrary
// GeoJSON.
type Shape struct {
	Type               string              `json:"type"`
	Point              *Point              `json:"-"`
	MultiPoint         *MultiPoint         `json:"-"`
	LineString         *LineString         `json:"-"`
	MultiLineString    *MultiLineString    `json:"-"`
	Polygon            *Polygon            `json:"-"`
	MultiPolygon       *MultiPolygon       `json:"-"`
	GeometryCollection *GeometryCollection `json:"-"`
	Envelope           *Envelope           `json:"-"`
	Circle             *Circle             `json:"-"`
}

// NewShape builds an empty generic Shape, then applies each option in order to
// set its Type and the corresponding geometry field. Common choices are the
// WithPoint, WithLineString, WithPolygon, ... helpers. NewShape returns a Shape
// whose Type is empty until at least one option is supplied.
func NewShape(opts ...Option) *Shape {
	shape := &Shape{}
	for _, opt := range opts {
		opt(shape)
	}
	return shape
}

// IsPoint reports whether m holds a valid Point, i.e. its type is TypePoint
// and its Point field is set.
func (m *Shape) IsPoint() bool {
	return m.Type == TypePoint && m.Point != nil
}

// IsMultiPoint reports whether m holds a valid MultiPoint, i.e. its type is
// TypeMultiPoint and its MultiPoint field is set.
func (m *Shape) IsMultiPoint() bool {
	return m.Type == TypeMultiPoint && m.MultiPoint != nil
}

// IsLineString reports whether m holds a valid LineString, i.e. its type is
// TypeLineString and its LineString field is set.
func (m *Shape) IsLineString() bool {
	return m.Type == TypeLineString && m.LineString != nil
}

// IsMultiLineString reports whether m holds a valid MultiLineString, i.e. its
// type is TypeMultiLineString and its MultiLineString field is set.
func (m *Shape) IsMultiLineString() bool {
	return m.Type == TypeMultiLineString && m.MultiLineString != nil
}

// IsPolygon reports whether m holds a valid Polygon, i.e. its type is
// TypePolygon and its Polygon field is set.
func (m *Shape) IsPolygon() bool {
	return m.Type == TypePolygon && m.Polygon != nil
}

// IsMultiPolygon reports whether m holds a valid MultiPolygon, i.e. its type
// is TypeMultiPolygon and its MultiPolygon field is set.
func (m *Shape) IsMultiPolygon() bool {
	return m.Type == TypeMultiPolygon && m.MultiPolygon != nil
}

// IsGeometryCollection reports whether m holds a valid GeometryCollection, i.e.
// its type is TypeGeometryCollection and its GeometryCollection field is set.
func (m *Shape) IsGeometryCollection() bool {
	return m.Type == TypeGeometryCollection && m.GeometryCollection != nil
}

// IsEnvelope reports whether m holds a valid Envelope, i.e. its type is
// TypeEnvelope and its Envelope field is set.
func (m *Shape) IsEnvelope() bool {
	return m.Type == TypeEnvelope && m.Envelope != nil
}

// IsCircle reports whether m holds a valid Circle, i.e. its type is TypeCircle
// and its Circle field is set.
func (m *Shape) IsCircle() bool {
	return m.Type == TypeCircle && m.Circle != nil
}

// UnmarshalJSON decodes a GeoJSON geometry into m. It reads the "type" field
// and decodes the remaining fields into the matching concrete geometry,
// storing the result in the corresponding Shape field. Returns an error if the
// "type" is not a recognized geometry kind.
func (m *Shape) UnmarshalJSON(data []byte) error {
	raw := &rawGeometry{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	return m.decode(raw)
}

// MarshalJSON serializes the contained geometry as a GeoJSON object, emitting
// the raw fields of whichever concrete geometry is set. Returns an error if no
// geometry type is matched by m.
func (m *Shape) MarshalJSON() ([]byte, error) {
	switch {
	case m.IsPoint():
		return json.Marshal(m.Point)
	case m.IsMultiPoint():
		return json.Marshal(m.MultiPoint)
	case m.IsLineString():
		return json.Marshal(m.LineString)
	case m.IsMultiLineString():
		return json.Marshal(m.MultiLineString)
	case m.IsPolygon():
		return json.Marshal(m.Polygon)
	case m.IsMultiPolygon():
		return json.Marshal(m.MultiPolygon)
	case m.IsGeometryCollection():
		return json.Marshal(m.GeometryCollection)
	case m.IsEnvelope():
		return json.Marshal(m.Envelope)
	case m.IsCircle():
		return json.Marshal(m.Circle)
	default:
		return nil, fmt.Errorf("geo: unknown type `%s`", m.Type)
	}
}

func (m *Shape) decode(raw *rawGeometry) error {
	geom, err := decodeGeometry(raw)
	if err != nil {
		return err
	}
	m.Type = assignToShape(m, geom)
	return nil
}

// assignToShape stores g in the Shape field matching its concrete type, clears
// any geometry previously held by the shape, and returns the geometry type
// name. It is the single mapping from a concrete geometry to its Shape field,
// shared by Shape.UnmarshalJSON and the NewShape options where reusing a Shape
// must not leave stale fields behind.
func assignToShape(s *Shape, g any) string {
	s.Point = nil
	s.MultiPoint = nil
	s.LineString = nil
	s.MultiLineString = nil
	s.Polygon = nil
	s.MultiPolygon = nil
	s.GeometryCollection = nil
	s.Envelope = nil
	s.Circle = nil

	switch v := g.(type) {
	case *Point:
		s.Point = v
		return TypePoint
	case *MultiPoint:
		s.MultiPoint = v
		return TypeMultiPoint
	case *LineString:
		s.LineString = v
		return TypeLineString
	case *MultiLineString:
		s.MultiLineString = v
		return TypeMultiLineString
	case *Polygon:
		s.Polygon = v
		return TypePolygon
	case *MultiPolygon:
		s.MultiPolygon = v
		return TypeMultiPolygon
	case *GeometryCollection:
		s.GeometryCollection = v
		return TypeGeometryCollection
	case *Envelope:
		s.Envelope = v
		return TypeEnvelope
	case *Circle:
		s.Circle = v
		return TypeCircle
	}
	return ""
}
