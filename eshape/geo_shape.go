// Package eshape adapts the author's GeoShape Query implementation to Elastic's
// official Go client (elastic/go-elasticsearch). The GeoJSON geometry modeling
// lives in this module's spot package; this package builds the ES `geo_shape`
// query DSL so it can be sent through the low-level/raw API.
package eshape

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	spot "github.com/heltonmarx/spot"
)

// Relation is one of the geo-shape spatial relation operators accepted at
// search time (intersects, disjoint, within, contains).
const (
	RelationIntersects = "intersects"
	RelationDisjoint   = "disjoint"
	RelationWithin     = "within"
	RelationContains   = "contains"
)

// GeoShapeQuery finds documents whose shape field matches a shape using the
// given spatial relation. It mirrors the builder API of the olivere/elastic
// GeoShapeQuery, but emits plain JSON suitable for the official client's raw
// request API.
type GeoShapeQuery struct {
	path         string
	relation     string
	shape        *spot.Shape
	indexedShape *IndexedShape
}

// IndexedShape references a shape that has already been indexed in another
// index and/or type, so the query doesn't have to inline the geometry.
type IndexedShape struct {
	Index   string `json:"index,omitempty"`
	ID      string `json:"id,omitempty"`
	Path    string `json:"path,omitempty"`
	Routing string `json:"routing,omitempty"`
}

// NewGeoShapeQuery creates a geo_shape query on the given field path, e.g.
// "location".
func NewGeoShapeQuery(path string) *GeoShapeQuery {
	return &GeoShapeQuery{path: path}
}

// IndexedShape points the query at a pre-indexed shape instead of an inline
// geometry.
func (q *GeoShapeQuery) IndexedShape(is *IndexedShape) *GeoShapeQuery {
	q.indexedShape = is
	return q
}

// Shape sets an inline shape built with go-spot.
func (q *GeoShapeQuery) Shape(shape *spot.Shape) *GeoShapeQuery {
	q.shape = shape
	return q
}

// Relation sets the spatial relation operator (see the Relation* constants).
func (q *GeoShapeQuery) Relation(relation string) *GeoShapeQuery {
	q.relation = relation
	return q
}

// MarshalJSON emits the `geo_shape` query DSL.
//
//	e.g. {"geo_shape":{"location":{"shape":{"type":"point","coordinates":[...]},"relation":"within"}}}
func (q *GeoShapeQuery) MarshalJSON() ([]byte, error) {
	if q.path == "" {
		return nil, errors.New("eshape: geo_shape query path must not be empty")
	}

	path := map[string]any{}

	switch {
	case q.indexedShape != nil:
		path["indexed_shape"] = q.indexedShape
	case q.shape != nil:
		path["shape"] = q.shape
		if q.relation != "" {
			path["relation"] = q.relation
		}
	default:
		return nil, errors.New("eshape: geo_shape query requires either a shape or an indexed_shape")
	}

	params := map[string]any{q.path: path}
	source := map[string]any{"geo_shape": params}

	return json.Marshal(source)
}

// Body renders the query as a JSON reader, ready to be passed as
// esapi.SearchRequest.Body (which accepts io.Reader).
func (q *GeoShapeQuery) Body() (io.Reader, error) {
	data, err := q.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}
