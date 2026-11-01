package cache

import (
	"container/heap"
	"fmt"
)

// set - new entry if not existing, update with new expiry if existing
// get - nil if not in the cache, remove from the cache if after expiry, return if in cache and not expired

// map - quick access - points to the item in minHeap
// minHeap - for the item that is most expired

type CacheEntry struct {
	Key       string
	Value     int
	ExpiresAt int
}

type ExpiresHeap []*CacheEntry

func (this ExpiresHeap) Len() int           { return len(this) }
func (this ExpiresHeap) Less(i, j int) bool { return this[i].ExpiresAt < this[j].ExpiresAt }
func (this ExpiresHeap) Swap(i, j int)      { this[i], this[j] = this[j], this[i] }

func (this *ExpiresHeap) Push(entry any) {
	*this = append(*this, entry.(*CacheEntry))
}

func (this *ExpiresHeap) Pop() any {
	last := len(*this) - 1
	entry := (*this)[last]
	(*this)[last] = nil
	*this = (*this)[:last]
	return entry
}

type TtlCache struct {
	Cache map[string]*CacheEntry
	Mp    ExpiresHeap
	now   func() int
}

func NewTtlCache(c func() int) TtlCache {
	return TtlCache{
		Cache: map[string]*CacheEntry{},
		Mp:    ExpiresHeap{},
		now:   c,
	}
}

func (this *TtlCache) Set(key string, value int, ttl int) error {
	if value < 0 {
		return fmt.Errorf("Error")
	}

	entry := &CacheEntry{
		Key:       key,
		Value:     value,
		ExpiresAt: this.now() + ttl,
	}
	this.Cache[key] = entry
	heap.Push(&this.Mp, entry)
	return nil
}

func (this *TtlCache) Get(key string) int {
	this.cleanup()
	if entry, ok := this.Cache[key]; ok {
		return entry.Value
	}

	return -1
}

func (this *TtlCache) cleanup() {
	now := this.now()
	for len(this.Mp) > 0 && this.Mp[0].ExpiresAt <= now {
		entry := heap.Pop(&this.Mp).(*CacheEntry)
		if this.Cache[entry.Key] == entry {
			delete(this.Cache, entry.Key)
		}
	}
}
