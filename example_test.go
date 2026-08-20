package geo

import (
	"encoding/json"
	"fmt"
)

// Coordinates are longitude-first, per the GeoJSON specification -- not the
// "latitude, longitude" order that people usually say aloud.
func ExampleNewPosition() {

	position := NewPosition(-118.40853, 34.02116)

	fmt.Println(position.String())
	// Output: -118.40853,34.02116
}

func ExampleNewPositionWithAltitude() {

	position := NewPositionWithAltitude(-118.40853, 34.02116, 38.4)

	data, _ := json.Marshal(position)
	fmt.Println(string(data))
	// Output: [-118.40853,34.02116,38.4]
}

// A Position marshals as a bare GeoJSON coordinate array, not an object.
func ExamplePosition_MarshalJSON() {

	data, _ := json.Marshal(NewPosition(1, 2))

	fmt.Println(string(data))
	// Output: [1,2]
}

func ExampleNewPoint() {

	point := NewPoint(-118.40853, 34.02116)

	data, _ := json.Marshal(point)
	fmt.Println(string(data))
	// Output: {"type":"Point","coordinates":[-118.40853,34.02116]}
}

// A zero Point marshals to null, so it disappears from a field tagged
// `json:",omitzero"`.
func ExamplePoint_MarshalJSON_zero() {

	data, _ := json.Marshal(Point{})

	fmt.Println(string(data))
	// Output: null
}

// LonLat and LatLon exist because consumers disagree about coordinate order.
func ExamplePoint_LonLat() {

	point := NewPoint(-118.40853, 34.02116)

	fmt.Println(point.LonLat())
	fmt.Println(point.LatLon())
	// Output:
	// -118.4085300000,34.0211600000
	// 34.0211600000,-118.4085300000
}

func ExampleNewPolygon() {

	polygon := NewPolygon(
		NewPosition(-118.5, 34.0),
		NewPosition(-118.0, 34.0),
		NewPosition(-118.0, 34.5),
	)

	data, _ := json.Marshal(polygon)
	fmt.Println(string(data))
	// Output: {"type":"Polygon","coordinates":[[[-118.5,34],[-118,34],[-118,34.5]]]}
}

// Unparseable values and an unpaired trailing coordinate are silently dropped.
func ExampleNewPolygonFromString() {

	polygon := NewPolygonFromString("-118.5,34.0,-118.0,34.0")

	fmt.Println(polygon.String())
	// Output: -118.5,34,-118,34
}

func ExampleAddress_LatLon() {

	address := Address{Latitude: 34.02116, Longitude: -118.40853}

	fmt.Println(address.LatLon())
	// Output: 34.0211600000,-118.4085300000
}

// JSONLD emits a schema.org Place, omitting every field that is empty.
func ExampleAddress_JSONLD() {

	address := Address{
		Name:     "Santa Monica Pier",
		Street1:  "200 Santa Monica Pier",
		Locality: "Santa Monica",
		Region:   "CA",
	}

	data, _ := json.Marshal(address.JSONLD())
	fmt.Println(string(data))
	// Output: {"address":{"locality":"Santa Monica","region":"CA","street1":"200 Santa Monica Pier"},"name":"Santa Monica Pier","type":"Place"}
}

// Setting "formatted" to a new value clears every geocoded field, because a
// new formatted address invalidates the old geocode.
func ExampleAddress_SetString() {

	address := Address{
		Formatted: "123 Main St",
		Locality:  "Springfield",
		Latitude:  39.7817,
	}

	address.SetString(AddressPropertyFormatted, "456 Oak Ave")

	fmt.Println(address.Formatted)
	fmt.Printf("%q %v\n", address.Locality, address.Latitude)
	// Output:
	// 456 Oak Ave
	// "" 0
}
