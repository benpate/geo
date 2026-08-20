package geo

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestValidateBSONDocument(t *testing.T) {

	valid, err := bson.Marshal(bson.M{"type": "Point"})
	require.Nil(t, err)
	require.Nil(t, validateBSONDocument(valid))
}

func TestValidateBSONDocument_TooShort(t *testing.T) {

	require.NotNil(t, validateBSONDocument(nil))
	require.NotNil(t, validateBSONDocument([]byte{}))
	require.NotNil(t, validateBSONDocument([]byte{5, 0, 0, 0}))
}

// A length prefix of zero is the input that panics the driver's own validator,
// because it indexes data[length-1] without checking the sign.
func TestValidateBSONDocument_ZeroLengthPrefix(t *testing.T) {

	require.NotNil(t, validateBSONDocument([]byte{0, 0, 0, 0, 0}))
}

// A length prefix above MaxInt32 wraps negative, and must be rejected for the
// same reason.
func TestValidateBSONDocument_NegativeLengthPrefix(t *testing.T) {

	require.NotNil(t, validateBSONDocument([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0}))
}

// A sane prefix that overruns the buffer is caught by the driver's validator,
// which is safe to reach once the prefix itself has been checked.
func TestValidateBSONDocument_Malformed(t *testing.T) {

	require.NotNil(t, validateBSONDocument([]byte{99, 0, 0, 0, 1, 2, 3}))
}

// The document that crashed FuzzPolygon_UnmarshalBSON before the guard existed.
func TestValidateBSONDocument_FuzzCrashers(t *testing.T) {

	require.NotNil(t, validateBSONDocument([]byte("\x00\x00\x00\x00")))

	crasher := []byte("1\x00\x00\x00\x02type\x00\b\x00\x00\x000000000\x00\x04coordinates0000\x00000\xa9000000")
	require.NotNil(t, validateBSONDocument(crasher))

	polygon := Polygon{}
	require.NotNil(t, polygon.UnmarshalBSON(crasher))

	point := Point{}
	require.NotNil(t, point.UnmarshalBSON(crasher))
}
