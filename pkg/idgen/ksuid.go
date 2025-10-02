package idgen

// KSUID (K-Sortable Unique Identifier)
//
// KSUID is designed to be:
// - K-sortable (lexicographically sortable)
// - Distributed system friendly
// - 160-bit (20 bytes)
// - Base62 encoded (27 characters)
//
// Structure:
// - 32 bits: Timestamp (seconds since custom epoch)
// - 128 bits: Random payload
//
// Format: 0ujtsYcgvSTl8PAuAdqWYSMnLOv
// Length: 27 characters
//
// Epoch: 2014-05-13T16:53:20Z (custom epoch for better resolution)
//
// Benefits:
// - Naturally ordered by creation time
// - No coordination required
// - 134 years of usable life (from epoch)
// - URL-safe

const (
	// KSUIDEpoch is the KSUID epoch (2014-05-13T16:53:20Z)
	KSUIDEpoch int64 = 1400000000
)

// GenerateKSUID generates a new KSUID
// TODO: Implement KSUID generation algorithm
func GenerateKSUID() (string, error) {
	return "", ErrKSUIDNotImplemented
}

// GenerateKSUIDBatch generates multiple KSUIDs
func GenerateKSUIDBatch(count int) ([]string, error) {
	if count <= 0 {
		return nil, ErrKSUIDNotImplemented
	}
	return nil, ErrKSUIDNotImplemented
}
