package processor

import (
	"testing"
)

func TestProcessLogFile_Integration(t *testing.T) {
	t.Parallel()

	keywords := []string{"INFO", "ERROR", "DEBUG"}
	filePath := "log.txt"
	numWorkers := 5

	result, err := ProcessLogFile(filePath, keywords, numWorkers)
	if err != nil {
		t.Fatalf("processLogFile failed: %v", err)
	}

	// Replace with the actual expected values from your log file
	expected := map[string]int{
		"INFO":  3,
		"ERROR": 2,
		"DEBUG": 1,
	}

	for keyword, expectedCount := range expected {
		if result[keyword] != expectedCount {
			t.Errorf("Expected %s count = %d, got = %d", keyword, expectedCount, result[keyword])
		}
	}
}
