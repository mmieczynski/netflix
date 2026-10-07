// Code generated from reading Markdown; DO NOT EDIT.
package readingcode

import (
	"container/heap"
	"sort"
	"sync"
)

// Source: reading/01-go.md

type Movie struct {
	Title   string
	Minutes int
}

func LookupMovie(movies map[string]Movie,
	id string) (Movie, bool) {
	movie, found := movies[id]
	return movie, found
}

func CopyIDs(ids []string) []string {
	// Give the caller independent element storage.
	result := make([]string, len(ids))
	copy(result, ids)
	return result
}

type Registry struct {
	seen  map[string]bool
	order []string
}

func (r *Registry) Accept(id string) bool {
	if r.seen[id] {
		return false
	}
	if r.seen == nil {
		// Initialize on first write; the zero value works.
		r.seen = make(map[string]bool)
	}
	r.seen[id] = true
	r.order = append(r.order, id)
	return true
}

func (r *Registry) History() []string {
	return CopyIDs(r.order)
}

// Source: reading/02-patterns.md

func PairSum(nums []int, target int) ([2]int, bool) {
	earlier := make(map[int]int)
	for i, value := range nums {
		// Search earlier positions before storing this one.
		if j, found := earlier[target-value]; found {
			return [2]int{j, i}, true
		}
		earlier[value] = i
	}
	return [2]int{}, false
}

func AnagramGroups(words []string) [][]string {
	groups := make(map[[26]int][]string)
	for _, word := range words {
		// Counts preserve each letter's multiplicity.
		var key [26]int
		for i := 0; i < len(word); i++ {
			key[word[i]-'a']++
		}
		groups[key] = append(groups[key], word)
	}
	result := make([][]string, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}
	return result
}

func UniqueLength(text string) int {
	runes := []rune(text)
	last := make(map[rune]int)
	left, best := 0, 0
	for right, r := range runes {
		// An old repeat must not move the window backward.
		if previous, found := last[r]; found &&
			previous >= left {
			left = previous + 1
		}
		last[r] = right
		if length := right - left + 1; length > best {
			best = length
		}
	}
	return best
}

func SubarrayCount(nums []int, k int) int {
	// The empty prefix lets a subarray start at index zero.
	frequency := map[int]int{0: 1}
	prefix, result := 0, 0
	for _, value := range nums {
		prefix += value
		// Count earlier boundaries before recording this one.
		result += frequency[prefix-k]
		frequency[prefix]++
	}
	return result
}

func ProductsExceptSelf(nums []int64) []int64 {
	out := make([]int64, len(nums))
	left := int64(1)
	for i, value := range nums {
		// Write the product before including this element.
		out[i] = left
		left *= value
	}
	right := int64(1)
	for i := len(nums) - 1; i >= 0; i-- {
		// Combine products strictly on either side of i.
		out[i] *= right
		right *= nums[i]
	}
	return out
}

func LongestTwo(titles []string) int {
	count := make(map[string]int)
	left, best := 0, 0
	for right, title := range titles {
		count[title]++
		for len(count) > 2 {
			old := titles[left]
			count[old]--
			if count[old] == 0 {
				// Only live keys count toward the limit.
				delete(count, old)
			}
			left++
		}
		if length := right - left + 1; length > best {
			best = length
		}
	}
	return best
}

// Source: reading/03-order.md

func ThreeZero(nums []int) [][3]int {
	// Sorting a copy preserves the caller's input order.
	a := append([]int(nil), nums...)
	sort.Ints(a)
	var out [][3]int
	for i := 0; i+2 < len(a); i++ {
		if i > 0 && a[i] == a[i-1] {
			continue
		}
		left, right := i+1, len(a)-1
		for left < right {
			sum := a[i] + a[left] + a[right]
			if sum < 0 {
				left++
			} else if sum > 0 {
				right--
			} else {
				out = append(out,
					[3]int{a[i], a[left], a[right]})
				// Skip values that would repeat this triple.
				x, y := a[left], a[right]
				for left < right && a[left] == x {
					left++
				}
				for left < right && a[right] == y {
					right--
				}
			}
		}
	}
	return out
}

func MergeCoverage(input [][2]int) [][2]int {
	a := append([][2]int(nil), input...)
	sort.Slice(a, func(i, j int) bool {
		return a[i][0] < a[j][0]
	})
	var out [][2]int
	for _, interval := range a {
		last := len(out) - 1
		if last < 0 || interval[0] > out[last][1] {
			out = append(out, interval)
		} else if interval[1] > out[last][1] {
			// A nested interval must not shorten coverage.
			out[last][1] = interval[1]
		}
	}
	return out
}

