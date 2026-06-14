package geo

import (
	"encoding/json"
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

func TestPosition_MarshalBSON(t *testing.T) {

	// Position.MarshalBSON serializes to a bare array, which is NOT a valid
	// standalone BSON document. (Point/Polygon embed it via MarshalStruct
	// instead.) This documents that marshalling a top-level Position fails.
	_, err := bson.Marshal(NewPositionWithAltitude(1, 2, 3))
	require.NotNil(t, err)
}

// NOTE: Position.UnmarshalBSON is intentionally left uncovered. Position
// marshals to a bare BSON array (see TestPosition_MarshalBSON), which is not a
// valid standalone BSON document, so there is no supported path that produces
// bytes this method can consume. Point/Polygon decode through their own
// struct-based UnmarshalBSON methods instead.

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
