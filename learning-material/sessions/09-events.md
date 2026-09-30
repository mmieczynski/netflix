# Chapter 9 — Deduplication counters versions and scheduling

Many practical tasks reuse the same ideas: key identity, chronological removal, predecessor lookup, and versioned invalidation.

## Deduplication needs an identity and a horizon

An event has an event ID, user ID, title ID, duration, and event time. To suppress repeat deliveries, first define identity. The same event ID identifies the same logical event in our exercise; two different IDs with identical payloads are separate events. Store accepted IDs in a set. A repeated ID is rejected, while an unseen ID is recorded and accepted.

With ten-minute retention, use a deadline per ID. Decide whether a duplicate refreshes the deadline. Here it does not: expiration is measured from the first accepted processing time. A repeat just before expiry therefore does not suppress that ID forever. Use server processing time for this exercise, not an untrusted event timestamp. A TTL map with a cleanup heap implements the contract directly.

Now consider a callback that performs the actual side effect. If you mark the ID before the callback and crash, a retry can be suppressed even though the effect never completed. If you mark it afterward and crash in between, the effect can happen twice. A local set is duplicate suppression during its lifetime, not a proof of durable exactly-once effects. Stronger behavior needs an idempotent destination or a transaction connecting the marker and effect. Explain the failure window rather than promising impossible guarantees.

**Think aloud.** Two workers receive the same ID at once. Both call Get, see missing, then call Put. Why is a thread-safe cache alone insufficient?

#### Worked answer — Coach notes

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

## Questions and worked explanations on temporal services

**Round 1.** Two events have different IDs but identical payloads. Another pair has the same ID delivered twice. Under event-ID deduplication, which pair is merged and why?

#### Worked answer — Explanation and next challenge

Only deliveries sharing the same event ID are merged. Identical payloads can represent two legitimate actions. Deduplication depends on the contract's logical identity, not on superficial similarity. If the sender can reuse an ID for conflicting payloads, define rejection, validation, or overwrite behavior explicitly.

**Round 2.** An ID is remembered for ten minutes from its first accepted processing time. A duplicate arrives after nine minutes. Does its memory deadline move? What behavior would changing that rule create?

#### Worked answer — Explanation and next challenge

It does not move under this contract. Refresh-on-duplicate would create a sliding suppression horizon and could retain a noisy ID indefinitely. Both policies may be intentional, but their storage lifetime and acceptance behavior differ. Tests must include a duplicate near the boundary.

**Round 3.** Two workers check the same absent ID, then both record it. What exactly must become atomic to suppress one duplicate inside a process?

#### Worked answer — Explanation and next challenge

The combined check-and-record must be one protected transition. Locking Get and Put separately leaves a gap between them. An atomic claim can choose one winner, but the subsequent side effect still has its own success and crash semantics. Keep those guarantees separate.

**Round 4.** Explain the crash window when you mark an event processed before its side effect, and the different crash window when you mark it afterward.

#### Worked answer — Explanation and next challenge

Mark first and a crash can suppress an effect that never happened. Effect first and a crash can allow a retry to repeat an effect that did happen. Idempotent effects or coordinated transactions address the gap. A local set can be a correct answer to a bounded in-memory exercise without proving durable exactly-once execution.

**Round 5.** Hits occur at one hundred, one hundred one, and one hundred fifty. At now one hundred sixty with a sixty-second window excluding the left edge, which hits count?

#### Worked answer — Explanation and next challenge

One hundred one and one hundred fifty count; one hundred is exactly on the excluded left boundary. The query needs an explicit now. Repeating this boundary convention across counters and limiters helps you notice when a new task intentionally chooses a different one.

**Round 6.** A circular buffer reuses the same slot every sixty seconds. Why must each slot store its absolute timestamp as well as a count?

#### Worked answer — Explanation and next challenge

The slot index alone cannot distinguish a count from the current cycle from an older cycle. On reuse, compare the stored timestamp and reset if it represents another bucket. A ring bounds positions, not the semantic age of their contents. Time metadata tells you whether the stored count is relevant.

**Round 7.** Values for a key are stored at times ten and twenty. How do you answer a query at fifteen? What if historical writes at time twelve are newly allowed?

