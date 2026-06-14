package geo

import (
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/require"
)

func TestNewAddress(t *testing.T) {
	require.Equal(t, Address{}, NewAddress())
	require.True(t, NewAddress().IsZero())
}

func TestAddress_Reset(t *testing.T) {

	address := Address{
		Name:       "Home",
		Formatted:  "123 Main St, Anytown",
		Street1:    "123 Main St",
		Street2:    "Apt 4",
		Locality:   "Anytown",
		Region:     "CA",
		PostalCode: "90210",
		Country:    "USA",
		PlusCode:   "849VCWC8+R9",
		Timezone:   "America/Los_Angeles",
		Latitude:   34,
		Longitude:  -118,
	}

	address.Reset()

	// Reset clears all geocoded values, but leaves Name and Formatted intact
	require.Equal(t, "Home", address.Name)
	require.Equal(t, "123 Main St, Anytown", address.Formatted)
	require.Empty(t, address.Street1)
	require.Empty(t, address.Street2)
	require.Empty(t, address.Locality)
	require.Empty(t, address.Region)
	require.Empty(t, address.PostalCode)
	require.Empty(t, address.Country)
	require.Empty(t, address.PlusCode)
	require.Empty(t, address.Timezone)
	require.Zero(t, address.Latitude)
	require.Zero(t, address.Longitude)
}

func TestAddress_IsZero(t *testing.T) {

	// IsZero only inspects Latitude, Longitude, and Formatted
	require.True(t, Address{}.IsZero())
	require.True(t, Address{Name: "name only"}.IsZero())
	require.True(t, Address{Street1: "street only"}.IsZero())

	require.False(t, Address{Latitude: 1}.IsZero())
	require.False(t, Address{Longitude: 1}.IsZero())
	require.False(t, Address{Formatted: "123 Main St"}.IsZero())
}

func TestAddress_NotZero(t *testing.T) {
	require.False(t, Address{}.NotZero())
	require.True(t, Address{Formatted: "123 Main St"}.NotZero())
}

func TestAddress_HasGeocode(t *testing.T) {
	require.False(t, Address{}.HasGeocode())
	require.False(t, Address{Formatted: "123 Main St"}.HasGeocode())
	require.True(t, Address{Latitude: 1}.HasGeocode())
	require.True(t, Address{Longitude: 1}.HasGeocode())
}

func TestAddress_HasAddress(t *testing.T) {

	require.False(t, Address{}.HasAddress())
	require.False(t, Address{Name: "name", Formatted: "formatted"}.HasAddress())

	// Any single parsed component is enough to count as "has address"
	require.True(t, Address{Country: "USA"}.HasAddress())
	require.True(t, Address{PostalCode: "90210"}.HasAddress())
	require.True(t, Address{Region: "CA"}.HasAddress())
	require.True(t, Address{Locality: "Anytown"}.HasAddress())
	require.True(t, Address{Street1: "123 Main St"}.HasAddress())
}

func TestAddress_LonLat(t *testing.T) {
	address := Address{Longitude: -118.5, Latitude: 34.25}
	require.Equal(t, "-118.5000000000,34.2500000000", address.LonLat())
}

func TestAddress_LatLon(t *testing.T) {
	address := Address{Longitude: -118.5, Latitude: 34.25}
	require.Equal(t, "34.2500000000,-118.5000000000", address.LatLon())
}

func TestAddress_SetPoint(t *testing.T) {

	address := Address{}
	address.SetPoint(NewPoint(-118, 34))

	require.Equal(t, -118.0, address.Longitude)
	require.Equal(t, 34.0, address.Latitude)
}

func TestAddress_MarshalMap(t *testing.T) {

	address := Address{
		Name:       "Home",
		Formatted:  "123 Main St",
		Street1:    "123 Main St",
		Street2:    "Apt 4",
		Locality:   "Anytown",
		Region:     "CA",
		PostalCode: "90210",
		Country:    "USA",
		PlusCode:   "849VCWC8+R9",
		Timezone:   "America/Los_Angeles",
		Latitude:   34,
		Longitude:  -118,
	}

	result := address.MarshalMap()

	require.Equal(t, "Home", result[AddressPropertyName])
	require.Equal(t, "123 Main St", result[AddressPropertyFormatted])
	require.Equal(t, "123 Main St", result[AddressPropertyStreet1])
	require.Equal(t, "Apt 4", result[AddressPropertyStreet2])
	require.Equal(t, "Anytown", result[AddressPropertyLocality])
	require.Equal(t, "CA", result[AddressPropertyRegion])
	require.Equal(t, "90210", result[AddressPropertyPostalCode])
	require.Equal(t, "USA", result[AddressPropertyCountry])
	require.Equal(t, "849VCWC8+R9", result[AddressPropertyPlusCode])
	require.Equal(t, "America/Los_Angeles", result[AddressPropertyTimezone])
	require.Equal(t, 34.0, result[AddressPropertyLatitude])
	require.Equal(t, -118.0, result[AddressPropertyLongitude])
}

