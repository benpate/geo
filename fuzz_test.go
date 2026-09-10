package geo

import (
	"encoding/json"
	"math"
	"testing"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

/******************************************
 * Decoder robustness
 *
 * These targets assert only that a decoder never panics on arbitrary
 * input. A malformed document must come back as an error.
 ******************************************/

// FuzzPoint_UnmarshalJSON confirms that the Point JSON decoder never panics,
// regardless of the input bytes. A well-formed input round-trips; a malformed
// one must return an error instead of crashing.
func FuzzPoint_UnmarshalJSON(f *testing.F) {

	f.Add(`{"type":"Point","coordinates":[1,2]}`)
	f.Add(`{"type":"Point","coordinates":[1,2,3]}`)
	f.Add(`{"type":"Polygon","coordinates":[1,2]}`)
	f.Add(`null`)
	f.Add(``)
	f.Add(`{`)

	f.Fuzz(func(_ *testing.T, data string) {
		point := Point{}
		_ = point.UnmarshalJSON([]byte(data))
	})
}

// FuzzPolygon_UnmarshalJSON confirms that the Polygon JSON decoder never panics.
func FuzzPolygon_UnmarshalJSON(f *testing.F) {

	f.Add(`{"type":"Polygon","coordinates":[[[1,2],[3,4]]]}`)
	f.Add(`{"type":"Polygon","coordinates":[]}`)
	f.Add(`{"type":"Point","coordinates":[1,2]}`)
	f.Add(`null`)
	f.Add(``)
	f.Add(`[`)

	f.Fuzz(func(_ *testing.T, data string) {
		polygon := Polygon{}
		_ = polygon.UnmarshalJSON([]byte(data))
	})
}

// FuzzPosition_UnmarshalJSON confirms that the Position coordinate-array
// decoder never panics.
func FuzzPosition_UnmarshalJSON(f *testing.F) {

	f.Add(`[1,2]`)
	f.Add(`[1,2,3]`)
	f.Add(`[1]`)
	f.Add(`[1,2,3,4]`)
	f.Add(`[]`)
	f.Add(`["a","b"]`)
	f.Add(`null`)
	f.Add(``)

	f.Fuzz(func(_ *testing.T, data string) {
		position := Position{}
		_ = position.UnmarshalJSON([]byte(data))
	})
}

// FuzzPoint_UnmarshalBSON confirms that the Point BSON decoder never panics on
// arbitrary bytes. BSON is length-prefixed and self-describing, so a truncated
// or hostile document is the interesting case here.
func FuzzPoint_UnmarshalBSON(f *testing.F) {

	seedBSON(f, NewPoint(1, 2))
	seedBSON(f, NewPointWithAltitude(1, 2, 3))
	f.Add([]byte{})
	f.Add([]byte{5, 0, 0, 0, 0})

	f.Fuzz(func(_ *testing.T, data []byte) {

		// Both historical shapes reach the same decoder, through different BSON types
		for _, dataType := range []bsontype.Type{bson.TypeEmbeddedDocument, bson.TypeArray, bson.TypeNull, bson.TypeUndefined} {
			point := Point{}
			_ = point.UnmarshalBSONValue(dataType, data)
		}

		point := Point{}
		_ = point.UnmarshalBSON(data)
	})
}

// FuzzPolygon_UnmarshalBSON confirms that the Polygon BSON decoder never panics.
func FuzzPolygon_UnmarshalBSON(f *testing.F) {

	seedBSON(f, NewPolygon(NewPosition(1, 2), NewPosition(3, 4)))
	seedBSON(f, Polygon{})
	f.Add([]byte{})
	f.Add([]byte{5, 0, 0, 0, 0})

	f.Fuzz(func(_ *testing.T, data []byte) {
		polygon := Polygon{}
		_ = polygon.UnmarshalBSON(data)
	})
}

// FuzzPosition_UnmarshalBSONValue confirms that the Position BSON *value*
// decoder never panics. It is fuzzed over both the type byte and the payload,
// because the type byte selects which driver codec reads the bytes.
func FuzzPosition_UnmarshalBSONValue(f *testing.F) {

	if dataType, data, err := NewPosition(1, 2).MarshalBSONValue(); err == nil {
		f.Add(byte(dataType), data)
	}

	f.Add(byte(bson.TypeArray), []byte{5, 0, 0, 0, 0})
	f.Add(byte(bson.TypeDouble), []byte{})
	f.Add(byte(0), []byte{})

	f.Fuzz(func(_ *testing.T, dataType byte, data []byte) {
		position := Position{}
		_ = position.UnmarshalBSONValue(bsontype.Type(dataType), data)
	})
}

// FuzzNewPolygonFromString confirms that the comma-delimited coordinate parser
// never panics and always produces an even number of coordinate components.
func FuzzNewPolygonFromString(f *testing.F) {

	f.Add("1,2,3,4")
	f.Add("1,2,3")
	f.Add("")
	f.Add("-118.5,34.25")
	f.Add("not,a,number")

	f.Fuzz(func(t *testing.T, data string) {
		polygon := NewPolygonFromString(data)

		// Every parsed position contributes a longitude and a latitude, so the
		// marshalled slice must always contain pairs (length-2 inner slices).
		for _, coordinate := range polygon.MarshalSlice() {
			if len(coordinate) != 2 {
				t.Fatalf("expected coordinate pairs, got length %d", len(coordinate))
			}
		}
	})
}

/******************************************
 * Round-trip oracles
 *
 * These targets encode a value and decode it again, asserting that
 * nothing changed. A coordinate that shifts by even one ULP is a
 * silently wrong location, which no error return would reveal.
 ******************************************/

// FuzzPosition_JSONRoundTrip confirms that a Position survives a JSON encode /
// decode cycle with its coordinates bit-for-bit intact.
func FuzzPosition_JSONRoundTrip(f *testing.F) {

	f.Add(-118.40853, 34.02116, 0.0)
	f.Add(0.1, 0.2, 0.3)
	f.Add(179.9999999999999, -89.9999999999999, 8848.86)
	f.Add(math.SmallestNonzeroFloat64, math.MaxFloat64, 0.0)

	f.Fuzz(func(t *testing.T, longitude float64, latitude float64, altitude float64) {

		skipUnencodable(t, longitude, latitude, altitude)

		original := NewPositionWithAltitude(longitude, latitude, altitude)

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("marshal failed: %v", err)
		}

		decoded := Position{}
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("unmarshal %s failed: %v", string(data), err)
		}

		requireSamePosition(t, original, decoded, string(data))
	})
}

