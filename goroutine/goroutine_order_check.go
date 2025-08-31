package main

import (
	"fmt"
	"time"
)

func main() {
	// Start goroutines in order 1, 2, 3
	go func() {
		fmt.Println("Goroutine 1 executed")
	}()

	go func() {
		fmt.Println("Goroutine 2 executed")
	}()

	go func() {
		fmt.Println("Goroutine 3 executed")
	}()

	time.Sleep(100 * time.Millisecond)
}
