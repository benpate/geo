package geo

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

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

/*
func TestPosition_BSON(t *testing.T) {

	p1 := Position{
		Longitude: 1,
		Latitude:  2,
	}

	data, err1 := bson.Marshal(p1)
	spew.Dump(data, err1)
	require.Nil(t, err1)
	t.Log(data)
	// require.Equal(t, []byte("[1,2]"), data)

	p2 := Position{}
	err2 := bson.Unmarshal(data, &p2)

	require.Nil(t, err2)
	require.Equal(t, p1, p2)
}

func TestPosition_AltitudeBSON(t *testing.T) {

	p1 := Position{
		Longitude: 1,
		Latitude:  2,
		Altitude:  3,
	}

	data, err1 := bson.Marshal(p1)
	require.Nil(t, err1)
	t.Log(data)
	// require.Equal(t, []byte("[1,2,3]"), data)

	p2 := Position{}
	err2 := bson.Unmarshal(data, &p2)

	require.Nil(t, err2)
	require.Equal(t, p1, p2)
}
*/
