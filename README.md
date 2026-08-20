# spot

GeoJSON geometry types for Elasticsearch geo-shape matching, extracted from the
author's [GeoShape Query implementation][pr] for `olivere/elastic` and repackaged
as a standalone, client-agnostic Go module.

This module models all **nine GeoJSON geometry types** as first-class Go structs,
with strict RFC 7946 round-tripping (Marshal/Unmarshal) — `Point`, `MultiPoint`,
`LineString`, `MultiLineString`, `Polygon`, `MultiPolygon`, `GeometryCollection`,
`Envelope`, and `Circle`.

It has **no dependency** on any Elasticsearch client, so it can be reused by the
fluent `olivere/elastic` client, Elastic's official `go-elasticsearch` client, or
any other consumer that needs to serialize a geo shape into JSON.

This single repository also bundles the official-client **GeoShape query adapter**
in the [`eshape`][eshape] subpackage — it depends on `go-elasticsearch`, so it is
only pulled into builds of code that imports it.

[pr]: https://github.com/olivere/elastic/pull/1058
[eshape]: ./eshape

## Install

```sh
go get github.com/heltonmarx/spot
```

## Usage

### Build a shape and serialize the `geo_shape` query DSL

Build a shape with the functional `NewShape(...Option)` API, then marshal it to
the GeoJSON ES expects inside a `geo_shape` query — `spot` is client-agnostic, so
this works with any ES client:

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/heltonmarx/spot"
)

func main() {
	// Point shape (lon, lat order per GeoJSON)
	shape := spot.NewShape(spot.WithPoint([]float64{13.400544, 52.530286}))

	// Or a polygon (closed ring: first and last vertex match)
	// shape := spot.NewShape(spot.WithPolygon([][][]float64{
	//     {{100.0, 0.0}, {101.0, 0.0}, {101.0, 1.0}, {100.0, 1.0}, {100.0, 0.0}},
	// }))

	// And a circle with a radius
	// shape := spot.NewShape(spot.WithCircle("25m", []float64{-109.874838, 44.439550}))

	// Any client can serialize it into the query body:
	body, _ := json.Marshal(map[string]any{
		"query": map[string]any{
			"geo_shape": map[string]any{
				"location": map[string]any{
					"shape":    shape,
					"relation": "within",
				},
			},
		},
	})
	fmt.Println(string(body))
}
```

The example above emits the exact ES `geo_shape` query DSL, e.g.:

```json
{"query":{"geo_shape":{"location":{"shape":{"type":"point","coordinates":[13.400544,52.530286]},"relation":"within"}}}}
```

### Search Elasticsearch through the official client

For a full end-to-end search, the bundled [`eshape`][eshape] subpackage builds the
`geo_shape` query on top of Elastic's official `go-elasticsearch` client, so you
can run a real query without hand-assembling the JSON:

```go
package main

import (
	"context"
	"fmt"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/heltonmarx/spot"
	"github.com/heltonmarx/spot/eshape"
)

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		panic(err)
	}

	// A polygon bounding a region (closed ring, lon/lat order).
	shape := spot.NewShape(spot.WithPolygon([][][]float64{
		{{-74.1, 40.7}, {-73.9, 40.7}, {-73.9, 40.9}, {-74.1, 40.7}},
	}))

	// Build the geo_shape query body against the "location" field.
	body, err := eshape.NewGeoShapeQuery("location").
		Shape(shape).
		Relation(eshape.RelationIntersects).
		Body()
	if err != nil {
		panic(err)
	}

	// Send it through the official client's low-level Search API.
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("places"),
		es.Search.WithBody(body),
	)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	fmt.Println(res.Status())
}
```

> The `eshape` subpackage only pulls in `go-elasticsearch` when you import it —
> using `spot` alone for shape serialization adds no ES dependency to your build.

## Supported geometry types

All constants are exported on the package:

| `Type` constant | `Shape` helper | GeoJSON object |
|---|---|---|
| `TypePoint` | `WithPoint` | point |
| `TypeMultiPoint` | `WithMultiPoint` | multipoint |
| `TypeLineString` | `WithLineString` | linestring |
| `TypeMultiLineString` | `WithMultiLineString` | multilinestring |
| `TypePolygon` | `WithPolygon` | polygon |
| `TypeMultiPolygon` | `WithMultiPolygon` | multipolygon |
| `TypeGeometryCollection` | `WithGeometryCollection` | geometrycollection |
| `TypeEnvelope` | `WithEnvelope` | envelope |
| `TypeCircle` | `WithCircle` | circle |

`Shape` also exposes `IsPoint()`, `IsPolygon()`, … predicates and a
discriminated `MarshalJSON`/`UnmarshalJSON`, so a heterogeneous `Shape` can be
round-tripped without knowing its concrete type.

## Tests

Test fixtures are taken verbatim from **RFC 7946** (the GeoJSON specification),
each exercised for both marshal and unmarshal round-trips.

```sh
go test ./...
```

## Origin & license

The geometry types are derived from the author's [GeoShape Query contribution][pr]
to `olivere/elastic` (PR #1058). Original code Copyright (c) Oliver Eilhard;
packaging Copyright (c) Helton Marques. Both are released under the MIT License.
