# Chapter 11 — Modeling larger practical programs

A practical prompt becomes manageable when you separate observable behavior, state, and transitions.

## An in-memory filesystem

Directories have named children; files have contents. A tree mirrors this relationship. Resolve a path by splitting it into components and following child maps. Lookup costs O(number of components), plus string-processing cost. Listing a directory is proportional to its children, with additional sorting if the output must be lexical.

Before implementation, define whether repeated separators are allowed, whether dot and dot-dot are normalized, whether missing parents are created, and whether a file can be overwritten by a directory. You can deliberately restrict the first version to absolute normalized paths. State that restriction instead of silently implementing incomplete path semantics.

Move is not merely deleting one map entry and inserting another. Reject moving a directory inside its own descendant; otherwise you introduce a cycle or detach an entire subtree incorrectly. Validate destination rules before mutating either parent so failed moves leave the filesystem unchanged. Copy requires new nodes; sharing the old subtree would make edits to one copy change the other.

## A transactional key-value store

Start with a base map and a stack of transaction overlays. Reads search overlays newest to oldest, then the base. Each overlay stores either a value or a tombstone representing deletion. A missing overlay entry means “keep searching”; a tombstone means “the key is absent.” Confusing those two states makes deleted values reappear from lower layers.

Begin pushes an empty overlay. Rollback discards the top overlay. Commit merges it into its parent overlay, or into the base if no parent exists. For nested transactions, an inner commit does not necessarily make changes durable outside the outer transaction. Our contract keeps them within the outer transaction until that transaction commits.

Trace base A equals five. Begin outer and set A to ten. Begin inner and set A to twenty. Commit inner: the outer overlay now holds twenty. Roll back outer: base A is five again. If your implementation committed the inner overlay directly to base, it violates this contract. Reads cost O(depth) map lookups; committing costs proportional to the top overlay's changes.

**Think aloud.** Why can an overlay not represent deletion by simply removing the key from that overlay? Explain with a key that exists in the base map.

#### Worked answer — Coach notes

Removing the overlay entry exposes the older base value. A tombstone must stop lookup and report absence. This is the same distinction between “no update here” and “an explicit deletion” that appears in temporal storage and merge systems.

## Pub/sub and callback ownership

Use topic to subscriber-ID to callback maps. Return a subscriber ID from Subscribe; arbitrary Go function values cannot be compared to each other for equality, so callback identity alone is not a suitable map key. On Publish, copy the current subscriber callbacks while holding the lock, release it, then invoke the copy. This allows callbacks to subscribe or unsubscribe without deadlocking or invalidating map iteration.

Define snapshot semantics: a subscriber removed during a publication may still receive that publication if it was in the snapshot. A synchronous slow callback delays later callbacks. Asynchronous delivery needs bounded buffers and a backpressure or drop policy; spawning unbounded goroutines transfers the problem into memory growth. Decide whether callback panics propagate or are recovered at the API boundary. These are behavioral choices with tests, not just implementation details.

## A playback metadata service

Compose the previous lessons around GetPlayback. Validate user and video inputs. Check metadata cache. On a fresh hit, return the cached metadata. On a miss, join or initiate the per-key load. If loading succeeds, publish it with TTL and capacity accounting. If it fails, return an error or an explicitly allowed stale value. Expiration, eviction, and in-flight loading have different purposes.

Define metrics precisely. One request can be a cache miss while joining an already-running load, so backend-load count need not equal miss count. A stale value served under an outage policy needs its own metric if you care about freshness. Do not cache authorization decisions under video ID alone; user-specific data requires an appropriate key or a separate authorization check.

What if an object is larger than the memory budget? Return the successful backend result to the caller while declining to cache it under our chosen policy. Cache admission failure need not mean playback failure. What if twenty callers cancel while one remains? Shared load lifetime and individual waiting lifetimes are separate decisions.

## Model a small program as a sequence of valid transitions

A larger coding prompt often overwhelms because it introduces several methods at once. Reduce it to a state machine: what is the state before a call, which preconditions allow the call, which mutations happen, and what must be true when it returns? You do not need formal notation. You need an operation sequence precise enough that another person could act as the computer.

