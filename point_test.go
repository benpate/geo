package geo

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

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

func TestPoint_AltitudeJSON(t *testing.T) {

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

func TestPoint_AltitudeBSON(t *testing.T) {

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
