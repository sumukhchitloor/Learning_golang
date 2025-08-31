package main

import (
	"fmt"
	"math/rand"
	"time"
)

const numLoop = 3 // Configurable number of iterations

func boring(msg string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c
}

func main() {
	c := boring("Hello")           // Configurable message
	for i := 0; i < numLoop; i++ { // Configurable count
		fmt.Printf("You say: %q\n", <-c)
	}
	fmt.Println("You're boring; I'm leaving.")
}
