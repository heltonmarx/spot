package spot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

//nolint:funlen // table-driven round-trip test carries twelve verbose coordinate fixtures; decomposing would harm cohesion
func TestGeometryRoundTrip(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		expected any
		actual   any
		raw      []byte
		pred     func(v any) bool
	}{
		{
			name:     "point",
			expected: NewPoint([]float64{100.0, 0.0}),
			actual:   &Point{},
			raw:      point,
		},
		{
			name:     "multipoint",
			expected: NewMultiPoint([][]float64{{100.0, 0.0}, {101.0, 1.0}}),
			actual:   &MultiPoint{},
			raw:      multiPoint,
		},
		{
			name:     "linestring",
			expected: NewLineString([][]float64{{100.0, 0.0}, {101.0, 1.0}}),
			actual:   &LineString{},
			raw:      lineString,
		},
		{
			name: "multilinestring",
			expected: NewMultiLineString([][][]float64{
				{{100.0, 0.0}, {101.0, 1.0}},
				{{102.0, 2.0}, {103.0, 3.0}},
			}),
			actual: &MultiLineString{},
			raw:    multiLineString,
		},
		{
			name: "polygon",
			expected: NewPolygon([][][]float64{
				{{100.0, 0.0}, {101.0, 0.0}, {101.0, 1.0}, {100.0, 1.0}, {100.0, 0.0}},
			}),
			actual: &Polygon{},
			raw:    polygon,
		},
		{
			name: "polygon with holes",
			expected: NewPolygon([][][]float64{
				{{100.0, 0.0}, {101.0, 0.0}, {101.0, 1.0}, {100.0, 1.0}, {100.0, 0.0}},
				{{100.2, 0.2}, {100.8, 0.2}, {100.8, 0.8}, {100.2, 0.8}, {100.2, 0.2}},
			}),
			actual: &Polygon{},
			raw:    polygonWithHoles,
		},
		{
			name: "multipolygon",
			expected: NewMultiPolygon([][][][]float64{
				{
					{{102.0, 2.0}, {103.0, 2.0}, {103.0, 3.0}, {102.0, 3.0}, {102.0, 2.0}},
				},
				{
					{{100.0, 0.0}, {101.0, 0.0}, {101.0, 1.0}, {100.0, 1.0}, {100.0, 0.0}},
					{{100.2, 0.2}, {100.8, 0.2}, {100.8, 0.8}, {100.2, 0.8}, {100.2, 0.2}},
				},
			}),
			actual: &MultiPolygon{},
			raw:    multiPolygon,
		},
		{
			name: "geometry collection",
			expected: NewGeometryCollection([]any{
				&Point{Type: TypePoint, Coordinates: []float64{100.0, 0.0}},
				&LineString{Type: TypeLineString, Coordinates: [][]float64{{101.0, 0.0}, {102.0, 1.0}}},
			}),
			actual: &GeometryCollection{},
			raw:    geometryCollection,
		},
		{
			name:     "shape",
			expected: &Shape{Type: TypePolygon, Polygon: NewPolygon([][][]float64{{{100.0, 0.0}, {101.0, 0.0}, {101.0, 1.0}, {100.0, 1.0}, {100.0, 0.0}}})},
			actual:   &Shape{},
			raw:      polygon,
			pred:     func(v any) bool { s, _ := v.(*Shape); return s.IsPolygon() },
		},
		{
			name:     "envelope",
			expected: NewEnvelope([][]float64{{100.0, 1.0}, {101.0, 0.0}}),
			actual:   &Envelope{},
			raw:      envelope,
		},
		{
			name:     "circle",
			expected: NewCircle("25m", []float64{-109.874838, 44.439550}),
			actual:   &Circle{},
			raw:      circle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt := tt
			t.Run("unmarshal", func(t *testing.T) {
				t.Parallel()
				require := require.New(t)
				require.NoError(json.Unmarshal(tt.raw, tt.actual))
				if tt.pred != nil {
					require.True(tt.pred(tt.actual))
				}
				require.Equal(tt.expected, tt.actual)
			})

			t.Run("marshal", func(t *testing.T) {
				t.Parallel()
				require := require.New(t)
				actual, err := json.Marshal(tt.expected)
				require.NoError(err)
				require.JSONEq(string(tt.raw), string(actual))
			})
		})
	}
}
