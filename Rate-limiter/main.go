package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type TokenBucket struct {
	capacity   int
	tokens     int
	reFillRate int
	lastRefill time.Time
	mutex      sync.Mutex
}

func NewTokenBucket(capacity, refillRate int) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		reFillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Refill() {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()

	newToken := int(elapsed * float64(tb.reFillRate))

	if newToken > 0 {
		tb.tokens = min(tb.capacity, tb.tokens+newToken)
		tb.lastRefill = now
	}
}

func (tb *TokenBucket) AllowRequest() bool {
	tb.Refill()

	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

func MakeRequest(tb *TokenBucket, userID int, requestCount int) {
	for count := 1; count < requestCount; count++ {
		if tb.AllowRequest() {
			fmt.Printf("User %d: Request %d allowed at %s\n", userID, count, time.Now().Format("15:04:05"))
		} else {
			fmt.Printf("User %d: Request %d denied (Rate Limited) at %s\n", userID, requestCount, time.Now().Format("15:04:05"))
		}

		time.Sleep(time.Duration(rand.Intn(800)+200) * time.Millisecond)
	}

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {

	limiter := NewTokenBucket(5, 2)

	var wg sync.WaitGroup

	user := 10

	for u := 1; u <= user; u++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			MakeRequest(limiter, userID, 10)
		}(u)
	}

	wg.Wait()

	fmt.Println("Completed task")

}
