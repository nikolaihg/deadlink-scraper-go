package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {
	const runs = 5
	const sleepDuration = 150 * time.Second // adjust as needed

	executables := map[string]string{
		"Concurrent":    "C:\\Users\\nikol\\Desktop\\deadlink-scraper-go\\tools\\benchmarks\builds\\deadlink-scraper-go-concurrent.exe",
		"NonConcurrent": "C:\\Users\\nikol\\Desktop\\deadlink-scraper-go\\tools\\benchmarks\builds\\deadlink-scraper-go-nonconcurrent.exe",
	}

	for label, exePath := range executables {
		fmt.Printf("\nRunning %s version:\n", label)

		var durations []time.Duration

		for i := 0; i < runs; i++ {
			start := time.Now()

			cmd := exec.Command(exePath, "https://scrape-me.dreamsofcode.io/")
			if err := cmd.Run(); err != nil {
				fmt.Printf("Run %d failed: %v\n", i+1, err)
				continue
			}

			elapsed := time.Since(start)
			durations = append(durations, elapsed)
			fmt.Printf("Run %d: %s\n", i+1, elapsed)

			time.Sleep(sleepDuration)
		}

		// Calculate average
		var total time.Duration
		for _, d := range durations {
			total += d
		}
		avg := total / time.Duration(len(durations))
		fmt.Printf("\nAverage time for %s: %s\n", label, avg)
	}
}
