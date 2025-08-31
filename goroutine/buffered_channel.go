package main

import (
	"fmt"
)

func main() {
	ch := make(chan string, 2) // Buffer size 2

	fmt.Println("1. About to send first")
	ch <- "first" // Goes into slot 1
	fmt.Println("2. Sent first (didn't block!)")

	fmt.Println("3. About to send second")
	ch <- "second" // Goes into slot 2
	fmt.Println("4. Sent second (didn't block!)")
	ch <- "third" // This will block until a slot is free
	fmt.Println("5. Buffer is now full")
}
