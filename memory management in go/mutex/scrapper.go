package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func countLetters(url string, frequency *[26]int, mu *sync.Mutex) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic("Error: status code \n" + resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	mu.Lock()
	for _, b := range body {
		c := strings.ToLower(string(b))
		cIndex := strings.Index(allLetters, c)
		if cIndex != -1 {
			frequency[cIndex]++
		}
	}
	fmt.Println("Completed:", url)
	mu.Unlock()
}

func main() {
	mu := sync.Mutex{}
	var frequency [26]int
	for i := 1000; i <= 1010; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countLetters(url, &frequency, &mu)
	}
	time.Sleep(2 * time.Second)
	mu.Lock()
	for i, c := range allLetters {

		fmt.Printf("%c : %d\n", c, frequency[i])
	}
	mu.Unlock()

}
