# Chapter 7 — Expiration clocks and concurrency

Time introduces a second order. Concurrency makes a sequence of correct individual steps unsafe unless they form one protected operation.

## Expiration is a behavioral rule first

For our TTL cache, Put at time t with positive lifetime d stores an absolute deadline t plus d. Get at time now returns a value only when now is strictly before the deadline. At the deadline it is expired. A nonpositive lifetime removes the key immediately. Updating resets the lifetime. A Get does not extend it. These choices distinguish expire-after-write from expire-after-access.

Use an injected clock so tests control now directly. The learning code uses int64 ticks in one agreed unit and assumes monotone time and no arithmetic overflow. In a real Go process, a clock returning time.Time can preserve monotonic elapsed-time behavior; persisted timestamps need different care. The important point is to make the source and meaning of time explicit, and never rely on real sleeping in unit tests. [Go time documentation](https://pkg.go.dev/time#hdr-Monotonic_Clocks).

A map with per-entry deadlines gives expected constant-time lookup. On Get, delete an expired entry and return missing. This is lazy expiration: it prevents stale reads, but an entry never read again may remain allocated indefinitely. Correct visibility and prompt memory reclamation are different requirements. Start with the simple version, then explain the gap honestly.

## Choose cleanup by the required order

A periodic scan checks every entry in linear time per sweep. It is simple and may be sufficient. A min heap of deadlines makes the next expiry easy to find. Repeatedly pop while the smallest deadline is at or before now. Removing r records costs O(r log h) where h is heap size. A queue only works if expiry order matches insertion order, such as a uniform TTL with monotone writes. Different TTLs break that assumption.

On each Put, store the current value, deadline, and a new generation number in the map; push deadline, key, generation into the heap. On cleanup, delete the map entry only if its generation still matches the popped record. This prevents an old timer from deleting a newer value. Assign generations from a globally increasing sequence, or otherwise avoid reusing an old generation after a key is deleted and recreated.

Trace A written at time zero, expiring at ten, generation one. At time five, overwrite A to expire at thirty, generation two. At time ten, the old heap record is ready. The map holds generation two, so discard the record without deleting A. The same defense applies to cancelled or rescheduled tasks. Reference: `TTLCache`.

**Think aloud.** Why is comparing only the key insufficient? Could comparing only the deadline ever confuse an old and a new entry? What extra invariant does a generation number provide?

#### Worked answer — Coach notes

The same key can identify multiple lifetimes. Equal deadlines may be reused, and future changes can make deadline-only identity fragile. A unique generation identifies the particular write, so cleanup can prove it is acting on the current lifetime.

Lazy invalidation simplifies updates but can accumulate stale heap records. Ten million rewrites of one key can produce ten million heap entries despite a one-entry map. Bound this by using an indexed heap with one live record per key, or periodically rebuild from live entries when stale overhead crosses a threshold. Rebuilding with heap.Init is linear in live entries, but creates a latency spike. “O(number of keys) memory” is not true for an unbounded lazy heap.

## Combining TTL and LRU

The map still answers identity. The linked list answers least recently used. The deadline heap answers earliest expiration. These are independent orders. An expired node can be at the front of the recency list, so removing only expired nodes at the tail does not clean all expired entries.

Choose whether capacity counts all physically stored entries or only unexpired ones. For “evict expired entries before live ones,” a straightforward Put first drains all due expirations, then updates or inserts, then evicts LRU entries if capacity is exceeded. That Put is not constant time when many items expire together. Get can remain approximately constant time by checking only its own key and leaving global cleanup elsewhere.

Removal must consistently update the map, recency list, and weight counter if present. Heap records may remain as harmless stale versions under lazy invalidation. Reuse a central remove helper rather than implementing three subtly different removal paths for eviction, expiration, and explicit deletion.

In your plan's capacity-three example, assume all actions occur before any deadline: Put A, Put B, Get A, Put C, Put D leaves D, C, A and evicts B. If A expires before D is inserted and cleanup removes expired entries first, A leaves and B can remain. Without specified elapsed time and cleanup policy, the original trace has no unique answer. Noticing that ambiguity is part of the skill.

## Thread safety protects the whole invariant

For a first concurrent implementation, put one mutex around each public operation that touches shared state. A successful LRU Get mutates the list, so it needs an exclusive lock. A TTL Get may delete an expired entry, so it is not automatically a read-only operation either. Protecting the map alone while updating pointers outside the lock leaves races and structural corruption.

Do not copy a cache containing a used mutex. Return pointers from constructors and use pointer receivers. Do not try to upgrade an RWMutex read lock into a write lock. A specialized concurrent map would protect only map operations, not the combined map/list/weight invariant. These synchronization rules are grounded in the [Go sync documentation](https://pkg.go.dev/sync) and [memory model](https://go.dev/ref/mem).

Slow loading introduces another concern. Holding the global cache mutex while calling a metadata service blocks unrelated keys. Instead check the cache, coordinate a per-key in-flight call, release the lock, load, then publish the result and wake waiters. Exactly one leader loads a key; other callers wait on that call's completion. Clear the in-flight record on failure as well as success. Decide whether failures are cached, whether stale values may be served, and how cancellations affect shared work.

Single-flight loading prevents duplicate simultaneous work; it does not enforce TTL or capacity. If a key is invalidated while loading, decide whether the late result may repopulate it. An invalidation generation can prevent that. With sharding, locks can reduce contention, but a globally exact LRU order becomes harder; mention the trade-off rather than claiming sharding preserves every property for free.

## Separate freshness, storage, and coordination

Three questions often get bundled together in a cache prompt. Is this value fresh enough to return? Is this entry still taking space? Is another caller already fetching its replacement? TTL answers the first question, cleanup answers the second, and in-flight coordination answers the third. A correct answer to one does not automatically answer the others. An entry can be expired but allocated, and a missing entry can have a fetch already underway.

Use a bookstore analogy carefully. An edition can become outdated at a known time even while its book remains on the shelf. A customer asking for that book must not be handed it as current. Removing it from the shelf can happen during that request or in a separate cleanup pass. Ordering shelves by most recent customer use does not order books by when their editions expire. That is why recency and deadline often require separate structures.

Now consider an expiry message as a delayed instruction: “Delete A at time ten.” If A is replaced before time ten, the instruction lacks enough identity. It should say, “Delete the particular lifetime of A that I was created for.” A generation number provides that identity. Cleanup looks at the current generation and acts only if it matches. This prevents stale work from affecting newer state. The principle applies to rescheduled tasks, cancelled requests, and late background computations.

Concurrency creates an analogous problem in the present rather than the future. Two requests can observe the same apparently available capacity before either records its own change. Each individual read and write may be safe, but their combined decision is not atomic. You need a single protected transition from the old valid state to the new valid state. For a cache, that transition can involve map, list, counters, and deadlines together.

Think about the exact moment at which an operation takes effect. In a simple locked cache, it occurs within the critical section while other callers cannot interleave. If you release the lock halfway through promotion, another operation can see a state that is neither the old order nor the new order. Locks are therefore protecting a semantic invariant, not merely preventing a runtime complaint about concurrent maps.

## Questions and worked explanations on time and overlapping operations

**Round 1.** Put A at time zero with lifetime ten. Under expiry at the deadline, what should Get return at nine and at ten? Why must this be decided before coding?

#### Worked answer — Explanation and next challenge

At nine it is fresh; at ten it is expired. This choice determines the comparison and the tests. If the contract were inclusive at the deadline, the boundary would differ. “After ten seconds” can be ambiguous in ordinary language, so state the exact rule once and use it consistently.

**Round 2.** Get deletes expired entries when they are requested. An expired key is never requested again. Is stale-return correctness satisfied? Is memory reclamation satisfied?

#### Worked answer — Explanation and next challenge

The cache can still guarantee it will not return expired values, because every Get checks freshness. It has not guaranteed prompt physical removal of idle expired keys. These are separate requirements. A periodic scan or deadline-based cleanup addresses storage without changing the visibility rule.

**Round 3.** A is inserted before B, but A lives for one hundred ticks and B for one tick. Why does an insertion-order queue fail to expose the next expiry?

#### Worked answer — Explanation and next challenge

B expires first despite arriving later. If cleanup checks only the queue front, fresh A blocks removal of expired B. A deadline heap handles arbitrary lifetimes. A FIFO queue is valid only when insertion order implies expiry order, such as a uniform lifetime with monotone insertion time.

**Round 4.** A's original lifetime expires at ten. At five, A is replaced with a value expiring at thirty. What must the cleanup record at ten check before deleting anything?

#### Worked answer — Explanation and next challenge

It must identify the write it belongs to and compare that identity with the current entry. A matching key alone is insufficient. If its generation is old, discard only the stale cleanup record. The current value should survive until its own deadline or another valid removal event.

**Round 5.** A key is deleted, then recreated. Its generation counter restarts at one, matching an old heap record. What bug can follow, and how can you avoid it?

#### Worked answer — Explanation and next challenge

The old record can mistake the new lifetime for the old one and delete it. Use generation identities that are not reused while stale records can exist, such as a cache-wide increasing sequence with a documented non-wrap assumption. A version number is useful only if it actually distinguishes the relevant lifetimes.

**Round 6.** You overwrite one key a million times and use a new heap record for every write. Can you honestly call memory proportional only to the number of current keys?

#### Worked answer — Explanation and next challenge

No. The map may have one entry while the heap retains many stale records. Use a heap index to update in place or rebuild periodically from live entries, and account for the work and pause introduced by that choice. Lazy invalidation trades update simplicity for cleanup and storage obligations.

**Round 7.** In an LRU-plus-TTL cache, the most recently accessed entry expires first. Can inspecting only the least-recent endpoint remove all expired entries?

#### Worked answer — Explanation and next challenge

No. Recency order and expiry order are independent. The expired entry can be anywhere in the recency list. If capacity policy requires removing expired entries before any live eviction, Put needs access to all due expirations or a scan. Do not infer deadline order from recency.

**Round 8.** Why can a method named Get require an exclusive lock in an LRU cache? Give two state changes a cache lookup may make.

#### Worked answer — Explanation and next challenge

A successful lookup may promote the node, mutating list links. An expired lookup may delete the entry and update counters. Method names do not determine read-only behavior. Protect the actual state transition, and do not assume a read lock is enough simply because the caller receives a value.

**Round 9.** Two callers both observe a missing key and then start loading metadata. Why does locking the individual Get and Put methods not automatically prevent duplicate loads?

#### Worked answer — Explanation and next challenge

The shared decision spans multiple operations, with the load between them. Both can complete Get before either performs Put. Introduce per-key in-flight state so one caller becomes leader and others join that work. Protect creation and publication of that state, but avoid holding a global lock throughout the slow load.

**Round 10.** A leader's metadata fetch fails. What should happen to the waiters and the in-flight marker? Why is the failure path part of the concurrency design?

#### Worked answer — Explanation and next challenge

Waiters must be released with a defined result, typically the error or an allowed stale value, and the marker must be cleared or completed so future requests can retry. Leaving it active can make callers wait forever. Single-flight coordination must complete on success, failure, and any supported cancellation path.

**Round 11.** A key is invalidated while a background fetch is running. Later the fetch succeeds. Should the result automatically repopulate the cache? Explain one defensible contract and how to enforce it.

#### Worked answer — Explanation and next challenge

If invalidation means older in-flight work must not repopulate the key, capture an invalidation generation when loading starts and compare it before publication. If the generation changed, do not cache the old result. Other contracts are possible, but the late-result behavior must be explicit. This is again stale work acting on newer state.

**Round 12.** Describe a test for expiry after refresh without sleeping. Include the controllable state, operations, expected values, and the bug it catches.

#### Worked answer — Explanation and mastery check

Inject a clock initially at zero. Put A expiring at ten; advance to five and overwrite A to expire at thirty; advance to ten and run cleanup; Get must still return the new value. Advance to thirty and it must be missing. The test is deterministic and catches old expiry records deleting refreshed values. It tests a causal sequence rather than hoping a real-time delay is long enough.

Explain one invariant protected by a lock and one identity protected by a generation number. Then explain why neither mechanism automatically bounds memory or prevents slow backend calls.

## Worked program design: TTL storage with safe cleanup

Expose Put(key, value, lifetime), Get(key), and Cleanup. A clock is supplied to the constructor. The store owns a map of current entries, a min heap of expiry records, and a generation sequence. Each current entry contains value, deadline, and generation. Each heap record contains key, deadline, and generation. Time units are consistent, deadlines fit the numeric type, and generation identities are not reused during relevant lifetimes.

Put with a nonpositive lifetime removes the current map entry. Otherwise it reads the clock, computes the deadline, allocates a fresh generation, stores the current entry, and pushes an expiry record. An older record for that key can remain in the heap. Get returns missing if absent. If the clock is at or past the current deadline, it deletes the map entry and returns missing. Otherwise it returns the value. Get does not drain the whole heap.

Cleanup reads now once for its cutoff. While the heap root is due, pop it. Look up its key in the current map. If absent, do nothing further. If present with a different generation, discard this stale record. Only a matching current generation may be deleted. A returned removal count can count current entries actually removed, rather than all stale heap records examined.

Write A at zero with deadline ten. At five, overwrite it with deadline thirty and a new generation. Cleanup at ten pops the old record but leaves the map unchanged. Cleanup at thirty removes the new entry. Deleting and recreating A between those calls is safe only if recreation does not reuse the old generation.

Put costs logarithmic time in heap size; Get is expected constant time; Cleanup costs according to the number of records popped. Heap size includes stale records, so repeated overwrites can grow memory. An indexed heap or periodic rebuilding addresses that extension. Adding a single mutex around public operations protects state, provided the injected clock is safe and does not reenter this cache.

**Optional spoken walkthrough:** explain all three Cleanup branches after a pop: absent key, mismatched generation, and matching generation. Then describe a fake-clock test that distinguishes expiry after write from expiry after access.
