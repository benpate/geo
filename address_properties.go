package geo

// Property names used to read, write, and marshal the fields of an Address.
const (
	// AddressPropertyName is the human-readable name of the address.
	AddressPropertyName = "name"

	// AddressPropertyFormatted is the full, unparsed value of the address.
	AddressPropertyFormatted = "formatted"

	// AddressPropertyStreet1 is the first line of the street address.
	AddressPropertyStreet1 = "street1"

	// AddressPropertyStreet2 is the second line of the street address.
	AddressPropertyStreet2 = "street2"

	// AddressPropertyLocality is the city or town of the address.
	AddressPropertyLocality = "locality"

	// AddressPropertyRegion is the state or province of the address.
	AddressPropertyRegion = "region"

	// AddressPropertyPostalCode is the postal code of the address.
	AddressPropertyPostalCode = "postalCode"

	// AddressPropertyCountry is the country of the address.
	AddressPropertyCountry = "country"

	// AddressPropertyLongitude is the longitude of the address.
	AddressPropertyLongitude = "longitude"

	// AddressPropertyLatitude is the latitude of the address.
	AddressPropertyLatitude = "latitude"

	// AddressPropertyTimezone is the IANA time zone name of the address.
	AddressPropertyTimezone = "timezone"

	// AddressPropertyPlusCode is the Google Plus Code of the address.
	AddressPropertyPlusCode = "plusCode"
)
