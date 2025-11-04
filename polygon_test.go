package geo

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

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
