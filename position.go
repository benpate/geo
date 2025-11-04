package geo

import (
	"encoding/json"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/sliceof"
	"go.mongodb.org/mongo-driver/bson"
)

// Position represents a Longitude/Latitude pair
// https://datatracker.ietf.org/doc/html/rfc7946#section-3.1.1
type Position struct {
	Longitude float64
	Latitude  float64
	Altitude  float64
}

func NewPosition(longitude float64, latitude float64) Position {

	return Position{
		Longitude: longitude,
		Latitude:  latitude,
		Altitude:  0,
	}
}

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

/******************************************
 * Marshalling methods
 ******************************************/

// MarshalSlice returns a longitude/latitude coordinate pair
func (position *Position) MarshalSlice() []float64 {

	if position.Altitude == 0 {
		return []float64{position.Longitude, position.Latitude}
	}

	return []float64{position.Longitude, position.Latitude, position.Altitude}
}

// MarshalJSON is a custom JBSON marshaller that serializes this
// Position into a GeoJSON coordinate pair
func (position Position) MarshalJSON() ([]byte, error) {
	return json.Marshal(position.MarshalSlice())
}

// MarshalBSON is a custom BSON marshaller that serializes this
// Position into a GeoJSON coordinate pair
func (position Position) MarshalBSON() ([]byte, error) {
	return bson.Marshal(position.MarshalSlice())
}

/******************************************
 * Unmarshalling methods
 ******************************************/

func (position *Position) UnmarshalSlice(coordinates sliceof.Float) error {

	const location = "geo.Position.unmarshalSlice"

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

	return derp.Internal(location, "Invalid coordinate length.  Coordinates must be length 2", coordinates)
}

// UnmarshalBSON is a custom BSON unmarshaller that deserializes
// a BSON / GeoJSON coordinate pair into this Position structure.
func (position *Position) UnmarshalJSON(data []byte) error {

	const location = "geo.LatLng.UnmarshalJSON"

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

// UnmarshalBSON is a custom BSON unmarshaller that deserializes
// a BSON / GeoJSON coordinate pair into this Position structure.
func (position *Position) UnmarshalBSON(data []byte) error {

	const location = "geo.LatLng.UnmarshalBSON"

	// Unmarshal into a temporary array
	intermediate := make(sliceof.Float, 0, 3)

	if err := bson.Unmarshal(data, &intermediate); err != nil {
		return derp.Wrap(err, location, "Unable to unmarshal BSON")
	}

	// Unmarshal into Position
	if err := position.UnmarshalSlice(intermediate); err != nil {
		return derp.Wrap(err, location, "Unable to unmarshal from intermediate slice", intermediate)
	}

	return nil
}
