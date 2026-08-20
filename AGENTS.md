# geo — Notes for AI Agents

- **Coordinate order is `longitude, latitude`** — per the GeoJSON spec, *not* the `lat, lon` order humans usually say aloud. `NewPosition`, `NewPoint`, and the `coordinates` arrays all take longitude first. The `LatLon()` / `LonLat()` helpers exist precisely because both orders are needed depending on the consumer.

- **String forms must preserve precision — never route coordinates through `rosetta/convert.String`.** As of rosetta v0.27.0, `convert.String` rounds floats to two decimal places, which silently corrupts coordinates (~1 km of error). `Position.String()` uses `strconv.FormatFloat(f, 'f', -1, 64)` (minimal lossless digits) for this reason. `formatCoordinatePair` (backing `LonLat`/`LatLon` on `Point` and `Address`) instead uses fixed 10-decimal precision — a deliberately different format, so don't consolidate the two.

- **`Address` map keys must match the struct tags.** `MarshalMap`/`UnmarshalMap` key by the `AddressProperty*` constants, while JSON/BSON key by the struct tags in `address.go`. These two key sets must stay identical or data crossing between the map path and the JSON/BSON path is silently dropped — the historical `"pluscode"` vs `"plusCode"` mismatch did exactly that. When adding an Address field, update the constant, the struct tag, `MarshalMap`, `UnmarshalMap`, and the accessors together.

- **`AddressSchema()` intentionally omits `timezone` and `plusCode`.** They are derived/geocoded fields, not user-editable form inputs, so they're absent from the rosetta schema even though `MarshalMap` emits them. This is deliberate, not an oversight.

- **`SetString("formatted", …)` has a side effect: it clears all geocoded fields** when the value actually changes. The formatted address and the parsed/geocoded fields are kept consistent — a new formatted value invalidates the old geocode.

- **`Point`/`Polygon` marshal to `null` when zero — but will not read `null` back.** `MarshalJSON` returns `null` for a zero value (works with `omitzero`, not `omitempty`). `UnmarshalJSON` rejects `null` with an error, so a zero value does *not* survive a JSON round-trip; `omitzero` avoids this by dropping the field entirely. The BSON codec differs again: a zero `Point` marshals to `{"type":"Point","coordinates":[0,0]}`, not `null`.

- **`Address` follows the rosetta-style `Get*`/`Get*OK` accessor convention.** `GetString`/`GetFloat` return a bare value (zero on miss); `GetStringOK`/`GetFloatOK` add a boolean that reports whether the property name was recognized.
