package prep

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
)

func TestSlidingBoundary(t *testing.T) {
	l := NewSlidingLimiter(2, 10)
	for _, tc := range []struct {
		user string
		now  int64
		want bool
	}{
		{"a", 0, true}, {"a", 1, true}, {"a", 9, false}, {"b", 9, true}, {"a", 10, true}, {"a", 10, false}, {"a", 11, true}, {"a", 100, true},
	} {
		if got := l.Allow(tc.user, tc.now); got != tc.want {
			t.Fatalf("%+v got %v", tc, got)
		}
	}
}

func TestSlidingAgainstFullHistory(t *testing.T) {
	l := NewSlidingLimiter(3, 7)
	r := rand.New(rand.NewSource(17))
	history := map[string][]int64{}
	var now int64
	for i := 0; i < 1000; i++ {
		now += int64(r.Intn(3))
		user := fmt.Sprint(r.Intn(4))
		count := 0
		for _, ts := range history[user] {
			if ts > now-7 {
				count++
			}
		}
		want := count < 3
		if got := l.Allow(user, now); got != want {
			t.Fatalf("at %d user %s: got %v want %v", now, user, got, want)
		}
		if want {
			history[user] = append(history[user], now)
		}
	}
}

func TestTokenBucket(t *testing.T) {
	b := NewTokenBucket(3, 1, 0)
	for _, tc := range []struct {
		now  float64
		want bool
	}{{0, true}, {0, true}, {0, true}, {0, false}, {0.5, false}, {1, true}, {1, false}, {100, true}, {100, true}, {100, true}, {100, false}, {99, false}} {
		if got := b.Allow(tc.now); got != tc.want {
			t.Fatalf("at %v got %v want %v", tc.now, got, tc.want)
		}
	}
	if b.Allow(math.NaN()) {
		t.Fatal("NaN accepted")
	}
}

func TestConcurrentAdmission(t *testing.T) {
	l := NewSlidingLimiter(10, 60)
	b := NewTokenBucket(10, 0, 0)
	var acceptedLog, acceptedBucket int32
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if l.Allow("one-user", 0) {
				atomic.AddInt32(&acceptedLog, 1)
			}
			if b.Allow(0) {
				atomic.AddInt32(&acceptedBucket, 1)
			}
		}()
	}
	wg.Wait()
	if acceptedLog != 10 || acceptedBucket != 10 {
		t.Fatal("non-atomic admission", acceptedLog, acceptedBucket)
	}
}

func TestRecommendationAgainstFullSort(t *testing.T) {
	r := rand.New(rand.NewSource(23))
	for trial := 0; trial < 200; trial++ {
		input := make([]Candidate, r.Intn(50))
		for i := range input {
			input[i] = Candidate{fmt.Sprint(r.Intn(12)), fmt.Sprint(r.Intn(3)), int64(r.Intn(9) - 4), r.Intn(4) != 0}
		}
		watched := map[string]bool{"2": true, "5": true}
		k := r.Intn(16)
		// Independent baseline: sort all eligible observations, then retain first per ID.
		eligible := []Candidate{}
		for _, c := range input {
			if c.Eligible && !watched[c.ID] {
				eligible = append(eligible, c)
			}
		}
		sort.Slice(eligible, func(i, j int) bool {
			a, b := eligible[i], eligible[j]
			if a.Score != b.Score {
				return a.Score > b.Score
			}
			if a.ID != b.ID {
				return a.ID < b.ID
			}
			return a.Genre < b.Genre
		})
		want := []Candidate{}
		seen := map[string]bool{}
		for _, c := range eligible {
			if !seen[c.ID] && len(want) < k {
				seen[c.ID] = true
				want = append(want, c)
			}
		}
		got := Recommend(input, watched, k)
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Fatalf("got %v want %v", got, want)
		}
	}
}
