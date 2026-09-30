# Chapter 9 — Deduplication counters versions and scheduling

Many practical tasks reuse the same ideas: key identity, chronological removal, predecessor lookup, and versioned invalidation.

## Deduplication needs an identity and a horizon

An event has an event ID, user ID, title ID, duration, and event time. To suppress repeat deliveries, first define identity. The same event ID identifies the same logical event in our exercise; two different IDs with identical payloads are separate events. Store accepted IDs in a set. A repeated ID is rejected, while an unseen ID is recorded and accepted.

With ten-minute retention, use a deadline per ID. Decide whether a duplicate refreshes the deadline. Here it does not: expiration is measured from the first accepted processing time. A repeat just before expiry therefore does not suppress that ID forever. Use server processing time for this exercise, not an untrusted event timestamp. A TTL map with a cleanup heap implements the contract directly.

Now consider a callback that performs the actual side effect. If you mark the ID before the callback and crash, a retry can be suppressed even though the effect never completed. If you mark it afterward and crash in between, the effect can happen twice. A local set is duplicate suppression during its lifetime, not a proof of durable exactly-once effects. Stronger behavior needs an idempotent destination or a transaction connecting the marker and effect. Explain the failure window rather than promising impossible guarantees.

**Failure case.** Two workers can both observe an absent event before either records it.

Each individual operation may be safe while the combined check-and-record sequence races. It must be one atomic operation under a shared lock or equivalent coordination. The later side effect still has separate crash and retry semantics.

## Hit counters aggregate what the query needs

For a counter of hits in the last sixty seconds with integer-second timestamps, keep timestamp-and-count buckets. Multiple hits in the same second increment one bucket. Remove old buckets and maintain a running total. A queue works for monotone timestamps. A ring of sixty buckets can bound space when the query is always that fixed horizon, but each slot must also store its absolute timestamp to distinguish this minute's second from an older cycle.

At now one hundred fifty, the interval strictly after ninety through one hundred fifty contains events at one hundred, one hundred one, one hundred one, and one hundred fifty: four hits. At now one hundred sixty, the event at one hundred is outside a sixty-second window. State the query's now value explicitly; “last sixty seconds” without a clock cannot determine an answer.

A ring with integer-second buckets is exact for an integer-second event model. For arbitrary subsecond queries it may be approximate at the partial boundary bucket. If queries request many different horizons, a single running total is insufficient; scan retained buckets or use prefix sums over immutable history. Storage must retain at least the largest supported horizon. Out-of-order ingestion also needs an explicit lateness policy.

## Time-based values are predecessor queries

Store versions per key ordered by timestamp. For quality at times ten and twenty, a query at fifteen returns the time-ten value. Search for the first timestamp greater than the requested time, then step back. Equal timestamp writes replace in our contract. Earlier writes are rejected so appends preserve order. If historical writes must be accepted, sorted insertion costs linear time in a slice, or you need a structure suited to ordered updates.

Deletion can be a tombstone version: an explicit marker saying the key ceased to exist at that point. Omitting deleted keys from history would make older queries incorrect. Expiring versions is another policy: it may make historical answers unavailable. Always ask whether retention and temporal queries must coexist.

## Schedulers are priority queues plus lifecycle state

Schedule takes an ID, a task, and a deadline. RunReady removes tasks whose deadlines are at or before now. A min heap exposes the earliest deadline. Define ordering for equal deadlines, such as insertion sequence. Cancellation marks an ID inactive. Rescheduling creates a new generation and heap record. When an old record surfaces, ignore it unless the active generation still matches.

Do not execute arbitrary callbacks while holding the scheduler's state lock. Collect ready work or claim it under the lock, then release the lock and execute. Decide whether cancellation after a task is claimed can still stop it. If a callback schedules another task due immediately, decide whether it belongs in the current drain or the next call; either choice should avoid an accidental infinite run.

With multiple workers, claim each task atomically. With crashes, a claimed task may never complete. A lease can make abandoned work eligible for retry, which creates possible duplicate execution; idempotent task handlers then matter. That is a discussion extension, not a demand to build a distributed scheduler during a short coding exercise.

## Identity and time solve different parts of event processing

Imagine receiving a viewing event twice because a sender retried. The repeated event ID tells you it is the same logical event, but its timestamp tells you where it belongs in time. These are independent attributes. Two distinct events can happen at the same timestamp. The same event can arrive at different processing times. Using time alone as a deduplication key merges unrelated events; using event ID alone does not answer a time-window query.

There are also at least two clocks in the story. Event time says when the viewing action occurred. Processing time says when your service handled a message. A late message may be old by event time but new by processing time. A ten-minute duplicate-suppression horizon measured from first acceptance answers a different question from “include viewing events that occurred during the last ten minutes.” State which one each deadline represents.

A scheduler reuses the same separation. Task ID identifies what can be cancelled or replaced. Scheduled time establishes priority. Generation identifies which scheduling decision a heap record represents. A map, deadline heap, and generation check therefore appear naturally, just as in TTL cleanup. This is transfer of a reasoning pattern: delayed work must prove it still applies to the current version before mutating state.

Retries reveal the boundary between an in-memory coding problem and a durable service. Suppose you record “processed” before charging a customer, then crash before the charge. A retry is suppressed and the effect is missing. Record it after charging and a crash in between can allow the charge to repeat. The order of two independent actions cannot eliminate both failure windows. You need the destination to recognize repeated identity or a transaction that connects effect and marker. Even if the interview asks only for an in-memory set, accurately naming its guarantee shows good judgment.

