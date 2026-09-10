package geo

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestValidateBSONDocument(t *testing.T) {

	valid, err := bson.Marshal(bson.M{"type": "Point"})
	require.Nil(t, err)
	require.Nil(t, validateBSONDocument(valid))
}

func TestValidateBSONDocument_TooShort(t *testing.T) {

	require.NotNil(t, validateBSONDocument(nil))
	require.NotNil(t, validateBSONDocument([]byte{}))
	require.NotNil(t, validateBSONDocument([]byte{5, 0, 0, 0}))
}

// A length prefix of zero is the input that panics the driver's own validator,
// because it indexes data[length-1] without checking the sign.
func TestValidateBSONDocument_ZeroLengthPrefix(t *testing.T) {

	require.NotNil(t, validateBSONDocument([]byte{0, 0, 0, 0, 0}))
}

// A length prefix above MaxInt32 wraps negative, and must be rejected for the
// same reason.
func TestValidateBSONDocument_NegativeLengthPrefix(t *testing.T) {

	require.NotNil(t, validateBSONDocument([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0}))
}

// A sane prefix that overruns the buffer is caught by the driver's validator,
// which is safe to reach once the prefix itself has been checked.
func TestValidateBSONDocument_Malformed(t *testing.T) {

	require.NotNil(t, validateBSONDocument([]byte{99, 0, 0, 0, 1, 2, 3}))
}

// The document that crashed FuzzPolygon_UnmarshalBSON before the guard existed.
func TestValidateBSONDocument_FuzzCrashers(t *testing.T) {

	require.NotNil(t, validateBSONDocument([]byte("\x00\x00\x00\x00")))

	crasher := []byte("1\x00\x00\x00\x02type\x00\b\x00\x00\x000000000\x00\x04coordinates0000\x00000\xa9000000")
	require.NotNil(t, validateBSONDocument(crasher))

	polygon := Polygon{}
	require.NotNil(t, polygon.UnmarshalBSON(crasher))

	point := Point{}
	require.NotNil(t, point.UnmarshalBSON(crasher))
}

/******************************************
 * Whole-record format tests
 *
 * These pin the BYTES this package writes, not the round trip.
 * A round trip is symmetric, so it passes just as happily against
 * a format that moved on both sides at once -- which is exactly
 * how the v0.1.0 Point regression reached production.
 ******************************************/

// updateGolden rewrites the checked-in fixtures instead of asserting against them.
// Run `go test -update-golden` after an INTENDED format change, and read the diff.
var updateGolden = flag.Bool("update-golden", false, "rewrite the testdata/golden fixtures")

// bsonFormatRecord stands in for an application document that persists every geo
// type at once, so these tests pin whole records rather than isolated fields.
type bsonFormatRecord struct {
	Name     string  `bson:"name"`
	Location Point   `bson:"location"`
	Area     Polygon `bson:"area"`
	Where    Address `bson:"where"`
}

// fullBSONFormatRecord returns a record whose every field carries a distinctive
// non-zero value, because `omitempty` hides a zero field from these assertions.
func fullBSONFormatRecord() bsonFormatRecord {

	return bsonFormatRecord{
		Name:     "Bluebird Theater",
		Location: NewPointWithAltitude(-104.9903, 39.7392, 1609),
		Area:     NewPolygon(NewPosition(-105, 39), NewPosition(-104, 39), NewPosition(-104, 40)),
		Where: Address{
			Name:       "Bluebird Theater",
			Formatted:  "3317 E Colfax Ave, Denver, CO 80206",
			Street1:    "3317 E Colfax Ave",
			Street2:    "Suite 2",
			Locality:   "Denver",
			Region:     "CO",
			PostalCode: "80206",
			Country:    "US",
			PlusCode:   "85FPQXX5+9F",
			Timezone:   "America/Denver",
			Latitude:   39.7392,
			Longitude:  -104.9903,
		},
	}
}

// extendedJSON renders raw BSON as indented canonical Extended JSON, so a format
// change lands in a review diff as readable lines instead of a wall of bytes.
func extendedJSON(t *testing.T, data []byte) string {

	t.Helper()

	buffer := bytes.Buffer{}
	require.NoError(t, json.Indent(&buffer, []byte(bson.Raw(data).String()), "", "  "))

	return buffer.String() + "\n"
}

