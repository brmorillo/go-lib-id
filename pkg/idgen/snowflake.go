package idgen

import (
	"errors"
	"sync"
	"time"
)

// Snowflake ID structure (64 bits):
// Based on Discord/Twitter Snowflake implementation
//
// ┌─────────┬──────────────┬────────────┬───────────┬──────────┐
// │ Sign    │  Timestamp   │ Process ID │ Worker ID │ Sequence │
// │ 1 bit   │   41 bits    │   5 bits   │  5 bits   │ 12 bits  │
// │ (unused)│              │  (0-31)    │  (0-31)   │ (0-4095) │
// └─────────┴──────────────┴────────────┴───────────┴──────────┘
//
// Sign bit (1 bit): Always 0 to keep the ID positive (int64 compatibility)
// Timestamp (41 bits): Milliseconds since custom epoch (~69 years capacity)
// Process ID (5 bits): Unique process identifier (0-31 processes)
// Worker ID (5 bits): Unique worker/thread identifier per process (0-31 workers)
// Sequence (12 bits): Incremental sequence per millisecond (0-4095 IDs/ms)
//
// Total capacity: 32 processes × 32 workers × 4096 IDs/ms = ~4.1M IDs per millisecond

const (
	// Epoch is the custom epoch (January 1, 2025 00:00:00 UTC)
	// IDs generated before this date are not supported
	// IMPORTANT: Epoch must be in the PAST relative to current time
	// to avoid negative timestamp deltas (current_time - epoch must be positive)

	/*
		Unix Seconds
			1735689600
		Unix Milliseconds
			1735689600000
		Etc/GMT
			January 1st 2025, 12:00:00 am GMT+00:00
		UTC ISO 8601
			2025-01-01T00:00:00.000Z
		UTC RFC 2822
			Wed, 01 Jan 2025 00:00:00 GMT
	*/
	DefaultEpoch int64 = 1735689600000

	// Bit lengths for Snowflake ID components (Discord/Twitter standard)
	// timestampBits: 41 bits = ~69 years from epoch (2^41 ms = ~69.73 years)
	// processIDBits: 5 bits = 32 unique processes (0-31)
	// workerIDBits: 5 bits = 32 unique workers per process (0-31)
	// sequenceBits: 12 bits = 4096 IDs per millisecond per worker (0-4095)
	timestampBits = 41
	processIDBits = 5
	workerIDBits  = 5
	sequenceBits  = 12

	// Max values calculated from bit lengths
	maxProcessID = -1 ^ (-1 << processIDBits) // 31
	maxWorkerID  = -1 ^ (-1 << workerIDBits)  // 31
	maxSequence  = -1 ^ (-1 << sequenceBits)  // 4095

	// Bit shifts for ID construction
	sequenceShift  = 0
	workerIDShift  = sequenceBits                                // 12
	processIDShift = sequenceBits + workerIDBits                 // 17
	timestampShift = sequenceBits + workerIDBits + processIDBits // 22
)

var (
	// ErrInvalidProcessID is returned when process ID is out of range
	ErrInvalidProcessID = errors.New("process ID must be between 0 and 31")

	// ErrInvalidWorkerID is returned when worker ID is out of range
	ErrInvalidWorkerID = errors.New("worker ID must be between 0 and 31")

	// ErrClockMovedBackwards is returned when system clock moves backwards
	ErrClockMovedBackwards = errors.New("clock moved backwards")
)

// Snowflake generates unique 64-bit IDs in a distributed system
type Snowflake struct {
	mu            sync.Mutex
	epoch         int64
	processID     int64
	workerID      int64
	sequence      int64
	lastTimestamp int64
}

