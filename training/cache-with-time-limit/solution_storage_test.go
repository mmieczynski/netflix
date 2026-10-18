package cache

import "testing"

// Obsolete records may accumulate, but cleanup must release them when due.
func TestRepeatedWritesReleaseExpiredRecords(t *testing.T) {
	var now int64
	c := NewTimeLimitedCache(func() int64 { return now })
	for key := 0; key < 3; key++ {
		c.Set(key, key, 1_000)
	}
	for value := 0; value < 10_000; value++ {
		c.Set(1, value, int64(value%1_000+1))
	}
	if got := c.Count(); got != 3 {
		t.Fatalf("Count() after repeated writes = %d, want 3", got)
	}
	now = 999
	if got := c.Count(); got != 3 {
		t.Fatalf("Count() after obsolete deadlines = %d, want 3", got)
	}
	if got := c.Get(1); got != 9_999 {
		t.Fatalf("Get(1) = %d, want 9999", got)
	}
	now = 1_000
	if got := c.Count(); got != 0 {
		t.Fatalf("Count() at expiration = %d, want 0", got)
	}
	if got := len(c.expirations); got != 0 {
		t.Fatalf("cleanup retained %d expired records, want 0", got)
	}
}

// Cleanup must reclaim expired keys even when they are never read or counted.
func TestWritesReclaimUntouchedExpiredEntries(t *testing.T) {
	var now int64
	c := NewTimeLimitedCache(func() int64 { return now })
	for key := 0; key < 1_000; key++ {
		now = int64(key)
		c.Set(key, key, 1)
		if len(c.entries) != 1 || len(c.expirations) != 1 {
			t.Fatalf("after write %d: retained %d entries and %d expiration records, want 1 each",
				key, len(c.entries), len(c.expirations))
		}
	}
}