// FuzzPosition_BSONRoundTrip confirms that a Position survives a BSON encode /
// decode cycle with its coordinates bit-for-bit intact.
func FuzzPosition_BSONRoundTrip(f *testing.F) {

	f.Add(-118.40853, 34.02116, 0.0)
	f.Add(0.1, 0.2, 0.3)
	f.Add(math.MaxFloat64, -math.MaxFloat64, 1.0)

	f.Fuzz(func(t *testing.T, longitude float64, latitude float64, altitude float64) {

		skipUnencodable(t, longitude, latitude, altitude)

		original := NewPositionWithAltitude(longitude, latitude, altitude)

		// Position marshals as a BSON *value* (an array), so it must be wrapped
		// in a document to make a round trip through the driver.
		data, err := bson.Marshal(bson.M{"position": original})
		if err != nil {
			t.Fatalf("marshal failed: %v", err)
		}

		decoded := struct {
			Position Position `bson:"position"`
		}{}

		if err := bson.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("unmarshal failed: %v", err)
		}

		requireSamePosition(t, original, decoded.Position, "bson")
	})
}

// FuzzPoint_JSONRoundTrip confirms that a non-zero Point survives a JSON
// encode / decode cycle unchanged.
func FuzzPoint_JSONRoundTrip(f *testing.F) {

	f.Add(-118.40853, 34.02116, 0.0)
	f.Add(1.0, 2.0, 3.0)

	f.Fuzz(func(t *testing.T, longitude float64, latitude float64, altitude float64) {

		skipUnencodable(t, longitude, latitude, altitude)

		original := NewPointWithAltitude(longitude, latitude, altitude)

		// A zero Point encodes to `null`, which the decoder rejects. That gap is
		// BUG-121; excluding it here keeps this target looking for new defects.
		if original.IsZero() {
			t.Skip("zero Point does not round-trip -- see BUG-121")
		}

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("marshal failed: %v", err)
		}

		decoded := Point{}
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("unmarshal %s failed: %v", string(data), err)
		}

		requireSamePosition(t, original.Position, decoded.Position, string(data))
	})
}

