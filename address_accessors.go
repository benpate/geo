package geo

import "github.com/benpate/rosetta/schema"

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
	}

	return "", false
}

func (address Address) GetFloat(name string) (float64, bool) {

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
