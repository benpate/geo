package geo

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestNewPosition(t *testing.T) {
	require.Equal(t, Position{Longitude: 1, Latitude: 2, Altitude: 0}, NewPosition(1, 2))
	require.Equal(t, Position{Longitude: 1, Latitude: 2, Altitude: 3}, NewPositionWithAltitude(1, 2, 3))
}

func TestPosition_Zeroer(t *testing.T) {
	require.True(t, NewPosition(0, 0).IsZero())
	require.False(t, NewPosition(0, 0).NotZero())

	require.False(t, NewPosition(1, 0).IsZero())
	require.False(t, NewPosition(0, 1).IsZero())
	require.False(t, NewPositionWithAltitude(0, 0, 1).IsZero())
	require.True(t, NewPosition(1, 0).NotZero())
}

func TestPosition_String(t *testing.T) {
	require.Equal(t, "1,2", NewPosition(1, 2).String())
	require.Equal(t, "-118.5,34.25", NewPosition(-118.5, 34.25).String())
}

func TestPosition_MarshalSlice(t *testing.T) {
	// Altitude of zero is omitted from the coordinate pair
	require.Equal(t, []float64{1, 2}, NewPosition(1, 2).MarshalSlice())
	require.Equal(t, []float64{1, 2, 3}, NewPositionWithAltitude(1, 2, 3).MarshalSlice())
}

func TestPosition_UnmarshalSlice(t *testing.T) {

	// check confirms a coordinate slice unmarshals into the expected Position
	check := func(coordinates sliceof.Float, expected Position) {
		position := Position{}
		err := position.UnmarshalSlice(coordinates)
		require.Nil(t, err)
		require.Equal(t, expected, position)
	}

	check(sliceof.Float{1, 2}, Position{Longitude: 1, Latitude: 2})
	check(sliceof.Float{1, 2, 3}, Position{Longitude: 1, Latitude: 2, Altitude: 3})
}

func TestPosition_UnmarshalSlice_Errors(t *testing.T) {

	// Only lengths 2 and 3 are valid; everything else is an error
	checkError := func(coordinates sliceof.Float) {
		position := Position{}
		err := position.UnmarshalSlice(coordinates)
		require.NotNil(t, err)
	}

	checkError(sliceof.Float{})
	checkError(sliceof.Float{1})
	checkError(sliceof.Float{1, 2, 3, 4})
}

func TestPosition_BSON(t *testing.T) {

	// A Position is a BSON array value, so it round-trips as a field within a
	// document (via the bson.ValueMarshaler / ValueUnmarshaler interfaces).
	type wrapper struct {
		Position Position `bson:"position"`
	}

	check := func(position Position) {
		data, err := bson.Marshal(wrapper{Position: position})
		require.Nil(t, err)

		result := wrapper{}
		require.Nil(t, bson.Unmarshal(data, &result))
		require.Equal(t, position, result.Position)
	}

	check(NewPosition(1, 2))
	check(NewPositionWithAltitude(1, 2, 3))
}

func TestPosition_UnmarshalBSONValue_Error(t *testing.T) {

	position := Position{}

	// A BSON value of the wrong type (a string, not an array) should error
	stringType, stringData, err := bson.MarshalValue("not an array")
	require.Nil(t, err)
	require.NotNil(t, position.UnmarshalBSONValue(stringType, stringData))

	// A valid BSON array, but with an invalid coordinate length, should error
	shortType, shortData, err := bson.MarshalValue([]float64{1})
	require.Nil(t, err)
	require.NotNil(t, position.UnmarshalBSONValue(shortType, shortData))
}

func TestPosition_UnmarshalJSON_Error(t *testing.T) {

	// Malformed JSON should surface an error, not panic
	p := Position{}
	require.NotNil(t, p.UnmarshalJSON([]byte("not json")))

	// Valid JSON, but the wrong coordinate length
	require.NotNil(t, p.UnmarshalJSON([]byte("[1]")))
}

func TestPosition_JSON(t *testing.T) {

	p1 := Position{
		Longitude: 1,
		Latitude:  2,
	}

	data, err1 := json.Marshal(p1)
	require.Nil(t, err1)
	require.Equal(t, []byte("[1,2]"), data)

	p2 := Position{}
	err2 := json.Unmarshal(data, &p2)

	require.Nil(t, err2)
	require.Equal(t, p1, p2)
}

func TestPosition_AltitudeJSON(t *testing.T) {

	p1 := Position{
		Longitude: 1,
		Latitude:  2,
		Altitude:  3,
	}

	data, err1 := json.Marshal(p1)
	require.Nil(t, err1)
	require.Equal(t, []byte("[1,2,3]"), data)

	p2 := Position{}
	err2 := json.Unmarshal(data, &p2)

	require.Nil(t, err2)
	require.Equal(t, p1, p2)
}

