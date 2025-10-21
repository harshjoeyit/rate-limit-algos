package slidingwincounter

import (
	"sync"
	"time"
)

// SlidingWindowCounter represents a sliding window rate limiter.
type SlidingWindowCounter struct {
	mu          sync.Mutex
	windowSize  time.Duration
	maxRequests int
	windows     map[int64]int
}

// NewSlidingWindowCounter initializes a new SlidingWindowCounter.
func NewSlidingWindowCounter(windowSize time.Duration, maxRequests int) *SlidingWindowCounter {
	return &SlidingWindowCounter{
		windowSize:  windowSize,
		maxRequests: maxRequests,
		windows:     make(map[int64]int),
	}
}

// Allow checks if a request can proceed based on the sliding window.
func (swc *SlidingWindowCounter) Allow() bool {
	swc.mu.Lock()
	defer swc.mu.Unlock()

	currentTime := time.Now().Unix()
	currentWindow := currentTime / int64(swc.windowSize.Seconds())

	// Clean up old windows to avoid memory buildup.
	for window := range swc.windows {
		if window < currentWindow-1 {
			delete(swc.windows, window)
		}
	}

	// Calculate requests in the current and previous window.
	currentCount := swc.windows[currentWindow] + 1 // counter + 1 for current request
	prevCount := swc.windows[currentWindow-1]

	// Calculate weighted sum for a smoother sliding window effect.
	elapsed := float64(time.Now().Unix()%int64(swc.windowSize.Seconds())) / float64(swc.windowSize.Seconds())
	totalRequests := float64(prevCount)*(1.0-elapsed) + float64(currentCount)

	if totalRequests <= float64(swc.maxRequests) {
		// Increment the count for the current window.
		swc.windows[currentWindow]++
		return true
	}

	return false
}
