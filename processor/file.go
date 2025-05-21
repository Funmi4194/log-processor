package processor

import (
	"bufio"
	"fmt"
	"os"
)

// ReadFile opens a file from the specified filepath,
// returns an error if the file doesn't exist, is empty, or can't be accessed.
func ReadFile(filepath string) (*os.File, error) {
	// open text file
	file, err := os.Open(filepath)
	if err != nil {
		// check if file does not exist
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file path does not exist: %s", filepath)
		}
		return nil, fmt.Errorf("file could not be opened: %s", err.Error())
	}

	// check if file data is empty
	fileDetail, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("file detail could not be opened: %s", err.Error())
	}

	if fileDetail.Size() == 0 {
		file.Close()
		return nil, fmt.Errorf("file is empty")
	}

	return file, nil
}

// ReadFileLines reads a file line-by-line and sends each line to the provided channel.
func ReadFileLines(file *os.File, lineChan chan<- string) error {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineChan <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
