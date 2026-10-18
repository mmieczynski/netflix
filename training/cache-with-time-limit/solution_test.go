package cache_test

import (
	"fmt"
	"testing"

	"netflixprep/training/cache-with-time-limit"
)

type action struct {
	at        int64
	operation string
	key       int
	value     int
	duration  int64
	want      interface{}
}

func set(at int64, key, value int, duration int64, want bool) action {
	return action{at: at, operation: "set", key: key, value: value, duration: duration, want: want}
}

func get(at int64, key, want int) action {
	return action{at: at, operation: "get", key: key, want: want}
}

func count(at int64, want int) action {
	return action{at: at, operation: "count", want: want}
}

// runTrace controls time without sleeping and checks only the public API.
func runTrace(t *testing.T, actions []action) {
	t.Helper()
	var now int64
	var c cache.Cache = cache.NewTimeLimitedCache(func() int64 { return now })
	for i, a := range actions {
		if a.at < now {
			t.Fatalf("invalid test: time moved backwards at step %d", i+1)
		}
		now = a.at
		var got interface{}
		var call string
		switch a.operation {
		case "set":
			call = fmt.Sprintf("Set(%d, %d, %d)", a.key, a.value, a.duration)
			got = c.Set(a.key, a.value, a.duration)
		case "get":
			call = fmt.Sprintf("Get(%d)", a.key)
			got = c.Get(a.key)
		case "count":
			call = "Count()"
			got = c.Count()
		default:
			t.Fatalf("invalid test operation %q", a.operation)
		}
		if got != a.want {
			t.Errorf("step %d, t=%dms, %s: got %v, want %v", i+1, now, call, got, a.want)
		}
	}
}

func TestSamples(t *testing.T) {
	t.Run("expiration", func(t *testing.T) {
		runTrace(t, []action{
			set(0, 7, 80, 100, false),
			get(60, 7, 80),
			count(60, 1),
			get(100, 7, -1),
			count(100, 0),
		})
	})
	t.Run("replacement", func(t *testing.T) {
		runTrace(t, []action{
			set(0, 4, 12, 40, false),
			set(25, 4, 99, 80, true),
			get(40, 4, 99),
			get(104, 4, 99),
			get(105, 4, -1),
			count(105, 0),
		})
	})
}

