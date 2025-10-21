package main

import (
	"fmt"
	"time"

	"github.com/harshjoeyit/myratelimiter/leakybucket"
)

func main() {
	// rate := 10.0 / 60.0 // 10 request in 60 seconds
	// capacity := 10.0
	// tb := tockenbucket.NewTocketBucket(capacity, rate)

	// allowed := 0
	// now := time.Now()

	// for i := 1; i <= 120; i++ {
	// 	if tb.Allow() {
	// 		fmt.Printf("Request %d allowed at %v\n", i, time.Now())
	// 		allowed++
	// 	} else {
	// 		fmt.Printf("Request %d denied at %v\n", i, time.Now())
	// 	}
	// 	time.Sleep(400 * time.Millisecond) // Simulate 5 requests per second.
	// }

	// fmt.Println("allowed: ", allowed, "in: ", now.Sub(time.Now()).Seconds())

	// s := slidingwincounter.NewSlidingWindowCounter(60*time.Second, 10)
	// allowed := 0
	// now := time.Now()

	// for i := 1; i <= 120; i++ {
	// 	if s.Allow() {
	// 		fmt.Printf("Request %d allowed at %v\n", i, time.Now())
	// 		allowed++
	// 	} else {
	// 		fmt.Printf("Request %d denied at %v\n", i, time.Now())
	// 	}
	// 	time.Sleep(800 * time.Millisecond) // Simulate 5 requests per second.
	// }
	//
	// fmt.Println("allowed: ", allowed, "in: ", time.Until(now).Seconds())

	allowed := 0
	now := time.Now()

	lbq := leakybucket.NewLeakyBucketQueue(10, 200*time.Millisecond)
	for i := 1; i <= 30; i++ {
		if lbq.Allow() {
			allowed++
			fmt.Printf("Leaky Bucket: Request %d allowed at %v\n", i, time.Now())
		} else {
			fmt.Printf("Leaky Bucket: Request %d denied at %v\n", i, time.Now())
		}
		time.Sleep(50 * time.Millisecond) // Simulate incoming requests.
	}

	// Let the leaky bucket process remaining requests
	time.Sleep(5 * time.Second)

	fmt.Println("allowed: ", allowed, "in: ", time.Until(now).Seconds())
}
