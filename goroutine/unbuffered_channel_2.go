package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)

	go func() {
		fmt.Println("Goroutine: About to send")
		ch <- "message"                 // Blocks here until main receives
		fmt.Println("Goroutine: Sent!") // This prints AFTER main receives
	}()

	time.Sleep(1 * time.Second) // Let goroutine get to the send line
	fmt.Println("Main: About to receive")
	msg := <-ch // This unblocks the goroutine
	fmt.Println("Main: Received:", msg)
}
