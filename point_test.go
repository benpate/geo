package geo

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/benpate/rosetta/mapof"
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

func TestPoint_UnmarshalMap_Empty(t *testing.T) {

	// A nil or empty map has no "type", so it must be rejected -- and must
	// not panic reaching for keys that are not there
	point := Point{}
	require.NotNil(t, point.UnmarshalMap(nil))
	require.NotNil(t, point.UnmarshalMap(mapof.NewAny()))
}

func TestPoint_UnmarshalMap_MissingCoordinates(t *testing.T) {

	// The right type, but nothing to read -- GetSliceOfFloat returns an empty
	// slice, which is not a valid coordinate length
	point := Point{}
	require.NotNil(t, point.UnmarshalMap(mapof.Any{PropertyType: PropertyTypePoint}))
}

func TestPoint_Marshal_NaN(t *testing.T) {

	// A coordinate JSON cannot represent must fail the marshal rather than
	// emit a document with a silently wrong location
	point := NewPoint(math.NaN(), 0)
	_, err := json.Marshal(point)
	require.NotNil(t, err)
}

func TestPoint_LonLat_Extremes(t *testing.T) {

	// The fixed ten-decimal format must hold at the corners of the world
	require.Equal(t, "180.0000000000,90.0000000000", NewPoint(180, 90).LonLat())
	require.Equal(t, "90.0000000000,180.0000000000", NewPoint(180, 90).LatLon())
	require.Equal(t, "-180.0000000000,-90.0000000000", NewPoint(-180, -90).LonLat())

	// ...and must not lose sub-metre precision
	require.Equal(t, "-118.4085300000,34.0211600000", NewPoint(-118.40853, 34.02116).LonLat())
}

/******************************************
 * BSON format tests
 *
 * Point embeds Position, so Position's ValueMarshaller is promoted
 * into Point's method set. These tests pin which one wins, in every
 * context the driver can reach a Point through.
 ******************************************/

// Point must satisfy BOTH BSON interfaces. Dropping either one silently changes
// which encoder the driver resolves. See AGENTS.md.
var (
	_ bson.ValueMarshaler   = Point{}
	_ bson.Marshaler        = Point{}
	_ bson.ValueUnmarshaler = &Point{}
	_ bson.Unmarshaler      = &Point{}
)

// TestPoint_ShadowsPromotedPositionMarshaller asserts that Point's own DOCUMENT
// encoder wins over the bare ARRAY promoted from Position.
func TestPoint_ShadowsPromotedPositionMarshaller(t *testing.T) {

	pointType, pointData, err := NewPoint(1, 2).MarshalBSONValue()
	require.Nil(t, err)
	require.Equal(t, bson.TypeEmbeddedDocument, pointType, "Point must marshal as a GeoJSON document")

	// The promoted method it shadows still marshals a bare array, as it should
	positionType, positionData, err := NewPosition(1, 2).MarshalBSONValue()
	require.Nil(t, err)
	require.Equal(t, bson.TypeArray, positionType, "Position must marshal as a bare array")

	require.NotEqual(t, pointData, positionData, "Point must not fall through to Position")
}

// TestPoint_BSONInterfacesAgree pins the two encoders to identical bytes, so it
// cannot matter which one the driver resolves.
func TestPoint_BSONInterfacesAgree(t *testing.T) {

	point := NewPointWithAltitude(1, 2, 3)

	_, viaValue, err := point.MarshalBSONValue()
	require.Nil(t, err)

	viaDocument, err := point.MarshalBSON()
	require.Nil(t, err)

	require.Equal(t, viaDocument, viaValue)
}

// TestPoint_FieldContextMatrix marshals a Point through every context the driver
// resolves a codec for differently.
func TestPoint_FieldContextMatrix(t *testing.T) {

	point := NewPoint(1, 2)
	geoJSON := `{"type": "Point","coordinates": [{"$numberDouble":"1.0"},{"$numberDouble":"2.0"}]}`

	// check marshals one container and asserts the Extended JSON it produces
	check := func(name string, container any, expected string) {

		t.Run(name, func(t *testing.T) {
			data, err := bson.Marshal(container)
			require.Nil(t, err)
			require.Equal(t, expected, bson.Raw(data).String())
		})
	}

	check("StructField",
		struct {
			Location Point `bson:"location"`
		}{point},
		`{"location": `+geoJSON+`}`)

	check("PointerField",
		struct {
			Location *Point `bson:"location"`
		}{&point},
		`{"location": `+geoJSON+`}`)

	check("SliceElement",
		struct {
			Locations []Point `bson:"locations"`
		}{[]Point{point}},
		`{"locations": [`+geoJSON+`]}`)

	check("MapValue",
		struct {
			Locations map[string]Point `bson:"locations"`
		}{map[string]Point{"venue": point}},
		`{"locations": {"venue": `+geoJSON+`}}`)

	check("NestedStruct",
		struct {
			Venue struct {
				Location Point `bson:"location"`
			} `bson:"venue"`
		}{Venue: struct {
			Location Point `bson:"location"`
		}{point}},
		`{"venue": {"location": `+geoJSON+`}}`)

	// A standalone Point is a document in its own right
	t.Run("TopLevel", func(t *testing.T) {
		data, err := bson.Marshal(point)
		require.Nil(t, err)
		require.Equal(t, geoJSON, bson.Raw(data).String())
	})
}

