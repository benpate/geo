package geo

import (
	"encoding/json"
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/slice"
	"github.com/benpate/rosetta/sliceof"
	"go.mongodb.org/mongo-driver/bson"
)

// Polygon represents a GeoJSON "Polygon" object
// https://datatracker.ietf.org/doc/html/rfc7946#section-3.1.2
type Polygon struct {
	Coordinates sliceof.Object[Position]
}

func NewPolygon(coordinates ...Position) Polygon {
	return Polygon{
		Coordinates: coordinates,
	}
}

func NewPolygonFromString(data string) Polygon {

	// Parse the data as a slice of float64s
	coords := convert.SliceOfFloat(data)

	// Allocate a result that's half the length of the coordinate pairs
	result := make([]Position, 0, len(coords)/2)

	// Combiine coordinates into pairs
	for len(coords) > 1 {
		position := NewPosition(coords[0], coords[1])
		result = append(result, position)
		coords = coords[2:]
	}

	// UwU
	return NewPolygon(result...)
}

func (polygon Polygon) IsZero() bool {
	return polygon.Coordinates.IsZero()
}

func (polygon Polygon) NotZero() bool {
	return !polygon.IsZero()
}

/******************************************
 * Marhshalling methods
 ******************************************/

func (polygon Polygon) String() string {
	if polygon.IsZero() {
		return ""
	}

	// Combine all points into a single string and return
	result := slice.Map(polygon.Coordinates, position_string)
	return strings.Join(result, ",")
}

// GeoJSON returns a GeoJSON representation of this Polygon
func (polygon Polygon) GeoJSON() map[string]any {
	return map[string]any{
		PropertyType: PropertyTypePolygon,
		PropertyCoordinates: [][][]float64{
			polygon.MarshalSlice(),
		},
	}
}

// MarshalSlice returns (a slice of (a slice of floats)), which is the
// standard way of representing a GeoJSON polygon
func (polygon Polygon) MarshalSlice() [][]float64 {
	return slice.Map(polygon.Coordinates, position_slice)
}

// MarshalMap copies this Polygon into a mapof.Any
func (polygon Polygon) MarshalStruct() GeoJSONPolygon {

	return GeoJSONPolygon{
		Type: PropertyTypePolygon,
		Coordinates: [][][]float64{
			polygon.MarshalSlice(),
		},
	}
}

// MarshalJSON is a custom json.Marshaller that returns this Polygon
// as a GeoJSON object.
func (polygon Polygon) MarshalJSON() ([]byte, error) {

	if polygon.IsZero() {
		return json.Marshal(nil)
	}

	return json.Marshal(polygon.MarshalStruct())
}

// MarshalBSON is a custom BSON marshaller that serializes this
// Position into a GeoJSON coordinate pair
func (polygon Polygon) MarshalBSON() ([]byte, error) {
	return bson.Marshal(polygon.MarshalStruct())
}

/******************************************
 * Unmarhshalling methods
 ******************************************/

func (polygon *Polygon) UnmarshalStruct(data GeoJSONPolygon) error {

	const location = "geo.Polygon.UnmarshalStruct"

	// Validate Polygon length
	if len(data.Coordinates) != 1 {
		return derp.Internal(location, "Coordinates length must be 1", data.Coordinates)
	}

	// Initialize variable / clear existing values
	polygon.Coordinates = make(sliceof.Object[Position], len(data.Coordinates[0]))

	// Copy/translate coordinates into Position
	for index, coordinate := range data.Coordinates[0] {
		if err := polygon.Coordinates[index].UnmarshalSlice(coordinate); err != nil {
			return derp.Internal(location, "Invalid coordinate at index", index, coordinate)
		}
	}

	return nil
}

// UnmarshalJSON is a custom json.Unmarshaller that parses a GeoJSON
// object into this Polygon object.
func (polygon *Polygon) UnmarshalJSON(data []byte) error {

	const location = "geo.Polygon.UnmarshalJSON"

	// Unmarshall JSON into an intermediate object
	intermediate := GeoJSONPolygon{}

	if err := json.Unmarshal(data, &intermediate); err != nil {
		return derp.Wrap(err, location, "Unable to unmarshal original JSON", string(data))
	}

	// Unmarshal from intermediate object into this Polygon
	if err := polygon.UnmarshalStruct(intermediate); err != nil {
		return derp.Wrap(err, location, "Unable to unmarshal from struct", intermediate)
	}

	// I see this as an absolute win.
	return nil
}

// UnmarshalBSON is a custom BSON unmarshaller that deserializes
// a BSON / GeoJSON coordinate pair into this Position structure.
func (polygon *Polygon) UnmarshalBSON(data []byte) error {

	const location = "geo.LatLng.UnmarshalBSON"

	// Unmarshall BSON into an intermediate object
	intermediate := GeoJSONPolygon{}

	if err := bson.Unmarshal(data, &intermediate); err != nil {
		return derp.Wrap(err, location, "Unable to unmarshal original BSON", string(data))
	}

	// Unmarshal from intermediate object into this Polygon
	if err := polygon.UnmarshalStruct(intermediate); err != nil {
		return derp.Wrap(err, location, "Unable to unmarshal from struct", intermediate)
	}

	// I see this as an absolute win.
	return nil
}
