package geo

import (
	"strconv"

	"github.com/benpate/rosetta/mapof"
)

// Address represents a physical address on the planet
// It maps to https://www.w3.org/TR/activitystreams-vocabulary/#dfn-address
// and uses https://schema.org/PostalAddress to match Mobilizion
type Address struct {
	Name       string  `json:"name"        bson:"name,omitempty"`       // Human-readable name of the address
	Formatted  string  `json:"formatted"   bson:"formatted,omitempty"`  // Full, unparsed value of the address
	Street1    string  `json:"street1"     bson:"street1,omitempty"`    // Parsed street address line 1 of the address
	Street2    string  `json:"street2"     bson:"street2,omitempty"`    // Parsed street address line 2 of the address
	Locality   string  `json:"locality"    bson:"locality,omitempty"`   // Parsed city or town of the address
	Region     string  `json:"region"      bson:"region,omitempty"`     // Parsed state or province of the address
	PostalCode string  `json:"postalCode"  bson:"postalCode,omitempty"` // Parsed postal code of the address
	Country    string  `json:"country"     bson:"country,omitempty"`    // Parsed country of the address
	PlusCode   string  `json:"plusCode"    bson:"plusCode,omitempty"`   // PlusCode for this location https://maps.google.com/pluscodes/
	Timezone   string  `json:"timezone"    bson:"timezone,omitempty"`   // Time zone in tzdatabase format (https://en.wikipedia.org/wiki/Tz_database)
	Latitude   float64 `json:"latitude"    bson:"latitude,omitempty"`   // Latitude of the address
	Longitude  float64 `json:"longitude"   bson:"longitude,omitempty"`  // Longitude of the address
}

func NewAddress() Address {
	return Address{}
}

// Reset clears all geocoding information from this Address
func (address *Address) Reset() {
	address.Street1 = ""
	address.Street2 = ""
	address.Locality = ""
	address.Region = ""
	address.PostalCode = ""
	address.Country = ""
	address.Timezone = ""
	address.PlusCode = ""
	address.Latitude = 0
	address.Longitude = 0
}

func (address Address) IsZero() bool {

	if address.Latitude != 0 {
		return false
	}

	if address.Longitude != 0 {
		return false
	}

	if address.Formatted != "" {
		return false
	}

	return true
}

func (address Address) NotZero() bool {
	return !address.IsZero()
}

// HasGeocode returns TRUE if this Address has ANY Lat/Long information
func (address Address) HasGeocode() bool {

	if address.Latitude != 0 {
		return true
	}

	if address.Longitude != 0 {
		return true
	}

	return false
}

// HasAddress returns TRUE if this Address has ANY street adsress information
func (address Address) HasAddress() bool {

	if address.Country != "" {
		return true
	}

	if address.PostalCode != "" {
		return true
	}

	if address.Region != "" {
		return true
	}

	if address.Locality != "" {
		return true
	}

	if address.Street1 != "" {
		return true
	}

	return false
}

func (address Address) LonLat() string {
	return strconv.FormatFloat(address.Longitude, 'f', 10, 64) + "," + strconv.FormatFloat(address.Latitude, 'f', 10, 64)
}

func (address Address) LatLon() string {
	return strconv.FormatFloat(address.Latitude, 'f', 10, 64) + "," + strconv.FormatFloat(address.Longitude, 'f', 10, 64)
}

func (address *Address) SetPoint(point Point) {
	address.Longitude = point.Longitude
	address.Latitude = point.Latitude
}

/******************************************
 * Marshalling Methods
 ******************************************/

func (address *Address) MarshalMap() mapof.Any {

	return mapof.Any{
		AddressPropertyName:       address.Name,
		AddressPropertyFormatted:  address.Formatted,
		AddressPropertyStreet1:    address.Street1,
		AddressPropertyStreet2:    address.Street2,
		AddressPropertyLocality:   address.Locality,
		AddressPropertyRegion:     address.Region,
		AddressPropertyPostalCode: address.PostalCode,
		AddressPropertyCountry:    address.Country,
		AddressPropertyLongitude:  address.Longitude,
		AddressPropertyLatitude:   address.Latitude,
		AddressPropertyTimezone:   address.Timezone,
		AddressPropertyPlusCode:   address.PlusCode,
	}
}

/******************************************
 * Unmarshalling Methods
 ******************************************/

// UnmarshalMap populates this address with the properties in the `value` map
func (address *Address) UnmarshalMap(value mapof.Any) error {

	address.Name = value.GetString(AddressPropertyName)
	address.Formatted = value.GetString(AddressPropertyFormatted)
	address.Street1 = value.GetString(AddressPropertyStreet1)
	address.Street2 = value.GetString(AddressPropertyStreet2)
	address.Locality = value.GetString(AddressPropertyLocality)
	address.Region = value.GetString(AddressPropertyRegion)
	address.PostalCode = value.GetString(AddressPropertyPostalCode)
	address.Country = value.GetString(AddressPropertyCountry)
	address.Timezone = value.GetString(AddressPropertyTimezone)
	address.PlusCode = value.GetString(AddressPropertyPlusCode)
	address.Longitude = value.GetFloat(AddressPropertyLongitude)
	address.Latitude = value.GetFloat(AddressPropertyLongitude)

	return nil
}

/******************************************
 * Conversion Functions
 ******************************************/

// GeoJSON returns a GeoJSON object that matches the
// geo.GeoJSONer interface
// https://www.mongodb.com/docs/manual/reference/geojson/
func (address Address) GeoJSON() mapof.Any {
	return mapof.Any{
		PropertyType:        PropertyTypePoint,
		PropertyCoordinates: []float64{address.Longitude, address.Latitude},
	}
}

// GeoPoint returns a Point representation of this address
func (address Address) GeoPoint() Point {
	return NewPoint(address.Longitude, address.Latitude)
}

// JSONLD returns a JSON-LD representation of this object
func (address Address) JSONLD() mapof.Any {

	result := mapof.Any{
		"type": "Place",
	}

	if address.Name != "" {
		result["name"] = address.Name
	}

	if address.Latitude != 0 && address.Longitude != 0 {
		result["longitude"] = address.Longitude
		result["latitude"] = address.Latitude
	}

	if address := address.parsedValues(); address.NotEmpty() {
		result["address"] = address
	}

	return result
}

func (address Address) parsedValues() mapof.String {

	result := mapof.String{}

	if address.Street1 != "" {
		result["street1"] = address.Street1
	}

	if address.Street2 != "" {
		result["street2"] = address.Street2
	}

	if address.Locality != "" {
		result["locality"] = address.Locality
	}

	if address.Region != "" {
		result["region"] = address.Region
	}

	if address.PostalCode != "" {
		result["postalCode"] = address.PostalCode
	}

	if address.Country != "" {
		result["country"] = address.Country
	}

	return result
}
