package prep

import "sort"

// TwoSum returns distinct indices. Inputs and target must fit ordinary int arithmetic.
func TwoSum(nums []int, target int) (int, int, bool) {
	seen := make(map[int]int)
	for i, n := range nums {
		if j, ok := seen[target-n]; ok {
			return j, i, true
		}
		seen[n] = i
	}
	return 0, 0, false
}

// LongestUnique measures Unicode code points, not grapheme clusters or bytes.
func LongestUnique(s string) int {
	last := make(map[rune]int)
	left, best, pos := 0, 0, 0
	for _, ch := range s {
		if prev, ok := last[ch]; ok && prev >= left {
			left = prev + 1
		}
		last[ch] = pos
		if length := pos - left + 1; length > best {
			best = length
		}
		pos++
	}
	return best
}

// CountSubarrays assumes sums and result fit int64.
func CountSubarrays(nums []int64, target int64) int64 {
	counts := map[int64]int64{0: 1}
	var prefix, answer int64
	for _, x := range nums {
		prefix += x
		answer += counts[prefix-target]
		counts[prefix]++
	}
	return answer
}

// MaxNonAdjacent allows an empty selection; totals must fit int64.
func MaxNonAdjacent(values []int64) int64 {
	var next, afterNext int64
	for i := len(values) - 1; i >= 0; i-- {
		best := next
		if take := values[i] + afterNext; take > best {
			best = take
		}
		afterNext, next = next, best
	}
	return next
}

// CourseOrder edges are [prerequisite, dependent]. Invalid endpoints or cycles
// return false. Duplicate edges are counted consistently. Order need not be unique.
func CourseOrder(n int, edges [][2]int) ([]int, bool) {
	if n < 0 {
		return nil, false
	}
	adj := make([][]int, n)
	degree := make([]int, n)
	for _, e := range edges {
		if e[0] < 0 || e[0] >= n || e[1] < 0 || e[1] >= n {
			return nil, false
		}
		adj[e[0]] = append(adj[e[0]], e[1])
		degree[e[1]]++
	}
	ready := make([]int, 0, n)
	for v, d := range degree {
		if d == 0 {
			ready = append(ready, v)
		}
	}
	for head := 0; head < len(ready); head++ {
		for _, v := range adj[ready[head]] {
			degree[v]--
			if degree[v] == 0 {
				ready = append(ready, v)
			}
		}
	}
	if len(ready) != n {
		return nil, false
	}
	return ready, true
}

type version struct {
	Time  int64
	Value string
}

// TimeMap is single-threaded. Writes must be nondecreasing per key.
// Its zero value is ready to use. Equal timestamps replace the latest value.
type TimeMap struct{ values map[string][]version }

func (m *TimeMap) Set(key, value string, timestamp int64) bool {
	if m.values == nil {
		m.values = make(map[string][]version)
	}
	vs := m.values[key]
	if len(vs) > 0 {
		last := &vs[len(vs)-1]
		if timestamp < last.Time {
			return false
		}
		if timestamp == last.Time {
			last.Value = value
			return true
		}
	}
	m.values[key] = append(vs, version{timestamp, value})
	return true
}

func (m *TimeMap) Get(key string, timestamp int64) (string, bool) {
	vs := m.values[key]
	i := sort.Search(len(vs), func(i int) bool { return vs[i].Time > timestamp })
	if i == 0 {
		return "", false
	}
	return vs[i-1].Value, true
}
