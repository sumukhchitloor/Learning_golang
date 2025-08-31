package main

import (
	"fmt"
	"strconv"
)

// Stage 1: Generate numbers
func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out) // Ensure channel is closed
		for _, n := range nums {
			fmt.Printf("Gen: sending %d\n", n)
			out <- n
		}
		fmt.Println("Gen: finished")
	}()
	return out
}

// Stage 2: Square numbers
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			squared := n * n
			fmt.Printf("Square: %d -> %d\n", n, squared)
			out <- squared
		}
		fmt.Println("Square: finished")
	}()
	return out
}

// Stage 3: Convert to string and add prefix
func toString(in <-chan int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for n := range in {
			result := "Result: " + strconv.Itoa(n)
			fmt.Printf("ToString: %d -> %s\n", n, result)
			out <- result
		}
		fmt.Println("ToString: finished")
	}()
	return out
}

// Stage 4: Filter only large numbers (> 10)
func filter(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for s := range in {
			// Extract number from string for filtering
			if len(s) > 15 { // "Result: " + number > 10
				fmt.Printf("Filter: passing %s\n", s)
				out <- s
			} else {
				fmt.Printf("Filter: blocking %s\n", s)
			}
		}
		fmt.Println("Filter: finished")
	}()
	return out
}

func main() {
	fmt.Println("=== Simple Pipeline ===")
	// Simple pipeline: gen -> sq
	for n := range sq(gen(2, 3, 4)) {
		fmt.Printf("Final result: %d\n", n)
	}

	fmt.Println("\n=== Complex Pipeline ===")
	// Complex pipeline: gen -> sq -> toString -> filter
	pipeline := filter(toString(sq(gen(1, 2, 3, 4, 5))))

	for result := range pipeline {
		fmt.Printf("FINAL OUTPUT: %s\n", result)
	}

	fmt.Println("Pipeline complete!")
}