type Version struct {
	At    int64
	Value string
}

func VersionAt(v []Version, at int64) (string, bool) {
	// Find the first version strictly after the query.
	lo, hi := 0, len(v)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if v[mid].At <= at {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	if lo == 0 {
		return "", false
	}
	// Its predecessor is the latest eligible version.
	return v[lo-1].Value, true
}

func RotatedIndex(a []int, target int) int {
	lo, hi := 0, len(a)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if a[mid] == target {
			return mid
		}
		if a[lo] <= a[mid] {
			// The left half is sorted, so test its bounds.
			if a[lo] <= target && target < a[mid] {
				hi = mid - 1
			} else {
				lo = mid + 1
			}
		} else {
			if a[mid] < target && target <= a[hi] {
				lo = mid + 1
			} else {
				hi = mid - 1
			}
		}
	}
	return -1
}

func PushMin(h []int, value int) []int {
	h = append(h, value)
	for i := len(h) - 1; i > 0; {
		parent := (i - 1) / 2
		if h[parent] <= h[i] {
			break
		}
		h[parent], h[i] = h[i], h[parent]
		i = parent
	}
	return h
}

func RepairMinRoot(h []int) {
	for i := 0; ; {
		child := 2*i + 1
		if child >= len(h) {
			return
		}
		if child+1 < len(h) && h[child+1] < h[child] {
			// Repair through the smaller child.
			child++
		}
		if h[i] <= h[child] {
			return
		}
		h[i], h[child] = h[child], h[i]
		i = child
	}
}

func LargestK(values []int, k int) []int {
	if k <= 0 {
		return nil
	}
	var h []int
	for _, value := range values {
		if len(h) < k {
			h = PushMin(h, value)
		} else if value > h[0] {
			// The root is the weakest retained winner.
			h[0] = value
			RepairMinRoot(h)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(h)))
	return h
}

type TimeStore map[string][]Version

func (s TimeStore) Set(key, value string, at int64) bool {
	history := s[key]
	if len(history) > 0 {
		last := len(history) - 1
		if at < history[last].At {
			return false
		}
		if at == history[last].At {
			history[last].Value = value
			return true
		}
	}
	s[key] = append(history, Version{At: at, Value: value})
	return true
}

func (s TimeStore) Get(key string, at int64) (string, bool) {
	return VersionAt(s[key], at)
}

// Source: reading/04-graphs.md

func HopDistances(edges map[string][]string,
	start string) map[string]int {
	distance := map[string]int{start: 0}
	queue := []string{start}
	for head := 0; head < len(queue); head++ {
		from := queue[head]
		for _, to := range edges[from] {
			if _, seen := distance[to]; seen {
				continue
			}
			// Mark now so other parents skip this vertex.
			distance[to] = distance[from] + 1
			queue = append(queue, to)
		}
	}
	return distance
}

type TreeNode struct {
	Value       int
	Left, Right *TreeNode
}

func TreeLevels(root *TreeNode) [][]int {
	if root == nil {
		return nil
	}
	queue := []*TreeNode{root}
	var levels [][]int
	for head := 0; head < len(queue); {
		// Newly enqueued children belong to the next level.
		end := len(queue)
		var level []int
		for head < end {
			n := queue[head]
			head++
			level = append(level, n.Value)
			if n.Left != nil {
				queue = append(queue, n.Left)
			}
			if n.Right != nil {
				queue = append(queue, n.Right)
			}
		}
		levels = append(levels, level)
	}
	return levels
}

func IsBST(root *TreeNode) bool {
	previous, havePrevious := 0, false
	var visit func(*TreeNode) bool
	visit = func(n *TreeNode) bool {
		if n == nil {
			return true
		}
		if !visit(n.Left) {
			return false
		}
		// Inorder values must increase across subtree edges.
		if havePrevious && n.Value <= previous {
			return false
		}
		previous, havePrevious = n.Value, true
		return visit(n.Right)
	}
	return visit(root)
}

