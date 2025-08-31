package main

import (
	"fmt"
	"time"
)

// Simple rate limiter - sends a signal every 'd' duration
func rateLimiter(d time.Duration) <-chan time.Time {
	ch := make(chan time.Time)
	go func() {
		ticker := time.NewTicker(d)
		defer ticker.Stop()
		for {
			ch <- <-ticker.C // Send tick every 'd' duration
		}
	}()
	return ch
}

func main() {
	fmt.Println("=== Without Rate Limiting ===")
	start := time.Now()

	// Fast - no rate limiting
	for i := 1; i <= 5; i++ {
		fmt.Printf("Fast: Processing item %d at %v\n", i, time.Since(start).Round(time.Millisecond))
	}

	fmt.Println("\n=== With Rate Limiting ===")
	start = time.Now()
	limiter := rateLimiter(500 * time.Millisecond) // Every 500ms

	// Slow - rate limited to every 500ms
	for i := 1; i <= 5; i++ {
		<-limiter // Wait for permission
		fmt.Printf("Slow: Processing item %d at %v\n", i, time.Since(start).Round(time.Millisecond))
	}

	fmt.Println("\n=== API Call Example ===")
	start = time.Now()
	apiLimiter := rateLimiter(1 * time.Second) // 1 call per second

	urls := []string{"api1", "api2", "api3"}
	for _, url := range urls {
		<-apiLimiter // Wait for rate limit
		fmt.Printf("Calling %s at %v\n", url, time.Since(start).Round(time.Millisecond))
	}
}
