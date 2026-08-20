package geo

import (
	"encoding/binary"

	"github.com/benpate/derp"
	"go.mongodb.org/mongo-driver/bson"
)

// minimumBSONDocumentLength is the size of the smallest legal BSON document:
// a four-byte length prefix followed by the terminating null byte.
const minimumBSONDocumentLength = 5

// validateBSONDocument returns an error if the provided bytes are not a
// well-formed BSON document.
//
// It exists because the driver's own validator is not safe to call on
// arbitrary input: bsoncore.Document.Validate reads the length prefix and then
// indexes data[length-1] without checking the sign, so a document declaring a
// length of zero or less panics instead of returning an error. Both of this
// package's BSON entry points are exported and can be handed any bytes at all.
func validateBSONDocument(data []byte) error {

	const location = "geo.validateBSONDocument"

	// RULE: The prefix itself must be present before it can be read
	if len(data) < minimumBSONDocumentLength {
		return derp.BadRequest(location, "BSON document is too short to be valid", len(data))
	}

	// RULE: The declared length must be sane, or the driver's validator panics.
	// The signed conversion is the point of this check, not a hazard: BSON
	// lengths are int32, so a document declaring more than MaxInt32 must wrap to
	// a negative value and be rejected here, exactly as the driver would read it.
	// #nosec G115 -- the wraparound is deliberate and is what makes the guard correct
	if length := int32(binary.LittleEndian.Uint32(data)); length < minimumBSONDocumentLength {
		return derp.BadRequest(location, "BSON document declares an invalid length", length)
	}

	// Now that the prefix is trustworthy, the driver can check the rest
	if err := bson.Raw(data).Validate(); err != nil {
		return derp.Wrap(err, location, "Invalid BSON document")
	}

	// Bee-son. It's a BSON made of bees.
	return nil
}