func IslandCount(grid [][]byte) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}
	rows, cols := len(grid), len(grid[0])
	directions := [][2]int{{1, 0}, {-1, 0},
		{0, 1}, {0, -1}}
	count := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] != '1' {
				continue
			}
			count++
			grid[r][c] = '0'
			stack := [][2]int{{r, c}}
			for len(stack) > 0 {
				p := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				for _, d := range directions {
					nr, nc := p[0]+d[0], p[1]+d[1]
					if nr < 0 || nr >= rows ||
						nc < 0 || nc >= cols ||
						grid[nr][nc] != '1' {
						continue
					}
					grid[nr][nc] = '0'
					stack = append(stack, [2]int{nr, nc})
				}
			}
		}
	}
	return count
}

func JobOrder(n int, edges [][2]int) ([]int, bool) {
	if n < 0 {
		return nil, false
	}
	next := make([][]int, n)
	pending := make([]int, n)
	for _, e := range edges {
		if e[0] < 0 || e[0] >= n ||
			e[1] < 0 || e[1] >= n {
			return nil, false
		}
		next[e[0]] = append(next[e[0]], e[1])
		pending[e[1]]++
	}
	var ready []int
	for id, count := range pending {
		if count == 0 {
			ready = append(ready, id)
		}
	}
	for head := 0; head < len(ready); head++ {
		for _, id := range next[ready[head]] {
			pending[id]--
			if pending[id] == 0 {
				// All prerequisites have now been processed.
				ready = append(ready, id)
			}
		}
	}
	return ready, len(ready) == n
}

type GraphNode struct {
	Label     string
	Neighbors []*GraphNode
}

func CloneGraph(start *GraphNode) *GraphNode {
	copied := make(map[*GraphNode]*GraphNode)
	var clone func(*GraphNode) *GraphNode
	clone = func(n *GraphNode) *GraphNode {
		if n == nil {
			return nil
		}
		if existing, found := copied[n]; found {
			return existing
		}
		duplicate := &GraphNode{Label: n.Label}
		// Register before recursion to terminate cycles.
		copied[n] = duplicate
		for _, neighbor := range n.Neighbors {
			duplicate.Neighbors = append(duplicate.Neighbors,
				clone(neighbor))
		}
		return duplicate
	}
	return clone(start)
}

// Source: reading/05-states.md

func Balanced(text string) bool {
	pairs := map[byte]byte{')': '(', ']': '[', '}': '{'}
	var stack []byte
	for i := 0; i < len(text); i++ {
		ch := text[i]
		if ch == '(' || ch == '[' || ch == '{' {
			stack = append(stack, ch)
			continue
		}
		opener, valid := pairs[ch]
		if !valid || len(stack) == 0 ||
			stack[len(stack)-1] != opener {
			return false
		}
		stack = stack[:len(stack)-1]
	}
	return len(stack) == 0
}

func GreaterWait(values []int) []int {
	waits := make([]int, len(values))
	var stack []int
	for i, value := range values {
		for len(stack) > 0 {
			j := stack[len(stack)-1]
			if value <= values[j] {
				break
			}
			stack = stack[:len(stack)-1]
			// This is the first later value greater than j's.
			waits[j] = i - j
		}
		stack = append(stack, i)
	}
	return waits
}

func MostScreenings(input [][2]int) [][2]int {
	a := append([][2]int(nil), input...)
	sort.Slice(a, func(i, j int) bool {
		if a[i][1] != a[j][1] {
			return a[i][1] < a[j][1]
		}
		return a[i][0] < a[j][0]
	})
	var chosen [][2]int
	end, haveEnd := 0, false
	for _, interval := range a {
		// Earliest finish leaves room for later screenings.
		if !haveEnd || interval[0] >= end {
			chosen = append(chosen, interval)
			end, haveEnd = interval[1], true
		}
	}
	return chosen
}

func PlaylistChoices(duration []int,
	count, budget int) [][]int {
	if count < 0 || budget < 0 {
		return nil
	}
	var chosen []int
	var out [][]int
	var search func(int, int)
	search = func(start, remaining int) {
		if len(chosen) == count {
			// Later branches reuse chosen's backing array.
			snapshot := append([]int{}, chosen...)
			out = append(out, snapshot)
			return
		}
		needed := count - len(chosen)
		// Stop when too few titles remain to finish a choice.
		for i := start; i <= len(duration)-needed; i++ {
			if duration[i] > remaining {
				continue
			}
			chosen = append(chosen, i)
			search(i+1, remaining-duration[i])
			// Undo this choice before exploring its sibling.
			chosen = chosen[:len(chosen)-1]
		}
	}
	search(0, budget)
	return out
}

