package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func basicProducer(id int, jobs chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		num := rand.Intn(100)
		fmt.Printf("Producer : %d and producing job %d \n", id, num)
		jobs <- num
		time.Sleep(5 * time.Millisecond)
	}
}

func basicConsumer(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Consumer: %d processing job %d \n", id, job)
		time.Sleep(1 * time.Second)
	}
}

func MainExample() {
	now := time.Now()
	var wgProd sync.WaitGroup
	var wgCons sync.WaitGroup

	jobs := make(chan int, 10)

	for i := 0; i < 2; i++ {
		wgProd.Add(1)
		go basicProducer(i, jobs, &wgProd)
	}

	go func() {
		wgProd.Wait()
		close(jobs)
	}()

	for j := 0; j < 3; j++ {
		wgCons.Add(1)
		go basicConsumer(j, jobs, &wgCons)
	}

	wgCons.Wait()

	fmt.Printf("Taken time to complete task %v", time.Since(now))
}