// goldenBSON reads one fixture and converts it back into raw BSON bytes, so the
// tests that consume it exercise the real BSON decoder rather than the ExtJSON one.
func goldenBSON(t *testing.T, name string) []byte {

	t.Helper()

	text, err := os.ReadFile(filepath.Join("testdata", "golden", name))
	require.NoError(t, err, "missing fixture %s", name)

	raw := bson.Raw{}
	require.NoError(t, bson.UnmarshalExtJSON(text, true, &raw), "unreadable fixture %s", name)

	return raw
}

// requireGolden asserts that `actual` matches the named fixture, or rewrites it under -update-golden.
func requireGolden(t *testing.T, name string, actual string) {

	t.Helper()

	path := filepath.Join("testdata", "golden", name)

	if *updateGolden {
		require.NoError(t, os.WriteFile(path, []byte(actual), 0600))
		return
	}

	expected, err := os.ReadFile(path)
	require.NoError(t, err, "missing fixture %s -- run `go test -update-golden`", name)
	require.Equal(t, string(expected), actual, "format changed. If this is intended, run `go test -update-golden` and review the diff")
}

// TestBSON_GoldenRecord pins the exact bytes this package writes for a fully
// populated record.
func TestBSON_GoldenRecord(t *testing.T) {

	data, err := bson.Marshal(fullBSONFormatRecord())
	require.Nil(t, err)

	requireGolden(t, "record.json", extendedJSON(t, data))
}

// TestBSON_HistoricalCorpus decodes every on-disk shape this package has ever
// written. Once a shape is in this corpus it can never be silently dropped.
func TestBSON_HistoricalCorpus(t *testing.T) {

	// Each fixture is a real document shape, named for the version that wrote it
	eras := []string{
		"record.json",        // current: Point is a GeoJSON document
		"record.v0.0.3.json", // pre-2026-06-28: Point is a GeoJSON document
		"record.v0.1.1.json", // 2026-06-28 to v0.3.0: Point is a bare coordinate array
	}

	for _, era := range eras {

		t.Run(era, func(t *testing.T) {

			record := bsonFormatRecord{}
			require.NoError(t, bson.Unmarshal(goldenBSON(t, era), &record))

			// Every era must decode to the SAME values, whatever shape it used
			require.Equal(t, "Bluebird Theater", record.Name)
			require.InDelta(t, -104.9903, record.Location.Longitude, 0.00001)
			require.InDelta(t, 39.7392, record.Location.Latitude, 0.00001)
			require.Equal(t, 3, record.Area.Coordinates.Length())
			require.Equal(t, "Denver", record.Where.Locality)
		})
	}
}

// TestBSON_ByteStability asserts that bytes -> struct -> bytes is identity,
// which the symmetric struct -> bytes -> struct round trip cannot see.
func TestBSON_ByteStability(t *testing.T) {

	original := goldenBSON(t, "record.json")

	record := bsonFormatRecord{}
	require.NoError(t, bson.Unmarshal(original, &record))

	rewritten, err := bson.Marshal(record)
	require.Nil(t, err)

	require.Equal(t, extendedJSON(t, original), extendedJSON(t, rewritten))
}

// TestBSON_MigrationIsNoOpForCurrentShape confirms that re-writing a record already
// in the current shape changes nothing, so a normalization migration can run twice.
func TestBSON_MigrationIsNoOpForCurrentShape(t *testing.T) {

	first := bsonFormatRecord{}
	require.NoError(t, bson.Unmarshal(goldenBSON(t, "record.json"), &first))

	once, err := bson.Marshal(first)
	require.Nil(t, err)

	second := bsonFormatRecord{}
	require.NoError(t, bson.Unmarshal(once, &second))

	twice, err := bson.Marshal(second)
	require.Nil(t, err)

	require.Equal(t, once, twice)
}

// TestBSON_LegacyRecordUpgradesToCurrentShape walks the actual repair path: read
// a bare-array record, write it back, and land on the canonical shape.
func TestBSON_LegacyRecordUpgradesToCurrentShape(t *testing.T) {

	record := bsonFormatRecord{}
	require.NoError(t, bson.Unmarshal(goldenBSON(t, "record.v0.1.1.json"), &record))

	repaired, err := bson.Marshal(record)
	require.Nil(t, err)

	require.Equal(t, extendedJSON(t, goldenBSON(t, "record.json")), extendedJSON(t, repaired))
}
