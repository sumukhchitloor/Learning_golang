package main

import (
	"fmt"
	"sync"
	"time"
)

func stingy(money *int, mu *sync.Mutex) {
	for i := 0; i < 1000000; i++ {
		mu.Lock()
		*money += 10
		mu.Unlock()
	}
	fmt.Println("stingy done")
}

func spendy(money *int, mu *sync.Mutex) {
	for i := 0; i < 1000000; i++ {
		mu.Lock()
		*money -= 10
		mu.Unlock()
	}
	fmt.Println("spendy done")
}

func main() {
	mu := sync.Mutex{}
	money := 100
	go stingy(&money, &mu)
	go spendy(&money, &mu)
	fmt.Println("waiting...")
	time.Sleep(2 * time.Second)

	mu.Lock()
	fmt.Println("money:", money)
	mu.Unlock()

}
