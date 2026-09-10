# geo 🌎

[![Go Reference](https://pkg.go.dev/badge/github.com/benpate/geo.svg)](https://pkg.go.dev/github.com/benpate/geo)
[![Version](https://img.shields.io/github/v/release/benpate/geo?include_prereleases&style=flat-square&color=brightgreen)](https://github.com/benpate/geo/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/benpate/geo/go.yml?branch=main)](https://github.com/benpate/geo/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/geo?style=flat-square)](https://goreportcard.com/report/github.com/benpate/geo)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/geo.svg?style=flat-square)](https://codecov.io/gh/benpate/geo)

## GeoJSON Primitives for Go and MongoDB

This package provides four geographic value types — `Position`, `Point`, `Polygon`, and `Address` — each with JSON and BSON marshalling that follows the [GeoJSON specification (RFC 7946)](https://datatracker.ietf.org/doc/html/rfc7946). `Address` additionally models a human-readable postal address with optional geocoded coordinates.

The types are designed to be useful at their zero value and to round-trip cleanly through JSON, BSON, and `mapof.Any` maps.

## The types

- **`Position`** — a single `longitude, latitude[, altitude]` coordinate. The building block; `Point` embeds it. Marshals as a GeoJSON coordinate *array* (`[lon, lat]`), not an object.
- **`Point`** — a GeoJSON `Point` object (`{"type":"Point","coordinates":[lon,lat]}`).
- **`Polygon`** — a GeoJSON `Polygon`: a single ring of `Position` values.
- **`Address`** — a postal address (`schema.org/PostalAddress`-style) plus optional latitude/longitude, time zone, and Plus Code.

## BSON compatibility

`Point` writes a GeoJSON object and reads back either a GeoJSON object **or** a bare `[lon, lat]` array. The array shape is what v0.1.0 through v0.2.0 wrote, when `Position`'s promoted marshaller shadowed `Point`'s own, and records in that shape are still in the wild. The decoder accepts it permanently; nothing needs migrating to be readable.

`Point` also declares both BSON interfaces on purpose — `MarshalBSONValue`/`UnmarshalBSONValue` to shadow the pair promoted from the embedded `Position`, and `MarshalBSON`/`UnmarshalBSON` for callers that marshal a `Point` on its own. Both routes share one implementation, so they cannot disagree. See [AGENTS.md](AGENTS.md) before changing either.

## References

- GeoJSON: https://geojson.org
- RFC 7946: https://datatracker.ietf.org/doc/html/rfc7946
- MongoDB GeoJSON: https://www.mongodb.com/docs/manual/reference/geojson/

## Pull Requests Welcome

I'm trying to make geo the best it can be, and your help is greatly appreciated. If you find a bug or have an idea for a new feature, please open an issue or submit a pull request. We're all in this together! 🌎
