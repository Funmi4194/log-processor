package processor

import (
	"sync"
	"testing"
)

// TestCountLine tests the CountLine function with different keyword matching scenarios
func TestCountLine(t *testing.T) {
	t.Parallel()

	keywords := []string{"DEBUG", "ERROR", "INFO"}

	tests := []struct {
		name     string
		line     string
		expected map[string]int
	}{
		{"error", "ERROR, error", map[string]int{"ERROR": 1}},
		{"debug", "DEBUG, debug, INFO", map[string]int{"DEBUG": 1, "INFO": 1}},
		{"info", "INFO:info", map[string]int{"INFO": 1}},
		{"No Match", "no match", map[string]int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountLine(tt.line, keywords)

			// verify each expected keyword count
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("expected %v for %s, got %v", v, k, result[k])
				}
			}
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d keys, got %d", len(tt.expected), len(result))
			}
		})
	}
}

// TestAggregateCounts checks whether multiple maps sent over a channel are correctly merged
func TestAggregateCounts(t *testing.T) {
	countChan := make(chan map[string]int, 3)
	done := make(chan struct{})

	// simulate 3 batches of count maps
	countChan <- map[string]int{"ERROR": 1, "DEBUG": 2}
	countChan <- map[string]int{"ERROR": 3, "INFO": 1}
	countChan <- map[string]int{"DEBUG": 1}
	close(countChan)

	// call AggregateCounts and wait for the done signal
	result := AggregateCounts(countChan, done)
	<-done

	// expected combined result
	expected := map[string]int{
		"ERROR": 4,
		"DEBUG": 3,
		"INFO":  1,
	}

	// validate final aggregated map
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("expected %d for %s, got %d", v, k, result[k])
		}
	}
}

// TestStartWorkerPool tests the full worker pipeline: lines, counting, and aggregating
func TestStartWorkerPool(t *testing.T) {
	keywords := []string{"INFO", "ERROR", "DEBUG"}

	// Lines to simulate file lines
	lines := []string{
		"2023-10-28 12:00:01 - INFO - User logged in",
		"2023-10-28 12:00:03 - ERROR - Database connection failed",
		"2023-10-28 12:00:04 - DEBUG - Cache hit for request",
		"2023-10-28 12:00:05 - INFO - User logged out",
		"2023-10-28 12:00:07 - ERROR - Timeout reached while processing",
		"2023-10-28 12:00:08 - INFO - User logged in",
	}

	// Create input and output channels
	lineChan := make(chan string, len(lines))
	countChan := make(chan map[string]int, 2) // Small buffer for batching
	var wg sync.WaitGroup                     // WaitGroup to sync worker goroutines

	// Feed lines into the channel
	for _, line := range lines {
		lineChan <- line
	}
	close(lineChan) // Important: close input channel so workers can terminate

	// Start worker pool with 2 goroutines
	StartWorkerPool(2, keywords, lineChan, countChan, &wg)

	// Close countChan after all workers are done
	go func() {
		wg.Wait()
		close(countChan)
	}()

	done := make(chan struct{})
	final := AggregateCounts(countChan, done) // Aggregate results from workers
	<-done                                    // Wait for aggregation to complete

	// Expected final result
	expected := map[string]int{
		"DEBUG": 1,
		"ERROR": 2,
		"INFO":  3,
	}

	// Validate aggregated output
	for k, v := range expected {
		if final[k] != v {
			t.Errorf("expected %d for %s, got %d", v, k, final[k])
		}
	}
}
