package cache

// Cache is the supplied contract. All durations are in milliseconds.
type Cache interface {
	// Set replaces a key's value and lifetime, reporting whether it was active before the call.
	Set(key, value int, durationMillis int64) bool
	// Get returns the current value, or -1 for an absent or expired key.
	Get(key int) int
	// Count returns the number of active keys at the current time.
	Count() int
}