func TestCacheBehavior(t *testing.T) {
	tests := []struct {
		name    string
		actions []action
	}{
		{"empty_cache_and_missing_key", []action{
			count(0, 0),
			get(0, 8, -1),
			set(0, 1, 30, 10, false),
			get(0, 8, -1),
			count(0, 1),
		}},
		{"before_at_and_after_expiration", []action{
			set(10, 9, 21, 50, false),
			get(59, 9, 21),
			count(59, 1),
			get(60, 9, -1),
			count(60, 0),
			get(61, 9, -1),
		}},
		{"replacement_can_shorten_lifetime", []action{
			set(10, 2, 3, 100, false),
			set(20, 2, 4, 10, true),
			count(20, 1),
			get(29, 2, 4),
			get(30, 2, -1),
			count(30, 0),
			get(110, 2, -1),
		}},
		{"expired_key_is_new_without_prior_read", []action{
			set(0, 5, 7, 10, false),
			set(11, 5, 8, 10, false),
			count(11, 1),
			get(20, 5, 8),
			get(21, 5, -1),
		}},
		{"write_at_exact_expiration", []action{
			set(0, 3, 1, 10, false),
			set(10, 3, 2, 20, false),
			get(10, 3, 2),
			count(10, 1),
			get(30, 3, -1),
		}},
		{"zero_duration_on_missing_key", []action{
			set(0, 5, 6, 0, false),
			get(0, 5, -1),
			count(0, 0),
			set(0, 5, 7, 0, false),
			get(1, 5, -1),
		}},
		{"zero_duration_replaces_active_key", []action{
			set(0, 5, 6, 10, false),
			set(5, 5, 7, 0, true),
			get(5, 5, -1),
			count(5, 0),
			set(5, 5, 8, 10, false),
			get(14, 5, 8),
			get(15, 5, -1),
		}},
		{"count_expires_keys_without_reads", []action{
			set(0, 1, 10, 20, false),
			set(0, 2, 20, 30, false),
			set(10, 3, 30, 20, false),
			count(19, 3),
			count(20, 2),
			count(29, 2),
			count(30, 0),
			count(30, 0),
			count(31, 0),
		}},
		{"expiration_order_differs_from_insertion_order", []action{
			set(0, 1, 10, 100, false),
			set(1, 2, 20, 10, false),
			set(2, 3, 30, 20, false),
			count(10, 3),
			count(11, 2),
			get(11, 1, 10),
			get(11, 2, -1),
			count(22, 1),
			get(100, 1, -1),
			count(100, 0),
		}},
		{"extend_earliest_deadline_past_other_keys", []action{
			set(0, 1, 10, 10, false),
			set(0, 2, 20, 20, false),
			set(0, 3, 30, 30, false),
			set(5, 1, 11, 100, true),
			count(10, 3),
			count(20, 2),
			get(20, 1, 11),
			get(20, 2, -1),
			count(30, 1),
			get(104, 1, 11),
			get(105, 1, -1),
		}},
		{"shorten_later_deadline_to_expire_first", []action{
			set(0, 1, 10, 10, false),
			set(0, 2, 20, 20, false),
			set(0, 3, 30, 30, false),
			set(0, 4, 40, 40, false),
			set(1, 4, 41, 2, true),
			count(2, 4),
			count(3, 3),
			get(3, 4, -1),
			get(3, 1, 10),
			get(3, 2, 20),
			get(3, 3, 30),
			count(10, 2),
			count(30, 0),
		}},
		{"remove_middle_deadline_then_update_and_reinsert", []action{
			set(0, 1, 10, 10, false),
			set(0, 2, 20, 20, false),
			set(0, 3, 30, 30, false),
			set(0, 4, 40, 40, false),
			set(1, 2, 21, 0, true),
			count(1, 3),
			set(1, 4, 41, 1, true),
			count(2, 2),
			set(2, 2, 21, 3, false),
			count(5, 2),
			count(10, 1),
			get(10, 3, 30),
			count(30, 0),
		}},
		{"reads_and_counts_do_not_extend_lifetime", []action{
			set(0, 8, 22, 10, false),
			get(0, 8, 22),
			get(6, 8, 22),
			count(9, 1),
			get(9, 8, 22),
			get(10, 8, -1),
			count(10, 0),
		}},
		{"zero_and_maximum_keys_and_values", []action{
			set(0, 0, 0, 5, false),
			set(0, 1_000_000_000, 1_000_000_000, 5, false),
			get(4, 0, 0),
			get(4, 1_000_000_000, 1_000_000_000),
			count(4, 2),
			count(5, 0),
			get(5, 0, -1),
		}},
		{"same_value_still_gets_new_lifetime", []action{
			set(0, 6, 50, 10, false),
			set(8, 6, 50, 10, true),
			get(10, 6, 50),
			count(10, 1),
			get(17, 6, 50),
			get(18, 6, -1),
		}},
		{"reinsert_after_expired_read", []action{
			set(0, 2, 10, 5, false),
			get(5, 2, -1),
			set(5, 2, 20, 4, false),
			get(5, 2, 20),
			count(8, 1),
			count(9, 0),
		}},
		{"old_deadline_does_not_remove_reinserted_key", []action{
			set(0, 1, 10, 100, false),
			set(10, 1, 20, 10, true),
			count(20, 0),
			set(25, 1, 30, 100, false),
			get(100, 1, 30),
			count(100, 1),
			get(125, 1, -1),
			count(125, 0),
		}},
		{"large_timestamp_and_maximum_duration", []action{
			set(999_999_999_000, 1, 2, 1_000, false),
			get(999_999_999_999, 1, 2),
			count(999_999_999_999, 1),
			get(1_000_000_000_000, 1, -1),
			count(1_000_000_000_000, 0),
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runTrace(t, tc.actions)
		})
	}
}

func TestIndependentCaches(t *testing.T) {
	firstNow, secondNow := int64(0), int64(1_000)
	first := cache.NewTimeLimitedCache(func() int64 { return firstNow })
	if first.Set(1, 10, 10) {
		t.Error("first cache: initial Set returned true")
	}
	second := cache.NewTimeLimitedCache(func() int64 { return secondNow })
	if got := second.Count(); got != 0 {
		t.Errorf("new second cache: Count() = %d, want 0", got)
	}
	if second.Set(1, 20, 5) {
		t.Error("second cache: initial Set returned true")
	}
	if got := first.Get(1); got != 10 {
		t.Errorf("first cache after second write: Get(1) = %d, want 10", got)
	}
	if got := second.Get(1); got != 20 {
		t.Errorf("second cache: Get(1) = %d, want 20", got)
	}
	firstNow = 10
	if got := first.Count(); got != 0 {
		t.Errorf("expired first cache: Count() = %d, want 0", got)
	}
	if got := second.Count(); got != 1 {
		t.Errorf("second cache with unchanged clock: Count() = %d, want 1", got)
	}
	secondNow = 1_005
	if got := second.Count(); got != 0 {
		t.Errorf("expired second cache: Count() = %d, want 0", got)
	}
}
