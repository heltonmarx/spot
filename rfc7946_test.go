package spot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRFC7946AllStandardTypes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		geometry Geometry
		wantType string
	}{
		{TypePoint, NewPoint([]float64{1, 2}), "Point"},
		{TypeMultiPoint, NewMultiPoint([][]float64{{1, 2}, {3, 4}}), "MultiPoint"},
		{TypeLineString, NewLineString([][]float64{{1, 2}, {3, 4}}), "LineString"},
		{TypeMultiLineString, NewMultiLineString([][][]float64{{{1, 2}, {3, 4}}}), "MultiLineString"},
		{TypePolygon, NewPolygon([][][]float64{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}}), "Polygon"},
		{TypeMultiPolygon, NewMultiPolygon([][][][]float64{{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}}}), "MultiPolygon"},
		{TypeGeometryCollection, NewGeometryCollection([]Geometry{}), "GeometryCollection"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data, err := json.Marshal(RFC7946(tt.geometry))
			require.NoError(t, err)

			var out struct {
				Type string `json:"type"`
			}
			require.NoError(t, json.Unmarshal(data, &out))
			require.Equal(t, tt.wantType, out.Type)
		})
	}
}

func TestRFC7946GeometryCollectionRecurses(t *testing.T) {
	t.Parallel()
	g := NewGeometryCollection([]Geometry{
		NewPoint([]float64{100.0, 0.0}),
		NewLineString([][]float64{{101.0, 0.0}, {102.0, 1.0}}),
	})
	data, err := json.Marshal(RFC7946(g))
	require.NoError(t, err)
	require.JSONEq(t, `{
		"type":"GeometryCollection",
		"geometries":[
			{"type":"Point","coordinates":[100.0,0.0]},
			{"type":"LineString","coordinates":[[101.0,0.0],[102.0,1.0]]}
		]
	}`, string(data))
}

func TestRFC7946ESExtensionsUnchanged(t *testing.T) {
	t.Parallel()
	t.Run("envelope", func(t *testing.T) {
		t.Parallel()
		data, err := json.Marshal(RFC7946(NewEnvelope([][]float64{{100.0, 1.0}, {101.0, 0.0}})))
		require.NoError(t, err)
		require.JSONEq(t, `{"type":"envelope","coordinates":[[100.0,1.0],[101.0,0.0]]}`, string(data))
	})
	t.Run("circle", func(t *testing.T) {
		t.Parallel()
		data, err := json.Marshal(RFC7946(NewCircle("25m", []float64{-109.874838, 44.43955})))
		require.NoError(t, err)
		require.JSONEq(t, `{"type":"circle","radius":"25m","coordinates":[-109.874838,44.43955]}`, string(data))
	})
}
