package readingcode

import (
	"reflect"
	"sort"
	"testing"
)

func TestSignedPrefixCountsAndProductsAgainstBruteForce(t *testing.T) {
	// Distinct starting boundaries must count separately, including k=0.
	for n := 0; n <= 6; n++ {
		combinations := 1
		for i := 0; i < n; i++ {
			combinations *= 3
		}
		for mask := 0; mask < combinations; mask++ {
			a := make([]int, n)
			b := make([]int64, n)
			code := mask
			for i := range a {
				a[i] = code%3 - 1
				b[i] = int64(a[i])
				code /= 3
			}
			for target := -3; target <= 3; target++ {
				want := 0
				for i := range a {
					sum := 0
					for j := i; j < n; j++ {
						sum += a[j]
						if sum == target {
							want++
						}
					}
				}
				if got := SubarrayCount(a, target); got != want {
					t.Fatalf("%v k=%d: %d != %d", a, target, got, want)
				}
			}
			got := ProductsExceptSelf(b)
			for i := range b {
				want := int64(1)
				for j, value := range b {
					if j != i {
						want *= value
					}
				}
				if got[i] != want {
					t.Fatal(b, i, got, want)
				}
			}
		}
	}
}

func TestArrayAndStringContracts(t *testing.T) {
	if pair, ok := PairSum([]int{3, 3}, 6); !ok || pair != [2]int{0, 1} {
		t.Fatal(pair, ok)
	}
	if _, ok := PairSum([]int{3}, 6); ok {
		t.Fatal("reused same position")
	}
	if UniqueLength("ABBA") != 2 || UniqueLength("éAéB") != 3 {
		t.Fatal("rune boundaries")
	}
	groups := AnagramGroups([]string{"eat", "tea", "abb", "ab"})
	sizes := []int{}
	for _, group := range groups {
		sizes = append(sizes, len(group))
	}
	sort.Ints(sizes)
	if !reflect.DeepEqual(sizes, []int{1, 1, 2}) {
		t.Fatal(groups)
	}
	input := []int{-1, 0, 1, 2, -1, -4}
	original := append([]int(nil), input...)
	want := [][3]int{{-1, -1, 2}, {-1, 0, 1}}
	if got := ThreeZero(input); !reflect.DeepEqual(got, want) {
		t.Fatal(got)
	}
	if !reflect.DeepEqual(input, original) {
		t.Fatal("3Sum mutated input")
	}
	if got := MergeCoverage([][2]int{{2, 5}, {1, 3}, {7, 9}, {5, 6}}); !reflect.DeepEqual(got, [][2]int{{1, 6}, {7, 9}}) {
		t.Fatal(got)
	}
	a := []int{4, 5, 6, 7, 0, 1, 2}
	for i, value := range a {
		if got := RotatedIndex(a, value); got != i {
			t.Fatal(value, got)
		}
	}
	if RotatedIndex(a, 3) != -1 || RotatedIndex(nil, 1) != -1 {
		t.Fatal("missing rotated target")
	}
	for _, values := range [][]int{nil, {4, 9, 2, 7, 5}, {5, 5, 2}, {-3, -1, -2}} {
		sorted := append([]int(nil), values...)
		sort.Sort(sort.Reverse(sort.IntSlice(sorted)))
		for k := 0; k <= len(values)+2; k++ {
			got := LargestK(values, k)
			end := k
			if end > len(sorted) {
				end = len(sorted)
			}
			if len(got) != end {
				t.Fatal(values, k, got)
			}
			for i := range got {
				if got[i] != sorted[i] {
					t.Fatal(values, k, got)
				}
			}
		}
	}
}

