package geo

import (
	"strconv"

	"github.com/benpate/rosetta/schema"
)

// AddressSchema returns the rosetta schema that describes an Address.
func AddressSchema() schema.Element {
	return schema.Object{
		Properties: schema.ElementMap{
			AddressPropertyName:       schema.String{},
			AddressPropertyFormatted:  schema.String{},
			AddressPropertyStreet1:    schema.String{},
			AddressPropertyStreet2:    schema.String{},
			AddressPropertyLocality:   schema.String{},
			AddressPropertyRegion:     schema.String{},
			AddressPropertyPostalCode: schema.String{},
			AddressPropertyCountry:    schema.String{},
			AddressPropertyLongitude:  schema.Number{BitSize: 64},
			AddressPropertyLatitude:   schema.Number{BitSize: 64},
		},
	}
}

// GetString returns the named property as a string, or "" if it is not a string property.
func (address Address) GetString(name string) string {
	result, _ := address.GetStringOK(name)
	return result
}

// GetStringOK returns the named property as a string, and a boolean that is TRUE
// when the property name is recognized.
func (address Address) GetStringOK(name string) (string, bool) {

	switch name {

	case "name":
		return address.Name, true

	case "formatted":
		return address.Formatted, true

	case "street1":
		return address.Street1, true

	case "street2":
		return address.Street2, true

	case "locality":
		return address.Locality, true

	case "region":
		return address.Region, true

	case "postalCode":
		return address.PostalCode, true

	case "country":
		return address.Country, true

	// longitude/latitude can be read as strings, too...
	case "longitude":
		return strconv.FormatFloat(address.Longitude, 'f', 10, 64), true

	case "latitude":
		return strconv.FormatFloat(address.Latitude, 'f', 10, 64), true
	}

	return "", false
}

// GetFloat returns the named property as a float64, or 0 if it is not a float property.
func (address Address) GetFloat(name string) float64 {
	result, _ := address.GetFloatOK(name)
	return result
}

// GetFloatOK returns the named property as a float64, and a boolean that is TRUE
// when the property name is recognized ("latitude" or "longitude").
func (address Address) GetFloatOK(name string) (float64, bool) {

	switch name {

	case "latitude":
		return address.Latitude, true

	case "longitude":
		return address.Longitude, true
	}

	return 0, false
}

/******************************************
 * Setters
 ******************************************/

// SetString sets the "name" or "formatted" property and returns TRUE if the
// property name is writable. Changing "formatted" resets any geocoded fields.
func (address *Address) SetString(name string, value string) bool {

	switch name {

	case "name":
		address.Name = value
		return true

	case "formatted":
		if value != address.Formatted {
			address.Formatted = value
			address.Reset()
		}
		return true
	}

	return false
}

// SetFloat sets the "latitude" or "longitude" property and returns TRUE if the
// property name is writable.
func (address *Address) SetFloat(name string, value float64) bool {

	switch name {

	case "latitude":
		address.Latitude = value
		return true

	case "longitude":
		address.Longitude = value
		return true
	}

	return false
}
