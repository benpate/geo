package geo

import (
	"encoding/json"
	"strconv"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/sliceof"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// formatCoordinatePair formats two coordinates as "first,second", each to ten
// decimal places. It backs the LonLat/LatLon methods on Point and Address.
func formatCoordinatePair(first float64, second float64) string {
	return strconv.FormatFloat(first, 'f', 10, 64) + "," + strconv.FormatFloat(second, 'f', 10, 64)
}

// Position represents a Longitude/Latitude pair
// https://datatracker.ietf.org/doc/html/rfc7946#section-3.1.1
type Position struct {
	Longitude float64
	Latitude  float64
	Altitude  float64
}

// NewPosition returns a Position at the given longitude and latitude (no altitude).
func NewPosition(longitude float64, latitude float64) Position {

	return Position{
		Longitude: longitude,
		Latitude:  latitude,
		Altitude:  0,
	}
}

// NewPositionWithAltitude returns a Position at the given longitude, latitude, and altitude.
func NewPositionWithAltitude(longitude float64, latitude float64, altitude float64) Position {

	return Position{
		Longitude: longitude,
		Latitude:  latitude,
		Altitude:  altitude,
	}
}

// IsZero returns TRUE if this is a Zero position
func (position Position) IsZero() bool {
	return (position.Longitude == 0) && (position.Latitude == 0) && (position.Altitude == 0)
}

// NotZero returns TRUE if this Position is not Zero
func (position Position) NotZero() bool {
	return !position.IsZero()
}

/******************************************
 * Marshalling methods
 ******************************************/

// String returns a string representation of this coordinate pair
func (position Position) String() string {
	// Use strconv with -1 precision so coordinates round-trip losslessly,
	// instead of convert.String, which rounds floats to two decimal places.
	return strconv.FormatFloat(position.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(position.Latitude, 'f', -1, 64)
}

// MarshalSlice returns a longitude/latitude coordinate pair
func (position Position) MarshalSlice() []float64 {

	if position.Altitude == 0 {
		return []float64{position.Longitude, position.Latitude}
	}

	return []float64{position.Longitude, position.Latitude, position.Altitude}
}

// MarshalJSON is a custom JSON marshaller that serializes this
// Position into a GeoJSON coordinate pair
func (position Position) MarshalJSON() ([]byte, error) {
	return json.Marshal(position.MarshalSlice())
}

// MarshalBSONValue is a custom BSON marshaller that serializes this Position
// into a GeoJSON coordinate pair. It implements bson.ValueMarshaler (rather than
// bson.Marshaler) because a coordinate pair is a BSON array, not a document.
func (position Position) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(position.MarshalSlice())
}

/******************************************
 * Unmarshalling methods
 ******************************************/

// UnmarshalSlice populates this Position from a coordinate slice of length 2
// (longitude, latitude) or length 3 (longitude, latitude, altitude).
func (position *Position) UnmarshalSlice(coordinates sliceof.Float) error {

	const location = "geo.Position.UnmarshalSlice"

	switch coordinates.Length() {

	case 2:
		position.Longitude = coordinates[0]
		position.Latitude = coordinates[1]
		position.Altitude = 0
		return nil

	case 3:
		position.Longitude = coordinates[0]
		position.Latitude = coordinates[1]
		position.Altitude = coordinates[2]
		return nil
	}

	return derp.Internal(location, "Invalid coordinate length. Coordinates must be length 2 or 3", coordinates)
}

// UnmarshalJSON is a custom JSON unmarshaller that deserializes
// a GeoJSON coordinate pair into this Position structure.
func (position *Position) UnmarshalJSON(data []byte) error {

	const location = "geo.Position.UnmarshalJSON"

	// Unmarshal into a temporary array
	intermediate := make(sliceof.Float, 0, 3)

	if err := json.Unmarshal(data, &intermediate); err != nil {
		return derp.Wrap(err, location, "Unable to unmarshal JSON")
	}

	// Unmarshal into Position
	if err := position.UnmarshalSlice(intermediate); err != nil {
		return derp.Wrap(err, location, "Unable to unmarshal from intermediate slice", intermediate)
	}

	return nil
}

// UnmarshalBSONValue is a custom BSON unmarshaller that deserializes a BSON /
// GeoJSON coordinate pair into this Position structure. It implements
// bson.ValueUnmarshaler to match MarshalBSONValue's array encoding.
func (position *Position) UnmarshalBSONValue(dataType bsontype.Type, data []byte) error {

	const location = "geo.Position.UnmarshalBSONValue"

	// Unmarshal into a temporary array
	intermediate := make(sliceof.Float, 0, 3)

	if err := bson.UnmarshalValue(dataType, data, &intermediate); err != nil {
		return derp.Wrap(err, location, "Unable to unmarshal BSON")
	}

	// Unmarshal into Position
	if err := position.UnmarshalSlice(intermediate); err != nil {
		return derp.Wrap(err, location, "Unable to unmarshal from intermediate slice", intermediate)
	}

	return nil
}
