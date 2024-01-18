package main

import (
	"context"
	"fmt"
	"sync"
)

// Run
// go run main.go worker.go

func main() {
	inputSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Create channels for even and odd numbers
	evenChan := make(chan int, len(inputSlice))
	oddChan := make(chan int, len(inputSlice))

	workerSize := 5
	wp := NewWorkerPool(workerSize)
	wp.Run(context.Background())
	waitgroup := &sync.WaitGroup{}

	for _, num := range inputSlice {
		waitgroup.Add(1)
		data := num
		wp.AddTask(func() {
			if data%2 == 0 {
				evenChan <- data
			} else {
				oddChan <- data
			}
			waitgroup.Done()
		})
	}
	waitgroup.Wait()

	fmt.Println("EVEN")
	for len(evenChan) > 0 {
		fmt.Println(<-evenChan)
	}
	fmt.Println("ODD")
	for len(oddChan) > 0 {
		fmt.Println(<-oddChan)
	}

	defer func() {
		close(evenChan)
		close(oddChan)
	}()
}