func BestNonAdjacent(reward []int64) int64 {
	var next, afterNext int64
	for i := len(reward) - 1; i >= 0; i-- {
		// next is best[i+1]; afterNext is best[i+2].
		take := reward[i] + afterNext
		current := next
		if take > current {
			current = take
		}
		afterNext = next
		next = current
	}
	return next
}

// Source: reading/06-cache.md

type cacheNode struct {
	key        string
	value      int
	prev, next *cacheNode
}

type ReadCache struct {
	byKey      map[string]*cacheNode
	head, tail *cacheNode
	capacity   int
}

func NewReadCache(capacity int) *ReadCache {
	if capacity < 0 {
		capacity = 0
	}
	head, tail := &cacheNode{}, &cacheNode{}
	head.next, tail.prev = tail, head
	return &ReadCache{
		byKey: make(map[string]*cacheNode),
		head:  head, tail: tail, capacity: capacity,
	}
}

func detach(n *cacheNode) {
	// Join the neighbors before clearing this node's links.
	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev, n.next = nil, nil
}

func insertFront(head, n *cacheNode) {
	first := head.next
	n.prev, n.next = head, first
	head.next = n
	first.prev = n
}

func (c *ReadCache) remove(n *cacheNode) {
	// Membership must change in both the list and the map.
	detach(n)
	delete(c.byKey, n.key)
}

func (c *ReadCache) Get(key string) (int, bool) {
	n, found := c.byKey[key]
	if !found {
		return 0, false
	}
	detach(n)
	insertFront(c.head, n)
	return n.value, true
}

func (c *ReadCache) Put(key string, value int) {
	if c.capacity == 0 {
		return
	}
	if n, found := c.byKey[key]; found {
		// An overwrite reuses the existing list node.
		n.value = value
		detach(n)
		insertFront(c.head, n)
		return
	}
	n := &cacheNode{key: key, value: value}
	c.byKey[key] = n
	insertFront(c.head, n)
	if len(c.byKey) > c.capacity {
		// One insertion can require only one count eviction.
		c.remove(c.tail.prev)
	}
}

type WeightedCache struct {
	list         *ReadCache
	weights      map[string]int
	used, budget int
}

func NewWeightedCache(budget int) *WeightedCache {
	if budget < 0 {
		budget = 0
	}
	return &WeightedCache{
		list:    NewReadCache(0),
		weights: make(map[string]int), budget: budget,
	}
}

func (w *WeightedCache) Get(key string) (int, bool) {
	return w.list.Get(key)
}

func (w *WeightedCache) Put(key string,
	value, weight int) bool {
	// Rejection must preserve the old value and its recency.
	if weight <= 0 || weight > w.budget {
		return false
	}
	c := w.list
	n, found := c.byKey[key]
	if found {
		w.used -= w.weights[key]
		detach(n)
	} else {
		n = &cacheNode{key: key}
		c.byKey[key] = n
	}
	n.value = value
	insertFront(c.head, n)
	w.weights[key] = weight
	w.used += weight
	// A heavier overwrite may require several victims.
	for w.used > w.budget {
		victim := c.tail.prev
		w.used -= w.weights[victim.key]
		delete(w.weights, victim.key)
		c.remove(victim)
	}
	return true
}

// Source: reading/07-time.md

type TimedValue struct {
	Value      string
	Deadline   int64
	Generation uint64
}

type Expiry struct {
	Key        string
	Deadline   int64
	Generation uint64
}

type TTLMap struct {
	current map[string]TimedValue
	due     ExpiryHeap
	next    uint64
	clock   func() int64
}

func NewTTLMap(clock func() int64) *TTLMap {
	return &TTLMap{
		current: make(map[string]TimedValue),
		clock:   clock,
	}
}

func (c *TTLMap) Put(key, value string, ttl int64) {
	if ttl <= 0 {
		delete(c.current, key)
		return
	}
	deadline := c.clock() + ttl
	// A new generation distinguishes refreshes and reinserts.
	c.next++
	c.current[key] = TimedValue{
		Value: value, Deadline: deadline,
		Generation: c.next,
	}
	heap.Push(&c.due, Expiry{
		Key: key, Deadline: deadline, Generation: c.next,
	})
}

func (c *TTLMap) Get(key string) (string, bool) {
	entry, found := c.current[key]
	if !found {
		return "", false
	}
	if c.clock() >= entry.Deadline {
		delete(c.current, key)
		return "", false
	}
	return entry.Value, true
}

