package spot

// Option configures the geometry delivered by NewShape. Options are applied in
// order; later options overwrite the Type and geometry field set by earlier
// ones.
type Option func(*Shape)

// WithPoint sets the Shape to a Point with the given [lon, lat] coordinates.
func WithPoint(coordinates []float64) Option {
	return func(s *Shape) {
		s.Type = TypePoint
		s.Point = NewPoint(coordinates)
	}
}

// WithMultiPoint sets the Shape to a MultiPoint holding each [lon, lat] position.
func WithMultiPoint(coordinates [][]float64) Option {
	return func(s *Shape) {
		s.Type = TypeMultiPoint
		s.MultiPoint = NewMultiPoint(coordinates)
	}
}

// WithLineString sets the Shape to a LineString made of [lon, lat] positions.
func WithLineString(coordinates [][]float64) Option {
	return func(s *Shape) {
		s.Type = TypeLineString
		s.LineString = NewLineString(coordinates)
	}
}

// WithMultiLineString sets the Shape to a MultiLineString, one [lon, lat]
// position array per line.
func WithMultiLineString(coordinates [][][]float64) Option {
	return func(s *Shape) {
		s.Type = TypeMultiLineString
		s.MultiLineString = NewMultiLineString(coordinates)
	}
}

// WithPolygon sets the Shape to a Polygon. The outer slice is one ring per
// linear ring (first the exterior, then each hole); every ring is a closed
// sequence of [lon, lat] positions.
func WithPolygon(coordinates [][][]float64) Option {
	return func(s *Shape) {
		s.Type = TypePolygon
		s.Polygon = NewPolygon(coordinates)
	}
}

// WithMultiPolygon sets the Shape to a MultiPolygon, one polygon per top-level
// element.
func WithMultiPolygon(coordinates [][][][]float64) Option {
	return func(s *Shape) {
		s.Type = TypeMultiPolygon
		s.MultiPolygon = NewMultiPolygon(coordinates)
	}
}

// WithGeometryCollection sets the Shape to a GeometryCollection of the given
// decoded geometries.
func WithGeometryCollection(geometries []Geometry) Option {
	return func(s *Shape) {
		s.Type = TypeGeometryCollection
		s.GeometryCollection = NewGeometryCollection(geometries)
	}
}

// WithEnvelope sets the Shape to an Envelope from its two corner coordinates.
func WithEnvelope(coordinates [][]float64) Option {
	return func(s *Shape) {
		s.Type = TypeEnvelope
		s.Envelope = NewEnvelope(coordinates)
	}
}

// WithCircle sets the Shape to a Circle with the given radius and center
// coordinates. radius defaults to meters unless it carries a unit suffix.
func WithCircle(radius string, coordinates []float64) Option {
	return func(s *Shape) {
		s.Type = TypeCircle
		s.Circle = NewCircle(radius, coordinates)
	}
}
