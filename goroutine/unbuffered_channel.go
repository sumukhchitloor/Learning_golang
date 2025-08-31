package main

import (
	"fmt"
	"time"
)

// This program demonstrates an unbuffered channel where the sender and receiver
// must synchronize their actions. The goroutine sends a message to the main function,
// which waits to receive it.

func main() {
	ch := make(chan string)

	// This would BLOCK FOREVER because no one is receiving yet
	// ch <- "Hello"  // ❌ This line would deadlock!

	// Solution: Start a receiver first
	go func() {
		msg := <-ch // This goroutine waits to receive
		fmt.Println("Received:", msg)
	}()

	ch <- "Hello"                       // Now this can send (receiver is ready)
	time.Sleep(1000 * time.Millisecond) // Give time to print
}