// New creates a new Snowflake ID generator.
// This is the main constructor for creating ID generators.
//
// Parameters:
//   - processID: Unique process identifier (0-31). Should be unique per server/process
//   - workerID: Unique worker identifier within the process (0-31). Should be unique per thread/worker
//
// Returns:
//   - *Snowflake: A new ID generator instance
//   - error: ErrInvalidProcessID or ErrInvalidWorkerID if parameters are out of range
//
// Example:
//
//	generator, err := idgen.New(5, 12)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	id := generator.Generate()
func New(processID, workerID int64) (*Snowflake, error) {
	return NewWithEpoch(processID, workerID, DefaultEpoch)
}

// NewWithEpoch creates a new Snowflake ID generator with a custom epoch.
// This is useful for testing or when you need a different epoch than the default.
//
// Parameters:
//   - processID: Unique process identifier (0-31)
//   - workerID: Unique worker identifier within the process (0-31)
//   - epoch: Custom epoch in milliseconds since Unix epoch
//
// Returns:
//   - *Snowflake: A new ID generator instance
//   - error: ErrInvalidProcessID or ErrInvalidWorkerID if parameters are out of range
//
// Example:
//
//	customEpoch := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
//	generator, err := idgen.NewWithEpoch(5, 12, customEpoch)
func NewWithEpoch(processID, workerID int64, epoch int64) (*Snowflake, error) {
	if processID < 0 || processID > maxProcessID {
		return nil, ErrInvalidProcessID
	}
	if workerID < 0 || workerID > maxWorkerID {
		return nil, ErrInvalidWorkerID
	}

	return &Snowflake{
		epoch:         epoch,
		processID:     processID,
		workerID:      workerID,
		sequence:      0,
		lastTimestamp: 0,
	}, nil
}

// NewSnowflake creates a new Snowflake ID generator.
// Deprecated: Use New instead for cleaner API.
func NewSnowflake(processID, workerID int64) (*Snowflake, error) {
	return New(processID, workerID)
}

// NewSnowflakeWithEpoch creates a new Snowflake ID generator with custom epoch.
// Deprecated: Use NewWithEpoch instead for cleaner API.
func NewSnowflakeWithEpoch(processID, workerID int64, epoch int64) (*Snowflake, error) {
	return NewWithEpoch(processID, workerID, epoch)
}

