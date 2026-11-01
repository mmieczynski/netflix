package ratelimiter

import (
	"testing"
	"time"
)

// single call with empty limiter
// couple of calls
// calls exeeding the limit
// empty user input

func TestAllowSingleCall(t *testing.T) {
	now := time.Unix(0, 0)
	nowFunc := func() time.Time {
		return now
	}

	limiter := NewBucketRateLimiter(3, 10*time.Second, nowFunc)

	allowed, err := limiter.Allowed("foo")
	if err != nil {
		t.Fatalf("Error in the function: %v", err)
	}
	if !allowed {
		t.Fatal("Should allow but blocked")
	}
}

func TestAllowMultipleCalls(t *testing.T) {
	now := time.Unix(0, 0)
	nowFunc := func() time.Time {
		return now
	}

	limiter := NewBucketRateLimiter(3, 10*time.Second, nowFunc)

	allowed, err := limiter.Allowed("foo")
	if err != nil {
		t.Fatalf("Error in the function: %v", err)
	}
	if !allowed {
		t.Fatal("Should allow but blocked")
	}

	now = time.Unix(3, 0)
	allowed, err = limiter.Allowed("foo")
	if err != nil {
		t.Fatalf("Error in the function: %v", err)
	}
	if !allowed {
		t.Fatal("Should allow but blocked")
	}

	now = time.Unix(5, 0)
	allowed, err = limiter.Allowed("foo")
	if err != nil {
		t.Fatalf("Error in the function: %v", err)
	}
	if !allowed {
		t.Fatal("Should allow but blocked")
	}
}

func TestBlockAfterExceedingLimit(t *testing.T) {
	now := time.Unix(0, 0)
	nowFunc := func() time.Time {
		return now
	}

	limiter := NewBucketRateLimiter(3, 10*time.Second, nowFunc)

	limiter.Allowed("foo")

	now = time.Unix(3, 0)
	limiter.Allowed("foo")

	now = time.Unix(5, 0)
	limiter.Allowed("foo")

	now = time.Unix(9, 0)
	allowed, err := limiter.Allowed("foo")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("Should not be allowed")
	}
}
