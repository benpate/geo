package geo

import (
	"encoding/json"
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
