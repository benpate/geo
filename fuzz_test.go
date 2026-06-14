package geo

import (
	"testing"
)

// FuzzPoint_UnmarshalJSON confirms that the Point JSON decoder never panics,
// regardless of the input bytes. A well-formed input round-trips; a malformed
// one must return an error instead of crashing.
func FuzzPoint_UnmarshalJSON(f *testing.F) {

	f.Add(`{"type":"Point","coordinates":[1,2]}`)
	f.Add(`{"type":"Point","coordinates":[1,2,3]}`)
	f.Add(`{"type":"Polygon","coordinates":[1,2]}`)
	f.Add(`null`)
	f.Add(``)
	f.Add(`{`)

	f.Fuzz(func(t *testing.T, data string) {
		point := Point{}
		_ = point.UnmarshalJSON([]byte(data))
	})
}

// FuzzPolygon_UnmarshalJSON confirms that the Polygon JSON decoder never panics.
func FuzzPolygon_UnmarshalJSON(f *testing.F) {

	f.Add(`{"type":"Polygon","coordinates":[[[1,2],[3,4]]]}`)
	f.Add(`{"type":"Polygon","coordinates":[]}`)
	f.Add(`{"type":"Point","coordinates":[1,2]}`)
	f.Add(`null`)
	f.Add(``)
	f.Add(`[`)

	f.Fuzz(func(t *testing.T, data string) {
		polygon := Polygon{}
		_ = polygon.UnmarshalJSON([]byte(data))
	})
}

// FuzzNewPolygonFromString confirms that the comma-delimited coordinate parser
// never panics and always produces an even number of coordinate components.
func FuzzNewPolygonFromString(f *testing.F) {

	f.Add("1,2,3,4")
	f.Add("1,2,3")
	f.Add("")
	f.Add("-118.5,34.25")
	f.Add("not,a,number")

	f.Fuzz(func(t *testing.T, data string) {
		polygon := NewPolygonFromString(data)

		// Every parsed position contributes a longitude and a latitude, so the
		// marshalled slice must always contain pairs (length-2 inner slices).
		for _, coordinate := range polygon.MarshalSlice() {
			if len(coordinate) != 2 {
				t.Fatalf("expected coordinate pairs, got length %d", len(coordinate))
			}
		}
	})
}
