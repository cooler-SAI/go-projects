package main

import (
	"context"
	"fmt"
	"time"
)

func doSomething(ctx context.Context) {
	for {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("Doing something...")

		case <-ctx.Done():
			fmt.Println("Context is done")
		}
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	go doSomething(ctx)

	time.Sleep(5 * time.Second)

}
