package idgen

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// UUID represents a universally unique identifier
type UUID [16]byte

// String returns the string representation of the UUID in the format:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func (u UUID) String() string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:36], u[10:16])

	return string(buf)
}

// NewUUIDv4 generates a new UUID v4 (random)
// UUID v4 is a randomly generated UUID with 122 bits of randomness
// Format: xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx
// where x is any hexadecimal digit and y is one of 8, 9, a, or b
func NewUUIDv4() (UUID, error) {
	var uuid UUID

	// Generate 16 random bytes
	_, err := rand.Read(uuid[:])
	if err != nil {
		return uuid, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Set version (4) in the version field (bits 4-7 of byte 6)
	// Clear bits 4-7 and set to 0100 (4)
	uuid[6] = (uuid[6] & 0x0f) | 0x40

	// Set variant (RFC 4122) in the variant field (bits 6-7 of byte 8)
	// Clear bits 6-7 and set to 10
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return uuid, nil
}

// GenerateUUIDv4 generates a new UUID v4 and returns it as a string
func GenerateUUIDv4() (string, error) {
	uuid, err := NewUUIDv4()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

// GenerateUUIDv4Batch generates multiple UUID v4s
func GenerateUUIDv4Batch(count int) ([]string, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be greater than 0")
	}

	uuids := make([]string, count)
	for i := 0; i < count; i++ {
		uuid, err := GenerateUUIDv4()
		if err != nil {
			return nil, fmt.Errorf("failed to generate UUID at index %d: %w", i, err)
		}
		uuids[i] = uuid
	}

	return uuids, nil
}

// UUIDv7Generator generates time-ordered UUID v7
type UUIDv7Generator struct {
	mu            sync.Mutex
	lastTimestamp int64
	sequence      uint16
}

var globalUUIDv7Generator = &UUIDv7Generator{}

// NewUUIDv7 generates a new UUID v7 (time-ordered)
// UUID v7 uses Unix timestamp in milliseconds for time-ordering
// Format: xxxxxxxx-xxxx-7xxx-yxxx-xxxxxxxxxxxx
// Structure:
// - 48 bits: Unix timestamp in milliseconds
// - 4 bits: Version (7)
// - 12 bits: Random data
// - 2 bits: Variant (10)
// - 62 bits: Random data
func NewUUIDv7() (UUID, error) {
	return globalUUIDv7Generator.Generate()
}

// Generate creates a new UUID v7
func (g *UUIDv7Generator) Generate() (UUID, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	var uuid UUID

	// Get current timestamp in milliseconds
	now := time.Now().UnixMilli()

	// Handle clock regression or same millisecond
	if now <= g.lastTimestamp {
		g.sequence++
		// If sequence overflows, wait for next millisecond
		if g.sequence > 0x0FFF { // 12 bits max
			time.Sleep(time.Millisecond)
			now = time.Now().UnixMilli()
			g.sequence = 0
		}
	} else {
		g.sequence = 0
	}
	g.lastTimestamp = now

	// Fill timestamp (48 bits) - bytes 0-5
	binary.BigEndian.PutUint64(uuid[0:8], uint64(now)<<16)

	// Generate random bytes for the rest
	randomBytes := make([]byte, 10)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return uuid, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Set version 7 and sequence in bytes 6-7
	// Byte 6: version (4 bits) + sequence high (4 bits)
	uuid[6] = 0x70 | byte((g.sequence>>8)&0x0F)

	// Byte 7: sequence low (8 bits)
	uuid[7] = byte(g.sequence & 0xFF)

	// Copy random bytes to uuid[8:16]
	copy(uuid[8:], randomBytes[0:8])

	// Set variant (RFC 4122) in byte 8
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return uuid, nil
}

// GenerateUUIDv7 generates a new UUID v7 and returns it as a string
func GenerateUUIDv7() (string, error) {
	uuid, err := NewUUIDv7()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

// GenerateUUIDv7Batch generates multiple UUID v7s
func GenerateUUIDv7Batch(count int) ([]string, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be greater than 0")
	}

	uuids := make([]string, count)
	for i := 0; i < count; i++ {
		uuid, err := GenerateUUIDv7()
		if err != nil {
			return nil, fmt.Errorf("failed to generate UUID v7 at index %d: %w", i, err)
		}
		uuids[i] = uuid
	}

	return uuids, nil
}

// ExtractTimestampFromUUIDv7 extracts the timestamp from a UUID v7
func ExtractTimestampFromUUIDv7(uuid UUID) int64 {
	// Extract first 48 bits (6 bytes) as timestamp
	timestamp := binary.BigEndian.Uint64(append([]byte{0, 0}, uuid[0:6]...))
	return int64(timestamp)
}

// ExtractTimeFromUUIDv7 extracts the time.Time from a UUID v7
func ExtractTimeFromUUIDv7(uuid UUID) time.Time {
	timestamp := ExtractTimestampFromUUIDv7(uuid)
	return time.UnixMilli(timestamp)
}
