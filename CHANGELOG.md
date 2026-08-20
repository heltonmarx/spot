# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial release: RFC 7946 GeoJSON geometry types modeled as Go structs —
  `Point`, `MultiPoint`, `LineString`, `MultiLineString`, `Polygon`,
  `MultiPolygon`, and `GeometryCollection` — plus the Elasticsearch geo-shape
  extensions `Envelope` and `Circle`.
- Functional `NewXxx(...)` constructors for each geometry type and a
  `NewShape(...Option)` container builder.
- Type-aware `MarshalJSON`/`UnmarshalJSON` with strict RFC 7946 round-tripping.
- `IsXxx` predicates for the geometry types.
- `eshape` subpackage with the Elasticsearch `geo_shape` query builder:
  `NewGeoShapeQuery`, `IndexedShape`, `Shape`, `Relation`, and the `Relation`
  constants `Intersects`, `Disjoint`, `Within`, and `Contains`.
- Linting via `golangci-lint` (`.golangci.yml`).
