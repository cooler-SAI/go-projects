package main

import (
	"log"
	"net/http"
	_ "net/http/pprof" // Automatically registers pprof handlers with http.DefaultServeMux
	"time"
)

// CPU-intensive function that performs useless calculations to consume CPU time
func cpuIntensiveTask() {
	for i := 0; i < 1000000; i++ {
		_ = i * i // Calculate square (result ignored)
	}
}

// Memory-intensive function that allocates and returns a 1MB byte slice
func memoryIntensiveTask() []byte {
	return make([]byte, 1024*1024) // 1 MB
}

func main() {
	// Start pprof server in a separate goroutine
	// This allows profiling the application without blocking the main thread
	go func() {
		// pprof will be available at http://localhost:6060/debug/pprof/
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Give the server time to start before beginning profiling
	time.Sleep(1 * time.Second)

	// Infinite loop to create controlled load
	// Runs in a separate goroutine to avoid blocking main
	go func() {
		var data [][]byte // Slice to store allocated memory (simulating leak)

		for {
			// Create CPU load
			cpuIntensiveTask()

			// Create memory load (simulating memory leak)
			// Each iteration adds 1MB to the slice, which is never cleaned up
			data = append(data, memoryIntensiveTask())

			// Pause between iterations to control load intensity
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Informational message about profiling server startup
	log.Println("Pprof server started on http://localhost:6060")
	log.Println("Available profiles:")
	log.Println("  http://localhost:6060/debug/pprof/")
	log.Println("  http://localhost:6060/debug/pprof/heap")
	log.Println("  http://localhost:6060/debug/pprof/profile")
	log.Println("  http://localhost:6060/debug/pprof/goroutine")
	log.Println("  http://localhost:6060/debug/pprof/block")
	log.Println("  http://localhost:6060/debug/pprof/mutex")
	log.Println("  http://localhost:6060/debug/pprof/trace?seconds=5")

	// Verify that the server is responding
	go func() {
		time.Sleep(2 * time.Second)
		resp, err := http.Get("http://localhost:6060/debug/pprof/")
		if err != nil {
			log.Printf("Error checking pprof: %v", err)
		} else {
			err := resp.Body.Close()
			if err != nil {
				log.Printf("Error closing response body: %v", err)
				return
			}
			log.Printf("✓ Pprof server is responding (status: %d)", resp.StatusCode)
		}
	}()

	// Infinite blocking to keep the program running
	select {} // Blocks forever without consuming CPU
}
