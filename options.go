package spot

// Option configures the geometry delivered by NewShape. Options are applied in
// order; later options overwrite the Type and geometry field set by earlier
// ones.
type Option func(*Shape)

// applyOption applies geometry to s: it records the geometry's type and stores
// it in the field matching its concrete type, clearing any geometry previously
// held by the shape. Later options therefore overwrite the geometry set by
// earlier ones.
func applyOption(s *Shape, g any) {
	s.Type = assignToShape(s, g)
}

// WithPoint sets the Shape to a Point with the given [lon, lat] coordinates.
func WithPoint(coordinates []float64) Option {
	return func(s *Shape) {
		applyOption(s, NewPoint(coordinates))
	}
}

// WithMultiPoint sets the Shape to a MultiPoint holding each [lon, lat] position.
func WithMultiPoint(coordinates [][]float64) Option {
	return func(s *Shape) {
		applyOption(s, NewMultiPoint(coordinates))
	}
}

// WithLineString sets the Shape to a LineString made of [lon, lat] positions.
func WithLineString(coordinates [][]float64) Option {
	return func(s *Shape) {
		applyOption(s, NewLineString(coordinates))
	}
}

// WithMultiLineString sets the Shape to a MultiLineString, one [lon, lat]
// position array per line.
func WithMultiLineString(coordinates [][][]float64) Option {
	return func(s *Shape) {
		applyOption(s, NewMultiLineString(coordinates))
	}
}

// WithPolygon sets the Shape to a Polygon. The outer slice is one ring per
// linear ring (first the exterior, then each hole); every ring is a closed
// sequence of [lon, lat] positions.
func WithPolygon(coordinates [][][]float64) Option {
	return func(s *Shape) {
		applyOption(s, NewPolygon(coordinates))
	}
}

// WithMultiPolygon sets the Shape to a MultiPolygon, one polygon per top-level
// element.
func WithMultiPolygon(coordinates [][][][]float64) Option {
	return func(s *Shape) {
		applyOption(s, NewMultiPolygon(coordinates))
	}
}

// WithGeometryCollection sets the Shape to a GeometryCollection of the given
// decoded geometries.
func WithGeometryCollection(geometries []any) Option {
	return func(s *Shape) {
		applyOption(s, NewGeometryCollection(geometries))
	}
}

// WithEnvelope sets the Shape to an Envelope from its two corner coordinates.
func WithEnvelope(coordinates [][]float64) Option {
	return func(s *Shape) {
		applyOption(s, NewEnvelope(coordinates))
	}
}

// WithCircle sets the Shape to a Circle with the given radius and center
// coordinates. radius defaults to meters unless it carries a unit suffix.
func WithCircle(radius string, coordinates []float64) Option {
	return func(s *Shape) {
		applyOption(s, NewCircle(radius, coordinates))
	}
}