type ExpiryHeap []Expiry

func (h ExpiryHeap) Len() int { return len(h) }
func (h ExpiryHeap) Less(i, j int) bool {
	return h[i].Deadline < h[j].Deadline
}
func (h ExpiryHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *ExpiryHeap) Push(value any) {
	*h = append(*h, value.(Expiry))
}
func (h *ExpiryHeap) Pop() any {
	last := len(*h) - 1
	value := (*h)[last]
	// Release the key held in the unused backing-array slot.
	(*h)[last] = Expiry{}
	*h = (*h)[:last]
	return value
}

func ApplyExpiry(current map[string]TimedValue,
	record Expiry, now int64) bool {
	if record.Deadline > now {
		return false
	}
	entry, found := current[record.Key]
	// An old timer must not remove a newer write.
	if !found || entry.Generation != record.Generation {
		return false
	}
	delete(current, record.Key)
	return true
}

func (c *TTLMap) Cleanup() {
	now := c.clock()
	for len(c.due) > 0 && c.due[0].Deadline <= now {
		record := heap.Pop(&c.due).(Expiry)
		ApplyExpiry(c.current, record, now)
	}
}

type SafeReadCache struct {
	mu    sync.Mutex
	cache *ReadCache
}

func NewSafeReadCache(capacity int) *SafeReadCache {
	return &SafeReadCache{cache: NewReadCache(capacity)}
}

func (s *SafeReadCache) Get(key string) (int, bool) {
	// A cache hit changes recency, so Get needs a write lock.
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cache.Get(key)
}

func (s *SafeReadCache) Put(key string, value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache.Put(key, value)
}

type loadCall struct {
	done  chan struct{}
	value string
	err   error
}

type SharedLoader struct {
	mu     sync.Mutex
	active map[string]*loadCall
}

func (l *SharedLoader) Do(key string,
	load func() (string, error)) (string, error) {
	l.mu.Lock()
	if call, found := l.active[key]; found {
		l.mu.Unlock()
		// Wait without holding the lock needed by the leader.
		<-call.done
		return call.value, call.err
	}
	if l.active == nil {
		l.active = make(map[string]*loadCall)
	}
	call := &loadCall{done: make(chan struct{})}
	l.active[key] = call
	l.mu.Unlock()

	// Other keys can proceed while this backend call runs.
	value, err := load()
	l.mu.Lock()
	call.value, call.err = value, err
	delete(l.active, key)
	// Publish the result before waking every waiter.
	close(call.done)
	l.mu.Unlock()
	return value, err
}

// Source: reading/08-limits.md

type WindowLimiter struct {
	Window   int64
	Limit    int
	accepted []int64
	head     int
}

func (l *WindowLimiter) Allow(now int64) bool {
	if l.Window <= 0 || l.Limit <= 0 {
		return false
	}
	cutoff := now - l.Window
	// The left window boundary is excluded.
	for l.head < len(l.accepted) &&
		l.accepted[l.head] <= cutoff {
		l.head++
	}
	if l.head > 0 && l.head*2 >= len(l.accepted) {
		// Reclaim consumed storage without copying each time.
		l.accepted = append([]int64(nil),
			l.accepted[l.head:]...)
		l.head = 0
	}
	if len(l.accepted)-l.head >= l.Limit {
		return false
	}
	// Rejected attempts never enter the accepted log.
	l.accepted = append(l.accepted, now)
	return true
}

type UserLimiter struct {
	mu     sync.Mutex
	users  map[string]*WindowLimiter
	Window int64
	Limit  int
}

func (u *UserLimiter) Allow(user string, now int64) bool {
	// Creation and admission share one atomic decision.
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.users == nil {
		u.users = make(map[string]*WindowLimiter)
	}
	limiter := u.users[user]
	if limiter == nil {
		limiter = &WindowLimiter{
			Window: u.Window, Limit: u.Limit,
		}
		u.users[user] = limiter
	}
	return limiter.Allow(now)
}

type FixedLimiter struct {
	Window      int64
	Limit       int
	bucket      int64
	count       int
	initialized bool
}

func (f *FixedLimiter) Allow(now int64) bool {
	if f.Window <= 0 || f.Limit <= 0 || now < 0 {
		return false
	}
	bucket := now / f.Window
	if !f.initialized || bucket != f.bucket {
		// Each aligned window starts a separate allowance.
		f.bucket, f.count = bucket, 0
		f.initialized = true
	}
	if f.count >= f.Limit {
		return false
	}
	f.count++
	return true
}

