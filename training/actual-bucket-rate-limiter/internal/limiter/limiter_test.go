package limiter

import (
	"testing"
	"testing/synctest"
	"time"
)

func TestAllowedAfterRefill(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		limiter := NewBucketRateLimiter(1, 10*time.Second)

		if allowed := limiter.Allowed(); !allowed {
			t.Fatal("Should be allowed")
		}

		time.Sleep(10 * time.Second)

		if allowed := limiter.Allowed(); !allowed {
			t.Fatal("Should be allowed")
		}
	})
}
