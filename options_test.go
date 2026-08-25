package spot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewShapeOptions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		opts  []Option
		check func(t *testing.T, s *Shape)
	}{
		{
			name: "no options leaves shape empty",
			opts: nil,
			check: func(t *testing.T, s *Shape) {
				t.Helper()
				is := assert.New(t)
				is.Empty(s.Type)
				is.Nil(s.Point)
			},
		},
		{
			name: "with point",
			opts: []Option{WithPoint([]float64{100.0, 0.0})},
			check: func(t *testing.T, s *Shape) {
				t.Helper()
				is := assert.New(t)
				is.True(s.IsPoint())
				is.Equal(NewPoint([]float64{100.0, 0.0}), s.Point)
			},
		},
		{
			name: "with multipoint",
			opts: []Option{WithMultiPoint([][]float64{{100.0, 0.0}, {101.0, 1.0}})},
			check: func(t *testing.T, s *Shape) {
				t.Helper()
				is := assert.New(t)
				is.True(s.IsMultiPoint())
				is.Equal(NewMultiPoint([][]float64{{100.0, 0.0}, {101.0, 1.0}}), s.MultiPoint)
			},
		},
		{
			name: "with linestring",
			opts: []Option{WithLineString([][]float64{{100.0, 0.0}, {101.0, 1.0}})},
			check: func(t *testing.T, s *Shape) {
				t.Helper()
				is := assert.New(t)
				is.True(s.IsLineString())
				is.Equal(NewLineString([][]float64{{100.0, 0.0}, {101.0, 1.0}}), s.LineString)
			},
		},
		{
			name: "with multilinestring",
			opts: []Option{WithMultiLineString([][][]float64{{{100.0, 0.0}, {101.0, 1.0}}})},
			check: func(t *testing.T, s *Shape) {
				t.Helper()
				is := assert.New(t)
				is.True(s.IsMultiLineString())
				is.Equal(NewMultiLineString([][][]float64{{{100.0, 0.0}, {101.0, 1.0}}}), s.MultiLineString)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.check(t, NewShape(tt.opts...))
		})
	}
}

func TestNewShapeOptionsMore(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		opts  []Option
		check func(t *testing.T, s *Shape)
	}{
		{
			name: "with polygon",
			opts: []Option{WithPolygon([][][]float64{{{100.0, 0.0}, {101.0, 1.0}}})},
			check: func(t *testing.T, s *Shape) {
				t.Helper()
				is := assert.New(t)
				is.True(s.IsPolygon())
				is.Equal(NewPolygon([][][]float64{{{100.0, 0.0}, {101.0, 1.0}}}), s.Polygon)
			},
		},
		{
			name: "with multipolygon",
			opts: []Option{WithMultiPolygon([][][][]float64{{{{100.0, 0.0}}}})},
			check: func(t *testing.T, s *Shape) {
				t.Helper()
				is := assert.New(t)
				is.True(s.IsMultiPolygon())
				is.Equal(NewMultiPolygon([][][][]float64{{{{100.0, 0.0}}}}), s.MultiPolygon)
			},
		},
		{
			name: "with geometry collection",
			opts: []Option{WithGeometryCollection([]Geometry{NewPoint([]float64{100.0, 0.0})})},
			check: func(t *testing.T, s *Shape) {
				t.Helper()
				is := assert.New(t)
				is.True(s.IsGeometryCollection())
				is.Equal(NewGeometryCollection([]Geometry{NewPoint([]float64{100.0, 0.0})}), s.GeometryCollection)
			},
		},
		{
			name: "with envelope",
			opts: []Option{WithEnvelope([][]float64{{100.0, 1.0}, {101.0, 0.0}})},
			check: func(t *testing.T, s *Shape) {
				t.Helper()
				is := assert.New(t)
				is.True(s.IsEnvelope())
				is.Equal(NewEnvelope([][]float64{{100.0, 1.0}, {101.0, 0.0}}), s.Envelope)
			},
		},
		{
			name: "with circle",
			opts: []Option{WithCircle("25m", []float64{-109.874838, 44.439550})},
			check: func(t *testing.T, s *Shape) {
				t.Helper()
				is := assert.New(t)
				is.True(s.IsCircle())
				is.Equal(NewCircle("25m", []float64{-109.874838, 44.439550}), s.Circle)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.check(t, NewShape(tt.opts...))
		})
	}
}
