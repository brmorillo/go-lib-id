package idgen

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestNewUUIDv4(t *testing.T) {
	uuid, err := NewUUIDv4()
	if err != nil {
		t.Fatalf("NewUUIDv4() error = %v", err)
	}

	// Check that UUID is not all zeros
	allZeros := true
	for _, b := range uuid {
		if b != 0 {
			allZeros = false
			break
		}
	}
	if allZeros {
		t.Error("NewUUIDv4() returned all zeros")
	}

	// Check version bits (byte 6, bits 4-7 should be 0100)
	version := (uuid[6] >> 4) & 0x0f
	if version != 4 {
		t.Errorf("NewUUIDv4() version = %d, want 4", version)
	}

	// Check variant bits (byte 8, bits 6-7 should be 10)
	variant := (uuid[8] >> 6) & 0x03
	if variant != 2 {
		t.Errorf("NewUUIDv4() variant = %d, want 2 (RFC 4122)", variant)
	}
}

func TestUUIDString(t *testing.T) {
	uuid, err := NewUUIDv4()
	if err != nil {
		t.Fatalf("NewUUIDv4() error = %v", err)
	}

	str := uuid.String()

	// Check format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidPattern.MatchString(str) {
		t.Errorf("UUID.String() = %q, does not match UUID format", str)
	}

	// Check length
	if len(str) != 36 {
		t.Errorf("UUID.String() length = %d, want 36", len(str))
	}

	// Check version digit (13th character should be '4')
	if str[14] != '4' {
		t.Errorf("UUID.String() version digit = %c, want '4'", str[14])
	}

	// Check variant digit (17th character should be 8, 9, a, or b)
	variantChar := str[19]
	if variantChar != '8' && variantChar != '9' && variantChar != 'a' && variantChar != 'b' {
		t.Errorf("UUID.String() variant digit = %c, want one of [8, 9, a, b]", variantChar)
	}
}

func TestGenerateUUIDv4(t *testing.T) {
	uuid, err := GenerateUUIDv4()
	if err != nil {
		t.Fatalf("GenerateUUIDv4() error = %v", err)
	}

	// Check format
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	if !uuidPattern.MatchString(uuid) {
		t.Errorf("GenerateUUIDv4() = %q, does not match UUID v4 format", uuid)
	}
}

func TestGenerateUUIDv4Uniqueness(t *testing.T) {
	count := 1000
	uuids := make(map[string]bool, count)

	for i := 0; i < count; i++ {
		uuid, err := GenerateUUIDv4()
		if err != nil {
			t.Fatalf("GenerateUUIDv4() error = %v at iteration %d", err, i)
		}

		if uuids[uuid] {
			t.Errorf("GenerateUUIDv4() generated duplicate UUID: %s", uuid)
		}
		uuids[uuid] = true
	}

	if len(uuids) != count {
		t.Errorf("GenerateUUIDv4() unique count = %d, want %d", len(uuids), count)
	}
}

func TestGenerateUUIDv4Batch(t *testing.T) {
	tests := []struct {
		name    string
		count   int
		wantErr bool
	}{
		{
			name:    "generate 1 UUID",
			count:   1,
			wantErr: false,
		},
		{
			name:    "generate 10 UUIDs",
			count:   10,
			wantErr: false,
		},
		{
			name:    "generate 100 UUIDs",
			count:   100,
			wantErr: false,
		},
		{
			name:    "invalid count 0",
			count:   0,
			wantErr: true,
		},
		{
			name:    "invalid count negative",
			count:   -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuids, err := GenerateUUIDv4Batch(tt.count)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateUUIDv4Batch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if len(uuids) != tt.count {
				t.Errorf("GenerateUUIDv4Batch() count = %d, want %d", len(uuids), tt.count)
			}

			// Check uniqueness
			seen := make(map[string]bool, tt.count)
			for i, uuid := range uuids {
				if seen[uuid] {
					t.Errorf("GenerateUUIDv4Batch() duplicate UUID at index %d: %s", i, uuid)
				}
				seen[uuid] = true

				// Check format
				uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
				if !uuidPattern.MatchString(uuid) {
					t.Errorf("GenerateUUIDv4Batch() UUID at index %d = %q, invalid format", i, uuid)
				}
			}
		})
	}
}

