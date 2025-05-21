package processor

import (
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {

}

// helper function to create a temporary file
func createTempFile(t *testing.T, content string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("could not write to temp file: %v", err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}

func TestReadFile_FileNotExist(t *testing.T) {
	t.Parallel()
	_, err := ReadFile("nonexistent.txt")
	if err == nil || err.Error() != "file path does not exist: nonexistent.txt" {
		t.Errorf("expected file path does not exist error, got: %v", err)
	}
}

func TestReadFile_EmptyFile(t *testing.T) {
	t.Parallel()
	path := createTempFile(t, "")
	defer os.Remove(path)
	_, err := ReadFile(path)
	if err == nil || err.Error() != "file is empty" {
		t.Errorf("expected file is empty error, got: %v", err)
	}
}

func TestReadFileLines(t *testing.T) {
	t.Parallel()
	content := "2025-05-21 10:01:20 - DEBUG - Admin token issued\n2025-05-21 10:01:22 - INFO - Logs rotated\n2025-05-21 10:01:21 - ERROR - Unauthorized admin access"
	path := createTempFile(t, content)
	defer os.Remove(path)

	file, err := ReadFile(path)
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	defer file.Close()

	lines := make(chan string, 100)
	err = ReadFileLines(file, lines)
	if err != nil {
		t.Errorf("expected no error reading lines, got: %v", err)
	}

	close(lines)
	expected := []string{"2025-05-21 10:01:20 - DEBUG - Admin token issued", "2025-05-21 10:01:22 - INFO - Logs rotated", "2025-05-21 10:01:21 - ERROR - Unauthorized admin access"}
	i := 0
	for line := range lines {
		if line != expected[i] {
			t.Errorf("expected line '%s', got '%s'", expected[i], line)
		}
		i++
	}
}
