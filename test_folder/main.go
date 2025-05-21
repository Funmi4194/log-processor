package main

import (
	"fmt"

	"github.com/funmi4194/log-processor/processor"
)

func main() {
	keywords := []string{"INFO", "ERROR", "DEBUG"}
	filePath := "processor/log.txt"
	numWorkers := 4

	counts, err := processor.ProcessLogFile(filePath, keywords, numWorkers)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	processor.PrintSortedCounts(counts)
}
