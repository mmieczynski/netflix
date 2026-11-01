package cache

import "testing"

// test cases:
// set:
// happy path (insert x3, get x3)
// eviction of the last element
// insert whe key exists
// insert when key does not exist
// insert with wrong key and value

func CorrectlyInsertsAndReturnsCacheItems(t *testing.T) {
	cache := NewLruCache(3)

	if error := cache.Set("foo", 1); error != nil {
		t.fa
	}

}