// TestPoint_UnmarshalBSONValue_AllShapes covers every BSON type the decoder can
// be handed, including the two historical on-disk shapes.
func TestPoint_UnmarshalBSONValue_AllShapes(t *testing.T) {

	// decode runs one raw BSON value through a Point field on a real document
	decode := func(value any) (Point, error) {

		data, err := bson.Marshal(bson.M{"location": value})
		require.Nil(t, err)

		container := struct {
			Location Point `bson:"location"`
		}{}

		return container.Location, bson.Unmarshal(data, &container)
	}

	t.Run("GeoJSONDocument", func(t *testing.T) {
		point, err := decode(bson.M{"type": "Point", "coordinates": bson.A{1.0, 2.0}})
		require.Nil(t, err)
		require.Equal(t, NewPoint(1, 2), point)
	})

	t.Run("GeoJSONDocumentWithAltitude", func(t *testing.T) {
		point, err := decode(bson.M{"type": "Point", "coordinates": bson.A{1.0, 2.0, 3.0}})
		require.Nil(t, err)
		require.Equal(t, NewPointWithAltitude(1, 2, 3), point)
	})

	t.Run("LegacyBareArray", func(t *testing.T) {
		point, err := decode(bson.A{1.0, 2.0})
		require.Nil(t, err)
		require.Equal(t, NewPoint(1, 2), point)
	})

	t.Run("LegacyBareArrayWithAltitude", func(t *testing.T) {
		point, err := decode(bson.A{1.0, 2.0, 3.0})
		require.Nil(t, err)
		require.Equal(t, NewPointWithAltitude(1, 2, 3), point)
	})

	t.Run("Null", func(t *testing.T) {
		point, err := decode(nil)
		require.Nil(t, err)
		require.True(t, point.IsZero())
	})

	t.Run("MissingField", func(t *testing.T) {
		data, err := bson.Marshal(bson.M{"other": 1})
		require.Nil(t, err)

		container := struct {
			Location Point `bson:"location"`
		}{}
		require.Nil(t, bson.Unmarshal(data, &container))
		require.True(t, container.Location.IsZero())
	})

	t.Run("WrongBSONType", func(t *testing.T) {
		_, err := decode("not a point")
		require.NotNil(t, err)
	})

	t.Run("ArrayOfWrongLength", func(t *testing.T) {
		_, err := decode(bson.A{1.0})
		require.NotNil(t, err)
	})

	t.Run("DocumentOfWrongGeoJSONType", func(t *testing.T) {
		_, err := decode(bson.M{"type": "Polygon", "coordinates": bson.A{1.0, 2.0}})
		require.NotNil(t, err)
	})

	t.Run("Undefined", func(t *testing.T) {
		point := Point{}
		require.Nil(t, point.UnmarshalBSONValue(bson.TypeUndefined, nil))
		require.True(t, point.IsZero())
	})
}

// TestPoint_JSONAndBSONAgree pins the two wire formats side by side. They are
// allowed to differ; they are not allowed to differ by accident.
func TestPoint_JSONAndBSONAgree(t *testing.T) {

	point := NewPoint(1, 2)

	asJSON, err := json.Marshal(point)
	require.Nil(t, err)
	require.JSONEq(t, `{"type":"Point","coordinates":[1,2]}`, string(asJSON))

	asBSON, err := bson.Marshal(struct {
		Location Point `bson:"location"`
	}{point})
	require.Nil(t, err)
	require.Equal(t,
		`{"location": {"type": "Point","coordinates": [{"$numberDouble":"1.0"},{"$numberDouble":"2.0"}]}}`,
		bson.Raw(asBSON).String())
}
