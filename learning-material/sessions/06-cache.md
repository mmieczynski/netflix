# Chapter 6 — LRU and weighted caches

A cache exercise is about coordinating identity, order, and capacity without letting those views disagree.

## Begin with a contract small enough to implement

Our first cache stores string keys and integer values. Get returns a value and a found flag. Put inserts or overwrites. Successful Get and Put both make the entry most recently used. The capacity is the maximum number of entries. Zero capacity stores nothing. A missing Get does not change order. Expiration and concurrency are separate follow-ups.

LRU means least recently used. It is not least frequently used. A title accessed a thousand times yesterday can be evicted before one accessed once just now. This policy is a useful heuristic for temporal locality, not a guarantee that the most popular items survive. A one-time scan of many unique items can displace frequently reused entries.

## Derive the two structures

A map locates a value by key quickly but does not identify the least recently used key. A list ordered by recency identifies the oldest entry quickly but cannot locate an arbitrary key quickly. Combine them: map each key to a node in a doubly linked list. Keep the most recent node at the front and least recent at the back. Each node has its key, value, previous pointer, and next pointer.

Why doubly linked? Removing a known node requires access to both neighbors. With only a next pointer, you would generally have to search for the predecessor. With both pointers, detach by joining the neighbors, then insert at the front. Sentinel head and tail nodes make empty-list and endpoint operations follow the same pointer rules.

There are three core invariants. Every map entry points to exactly one real list node. Every real list node has exactly one map entry under its key. Neighbor pointers agree, and the real-node count never exceeds capacity when a public operation returns. Express these before coding; they tell you what every helper must preserve.

## A full trace you can narrate

With capacity two, Put A then Put B gives recency order B, A, spoken from most recent to least recent. Get A detaches A and moves it to the front, giving A, B. Put C inserts C, temporarily giving C, A, B. Capacity is exceeded, so evict B from the tail and delete B from the map. Overwrite A with a new value: update the existing node and promote it, giving A, C. Do not create a second A node.

**Think aloud.** You delete the tail node from the list but forget to remove its map entry. What happens on the next Get for that key? Which invariant would have caught the bug first?

#### Worked answer — Coach notes

The map can return an evicted value and may attempt to detach a node that is no longer linked, corrupting the list. The map-to-list correspondence fails immediately at eviction, before any later lookup. Test internal invariants after operations, not only returned values.

## Implement the helpers before the policy

Write detach, insert-at-front, and remove. Detach updates only list links. Insert-at-front links an already detached node. Remove detaches and deletes from the map. Get then becomes lookup, promote, return. Put becomes update and promote, or allocate and insert, followed by eviction when necessary. Avoid duplicating pointer manipulation in every branch.

Expected Get and Put time is constant for a count-limited cache, assuming fixed-cost keys and values. Space is linear in capacity, including map entries and list nodes. Go's built-in container/list can be appropriate in ordinary code, but your plan specifically calls for implementing the linked list yourself; the provided reference `LRU` does that.

## Weighted capacity changes the accounting

Now each entry has a weight, such as payload bytes. Maintain total weight as an invariant equal to the sum of resident entry weights. On overwrite subtract the old weight and add the new weight. Then evict from the back repeatedly while over budget. One Put can evict many items, so worst-case time is linear in the number of entries removed, not constant. Across an operation sequence, each inserted node can be evicted only after being inserted, which helps amortized analysis.

Choose an oversized-item policy explicitly. In this course, reject a weight greater than the entire budget and leave an existing value under that key unchanged. Check this before mutating. Another reasonable contract removes the old entry but refuses the new one; it must be stated and tested. Negative weights are invalid. Zero-weight entries can defeat a byte-only entry bound, so either disallow them, charge overhead, or use a second count limit.

Suppose budget is ten. Put A weight four, B weight four; Get A; Put C weight five. Before eviction, weight is thirteen and order is C, A, B. Evict B to reach nine. If C instead weighs eleven, reject it before touching A or B. A byte budget for values alone does not bound actual process memory: maps, nodes, allocator overhead, and stale heap records also consume space.

## Build the cache from the two questions it must answer