#### Worked answer — Explanation and next challenge

Return the latest version at or before fifteen, which is time ten. Appending a time-twelve write after time twenty would break sorted history. Either reject it under the old contract, insert it into order with the associated cost, or change the structure. A new write policy can invalidate an otherwise correct search.

**Round 8.** A key is deleted at time twenty but existed at ten. Why is simply erasing all versions inadequate if historical queries must still work?

#### Worked answer — Explanation and next challenge

A query at fifteen should still see the old value. A tombstone at twenty marks absence from that point while preserving history. Removing everything loses valid earlier answers. A tombstone is explicit information about a state transition, not merely an empty slot.

**Round 9.** Schedule A at ten and B at five, then reschedule A to three and cancel B. What should draining at four return? What should draining at twenty return?

#### Worked answer — Explanation and next challenge

At four return A once. At twenty return nothing: B is cancelled and the old A-at-ten record is stale. The active-ID map and generation checks determine validity; the heap determines which deadline records are ready to inspect. Neither structure alone answers both questions.

**Round 10.** A ready task callback schedules another task while the scheduler still holds its state mutex. What can go wrong, and how would you separate state mutation from callback execution?

#### Worked answer — Explanation and next challenge

The callback may reenter an operation requiring the same mutex and deadlock, or simply hold up all scheduling while it runs. Claim or collect ready work under the lock, release it, and then execute callbacks. Define whether newly scheduled immediate work belongs to this drain or the next call.

**Round 11.** A worker claims a task, then crashes. A lease later allows another worker to retry it. Why can this improve completion while still permitting duplicate execution?

#### Worked answer — Explanation and next challenge

The first worker may have performed the effect before crashing or losing its lease, without recording completion. The retry then repeats the effect. Leases restore liveness for abandoned work but do not by themselves make effects exactly once. An idempotent handler or destination-level identity check may still be required.

**Round 12.** Explain the common pattern linking TTL cleanup, rescheduled tasks, and late background fetches. Name what the identifier, deadline, and version each contribute.

#### Worked answer — Explanation and mastery check

The identifier locates the logical object; a deadline or readiness condition orders delayed work; a version proves that delayed work still belongs to the current object lifetime. Old work can be discarded without changing new state. The memory cost of retained stale work must still be bounded or acknowledged.

Ask for a task with a different story but the same delayed-work problem, such as expiring reservations. Identify the current state, delayed record, validity check, and cleanup policy before choosing an implementation.

## Worked program design: a cancellable scheduler

Use Schedule(id, deadline), Cancel(id), and Drain(now), where Drain returns ready IDs instead of running callbacks. Scheduling an already-active ID replaces its previous schedule. Cancellation of an unknown ID returns false. Deadlines equal to now are ready. Ties execute in scheduling-sequence order. This deterministic first version makes behavior easy to trace without introducing workers.

Store a map from task ID to its active generation, and a min heap of records containing deadline, scheduling sequence, ID, and generation. A fresh sequence can also serve as a nonreused generation. Heap order compares deadline first and scheduling sequence second. Schedule updates the active map and pushes a fresh record. Cancel deletes the active map entry and reports whether it existed.

Drain repeatedly pops records whose deadlines are due. A missing active ID means cancellation or prior completion, so skip the record. A mismatched generation means replacement, so skip it. A matching generation is current: delete its active map entry and append its ID to the returned list. Stop when the root is in the future or the heap is empty.

Schedule A for ten, B for five, replace A with three, and cancel B. Drain at four returns A. Drain at twenty returns nothing because the remaining records are cancelled or stale. Schedule two new tasks for the same deadline and their sequence numbers determine stable order. Old records for a reused ID cannot claim a newly scheduled task's lifetime.

Like lazy TTL cleanup, this design can retain stale heap records until their deadlines arrive. State the memory cost. If Drain later runs callbacks, collect or claim work under the lock, release the lock, then execute callbacks. Define whether cancellation after claiming can stop execution. Returning IDs first keeps lifecycle logic separate from arbitrary callback behavior.

**Optional spoken walkthrough:** describe the record comparator and the three validity branches in Drain. Then explain what new state would be needed for retries after a worker crash, and which execution guarantee a lease alone does not provide.
