package limiter

import (
	"context"
	"time"
)

type LeakyLimiter struct {
	queue    chan func()
	interval time.Duration
}

func NewLeakyLimiter(interval time.Duration, capacity int) LeakyLimiter {
	if interval <= 0 && capacity <= 0 {
		panic("Wrong capacity or interval")
	}

	return LeakyLimiter{queue: make(chan func(), capacity), interval: interval}
}

func (l *LeakyLimiter) Submit(job func()) bool {
	select {
	case l.queue <- job:
		return true
	default:
		return false
	}
}

func (l LeakyLimiter) Run(ctx context.Context) {
	timer := time.NewTimer(l.interval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}
		if ctx.Err() != nil {
			return
		}

		started := time.Now()

		select {
		case job := <-l.queue:
			job()
		default:
		}

		timer.Reset(max(0, l.interval-time.Since(started)))
	}
}
