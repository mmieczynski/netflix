package prep

import (
	"container/heap"
	"sync"
)

type node struct {
	key        string
	value      int
	prev, next *node
}

// LRU uses a manual doubly linked list. Use NewLRU; do not copy after use.
type LRU struct {
	mu         sync.Mutex
	capacity   int
	items      map[string]*node
	head, tail *node
}

func NewLRU(capacity int) *LRU {
	if capacity < 0 {
		panic("negative capacity")
	}
	c := &LRU{capacity: capacity, items: make(map[string]*node), head: &node{}, tail: &node{}}
	c.head.next, c.tail.prev = c.tail, c.head
	return c
}

// Helpers require the caller to hold mu.
func (c *LRU) detach(n *node) {
	n.prev.next, n.next.prev = n.next, n.prev
	n.prev, n.next = nil, nil
}

func (c *LRU) front(n *node) {
	n.prev, n.next = c.head, c.head.next
	c.head.next.prev = n
	c.head.next = n
}

func (c *LRU) Get(key string) (int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	n, ok := c.items[key]
	if !ok {
		return 0, false
	}
	c.detach(n)
	c.front(n)
	return n.value, true
}

func (c *LRU) Put(key string, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.capacity == 0 {
		return
	}
	if n, ok := c.items[key]; ok {
		n.value = value
		c.detach(n)
		c.front(n)
		return
	}
	n := &node{key: key, value: value}
	c.items[key] = n
	c.front(n)
	if len(c.items) > c.capacity {
		victim := c.tail.prev
		c.detach(victim)
		delete(c.items, victim.key)
	}
}

type ttlEntry struct {
	value      int
	deadline   int64
	generation uint64
}
type expiry struct {
	key        string
	deadline   int64
	generation uint64
}
type expiryHeap []expiry

func (h expiryHeap) Len() int            { return len(h) }
func (h expiryHeap) Less(i, j int) bool  { return h[i].deadline < h[j].deadline }
func (h expiryHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *expiryHeap) Push(x interface{}) { *h = append(*h, x.(expiry)) }
func (h *expiryHeap) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	old[len(old)-1] = expiry{} // release string reference
	*h = old[:len(old)-1]
	return x
}

// TTLCache demonstrates lazy versioned expiry, not bounded memory or LRU.
// now must be quick, concurrency-safe, monotone, and must not call this cache.
// Deadlines must fit int64; generations must not wrap during the cache lifetime.
// Call Cleanup explicitly: there is intentionally no background goroutine.
type TTLCache struct {
	mu    sync.Mutex
	now   func() int64
	items map[string]ttlEntry
	heap  expiryHeap
	seq   uint64
}

func NewTTLCache(now func() int64) *TTLCache {
	if now == nil {
		panic("nil clock")
	}
	return &TTLCache{now: now, items: make(map[string]ttlEntry)}
}

func (c *TTLCache) Put(key string, value int, ttl int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ttl <= 0 {
		delete(c.items, key)
		return
	}
	c.seq++
	entry := ttlEntry{value, c.now() + ttl, c.seq}
	c.items[key] = entry
	heap.Push(&c.heap, expiry{key, entry.deadline, entry.generation})
}

func (c *TTLCache) Get(key string) (int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.items[key]
	if !ok {
		return 0, false
	}
	if c.now() >= e.deadline {
		delete(c.items, key)
		return 0, false
	}
	return e.value, true
}

// Cleanup returns the number of current entries removed, not stale records popped.
func (c *TTLCache) Cleanup() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	now, removed := c.now(), 0
	for len(c.heap) > 0 && c.heap[0].deadline <= now {
		old := heap.Pop(&c.heap).(expiry)
		if current, ok := c.items[old.key]; ok && current.generation == old.generation {
			delete(c.items, old.key)
			removed++
		}
	}
	return removed
}
