package cache

import (
	"errors"
	"strings"
)

// assumptions/constraints:
// - single threaded access
// - max 3 items in the cache
// - all items are the same size
// - one key one value in the cache
// - get does not fetch the db

// operations:
// get - check if the item exists in the cache, if it does, return its value and move it to the fron of the cache
// set - set given value in the front of the cache

// data structures:
// map - storing keys in the cache and pointers to the items
// doubly linked list - for items storage and lru queue

type Node struct {
	Prev, Next *Node
	Key        string
	Value      int
}

type LruCache struct {
	Cache      map[string]*Node
	Head, Tail *Node
	MaxSize    int
}

func NewLruCache(maxSize int) LruCache {
	head := &Node{}
	tail := &Node{}
	head.Next = tail
	tail.Prev = head
	return LruCache{
		Head:    head,
		Tail:    tail,
		Cache:   map[string]*Node{},
		MaxSize: maxSize,
	}
}

func (cache *LruCache) Set(key string, value int) error {
	if strings.TrimSpace(key) == "" || value <= 0 {
		return errors.New("Not accepted")
	}

	if node, ok := cache.Cache[key]; ok {
		node.Value = value
		node.Prev.Next = node.Next
		node.Next.Prev = node.Prev
		oldFront := cache.Head.Next
		cache.Head.Next = node
		oldFront.Prev = node
		node.Next = oldFront
		node.Prev = cache.Head
	} else {
		if len(cache.Cache) >= cache.MaxSize {
			last := cache.Tail.Prev
			cache.Tail.Prev = last.Prev
			last.Prev.Next = cache.Tail
			delete(cache.Cache, last.Key)
		}
		node := Node{Key: key, Value: value, Prev: cache.Head, Next: cache.Head.Next}

		// move first element to the back
		oldFirst := cache.Head.Next
		oldFirst.Prev = &node

		// connect new node to head
		cache.Head.Next = &node

		// add to the map
		cache.Cache[key] = &node
	}

	return nil
}

func (cache *LruCache) Get(key string) (value int, found bool) {
	if _, ok := cache.Cache[key]; !ok {
		return 0, false
	}

	node := cache.Cache[key]
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
	oldFront := cache.Head.Next
	cache.Head.Next = node
	oldFront.Prev = node
	node.Next = oldFront
	node.Prev = cache.Head

	return node.Value, true
}
