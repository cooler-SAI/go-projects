package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "userID", 123)

	printUserID(ctx)
}

func printUserID(ctx context.Context) {
	if userID, ok := ctx.Value("userID").(int); ok {
		fmt.Printf("User ID: %d\n", userID)
	} else {
		fmt.Println("User ID not found or wrong type")
	}
}