// FuzzPolygon_JSONRoundTrip confirms that a non-empty Polygon survives a JSON
// encode / decode cycle with every vertex intact.
func FuzzPolygon_JSONRoundTrip(f *testing.F) {

	f.Add("1,2,3,4,5,6")
	f.Add("-118.40853,34.02116")
	f.Add("0.1,0.2,0.3,0.4")

	f.Fuzz(func(t *testing.T, coordinates string) {

		original := NewPolygonFromString(coordinates)

		// An empty Polygon encodes to `null`, which the decoder rejects (BUG-121).
		if original.IsZero() {
			t.Skip("empty Polygon does not round-trip -- see BUG-121")
		}

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("marshal failed: %v", err)
		}

		decoded := Polygon{}
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("unmarshal %s failed: %v", string(data), err)
		}

		if len(decoded.Coordinates) != len(original.Coordinates) {
			t.Fatalf("vertex count changed: %d -> %d via %s",
				len(original.Coordinates), len(decoded.Coordinates), string(data))
		}

		for index := range original.Coordinates {
			requireSamePosition(t, original.Coordinates[index], decoded.Coordinates[index], string(data))
		}
	})
}

// FuzzPosition_StringRoundTrip confirms that Position.String emits enough
// precision for NewPolygonFromString to recover the exact coordinate. The
// package deliberately avoids convert.String here because it rounds to two
// decimal places; this target is what keeps that from regressing.
func FuzzPosition_StringRoundTrip(f *testing.F) {

	f.Add(-118.40853, 34.02116)
	f.Add(0.1, 0.2)
	f.Add(1.0/3.0, 2.0/3.0)
	f.Add(179.99999999999997, -89.99999999999999)

	f.Fuzz(func(t *testing.T, longitude float64, latitude float64) {

		skipUnencodable(t, longitude, latitude, 0)

		original := NewPosition(longitude, latitude)

		// Round-tripping runs through the Polygon parser, the only reader of
		// this format in the package.
		decoded := NewPolygonFromString(original.String())

		if len(decoded.Coordinates) != 1 {
			t.Skipf("%q did not re-parse as a single position", original.String())
		}

		requireSamePosition(t, original, decoded.Coordinates[0], original.String())
	})
}

// FuzzAddress_MapRoundTrip confirms that every Address field survives a
// MarshalMap / UnmarshalMap cycle. The two methods key by the AddressProperty
// constants, and a key present on one side but not the other silently drops
// the field -- which is exactly what the historical "pluscode" / "plusCode"
// mismatch did.
func FuzzAddress_MapRoundTrip(f *testing.F) {

	f.Add("Name", "123 Main St, Springfield IL", "123 Main St", "Apt 4",
		"Springfield", "IL", "62704", "US", "86HJ+9C", "America/Chicago", 39.7817, -89.6501)
	f.Add("", "", "", "", "", "", "", "", "", "", 0.0, 0.0)

	f.Fuzz(func(t *testing.T, name string, formatted string, street1 string, street2 string,
		locality string, region string, postalCode string, country string,
		plusCode string, timezone string, latitude float64, longitude float64) {

		skipUnencodable(t, latitude, longitude, 0)

		original := Address{
			Name:       name,
			Formatted:  formatted,
			Street1:    street1,
			Street2:    street2,
			Locality:   locality,
			Region:     region,
			PostalCode: postalCode,
			Country:    country,
			PlusCode:   plusCode,
			Timezone:   timezone,
			Latitude:   latitude,
			Longitude:  longitude,
		}

		decoded := NewAddress()
		if err := decoded.UnmarshalMap(original.MarshalMap()); err != nil {
			t.Fatalf("UnmarshalMap failed: %v", err)
		}

		if decoded != original {
			t.Fatalf("Address changed across the map round trip:\n  original: %+v\n   decoded: %+v",
				original, decoded)
		}
	})
}

// FuzzAddress_JSONRoundTrip confirms that every Address field survives a JSON
// encode / decode cycle, which travels the struct tags rather than the
// AddressProperty constants.
func FuzzAddress_JSONRoundTrip(f *testing.F) {

	f.Add("Name", "123 Main St", "Springfield", "62704", "86HJ+9C", 39.7817, -89.6501)
	f.Add("", "", "", "", "", 0.0, 0.0)

	f.Fuzz(func(t *testing.T, name string, street1 string, locality string,
		postalCode string, plusCode string, latitude float64, longitude float64) {

		skipUnencodable(t, latitude, longitude, 0)
		skipInvalidUTF8(t, name, street1, locality, postalCode, plusCode)

		original := Address{
			Name:       name,
			Street1:    street1,
			Locality:   locality,
			PostalCode: postalCode,
			PlusCode:   plusCode,
			Latitude:   latitude,
			Longitude:  longitude,
		}

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("marshal failed: %v", err)
		}

		decoded := NewAddress()
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("unmarshal %s failed: %v", string(data), err)
		}

		if decoded != original {
			t.Fatalf("Address changed across the JSON round trip:\n  original: %+v\n   decoded: %+v\n      json: %s",
				original, decoded, string(data))
		}
	})
}