// Generate creates a new unique Snowflake ID.
// This method is thread-safe and can be called concurrently from multiple goroutines.
//
// The generated ID is guaranteed to be:
//   - Unique across all generators with different processID/workerID combinations
//   - Approximately sortable by creation time
//   - Positive (fits in int64 without issues)
//
// Returns:
//   - int64: A unique 64-bit Snowflake ID
//
// Example:
//
//	generator, _ := idgen.New(5, 12)
//	id := generator.Generate()
//	fmt.Printf("Generated ID: %d\n", id)
func (s *Snowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := s.currentTimestamp()

	// Clock moved backwards - wait until it catches up
	if timestamp < s.lastTimestamp {
		// In production, you might want to return an error here
		// For now, we'll wait
		time.Sleep(time.Duration(s.lastTimestamp-timestamp) * time.Millisecond)
		timestamp = s.currentTimestamp()
	}

	// Same millisecond - increment sequence
	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence

		// Sequence overflow - wait for next millisecond
		if s.sequence == 0 {
			timestamp = s.waitNextMillis(s.lastTimestamp)
		}
	} else {
		// New millisecond - reset sequence
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	// Construct the ID (Discord/Twitter Snowflake format)
	// [1 bit sign (0)] [41 bits timestamp] [5 bits processID] [5 bits workerID] [12 bits sequence]
	id := ((timestamp - s.epoch) << timestampShift) |
		(s.processID << processIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return id
}

// GenerateBatch generates multiple IDs at once for better performance.
// This method is more efficient than calling Generate() multiple times
// when you need many IDs at once.
//
// Parameters:
//   - count: Number of IDs to generate
//
// Returns:
//   - []int64: Slice of unique Snowflake IDs
//
// Example:
//
//	generator, _ := idgen.New(5, 12)
//	ids := generator.GenerateBatch(100)
//	fmt.Printf("Generated %d IDs\n", len(ids))
func (s *Snowflake) GenerateBatch(count int) []int64 {
	ids := make([]int64, count)
	for i := 0; i < count; i++ {
		ids[i] = s.Generate()
	}
	return ids
}

// ExtractTimestamp extracts the timestamp component from a Snowflake ID.
// Returns the timestamp in milliseconds since Unix epoch.
//
// Parameters:
//   - id: A Snowflake ID to extract timestamp from
//
// Returns:
//   - int64: Timestamp in milliseconds since Unix epoch
//
// Example:
//
//	generator, _ := idgen.New(5, 12)
//	id := generator.Generate()
//	timestamp := generator.ExtractTimestamp(id)
//	fmt.Printf("ID was created at: %d ms\n", timestamp)
func (s *Snowflake) ExtractTimestamp(id int64) int64 {
	return (id >> timestampShift) + s.epoch
}

// ExtractProcessID extracts the process ID component from a Snowflake ID.
//
// Parameters:
//   - id: A Snowflake ID to extract process ID from
//
// Returns:
//   - int64: Process ID (0-31)
func (s *Snowflake) ExtractProcessID(id int64) int64 {
	return (id >> processIDShift) & maxProcessID
}

// ExtractWorkerID extracts the worker ID component from a Snowflake ID.
//
// Parameters:
//   - id: A Snowflake ID to extract worker ID from
//
// Returns:
//   - int64: Worker ID (0-31)
func (s *Snowflake) ExtractWorkerID(id int64) int64 {
	return (id >> workerIDShift) & maxWorkerID
}

// ExtractSequence extracts the sequence number from a Snowflake ID.
// The sequence represents the order of IDs generated within the same millisecond.
//
// Parameters:
//   - id: A Snowflake ID to extract sequence from
//
// Returns:
//   - int64: Sequence number (0-4095)
func (s *Snowflake) ExtractSequence(id int64) int64 {
	return id & maxSequence
}

// ExtractTime converts the Snowflake ID timestamp to a time.Time object.
// This is a convenience method that combines ExtractTimestamp with time conversion.
//
// Parameters:
//   - id: A Snowflake ID to extract time from
//
// Returns:
//   - time.Time: The creation time of the ID in UTC
//
// Example:
//
//	generator, _ := idgen.New(5, 12)
//	id := generator.Generate()
//	createdAt := generator.ExtractTime(id)
//	fmt.Printf("ID created at: %s\n", createdAt.Format(time.RFC3339))
func (s *Snowflake) ExtractTime(id int64) time.Time {
	timestamp := s.ExtractTimestamp(id)
	return time.Unix(timestamp/1000, (timestamp%1000)*1000000).UTC()
}

// currentTimestamp returns current timestamp in milliseconds
func (s *Snowflake) currentTimestamp() int64 {
	return time.Now().UnixMilli()
}

// waitNextMillis waits until next millisecond
func (s *Snowflake) waitNextMillis(lastTimestamp int64) int64 {
	timestamp := s.currentTimestamp()
	for timestamp <= lastTimestamp {
		time.Sleep(100 * time.Microsecond)
		timestamp = s.currentTimestamp()
	}
	return timestamp
}

// ProcessID returns the process ID configured for this generator.
//
// Returns:
//   - int64: Process ID (0-31)
func (s *Snowflake) ProcessID() int64 {
	return s.processID
}

// WorkerID returns the worker ID configured for this generator.
//
// Returns:
//   - int64: Worker ID (0-31)
func (s *Snowflake) WorkerID() int64 {
	return s.workerID
}

// Epoch returns the epoch configured for this generator.
//
// Returns:
//   - int64: Epoch timestamp in milliseconds since Unix epoch
func (s *Snowflake) Epoch() int64 {
	return s.epoch
}
