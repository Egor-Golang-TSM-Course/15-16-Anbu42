package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func makeRequest(ctx context.Context, client *http.Client, url string, wg *sync.WaitGroup) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Error creating request for %s: %v\n", url, err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request to %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Response from %s\n", url)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	client := &http.Client{}

	var wg sync.WaitGroup

	urls := []string{
		"https://www.example.com",
		"https://www.example.org",
		"https://www.example.net",
		"https://www.google.com/",
		"https://www.netflix.com/",
	}

	for _, url := range urls {
		wg.Add(1)
		go makeRequest(ctx, client, url, &wg)
	}

	wg.Wait()
}