func TestGraphAndChoiceExamples(t *testing.T) {
	tree := &TreeNode{Value: 10, Left: &TreeNode{Value: 5, Right: &TreeNode{Value: 12}}, Right: &TreeNode{Value: 15}}
	if IsBST(tree) {
		t.Fatal("ignored ancestor bound")
	}
	if got := TreeLevels(tree); !reflect.DeepEqual(got, [][]int{{10}, {5, 15}, {12}}) {
		t.Fatal(got)
	}
	grid := [][]byte{[]byte("110"), []byte("010"), []byte("001")}
	if got := IslandCount(grid); got != 2 {
		t.Fatal(got)
	}
	order, ok := JobOrder(5, [][2]int{{0, 1}, {0, 2}, {1, 3}, {2, 3}})
	if !ok || len(order) != 5 {
		t.Fatal(order, ok)
	}
	position := make(map[int]int)
	for i, id := range order {
		position[id] = i
	}
	for _, edge := range [][2]int{{0, 1}, {0, 2}, {1, 3}, {2, 3}} {
		if position[edge[0]] >= position[edge[1]] {
			t.Fatal(order)
		}
	}
	if _, ok := JobOrder(2, [][2]int{{0, 1}, {1, 0}}); ok {
		t.Fatal("cycle")
	}
	a, b := &GraphNode{Label: "same"}, &GraphNode{Label: "same"}
	a.Neighbors, b.Neighbors = []*GraphNode{b}, []*GraphNode{a}
	copy := CloneGraph(a)
	if copy == a || copy.Neighbors[0] == b || copy.Neighbors[0] == copy || copy.Neighbors[0].Neighbors[0] != copy {
		t.Fatal("clone identity or cycle")
	}
	if !Balanced("([])") || Balanced("([)]") || Balanced(")") || !Balanced("") {
		t.Fatal("brackets")
	}
	if got := GreaterWait([]int{70, 73, 71, 74}); !reflect.DeepEqual(got, []int{1, 2, 1, 0}) {
		t.Fatal(got)
	}
	if got := PlaylistChoices([]int{2, 3, 4}, 2, 5); !reflect.DeepEqual(got, [][]int{{0, 1}}) {
		t.Fatal(got)
	}
	if got := MostScreenings([][2]int{{0, 10}, {1, 2}, {2, 3}, {3, 4}}); len(got) != 3 {
		t.Fatal(got)
	}
}

func TestFullCacheContracts(t *testing.T) {
	c := NewReadCache(2)
	c.Put("A", 4)
	c.Put("B", 7)
	c.Get("A")
	c.Put("C", 9)
	c.Put("A", 0)
	if _, found := c.Get("B"); found {
		t.Fatal("wrong LRU victim")
	}
	if v, found := c.Get("A"); !found || v != 0 || len(c.byKey) != 2 {
		t.Fatal(v, found)
	}
	w := NewWeightedCache(10)
	w.Put("A", 4, 4)
	w.Put("B", 3, 3)
	w.Put("C", 2, 2)
	if w.Put("A", 11, 11) || w.used != 9 || w.list.head.next.key != "C" {
		t.Fatal("invalid overwrite mutated state")
	}
	if !w.Put("A", 9, 9) || w.used != 9 || len(w.weights) != 1 {
		t.Fatal("repeated eviction")
	}
	if v, ok := w.Get("A"); !ok || v != 9 {
		t.Fatal(v, ok)
	}
	now := int64(0)
	ttl := NewTTLMap(func() int64 { return now })
	ttl.Put("A", "old", 5)
	now = 3
	ttl.Put("A", "new", 7)
	now = 5
	ttl.Cleanup()
	if v, ok := ttl.Get("A"); !ok || v != "new" {
		t.Fatal("stale timer removed refresh")
	}
	now = 10
	if _, ok := ttl.Get("A"); ok {
		t.Fatal("exact deadline was live")
	}
	ttl.Cleanup()
	if len(ttl.due) != 0 {
		t.Fatal("due record retained")
	}
	now = 0
	m := NewMetadataCache(2, func() int64 { return now })
	m.Put("A", "", 5)
	m.Put("B", "UHD", 20)
	if v, ok := m.Get("A"); !ok || v != "" {
		t.Fatal(v, ok)
	}
	now = 5
	m.Put("C", "HDR", 20)
	if _, ok := m.Get("A"); ok {
		t.Fatal("expired MRU survived")
	}
	if _, ok := m.Get("B"); !ok {
		t.Fatal("evicted live LRU before expired MRU")
	}
	now = 6
	m.Put("C", "new", 50)
	now = 25
	if v, ok := m.Get("C"); !ok || v != "new" {
		t.Fatal(v, ok)
	}
}

