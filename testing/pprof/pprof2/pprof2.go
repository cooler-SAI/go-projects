package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"strings"
)

func processStringConcatenationsWithBuilder(count int) string {
	var sb strings.Builder

	sb.Grow(count * len("some_string_part"))
	for i := 0; i < count; i++ {
		sb.WriteString("some_string_part")
	}
	return sb.String()
}

func main() {

	cpuProfileFile, err := os.Create("./pprof/pprof2/cpu_optimized.prof")
	if err != nil {
		panic(err)
	}
	defer func(cpuProfileFile *os.File) {
		err := cpuProfileFile.Close()
		if err != nil {
			fmt.Println("Could not close CPU profile file:", err)
		}
	}(cpuProfileFile)

	err = pprof.StartCPUProfile(cpuProfileFile)
	if err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()
	fmt.Println("Starting performance-intensive task (optimized)...")

	processStringConcatenationsWithBuilder(50000000)

	fmt.Println("Task finished. Profile data written to cpu_optimized.prof")

}
