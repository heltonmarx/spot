package spot

import "encoding/json"

// rfc7946TypeRaw maps Elasticsearch lowercase type names to their pre-encoded
// RFC 7946 PascalCase JSON string values. Envelope and Circle are Elasticsearch
// extensions with no RFC 7946 counterpart and are intentionally absent.
var rfc7946TypeRaw = map[string]json.RawMessage{
	TypePoint:              json.RawMessage(`"Point"`),
	TypeMultiPoint:         json.RawMessage(`"MultiPoint"`),
	TypeLineString:         json.RawMessage(`"LineString"`),
	TypeMultiLineString:    json.RawMessage(`"MultiLineString"`),
	TypePolygon:            json.RawMessage(`"Polygon"`),
	TypeMultiPolygon:       json.RawMessage(`"MultiPolygon"`),
	TypeGeometryCollection: json.RawMessage(`"GeometryCollection"`),
}

// RFC7946Geometry wraps a Geometry and marshals it with PascalCase type names
// as required by RFC 7946. Use it when targeting GeoJSON consumers that follow
// the standard: MongoDB, PostGIS, Solr, and any RFC 7946-strict API.
//
// Envelope and Circle are Elasticsearch extensions with no RFC 7946 equivalent;
// their type names are left unchanged.
type RFC7946Geometry struct {
	g Geometry
}

// RFC7946 wraps g so that MarshalJSON emits RFC 7946-compliant PascalCase type
// names instead of the Elasticsearch-compatible lowercase names.
func RFC7946(g Geometry) RFC7946Geometry {
	return RFC7946Geometry{g: g}
}

// MarshalJSON marshals the wrapped geometry with PascalCase type names.
// GeometryCollection members are recursed through the RFC7946 wrapper via a
// type assertion, avoiding JSON field-name inspection.
func (r RFC7946Geometry) MarshalJSON() ([]byte, error) {
	pascalRaw, hasRFC7946 := rfc7946TypeRaw[r.g.GeoType()]
	if !hasRFC7946 {
		// ES extension (Envelope, Circle): no RFC 7946 name, marshal as-is.
		return json.Marshal(r.g)
	}

	if gc, ok := r.g.(*GeometryCollection); ok {
		wrapped := make([]RFC7946Geometry, len(gc.Geometries))
		for i, g := range gc.Geometries {
			wrapped[i] = RFC7946(g)
		}
		return json.Marshal(struct {
			Type       json.RawMessage   `json:"type"`
			Geometries []RFC7946Geometry `json:"geometries"`
		}{Type: pascalRaw, Geometries: wrapped})
	}

	data, err := json.Marshal(r.g)
	if err != nil {
		return nil, err
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	m["type"] = pascalRaw
	return json.Marshal(m)
}
