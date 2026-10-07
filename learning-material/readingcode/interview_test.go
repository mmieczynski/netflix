package readingcode

import (
	"container/heap"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestRegistryOwnsItsHistory(t *testing.T) {
	var registry Registry
	for _, id := range []string{"", "A", "B"} {
		if !registry.Accept(id) || registry.Accept(id) {
			t.Fatalf("ID %q was not accepted exactly once", id)
		}
	}
	history := registry.History()
	history[0] = "changed"
	if got := registry.History(); !reflect.DeepEqual(got, []string{"", "A", "B"}) {
		t.Fatalf("caller changed registry history: %v", got)
	}
}

func TestTimeStoreWriteAndQueryContract(t *testing.T) {
	store := make(TimeStore)
	if !store.Set("quality", "HD", 10) || !store.Set("quality", "UHD", 20) {
		t.Fatal("ordered writes failed")
	}
	if store.Set("quality", "obsolete", 15) {
		t.Fatal("out-of-order write was accepted")
	}
	if !store.Set("quality", "", 20) {
		t.Fatal("equal-time replacement failed")
	}
	for _, tc := range []struct {
		key   string
		at    int64
		want  string
		found bool
	}{
		{"quality", 9, "", false},
		{"quality", 10, "HD", true},
		{"quality", 15, "HD", true},
		{"quality", 20, "", true},
		{"quality", 30, "", true},
		{"missing", 30, "", false},
	} {
		if got, found := store.Get(tc.key, tc.at); got != tc.want || found != tc.found {
			t.Fatalf("Get(%q, %d) = %q, %v; want %q, %v", tc.key, tc.at, got, found, tc.want, tc.found)
		}
	}
}

func TestPlaylistChoicesAgainstAllSubsets(t *testing.T) {
	// Subset enumeration checks pruning independently of recursive choices.
	for _, durations := range [][]int{nil, {0}, {0, 0, 1}, {2, 3, 4}, {1, 1, 1, 1}} {
		for count := 0; count <= len(durations)+1; count++ {
			for budget := 0; budget <= 6; budget++ {
				want := make(map[string]bool)
				for mask := 0; mask < 1<<len(durations); mask++ {
					var indices []int
					total := 0
					for i, duration := range durations {
						if mask&(1<<i) != 0 {
							indices = append(indices, i)
							total += duration
						}
					}
					if len(indices) == count && total <= budget {
						want[fmt.Sprint(indices)] = true
					}
				}
				got := PlaylistChoices(durations, count, budget)
				if len(got) != len(want) {
					t.Fatalf("%v count=%d budget=%d: got %v, want %v", durations, count, budget, got, want)
				}
				for _, indices := range got {
					key := fmt.Sprint(indices)
					if !want[key] {
						t.Fatalf("unexpected or repeated choice %v", indices)
					}
					delete(want, key)
				}
			}
		}
	}
	choices := PlaylistChoices([]int{1, 1, 1}, 2, 2)
	choices[0][0] = 99
	if !reflect.DeepEqual(choices[1:], [][]int{{0, 2}, {1, 2}}) {
		t.Fatalf("completed choices share storage: %v", choices)
	}
}

func runTogether(n int, action func(int)) {
	var workers sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < n; i++ {
		workers.Add(1)
		go func(i int) {
			defer workers.Done()
			<-start
			action(i)
		}(i)
	}
	close(start)
	workers.Wait()
}

func TestConcurrentCacheAdmissionAndDeduplication(t *testing.T) {
	t.Run("cache", func(t *testing.T) {
		cache := NewSafeReadCache(32)
		runTogether(32, func(i int) {
			key := fmt.Sprint(i)
			cache.Put(key, i)
			if got, found := cache.Get(key); !found || got != i {
				t.Errorf("Get(%q) = %d, %v", key, got, found)
			}
		})
		if _, found := cache.Get("missing"); found {
			t.Fatal("missing cache key was found")
		}
	})
	t.Run("per-user limit", func(t *testing.T) {
		limiter := UserLimiter{Window: 10, Limit: 3}
		accepted := make(chan string, 64)
		runTogether(64, func(i int) {
			user := fmt.Sprint(i % 2)
			if limiter.Allow(user, 0) {
				accepted <- user
			}
		})
		close(accepted)
		counts := make(map[string]int)
		for user := range accepted {
			counts[user]++
		}
		if !reflect.DeepEqual(counts, map[string]int{"0": 3, "1": 3}) {
			t.Fatalf("concurrent admission violated per-user limits: %v", counts)
		}
	})
	t.Run("deduplication", func(t *testing.T) {
		deduper := Deduper{Horizon: 10}
		accepted := make(chan bool, 32)
		runTogether(32, func(_ int) {
			if deduper.Accept("event", 0) {
				accepted <- true
			}
		})
		if len(accepted) != 1 {
			t.Fatalf("accepted the same event %d times", len(accepted))
		}
		if !deduper.Accept("event", 10) {
			t.Fatal("event was not accepted at the expiry boundary")
		}
	})
}

func TestFixedWindowAndScheduledClaims(t *testing.T) {
	limiter := FixedLimiter{Window: 10, Limit: 2}
	for _, now := range []int64{9, 10} {
		if !limiter.Allow(now) || !limiter.Allow(now) || limiter.Allow(now) {
			t.Fatalf("fixed allowance failed at %d", now)
		}
	}
	due := ExpiryHeap{
		{Key: "A", Deadline: 5, Generation: 1},
		{Key: "cancelled", Deadline: 6, Generation: 2},
		{Key: "A", Deadline: 9, Generation: 3},
		{Key: "B", Deadline: 10, Generation: 4},
	}
	heap.Init(&due)
	current := map[string]uint64{"A": 3, "B": 4}
	if got := ReadyJobs(&due, current, 8); len(got) != 0 {
		t.Fatalf("stale or cancelled jobs were returned: %v", got)
	}
	if got := ReadyJobs(&due, current, 9); !reflect.DeepEqual(got, []string{"A"}) {
		t.Fatalf("expected A at its new deadline: %v", got)
	}
	if got := ReadyJobs(&due, current, 9); len(got) != 0 {
		t.Fatalf("job was claimed twice: %v", got)
	}
	if got := ReadyJobs(&due, current, 10); !reflect.DeepEqual(got, []string{"B"}) {
		t.Fatalf("future job was lost: %v", got)
	}
}

func TestRankingFiltersAndFillsGenreSlots(t *testing.T) {
	input := []CatalogCandidate{
		{ID: "A", Genre: "action", Score: 8, Eligible: true},
		{ID: "B", Genre: "action", Score: 9, Eligible: true},
		{ID: "A", Genre: "action", Score: 10, Eligible: true},
		{ID: "C", Genre: "drama", Score: 9, Eligible: true},
		{ID: "C", Genre: "drama", Score: 99, Eligible: false},
	}
	original := append([]CatalogCandidate(nil), input...)
	got := RankEligible(input, map[string]bool{"B": true}, 2)
	if !reflect.DeepEqual(got, []Candidate{{ID: "A", Score: 10}, {ID: "C", Score: 9}}) {
		t.Fatalf("filter/deduplicate/rank order: %v", got)
	}
	if !reflect.DeepEqual(input, original) {
		t.Fatal("ranking changed input")
	}
	distinct := []CatalogCandidate{
		{ID: "B", Genre: "action", Score: 10, Eligible: true},
		{ID: "A", Genre: "action", Score: 10, Eligible: true},
		{ID: "C", Genre: "drama", Score: -2, Eligible: true},
	}
	for _, k := range []int{2, 5} {
		if got := GenreWinners(distinct, k); !reflect.DeepEqual(got, []CatalogCandidate{distinct[1], distinct[2]}) {
			t.Fatalf("genre selection at k=%d: %v", k, got)
		}
	}
}

func TestTransactionBoundariesAndTombstones(t *testing.T) {
	var store TransactionStore
	if store.Commit() || store.Rollback() {
		t.Fatal("completed a nonexistent transaction")
	}
	store.Change("A", Edit{Value: "original"})
	store.Begin()
	store.Change("A", Edit{Deleted: true})
	store.Begin()
	store.Change("A", Edit{Value: ""})
	if !store.Commit() {
		t.Fatal("inner commit failed")
	}
	if value, found := store.Get("A"); !found || value != "" {
		t.Fatal("empty value failed to replace the parent tombstone")
	}
	if !store.Rollback() {
		t.Fatal("outer rollback failed")
	}
	if value, found := store.Get("A"); !found || value != "original" {
		t.Fatal("outer rollback did not restore the base")
	}
	store.Begin()
	store.Change("A", Edit{Deleted: true})
	if !store.Commit() {
		t.Fatal("base commit failed")
	}
	if _, found := store.Get("A"); found {
		t.Fatal("committed deletion did not reach the base")
	}
}

func TestFilesystemResolutionAndPublishSnapshot(t *testing.T) {
	file := &FSNode{IsFile: true, Content: "HD"}
	root := &FSNode{Children: map[string]*FSNode{"movie": file}}
	if got, found := Resolve(root, nil); !found || got != root {
		t.Fatal("empty path did not resolve root")
	}
	if got, found := Resolve(root, []string{"movie"}); !found || got != file {
		t.Fatal("file did not resolve")
	}
	for _, path := range [][]string{{"missing"}, {"movie", "child"}} {
		if _, found := Resolve(root, path); found {
			t.Fatalf("invalid path resolved: %v", path)
		}
	}
	bus := TopicBus{callbacks: map[string]map[int]func(string){"movies": {}}}
	var calls [2]int
	for i := range calls {
		i := i
		bus.callbacks["movies"][i] = func(message string) {
			if message != "new title" {
				t.Errorf("unexpected message %q", message)
			}
			calls[i]++
			// Either callback may run first; both remove the live topic.
			bus.mu.Lock()
			delete(bus.callbacks, "movies")
			bus.mu.Unlock()
		}
	}
	bus.Publish("movies", "new title")
	bus.Publish("movies", "new title")
	if calls != [2]int{1, 1} {
		t.Fatalf("snapshot membership changed during delivery: %v", calls)
	}
}

func TestSharedLoaderFailureAndIndependentKeys(t *testing.T) {
	var loader SharedLoader
	failure := errors.New("backend unavailable")
	_, err := loader.Do("A", func() (string, error) {
		value, err := loader.Do("B", func() (string, error) { return "UHD", nil })
		if err != nil || value != "UHD" {
			t.Errorf("independent key failed: %q, %v", value, err)
		}
		return "", failure
	})
	if err != failure {
		t.Fatalf("backend error was lost: %v", err)
	}
	value, err := loader.Do("A", func() (string, error) { return "HD", nil })
	if err != nil || value != "HD" {
		t.Fatalf("failed call prevented retry: %q, %v", value, err)
	}
}
