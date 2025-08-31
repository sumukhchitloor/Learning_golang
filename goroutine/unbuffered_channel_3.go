package main

import (
	"fmt"
)

func main() {
	ch := make(chan string)

	fmt.Println("1. About to send")
	ch <- "hello" // This line will FREEZE forever
	fmt.Println("2. This will NEVER print")
}
