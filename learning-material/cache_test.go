package prep

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

func checkLRU(t *testing.T, c *LRU) []string {
	t.Helper()
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := []string{}
	seen := map[*node]bool{}
	previous := c.head
	for n := c.head.next; n != c.tail; n = n.next {
		if n == nil || seen[n] {
			t.Fatal("broken list or cycle")
		}
		seen[n] = true
		if n.prev != previous || previous.next != n || c.items[n.key] != n {
			t.Fatal("list/map invariant")
		}
		keys = append(keys, n.key)
		previous = n
	}
	if c.tail.prev != previous || len(keys) != len(c.items) || len(keys) > c.capacity {
		t.Fatal("tail or size invariant")
	}
	return keys
}

func TestLRUAgainstSimpleModel(t *testing.T) {
	r := rand.New(rand.NewSource(13))
	for capacity := 0; capacity <= 5; capacity++ {
		c := NewLRU(capacity)
		values := map[string]int{}
		order := []string{} // slow model: most recent first
		promote := func(key string) {
			next := []string{key}
			for _, k := range order {
				if k != key {
					next = append(next, k)
				}
			}
			order = next
		}
		for step := 0; step < 500; step++ {
			key := fmt.Sprint(r.Intn(8))
			if r.Intn(2) == 0 {
				v := r.Intn(10)
				c.Put(key, v)
				if capacity > 0 {
					values[key] = v
					promote(key)
					if len(order) > capacity {
						delete(values, order[len(order)-1])
						order = order[:len(order)-1]
					}
				}
			} else {
				v, ok := c.Get(key)
				want, exists := values[key]
				if v != want || ok != exists {
					t.Fatal("model value mismatch")
				}
				if exists {
					promote(key)
				}
			}
			got := checkLRU(t, c)
			if fmt.Sprint(got) != fmt.Sprint(order) {
				t.Fatalf("order %v != %v", got, order)
			}
		}
	}
}

func TestTTLBoundariesAndOldGenerations(t *testing.T) {
	now := int64(0)
	c := NewTTLCache(func() int64 { return now })
	c.Put("a", 1, 10)
	now = 9
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatal("expired too early")
	}
	now = 10
	if _, ok := c.Get("a"); ok {
		t.Fatal("deadline must be expired")
	}
	c.Put("a", 2, 20)
	if c.Cleanup() != 0 {
		t.Fatal("old generation deleted replacement")
	}
	now = 15
	c.Put("a", 3, 30)
	now = 30
	if c.Cleanup() != 0 {
		t.Fatal("overwritten record deleted current value")
	}
	if v, ok := c.Get("a"); !ok || v != 3 {
		t.Fatal("lost refreshed value")
	}
	c.Put("a", 4, 0)
	if _, ok := c.Get("a"); ok {
		t.Fatal("nonpositive TTL must remove")
	}
	c.Put("a", 5, 30)
	now = 45
	if c.Cleanup() != 0 {
		t.Fatal("recreation reused generation")
	}
	now = 60
	if c.Cleanup() != 1 || len(c.items) != 0 {
		t.Fatal("idle cleanup failed")
	}
	if c.Cleanup() != 0 {
		t.Fatal("cleanup must be idempotent")
	}
}

func TestConcurrentCacheOperations(t *testing.T) {
	c := NewLRU(10)
	ttl := NewTTLCache(func() int64 { return 100 }) // immutable clock is concurrency-safe
	var wg sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				key := fmt.Sprint((id + i) % 20)
				c.Put(key, i)
				c.Get(key)
				ttl.Put(key, i, 10)
				ttl.Get(key)
				ttl.Cleanup()
			}
		}(worker)
	}
	wg.Wait()
	checkLRU(t, c)
}
