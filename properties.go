package geo

// GeoJSON property names used in marshalled output.
const (
	// PropertyType is the GeoJSON "type" property.
	PropertyType = "type"

	// PropertyCoordinates is the GeoJSON "coordinates" property.
	PropertyCoordinates = "coordinates"
)

// GeoJSON "type" values supported by this package.
const (
	// PropertyTypePoint is the GeoJSON type value for a Point.
	PropertyTypePoint = "Point"

	// PropertyTypePolygon is the GeoJSON type value for a Polygon.
	PropertyTypePolygon = "Polygon"
)