func TestAddress_UnmarshalMap(t *testing.T) {

	value := mapof.Any{
		AddressPropertyName:       "Home",
		AddressPropertyFormatted:  "123 Main St",
		AddressPropertyStreet1:    "123 Main St",
		AddressPropertyStreet2:    "Apt 4",
		AddressPropertyLocality:   "Anytown",
		AddressPropertyRegion:     "CA",
		AddressPropertyPostalCode: "90210",
		AddressPropertyCountry:    "USA",
		AddressPropertyPlusCode:   "849VCWC8+R9",
		AddressPropertyTimezone:   "America/Los_Angeles",
		AddressPropertyLongitude:  -118.0,
		AddressPropertyLatitude:   34.0,
	}

	address := Address{}
	err := address.UnmarshalMap(value)

	require.Nil(t, err)
	require.Equal(t, "Home", address.Name)
	require.Equal(t, "123 Main St", address.Formatted)
	require.Equal(t, "123 Main St", address.Street1)
	require.Equal(t, "Apt 4", address.Street2)
	require.Equal(t, "Anytown", address.Locality)
	require.Equal(t, "CA", address.Region)
	require.Equal(t, "90210", address.PostalCode)
	require.Equal(t, "USA", address.Country)
	require.Equal(t, "849VCWC8+R9", address.PlusCode)
	require.Equal(t, "America/Los_Angeles", address.Timezone)
	require.Equal(t, -118.0, address.Longitude)
	require.Equal(t, 34.0, address.Latitude)
}

func TestAddress_GeoJSON(t *testing.T) {

	address := Address{Longitude: -118, Latitude: 34}
	result := address.GeoJSON()

	require.Equal(t, PropertyTypePoint, result[PropertyType])
	require.Equal(t, []float64{-118, 34}, result[PropertyCoordinates])
}

func TestAddress_GeoPoint(t *testing.T) {

	address := Address{Longitude: -118, Latitude: 34}
	point := address.GeoPoint()

	require.Equal(t, NewPoint(-118, 34), point)
}

func TestAddress_JSONLD_Empty(t *testing.T) {

	result := Address{}.JSONLD()

	// An empty address still carries its "Place" type, but nothing else
	require.Equal(t, "Place", result["type"])
	require.NotContains(t, result, "name")
	require.NotContains(t, result, "latitude")
	require.NotContains(t, result, "longitude")
	require.NotContains(t, result, "address")
}

func TestAddress_JSONLD_Full(t *testing.T) {

	address := Address{
		Name:       "Home",
		Street1:    "123 Main St",
		Street2:    "Apt 4",
		Locality:   "Anytown",
		Region:     "CA",
		PostalCode: "90210",
		Country:    "USA",
		Latitude:   34,
		Longitude:  -118,
	}

	result := address.JSONLD()

	require.Equal(t, "Place", result["type"])
	require.Equal(t, "Home", result["name"])
	require.Equal(t, -118.0, result["longitude"])
	require.Equal(t, 34.0, result["latitude"])

	parsed := result["address"].(mapof.String)
	require.Equal(t, "123 Main St", parsed["street1"])
	require.Equal(t, "Apt 4", parsed["street2"])
	require.Equal(t, "Anytown", parsed["locality"])
	require.Equal(t, "CA", parsed["region"])
	require.Equal(t, "90210", parsed["postalCode"])
	require.Equal(t, "USA", parsed["country"])
}

func TestAddress_JSONLD_PartialGeocode(t *testing.T) {

	// JSONLD only emits coordinates when BOTH latitude and longitude are non-zero
	withLatOnly := Address{Latitude: 34}.JSONLD()
	require.NotContains(t, withLatOnly, "latitude")
	require.NotContains(t, withLatOnly, "longitude")

	withLonOnly := Address{Longitude: -118}.JSONLD()
	require.NotContains(t, withLonOnly, "latitude")
	require.NotContains(t, withLonOnly, "longitude")
}
