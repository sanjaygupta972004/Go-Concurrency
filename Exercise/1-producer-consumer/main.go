// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently

package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream, jobs chan<- *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(jobs)
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			fmt.Println("Error: File is ended", err)
			return
		}

		jobs <- tweet
	}
}

func consumer(jobs <-chan *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range jobs {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {

	jobs := make(chan *Tweet)
	start := time.Now()
	stream := GetMockStream()

	var wg sync.WaitGroup
	wg.Add(1)

	// Producer
	go producer(stream, jobs, &wg)

	// consumer
	wg.Add(1)
	go consumer(jobs, &wg)
	wg.Wait()

	// running another example

	MainExample()

	fmt.Printf("Process took %s\n", time.Since(start))
}
