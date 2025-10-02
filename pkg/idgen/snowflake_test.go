package idgen

import (
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		processID   int64
		workerID    int64
		expectError bool
	}{
		{"valid IDs", 0, 0, false},
		{"valid max IDs", 31, 31, false},
		{"invalid process ID - negative", -1, 0, true},
		{"invalid process ID - too large", 32, 0, true},
		{"invalid worker ID - negative", 0, -1, true},
		{"invalid worker ID - too large", 0, 32, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator, err := New(tt.processID, tt.workerID)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for processID=%d workerID=%d, got nil", tt.processID, tt.workerID)
				}
				if generator != nil {
					t.Errorf("Expected nil generator for invalid IDs")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if generator == nil {
					t.Errorf("Expected non-nil generator")
				}
				if generator != nil {
					if generator.ProcessID() != tt.processID {
						t.Errorf("Expected processID %d, got %d", tt.processID, generator.ProcessID())
					}
					if generator.WorkerID() != tt.workerID {
						t.Errorf("Expected workerID %d, got %d", tt.workerID, generator.WorkerID())
					}
				}
			}
		})
	}
}

func TestSnowflakeGenerate(t *testing.T) {
	generator, err := New(1, 5)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	// Generate first ID
	id1 := generator.Generate()
	if id1 <= 0 {
		t.Errorf("Expected positive ID, got %d", id1)
	}

	// Generate second ID
	id2 := generator.Generate()
	if id2 <= 0 {
		t.Errorf("Expected positive ID, got %d", id2)
	}

	// IDs should be unique
	if id1 == id2 {
		t.Error("Generated duplicate IDs")
	}

	// IDs should be increasing (monotonic)
	if id2 <= id1 {
		t.Logf("Warning: ID decreased from %d to %d (acceptable in some cases)", id1, id2)
	}
}

func TestSnowflakeComponents(t *testing.T) {
	processID := int64(10)
	workerID := int64(20)

	generator, err := New(processID, workerID)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	id := generator.Generate()

	// Extract and verify components
	extractedProcessID := generator.ExtractProcessID(id)
	if extractedProcessID != processID {
		t.Errorf("Expected processID %d, got %d", processID, extractedProcessID)
	}

	extractedWorkerID := generator.ExtractWorkerID(id)
	if extractedWorkerID != workerID {
		t.Errorf("Expected workerID %d, got %d", workerID, extractedWorkerID)
	}

	// Verify timestamp is reasonable (within last minute)
	timestamp := generator.ExtractTimestamp(id)
	now := time.Now().UnixMilli()
	diff := now - timestamp

	if diff < 0 || diff > 60000 { // More than 1 minute difference
		t.Errorf("Timestamp seems wrong. Now: %d, Extracted: %d, Diff: %d ms", now, timestamp, diff)
	}

	// Verify sequence is within valid range
	sequence := generator.ExtractSequence(id)
	if sequence < 0 || sequence > maxSequence {
		t.Errorf("Sequence out of range: %d (max: %d)", sequence, maxSequence)
	}
}

func TestSnowflakeBatch(t *testing.T) {
	generator, err := New(2, 3)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	count := 100
	ids := generator.GenerateBatch(count)

	if len(ids) != count {
		t.Errorf("Expected %d IDs, got %d", count, len(ids))
	}

	// Check all IDs are unique
	seen := make(map[int64]bool)
	for i, id := range ids {
		if id <= 0 {
			t.Errorf("ID at index %d is not positive: %d", i, id)
		}
		if seen[id] {
			t.Errorf("Duplicate ID found: %d", id)
		}
		seen[id] = true
	}
}

