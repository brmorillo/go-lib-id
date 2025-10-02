package idgen

import (
	"sync"
)

var (
	// defaultGenerator is the global Snowflake generator
	defaultGenerator *Snowflake
	defaultMu        sync.RWMutex
)

// SetDefaultMachineID configures the global Snowflake generator with processID and workerID
// This should be called once at application startup
// processID: unique identifier for the process (0-31)
// workerID: unique identifier for the worker thread (0-31)
func SetDefaultMachineID(processID, workerID int64) error {
	generator, err := NewSnowflake(processID, workerID)
	if err != nil {
		return err
	}

	defaultMu.Lock()
	defaultGenerator = generator
	defaultMu.Unlock()

	return nil
}

// SetDefaultGenerator sets a custom Snowflake generator as the global default
func SetDefaultGenerator(generator *Snowflake) {
	defaultMu.Lock()
	defaultGenerator = generator
	defaultMu.Unlock()
}

// GenerateSnowflake generates a Snowflake ID using the global generator
// Panics if SetDefaultMachineID was not called first
func GenerateSnowflake() int64 {
	defaultMu.RLock()
	gen := defaultGenerator
	defaultMu.RUnlock()

	if gen == nil {
		panic("idgen: default generator not initialized. Call SetDefaultMachineID(processID, workerID) first")
	}

	return gen.Generate()
}

// GenerateSnowflakeBatch generates multiple Snowflake IDs using the global generator
func GenerateSnowflakeBatch(count int) []int64 {
	defaultMu.RLock()
	gen := defaultGenerator
	defaultMu.RUnlock()

	if gen == nil {
		panic("idgen: default generator not initialized. Call SetDefaultMachineID(processID, workerID) first")
	}

	return gen.GenerateBatch(count)
}

// GetDefaultGenerator returns the global Snowflake generator
func GetDefaultGenerator() *Snowflake {
	defaultMu.RLock()
	defer defaultMu.RUnlock()
	return defaultGenerator
}
