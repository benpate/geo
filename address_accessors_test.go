package geo

import (
	"testing"

	"github.com/benpate/rosetta/schema"
	"github.com/stretchr/testify/require"
)

func TestAddressSchema(t *testing.T) {

	element := AddressSchema()

	object, ok := element.(schema.Object)
	require.True(t, ok)

	// All ten string/number properties should be present in the schema
	require.Contains(t, object.Properties, AddressPropertyName)
	require.Contains(t, object.Properties, AddressPropertyFormatted)
	require.Contains(t, object.Properties, AddressPropertyStreet1)
	require.Contains(t, object.Properties, AddressPropertyStreet2)
	require.Contains(t, object.Properties, AddressPropertyLocality)
	require.Contains(t, object.Properties, AddressPropertyRegion)
	require.Contains(t, object.Properties, AddressPropertyPostalCode)
	require.Contains(t, object.Properties, AddressPropertyCountry)
	require.Contains(t, object.Properties, AddressPropertyLongitude)
	require.Contains(t, object.Properties, AddressPropertyLatitude)

	require.IsType(t, schema.String{}, object.Properties[AddressPropertyName])
	require.Equal(t, schema.Number{BitSize: 64}, object.Properties[AddressPropertyLongitude])
	require.Equal(t, schema.Number{BitSize: 64}, object.Properties[AddressPropertyLatitude])
}

func TestAddress_GetStringOK(t *testing.T) {

	address := Address{
		Name:       "Home",
		Formatted:  "123 Main St",
		Street1:    "123 Main St",
		Street2:    "Apt 4",
		Locality:   "Anytown",
		Region:     "CA",
		PostalCode: "90210",
		Country:    "USA",
		Longitude:  -118,
		Latitude:   34,
	}

	// check confirms both the returned value and the "ok" flag for a property name
	check := func(name string, expectedValue string, expectedOK bool) {
		value, ok := address.GetStringOK(name)
		require.Equal(t, expectedValue, value, "value for %q", name)
		require.Equal(t, expectedOK, ok, "ok for %q", name)
	}

	check("name", "Home", true)
	check("formatted", "123 Main St", true)
	check("street1", "123 Main St", true)
	check("street2", "Apt 4", true)
	check("locality", "Anytown", true)
	check("region", "CA", true)
	check("postalCode", "90210", true)
	check("country", "USA", true)

	// longitude/latitude are returned as their formatted numeric value
	check("longitude", "-118.0000000000", true)
	check("latitude", "34.0000000000", true)

	// Unknown property names return empty + false
	check("unknown", "", false)
	check("", "", false)
}

func TestAddress_GetString(t *testing.T) {

	address := Address{Name: "Home", Longitude: -118}

	require.Equal(t, "Home", address.GetString("name"))
	require.Equal(t, "-118.0000000000", address.GetString("longitude"))
	require.Equal(t, "", address.GetString("unknown"))
}

func TestAddress_GetFloat(t *testing.T) {

	address := Address{Longitude: -118, Latitude: 34}

	check := func(name string, expectedValue float64, expectedOK bool) {
		value, ok := address.GetFloat(name)
		require.Equal(t, expectedValue, value, "value for %q", name)
		require.Equal(t, expectedOK, ok, "ok for %q", name)
	}

	check("latitude", 34, true)
	check("longitude", -118, true)
	check("name", 0, false)
	check("unknown", 0, false)
}

func TestAddress_SetString(t *testing.T) {

	address := Address{}

	require.True(t, address.SetString("name", "Home"))
	require.Equal(t, "Home", address.Name)

	require.True(t, address.SetString("formatted", "123 Main St"))
	require.Equal(t, "123 Main St", address.Formatted)

	// Unknown / read-only properties cannot be set
	require.False(t, address.SetString("street1", "123 Main St"))
	require.Empty(t, address.Street1)
	require.False(t, address.SetString("unknown", "value"))
}

func TestAddress_SetString_FormattedResets(t *testing.T) {

	address := Address{
		Formatted: "123 Main St",
		Locality:  "Anytown",
		Latitude:  34,
		Longitude: -118,
	}

	// Setting "formatted" to a DIFFERENT value clears the geocoded fields
	require.True(t, address.SetString("formatted", "456 Elm St"))
	require.Equal(t, "456 Elm St", address.Formatted)
	require.Empty(t, address.Locality)
	require.Zero(t, address.Latitude)
	require.Zero(t, address.Longitude)
}

func TestAddress_SetString_FormattedUnchanged(t *testing.T) {

	address := Address{
		Formatted: "123 Main St",
		Locality:  "Anytown",
		Latitude:  34,
	}

	// Setting "formatted" to the SAME value leaves the geocoded fields intact
	require.True(t, address.SetString("formatted", "123 Main St"))
	require.Equal(t, "Anytown", address.Locality)
	require.Equal(t, 34.0, address.Latitude)
}

func TestAddress_SetFloat(t *testing.T) {

	address := Address{}

	require.True(t, address.SetFloat("latitude", 34))
	require.Equal(t, 34.0, address.Latitude)

	require.True(t, address.SetFloat("longitude", -118))
	require.Equal(t, -118.0, address.Longitude)

	require.False(t, address.SetFloat("altitude", 100))
	require.False(t, address.SetFloat("unknown", 1))
}
