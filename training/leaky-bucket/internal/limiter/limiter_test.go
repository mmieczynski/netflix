package limiter

import (
	"context"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

func TestCorrectlyDrainTheBucket(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const interval = 5 * time.Second
		limiter := NewLeakyLimiter(interval, 2)
		executed := make(chan int, 3)

		// Fill the queue before starting the worker so capacity is deterministic.
		if !limiter.Submit(func() { executed <- 1 }) {
			t.Fatal("first job should be accepted")
		}
		if !limiter.Submit(func() { executed <- 2 }) {
			t.Fatal("second job should be accepted")
		}
		if limiter.Submit(func() { executed <- 3 }) {
			t.Fatal("third job should be rejected because the queue is full")
		}

		ctx, cancel := context.WithCancel(context.Background())
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			limiter.Run(ctx)
		}()
		defer func() {
			cancel()
			wg.Wait()
		}()

		start := time.Now()
		assertCount := func(want int) {
			t.Helper()
			synctest.Wait()
			if got := len(executed); got != want {
				t.Fatalf("at %v: executed %d jobs, want %d", time.Since(start), got, want)
			}
		}

		// Let Run initialize its timer before advancing fake time.
		assertCount(0)
		time.Sleep(interval - time.Nanosecond)
		assertCount(0)
		time.Sleep(time.Nanosecond)
		assertCount(1)

		time.Sleep(interval - time.Nanosecond)
		assertCount(1)
		time.Sleep(time.Nanosecond)
		assertCount(2)

		// The rejected job must not execute on a later tick.
		time.Sleep(interval)
		assertCount(2)

		for want := 1; want <= 2; want++ {
			if got := <-executed; got != want {
				t.Fatalf("execution order: got job %d, want job %d", got, want)
			}
		}
	})
}