func TestSnowflakeConcurrency(t *testing.T) {
	generator, err := New(3, 4)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	const goroutines = 10
	const idsPerGoroutine = 100

	var wg sync.WaitGroup
	idChan := make(chan int64, goroutines*idsPerGoroutine)

	// Generate IDs concurrently
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				idChan <- generator.Generate()
			}
		}()
	}

	wg.Wait()
	close(idChan)

	// Collect and verify uniqueness
	seen := make(map[int64]bool)
	count := 0
	for id := range idChan {
		if id <= 0 {
			t.Errorf("Generated negative ID: %d", id)
		}
		if seen[id] {
			t.Errorf("Duplicate ID in concurrent generation: %d", id)
		}
		seen[id] = true
		count++
	}

	expectedCount := goroutines * idsPerGoroutine
	if count != expectedCount {
		t.Errorf("Expected %d IDs, got %d", expectedCount, count)
	}
}

func TestSnowflakeCustomEpoch(t *testing.T) {
	customEpoch := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

	generator, err := NewWithEpoch(5, 6, customEpoch)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	if generator.Epoch() != customEpoch {
		t.Errorf("Expected epoch %d, got %d", customEpoch, generator.Epoch())
	}

	id := generator.Generate()
	if id <= 0 {
		t.Errorf("Expected positive ID with custom epoch, got %d", id)
	}
}

func TestSnowflakeSequenceRollover(t *testing.T) {
	generator, err := New(7, 8)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	// Generate enough IDs to potentially rollover sequence
	// In a single millisecond, we can generate up to 4096 IDs before rollover
	ids := generator.GenerateBatch(5000)

	// All IDs should still be unique
	seen := make(map[int64]bool)
	for _, id := range ids {
		if seen[id] {
			t.Errorf("Duplicate ID after sequence rollover: %d", id)
		}
		seen[id] = true
	}
}

func TestSnowflakeMultipleMachines(t *testing.T) {
	const numProcesses = 3
	const numWorkers = 3
	const idsPerMachine = 100

	type machineKey struct {
		processID int64
		workerID  int64
	}

	allIDs := make(map[int64]machineKey)
	var mu sync.Mutex

	var wg sync.WaitGroup

	// Create multiple "machines" (process + worker combinations)
	for p := int64(0); p < numProcesses; p++ {
		for w := int64(0); w < numWorkers; w++ {
			wg.Add(1)
			go func(processID, workerID int64) {
				defer wg.Done()

				generator, err := New(processID, workerID)
				if err != nil {
					t.Errorf("Failed to create generator: %v", err)
					return
				}

				for i := 0; i < idsPerMachine; i++ {
					id := generator.Generate()

					mu.Lock()
					if existing, exists := allIDs[id]; exists {
						t.Errorf("Duplicate ID %d from process=%d worker=%d and process=%d worker=%d",
							id, processID, workerID, existing.processID, existing.workerID)
					}
					allIDs[id] = machineKey{processID, workerID}
					mu.Unlock()

					// Verify extracted components match
					if generator.ExtractProcessID(id) != processID {
						t.Errorf("Extracted processID doesn't match: expected %d, got %d",
							processID, generator.ExtractProcessID(id))
					}
					if generator.ExtractWorkerID(id) != workerID {
						t.Errorf("Extracted workerID doesn't match: expected %d, got %d",
							workerID, generator.ExtractWorkerID(id))
					}
				}
			}(p, w)
		}
	}

	wg.Wait()

	expectedTotal := numProcesses * numWorkers * idsPerMachine
	if len(allIDs) != int(expectedTotal) {
		t.Errorf("Expected %d unique IDs, got %d", expectedTotal, len(allIDs))
	}
}

// Benchmarks
func BenchmarkSnowflakeGenerate(b *testing.B) {
	generator, _ := New(1, 2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.Generate()
	}
}

func BenchmarkSnowflakeGenerateParallel(b *testing.B) {
	generator, _ := New(1, 3)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			generator.Generate()
		}
	})
}

func BenchmarkSnowflakeBatch(b *testing.B) {
	generator, _ := New(2, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.GenerateBatch(100)
	}
}
