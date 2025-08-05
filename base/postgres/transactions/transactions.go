package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	doneChan := make(chan struct{})

	go func() {
		<-stopChan
		fmt.Println("\nReceived shutdown signal...")
		close(doneChan)
		time.Sleep(300 * time.Millisecond)
	}()

	for {
		select {
		case <-doneChan:
			fmt.Println("Shutting down gracefully...")
			os.Exit(0)
		default:
			fmt.Println("Program Running.....")
			time.Sleep(1 * time.Second)
		}
	}
}