Before imagining pointers, imagine two separate conversations with the cache. A caller asks, “Do you have movie A?” Capacity management asks, “Which resident movie has gone longest without being used?” The first question is about identity; the second is about order. A map answers identity but has no meaningful recency order. A recency list answers oldest and newest but cannot immediately find an arbitrary movie. The combined solution is a small example of maintaining multiple indexes over the same objects.

Now imagine a node as a card containing a key and value, with a connection to the card before it and the card after it. The map does not need another copy of the value; it can identify the card directly. When the card moves, its identity remains the same, so the map continues to point to it. Moving a value inside a slice, by contrast, can change positional indexes. A stable node reference lets identity and position vary independently.

Detaching a card means making its two neighbors point to each other. Inserting it at the front means placing it between the head marker and the previous first card. You do not search through the list during either operation because the map already found the card and the card knows its neighbors. That is the actual reason promotion is constant time. A linked list without fast key lookup would still require a search; a key map without backward links would not immediately provide the predecessor needed for removal.

Sentinel nodes are deliberately empty boundary markers. An empty list consists of head connected to tail. The first insertion uses the same “insert between two neighbors” rule as later insertions. They reduce exceptional pointer cases; they do not count toward capacity or belong in the key map. A useful spoken invariant is: every real card is reachable exactly once between the markers and is named exactly once by the map.

Weighted capacity adds another view: the total charge for all resident cards. That number is a cached summary of the list contents. If you overwrite an entry, you must remove the old charge before adding the new one. If you evict several entries, subtract each exactly once. This is not merely arithmetic bookkeeping; it is an invariant connecting state representations. The same style of reasoning later governs streaming sums when events expire.

## Questions and worked explanations on LRU state

**Round 1.** A cache holds A and B. You need fast lookup and fast identification of the least recently used entry. Explain why a map alone and a list alone each leave one expensive operation.

#### Worked answer — Explanation and next challenge

A map finds A directly but does not say which key is oldest by access. A recency list exposes an endpoint but finding A inside it may require a scan. Mapping keys to list nodes joins the strengths. Name the missing capability first, then the added structure that supplies it.

**Round 2.** Capacity is two. Put A, then B. Say the order from most recent to least recent. Now read A successfully and say the new order.

#### Worked answer — Explanation and next challenge

After the writes the order is B, A. After reading A it is A, B. A successful read counts as use under this contract. A missing read would not add an entry or alter order. Always state which actions update recency before tracing an LRU cache.

**Round 3.** From most recent A, then B, insert C. Which key leaves and why? Could you decide from the values stored under the keys?

#### Worked answer — Explanation and next challenge

B leaves because it is least recent, leaving C, A. Its stored value is irrelevant to recency unless the policy explicitly uses value. The list remembers operation history that cannot be recovered from the payload. This illustrates why the cache maintains metadata in addition to the cached data.

**Round 4.** Overwrite A with a new value while A is already resident. Why should you normally reuse its existing node instead of adding a second A node?

#### Worked answer — Explanation and next challenge

There must be one real node per key. Reusing and promoting the node preserves that correspondence. A second node could leave an old A reachable in the list while the map identifies only the new A. Later eviction might then remove the wrong map entry or return an old value. Identity must remain unique across both views.

**Round 5.** Without writing pointer syntax, describe removing the middle card B from the list A, B, C. What do A and C need to know afterward?

#### Worked answer — Explanation and next challenge

A's next card becomes C, and C's previous card becomes A. B should no longer be linked into that chain. To promote B, then place it between the head marker and the old first card. Thinking in neighbor relationships is enough to understand the operation before translating it into Go assignments.

**Round 6.** You remove an evicted node from the list but leave its key in the map. What observable bug can occur before any obvious crash?

#### Worked answer — Explanation and next challenge

A later lookup may return an entry that should have been evicted. It can then attempt to promote a detached node and corrupt links. The first violated invariant is that every map entry names a resident list node. Correct output on the previous operation does not prove internal state is valid.

**Round 7.** How would you verbally test the list-map invariant after an arbitrary operation sequence? Give a method, not just individual Get expectations.

#### Worked answer — Explanation and next challenge