## Worked example: cancellation invalidates a heap record

A scheduler stores jobs by ID and orders work by due time. Schedule A at time five with generation one, and B at time seven with generation two. Reschedule A to time nine with generation three. The map now says A is due at nine, but the heap can retain the obsolete time-five record.

![The heap supplies the next candidate deadline; the current-job map decides whether that candidate is still authorized to run.](figures/scheduler.svg)

| Clock or action | Record examined | Outcome |
| --- | --- | --- |
| Drain at 5 | A, generation 1 | Ignore: current generation is 3 |
| Cancel B | Current B removed | Its heap record becomes stale |
| Drain at 7 | B, generation 2 | Ignore: no current job |
| Drain at 9 | A, generation 3 | Claim and emit A |

To claim a valid due job, remove it from the current-job map before returning it to the caller. A second drain then cannot emit the same generation again. For equal deadlines, use an explicit tie rule such as insertion sequence if reproducible order is required. A job ID alone does not distinguish an old schedule from its replacement.

The heap may hold more records than there are active jobs. If thousands of reschedules target a distant deadline, none of those stale records is removed by draining near-term work. Rebuilding from the current map once stale overhead becomes too large is one possible bound. An indexed heap trades more update bookkeeping for immediate removal instead.

## Worked example: aggregation forgets what expiration later needs

Suppose title A receives a four-minute event at time two and a six-minute event at time eight. Its total is ten. For a trailing ten-second query at time twelve, the first event is exactly on the excluded left boundary and must leave. The answer is six.

A totals map containing only A:10 cannot tell which contribution expired. Keep the individual active contributions in timestamp order as well as the aggregate. Removing the time-two record subtracts four from A's total. If a title's total becomes zero and durations are positive, remove its aggregate entry.

| Time of query | Active interval | Active contributions | A total |
| --- | --- | --- | --- |
| 9 | (-1, 9] | A:4 at 2; A:6 at 8 | 10 |
| 12 | (2, 12] | A:6 at 8 | 6 |
| 18 | (8, 18] | None | 0 |

Deduplication is a separate history. Forgetting an expired contribution does not necessarily permit forgetting its event ID: a late retry can otherwise be mistaken for new work. Conversely, retaining every ID forever bounds duplicates but not memory. The contract must establish a retry horizon, stable event timestamps, or a durable identity mechanism. The queue answers which contributions are active; the identity store answers which deliveries have already been accounted for.

## Putting the program together: a cancellable scheduler

Use Schedule(id, deadline), Cancel(id), and Drain(now), where Drain returns ready IDs instead of running callbacks. Scheduling an already-active ID replaces its previous schedule. Cancellation of an unknown ID returns false. Deadlines equal to now are ready. Ties execute in scheduling-sequence order. This deterministic first version makes behavior easy to trace without introducing workers.

Store a map from task ID to its active generation, and a min heap of records containing deadline, scheduling sequence, ID, and generation. A fresh sequence can also serve as a nonreused generation. Heap order compares deadline first and scheduling sequence second. Schedule updates the active map and pushes a fresh record. Cancel deletes the active map entry and reports whether it existed.

Drain repeatedly pops records whose deadlines are due. A missing active ID means cancellation or prior completion, so skip the record. A mismatched generation means replacement, so skip it. A matching generation is current: delete its active map entry and append its ID to the returned list. Stop when the root is in the future or the heap is empty.

Schedule A for ten, B for five, replace A with three, and cancel B. Drain at four returns A. Drain at twenty returns nothing because the remaining records are cancelled or stale. Schedule two new tasks for the same deadline and their sequence numbers determine stable order. Old records for a reused ID cannot claim a newly scheduled task's lifetime.

Like lazy TTL cleanup, this design can retain stale heap records until their deadlines arrive. State the memory cost. If Drain later runs callbacks, collect or claim work under the lock, release the lock, then execute callbacks. Define whether cancellation after claiming can stop execution. Returning IDs first keeps lifecycle logic separate from arbitrary callback behavior.

## Go example: claim a scheduled generation before returning it

This helper handles a record already popped from a due-time heap. It is not the complete scheduler. The map stores each job's current generation, assigned freshly on scheduling or rescheduling. The caller serializes the check-and-delete sequence and invokes external work only after this helper succeeds.

```go
type JobRecord struct {
    ID string
    Due int64
    Generation uint64
}

func ClaimJob(current map[string]uint64,
    record JobRecord, now int64) bool {
    if record.Due > now {
        return false
    }
    generation, found := current[record.ID]
    if !found || generation != record.Generation {
        return false
    }
    delete(current, record.ID)
    return true
}
```

After rescheduling A from generation one to three, a due record for one returns false and leaves A scheduled. A record for three at its deadline returns true and removes the claim. Calling again with that same record returns false because the current entry is gone. Cancellation uses the same absence branch.

The caller must peek at the heap deadline before popping: a not-yet-due record must remain queued, even though this defensive helper would reject it. This gives at-most-one claim of a generation within this in-memory state, not exactly-once execution of an external effect. If a worker crashes after the successful claim, durable retries and idempotent effects need a larger lifecycle model.
