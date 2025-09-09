package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func deliverPizza(ctx context.Context, pizzaName string) {
	fmt.Printf("Pizza '%s' is being prepared...\n", pizzaName)

	prepareRand := rand.Intn(10) + 25 // Random between 25-34 seconds
	select {
	case <-time.After(time.Duration(prepareRand) * time.Second):
		fmt.Printf("Pizza '%s' delivered! 🍕\n", pizzaName)
	case <-ctx.Done():
		fmt.Printf("Pizza '%s' cancelled: %v\n", pizzaName, ctx.Err())
	}
}

func main() {
	fmt.Println("Preparing to order pizza...")

	waitingRand := rand.Intn(10) + 25
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(waitingRand)*time.Second) // 20 seconds timeout
	defer cancel()

	// Start pizza delivery in a separate goroutine
	go deliverPizza(ctx, "Pepperoni")

	fmt.Println("Wait to see if pizza is delivered or cancelled.....")
	time.Sleep(30 * time.Second)
	fmt.Println("Finished waiting.")

}
