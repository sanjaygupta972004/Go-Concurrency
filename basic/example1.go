package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, res chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker ID %d is doing job : %d \n ", id, job)
		time.Sleep(2 * time.Second)
		result := fmt.Sprintf("Completed task %d", job)
		res <- result
	}

}

func WorkJobProblem() {

	numOfJobs := 10
	numOfWorkers := 3

	var wg sync.WaitGroup

	job := make(chan int, numOfJobs)
	res := make(chan string, numOfJobs)

	for i := 0; i < numOfWorkers; i++ {
		wg.Add(1)
		go worker(i, job, res, &wg)
	}

	for i := 0; i < numOfJobs; i++ {
		job <- i
	}
	close(job)

	go func() {
		wg.Wait()
		close(res)
	}()

	for re := range res {
		fmt.Println(re)
	}
}
