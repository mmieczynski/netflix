package limiter

import (
	"sync"
	"time"
)

// assumption/constraints:
// - single threaded access
// - max size set by the consumer
// - time based token regeneration (not by separate process)
// - no out of band token refill

// Allowed
// refill the bucket based on current time (up to max)
// check if bucket has tokens
// if there is, remove one and return yes
// if no tokens, return no

type BucketRateLimiter struct {
	maxTokens      int
	refillInterval time.Duration
	lastRefill     time.Time
	tokens         int
	mu             sync.Mutex
}

func NewBucketRateLimiter(maxTokens int, refillInterval time.Duration) BucketRateLimiter {
	return BucketRateLimiter{maxTokens: maxTokens, refillInterval: refillInterval, lastRefill: time.Now(), tokens: maxTokens}
}

func (r *BucketRateLimiter) Allowed() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	refill := int(time.Since(r.lastRefill) / r.refillInterval)
	if r.tokens < r.maxTokens {
		missingTokens := r.maxTokens - r.tokens
		if missingTokens <= refill {
			r.tokens += missingTokens
			r.lastRefill = time.Now()
		} else {
			r.tokens += refill
			r.lastRefill = r.lastRefill.Add(r.refillInterval * time.Duration(refill))
		}
	} else {
		r.lastRefill = time.Now()
	}

	if r.tokens <= 0 {
		return false
	}

	r.tokens--
	return true
}
