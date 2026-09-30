package readingcode

import (
	"reflect"
	"testing"
)

func TestLookupAndVersions(t *testing.T) {
	m := map[string]Movie{"A": {Title: "A"}}
	if got, ok := LookupMovie(m, "A"); !ok || got.Minutes != 0 {
		t.Fatal(got, ok)
	}
	if _, ok := LookupMovie(m, "B"); ok {
		t.Fatal("missing movie")
	}
	if _, ok := LookupMovie(nil, "A"); ok {
		t.Fatal("nil map")
	}
	versions := []Version{{3, "a"}, {7, "b"}, {12, ""}}
	for _, tc := range []struct {
		at    int64
		want  string
		found bool
	}{
		{2, "", false}, {3, "a", true}, {7, "b", true}, {9, "b", true}, {12, "", true}, {20, "", true},
	} {
		got, ok := VersionAt(versions, tc.at)
		if got != tc.want || ok != tc.found {
			t.Fatalf("at %d: %q %v", tc.at, got, ok)
		}
	}
	if _, ok := VersionAt(nil, 9); ok {
		t.Fatal("empty history")
	}
}

func TestWindowAgainstAllSubarrays(t *testing.T) {
	// Exhaustive four-symbol inputs, including repeated repair steps.
	for n := 0; n <= 6; n++ {
		combinations := 1
		for i := 0; i < n; i++ {
			combinations *= 4
		}
		for mask := 0; mask < combinations; mask++ {
			a := make([]string, n)
			code := mask
			for i := range a {
				a[i] = []string{"A", "B", "C", "D"}[code%4]
				code /= 4
			}
			want := 0
			for start := range a {
				for end := start; end < n; end++ {
					distinct := map[string]bool{}
					for _, v := range a[start : end+1] {
						distinct[v] = true
					}
					if len(distinct) <= 2 && end-start+1 > want {
						want = end - start + 1
					}
				}
			}
			if got := LongestTwo(a); got != want {
				t.Fatalf("%v: got %d want %d", a, got, want)
			}
		}
	}
}

func TestHopsAndRewards(t *testing.T) {
	edges := map[string][]string{"A": {"B", "C"}, "B": {"D"}, "C": {"D"}, "D": {"A"}}
	want := map[string]int{"A": 0, "B": 1, "C": 1, "D": 2}
	if got := HopDistances(edges, "A"); !reflect.DeepEqual(got, want) {
		t.Fatal(got)
	}
	if got := HopDistances(nil, "Z"); !reflect.DeepEqual(got, map[string]int{"Z": 0}) {
		t.Fatal(got)
	}
	for _, a := range [][]int64{nil, {6, 10, 6}, {-5, -1}, {2, 7, 9, 3, 1}, {0, 0, 0}, {8, -2, 1, 9}} {
		var best int64
		for mask := 0; mask < (1 << len(a)); mask++ {
			if mask&(mask<<1) != 0 {
				continue
			}
			var sum int64
			for i, v := range a {
				if mask&(1<<i) != 0 {
					sum += v
				}
			}
			if sum > best {
				best = sum
			}
		}
		if got := BestNonAdjacent(a); got != best {
			t.Fatal(a, got, best)
		}
	}
}

func TestCacheListInvariants(t *testing.T) {
	head, tail := &cacheNode{}, &cacheNode{}
	head.next = tail
	tail.prev = head
	a, b := &cacheNode{key: "A", value: 0}, &cacheNode{key: "B", value: 7}
	insertFront(head, a)
	insertFront(head, b)
	c := ReadCache{map[string]*cacheNode{"A": a, "B": b}, head, tail}
	for _, key := range []string{"A", "A", "B", "absent", "A"} {
		value, ok := c.Get(key)
		if key == "absent" {
			if ok {
				t.Fatal("miss was hit")
			}
		} else if !ok || value != c.byKey[key].value {
			t.Fatal(key, value, ok)
		}
		seen := map[*cacheNode]bool{}
		previous := head
		for n := head.next; n != tail; n = n.next {
			if n == nil || seen[n] || n.prev != previous || c.byKey[n.key] != n {
				t.Fatal("broken list")
			}
			seen[n] = true
			previous = n
		}
		if tail.prev != previous || len(seen) != len(c.byKey) {
			t.Fatal("membership mismatch")
		}
	}
	if head.next != a || a.next != b {
		t.Fatal("incorrect recency")
	}
}

