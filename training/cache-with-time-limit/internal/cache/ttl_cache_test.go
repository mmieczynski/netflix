package cache

import "testing"


func ShouldAddNewEntry(t *testing.T) {
	now := 0
	cache := NewTtlCache(func() int {return now})

	cache.Set("foo", 1, 10)
	if 
}