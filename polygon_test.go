package geo

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestPolygon_Zeroer(t *testing.T) {
	require.True(t, NewPolygon().IsZero())
	require.False(t, NewPolygon().NotZero())

	require.False(t, NewPolygon(NewPosition(0, 0)).IsZero())
	require.True(t, NewPolygon(NewPosition(0, 0)).NotZero())
}

func TestNewPolygonFromString(t *testing.T) {

	// check confirms a coordinate string parses into the expected positions
	check := func(data string, expected ...Position) {
		require.Equal(t, NewPolygon(expected...), NewPolygonFromString(data))
	}

	check("1,2,3,4",
		NewPosition(1, 2),
		NewPosition(3, 4),
	)

	check("-118.5,34.25,-119,35",
		NewPosition(-118.5, 34.25),
		NewPosition(-119, 35),
	)

	// An odd trailing coordinate is dropped (no partner to pair with)
	check("1,2,3",
		NewPosition(1, 2),
	)
}

func TestNewPolygonFromString_Empty(t *testing.T) {
	require.True(t, NewPolygonFromString("").IsZero())
}

func TestPolygon_String(t *testing.T) {

	require.Equal(t, "", NewPolygon().String())

	polygon := NewPolygon(
		NewPosition(1, 2),
		NewPosition(3, 4),
	)
	require.Equal(t, "1,2,3,4", polygon.String())
}

func TestPolygon_GeoJSON(t *testing.T) {

	polygon := NewPolygon(
		NewPosition(1, 2),
		NewPosition(3, 4),
	)

	result := polygon.GeoJSON()
	require.Equal(t, PropertyTypePolygon, result[PropertyType])
	require.Equal(t, [][][]float64{{{1, 2}, {3, 4}}}, result[PropertyCoordinates])
}

func TestPolygon_MarshalSlice(t *testing.T) {

	polygon := NewPolygon(
		NewPosition(1, 2),
		NewPosition(3, 4),
	)

	require.Equal(t, [][]float64{{1, 2}, {3, 4}}, polygon.MarshalSlice())
}

func TestPolygon_MarshalStruct(t *testing.T) {

	polygon := NewPolygon(
		NewPosition(1, 2),
		NewPosition(3, 4),
	)

	result := polygon.MarshalStruct()
	require.Equal(t, PropertyTypePolygon, result.Type)
	require.Equal(t, [][][]float64{{{1, 2}, {3, 4}}}, result.Coordinates)
}

func TestPolygon_UnmarshalStruct_Errors(t *testing.T) {

	// Coordinates must contain exactly one ring
	polygon := Polygon{}
	require.NotNil(t, polygon.UnmarshalStruct(GeoJSONPolygon{
		Type:        PropertyTypePolygon,
		Coordinates: [][][]float64{},
	}))

	require.NotNil(t, polygon.UnmarshalStruct(GeoJSONPolygon{
		Type: PropertyTypePolygon,
		Coordinates: [][][]float64{
			{{1, 2}},
			{{3, 4}},
		},
	}))

	// A coordinate inside the ring has an invalid length
	require.NotNil(t, polygon.UnmarshalStruct(GeoJSONPolygon{
		Type: PropertyTypePolygon,
		Coordinates: [][][]float64{
			{{1}},
		},
	}))
}

func TestPolygon_UnmarshalJSON_Errors(t *testing.T) {

	polygon := Polygon{}

	// Malformed JSON
	require.NotNil(t, polygon.UnmarshalJSON([]byte("not json")))

	// Valid JSON, but invalid Polygon structure (zero rings)
	require.NotNil(t, polygon.UnmarshalJSON([]byte(`{"type":"Polygon","coordinates":[]}`)))
}

func TestPolygon_UnmarshalBSON_Error(t *testing.T) {

	polygon := Polygon{}
	require.NotNil(t, polygon.UnmarshalBSON([]byte("not bson")))
}

