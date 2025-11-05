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
