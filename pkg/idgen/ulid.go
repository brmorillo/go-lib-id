package idgen

// ULID (Universally Unique Lexicographically Sortable Identifier)
//
// ULID is designed to be:
// - Lexicographically sortable
// - Case insensitive
// - URL safe (Base32 encoded)
// - Monotonically increasing within the same millisecond
// - 128-bit compatibility with UUID
//
// Structure:
// - 48 bits: Timestamp (milliseconds since Unix epoch)
// - 80 bits: Randomness
//
// Format: 01AN4Z07BY79KA1307SR9X4MV3
// Length: 26 characters (Base32 encoded)
//
// Benefits over UUID:
// - More compact (26 chars vs 36)
// - Case insensitive
// - Better database index performance (time-ordered)
// - URL-safe by default

// GenerateULID generates a new ULID
// TODO: Implement ULID generation algorithm
func GenerateULID() (string, error) {
	return "", ErrULIDNotImplemented
}

// GenerateULIDBatch generates multiple ULIDs
func GenerateULIDBatch(count int) ([]string, error) {
	if count <= 0 {
		return nil, ErrULIDNotImplemented
	}
	return nil, ErrULIDNotImplemented
}
