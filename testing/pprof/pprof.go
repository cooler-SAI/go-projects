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
	cpuProfileFile, err := os.Create("./pprof/cpu_optimized.prof")
	if err != nil {
		panic(err)
	}
	defer func(cpuProfileFile *os.File) {
		err := cpuProfileFile.Close()
		if err != nil {
			fmt.Println("Could not close CPU profile file:", err)
		}
	}(cpuProfileFile)

	err2 := pprof.StartCPUProfile(cpuProfileFile)
	if err2 != nil {
		fmt.Println("Could not start CPU profile:", err2)
		return
	}
	defer pprof.StopCPUProfile()

	fmt.Println("Starting performance-intensive task (optimized)...")

	processStringConcatenationsWithBuilder(1000000)

	fmt.Println("Task finished. Profile data written to cpu_optimized.prof")

}
