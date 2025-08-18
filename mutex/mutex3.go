package main

import "sync"

func withoutMutexNew() {
	counter := 10000
	var wg sync.WaitGroup

	wg.Add(counter)
	for i := 0; i < counter; i++ {
		go func() {
			defer wg.Done()
			counter++
		}()
	}
	wg.Wait()
	println(counter)

}

func main() {
	withoutMutexNew()
}
