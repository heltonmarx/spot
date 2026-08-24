package spot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShapeIsPredicates(t *testing.T) {
	t.Parallel()
	shape := &Shape{Type: TypePoint, Point: NewPoint([]float64{1, 2})}

	is := assert.New(t)
	is.True(shape.IsPoint())
	is.False(shape.IsMultiPoint())
	is.False(shape.IsLineString())
	is.False(shape.IsMultiLineString())
	is.False(shape.IsPolygon())
	is.False(shape.IsMultiPolygon())
	is.False(shape.IsGeometryCollection())
	is.False(shape.IsEnvelope())
	is.False(shape.IsCircle())
}

func TestShapeIsPredicatesRequireMatchingField(t *testing.T) {
	t.Parallel()
	// A shape whose Type and field disagree (e.g. built by direct struct
	// literal) must not report itself as valid, since the type and the set
	// field must line up.
	shape := &Shape{Type: TypePoint, Polygon: NewPolygon([][][]float64{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}})}

	is := assert.New(t)
	is.False(shape.IsPoint())
	is.False(shape.IsPolygon())
}

func TestShapeUnmarshalDispatchesToConcreteField(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		raw      string
		typeName string
		check    func(t *testing.T, s *Shape)
	}{
		{
			name:     "multipoint",
			raw:      `{"type":"multipoint","coordinates":[[1,2],[3,4]]}`,
			typeName: TypeMultiPoint,
			check: func(t *testing.T, s *Shape) {
				assert.True(t, s.IsMultiPoint())
				assert.Equal(t, [][]float64{{1, 2}, {3, 4}}, s.MultiPoint.Coordinates)
			},
		},
		{
			name:     "multilinestring",
			raw:      `{"type":"multilinestring","coordinates":[[[1,2],[3,4]]]}`,
			typeName: TypeMultiLineString,
			check: func(t *testing.T, s *Shape) {
				assert.True(t, s.IsMultiLineString())
			},
		},
		{
			name:     "multipolygon",
			raw:      `{"type":"multipolygon","coordinates":[[[[0,0],[1,0],[1,1],[0,0]]]]}`,
			typeName: TypeMultiPolygon,
			check: func(t *testing.T, s *Shape) {
				assert.True(t, s.IsMultiPolygon())
			},
		},
		{
			name:     "geometrycollection",
			raw:      `{"type":"geometrycollection","geometries":[{"type":"point","coordinates":[1,2]}]}`,
			typeName: TypeGeometryCollection,
			check: func(t *testing.T, s *Shape) {
				assert.True(t, s.IsGeometryCollection())
				is := assert.New(t)
				is.Len(s.GeometryCollection.Geometries, 1)
			},
		},
		{
			name:     "envelope",
			raw:      `{"type":"envelope","coordinates":[[100,1],[101,0]]}`,
			typeName: TypeEnvelope,
			check: func(t *testing.T, s *Shape) {
				assert.True(t, s.IsEnvelope())
			},
		},
		{
			name:     "circle",
			raw:      `{"type":"circle","radius":"25m","coordinates":[-109.87,44.43]}`,
			typeName: TypeCircle,
			check: func(t *testing.T, s *Shape) {
				assert.True(t, s.IsCircle())
				assert.Equal(t, "25m", s.Circle.Radius)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			shape := &Shape{}
			require := require.New(t)
			require.NoError(json.Unmarshal([]byte(tt.raw), shape))
			require.Equal(tt.typeName, shape.Type)
			tt.check(t, shape)
		})
	}
}

func TestShapeUnmarshalUnknownType(t *testing.T) {
	t.Parallel()
	shape := &Shape{}
	err := json.Unmarshal([]byte(`{"type":"teapot"}`), shape)
	require.Error(t, err)

	require := require.New(t)
	require.Contains(err.Error(), "unknown type `teapot`")
}

func TestShapeUnmarshalInvalidJSON(t *testing.T) {
	t.Parallel()
	shape := &Shape{}
	err := json.Unmarshal([]byte(`{"type":"point","coordinates":`), shape)
	require.Error(t, err)
}

func TestShapeUnmarshalCoordinateTypeMismatch(t *testing.T) {
	t.Parallel()
	// A shape dispatching to a point whose coordinates carry an object fails at
	// the concrete type's decode, and the error must propagate out of Shape.
	shape := &Shape{}
	err := json.Unmarshal([]byte(`{"type":"point","coordinates":{"x":1}}`), shape)
	require.Error(t, err)
}

func TestShapeMarshalUnknownType(t *testing.T) {
	t.Parallel()
	shape := &Shape{Type: "teapot"}
	_, err := json.Marshal(shape)
	require.Error(t, err)
}

func TestShapeMarshalEmptyType(t *testing.T) {
	t.Parallel()
	shape := &Shape{}
	_, err := json.Marshal(shape)
	require.Error(t, err)
}

func TestShapeMarshalRoundTripForAllTypes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		shape    *Shape
		expected string
	}{
		{
			name:     "multipoint",
			shape:    NewShape(WithMultiPoint([][]float64{{1, 2}, {3, 4}})),
			expected: `{"type":"multipoint","coordinates":[[1,2],[3,4]]}`,
		},
		{
			name:     "multilinestring",
			shape:    NewShape(WithMultiLineString([][][]float64{{{1, 2}, {3, 4}}})),
			expected: `{"type":"multilinestring","coordinates":[[[1,2],[3,4]]]}`,
		},
		{
			name:     "multipolygon",
			shape:    NewShape(WithMultiPolygon([][][][]float64{{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}}})),
			expected: `{"type":"multipolygon","coordinates":[[[[0,0],[1,0],[1,1],[0,0]]]]}`,
		},
		{
			name:     "geometrycollection",
			shape:    NewShape(WithGeometryCollection([]Geometry{NewPoint([]float64{1, 2})})),
			expected: `{"type":"geometrycollection","geometries":[{"type":"point","coordinates":[1,2]}]}`,
		},
		{
			name:     "envelope",
			shape:    NewShape(WithEnvelope([][]float64{{100, 1}, {101, 0}})),
			expected: `{"type":"envelope","coordinates":[[100,1],[101,0]]}`,
		},
		{
			name:     "circle",
			shape:    NewShape(WithCircle("25m", []float64{-109.87, 44.43})),
			expected: `{"type":"circle","radius":"25m","coordinates":[-109.87,44.43]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := json.Marshal(tt.shape)
			require := require.New(t)
			require.NoError(err)
			require.JSONEq(tt.expected, string(actual))
		})
	}
}
