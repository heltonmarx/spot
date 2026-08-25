# spot

[![Go Reference](https://pkg.go.dev/badge/github.com/heltonmarx/spot.svg)](https://pkg.go.dev/github.com/heltonmarx/spot)

RFC7946 GeoJSON geometry types for Go, with Elasticsearch geo-shape support.

This module models all nine geometry types as first-class Go structs with
typed marshal/unmarshal: `Point`, `MultiPoint`, `LineString`, `MultiLineString`,
`Polygon`, `MultiPolygon`, `GeometryCollection`, `Envelope`, and `Circle`.

All types implement the `Geometry` interface (`GeoType() string`), so a
`GeometryCollection` carries `[]Geometry` rather than `[]any`.

> **Type name casing:** this library uses lowercase type names (e.g. `"point"`,
> `"polygon"`) as required by Elasticsearch. This differs from RFC 7946, which
> specifies PascalCase (`"Point"`, `"Polygon"`).

This module has **no runtime dependency** on any Elasticsearch client.

## Install

```sh
go get github.com/heltonmarx/spot
```

## Usage

### 1. Indexing a document with a geo_shape field

When you index a document, Elasticsearch expects a GeoJSON geometry object for
any `geo_shape` mapped field. `spot` gives you typed constructors instead of
hand-assembled `map[string]any`:

```go
doc := map[string]any{
    "name": "Central Park",
    "location": spot.NewPolygon([][][]float64{
        {
            {-73.98, 40.77}, {-73.95, 40.77},
            {-73.95, 40.76}, {-73.98, 40.76},
            {-73.98, 40.77},
        },
    }),
}

body, _ := json.Marshal(doc)
es.Index("places", strings.NewReader(string(body)))
```

### 2. Geo-shape queries

Find documents whose shape intersects, contains, or is within a query shape.
Marshal a `spot.Shape` into `GeoShapeFieldQuery.Shape` from the official client:

```go
import (
    "context"
    "encoding/json"

    elasticsearch "github.com/elastic/go-elasticsearch/v8"
    estypes "github.com/elastic/go-elasticsearch/v8/typedapi/types"
    "github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/geoshaperelation"
    "github.com/heltonmarx/spot"
)

func search(es *elasticsearch.TypedClient) {
    shapeBytes, _ := json.Marshal(spot.NewShape(spot.WithEnvelope([][]float64{
        {-74.1, 40.9}, {-73.9, 40.7},
    })))

    q := estypes.NewGeoShapeQuery()
    q.GeoShapeQuery["location"] = estypes.GeoShapeFieldQuery{
        Shape:    shapeBytes,
        Relation: &geoshaperelation.Intersects,
    }

    es.Search().Index("places").
        Request(&estypes.SearchRequest{
            Query: &estypes.Query{GeoShape: q},
        }).
        Do(context.Background())
}
```

Spatial relations map to real questions:

| Relation | Meaning |
|---|---|
| `intersects` | shapes overlap at all (default) |
| `within` | document shape is fully inside the query shape |
| `contains` | document shape fully encloses the query shape |
| `disjoint` | shapes have no overlap |

### 3. Parsing shapes from Elasticsearch responses

When reading a document back from ES, the `location` field arrives as raw
GeoJSON. `Shape.UnmarshalJSON` dispatches to the correct concrete type
automatically (no manual type-switching required):

```go
var shape spot.Shape
json.Unmarshal(hit["location"], &shape)

switch {
case shape.IsPolygon():
    rings := shape.Polygon.Coordinates
    // ...
case shape.IsPoint():
    lon, lat := shape.Point.Coordinates[0], shape.Point.Coordinates[1]
    // ...
}
```

### 4. GeometryCollection for mixed-type fields

Elasticsearch supports indexing a `geometrycollection`, for example a venue
that has both a polygon boundary and a point entrance. `GeometryCollection`
carries `[]Geometry`, so you can build and inspect collections without type
assertions:

```go
collection := spot.NewGeometryCollection([]spot.Geometry{
    spot.NewPoint([]float64{-73.98, 40.76}),
    spot.NewPolygon([][][]float64{
        {
            {-73.99, 40.77}, {-73.97, 40.77},
            {-73.97, 40.75}, {-73.99, 40.75},
            {-73.99, 40.77},
        },
    }),
})
```

### 5. Envelope and Circle for bounding-box and radius queries

`Envelope` (top-left + bottom-right corners) and `Circle` (center + radius)
are Elasticsearch extensions not in standard GeoJSON. They are the most
efficient shapes for bounding-box and proximity queries:

```go
// Bounding box: upper-left corner first, lower-right second.
envelope := spot.NewShape(spot.WithEnvelope([][]float64{
    {-74.1, 40.9},
    {-73.9, 40.7},
}))

// Circle: radius defaults to meters when no unit suffix is given.
circle := spot.NewShape(spot.WithCircle("5km", []float64{-73.98, 40.76}))
```

## Supported geometry types

| `Type` constant | `Shape` helper | GeoJSON object |
|---|---|---|
| `TypePoint` | `WithPoint` | point |
| `TypeMultiPoint` | `WithMultiPoint` | multipoint |
| `TypeLineString` | `WithLineString` | linestring |
| `TypeMultiLineString` | `WithMultiLineString` | multilinestring |
| `TypePolygon` | `WithPolygon` | polygon |
| `TypeMultiPolygon` | `WithMultiPolygon` | multipolygon |
| `TypeGeometryCollection` | `WithGeometryCollection` | geometrycollection |
| `TypeEnvelope` | `WithEnvelope` | envelope (ES extension) |
| `TypeCircle` | `WithCircle` | circle (ES extension) |

`Shape` also exposes `IsPoint()`, `IsPolygon()`, ... predicates and a
discriminated `MarshalJSON`/`UnmarshalJSON`, so a heterogeneous `Shape` can be
round-tripped without knowing its concrete type.

## Compatibility

`spot` uses lowercase type names (`"point"`, `"polygon"`) as required by
Elasticsearch. Other GeoJSON consumers follow RFC 7946, which uses PascalCase
(`"Point"`, `"Polygon"`). Use the `RFC7946` wrapper for those systems.

```go
// Elasticsearch / OpenSearch (default, lowercase)
shapeBytes, _ := json.Marshal(shape)

// MongoDB, PostGIS, Solr, and any RFC 7946 consumer (PascalCase)
shapeBytes, _ := json.Marshal(spot.RFC7946(shape))
```

`RFC7946` recurses into `GeometryCollection` members automatically. `Envelope`
and `Circle` are Elasticsearch extensions with no RFC 7946 equivalent; their
type names are left unchanged by the wrapper.

| Consumer | Type name format | How to marshal |
|---|---|---|
| Elasticsearch | lowercase | `json.Marshal(shape)` |
| OpenSearch | lowercase | `json.Marshal(shape)` |
| CrateDB | lowercase | `json.Marshal(shape)` |
| MongoDB | PascalCase (RFC 7946) | `json.Marshal(spot.RFC7946(shape))` |
| PostgreSQL + PostGIS | PascalCase (RFC 7946) | `json.Marshal(spot.RFC7946(shape))` |
| Solr | PascalCase (RFC 7946) | `json.Marshal(spot.RFC7946(shape))` |

## Tests

```sh
go test ./...
```

## License

Copyright (c) Helton Marques. Released under the MIT License.
