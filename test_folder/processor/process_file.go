package processor

import (
	"bufio"
	"fmt"
	"sort"
	"strings"
	"sync"
)

func ProcessLogFile(filepath string, keywords []string, numWorkers int) (map[string]int, error) {
	file, err := ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// use a buffered channel to allow file reader goroutine to push up to 100
	// items into the channel without blocking as long as the buffer isn't full
	lineChan := make(chan string, 100)
	countChan := make(chan map[string]int, numWorkers)
	done := make(chan struct{})
	var wg sync.WaitGroup

	// start the aggregator in the background
	var finalCount map[string]int
	go func() {
		finalCount = AggregateCounts(countChan, done)
	}()

	// start workers
	StartWorkerPool(numWorkers, keywords, lineChan, countChan, &wg)

	// read file line by line and stream to workers
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineChan <- scanner.Text()
	}
	close(lineChan)

	// wait for workers to finish and close countChan
	wg.Wait()
	close(countChan)

	// wait for aggregator to finish
	<-done
	return finalCount, nil
}

// StartWorkerPool spins up multiple goroutines to process lines from lineChan,
// counts keyword occurrences using CountLine, and sends the results to countChan in batches.
func StartWorkerPool(numWorkers int, keywords []string, lineChan <-chan string, countChan chan<- map[string]int, wg *sync.WaitGroup) {
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			localCount := make(map[string]int)
			batchSize := 1000
			count := 0

			// continuously read lines from the input channel until it is closed
			for line := range lineChan {
				lineCount := CountLine(line, keywords)

				// add counts from lineCount to the local batch count
				for k, v := range lineCount {
					localCount[k] += v
				}

				// increment the batch line count
				count++

				// when batch size is reached, send the batch counts to the output channel
				if count >= batchSize {
					countChan <- localCount
					localCount = make(map[string]int)
					count = 0
				}
			}

			// after all lines are processed, if any counts remain in the local map, send them
			if count > 0 {
				countChan <- localCount
			}
		}()
	}
}

// CountLine processes a single line and counts keyword matches (case-insensitive).
func CountLine(line string, keywords []string) map[string]int {
	counts := make(map[string]int)
	line = strings.ToUpper(line)

	for _, keyword := range keywords {
		if strings.Contains(line, strings.ToUpper(keyword)) {
			counts[strings.ToUpper(keyword)]++
		}
	}
	return counts
}

// AggregateCounts merges all maps from countChan into a final count map.
func AggregateCounts(countChan <-chan map[string]int, done chan struct{}) map[string]int {
	// map to store the aggregated counts from all workers
	finalCount := make(map[string]int)

	// loop over maps received on countChan until it is closed
	for counts := range countChan {
		for k, v := range counts {
			finalCount[k] += v
		}
	}
	close(done)
	return finalCount
}

// PrintSortedCounts prints the final result in descending order.
func PrintSortedCounts(counts map[string]int) {
	type kv struct {
		Key   string
		Value int
	}
	var sorted []kv
	for k, v := range counts {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})
	for _, entry := range sorted {
		fmt.Printf("%s: %d\n", entry.Key, entry.Value)
	}
}
