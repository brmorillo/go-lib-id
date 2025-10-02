package main

import (
	"fmt"
	"time"

	"github.com/brmorillo/go-lib-id/pkg/idgen"
)

func main() {
	fmt.Println("🔢 Capacity Demonstration - Snowflake ID")
	fmt.Println("=========================================")
	fmt.Println()

	// 1. Test how many IDs we can generate in 1 millisecond
	fmt.Println("1️⃣  Testing capacity in 1 millisecond:")
	generator, _ := idgen.New(0, 0)

	startTime := time.Now()
	currentMs := startTime.UnixMilli()
	count := 0

	// Generate IDs until millisecond changes
	for {
		id := generator.Generate()
		count++

		// Extract timestamp from ID to see when it changed
		if generator.ExtractTimestamp(id) > currentMs {
			break
		}

		// Protection against infinite loop (theoretically should never happen)
		if count > 5000 {
			break
		}
	}

	fmt.Printf("   ✅ Generated %d IDs in 1 millisecond\n", count)
	fmt.Printf("   📊 Theoretical limit: 4096 IDs/ms\n")
	fmt.Println()

	// 2. Demonstrate continuous generation
	fmt.Println("2️⃣  2-second continuous generation test:")
	gen2, _ := idgen.New(5, 10)

	testDuration := 2 * time.Second
	start := time.Now()
	totalCount := 0

	// Generate for exactly 2 seconds
	for time.Since(start) < testDuration {
		gen2.Generate()
		totalCount++
	}

	elapsed := time.Since(start)
	idsPerSecond := float64(totalCount) / elapsed.Seconds()

	fmt.Printf("   ✅ Test duration: %v\n", elapsed)
	fmt.Printf("   📊 IDs generated: %d\n", totalCount)
	fmt.Printf("   📊 Average rate: %.0f IDs/second\n", idsPerSecond)
	fmt.Printf("   📊 Theoretical limit (1 worker): 4,096,000 IDs/second\n")
	fmt.Printf("   📊 Efficiency: %.1f%%\n", (idsPerSecond/4096000.0)*100)
	fmt.Println()

	// 3. Demonstrate with multiple workers
	fmt.Println("3️⃣  Simulating multiple workers (5 workers × 2 seconds):")

	numWorkers := 5
	testDuration = 2 * time.Second

	start = time.Now()

	// Channel to collect IDs from all workers
	type workerResult struct {
		workerID int
		count    int
	}
	resultsChan := make(chan workerResult, numWorkers)

	// Launch workers
	for w := 0; w < numWorkers; w++ {
		go func(workerID int) {
			gen, _ := idgen.New(0, int64(workerID))
			count := 0
			startTime := time.Now()

			for time.Since(startTime) < testDuration {
				gen.Generate()
				count++
			}

			resultsChan <- workerResult{workerID: workerID, count: count}
		}(w)
	}

	// Collect results
	totalGenerated := 0
	for i := 0; i < numWorkers; i++ {
		result := <-resultsChan
		totalGenerated += result.count
		fmt.Printf("   • Worker %d: %d IDs generated\n", result.workerID, result.count)
	}
	close(resultsChan)

	elapsed = time.Since(start)
	idsPerSecond = float64(totalGenerated) / elapsed.Seconds()

	fmt.Printf("\n   ✅ Test duration: %v\n", elapsed)
	fmt.Printf("   📊 Total IDs generated: %d\n", totalGenerated)
	fmt.Printf("   📊 Combined rate: %.0f IDs/second\n", idsPerSecond)
	fmt.Printf("   📊 Average per worker: %.0f IDs/second\n", idsPerSecond/float64(numWorkers))
	fmt.Printf("   📊 Theoretical limit (%d workers): %d IDs/second\n",
		numWorkers, numWorkers*4096000)
	fmt.Printf("   📊 Efficiency: %.1f%%\n", (idsPerSecond/float64(numWorkers*4096000))*100)
	fmt.Println()

	// 4. Theoretical capacity explanation
	fmt.Println("4️⃣  Total Theoretical Capacity:")
	fmt.Println("   ┌─────────────────────────────────────────────────┐")
	fmt.Println("   │ Sequence: 12 bits = 4096 IDs per millisecond   │")
	fmt.Println("   │ Workers: 32 workers × 32 processes = 1024      │")
	fmt.Println("   └─────────────────────────────────────────────────┘")
	fmt.Println()
	fmt.Println("   PER WORKER:")
	fmt.Println("   • 4,096 IDs per millisecond")
	fmt.Println("   • 4,096 × 1,000 ms = 4,096,000 IDs/second")
	fmt.Println()
	fmt.Println("   COMPLETE SYSTEM (1024 workers):")
	fmt.Println("   • 1,024 workers × 4,096 IDs/ms = 4,194,304 IDs/ms")
	fmt.Println("   • 4,194,304 × 1,000 = 4,194,304,000 IDs/second")
	fmt.Println("   • = ~4.2 BILLION IDs per second!")
	fmt.Println()

	// 5. Scalability comparison
	fmt.Println("5️⃣  Scalability Comparison:")
	fmt.Println("   ┌────────────────┬──────────────────────────┐")
	fmt.Println("   │ Configuration  │ IDs per second           │")
	fmt.Println("   ├────────────────┼──────────────────────────┤")
	fmt.Println("   │ 1 worker       │        ~4 million        │")
	fmt.Println("   │ 10 workers     │       ~40 million        │")
	fmt.Println("   │ 100 workers    │      ~400 million        │")
	fmt.Println("   │ 1024 workers   │      ~4.2 BILLION        │")
	fmt.Println("   └────────────────┴──────────────────────────┘")
	fmt.Println()

	fmt.Println("✨ The limitation is NOT per second, it's per MILLISECOND!")
	fmt.Println("   Sequence resets every millisecond, not every second.")
}
