package idgen

import (
	"testing"
)

func TestSetDefaultMachineID(t *testing.T) {
	err := SetDefaultMachineID(5, 10)
	if err != nil {
		t.Fatalf("Failed to set default machine ID: %v", err)
	}

	gen := GetDefaultGenerator()
	if gen == nil {
		t.Fatal("Expected non-nil generator after SetDefaultMachineID")
	}

	if gen.ProcessID() != 5 {
		t.Errorf("Expected processID 5, got %d", gen.ProcessID())
	}

	if gen.WorkerID() != 10 {
		t.Errorf("Expected workerID 10, got %d", gen.WorkerID())
	}
}

func TestSetDefaultMachineIDInvalid(t *testing.T) {
	tests := []struct {
		name      string
		processID int64
		workerID  int64
	}{
		{"Invalid processID", -1, 0},
		{"Invalid workerID", 0, -1},
		{"ProcessID too large", 32, 0},
		{"WorkerID too large", 0, 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetDefaultMachineID(tt.processID, tt.workerID)
			if err == nil {
				t.Errorf("Expected error for processID=%d workerID=%d", tt.processID, tt.workerID)
			}
		})
	}
}

func TestGenerateSnowflake(t *testing.T) {
	err := SetDefaultMachineID(3, 7)
	if err != nil {
		t.Fatalf("Failed to set default machine ID: %v", err)
	}

	id := GenerateSnowflake()
	if id <= 0 {
		t.Errorf("Expected positive ID, got %d", id)
	}

	// Generate another and ensure uniqueness
	id2 := GenerateSnowflake()
	if id2 <= 0 {
		t.Errorf("Expected positive ID, got %d", id2)
	}

	if id == id2 {
		t.Error("Generated duplicate IDs")
	}
}

func TestGenerateSnowflakeWithoutSetup(t *testing.T) {
	// Reset generator
	SetDefaultGenerator(nil)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when generating without setup")
		}
	}()

	GenerateSnowflake()
}

func TestGenerateSnowflakeBatchGlobal(t *testing.T) {
	err := SetDefaultMachineID(4, 8)
	if err != nil {
		t.Fatalf("Failed to set default machine ID: %v", err)
	}

	count := 50
	ids := GenerateSnowflakeBatch(count)

	if len(ids) != count {
		t.Errorf("Expected %d IDs, got %d", count, len(ids))
	}

	// Check uniqueness
	seen := make(map[int64]bool)
	for _, id := range ids {
		if id <= 0 {
			t.Errorf("Generated non-positive ID: %d", id)
		}
		if seen[id] {
			t.Errorf("Duplicate ID in batch: %d", id)
		}
		seen[id] = true
	}
}

func TestSetDefaultGenerator(t *testing.T) {
	customGen, err := New(15, 20)
	if err != nil {
		t.Fatalf("Failed to create custom generator: %v", err)
	}

	SetDefaultGenerator(customGen)

	retrievedGen := GetDefaultGenerator()
	if retrievedGen != customGen {
		t.Error("Retrieved generator is not the same as set generator")
	}

	if retrievedGen.ProcessID() != 15 {
		t.Errorf("Expected processID 15, got %d", retrievedGen.ProcessID())
	}

	if retrievedGen.WorkerID() != 20 {
		t.Errorf("Expected workerID 20, got %d", retrievedGen.WorkerID())
	}
}

func TestGetDefaultGenerator(t *testing.T) {
	// Set a generator
	err := SetDefaultMachineID(9, 11)
	if err != nil {
		t.Fatalf("Failed to set default machine ID: %v", err)
	}

	gen := GetDefaultGenerator()
	if gen == nil {
		t.Fatal("Expected non-nil generator")
	}

	// Verify it works
	id := gen.Generate()
	if id <= 0 {
		t.Errorf("Expected positive ID from retrieved generator, got %d", id)
	}
}

// Benchmarks
func BenchmarkGenerateSnowflakeGlobal(b *testing.B) {
	_ = SetDefaultMachineID(1, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateSnowflake()
	}
}

func BenchmarkGenerateSnowflakeGlobalParallel(b *testing.B) {
	_ = SetDefaultMachineID(2, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GenerateSnowflake()
		}
	})
}