type TokenBucket struct {
	Capacity, Rate float64
	tokens, last   float64
}

func NewTokenBucket(capacity, rate,
	now float64) *TokenBucket {
	return &TokenBucket{
		Capacity: capacity, Rate: rate,
		tokens: capacity, last: now,
	}
}

func (b *TokenBucket) Allow(now float64) bool {
	if now < b.last {
		now = b.last
	}
	b.tokens += (now - b.last) * b.Rate
	if b.tokens > b.Capacity {
		b.tokens = b.Capacity
	}
	// Advance on rejection too; never count refill twice.
	b.last = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// Source: reading/09-events.md

type Deduper struct {
	mu      sync.Mutex
	expiry  map[string]int64
	Horizon int64
}

func (d *Deduper) Accept(id string, now int64) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.Horizon <= 0 {
		return true
	}
	if deadline, found := d.expiry[id]; found &&
		now < deadline {
		// A duplicate does not extend its suppression window.
		return false
	}
	if d.expiry == nil {
		d.expiry = make(map[string]int64)
	}
	d.expiry[id] = now + d.Horizon
	return true
}

type HitBucket struct {
	At    int64
	Count int
}

type HitCounter struct {
	Window  int64
	buckets []HitBucket
	total   int
}

func (c *HitCounter) prune(now int64) {
	expired := 0
	for expired < len(c.buckets) &&
		c.buckets[expired].At <= now-c.Window {
		// Reverse each expired bucket's contribution once.
		c.total -= c.buckets[expired].Count
		expired++
	}
	c.buckets = c.buckets[expired:]
}

func (c *HitCounter) Record(now int64) {
	c.prune(now)
	last := len(c.buckets) - 1
	if last >= 0 && c.buckets[last].At == now {
		c.buckets[last].Count++
	} else {
		c.buckets = append(c.buckets, HitBucket{now, 1})
	}
	c.total++
}

func (c *HitCounter) Count(now int64) int {
	c.prune(now)
	return c.total
}

type JobRecord struct {
	ID         string
	Due        int64
	Generation uint64
}

func ClaimJob(current map[string]uint64,
	record JobRecord, now int64) bool {
	if record.Due > now {
		return false
	}
	generation, found := current[record.ID]
	if !found || generation != record.Generation {
		return false
	}
	// Remove authorization before handing work to a caller.
	delete(current, record.ID)
	return true
}

func ReadyJobs(due *ExpiryHeap,
	current map[string]uint64, now int64) []string {
	var ready []string
	for due.Len() > 0 && (*due)[0].Deadline <= now {
		record := heap.Pop(due).(Expiry)
		job := JobRecord{
			ID: record.Key, Due: record.Deadline,
			Generation: record.Generation,
		}
		// Cancelled or replaced generations are discarded.
		if ClaimJob(current, job, now) {
			ready = append(ready, job.ID)
		}
	}
	return ready
}

// Source: reading/10-recommend.md

type Candidate struct {
	ID    string
	Score int64
}

func TopTitles(candidates []Candidate, k int) []Candidate {
	if k <= 0 {
		return nil
	}
	ranked := append([]Candidate(nil), candidates...)
	sort.Slice(ranked, func(i, j int) bool {
		a, b := ranked[i], ranked[j]
		if a.Score != b.Score {
			return a.Score > b.Score
		}
		// Resolve ties before truncating to k winners.
		return a.ID < b.ID
	})
	if k < len(ranked) {
		ranked = ranked[:k]
	}
	return ranked
}

type CatalogCandidate struct {
	ID, Genre string
	Score     int64
	Eligible  bool
}

func RankEligible(input []CatalogCandidate,
	watched map[string]bool, k int) []Candidate {
	if k <= 0 {
		return nil
	}
	unique := make(map[string]CatalogCandidate)
	for _, item := range input {
		// Filter before an ID can occupy a result slot.
		if !item.Eligible || watched[item.ID] {
			continue
		}
		old, found := unique[item.ID]
		if !found || item.Score > old.Score ||
			(item.Score == old.Score &&
				item.Genre < old.Genre) {
			unique[item.ID] = item
		}
	}
	candidates := make([]Candidate, 0, len(unique))
	for _, item := range unique {
		candidates = append(candidates,
			Candidate{ID: item.ID, Score: item.Score})
	}
	return TopTitles(candidates, k)
}

