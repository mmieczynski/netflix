Given the patterns above, I'd work through these **in this order**. The goal isn't just to solve them: for the later exercises, practise handling follow-up requirements without rewriting everything.

### Phase 1 — Core LeetCode patterns

| # | Exercise | Main skill |
|---|---|---|
| 1 | Two Sum | HashMap |
| 2 | Group Anagrams | HashMap / strings |
| 3 | Longest Substring Without Repeating Characters | Sliding window |
| 4 | Product of Array Except Self | Arrays |
| 5 | 3Sum | Sorting + two pointers |
| 6 | Merge Intervals | Sorting / intervals |
| 7 | Binary Search + Search in Rotated Sorted Array | Binary search |
| 8 | Top K Frequent Elements | Heap / hashmap |
| 9 | Binary Tree Level Order Traversal | BFS |
| 10 | Validate Binary Search Tree | Trees / recursion |
| 11 | Number of Islands | Graph / DFS-BFS |
| 12 | Course Schedule | Graph / cycle detection |
| 13 | Clone Graph | Graph + hashmap |
| 14 | Kth Largest Element | Heap |
| 15 | LRU Cache | HashMap + linked list |

For #15, **don't use your language's built-in ordered dictionary**. Implement the hashmap + doubly linked list yourself.

At this point, I'd expect you to be comfortable with most common Medium questions in roughly **25–35 minutes**.

---

## Phase 2 — The important Netflix-style exercises

This is where I'd spend serious time.

### 16. TTL Cache

Implement:

```typescript
class Cache<K, V> {
    put(key: K, value: V, ttlMs: number): void
    get(key: K): V | undefined
}
```

Then have these follow-ups ready:

**A.** What happens when an item expires?

**B.** Expired entries currently consume memory. Fix that.

**C.** `get()` must remain approximately O(1).

**D.** What if there are 10 million entries?

**E.** How would you test expiration without `sleep()`?

That last question is particularly useful: inject a `Clock`/time provider.

---

### 17. LRU + TTL Cache

Combine #15 and #16.

```text
capacity = 3

put("A", 1, 10s)
put("B", 2, 60s)
get("A")
put("C", 3, 20s)
put("D", 4, 20s)
```

Now determine what gets evicted.

Then consider:

> What if the LRU entry has already expired?

You need to reason about **two independent ordering mechanisms**:

```text
LRU ordering
A ←→ C ←→ D

Expiration ordering
A: 10:01
D: 10:03
C: 10:05
```

This is excellent interview practice.

---

### 18. Weighted LRU Cache

Instead of:

```typescript
capacity = 100 // objects
```

you have:

```typescript
capacity = 1_000_000 // bytes
```

and:

```typescript
put(key, value, size)
```

Evict LRU objects until:

```text
totalSize <= capacity
```

Then ask yourself what happens if **one object is larger than the entire cache**.

---

### 19. Thread-safe cache

Take #17.

Now assume:

```text
Thread A → get()
Thread B → put()
Thread C → cleanup()
```

Explain where races can occur.

You don't necessarily need a perfect lock-free implementation. You should be able to discuss:

- mutexes/locks
- lock granularity
- atomic operations
- contention
- deadlocks
- concurrent collections.

---

### 20. Sliding-window rate limiter

Implement:

```typescript
allow(userId: string): boolean
```

Rule:

> Maximum 100 requests per user during any rolling 60-second period.

Start with:

```text
Map<UserId, Queue<Timestamp>>
```

Then follow up:

> We have 20 million users.

Then:

> Requests arrive on multiple servers.

Then:

> How accurate does the limit actually need to be?

This transitions beautifully from coding → system design.

---

### 21. Token-bucket rate limiter

Implement a token bucket with:

```text
capacity = 100
refill = 10 tokens / second
```

Be able to explain why this behaves differently from #20.

Then make it thread-safe.

---

### 22. Time-based key-value store

API:

```typescript
set(key, value, timestamp)

get(key, timestamp)
```

`get()` returns the newest value at or before the requested timestamp.

Example:

```text
set("quality", "1080p", 10)
set("quality", "4k",    20)

get("quality", 15) → "1080p"
get("quality", 25) → "4k"
```

The interesting solution involves storing ordered values per key and using **binary search**.

There's also a LeetCode problem for this: **981. Time Based Key-Value Store**.

---

### 23. Event deduplication

Netflix receives:

