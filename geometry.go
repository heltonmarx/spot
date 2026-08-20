// Package spot models GeoJSON geometry types (Point, LineString, Polygon, and
// their multi/collection variants) plus the additional shape types that the
// Elasticsearch geo-shape spatial strategy accepts: Envelope and Circle.
// Shapes unmarshal GeoJSON with type-aware dispatch and marshal back to the
// canonical form.
package spot

import (
	"encoding/json"
	"fmt"
)

// The geometry types supported by Elasticsearch.
//
// For more details, see:
// https://www.elastic.co/guide/en/elasticsearch/reference/current/geo-shape.html#spatial-strategy
const (
	TypePoint              = "point"
	TypeMultiPoint         = "multipoint"
	TypeLineString         = "linestring"
	TypeMultiLineString    = "multilinestring"
	TypePolygon            = "polygon"
	TypeMultiPolygon       = "multipolygon"
	TypeGeometryCollection = "geometrycollection"
	TypeEnvelope           = "envelope"
	TypeCircle             = "circle"
)

// rawGeometry holds generic data used to unmarshal GeoJSON information.
type rawGeometry struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
	Geometries  []rawGeometry   `json:"geometries"`
	Radius      string          `json:"radius,omitempty"`
}

// rawGeometryDecoder is implemented by every concrete geometry type; each
// decodes a rawGeometry into itself.
type rawGeometryDecoder interface {
	decode(*rawGeometry) error
}

// decoded decodes raw into g and returns g. It factors out the identical
// construct-then-decode pattern shared by every case of decodeGeometry.
func decoded(raw *rawGeometry, g rawGeometryDecoder) (any, error) {
	if err := g.decode(raw); err != nil {
		return nil, err
	}
	return g, nil
}

// decodeGeometry decodes a rawGeometry into the concrete geometry type
// matching its "type" field. It is the single source of truth for the
// type-to-struct mapping, shared by Shape and GeometryCollection.
func decodeGeometry(raw *rawGeometry) (any, error) {
	switch raw.Type {
	case TypePoint:
		return decoded(raw, &Point{})
	case TypeMultiPoint:
		return decoded(raw, &MultiPoint{})
	case TypeLineString:
		return decoded(raw, &LineString{})
	case TypeMultiLineString:
		return decoded(raw, &MultiLineString{})
	case TypePolygon:
		return decoded(raw, &Polygon{})
	case TypeMultiPolygon:
		return decoded(raw, &MultiPolygon{})
	case TypeGeometryCollection:
		return decoded(raw, &GeometryCollection{})
	case TypeEnvelope:
		return decoded(raw, &Envelope{})
	case TypeCircle:
		return decoded(raw, &Circle{})
	default:
		return nil, fmt.Errorf("geo: unknown type `%s`", raw.Type)
	}
}
