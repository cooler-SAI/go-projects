package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func cpuIntensiveTask() {
	for i := 0; i < 1000000; i++ {
		_ = i * i
	}
}

func memoryIntensiveTask() []byte {
	return make([]byte, 1024*1024) // 1 MB
}

func main() {
	cpuFile, err := os.Create("./pprof/pprof4/cpu.prof")
	if err != nil {
		panic(err)
	}
	defer func(cpuFile *os.File) {
		err := cpuFile.Close()
		if err != nil {
			log.Println("Could not close CPU profile file:", err)
		}
	}(cpuFile)

	err2 := pprof.StartCPUProfile(cpuFile)
	if err2 != nil {
		return
	}
	defer pprof.StopCPUProfile()

	log.Println("Starting workload for 30 seconds...")

	start := time.Now()
	var data [][]byte

	for time.Since(start) < 30*time.Second {
		cpuIntensiveTask()
		data = append(data, memoryIntensiveTask())
		time.Sleep(100 * time.Millisecond)
	}
	log.Println("Workload completed. Saving memory profile...")
	memFile, err := os.Create("./pprof/pprof4/mem.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer func(memFile *os.File) {
		err := memFile.Close()
		if err != nil {
			log.Println("Could not close memory profile file:", err)
		}
	}(memFile)

	runtime.GC()
	err3 := pprof.WriteHeapProfile(memFile)
	if err3 != nil {
		log.Fatal(err)
		return
	}

	log.Println("Profiles saved successfully")
}