func TestExpiryAndClaims(t *testing.T) {
	m := map[string]TimedValue{"A": {Value: "new", Deadline: 5, Generation: 2}}
	if ApplyExpiry(m, Expiry{"A", 5, 1}, 5) {
		t.Fatal("stale generation removed refresh")
	}
	if ApplyExpiry(m, Expiry{"A", 5, 2}, 4) {
		t.Fatal("early expiry")
	}
	if !ApplyExpiry(m, Expiry{"A", 5, 2}, 5) || len(m) != 0 {
		t.Fatal("exact boundary")
	}
	if ApplyExpiry(m, Expiry{"A", 5, 2}, 6) {
		t.Fatal("absent entry")
	}
	jobs := map[string]uint64{"A": 3}
	if ClaimJob(jobs, JobRecord{"A", 5, 1}, 5) {
		t.Fatal("old job")
	}
	if ClaimJob(jobs, JobRecord{"A", 9, 3}, 8) {
		t.Fatal("early job")
	}
	if !ClaimJob(jobs, JobRecord{"A", 9, 3}, 9) {
		t.Fatal("due job")
	}
	if ClaimJob(jobs, JobRecord{"A", 9, 3}, 10) {
		t.Fatal("double claim")
	}
}

func TestLimiterBoundaryAndRejectedAttempts(t *testing.T) {
	limiter := WindowLimiter{Window: 10, Limit: 3}
	times := []int64{0, 4, 9, 10, 10, 14, 14, 20}
	want := []bool{true, true, true, true, false, true, false, true}
	for i, now := range times {
		if got := limiter.Allow(now); got != want[i] {
			t.Fatal(i, now, got)
		}
	}
	zero := WindowLimiter{Window: 10, Limit: 0}
	if zero.Allow(0) || len(zero.accepted) != 0 {
		t.Fatal("zero limit")
	}
}

func TestRankingAndLayeredReads(t *testing.T) {
	input := []Candidate{{"D", 7}, {"A", 9}, {"C", 7}}
	copyInput := append([]Candidate(nil), input...)
	if got := TopTitles(input, 2); !reflect.DeepEqual(got, []Candidate{{"A", 9}, {"C", 7}}) {
		t.Fatal(got)
	}
	if !reflect.DeepEqual(input, copyInput) {
		t.Fatal("mutated input")
	}
	if TopTitles(input, 0) != nil || len(TopTitles(input, 10)) != 3 {
		t.Fatal("k boundary")
	}
	base := map[string]string{"A": "1"}
	outer := map[string]Edit{"A": {Value: "2"}}
	inner := map[string]Edit{"A": {Deleted: true}}
	if _, ok := ReadLayered(base, []map[string]Edit{outer, inner}, "A"); ok {
		t.Fatal("ignored tombstone")
	}
	if v, ok := ReadLayered(base, []map[string]Edit{outer}, "A"); !ok || v != "2" {
		t.Fatal(v, ok)
	}
	inner["A"] = Edit{Value: ""}
	if v, ok := ReadLayered(base, []map[string]Edit{outer, inner}, "A"); !ok || v != "" {
		t.Fatal("empty value lost")
	}
	if v, ok := ReadLayered(base, nil, "A"); !ok || v != "1" {
		t.Fatal(v, ok)
	}
}

func TestRecentTotals(t *testing.T) {
	events := []WatchEvent{{"e3", "A", 9, 3}, {"e2", "B", 8, 6}, {"e1", "A", 2, 4}, {"e2", "B", 8, 6}, {"future", "C", 13, 100}}
	got := RecentTotals(events, 12, 10)
	if !reflect.DeepEqual(got, map[string]int64{"A": 3, "B": 6}) {
		t.Fatal(got)
	}
	if len(RecentTotals(nil, 12, 10)) != 0 {
		t.Fatal("empty batch")
	}
}