func TestPolygon_UnmarshalBSON_InvalidStructure(t *testing.T) {

	// Well-formed BSON that decodes, but has the wrong number of rings
	data, err := bson.Marshal(GeoJSONPolygon{Type: PropertyTypePolygon, Coordinates: [][][]float64{}})
	require.Nil(t, err)

	polygon := Polygon{}
	require.NotNil(t, polygon.UnmarshalBSON(data))
}

func TestPolygon_JSON(t *testing.T) {

	p1 := NewPolygon(
		NewPosition(1, 2),
		NewPosition(3, 4),
		NewPosition(5, 6),
		NewPosition(7, 8),
	)

	data, err1 := json.Marshal(p1)
	require.Nil(t, err1)

	p2 := Polygon{}

	err2 := json.Unmarshal(data, &p2)
	require.Nil(t, err2)
	require.Equal(t, p1, p2)
	require.Equal(t, 4, len(p2.Coordinates))
	require.Equal(
		t,
		Position{Longitude: 1, Latitude: 2},
		p2.Coordinates[0],
	)

	require.Equal(
		t,
		Position{Longitude: 3, Latitude: 4},
		p2.Coordinates[1],
	)

	require.Equal(
		t,
		Position{Longitude: 5, Latitude: 6},
		p2.Coordinates[2],
	)

	require.Equal(
		t,
		Position{Longitude: 7, Latitude: 8},
		p2.Coordinates[3],
	)
}

func TestPolygon_JSON_OmitZero(t *testing.T) {

	p1 := NewPolygon()

	data, err1 := json.Marshal(p1)
	require.Nil(t, err1)
	require.Equal(t, "null", string(data))
}

func TestPolygon_JSON_OmitZero_Struct(t *testing.T) {

	mystruct := struct {
		Title   string  `json:"title"`
		Polygon Polygon `json:"polygon,omitzero"`
	}{
		Title: "test",
	}

	data, err1 := json.Marshal(mystruct)
	require.Nil(t, err1)
	require.Equal(t, `{"title":"test"}`, string(data))
}

func TestPolygon_BSON(t *testing.T) {

	p1 := NewPolygon(
		NewPosition(1, 2),
		NewPosition(3, 4),
		NewPosition(5, 6),
		NewPosition(7, 8),
	)

	data, err1 := bson.Marshal(p1)
	require.Nil(t, err1)

	p2 := Polygon{}

	err2 := bson.Unmarshal(data, &p2)
	require.Nil(t, err2)
	require.Equal(t, p1, p2)
	require.Equal(t, 4, len(p2.Coordinates))
	require.Equal(
		t,
		Position{Longitude: 1, Latitude: 2},
		p2.Coordinates[0],
	)

	require.Equal(
		t,
		Position{Longitude: 3, Latitude: 4},
		p2.Coordinates[1],
	)

	require.Equal(
		t,
		Position{Longitude: 5, Latitude: 6},
		p2.Coordinates[2],
	)

	require.Equal(
		t,
		Position{Longitude: 7, Latitude: 8},
		p2.Coordinates[3],
	)
}

func TestPolygon_BSON_Empty(t *testing.T) {

	mystruct := struct {
		Title   string  `bson:"title"`
		Polygon Polygon `bson:"polygon"`
	}{
		Title: "test",
	}

	data, err1 := bson.Marshal(mystruct)
	require.Nil(t, err1)
	t.Log(data)
}

func TestPolygon_BSON_OmitEmpty(t *testing.T) {

	mystruct := struct {
		Title   string  `bson:"title"`
		Polygon Polygon `bson:"polygon,omitempty"`
	}{
		Title: "test",
	}

	data, err1 := bson.Marshal(mystruct)
	require.Nil(t, err1)
	t.Log(data)
}

