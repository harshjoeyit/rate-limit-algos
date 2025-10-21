package leakybucket

import (
	"fmt"
	"sync"
	"time"
)

type LeakyBucketQueue struct {
	capacity int           // maximum number of requests in the bucket
	rate     time.Duration // time duration between processing each request
	queue    chan struct{}
	mu       sync.Mutex
	ticker   *time.Ticker
}

func NewLeakyBucketQueue(capacity int, rate time.Duration) *LeakyBucketQueue {
	lbq := &LeakyBucketQueue{
		capacity: capacity,
		rate:     rate,
		queue:    make(chan struct{}, capacity),
		ticker:   time.NewTicker(rate),
	}

	go lbq.leak()

	return lbq
}

func (lbq *LeakyBucketQueue) leak() {
	for range lbq.ticker.C {
		lbq.mu.Lock()
		select {
		case <-lbq.queue:
			// Process the request (in this case, just removing it from the queue)
			fmt.Println("Request processed")
		default:
			// No requests to process
			fmt.Println("No requests to process, bucket empty")
		}
		lbq.mu.Unlock()
	}
}

func (lbq *LeakyBucketQueue) Allow() bool {
	lbq.mu.Lock()
	defer lbq.mu.Unlock()

	select {
	case lbq.queue <- struct{}{}:
		// Request added to the bucket
		fmt.Println("Request allowed")
		return true
	default:
		// Bucket is full, request denied
		fmt.Println("Request denied, bucket full")
		return false
	}
}