func GenreWinners(input []CatalogCandidate,
	k int) []CatalogCandidate {
	if k <= 0 {
		return nil
	}
	ranked := append([]CatalogCandidate(nil), input...)
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].Score != ranked[j].Score {
			return ranked[i].Score > ranked[j].Score
		}
		return ranked[i].ID < ranked[j].ID
	})
	used := make(map[string]bool)
	var out []CatalogCandidate
	for _, item := range ranked {
		if used[item.Genre] {
			continue
		}
		used[item.Genre] = true
		out = append(out, item)
		if len(out) == k {
			break
		}
	}
	return out
}

func CollaborativeScores(target []string,
	neighbors [][]string) map[string]int64 {
	watched := make(map[string]bool)
	for _, id := range target {
		watched[id] = true
	}
	scores := make(map[string]int64)
	for _, history := range neighbors {
		distinct := make(map[string]bool)
		var weight int64
		for _, id := range history {
			// Repeated views count once toward similarity.
			if !distinct[id] && watched[id] {
				weight++
			}
			distinct[id] = true
		}
		if weight == 0 {
			continue
		}
		for id := range distinct {
			if !watched[id] {
				scores[id] += weight
			}
		}
	}
	return scores
}

type StatEvent struct {
	At      int64
	Title   string
	Minutes int64
}

type VideoTotal struct {
	Sum   int64
	Count int
}

type RollingStats struct {
	Window int64
	events []StatEvent
	totals map[string]VideoTotal
}

func (s *RollingStats) Expire(now int64) {
	expired := 0
	for expired < len(s.events) &&
		s.events[expired].At <= now-s.Window {
		e := s.events[expired]
		// Undo both parts of the average, not just the sum.
		total := s.totals[e.Title]
		total.Sum -= e.Minutes
		total.Count--
		if total.Count == 0 {
			delete(s.totals, e.Title)
		} else {
			s.totals[e.Title] = total
		}
		s.events[expired] = StatEvent{}
		expired++
	}
	s.events = s.events[expired:]
}

func (s *RollingStats) Record(e StatEvent) bool {
	if e.Minutes <= 0 {
		return false
	}
	s.Expire(e.At)
	if s.totals == nil {
		s.totals = make(map[string]VideoTotal)
	}
	total := s.totals[e.Title]
	total.Sum += e.Minutes
	total.Count++
	s.totals[e.Title] = total
	s.events = append(s.events, e)
	return true
}

func (s *RollingStats) Average(title string,
	now int64) (float64, bool) {
	s.Expire(now)
	total, found := s.totals[title]
	if !found {
		return 0, false
	}
	return float64(total.Sum) / float64(total.Count), true
}

// Source: reading/11-practical.md

type FSNode struct {
	IsFile   bool
	Content  string
	Children map[string]*FSNode
}

func Resolve(root *FSNode,
	parts []string) (*FSNode, bool) {
	current := root
	for _, name := range parts {
		// A file cannot serve as an intermediate directory.
		if current == nil || current.IsFile {
			return nil, false
		}
		current = current.Children[name]
		if current == nil {
			return nil, false
		}
	}
	return current, current != nil
}

type Edit struct {
	Value   string
	Deleted bool
}

func ReadLayered(base map[string]string,
	layers []map[string]Edit, key string) (string, bool) {
	for i := len(layers) - 1; i >= 0; i-- {
		edit, found := layers[i][key]
		if !found {
			continue
		}
		if edit.Deleted {
			// A tombstone hides every older value.
			return "", false
		}
		return edit.Value, true
	}
	value, found := base[key]
	return value, found
}

type TransactionStore struct {
	base   map[string]string
	layers []map[string]Edit
}

func (s *TransactionStore) Begin() {
	s.layers = append(s.layers, make(map[string]Edit))
}

func (s *TransactionStore) Change(key string, edit Edit) {
	if len(s.layers) > 0 {
		s.layers[len(s.layers)-1][key] = edit
		return
	}
	if s.base == nil {
		s.base = make(map[string]string)
	}
	if edit.Deleted {
		delete(s.base, key)
	} else {
		s.base[key] = edit.Value
	}
}

func (s *TransactionStore) Rollback() bool {
	if len(s.layers) == 0 {
		return false
	}
	last := len(s.layers) - 1
	// Release the discarded map from the backing array.
	s.layers[last] = nil
	s.layers = s.layers[:last]
	return true
}