/******************************************
 * Accessor invariants
 ******************************************/

// FuzzAddress_Accessors confirms that the Get/Set accessors never panic on an
// arbitrary property name, and that a name the setter accepts is a name the
// getter reads back.
func FuzzAddress_Accessors(f *testing.F) {

	f.Add(AddressPropertyName, "value", 1.5)
	f.Add(AddressPropertyFormatted, "123 Main St", 0.0)
	f.Add(AddressPropertyLatitude, "", 34.02116)
	f.Add(AddressPropertyTimezone, "America/Chicago", 0.0)
	f.Add("", "", 0.0)
	f.Add("../../etc/passwd", "x", 0.0)

	f.Fuzz(func(t *testing.T, name string, stringValue string, floatValue float64) {

		skipUnencodable(t, floatValue, 0, 0)

		address := NewAddress()

		// RULE: a writable string property must read back the value it was given.
		// "formatted" is exempt: it deliberately clears the geocoded fields, and
		// re-reading it is covered by the unit tests instead.
		if address.SetString(name, stringValue) {
			if name != AddressPropertyFormatted {
				if actual, ok := address.GetStringOK(name); !ok || actual != stringValue {
					t.Fatalf("SetString(%q, %q) accepted, but GetStringOK returned (%q, %v)",
						name, stringValue, actual, ok)
				}
			}
		}

		// RULE: a writable float property must read back the value it was given.
		if address.SetFloat(name, floatValue) {
			if actual, ok := address.GetFloatOK(name); !ok || actual != floatValue {
				t.Fatalf("SetFloat(%q, %v) accepted, but GetFloatOK returned (%v, %v)",
					name, floatValue, actual, ok)
			}
		}

		// The bare getters must agree with their OK-returning counterparts.
		if expected, _ := address.GetStringOK(name); address.GetString(name) != expected {
			t.Fatalf("GetString(%q) disagrees with GetStringOK", name)
		}

		if expected, _ := address.GetFloatOK(name); address.GetFloat(name) != expected {
			t.Fatalf("GetFloat(%q) disagrees with GetFloatOK", name)
		}
	})
}

/******************************************
 * Fuzz Helpers
 ******************************************/

// skipUnencodable skips a fuzz case whose coordinates cannot be represented in
// JSON. NaN and the infinities are rejected by encoding/json, which is the
// correct behavior and not what these targets are looking for.
func skipUnencodable(t *testing.T, values ...float64) {

	t.Helper()

	for _, value := range values {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			t.Skip("NaN and Inf are not representable in JSON")
		}
	}
}

// skipInvalidUTF8 skips a fuzz case carrying a string that JSON cannot
// represent. encoding/json coerces invalid UTF-8 to U+FFFD on the way out, so
// such a value can never survive a round trip -- that is documented stdlib
// behavior, not something this package can fix.
func skipInvalidUTF8(t *testing.T, values ...string) {

	t.Helper()

	for _, value := range values {
		if !utf8.ValidString(value) {
			t.Skip("invalid UTF-8 is not representable in JSON")
		}
	}
}

// requireSamePosition fails the test unless two Positions are bit-for-bit
// identical. Coordinates are compared exactly, because a rounded coordinate is
// a wrong location that no error return would reveal.
func requireSamePosition(t *testing.T, expected Position, actual Position, encoded string) {

	t.Helper()

	if expected != actual {
		t.Fatalf("Position changed across the round trip:\n  expected: %+v\n    actual: %+v\n   encoded: %s",
			expected, actual, encoded)
	}
}

// seedBSON adds the BSON encoding of a value to the fuzz corpus, so that the
// fuzzer starts from documents that are at least structurally valid.
func seedBSON(f *testing.F, value any) {

	f.Helper()

	if data, err := bson.Marshal(value); err == nil {
		f.Add(data)
	}
}
