# geo — Notes for AI Agents

- **Coordinate order is `longitude, latitude`** — per the GeoJSON spec, *not* the `lat, lon` order humans usually say aloud. `NewPosition`, `NewPoint`, and the `coordinates` arrays all take longitude first. The `LatLon()` / `LonLat()` helpers exist precisely because both orders are needed depending on the consumer.

- **String forms must preserve precision — never route coordinates through `rosetta/convert.String`.** As of rosetta v0.27.0, `convert.String` rounds floats to two decimal places, which silently corrupts coordinates (~1 km of error). `Position.String()` uses `strconv.FormatFloat(f, 'f', -1, 64)` (minimal lossless digits) for this reason. `formatCoordinatePair` (backing `LonLat`/`LatLon` on `Point` and `Address`) instead uses fixed 10-decimal precision — a deliberately different format, so don't consolidate the two.

- **`Address` map keys must match the struct tags.** `MarshalMap`/`UnmarshalMap` key by the `AddressProperty*` constants, while JSON/BSON key by the struct tags in `address.go`. These two key sets must stay identical or data crossing between the map path and the JSON/BSON path is silently dropped — the historical `"pluscode"` vs `"plusCode"` mismatch did exactly that. When adding an Address field, update the constant, the struct tag, `MarshalMap`, `UnmarshalMap`, and the accessors together.

- **`AddressSchema()` intentionally omits `timezone` and `plusCode`.** They are derived/geocoded fields, not user-editable form inputs, so they're absent from the rosetta schema even though `MarshalMap` emits them. This is deliberate, not an oversight.

- **`SetString("formatted", …)` has a side effect: it clears all geocoded fields** when the value actually changes. The formatted address and the parsed/geocoded fields are kept consistent — a new formatted value invalidates the old geocode.

- **`Point`/`Polygon` marshal to JSON `null` when zero — and only `Point` reads `null` back.** `MarshalJSON` returns `null` for a zero value (works with `omitzero`, not `omitempty`). `UnmarshalJSON` rejects `null` with an error, so a zero value does *not* survive a JSON round-trip; `omitzero` avoids this by dropping the field entirely. BSON differs on both halves: a zero `Point` marshals to `{"type":"Point","coordinates":[0,0]}` rather than `null` (and vanishes entirely under `omitempty`, because `IsZero` satisfies the driver's `Zeroer`), while `Point.UnmarshalBSONValue` *does* accept `null` and `undefined` as the zero `Point`.

- **`Address` follows the rosetta-style `Get*`/`Get*OK` accessor convention.** `GetString`/`GetFloat` return a bare value (zero on miss); `GetStringOK`/`GetFloatOK` add a boolean that reports whether the property name was recognized.

## A type that embeds another type's BSON methods must redeclare them

`Point` embeds `Position`, so `Position`'s methods are promoted into `Point`'s method set. **Go's shadowing rule only protects you when the outer and inner methods share a name**, and BSON offers two different names for the same job:

| Interface | Marshal | Unmarshal |
|---|---|---|
| `bson.Marshaler` / `bson.Unmarshaler` | `MarshalBSON` | `UnmarshalBSON` |
| `bson.ValueMarshaler` / `bson.ValueUnmarshaler` | `MarshalBSONValue` | `UnmarshalBSONValue` |

The mongo-driver registers the **`*Value` hooks first** and `lookupInterfaceDecoder` returns the first match, so `ValueMarshaler` wins whenever a type implements both. In v0.1.0 `Position` moved to the `*Value` names while `Point` kept the plain ones. Nothing shadowed anything any more, the promoted array encoder won, and every stored `Point` silently changed from `{"type":"Point","coordinates":[…]}` to a bare `[lon,lat]` — for two months, until a decode failure took a production page down. See BUG-139.

So: **`Point` declares `MarshalBSONValue`/`UnmarshalBSONValue` at depth 0 specifically to shadow the promoted pair.** It also keeps `MarshalBSON`/`UnmarshalBSON`, delegating to the same `MarshalStruct` and to `UnmarshalBSONValue`, so the two interfaces cannot disagree no matter which one a caller or codec picks. Renaming, removing, or diverging either pair changes the format on disk.

Each type also carries compile-time assertions naming the interfaces it must keep satisfying (`var _ bson.ValueMarshaler = ...`, beside the BSON tests in `point_test.go` and `position_test.go`). They exist because these interfaces are satisfied *structurally*: a signature that drifts out of step with the driver's does not fail to compile at the method itself, it just stops being recognized. `Position` needs them most -- it has no plain `MarshalBSON` to fall back on.

`Polygon` is safe for one reason only: it holds `Position` in a **named field** rather than embedding it, so nothing is promoted. `TestPolygon_DoesNotImplementValueMarshaler` fails if anyone gives it a `*Value` method, because that would change the stored format of every `Polygon`.

`Point.UnmarshalBSONValue` accepts a GeoJSON document *and* a bare coordinate array, permanently. The array shape is what v0.1.0–v0.2.0 wrote, and those records are still out there.

## Format tests assert bytes, not round trips

A struct → bytes → struct round trip is symmetric: it passes just as happily against a format that moved on both sides at once, which is exactly how the v0.1.0 regression survived 100% coverage. The checked-in fixtures in [testdata/golden](testdata/golden/) pin the actual bytes for a fully populated record, one file per historical shape.

Two rules keep them useful. **Marshal a type as a struct field, never only at the top level** — `bson.Marshal(point)` takes the `bson.Marshaler` path directly and never consults the codec registry, so it cannot see this class of bug at all. And **populate every field**: `omitempty` hides a zero value, and a hidden field cannot regress visibly.

After an intended format change, run `go test -update-golden` and read the diff.