func (s *TransactionStore) Commit() bool {
	if len(s.layers) == 0 {
		return false
	}
	last := len(s.layers) - 1
	top := s.layers[last]
	s.layers[last] = nil
	s.layers = s.layers[:last]
	// Change now targets the parent overlay or the base.
	for key, edit := range top {
		s.Change(key, edit)
	}
	return true
}

func (s *TransactionStore) Get(key string) (string, bool) {
	return ReadLayered(s.base, s.layers, key)
}

type TopicBus struct {
	mu        sync.Mutex
	callbacks map[string]map[int]func(string)
}

func (b *TopicBus) Publish(topic, message string) {
	b.mu.Lock()
	var snapshot []func(string)
	for _, callback := range b.callbacks[topic] {
		snapshot = append(snapshot, callback)
	}
	b.mu.Unlock()
	// Callbacks may safely change subscriptions after unlock.
	for _, callback := range snapshot {
		callback(message)
	}
}

// Source: reading/12-case-studies.md

type metadataNode struct {
	key, value string
	deadline   int64
	prev, next *metadataNode
}

type MetadataCache struct {
	capacity   int
	clock      func() int64
	byKey      map[string]*metadataNode
	head, tail *metadataNode
}

func NewMetadataCache(capacity int,
	clock func() int64) *MetadataCache {
	if capacity < 0 {
		capacity = 0
	}
	head, tail := &metadataNode{}, &metadataNode{}
	head.next, tail.prev = tail, head
	return &MetadataCache{
		capacity: capacity, clock: clock,
		byKey: make(map[string]*metadataNode),
		head:  head, tail: tail,
	}
}

func unlinkMetadata(n *metadataNode) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

func (c *MetadataCache) front(n *metadataNode) {
	first := c.head.next
	n.prev, n.next = c.head, first
	c.head.next, first.prev = n, n
}

func (c *MetadataCache) remove(n *metadataNode) {
	unlinkMetadata(n)
	delete(c.byKey, n.key)
}

func (c *MetadataCache) expire(now int64) {
	for _, n := range c.byKey {
		if n.deadline <= now {
			c.remove(n)
		}
	}
}

func (c *MetadataCache) Get(key string) (string, bool) {
	n, found := c.byKey[key]
	if !found {
		return "", false
	}
	if c.clock() >= n.deadline {
		// Expiration takes precedence over recency promotion.
		c.remove(n)
		return "", false
	}
	unlinkMetadata(n)
	c.front(n)
	return n.value, true
}

func (c *MetadataCache) Put(key, value string,
	ttl int64) {
	if ttl <= 0 {
		if n, found := c.byKey[key]; found {
			c.remove(n)
		}
		return
	}
	if c.capacity == 0 {
		return
	}
	now := c.clock()
	// Reclaim expired residents before evicting a live one.
	c.expire(now)
	n, found := c.byKey[key]
	if found {
		unlinkMetadata(n)
	} else {
		n = &metadataNode{key: key}
		c.byKey[key] = n
	}
	n.value, n.deadline = value, now+ttl
	c.front(n)
	if len(c.byKey) > c.capacity {
		c.remove(c.tail.prev)
	}
}

type WatchEvent struct {
	ID, Title   string
	At, Minutes int64
}

func RecentTotals(events []WatchEvent,
	now, window int64) map[string]int64 {
	seen := make(map[string]bool)
	totals := make(map[string]int64)
	for _, event := range events {
		// Exclude the left boundary and all future events.
		if event.At <= now-window || event.At > now {
			continue
		}
		if seen[event.ID] {
			continue
		}
		// Deduplicate deliveries by event ID, not title ID.
		seen[event.ID] = true
		totals[event.Title] += event.Minutes
	}
	return totals
}

func RecentRanking(events []WatchEvent,
	now, window int64, k int) ([]Candidate, bool) {
	if window <= 0 {
		return nil, false
	}
	// Validate the whole batch even when k is nonpositive.
	for _, event := range events {
		if event.Minutes <= 0 {
			return nil, false
		}
	}
	if k <= 0 {
		return nil, true
	}
	totals := RecentTotals(events, now, window)
	candidates := make([]Candidate, 0, len(totals))
	for title, minutes := range totals {
		candidates = append(candidates,
			Candidate{ID: title, Score: minutes})
	}
	return TopTitles(candidates, k), true
}
