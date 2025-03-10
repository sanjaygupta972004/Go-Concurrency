package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	MaximumWorkers = 3
	RequestTimeOut = 5 * time.Second
	RateLimit      = 2
)

type Crawler struct {
	queue       chan string
	visitedURLs map[string]bool
	mu          sync.Mutex
	wg          sync.WaitGroup
	limiter     *rate.Limiter
}

func NewCrawler() *Crawler {
	return &Crawler{
		queue:       make(chan string, MaximumWorkers),
		visitedURLs: map[string]bool{},
		limiter:     rate.NewLimiter(rate.Every(time.Second/RateLimit), 1),
	}
}

func (c *Crawler) crawl(url string) {
	defer c.wg.Done()

	c.mu.Lock()

	if c.visitedURLs[url] {
		c.mu.Unlock()
		return
	}
	c.visitedURLs[url] = true

	c.mu.Unlock()

	// implement rate limit

	if err := c.limiter.Wait(context.Background()); err != nil {
		fmt.Println("Rate limiter error:", err)
		return
	}

	// make http request

	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeOut)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error fetching:", url, err)
		return
	}

	fmt.Printf("Response from url %v \n", resp.Body)

	defer resp.Body.Close()

	fmt.Printf("Fetched url %s: status code : %d", url, resp.StatusCode)
}

func (c *Crawler) Worker() {
	for url := range c.queue {
		c.crawl(url)
	}
}

func (c *Crawler) Start(startURLs []string) {

	for i := 0; i < MaximumWorkers; i++ {
		go c.Worker()
	}

	for _, url := range startURLs {
		c.wg.Add(1)
		c.mu.Lock()
		c.queue <- url
		c.mu.Unlock()
	}

	go func() {
		c.wg.Wait()
		close(c.queue)
	}()
}

func MainCrawl() {
	crawler := NewCrawler()

	urlStrings := []string{"https://google.com", "https://bewanderic.com", "https://carveo.site", "https://github.com"}

	crawler.Start(urlStrings)
}
