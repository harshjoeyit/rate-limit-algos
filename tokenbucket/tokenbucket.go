package tockenbucket

import (
	"fmt"
	"sync"
	"time"
)

// TokenBucket struct represents the token bucket.
type TokenBucket struct {
	mu         sync.Mutex
	tokens     float64 // current number of tokens
	capacity   float64 // maximum number of tokens
	refillRate float64 // tokens added per second
	lastRefill time.Time
}

func NewTokenBucket(capacity, refillRate float64) *TokenBucket {
	return &TokenBucket{
		mu:         sync.Mutex{},
		tokens:     capacity,
		capacity:   capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (t *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(t.lastRefill).Seconds()
	t.tokens += t.refillRate * elapsed
	if t.tokens > t.capacity {
		t.tokens = t.capacity
	}
	t.lastRefill = now
}

func (t *TokenBucket) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.refill()

	if t.tokens >= 1 {
		t.tokens--
		fmt.Println("Tokens", t.tokens)
		return true
	}

	return false
}
