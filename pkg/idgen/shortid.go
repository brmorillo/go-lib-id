package idgen

// ShortID - Short ID generator (UUID compression)
//
// ShortID is designed to be:
// - Short (22 characters vs UUID's 36)
// - URL-safe
// - Non-sequential (secure)
// - Base62 encoded
//
// Format: ppBqWA9fuP3FcvjJHQxNz3
// Length: 22 characters
//
// Encoding: Base62 (0-9, a-z, A-Z)
//
// Use cases:
// - Shortened URLs
// - Compact database keys
// - User-facing IDs
// - API resources

// GenerateShortID generates a new ShortID
// TODO: Implement ShortID generation algorithm
func GenerateShortID() (string, error) {
	return "", ErrShortIDNotImplemented
}

// GenerateShortIDBatch generates multiple ShortIDs
func GenerateShortIDBatch(count int) ([]string, error) {
	if count <= 0 {
		return nil, ErrShortIDNotImplemented
	}
	return nil, ErrShortIDNotImplemented
}
