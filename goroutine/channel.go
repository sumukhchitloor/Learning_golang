package main

import (
	"fmt"
)

func main() {
	ch := make(chan string)

	// Start a goroutine
	go func() {
		ch <- "Hello from goroutine!" // Send message
	}()

	message := <-ch // Receive message
	fmt.Println(message)
}
