# geo 🌎

GeoJSON primitives for Go and MongoDB. This package provides four geographic value types — `Position`, `Point`, `Polygon`, and `Address` — each with JSON and BSON marshalling that follows the [GeoJSON specification (RFC 7946)](https://datatracker.ietf.org/doc/html/rfc7946). `Address` additionally models a human-readable postal address with optional geocoded coordinates.

The types are designed to be useful at their zero value and to round-trip cleanly through JSON, BSON, and `mapof.Any` maps.

## The types

- **`Position`** — a single `longitude, latitude[, altitude]` coordinate. The building block; `Point` embeds it. Marshals as a GeoJSON coordinate *array* (`[lon, lat]`), not an object.
- **`Point`** — a GeoJSON `Point` object (`{"type":"Point","coordinates":[lon,lat]}`).
- **`Polygon`** — a GeoJSON `Polygon`: a single ring of `Position` values.
- **`Address`** — a postal address (`schema.org/PostalAddress`-style) plus optional latitude/longitude, time zone, and Plus Code.

## What matters here

- **Coordinate order is `longitude, latitude`** — per the GeoJSON spec, *not* the `lat, lon` order humans usually say aloud. `NewPosition`, `NewPoint`, and the `coordinates` arrays all take longitude first. The `LatLon()` / `LonLat()` helpers exist precisely because both orders are needed depending on the consumer.

- **String forms must preserve precision — never route coordinates through `rosetta/convert.String`.** As of rosetta v0.27.0, `convert.String` rounds floats to two decimal places, which silently corrupts coordinates (~1 km of error). `Position.String()` uses `strconv.FormatFloat(f, 'f', -1, 64)` (minimal lossless digits) for this reason. `formatCoordinatePair` (backing `LonLat`/`LatLon` on `Point` and `Address`) instead uses fixed 10-decimal precision — a deliberately different format, so don't consolidate the two.

- **`Address` map keys must match the struct tags.** `MarshalMap`/`UnmarshalMap` key by the `AddressProperty*` constants, while JSON/BSON key by the struct tags in `address.go`. These two key sets must stay identical or data crossing between the map path and the JSON/BSON path is silently dropped — the historical `"pluscode"` vs `"plusCode"` mismatch did exactly that. When adding an Address field, update the constant, the struct tag, `MarshalMap`, `UnmarshalMap`, and the accessors together.

- **`AddressSchema()` intentionally omits `timezone` and `plusCode`.** They are derived/geocoded fields, not user-editable form inputs, so they're absent from the rosetta schema even though `MarshalMap` emits them. This is deliberate, not an oversight.

- **`SetString("formatted", …)` has a side effect: it clears all geocoded fields** when the value actually changes. The formatted address and the parsed/geocoded fields are kept consistent — a new formatted value invalidates the old geocode.

- **`Point`/`Polygon` marshal to `null` when zero.** `MarshalJSON` returns `null` for a zero value (works with `omitzero`, not `omitempty`). Round-tripping a zero `Point` through JSON yields a zero `Point`, not an error.

## Accessor pattern

`Address` follows the rosetta-style `Get*`/`Get*OK` convention: `GetString`/`GetFloat` return a bare value (zero on miss); `GetStringOK`/`GetFloatOK` add a boolean that reports whether the property name was recognized.

## References

- GeoJSON: https://geojson.org
- RFC 7946: https://datatracker.ietf.org/doc/html/rfc7946
- MongoDB GeoJSON: https://www.mongodb.com/docs/manual/reference/geojson/
