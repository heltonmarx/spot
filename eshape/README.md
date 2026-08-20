# eshape — GeoShape Query adapter for the official Go client

The `eshape` subpackage ports the [GeoShape Query][pr] (originally for the fluent
`olivere/elastic` client) onto Elastic's **official Go client**
`github.com/elastic/go-elasticsearch`. It builds the ES `geo_shape` query DSL as
plain JSON, using this module's client-agnostic `spot` types, so the query composes
with the official client's raw / low-level API.

[pr]: https://github.com/olivere/elastic/pull/1058

## Why adapter + core split

`olivere/elastic` is a high-level builder-style client; `go-elasticsearch` is
low-level with a generated typed DSL. Rather than fight the generator, this package
keeps the reusable, hard part — the GeoJSON modeling and RFC 7946 compliance — in
the `spot` package, and glues it into the official client via its ability to accept
arbitrary JSON bodies. That's the pragmatic way geo-shape queries are sent through
the low-level client in production.

## Usage

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

	// A polygon bounding box around a region (closed ring, lon/lat order).
	shape := spot.NewShape(spot.WithPolygon([][][]float64{
		{{-74.1, 40.7}, {-73.9, 40.7}, {-73.9, 40.9}, {-74.1, 40.7}},
	}))

	// Build the geo_shape query body.
	body, err := eshape.NewGeoShapeQuery("location").
		Shape(shape).
		Relation(eshape.RelationWithin).
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

## Supported relations

The geo-shape query accepts the `intersects`, `disjoint`, `within`, and `contains`
spatial relations (see the `eshape.Relation*` constants).

## License

Derived from the author's GeoShape contribution to `olivere/elastic` (PR #1058).
Original code Copyright (c) Oliver Eilhard; adapter/packaging Copyright (c)
Helton Marques. MIT licensed.
