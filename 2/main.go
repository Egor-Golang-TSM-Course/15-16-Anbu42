package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func worker(ctx context.Context, name string, input <-chan int, output chan<- int) {
	defer func() {
		if output != nil {
			close(output)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s: Context is canceled. Exiting.\n", name)
			return
		case data, ok := <-input:
			if !ok {
				fmt.Printf("%s: Input channel closed. Exiting.\n", name)
				return
			}

			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

			result := data * 2
			fmt.Printf("%s: Processed %d, result %d\n", name, data, result)
			if output != nil {
				output <- result
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)
	ch4 := make(chan int)

	var wg sync.WaitGroup

	wg.Add(4)
	go func() {
		defer wg.Done()
		worker(ctx, "Worker1", ch1, ch2)
	}()
	go func() {
		defer wg.Done()
		worker(ctx, "Worker2", ch2, ch3)
	}()
	go func() {
		defer wg.Done()
		worker(ctx, "Worker3", ch3, ch4)
	}()
	go func() {
		defer wg.Done()
		worker(ctx, "Worker4", ch4, nil)
	}()

	ch1 <- 10

	close(ch1) // Закрываем ch1, чтобы все горутины завершили выполнение

	wg.Wait()

	fmt.Println("All workers have completed.")
}
