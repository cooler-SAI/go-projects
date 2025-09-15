package main

import (
	"fmt"
	"os"
	"runtime/pprof"
)

func processStringConcatenations(count int) string {
	var result string
	for i := 0; i < count; i++ {
		result += "some_string_part"
	}
	return result
}

func main() {
	cpuProfileF, err := os.Create("./pprof/cpu_profile.prof")
	if err != nil {
		fmt.Println("Error exist creating CPU profile file:", err)
		panic(err)

	}
	defer func(cpuProfileF *os.File) {
		err := cpuProfileF.Close()
		if err != nil {
			fmt.Println("Error exist closing CPU profile file:", err)
		}
	}(cpuProfileF)

	err2 := pprof.StartCPUProfile(cpuProfileF)
	if err2 != nil {
		return
	}
	defer pprof.StopCPUProfile()
	processStringConcatenations(100000)

	fmt.Println("CPU profiling completed and saved to cpu_profile.prof")

}
