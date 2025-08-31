package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Result represents a search result
type Result string

// Search is a function type that takes a query and returns a Result
type Search func(query string) Result

// fakeSearch returns a Search function that simulates a delay and returns a formatted result
func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

// Create three fake search functions for different types
var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")
)

// Google performs all three searches sequentially and returns the results
func Google(query string) []Result {
	return []Result{
		Web(query),
		Image(query),
		Video(query),
	}
}

// Main function to test the framework
func main() {
	start := time.Now()
	results := Google("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