Walk from head to tail, verify matching forward and backward links, detect repeated nodes, and ensure each key maps to that exact node. Compare the number of real nodes with map size and capacity. This checks structure, while returned-value tests check observable behavior. A small simple model can independently track expected recency order.

**Round 8.** A cache has a byte budget of ten. A weighs four and B weighs four; A was used most recently. Insert C weighing five. Explain every weight and order update.

#### Worked answer — Explanation and next challenge

Insertion temporarily gives total thirteen and order C, A, B. Remove B and subtract four, leaving total nine and order C, A. Under a byte budget, one insertion may need several evictions. The accounting invariant is the sum of weights of resident entries, not merely the sum of weights ever inserted.

**Round 9.** An incoming item weighs eleven while the entire budget is ten. What policy would you choose, and why should the check happen before unrelated entries are evicted?

#### Worked answer — Explanation and next challenge

One sensible policy rejects admission while preserving the existing cache. Since the item cannot fit even in an empty cache, evicting everything first gains nothing. If it overwrites an existing key, explicitly decide whether the old value remains. A clear failure contract avoids partial destructive updates.

**Round 10.** Does a budget based only on payload bytes bound the process's actual memory? Name concrete missing costs.

#### Worked answer — Explanation and next challenge

Map entries, node pointers, keys, allocation overhead, and auxiliary expiration records also consume memory. Zero-weight payloads can allow an unbounded number of entries unless overhead or count is separately bounded. A logical cache budget is useful but should not be mislabeled as an exact process-memory guarantee.

**Round 11.** A movie was accessed a thousand times yesterday. Another was accessed once just now. Which can LRU evict first, and why is “keep popular movies” an imprecise description of LRU?

#### Worked answer — Explanation and next challenge

LRU can evict the once-popular movie because its last access is older. It measures recency, not frequency. A scan of one-time objects can displace reusable entries. If the requirement is frequency-based retention, a different policy or an admission strategy may be needed; do not quietly change LRU semantics.

**Round 12.** Explain why count-limited Put can be expected constant time while a weighted Put may take longer. Include the work hidden by “evict until it fits.”

#### Worked answer — Explanation and mastery check

A count-limited new insertion exceeds capacity by at most one, so it removes at most one node. A weighted insertion can require many small entries to be removed. Its individual time is proportional to the number evicted, even if aggregate eviction work is favorable across many operations. A senior-level explanation distinguishes individual latency from amortized cost.

Describe the full cache contract and its three linked views: key identity, recency order, and capacity accounting. Ask the tutor to give you a four-operation sequence and justify every resulting state without writing code.

## Worked program design: an LRU cache without syntax

The public methods are Get and Put, with a nonnegative entry capacity. Get returns a value and found flag; a hit updates recency. Put overwrites or inserts and also updates recency. Zero capacity retains nothing. The cache struct contains capacity, a map from key to node, and head and tail sentinel nodes. Each real node contains key, value, previous node, and next node.

Define three helpers. Detach joins a node's two neighbors and clears its own links. InsertFront places a detached node between head and the current first node, repairing both directions. Remove calls Detach and deletes the key from the map. Helpers have preconditions: do not detach a sentinel or a node that is not linked, and do not insert a node that is still linked elsewhere.

Get looks up the key. A miss returns missing without mutation. A hit detaches the node, inserts it at the front, and returns its value. Put first handles zero capacity. If the key exists, update that node's value, promote it, and return. Otherwise allocate a node, register it in the map, and insert it at the front. If capacity is exceeded, remove the node immediately before tail.

For capacity two, Put A then B produces B, A. Get A produces A, B. Put C produces C, A after evicting B. Overwrite A with value zero produces A, C; Get A must return zero with found true. This last observation tests the absence representation independently of the ordering.

The public methods preserve one-to-one map/list membership and valid neighboring links. Every operation has expected constant work under the count capacity. If concurrency is added, protect the entire public transition with one mutex first; Get requires the exclusive lock because promotion mutates links. Do not separately lock each helper in a way that leaves inconsistent intermediate states visible or recursively reacquires the same mutex.

**Optional spoken walkthrough:** let the tutor play a caller and name operations. Describe the exact helper calls and map mutations after each one. Then explain how weighted capacity adds a weight field, a total counter, and a loop that may remove more than one victim.