func TestBucketsAndServiceExamples(t *testing.T) {
	bucket := NewTokenBucket(3, 1, 0)
	for i := 0; i < 3; i++ {
		if !bucket.Allow(0) {
			t.Fatal(i)
		}
	}
	if bucket.Allow(0) || bucket.Allow(.5) || !bucket.Allow(1) {
		t.Fatal("fractional refill counted twice")
	}
	d := Deduper{Horizon: 600}
	if !d.Accept("e1", 0) || d.Accept("e1", 599) || !d.Accept("e1", 600) {
		t.Fatal("dedup horizon")
	}
	counter := HitCounter{Window: 60}
	for _, at := range []int64{100, 101, 101, 150} {
		counter.Record(at)
	}
	if counter.Count(150) != 4 || counter.Count(160) != 3 {
		t.Fatal("counter boundary")
	}
	scores := CollaborativeScores([]string{"A", "B"}, [][]string{{"A", "C", "D", "C"}, {"B", "C"}, {"E"}})
	if !reflect.DeepEqual(scores, map[string]int64{"C": 2, "D": 1}) {
		t.Fatal(scores)
	}
	stats := RollingStats{Window: 10}
	stats.Record(StatEvent{0, "A", 10})
	stats.Record(StatEvent{9, "A", 2})
	if v, ok := stats.Average("A", 9); !ok || v != 6 {
		t.Fatal(v, ok)
	}
	if v, ok := stats.Average("A", 10); !ok || v != 2 {
		t.Fatal(v, ok)
	}
	if _, ok := stats.Average("A", 19); ok {
		t.Fatal("empty average")
	}
	s := TransactionStore{}
	s.Change("A", Edit{Value: "1"})
	s.Begin()
	s.Change("A", Edit{Value: "2"})
	s.Begin()
	s.Change("A", Edit{Deleted: true})
	s.Commit()
	if _, ok := s.Get("A"); ok {
		t.Fatal("inner tombstone lost")
	}
	s.Rollback()
	if v, ok := s.Get("A"); !ok || v != "1" {
		t.Fatal(v, ok)
	}
	events := []WatchEvent{{"e1", "A", 6, 3}, {"e2", "A", 9, 4}, {"e2", "A", 9, 4}, {"e3", "B", 10, 7}, {"e4", "C", 5, 100}, {"e5", "D", 11, 20}}
	got, valid := RecentRanking(events, 10, 5, 2)
	if !valid || !reflect.DeepEqual(got, []Candidate{{"A", 7}, {"B", 7}}) {
		t.Fatal(got, valid)
	}
	if _, valid := RecentRanking([]WatchEvent{{"bad", "A", 10, 0}}, 10, 5, 2); valid {
		t.Fatal("invalid duration")
	}
}

func TestPendingLoadUsesPublishedResult(t *testing.T) {
	// Pre-register a pending leader to deterministically exercise the join path.
	call := &loadCall{done: make(chan struct{})}
	loader := SharedLoader{active: map[string]*loadCall{"A": call}}
	result := make(chan string, 1)
	go func() {
		value, err := loader.Do("A", func() (string, error) {
			panic("a follower must not load")
		})
		if err != nil {
			value = "unexpected error"
		}
		result <- value
	}()
	call.value = "HD"
	close(call.done)
	if got := <-result; got != "HD" {
		t.Fatal(got)
	}
	leader := SharedLoader{}
	value, err := leader.Do("B", func() (string, error) { return "UHD", nil })
	if err != nil || value != "UHD" || len(leader.active) != 0 {
		t.Fatal(value, err)
	}
}