func TestUUIDv4Format(t *testing.T) {
	// Generate multiple UUIDs to test consistency
	for i := 0; i < 100; i++ {
		uuid, err := GenerateUUIDv4()
		if err != nil {
			t.Fatalf("GenerateUUIDv4() error = %v at iteration %d", err, i)
		}

		// Split by hyphen
		parts := strings.Split(uuid, "-")
		if len(parts) != 5 {
			t.Errorf("UUID has %d parts, want 5: %s", len(parts), uuid)
			continue
		}

		// Check each part length
		expectedLengths := []int{8, 4, 4, 4, 12}
		for j, part := range parts {
			if len(part) != expectedLengths[j] {
				t.Errorf("UUID part %d length = %d, want %d: %s", j, len(part), expectedLengths[j], uuid)
			}
		}

		// Verify version and variant
		if parts[2][0] != '4' {
			t.Errorf("UUID version = %c, want '4': %s", parts[2][0], uuid)
		}

		variantChar := parts[3][0]
		if variantChar != '8' && variantChar != '9' && variantChar != 'a' && variantChar != 'b' {
			t.Errorf("UUID variant = %c, want one of [8, 9, a, b]: %s", variantChar, uuid)
		}
	}
}

func BenchmarkNewUUIDv4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewUUIDv4()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGenerateUUIDv4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GenerateUUIDv4()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGenerateUUIDv4Batch(b *testing.B) {
	counts := []int{10, 100, 1000}

	for _, count := range counts {
		b.Run(string(rune(count)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := GenerateUUIDv4Batch(count)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// ==================== UUID v7 Tests ====================

func TestNewUUIDv7(t *testing.T) {
	uuid, err := NewUUIDv7()
	if err != nil {
		t.Fatalf("NewUUIDv7() error = %v", err)
	}

	// Check that UUID is not all zeros
	allZeros := true
	for _, b := range uuid {
		if b != 0 {
			allZeros = false
			break
		}
	}
	if allZeros {
		t.Error("NewUUIDv7() returned all zeros")
	}

	// Check version bits (byte 6, bits 4-7 should be 0111 = 7)
	version := (uuid[6] >> 4) & 0x0f
	if version != 7 {
		t.Errorf("NewUUIDv7() version = %d, want 7", version)
	}

	// Check variant bits (byte 8, bits 6-7 should be 10)
	variant := (uuid[8] >> 6) & 0x03
	if variant != 2 {
		t.Errorf("NewUUIDv7() variant = %d, want 2 (RFC 4122)", variant)
	}
}

func TestGenerateUUIDv7(t *testing.T) {
	uuid, err := GenerateUUIDv7()
	if err != nil {
		t.Fatalf("GenerateUUIDv7() error = %v", err)
	}

	// Check format
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	if !uuidPattern.MatchString(uuid) {
		t.Errorf("GenerateUUIDv7() = %q, does not match UUID v7 format", uuid)
	}
}

func TestGenerateUUIDv7Uniqueness(t *testing.T) {
	count := 1000
	uuids := make(map[string]bool, count)

	for i := 0; i < count; i++ {
		uuid, err := GenerateUUIDv7()
		if err != nil {
			t.Fatalf("GenerateUUIDv7() error = %v at iteration %d", err, i)
		}

		if uuids[uuid] {
			t.Errorf("GenerateUUIDv7() generated duplicate UUID: %s", uuid)
		}
		uuids[uuid] = true
	}

	if len(uuids) != count {
		t.Errorf("GenerateUUIDv7() unique count = %d, want %d", len(uuids), count)
	}
}

func TestGenerateUUIDv7TimeOrdering(t *testing.T) {
	// Generate UUIDs and check they are time-ordered
	count := 100
	uuids := make([]UUID, count)

	for i := 0; i < count; i++ {
		uuid, err := NewUUIDv7()
		if err != nil {
			t.Fatalf("NewUUIDv7() error = %v at iteration %d", err, i)
		}
		uuids[i] = uuid

		// Small sleep to ensure different timestamps
		if i%10 == 0 {
			time.Sleep(time.Millisecond)
		}
	}

	// Check that timestamps are non-decreasing
	for i := 1; i < count; i++ {
		ts1 := ExtractTimestampFromUUIDv7(uuids[i-1])
		ts2 := ExtractTimestampFromUUIDv7(uuids[i])

		if ts2 < ts1 {
			t.Errorf("UUID v7 timestamps not ordered: uuid[%d] timestamp %d >= uuid[%d] timestamp %d",
				i-1, ts1, i, ts2)
		}
	}
}

func TestGenerateUUIDv7Batch(t *testing.T) {
	tests := []struct {
		name    string
		count   int
		wantErr bool
	}{
		{
			name:    "generate 1 UUID",
			count:   1,
			wantErr: false,
		},
		{
			name:    "generate 10 UUIDs",
			count:   10,
			wantErr: false,
		},
		{
			name:    "generate 100 UUIDs",
			count:   100,
			wantErr: false,
		},
		{
			name:    "invalid count 0",
			count:   0,
			wantErr: true,
		},
		{
			name:    "invalid count negative",
			count:   -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuids, err := GenerateUUIDv7Batch(tt.count)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateUUIDv7Batch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if len(uuids) != tt.count {
				t.Errorf("GenerateUUIDv7Batch() count = %d, want %d", len(uuids), tt.count)
			}

			// Check uniqueness
			seen := make(map[string]bool, tt.count)
			for i, uuid := range uuids {
				if seen[uuid] {
					t.Errorf("GenerateUUIDv7Batch() duplicate UUID at index %d: %s", i, uuid)
				}
				seen[uuid] = true

				// Check format
				uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
				if !uuidPattern.MatchString(uuid) {
					t.Errorf("GenerateUUIDv7Batch() UUID at index %d = %q, invalid format", i, uuid)
				}
			}
		})
	}
}

func TestExtractTimestampFromUUIDv7(t *testing.T) {
	before := time.Now().UnixMilli()

	uuid, err := NewUUIDv7()
	if err != nil {
		t.Fatalf("NewUUIDv7() error = %v", err)
	}

	after := time.Now().UnixMilli()

	timestamp := ExtractTimestampFromUUIDv7(uuid)

	// Timestamp should be between before and after
	if timestamp < before || timestamp > after {
		t.Errorf("ExtractTimestampFromUUIDv7() = %d, want between %d and %d", timestamp, before, after)
	}
}

func TestExtractTimeFromUUIDv7(t *testing.T) {
	before := time.Now()

	uuid, err := NewUUIDv7()
	if err != nil {
		t.Fatalf("NewUUIDv7() error = %v", err)
	}

	after := time.Now()

	extractedTime := ExtractTimeFromUUIDv7(uuid)

	// Time should be between before and after (with some tolerance)
	if extractedTime.Before(before.Add(-time.Second)) || extractedTime.After(after.Add(time.Second)) {
		t.Errorf("ExtractTimeFromUUIDv7() = %v, want between %v and %v", extractedTime, before, after)
	}
}

func TestUUIDv7Format(t *testing.T) {
	// Generate multiple UUIDs to test consistency
	for i := 0; i < 100; i++ {
		uuid, err := GenerateUUIDv7()
		if err != nil {
			t.Fatalf("GenerateUUIDv7() error = %v at iteration %d", err, i)
		}

		// Split by hyphen
		parts := strings.Split(uuid, "-")
		if len(parts) != 5 {
			t.Errorf("UUID has %d parts, want 5: %s", len(parts), uuid)
			continue
		}

		// Check each part length
		expectedLengths := []int{8, 4, 4, 4, 12}
		for j, part := range parts {
			if len(part) != expectedLengths[j] {
				t.Errorf("UUID part %d length = %d, want %d: %s", j, len(part), expectedLengths[j], uuid)
			}
		}

		// Verify version (should be '7')
		if parts[2][0] != '7' {
			t.Errorf("UUID version = %c, want '7': %s", parts[2][0], uuid)
		}

		// Verify variant
		variantChar := parts[3][0]
		if variantChar != '8' && variantChar != '9' && variantChar != 'a' && variantChar != 'b' {
			t.Errorf("UUID variant = %c, want one of [8, 9, a, b]: %s", variantChar, uuid)
		}
	}
}

func BenchmarkNewUUIDv7(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewUUIDv7()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGenerateUUIDv7(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GenerateUUIDv7()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGenerateUUIDv7Batch(b *testing.B) {
	counts := []int{10, 100, 1000}

	for _, count := range counts {
		b.Run(string(rune(count)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := GenerateUUIDv7Batch(count)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
