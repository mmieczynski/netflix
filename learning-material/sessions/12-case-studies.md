# Chapter 12 - Complete interview case studies

These cases combine the book's ideas into programs that can be designed entirely in words. Read the brief, pause to reason if you wish, then follow the complete worked design. In conversation, the tutor can withhold the worked answer and introduce the follow-ups one at a time. There is no requirement to type or execute code.

## Case one: a metadata cache with recency and expiration

The brief is to support Put(key, value, lifetime) and Get(key) for video metadata, with a maximum number of resident entries. Reads must never return expired values. Successful reads and writes update recency. When there is no space, remove expired entries before evicting a live least-recently-used entry. A clock is provided so time can be controlled.

### Resolve the observable rules

Choose these explicit rules for the case. Capacity is nonnegative; zero stores nothing. At the exact deadline, an entry is expired. A nonpositive lifetime removes the key. Overwriting a key changes its value, resets its deadline, and makes it most recent. Get does not extend the deadline. A missing Get does not change recency. A live value can be empty, so Get returns a found flag separately from the value.

These questions matter because they change branches and bookkeeping. They are not a request to design a globally distributed cache. Start with one process and one thread, then add concurrency after the state transitions are correct.

### Build a correct baseline first

A map of entries plus a recency list can implement the baseline. Before a capacity-sensitive Put, scan the map and remove every entry whose deadline has passed. This is linear cleanup, but it is correct and easy to explain. Get checks only its own entry's deadline, so it need not scan the map. Identify cleanup as the expensive operation before introducing a heap to optimize it.

For the optimized version, the cache struct contains the key-to-node map, a doubly linked recency list, a deadline min heap, a generation sequence, capacity, and the clock. Each node contains key, value, deadline, generation, and previous/next links. Heap records contain key, deadline, and generation. The same live node participates in identity and recency; the heap can retain stale records that are checked before they act.

### Define helpers around invariants

Detach joins neighboring list nodes. InsertFront adds a detached node immediately after the head marker. Remove unlinks a real node and deletes its map entry. If weights are later added, Remove also subtracts its weight exactly once. DrainExpired reads due heap records and removes a current node only when its generation matches.

Every current map entry must correspond to exactly one real list node. Every real node must have one map entry under its key. Neighbor links must agree. No expired node may be returned. On return from Put, resident count must be within capacity. Heap records are not required to correspond one-to-one with map entries in this lazy design, but stale records must never remove a newer generation.

### Describe Get as an executable sequence

Look up the key. If absent, return missing. If present, read now and compare with the node's deadline. At or beyond the deadline, Remove the node and return missing. Otherwise Detach it, InsertFront it, and return its value with found true. A stale heap record may remain after lazy removal; later cleanup will see that its key is absent or its generation changed and ignore it.

This Get has expected constant map and list work. It does not perform global expiry cleanup. The guarantee is fresh visibility, not that a Get for one key instantly reclaims every other expired key.

### Describe Put and the order of its changes

For zero capacity, retain nothing. For a nonpositive lifetime, remove an existing key and return. For a positive lifetime, read now and drain all due expirations before making a live capacity eviction. Find the key after cleanup. If it exists, update its value and deadline, assign a fresh generation, and promote it. If it does not, create and register a new node at the front. Push a matching deadline record. If count exceeds capacity, Remove the least-recent node at the tail.

After draining due records, count-limited insertion can require at most one live eviction. Draining itself can remove many records. Therefore the complete Put does not have a constant worst-case bound. Honest complexity names the work performed, including old heap records that do not delete anything.

### Walk a boundary case

Capacity is two. At time zero, Put A with lifetime five and B with lifetime twenty. Recency is B, A. Get A makes it A, B. Advance to five and Put C with lifetime twenty. A has expired even though it is most recent. Cleanup removes A first. C can then be inserted without evicting B, leaving C, B.

At time six, overwrite C with lifetime fifty, so its new deadline is fifty-six. At time twenty-five, the old C deadline record is due. Its generation no longer matches, so the new C survives. The same key name is not enough to authorize deletion.

### Add concurrency and loading carefully

A first concurrency extension uses one mutex around each complete public operation. Get takes an exclusive lock because it promotes or removes nodes. The lock protects map/list/generation consistency, not just the map. Helpers assume the caller holds the lock rather than locking recursively.

If metadata is loaded from a slow service on a miss, do not hold this global mutex through the service call. Create or join per-key in-flight work under the lock, release the lock while fetching, then publish or complete the result under a protected transition. Release waiters on both success and failure. If invalidation occurred during the fetch, a separate generation check can prevent an obsolete result from repopulating the key.

### Tests you can reason through without running

Capacity zero must remain empty. A stored empty value must still be a hit. Expiry exactly at the deadline must miss. Refresh before an old deadline must survive old cleanup. Overwrite must not create duplicate list nodes. An expired most-recent entry must be removed before a live least-recent eviction under this policy. Repeated refreshes must be included in memory analysis because their stale heap records may accumulate.

