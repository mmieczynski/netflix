package ratelimiter

import (
	"errors"
	"strings"
	"time"
)

// assumtions/contraints:
// - window needs to be bigger than 0
// - single threaded
// - single entry point updating the count and answering if the request should be blocked
// - no way to update the counts out of band
// - multiple different users have different counts
// - we do not count rejected requests into limits

// operations:
// - check if blocked
// - update count for a given user
// - remove stale entries outside of window

// [user] - [3, 4, 9] -> now == 10 -> blocked
// [user] - [3, 4, 9] -> now == 15 -> remove all outside of window (remove 3, 4) -> [9, 15] -> allowed

type BucketRateLimiter struct {
	users  map[string][]time.Time
	limit  int
	window time.Duration
	now    func() time.Time
}

func NewBucketRateLimiter(limit int, window time.Duration, now func() time.Time) BucketRateLimiter {
	return BucketRateLimiter{
		limit:  limit,
		window: window,
		users:  map[string][]time.Time{},
		now:    now,
	}
}

func (r BucketRateLimiter) Allowed(user string) (allowed bool, err error) {
	if strings.TrimSpace(user) == "" {
		return false, errors.New("User should not be an empty string")
	}

	limits := r.users[user]
	for len(limits) > 0 && time.Since(limits[0]) >= r.window {
		limits = limits[1:]
	}

	if len(limits) >= r.limit {
		return false, nil
	}

	r.users[user] = append(limits, r.now())
	return true, nil
}
