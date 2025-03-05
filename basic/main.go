package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func sayHello() {
	fmt.Println("Hello concurrency")
}

// how to create multiple go
func downLoadfile(file string, message <-chan string) {
	defer wg.Done() // decrement work

	fmt.Println(<-message)

	time.Sleep(2 * time.Second) // simulating the time to complete the work

	fmt.Println("Downloaded file", file)

}

// func printNumber(i int) {
// 	defer wg.Done()

// 	fmt.Printf("Start Working : %d\n", i)

// 	time.Sleep(2 * time.Second)

// 	fmt.Printf("Completed task : %d", i)
// }

func main() {
	fmt.Println("Hello go")
	go sayHello()

	start := time.Now()

	msg := make(chan string, 3)

	files := []string{"file1.png", "file2.png", "file3.png"}

	for _, file := range files {
		wg.Add(1)
		go downLoadfile(file, msg)
		msg <- "Downloading file"
	}

	go func() {
		wg.Wait()
		close(msg)
	}()

	// for i := 0; i < 500; i++ {
	// 	wg.Add(1)
	// 	go printNumber(i)
	// }

	// example number one

	WorkJobProblem()
	PipelineExa()
	fmt.Println("Total time taken ", time.Since(start))
}