func TestPosition_UnmarshalJSON_Null(t *testing.T) {

	// JSON null decodes to a nil coordinate slice, which is how a nil actually
	// reaches UnmarshalSlice. A nil slice has length zero, so it is rejected
	// like any other invalid length -- and must not panic on the way.
	position := Position{}
	require.NotNil(t, position.UnmarshalJSON([]byte("null")))

	// The same path through BSON
	require.NotNil(t, position.UnmarshalBSONValue(bson.TypeNull, nil))
}

func TestPosition_UnmarshalSlice_ClearsAltitude(t *testing.T) {

	// A two-element coordinate must CLEAR a previously-set altitude rather
	// than leaving it behind, so a reused Position cannot leak the old value
	position := NewPositionWithAltitude(9, 9, 9)

	require.Nil(t, position.UnmarshalSlice(sliceof.Float{1, 2}))
	require.Equal(t, NewPosition(1, 2), position)
	require.Zero(t, position.Altitude)
}

func TestPosition_String_DropsAltitude(t *testing.T) {

	// String() is a coordinate PAIR, so altitude is deliberately not included
	require.Equal(t, "1,2", NewPositionWithAltitude(1, 2, 3).String())
}

func TestPosition_String_Precision(t *testing.T) {

	// String() must emit the minimal lossless representation. Routing this
	// through rosetta's convert.String would round to two decimals and move
	// the coordinate by roughly a kilometer.
	require.Equal(t, "-118.40853,34.02116", NewPosition(-118.40853, 34.02116).String())
	require.Equal(t, "0.3333333333333333,0.6666666666666666", NewPosition(1.0/3.0, 2.0/3.0).String())
	require.Equal(t, "179.99999999999997,-89.99999999999999", NewPosition(179.99999999999997, -89.99999999999999).String())
}

func TestPosition_String_IsNotFormatCoordinatePair(t *testing.T) {

	// These two formats are deliberately different and must not be merged:
	// String() is minimal-lossless, formatCoordinatePair is fixed to ten
	// decimal places. A change that consolidates them breaks this test.
	position := NewPosition(-118.40853, 34.02116)

	require.Equal(t, "-118.40853,34.02116", position.String())
	require.Equal(t, "-118.4085300000,34.0211600000", formatCoordinatePair(position.Longitude, position.Latitude))
	require.NotEqual(t, position.String(), formatCoordinatePair(position.Longitude, position.Latitude))
}

func TestPosition_ExtremeCoordinates(t *testing.T) {

	// check confirms a coordinate pair survives a JSON round trip exactly
	check := func(longitude float64, latitude float64, altitude float64) {
		original := NewPositionWithAltitude(longitude, latitude, altitude)

		data, err := json.Marshal(original)
		require.Nil(t, err, "marshalling %v", original)

		decoded := Position{}
		require.Nil(t, decoded.UnmarshalJSON(data), "unmarshalling %s", string(data))
		require.Equal(t, original, decoded, "round-tripping %s", string(data))
	}

	check(180, 90, 0)                                                   // the corners of the world
	check(-180, -90, 0)                                                 //
	check(math.MaxFloat64, -math.MaxFloat64, math.MaxFloat64)           // the corners of float64
	check(math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64, 0) // barely not zero
	check(math.Pi, math.E, math.Sqrt2)                                  // irrational, full precision
}

func TestPosition_NegativeZero(t *testing.T) {

	// Negative zero compares equal to zero, so it counts as a zero Position
	position := NewPosition(math.Copysign(0, -1), 0)
	require.True(t, position.IsZero())

	// ...but it is preserved in the string form, and survives a round trip
	require.Equal(t, "-0,0", position.String())

	data, err := json.Marshal(position)
	require.Nil(t, err)
	require.Equal(t, "[-0,0]", string(data))
}

func TestPosition_NaN_And_Infinity(t *testing.T) {

	// JSON cannot represent these, so marshalling must fail loudly rather
	// than emit a coordinate that is silently wrong
	check := func(value float64, label string) {
		_, err := json.Marshal(NewPosition(value, 0))
		require.NotNil(t, err, "marshalling %s should fail", label)
	}

	check(math.NaN(), "NaN")
	check(math.Inf(1), "+Inf")
	check(math.Inf(-1), "-Inf")
}

/******************************************
 * BSON format tests
 ******************************************/

// TestPosition_FieldFormatLock pins the bytes a bare Position writes as a struct
// field. A coordinate pair is a BSON ARRAY, not a document.
func TestPosition_FieldFormatLock(t *testing.T) {

	data, err := bson.Marshal(struct {
		Position Position `bson:"position"`
	}{NewPosition(1, 2)})

	require.Nil(t, err)
	require.Equal(t,
		`{"position": [{"$numberDouble":"1.0"},{"$numberDouble":"2.0"}]}`,
		bson.Raw(data).String())
}
