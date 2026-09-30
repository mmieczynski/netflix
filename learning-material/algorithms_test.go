package prep

import (
	"math/rand"
	"testing"
)

func TestTwoSum(t *testing.T) {
	for _, tc := range []struct {
		xs     []int
		target int
		found  bool
	}{
		{nil, 0, false}, {[]int{3}, 6, false}, {[]int{3, 3}, 6, true},
		{[]int{2, 7, 11, 15}, 9, true}, {[]int{-3, 4, 2}, 1, true},
	} {
		i, j, ok := TwoSum(tc.xs, tc.target)
		if ok != tc.found || (ok && (i == j || tc.xs[i]+tc.xs[j] != tc.target)) {
			t.Fatalf("TwoSum(%v,%d) = %d,%d,%v", tc.xs, tc.target, i, j, ok)
		}
	}
}

func TestLongestUnique(t *testing.T) {
	for s, want := range map[string]int{"": 0, "aaaa": 1, "abba": 2, "abcabcbb": 3, "界a界b": 3, "🙂é🙂": 2} {
		if got := LongestUnique(s); got != want {
			t.Fatalf("%q: got %d want %d", s, got, want)
		}
	}
}

func TestPrefixCountAgainstBruteForce(t *testing.T) {
	r := rand.New(rand.NewSource(7))
	for trial := 0; trial < 300; trial++ {
		xs := make([]int64, r.Intn(15))
		for i := range xs {
			xs[i] = int64(r.Intn(7) - 3)
		}
		k := int64(r.Intn(9) - 4)
		var want int64
		for i := range xs {
			var sum int64
			for j := i; j < len(xs); j++ {
				sum += xs[j]
				if sum == k {
					want++
				}
			}
		}
		if got := CountSubarrays(xs, k); got != want {
			t.Fatalf("%v target %d: %d != %d", xs, k, got, want)
		}
	}
}

func TestMaxNonAdjacentAgainstEnumeration(t *testing.T) {
	r := rand.New(rand.NewSource(9))
	for trial := 0; trial < 100; trial++ {
		xs := make([]int64, r.Intn(10))
		for i := range xs {
			xs[i] = int64(r.Intn(15) - 5)
		}
		var want int64
		for mask := 0; mask < 1<<len(xs); mask++ {
			if mask&(mask<<1) != 0 {
				continue
			}
			var total int64
			for i, x := range xs {
				if mask&(1<<i) != 0 {
					total += x
				}
			}
			if total > want {
				want = total
			}
		}
		if got := MaxNonAdjacent(xs); got != want {
			t.Fatalf("%v: %d != %d", xs, got, want)
		}
	}
}

func TestCourseOrder(t *testing.T) {
	edges := [][2]int{{0, 2}, {1, 2}, {2, 3}, {0, 2}}
	order, ok := CourseOrder(5, edges)
	if !ok || len(order) != 5 {
		t.Fatal(order, ok)
	}
	position := make(map[int]int)
	for i, v := range order {
		if _, exists := position[v]; exists {
			t.Fatal("duplicate vertex")
		}
		position[v] = i
	}
	for _, e := range edges {
		if position[e[0]] >= position[e[1]] {
			t.Fatal("prerequisite violated")
		}
	}
	for _, edges := range [][][2]int{{{0, 0}}, {{0, 1}, {1, 0}}, {{0, 8}}} {
		if _, ok := CourseOrder(3, edges); ok {
			t.Fatal("accepted cycle or invalid edge", edges)
		}
	}
	if order, ok := CourseOrder(0, nil); !ok || len(order) != 0 {
		t.Fatal("empty graph")
	}
}

func TestTimeMap(t *testing.T) {
	var m TimeMap
	if _, ok := m.Get("missing", 20); ok {
		t.Fatal("missing key")
	}
	if !m.Set("q", "1080p", 10) || !m.Set("q", "4k", 20) || m.Set("q", "old", 15) {
		t.Fatal("write policy")
	}
	if !m.Set("q", "8k", 20) {
		t.Fatal("replace equal timestamp")
	}
	for _, tc := range []struct {
		time  int64
		value string
		ok    bool
	}{{9, "", false}, {10, "1080p", true}, {15, "1080p", true}, {20, "8k", true}, {99, "8k", true}} {
		v, ok := m.Get("q", tc.time)
		if v != tc.value || ok != tc.ok {
			t.Fatalf("at %d got %q,%v", tc.time, v, ok)
		}
	}
}