```typescript
{
    eventId,
    userId,
    type,
    timestamp
}
```

Multiple identical events can arrive.

Implement:

```typescript
process(event)
```

such that the same `eventId` is processed only once.

Then:

> IDs only need to be remembered for 10 minutes.

Then:

> You're processing 1M events/sec.

Then:

> Events are processed across 50 servers.

Again: coding → data structures → distributed systems.

---

### 24. Hit counter

Implement:

```typescript
record(timestamp)
count(lastNSeconds)
```

For example:

```text
record(100)
record(101)
record(101)
record(150)

count(last 60 seconds)
```

Try implementations using:

- queue
- circular buffer
- timestamp → count aggregation.

Compare memory/time trade-offs.

---

### 25. Task scheduler

Given:

```typescript
schedule(task, executeAt)
```

and:

```typescript
runReadyTasks()
```

execute everything whose scheduled time has arrived.

Start with a **priority queue/min-heap**.

Then:

> A task can be cancelled.

Then:

> Tasks can be rescheduled.

Then:

> Multiple workers execute tasks.

Then:

> Workers can crash.

You don't have to implement the distributed version. Explain how your architecture would evolve.

---

## Phase 3 — Larger practical coding

These should take roughly **45–60 minutes**, rather than 20-minute LeetCode sessions.

### 26. In-memory filesystem

Implement:

```typescript
mkdir("/movies/action")
write("/movies/action/movie.txt", "...")
read("/movies/action/movie.txt")
ls("/movies")
delete("/movies/action/movie.txt")
```

Think about whether your internal representation should be a tree.

Then add:

```text
move()
copy()
```

And ask what happens with concurrent modification.

---

### 27. In-memory database

Start with:

```typescript
set(key, value)
get(key)
delete(key)
```

Add transactions:

```text
BEGIN
SET A 10
SET B 20
ROLLBACK
```

Then:

```text
BEGIN
SET A 10

    BEGIN
    SET A 20
    COMMIT

ROLLBACK
```

Nested transactions make this considerably more interesting.

---

### 28. Pub/sub

Implement:

```typescript
subscribe(topic, callback)
unsubscribe(topic, callback)
publish(topic, message)
```

Then ask:

> What if a subscriber is slow?

> What if a callback throws?

> Can subscribers modify subscriptions while `publish()` is executing?

> What happens concurrently?

This tests actual software engineering rather than algorithm recognition.

---

### 29. Streaming statistics

Events arrive continuously:

```text
(videoId, watchDuration)
```

Expose:

```typescript
average(videoId)
topKVideos(k)
```

Then:

> Only events from the last hour count.

Now you need windowing/expiration.

Then:

> There are 100M videos.

Now memory matters.

---

### 30. Mini Netflix playback service

Finally I'd combine several ideas.

Implement a simplified service:

```typescript
getPlayback(userId, videoId)
```

You have an expensive metadata service:

```typescript
metadataService.get(videoId)
```

and want to cache responses.

Requirements arrive incrementally:

```text
1. Cache metadata.

2. Entries expire after 5 minutes.

3. Cache has a memory limit.

4. Popular videos should remain cached.

5. Multiple requests for the same uncached
   video can arrive simultaneously.

6. Metadata service sometimes fails.

7. Collect cache hit/miss metrics.
```

That forces you to think about:

**cache stampedes → TTL → eviction → concurrency → failures → observability.**

That's exactly the transition from "can you code?" to "can you engineer?" that I'd want to practise.

---

## How I'd actually use these

Don't do all 30 exactly once.

For **#1–15**, solve them in your interview language under a timer. If you can't solve one, study the solution and then redo it **2–3 days later from scratch**.

For **#16–25**, don't look up solutions immediately. Give yourself ~45 minutes and actually write runnable code and tests. After you've implemented the basic requirement, use the follow-ups above as though an interviewer had just interrupted you.

For **#26–30**, I'd do perhaps **one per day** and spend 60–90 minutes on it.

Most importantly, I'd repeatedly redo **#16–20**:

> **TTL Cache → LRU → LRU+TTL → concurrent cache → rate limiter**

If your Netflix interview were **two weeks away**, those five exercises would probably get more of my attention than grinding another 50 random LeetCode questions.

And if you tell me **your Netflix role (backend/frontend/full-stack), level (L4/L5/etc.), interview date, and language you'll use**, I can narrow these 30 down into a day-by-day preparation schedule and tell you exactly which ones I'd skip.