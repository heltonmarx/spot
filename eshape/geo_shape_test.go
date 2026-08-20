package eshape

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	spot "github.com/heltonmarx/spot"
	"github.com/stretchr/testify/require"
)

func requireBody(t *testing.T, q *GeoShapeQuery) []byte {
	t.Helper()
	r, err := q.Body()
	require.NoError(t, err)
	b, err := io.ReadAll(r)
	require.NoError(t, err)
	return b
}

func assertJSONEq(t *testing.T, expected []byte, actual []byte) {
	t.Helper()
	var e, a any
	require.NoError(t, json.Unmarshal(expected, &e))
	require.NoError(t, json.Unmarshal(actual, &a))
	require.Equal(t, e, a)
}

func TestGeoShapeQueryIndexedShape(t *testing.T) {
	t.Parallel()
	is := &IndexedShape{Index: "shapes", ID: "deu", Path: "location"}
	expected := []byte(`{"geo_shape":{"location":{"indexed_shape":{"index":"shapes","id":"deu","path":"location"}}}}`)

	actual := requireBody(t, NewGeoShapeQuery("location").IndexedShape(is))
	assertJSONEq(t, expected, actual)
}

func TestGeoShapeQueryPoint(t *testing.T) {
	t.Parallel()
	shape := spot.NewShape(spot.WithPoint([]float64{13.400544, 52.530286}))
	expected := []byte(`{"geo_shape":{"location":{"shape":{"type":"point","coordinates":[13.400544,52.530286]},"relation":"within"}}}`)

	actual := requireBody(t, NewGeoShapeQuery("location").Shape(shape).Relation(RelationWithin))
	assertJSONEq(t, expected, actual)
}

func TestGeoShapeQueryPolygon(t *testing.T) {
	t.Parallel()
	shape := spot.NewShape(spot.WithPolygon([][][]float64{
		{{100.0, 0.0}, {101.0, 0.0}, {101.0, 1.0}, {100.0, 1.0}, {100.0, 0.0}},
	}))
	expected := []byte(`{"geo_shape":{"location":{"shape":{"type":"polygon","coordinates":[[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]},"relation":"within"}}}`)

	actual := requireBody(t, NewGeoShapeQuery("location").Shape(shape).Relation(RelationWithin))
	assertJSONEq(t, expected, actual)
}

func TestGeoShapeQueryCircle(t *testing.T) {
	t.Parallel()
	shape := spot.NewShape(spot.WithCircle("25m", []float64{-109.874838, 44.439550}))
	expected := []byte(`{"geo_shape":{"location":{"shape":{"type":"circle","radius":"25m","coordinates":[-109.874838,44.439550]}}}}`)

	actual := requireBody(t, NewGeoShapeQuery("location").Shape(shape))
	assertJSONEq(t, expected, actual)
}

func TestGeoShapeQueryErrors(t *testing.T) {
	t.Parallel()
	t.Run("empty path", func(t *testing.T) {
		t.Parallel()
		q := NewGeoShapeQuery("")
		q.shape = spot.NewShape(spot.WithPoint([]float64{1, 2}))
		_, err := q.MarshalJSON()
		require.Error(t, err)
	})

	t.Run("no shape", func(t *testing.T) {
		t.Parallel()
		q := NewGeoShapeQuery("location")
		_, err := q.MarshalJSON()
		require.Error(t, err)
	})
}

func TestGeoShapeQueryWithOfficialSearchRequest(t *testing.T) {
	t.Parallel()
	// Demonstrates the shape composes with the official client's raw search API.
	shape := spot.NewShape(spot.WithEnvelope([][]float64{{13, 53}, {14, 52}}))
	body, err := NewGeoShapeQuery("location").Shape(shape).Relation(RelationWithin).Body()
	require.NoError(t, err)

	// Build an actual esapi.SearchRequest using the body, exactly as a caller
	// would with the official client. This proves Body() satisfies the
	// io.Reader contract the client expects.
	req := esapi.SearchRequest{
		Index: []string{"places"},
		Body:  body,
	}
	if len(req.Index) != 1 || req.Body == nil {
		t.Fatal("search request not assembled from the geo_shape body")
	}

	var dst map[string]any
	require.NoError(t, json.NewDecoder(req.Body).Decode(&dst))
	require.Contains(t, dst, "geo_shape")
}