**Reflection:** Which of these tests checks visibility, which checks structural consistency, and which checks lifecycle identity? Separating those categories makes it easier to find a missing invariant when a new follow-up arrives.

## Case two: recent viewing totals and ranked recommendations

The brief is to accept a finite collection of viewing events and return the top k titles by total watch duration in a recent window. Each event has event ID, user ID, title ID, duration, and event timestamp. Redeliveries with the same event ID count once. Tied totals rank by ascending title ID. The caller supplies now and window width.

### Define the batch contract

Use the interval strictly after now minus width and at or before now. Future events are excluded. The first version assumes repeated IDs have identical payloads. Durations are positive; choose to reject invalid duration records rather than let them silently affect totals. k at or below zero gives an empty result. Input order is arbitrary and preserved. This is a batch query, so event order need not support a queue.

The program separates validation, filtering, identity, aggregation, and ranking. This separation helps you explain a change without rewriting the whole solution. If the prompt wants invalid records skipped instead of rejected, that changes validation policy while leaving the valid-data aggregation argument intact.

### Choose fields and helper responsibilities

A seen-event set suppresses duplicates. A totals map accumulates duration by title. A result record contains title ID and total duration. A comparison helper orders larger totals first, then smaller IDs. The main function validates records, applies the window rule, suppresses duplicates among relevant identical deliveries, and adds each accepted contribution to the totals map.

Because duplicate payloads are assumed identical, filtering an out-of-window delivery before deduplication cannot hide an in-window version of that same event. If duplicate IDs can have different timestamps or values, this reasoning fails. You then need a canonical record policy, conflict rejection, or a defined correction model. State that changed requirement before modifying the code structure.

### Use a simple ranking baseline

Collect the totals map into a slice of result records. Sort by the complete comparison and return up to k. If n is events and u is distinct qualifying titles, expected scan work is linear in n, sorting is u log u, and retained identity plus aggregation state can be linear in n plus u. A map being fast does not make its memory free.

For small k, use a heap containing at most k winners and exposing the worst retained one. Compare every unique title's final total with that root. Sort the surviving winners before returning. This saves ranking work but does not remove the seen-event set or totals map. Do not claim the entire algorithm now needs only k memory.

### Execute the example mentally

At now ten with width five, event e1 gives A three duration at time six. Event e2 gives A four at nine and is delivered twice. Event e3 gives B seven at ten. Event e4 gives C one hundred at five. The interval excludes time five, so C contributes nothing. A totals seven because e2 counts once. B also totals seven. A wins the ID tie, so the ranking is A followed by B.

An event at eleven is in the future and is excluded. An empty collection gives no titles. k larger than the number of qualifying titles returns all of them without inventing missing positions. These observations follow directly from the contract and can be checked without executing a program.

### Change the input to a stream

Now events arrive in nondecreasing timestamp order, and queries also move forward in time. Retain a queue of accepted contributions. Before answering at now, remove contributions at or before the cutoff and subtract their durations from their titles' aggregates. Maintain counts as well as sums if needed to distinguish a title with no remaining events. Delete an aggregate when its last retained event leaves.

The queue works because expiry order matches event-time order. If late events are allowed, insertion order no longer proves that expired contributions form a prefix. Define a bounded-lateness policy, an ordered structure, or a separate event-time processing rule. Do not silently accept arbitrary late events into a FIFO queue.

Deduplication retention is not automatically identical to the ranking window. If the sender can redeliver an old event after its identity marker has expired, the program needs to decide whether the event is still eligible and whether it can be counted again. For corrections with reused IDs, you may need to remove an old contribution before adding its replacement. Identity, event eligibility, and processing lifetime remain separate decisions.

### Understand score decreases

Suppose A is the current winner and B is second. A's oldest contribution expires, reducing its score below B. A heap containing only A cannot discover B if B has been discarded from all other state. Keep full active aggregates and rank on query for a correct baseline, or maintain an indexed ordered representation covering all candidates. A static top-k selection proof cannot be reused unchanged for scores that decrease.

### Explain the final program to an interviewer

Begin with the batch contract and the simplest correct scan-and-sort program. State which map captures identity and which captures aggregation. Show the boundary example and deterministic tie rule. Then explain the small-k optimization and its true memory bound. For streaming, introduce reversible contributions, ordered expiry, and the limitation of winner-only state. Each extension has a specific changed operation and invariant.

## A review you can do in conversation or while reading

For either case, check whether you can state the contract without ambiguity, justify each structure by a required operation, explain an invariant, trace an exact boundary, and adapt to one changed requirement. If not, the worked answer shows where to revisit. No numeric score or fixed time is needed.

An oral mock can stop before the worked answer and let a tutor reveal requirements gradually. An independent reader can simply follow the design from start to finish. Both modes develop the same reasoning. Neither claims that a spoken program has been compiled or that the case is an actual Netflix interview question.
