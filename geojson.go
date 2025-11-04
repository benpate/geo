package geo

// GeoJSONPoint represents the "strict" format for a Point in GeoJSON
type GeoJSONPoint struct {
	Type        string    `json:"type"        bson:"type"`        // This should always be "Point"
	Coordinates []float64 `json:"coordinates" bson:"coordinates"` // Whatevs
}

// GeoJSONPolygon represents the "strict" format for a Polygon in GeoJSON.
// is is used here to simplify conversion to/from serialization formats
type GeoJSONPolygon struct {
	Type        string        `json:"type"        bson:"type"`        // this should always be "Polygon"
	Coordinates [][][]float64 `json:"coordinates" bson:"coordinates"` // ick. Thanks IETF.
}
