# Cache with a time limit

A Go practice project inspired by LeetCode's
[Cache With Time Limit](https://leetcode.com/problems/cache-with-time-limit/).
The layout follows the idea of CodeSignal's
[filesystem tasks](https://codesignal.com/blog/file-system-tasks/): supplied
interfaces, an implementation file, and unit tests that call your methods.
This is a local approximation for a live interview with another engineer;
the actual interview's starter files and format may differ.

## Project layout

| File | Role |
| --- | --- |
| `README.md` | Problem statement and examples |
| `interface.go` | Supplied method contract |
| `solution.go` | Complete reference implementation using a map and a min-heap with lazy invalidation |
| `solution_test.go` | Supplied examples and behavioral checks |
| `solution_storage_test.go` | Checks that cleanup releases due records, including obsolete writes |
| `go.mod` | Independent Go module, using only the standard library |

The reference solution uses a map to find the current entry for each key and a
min-heap to find the earliest expiration. The heap holds both current entries
and obsolete records from previous writes. All local tests are visible.

## Task

Implement a cache of integer keys and values. Each write gives its key a new
lifetime. The public API is:

```go
func NewTimeLimitedCache(now func() int64) *TimeLimitedCache
func (c *TimeLimitedCache) Set(key, value int, durationMillis int64) bool
func (c *TimeLimitedCache) Get(key int) int
func (c *TimeLimitedCache) Count() int
```

`Set` returns whether the key was active immediately before the write. Every
write replaces both its value and lifetime. `Get` returns the active value,
or `-1` when unavailable. `Count` includes only active keys.

A write at time `t` with duration `d` is active during `[t, t+d)`.
At exactly `t+d`, it has expired. A zero-duration write is immediately
unavailable, even to another operation at the same timestamp. If it replaces
an active key, `Set` still returns `true`, and the old value is gone.

Each new cache starts empty and has independent state. Reads and counts do
not extend lifetimes.

## Clock and inputs

The constructor receives a clock function instead of using real elapsed time.
The harness advances that clock before each operation, so tests need no sleeps.
Every time shown below is an **absolute timestamp in milliseconds**, not a delay
since the previous row.

- The clock is non-nil, never goes backwards, and stays fixed during a call.
- Multiple operations can share a timestamp; their listed order matters.
- Keys and values are between `0` and `1_000_000_000`, inclusive.
- Durations are between `0` and `1_000` milliseconds, inclusive.
- Clock values are between `0` and `1_000_000_000_000` milliseconds, inclusive.
- Calls are sequential. Concurrency, background timers, and invalid-input
  handling are outside this exercise.
- The harness constructs caches and calls methods directly. There is no
  `main`, input parser, or output formatting to implement.

Keep the supplied signatures. Choose your representation and be prepared to
explain the time and space cost of each operation.

## Example 1: expiration

Starting with a new cache:

| Time (ms) | Call | Result |
| --- | --- | --- |
| 0 | `Set(7, 80, 100)` | `false` |
| 60 | `Get(7)` | `80` |
| 60 | `Count()` | `1` |
| 100 | `Get(7)` | `-1` |
| 100 | `Count()` | `0` |

## Example 2: replacing a lifetime

Starting with another new cache:

| Time (ms) | Call | Result |
| --- | --- | --- |
| 0 | `Set(4, 12, 40)` | `false` |
| 25 | `Set(4, 99, 80)` | `true` |
| 40 | `Get(4)` | `99` |
| 104 | `Get(4)` | `99` |
| 105 | `Get(4)` | `-1` |
| 105 | `Count()` | `0` |

The second write gives key `4` a new expiration time of `105` milliseconds.
Its previous expiration time of `40` has no effect on the replacement.

## Reference design and cost

Each entry stores its key, value, and absolute expiration timestamp. Every write
with a positive duration creates a fresh entry object, updates the map to point
to it, and pushes it onto the heap. Entries already in the heap are never mutated.
The heap uses push/pop operations without tracking entry indices.

When cleanup pops a due record, it deletes the map entry only if the map still
points to that exact object. A record from an earlier write is otherwise discarded.
Pointer identity distinguishes writes even when keys, values, or deadlines match.
A zero-duration write removes the current map entry immediately; its older heap
records remain until their deadlines are processed.

Each public operation reads the clock once and removes due entries from the
heap's root before doing its own work. Cleanup stops at the first future deadline;
it never scans the map. After cleanup, every map entry is active, so `Count`
can return the map's length. Entries that expire while the cache is idle remain
stored until the next operation.

Let `n` be the number of current keys and `h` the number of heap records, including
obsolete writes. A positive-duration write takes O(log h) heap work; lookups,
counts, and zero-duration writes take expected O(1), plus that call's cleanup.
Removing `k` due records adds O(k log h); when no records are due, checking the
root takes O(1). If every record is due, one call can do O(h log h) cleanup work.

The logical storage cost is O(n + h). Repeated writes to one key can leave many
obsolete records in the heap before their deadlines pass, so heap size is not
bounded by the number of current keys. Cleanup releases due records; backing
allocations may still retain capacity afterward. This version does not rebuild
the heap to discard obsolete records early.

The implementation uses Go's standard
[`container/heap`](https://pkg.go.dev/container/heap) package. Call `heap.Push`
and `heap.Pop` to maintain heap order. The type's own
`Push` and `Pop` methods are callbacks used by that package; its `Pop` removes
the final slice element after the package has moved the selected item there.

## Practice and verification

This project also supports conversation practice: discuss the state, helper
functions, branches, and the example traces with the tutor. Writing syntax and
running code are optional. Scaffolding this exercise does not record lesson
completion.

For an agent or maintainer working in this directory, the checks are:

| Command | Purpose |
| --- | --- |
| `go test -run '^$' ./...` | Compile the implementation and tests without running cases |
| `go test -run '^TestSamples$' -v ./...` | Evaluate the two examples |
| `go test -v ./...` | Evaluate all supplied behavioral cases |

The test suite checks expiration boundaries, lifetime changes, zero duration,
counting without prior reads, independent caches, and large timestamps.
The reference implementation is expected to pass all supplied cases.
