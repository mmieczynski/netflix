package prep

import (
	"container/heap"
	"math"
	"sort"
	"sync"
)

type timeQueue struct {
	times []int64
	head  int
}

// SlidingLimiter counts accepted requests in (now-window, now].
// Use server time with nondecreasing ticks across calls; differences must fit int64.
// Per-user idle state is retained until that user next makes a request.
type SlidingLimiter struct {
	mu     sync.Mutex
	limit  int
	window int64
	users  map[string]*timeQueue
}

func NewSlidingLimiter(limit int, window int64) *SlidingLimiter {
	if limit <= 0 || window <= 0 {
		panic("limit and window must be positive")
	}
	return &SlidingLimiter{limit: limit, window: window, users: make(map[string]*timeQueue)}
}

func (l *SlidingLimiter) Allow(user string, now int64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	q := l.users[user]
	if q == nil {
		q = &timeQueue{}
		l.users[user] = q
	}
	for q.head < len(q.times) && q.times[q.head] <= now-l.window {
		q.head++
	}
	// Compact after enough removals; retained storage stays proportional to limit.
	if q.head > 0 && q.head >= len(q.times)-q.head {
		copy(q.times, q.times[q.head:])
		q.times = q.times[:len(q.times)-q.head]
		q.head = 0
	}
	if len(q.times)-q.head >= l.limit {
		return false
	}
	q.times = append(q.times, now)
	return true
}

// TokenBucket is one bucket (wrap in a per-user map for multiple users).
// Times are finite seconds; rate is tokens per second. Float arithmetic is approximate.
type TokenBucket struct {
	mu                           sync.Mutex
	capacity, rate, tokens, last float64
}

func NewTokenBucket(capacity, rate, now float64) *TokenBucket {
	if !finite(capacity) || !finite(rate) || !finite(now) || capacity < 1 || rate < 0 {
		panic("invalid token bucket parameters")
	}
	return &TokenBucket{capacity: capacity, rate: rate, tokens: capacity, last: now}
}

func finite(x float64) bool { return !math.IsNaN(x) && !math.IsInf(x, 0) }

func (b *TokenBucket) Allow(now float64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !finite(now) {
		return false
	}
	if now < b.last {
		now = b.last
	}
	if b.rate > 0 {
		b.tokens = math.Min(b.capacity, b.tokens+(now-b.last)*b.rate)
	}
	b.last = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

type Candidate struct {
	ID       string
	Genre    string
	Score    int64
	Eligible bool
}

func better(a, b Candidate) bool {
	if a.Score != b.Score {
		return a.Score > b.Score
	}
	return a.ID < b.ID
}

type candidateHeap []Candidate

func (h candidateHeap) Len() int { return len(h) }

// Worst retained candidate first: invert the final output order.
func (h candidateHeap) Less(i, j int) bool  { return better(h[j], h[i]) }
func (h candidateHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *candidateHeap) Push(x interface{}) { *h = append(*h, x.(Candidate)) }
func (h *candidateHeap) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	old[len(old)-1] = Candidate{}
	*h = old[:len(old)-1]
	return x
}

// Recommend filters before deduplicating. Equal scores for one ID use smaller genre.
// Output: descending score then ascending ID. Does not mutate inputs.
func Recommend(input []Candidate, watched map[string]bool, k int) []Candidate {
	if k <= 0 {
		return []Candidate{}
	}
	unique := make(map[string]Candidate)
	for _, c := range input {
		if !c.Eligible || watched[c.ID] {
			continue
		}
		old, ok := unique[c.ID]
		if !ok || c.Score > old.Score || (c.Score == old.Score && c.Genre < old.Genre) {
			unique[c.ID] = c
		}
	}
	h := candidateHeap{}
	for _, c := range unique {
		if len(h) < k {
			heap.Push(&h, c)
		} else if better(c, h[0]) {
			h[0] = c
			heap.Fix(&h, 0)
		}
	}
	sort.Slice(h, func(i, j int) bool { return better(h[i], h[j]) })
	return []Candidate(h)
}
