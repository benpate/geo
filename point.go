package geo

import (
	"encoding/json"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"go.mongodb.org/mongo-driver/bson"
)

// Point represents a GeoJSON "Point" object
// https://datatracker.ietf.org/doc/html/rfc7946#section-3.1.2
type Point struct {
	Position
}

func NewPoint(longitude float64, latitude float64) Point {
	return Point{
		Position: NewPosition(longitude, latitude),
	}
}

func NewPointWithAltitude(longitude float64, latitude float64, altitude float64) Point {
	return Point{
		Position: NewPositionWithAltitude(longitude, latitude, altitude),
	}
}

/******************************************
 * Marhshalling methods
 ******************************************/

func (point Point) GeoJSON() map[string]any {
	return map[string]any{
		PropertyType:        PropertyTypePoint,
		PropertyCoordinates: point.MarshalSlice(),
	}
}

func (point Point) MarshalStruct() GeoJSONPoint {
	return GeoJSONPoint{
		Type:        PropertyTypePoint,
		Coordinates: point.MarshalSlice(),
	}
}

// MarshalJSON is a custom json.Marshaller that returns this Point
// as a GeoJSON object. This marshaller works with `omitzero` but not
// `omitempty`
func (point Point) MarshalJSON() ([]byte, error) {

	if point.IsZero() {
		return json.Marshal(nil)
	}

	return json.Marshal(point.MarshalStruct())
}

// MarshalBSON is a custom BSON marshaller that serializes this
// Position into a GeoJSON coordinate pair
func (point Point) MarshalBSON() ([]byte, error) {
	return bson.Marshal(point.MarshalStruct())
}

/******************************************
 * Unmarhshalling methods
 ******************************************/

// UnmarshalMap populates this Point using the values from the provided data.
// If the data does not fit the correct GeoJSON format, then this method returns an error
func (point *Point) UnmarshalMap(data mapof.Any) error {

	const location = "geo.Point.UnmarshalMap"

	// Validate the "type" property
	if data.GetString(PropertyType) != PropertyTypePoint {
		return derp.Internal(location, "Invalid GeoJSON. Type must be 'Point'", data)
	}

	// Parse the coordinates
	coordinates := data.GetSliceOfFloat(PropertyCoordinates)

	if err := point.UnmarshalSlice(coordinates); err != nil {
		return derp.Wrap(err, location, "Unable to unmarshal coordinates", coordinates)
	}

	// Sucess
	return nil
}

// UnmarshalJSON is a custom json.Unmarshaller that parses a GeoJSON
// object into this Point object.
func (point *Point) UnmarshalJSON(data []byte) error {

	const location = "geo.Point.UnmarshalJSON"

	// Unmarshall JSON into an intermediate object
	intermediate := mapof.NewAny()

	if err := json.Unmarshal(data, &intermediate); err != nil {
		return derp.Wrap(err, location, "Unable to unmarshal original JSON", string(data))
	}

	// Unmarshal from intermediate object into this Point
	if err := point.UnmarshalMap(intermediate); err != nil {
		return derp.Wrap(err, location, "Unable to unmarshal from Map", intermediate)
	}

	// I see this as an absolute win.
	return nil
}

// UnmarshalBSON is a custom BSON unmarshaller that deserializes
// a BSON / GeoJSON coordinate pair into this Position structure.
func (point *Point) UnmarshalBSON(data []byte) error {

	const location = "geo.LatLng.UnmarshalBSON"

	// Unmarshall BSON into an intermediate object
	intermediate := mapof.NewAny()

	if err := bson.Unmarshal(data, &intermediate); err != nil {
		return derp.Wrap(err, location, "Unable to unmarshal original BSON", string(data))
	}

	// Unmarshal from intermediate object into this Point
	if err := point.UnmarshalMap(intermediate); err != nil {
		return derp.Wrap(err, location, "Unable to unmarshal from Map", intermediate)
	}

	// I see this as an absolute win.
	return nil
}