For a filesystem move, the state consists of nodes and parent-child relationships. A successful move changes two relationships: remove the source from its old parent and add it under the destination parent. But checking whether the destination exists, whether the new name conflicts, and whether a cycle would be created should normally happen first. Otherwise a failed validation can leave the source detached. This is the same “validate before destructive mutation” principle used for oversized cache admission.

A transaction overlay illustrates a different modeling skill: representing absence separately from no change. If an upper layer does not mention key A, read the lower layer. If the upper layer says A was deleted, stop and report missing. Both situations lack a normal value in the upper layer, but they have different behavior. A tombstone preserves that distinction. Many subtle bugs are really failures to represent two semantically different states separately.

Callbacks make ownership and reentrancy visible. A pub/sub system owns its subscriber registry, but a callback is arbitrary user-controlled work. Holding the registry lock while executing that callback extends the lock across code whose duration and behavior you do not control. Taking a snapshot of subscribers separates the protected registry operation from later delivery. The snapshot policy then becomes part of the contract: a subscriber removed after the snapshot may still receive the current message.

Composition adds another dimension. A playback service can validate access, look up metadata, join an in-flight load, and update metrics. Keeping each responsibility explicit helps you reason about errors. A successful backend response that is too large to cache can still be returned. A cache miss that joins another caller's load need not cause another backend request. “Request failed,” “cache did not admit,” and “load was shared” should not collapse into one undifferentiated status.

## Questions and worked explanations on API design and state transitions

**Round 1.** You are asked for mkdir, write, read, and list operations. What belongs in a directory node, and what belongs in a file node? What invalid combination should the representation prevent?

#### Worked answer — Explanation and next challenge

A directory needs a map of child names to nodes; a file needs content. A node type or explicit kind should distinguish them so file content and directory traversal are not confused. Define whether a path component can be a file when traversal expects a directory. A representation should make invalid states difficult to create or easy to reject.

**Round 2.** Before removing a source node during Move, what destination checks should happen? Why is their order observable if the move fails?

#### Worked answer — Explanation and next challenge

Resolve the destination parent, check its type, decide conflict behavior, and reject moving a directory inside itself or a descendant. If you detach first and then discover failure, the filesystem may lose a valid path. Validate first, then perform the small mutation sequence that constitutes success.

**Round 3.** Copying a directory simply points a new parent entry at the old directory node. What later behavior reveals that this was an alias rather than an independent copy?

#### Worked answer — Explanation and next challenge

Editing or deleting a child through one path affects the other because both paths reach the same object. An independent copy requires new nodes throughout the copied subtree. Shared immutable data or copy-on-write can be valid designs, but they need an explicit sharing contract.

**Round 4.** Base storage contains A equals five. A transaction deletes A. Why can the transaction's overlay not represent that by simply omitting A?

#### Worked answer — Explanation and next challenge

Omission means no update at that layer, so lookup would continue to the base and return five. A tombstone means an explicit deletion and stops the search. Describe the difference as “no statement about A” versus “A must be absent.” That distinction recurs in versioned databases and layered configuration.

**Round 5.** Base A is five. Outer transaction sets ten. Inner transaction sets twenty and commits. Outer transaction then rolls back. What is A under nested-overlay semantics?

#### Worked answer — Explanation and next challenge

A is five. Inner commit merges into the outer overlay, not directly into the base. The outer rollback discards all changes still contained in that transaction. If you instead want inner commit to escape outer rollback, that is a different contract and must be stated explicitly.

**Round 6.** Describe Get across three nested transaction layers. How does a tombstone change the stopping condition, and what is the cost in transaction depth?

#### Worked answer — Explanation and next challenge

Search newest overlay first. A value returns immediately; a tombstone returns missing immediately; no entry means continue downward. If no overlay mentions the key, consult the base. In the worst case, one map lookup per layer makes cost proportional to depth. A fast map at each layer does not make the whole layered operation constant time.

**Round 7.** A pub/sub API returns a subscription ID. Why is that often a better cancellation handle in Go than trying to compare callback functions?

#### Worked answer — Explanation and next challenge

