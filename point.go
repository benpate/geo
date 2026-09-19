package geo

import (
	"encoding/json"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// Point represents a GeoJSON "Point" object
// https://datatracker.ietf.org/doc/html/rfc7946#section-3.1.2
type Point struct {
	Position
}

// NewPoint returns a Point at the given longitude and latitude (no altitude).
func NewPoint(longitude float64, latitude float64) Point {
	return Point{
		Position: NewPosition(longitude, latitude),
	}
}

// NewPointWithAltitude returns a Point at the given longitude, latitude, and altitude.
func NewPointWithAltitude(longitude float64, latitude float64, altitude float64) Point {
	return Point{
		Position: NewPositionWithAltitude(longitude, latitude, altitude),
	}
}

// LonLat returns the coordinates as a "longitude,latitude" string.
func (point Point) LonLat() string {
	return formatCoordinatePair(point.Longitude, point.Latitude)
}

// LatLon returns the coordinates as a "latitude,longitude" string.
func (point Point) LatLon() string {
	return formatCoordinatePair(point.Latitude, point.Longitude)
}

/******************************************
 * Marshalling methods
 ******************************************/

// GeoJSON returns this Point as a GeoJSON object (a "type"/"coordinates" map).
func (point Point) GeoJSON() map[string]any {

	// The UNNAMED map type is what satisfies exp.GeoJSONer, and keeps `exp`
	// free of a rosetta dependency. See AGENTS.md.
	return map[string]any{
		PropertyType:        PropertyTypePoint,
		PropertyCoordinates: point.MarshalSlice(),
	}
}

// MarshalStruct returns this Point as a strongly-typed GeoJSONPoint.
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

// MarshalBSONValue serializes this Point into a GeoJSON object.
// It implements bson.ValueMarshaler.
func (point Point) MarshalBSONValue() (bsontype.Type, []byte, error) {

	// The *Value form is required: Position's ValueMarshaler is promoted into Point
	// and the driver prefers it, so only a method declared here shadows it (AGENTS.md).
	return bson.MarshalValue(point.MarshalStruct())
}

// MarshalBSON serializes this Point into a GeoJSON object.
// It implements bson.Marshaler, for callers that marshal a Point on its own.
func (point Point) MarshalBSON() ([]byte, error) {

	// Shares MarshalStruct with MarshalBSONValue so the two interfaces cannot disagree
	return bson.Marshal(point.MarshalStruct())
}

/******************************************
 * Unmarshalling methods
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

	// Success
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

// UnmarshalBSONValue deserializes a GeoJSON object, or a bare coordinate
// array, into this Point. It implements bson.ValueUnmarshaler.
func (point *Point) UnmarshalBSONValue(dataType bsontype.Type, data []byte) error {

	const location = "geo.Point.UnmarshalBSONValue"

	switch dataType {

	// A GeoJSON document is the shape this package writes
	case bson.TypeEmbeddedDocument:

		// RULE: Reject a malformed document before handing it to the driver, which
		// panics on some of them instead of returning an error
		if err := validateBSONDocument(data); err != nil {
			return derp.Wrap(err, location, "Unable to unmarshal original BSON", string(data))
		}

		// Unmarshall BSON into an intermediate object. mapof.Any has no custom decoder, so
		// the driver's map codec accepts anything validateBSONDocument already admitted --
		// this error is unreachable, and is kept only so a driver change cannot slip past.
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

	// A bare coordinate array is the shape geo v0.1.0 - v0.2.0 wrote (AGENTS.md)
	case bson.TypeArray:
		return point.Position.UnmarshalBSONValue(dataType, data)

	// RULE: An absent or NULL Point is the zero Point, not an error
	case bson.TypeNull, bson.TypeUndefined:
		*point = Point{}
		return nil
	}

	return derp.Internal(location, "Cannot unmarshal this BSON type into a Point", dataType.String())
}

// UnmarshalBSON deserializes a GeoJSON object into this Point.
// It implements bson.Unmarshaler, for callers that unmarshal a Point on its own.
func (point *Point) UnmarshalBSON(data []byte) error {

	// A standalone BSON document is an EmbeddedDocument value by any other name
	return point.UnmarshalBSONValue(bson.TypeEmbeddedDocument, data)
}
