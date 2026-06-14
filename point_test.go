package geo

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestPoint_Zeroer(t *testing.T) {
	require.True(t, NewPoint(0, 0).IsZero())
	require.False(t, NewPoint(1, 0).IsZero())
}

func TestNewPoint(t *testing.T) {
	require.Equal(t, Point{Position: NewPosition(1, 2)}, NewPoint(1, 2))
	require.Equal(t, Point{Position: NewPositionWithAltitude(1, 2, 3)}, NewPointWithAltitude(1, 2, 3))
}

func TestPoint_LonLat(t *testing.T) {
	point := NewPoint(-118.5, 34.25)
	require.Equal(t, "-118.5000000000,34.2500000000", point.LonLat())
}

func TestPoint_LatLon(t *testing.T) {
	point := NewPoint(-118.5, 34.25)
	require.Equal(t, "34.2500000000,-118.5000000000", point.LatLon())
}

func TestPoint_GeoJSON(t *testing.T) {
	result := NewPoint(1, 2).GeoJSON()
	require.Equal(t, PropertyTypePoint, result[PropertyType])
	require.Equal(t, []float64{1, 2}, result[PropertyCoordinates])
}

func TestPoint_MarshalStruct(t *testing.T) {
	result := NewPointWithAltitude(1, 2, 3).MarshalStruct()
	require.Equal(t, PropertyTypePoint, result.Type)
	require.Equal(t, []float64{1, 2, 3}, result.Coordinates)
}

func TestPoint_UnmarshalMap(t *testing.T) {

	point := Point{}
	err := point.UnmarshalMap(map[string]any{
		PropertyType:        PropertyTypePoint,
		PropertyCoordinates: []float64{1, 2},
	})

	require.Nil(t, err)
	require.Equal(t, NewPoint(1, 2), point)
}

func TestPoint_UnmarshalMap_WrongType(t *testing.T) {

	point := Point{}
	err := point.UnmarshalMap(map[string]any{
		PropertyType:        PropertyTypePolygon,
		PropertyCoordinates: []float64{1, 2},
	})

	require.NotNil(t, err)
}

func TestPoint_UnmarshalMap_BadCoordinates(t *testing.T) {

	point := Point{}
	err := point.UnmarshalMap(map[string]any{
		PropertyType:        PropertyTypePoint,
		PropertyCoordinates: []float64{1}, // too few coordinates
	})

	require.NotNil(t, err)
}

func TestPoint_UnmarshalJSON_Errors(t *testing.T) {

	point := Point{}

	// Malformed JSON
	require.NotNil(t, point.UnmarshalJSON([]byte("not json")))

	// Valid JSON, wrong GeoJSON type
	require.NotNil(t, point.UnmarshalJSON([]byte(`{"type":"Polygon","coordinates":[1,2]}`)))
}

func TestPoint_UnmarshalBSON_Error(t *testing.T) {

	point := Point{}
	require.NotNil(t, point.UnmarshalBSON([]byte("not bson")))
}

func TestPoint_UnmarshalBSON_WrongType(t *testing.T) {

	// Well-formed BSON that decodes into a map, but is not a valid Point
	data, err := bson.Marshal(GeoJSONPoint{Type: PropertyTypePolygon, Coordinates: []float64{1, 2}})
	require.Nil(t, err)

	point := Point{}
	require.NotNil(t, point.UnmarshalBSON(data))
}

func TestPoint_JSON(t *testing.T) {

	p1 := NewPoint(1, 2)

	data, err1 := json.Marshal(p1)
	require.Nil(t, err1)

	p2 := Point{}
	err2 := json.Unmarshal(data, &p2)

	require.Nil(t, err2)
	require.Equal(t, p1, p2)
	require.Equal(t, 1.0, p2.Longitude)
	require.Equal(t, 2.0, p2.Latitude)
	require.Equal(t, 0.0, p2.Altitude)
}

func TestPoint_JSON_WithAltitude(t *testing.T) {

	p1 := NewPointWithAltitude(1, 2, 3)

	data, err1 := json.Marshal(p1)
	require.Nil(t, err1)

	p2 := Point{}
	err2 := json.Unmarshal(data, &p2)

	require.Nil(t, err2)
	require.Equal(t, p1, p2)
	require.Equal(t, 1.0, p2.Longitude)
	require.Equal(t, 2.0, p2.Latitude)
	require.Equal(t, 3.0, p2.Altitude)
}

func TestPoint_JSON_Nil(t *testing.T) {

	p1 := Point{}

	data, err1 := json.Marshal(p1)
	require.Nil(t, err1)
	require.Equal(t, "null", string(data))
}

func TestPoint_JSON_OmitZero(t *testing.T) {

	mystruct := struct {
		Title string `json:"title"`
		Point Point  `json:"point,omitzero"`
	}{
		Title: "test",
	}

	data, err1 := json.Marshal(mystruct)
	require.Nil(t, err1)
	require.Equal(t, `{"title":"test"}`, string(data))
}

func TestPoint_BSON(t *testing.T) {

	p1 := NewPoint(1, 2)

	data, err1 := bson.Marshal(p1)
	require.Nil(t, err1)

	p2 := Point{}
	err2 := bson.Unmarshal(data, &p2)

	require.Nil(t, err2)
	require.Equal(t, p1, p2)
	require.Equal(t, 1.0, p2.Longitude)
	require.Equal(t, 2.0, p2.Latitude)
	require.Equal(t, 0.0, p2.Altitude)
}

func TestPoint_BSON_WithAltitude(t *testing.T) {

	p1 := NewPointWithAltitude(1, 2, 3)

	data, err1 := bson.Marshal(p1)
	require.Nil(t, err1)

	p2 := Point{}
	err2 := bson.Unmarshal(data, &p2)

	require.Nil(t, err2)
	require.Equal(t, p1, p2)
	require.Equal(t, 1.0, p2.Longitude)
	require.Equal(t, 2.0, p2.Latitude)
	require.Equal(t, 3.0, p2.Altitude)
}

func TestPoint_BSON_OmitEmpty(t *testing.T) {

	p1 := struct {
		Title string `bson:"title"`
		Point Point  `bson:"point,omitempty"`
	}{
		Title: "test bson",
	}

	data, err1 := bson.Marshal(p1)
	require.Nil(t, err1)

	p2 := struct {
		Title string `bson:"title"`
		Point Point  `bson:"point"`
	}{}

	err2 := bson.Unmarshal(data, &p2)
	require.Nil(t, err2)

	// Re-marshal the new data to show that it looks different w/o "omitempty"
	data2, err3 := bson.Marshal(p2)
	require.Nil(t, err3)
	require.NotEqual(t, data, data2)
}
