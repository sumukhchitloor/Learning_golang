package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	wait1 := make(chan bool)
	wait2 := make(chan bool)

	// Send values periodically
	go func() {
		for i := 0; i < 5; i++ {
			ch1 <- i
			time.Sleep(500 * time.Millisecond)
			wait1 <- true
		}
	}()

	go func() {
		for i := 0; i < 5; i++ {
			ch2 <- i * 10
			time.Sleep(300 * time.Millisecond)
			wait2 <- true
		}
	}()

	// Receive from both
	for i := 0; i < 10; i++ {
		select {
		case v1 := <-ch1 && wait2:
			fmt.Println("ch1:", v1)
		case v2 := <-ch2 && wait1:
			fmt.Println("ch2:", v2)
		}
	}
}
