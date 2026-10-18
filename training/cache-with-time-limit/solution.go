package cache

import "container/heap"

type entry struct {
	key       int
	value     int
	expiresAt int64
}

type expirationHeap []*entry

func (h expirationHeap) Len() int { return len(h) }

func (h expirationHeap) Less(i, j int) bool {
	return h[i].expiresAt < h[j].expiresAt
}

func (h expirationHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *expirationHeap) Push(x any) {
	*h = append(*h, x.(*entry))
}

// heap.Pop moves the minimum to the end before calling this method.
func (h *expirationHeap) Pop() any {
	last := len(*h) - 1
	item := (*h)[last]
	(*h)[last] = nil // Release the removed entry's reference.
	*h = (*h)[:last]
	return item
}

// TimeLimitedCache stores values until their expiration timestamps.
type TimeLimitedCache struct {
	entries     map[int]*entry
	expirations expirationHeap
	now         func() int64
}

var _ Cache = (*TimeLimitedCache)(nil)

// NewTimeLimitedCache creates an empty cache.
// now returns the current time in milliseconds; tests control this clock.
func NewTimeLimitedCache(now func() int64) *TimeLimitedCache {
	// The initially empty heap already satisfies the heap invariant.
	return &TimeLimitedCache{
		entries: make(map[int]*entry),
		now:     now,
	}
}

func (c *TimeLimitedCache) cleanup(now int64) {
	for len(c.expirations) > 0 && c.expirations[0].expiresAt <= now {
		item := heap.Pop(&c.expirations).(*entry)
		// An older write must not remove a newer entry for the same key.
		if c.entries[item.key] == item {
			delete(c.entries, item.key)
		}
	}
}

func (c *TimeLimitedCache) Set(key, value int, durationMillis int64) bool {
	now := c.now()
	c.cleanup(now)
	_, exists := c.entries[key] // Every remaining map entry is active.

	// A zero-duration write removes even an active previous value.
	if durationMillis == 0 {
		delete(c.entries, key)
		return exists
	}

	// Entries are immutable after insertion; each write gets a new identity.
	item := &entry{
		key:       key,
		value:     value,
		expiresAt: now + durationMillis,
	}
	c.entries[key] = item
	heap.Push(&c.expirations, item)
	return exists
}

func (c *TimeLimitedCache) Get(key int) int {
	c.cleanup(c.now())
	item, exists := c.entries[key]
	if !exists {
		return -1
	}
	return item.value
}

func (c *TimeLimitedCache) Count() int {
	c.cleanup(c.now())
	return len(c.entries)
}