func TestPolygon_UnmarshalStruct_ReportsCause(t *testing.T) {

	polygon := Polygon{}
	err := polygon.UnmarshalStruct(GeoJSONPolygon{
		Type:        PropertyTypePolygon,
		Coordinates: [][][]float64{{{1, 2}, {3}}},
	})

	// The error must name WHY the coordinate was rejected, not just that it was
	require.NotNil(t, err)

	cause := errors.Unwrap(err)
	require.NotNil(t, cause)
	require.Contains(t, cause.Error(), "Coordinates must be length 2 or 3")
}

func TestPolygon_UnmarshalBSON_WrongFieldType(t *testing.T) {

	// A well-formed BSON document whose "coordinates" field is a string. It
	// clears validation and fails in the decoder, which is the only way to
	// reach that branch now that malformed documents are rejected earlier.
	data, err := bson.Marshal(bson.M{"type": PropertyTypePolygon, "coordinates": "nope"})
	require.Nil(t, err)

	polygon := Polygon{}
	require.NotNil(t, polygon.UnmarshalBSON(data))
}

func TestPolygon_UnmarshalStruct_NilRing(t *testing.T) {

	// Exactly one ring, but the ring itself is nil. That is zero coordinates,
	// which is a legal (if empty) result -- and must not panic.
	polygon := Polygon{}
	require.Nil(t, polygon.UnmarshalStruct(GeoJSONPolygon{
		Type:        PropertyTypePolygon,
		Coordinates: [][][]float64{nil},
	}))
	require.Zero(t, len(polygon.Coordinates))
}

func TestPolygon_UnmarshalStruct_NilCoordinates(t *testing.T) {

	// No rings at all fails the single-ring check
	polygon := Polygon{}
	require.NotNil(t, polygon.UnmarshalStruct(GeoJSONPolygon{Type: PropertyTypePolygon}))
}

func TestPolygon_UnmarshalStruct_TooManyRings(t *testing.T) {

	// GeoJSON allows a polygon with holes; this package does not
	polygon := Polygon{}
	require.NotNil(t, polygon.UnmarshalStruct(GeoJSONPolygon{
		Type:        PropertyTypePolygon,
		Coordinates: [][][]float64{{{1, 2}}, {{3, 4}}},
	}))
}

func TestPolygon_MarshalSlice_Zero(t *testing.T) {

	// A zero Polygon has no vertices, and marshalling it must not panic
	require.Zero(t, len(Polygon{}.MarshalSlice()))
}

func TestNewPolygonFromString_Edges(t *testing.T) {

	// check confirms a coordinate string parses into the expected vertex count
	check := func(input string, expected int) {
		polygon := NewPolygonFromString(input)
		require.Equal(t, expected, len(polygon.Coordinates), "parsing %q", input)
	}

	check("", 0)             // nothing at all
	check("1", 0)            // a single unpaired value
	check("1,2", 1)          // one pair
	check("1,2,3", 1)        // trailing odd value is dropped
	check(",,,,", 2)         // empty values parse as zero
	check("not,a,number", 1) // unparseable values parse as zero
	check("1,2,3,4,5,6", 3)  // three pairs
}

func TestPolygon_String_RoundTrip(t *testing.T) {

	// The string form and its parser are inverses for coordinate pairs
	original := NewPolygon(NewPosition(-118.40853, 34.02116), NewPosition(1.5, -2.25))

	require.Equal(t, "-118.40853,34.02116,1.5,-2.25", original.String())
	require.Equal(t, original, NewPolygonFromString(original.String()))
}

func TestPolygon_String_DropsAltitude(t *testing.T) {

	// Position.String is a coordinate pair, so a Polygon's string form cannot
	// carry altitude and the round trip is lossy for 3D positions
	original := NewPolygon(NewPositionWithAltitude(1, 2, 300))

	require.Equal(t, "1,2", original.String())
	require.Equal(t, NewPolygon(NewPosition(1, 2)), NewPolygonFromString(original.String()))
}