Arbitrary function values are not generally equality-comparable to each other in Go. An explicit ID gives stable identity and supports a subscriber map. It also distinguishes two registrations of the same conceptual callback, if duplicate subscriptions are allowed. Identity should be modeled instead of inferred from executable code.

**Round 8.** Publish snapshots subscribers A and B, then A's callback unsubscribes B. Should B still receive the current message under snapshot semantics?

#### Worked answer — Explanation and next challenge

Yes, because B was part of the snapshot for that publication. It will not be in later snapshots. Other cancellation semantics are possible but require a further validity check or coordination. The important point is that behavior during mutation must be defined rather than depend on map-iteration accidents.

**Round 9.** A slow subscriber blocks all later synchronous callbacks. Would starting an unlimited goroutine for every delivery eliminate all problems?

#### Worked answer — Explanation and next challenge

It shifts the problem to potentially unbounded pending work, memory use, ordering, and error handling. A bounded queue with a defined block, drop, or disconnect policy makes resource behavior explicit. Concurrency is a tool for managing execution, not a substitute for a load policy.

**Round 10.** A playback metadata fetch succeeds, but the result exceeds the cache's byte budget. Should the playback request necessarily fail? Separate the two decisions.

#### Worked answer — Explanation and next challenge

The caller can receive the successful metadata while the cache declines to retain it. Serving the request and admitting a reusable copy are separate decisions. The error or metric for cache admission should not automatically become a backend failure. Define whether this behavior is acceptable for the API.

**Round 11.** Three requests miss the cache but share one in-flight load. What should cache-miss count and backend-load count each record under request-level misses?

#### Worked answer — Explanation and next challenge

Three misses and one backend load. A request can miss locally without initiating expensive work. If metrics are supposed to explain load reduction, preserve that distinction. Name units and observation points instead of treating all counters as interchangeable indicators of traffic.

**Round 12.** Choose one system from this session and describe a failing operation that must leave state unchanged. Give its preconditions, validation, mutation boundary, and a test that detects partial failure.

#### Worked answer — Explanation and mastery check

For example, a move to a nonexistent directory should report failure while the original source path still resolves to its contents. Validate destination before detaching. A transaction Commit with no active transaction can return an error without modifying the base. The mastery signal is a precise state transition, not merely a list of methods.

Use the worked program design to assemble a small program method by method. Let the tutor play a caller that changes requirements, and identify the smallest set of fields or helpers that must change.

## Worked program design: nested transactions

The database exposes Begin, Get, Set, Delete, Commit, and Rollback. Outside a transaction, writes affect a base map. Inside a transaction, they affect only the newest overlay. An overlay entry is either a value or a deletion marker. Commit merges the newest overlay into its parent if one exists, otherwise into the base. Commit or Rollback with no active transaction returns an error without mutation.

Begin pushes an empty map onto a stack. Set stores a value entry in the top overlay, or the base when the stack is empty. Delete stores a tombstone in the top overlay, or removes the base entry outside a transaction. Get searches overlays newest to oldest. A value returns immediately, a tombstone returns missing, and no entry continues the search. Only after all overlays are silent does it read the base.

Commit first checks that an overlay exists. Remove that overlay from the stack, then merge each change into the new top layer. Copy both values and tombstones into a parent overlay so deletions remain explicit. When committing into the base, a tombstone performs an actual delete and a value performs an assignment. Rollback simply discards the newest overlay. These operations define nesting without needing to copy the entire database at each Begin.

With base A equal to five, Begin then Set A to ten affects only the outer overlay. Begin again and Set A to twenty affects the inner one. Commit inner moves twenty into outer. Rollback outer discards it, leaving base five. If inner Delete A is committed instead, the outer overlay must contain a tombstone; merely deleting A from the overlay would reveal the original five too early.

Get costs up to one map lookup per active depth. Commit work is proportional to the number of changes in the top overlay. This is a single-threaded nesting model, not an implementation of all database isolation levels. Adding concurrent transactions requires a visibility and conflict policy beyond simply putting a mutex around each method.

**Optional spoken walkthrough:** describe every public method's target layer and the three-way Get decision. Have the tutor insert a deletion into the nested example and account for the value visible after every operation.
