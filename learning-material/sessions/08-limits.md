# Chapter 8 — Rate limiting and rolling windows

The policy determines the data structure. “One hundred per minute” is not precise enough to choose an algorithm.

## Define what counts

Our exact limiter allows at most L accepted requests per user in the rolling interval strictly after now minus W and at or before now. Rejected attempts do not count. Requests can share a timestamp. Time is monotone for the limiter and uses integer ticks. Users are independent. At the exact left boundary, an old request no longer counts.

For limit two and window ten, requests at time zero and time one are accepted; time nine is rejected. At time ten, the request at zero leaves the window, so the new request is accepted. The retained accepted timestamps are one and ten. Writing this trace before code settles the off-by-one choice.

## Exact sliding log

Store a map from user ID to a queue of accepted timestamps. Before deciding, discard timestamps at or before now minus W. If the remaining length is already L, reject. Otherwise append now and accept. Chronological order makes expired records a prefix, so no search is needed. A slice with a head index avoids shifting every remaining timestamp on every removal.

Each accepted timestamp is appended once and removed once, so work is amortized constant per request. A single request after a long idle period may remove L entries and cost O(L). Memory is O(U times L) for U active users, plus bookkeeping for idle users if you never delete them. Rejected attempts must not be appended under this contract, or an attacker could keep extending lockout by sending more requests.

With twenty million users and limit one hundred, even bare eight-byte timestamps could occupy sixteen billion bytes if every queue were full, before map and slice overhead. This is an illustrative worst-case calculation, not a measured Netflix workload. Aggregate timestamps within buckets when approximate precision is acceptable; retain exact timestamps when boundary accuracy is required. Delete empty idle users or expire their state, while ensuring cleanup itself does not create unbounded latency.

**Think aloud.** Why does the queue algorithm fail if client timestamps arrive out of order? What is the simplest contract that restores correctness?

#### Worked answer — Coach notes

An expired timestamp might appear behind a fresh one, so removing only a prefix fails. Use monotone server admission time for a request limiter. Event-time analytics with late arrivals is a different problem and may require ordered storage, lateness bounds, or explicit rejection of old data.

## Fixed windows are cheaper but allow boundary bursts

A fixed-window counter stores a bucket number and count per user, resetting when the bucket changes. With limit one hundred, a user can send one hundred requests just before a minute boundary and another hundred just after it. That is valid per fixed minute but violates a strict rolling-minute maximum. Say which guarantee the implementation provides.

A two-bucket sliding estimate weights the previous bucket's count by the fraction still overlapping the current window. It is compact but assumes something about the distribution of requests within a bucket. It cannot be exact without more information. At interview scale, describing the approximation and a burst counterexample is more valuable than asserting it is “basically the same.”

## Token buckets express burst capacity and sustained rate

A token bucket has capacity C, refill rate r tokens per tick, current tokens, and last refill time. On each request, add elapsed time times r, capped at C, then move last refill time to now. If at least one token exists, subtract one and accept. A bucket starts full in our contract, allowing an initial burst.

With capacity three and refill one token per second, three requests at time zero pass and a fourth fails. At time one, one token has refilled, so one request passes. Long inactivity never accumulates more than three tokens. Fractional tokens allow smooth refill; integer arithmetic with retained remainder is an alternative when exact units are important.

A token bucket with capacity one hundred and refill ten per second is not “at most one hundred per rolling minute.” It intentionally permits a burst plus continued refill. Over an interval of length T, accepted demand is bounded roughly by capacity plus rate times T, subject to initial state and endpoint conventions. Choose it when bursts are acceptable and long-run throughput must be controlled.

Refill must account for rejected requests too: update the timestamp whenever you materialize elapsed refill, or the same elapsed time can be counted again. Our reference clamps a backward clock reading to the previous time. That is defensive behavior, not permission to use arbitrary client timestamps. At very small scales, floating-point arithmetic can affect boundaries; tests use exactly representable examples and document this limitation.

## Concurrency and multiple servers

Two simultaneous requests can both observe one remaining slot and both accept unless prune, check, and append form one atomic operation. A mutex solves this inside a process. Per-user locks or shards can improve concurrency when required, but map creation and lock lifetime must still be coordinated. Token refill, check, and debit likewise need one critical section.

Time must also follow that serialized operation order. Two callers can sample increasing timestamps before locking but acquire the lock in reverse order. Sampling an injected server clock inside the same critical section preserves the ordering needed for FIFO expiry.

Independent per-server limiters can multiply a user's global allowance. A shared atomic decision, stable ownership of each user's state, or intentionally partitioned quotas can address that, with different availability and utilization trade-offs. A distributed system discussion should name failure policy: if the limiter's state service is unavailable, does the request pass or fail? No single answer is correct for every service.

## Start from the guarantee a caller should observe

“One hundred requests per minute” can describe several different products. A billing counter might group requests by calendar minute. An abuse rule might forbid any rolling sixty-second interval from containing more than one hundred accepted requests. A service protection rule might permit a burst after inactivity while enforcing a sustained rate. Those are not small implementation variations; they permit different request sequences. A correct implementation begins with the observable guarantee.

An exact sliding log retains the recent accepted requests because each timestamp can affect a future boundary decision. A fixed-window counter forgets where within the bucket those requests occurred. That lost detail is why it cannot distinguish one hundred requests spread throughout a minute from one hundred clustered at the end. Compression of state usually comes with a question: are the histories being merged truly equivalent for the required queries?

A token bucket summarizes a different history. Its token balance says how much accumulated permission remains. Waiting restores permission at a configured rate, but the cap prevents indefinite accumulation. You can think of the balance as unused capacity, rather than as a count of events inside a precise interval. That interpretation explains why a bucket naturally allows bursts and why it does not enforce an exact rolling-window count with the same parameters.

Atomicity matters at the decision point. Imagine the exact limiter has one slot remaining. Two callers independently read that fact, then both append an accepted request. Every individual access could use a thread-safe map and still produce an invalid combined result. The protected operation must include removal of obsolete events, the limit check, and recording acceptance. For a bucket it must include refill, checking balance, and debiting a token. The state invariant is what defines the lock boundary.

Large user populations add lifecycle questions. A user who disappears does not trigger their own cleanup. If state remains forever, the number of remembered users can grow even while each queue is bounded. Deleting idle state is safe only when reconstructing it gives the same future behavior. An empty exact sliding log can disappear after all counted requests leave the horizon. A token bucket can disappear once it would be fully replenished, assuming a recreated bucket starts full.

## Questions and worked explanations on admission policies

**Round 1.** State an exact two-request, ten-second rolling rule in words. Are rejected attempts counted, and what happens to a request exactly ten seconds old?

#### Worked answer — Explanation and next challenge

Our rule counts accepted requests strictly newer than now minus ten and no later than now. Rejected attempts are not added. An accepted request exactly ten seconds old leaves the window. Other policies can be valid, but their behavior must be stated before the queue logic is chosen.

**Round 2.** Requests arrive at zero, one, and nine for the same user. With limit two and window ten, decide each one. Now decide a request at ten.

#### Worked answer — Explanation and next challenge

Zero and one pass. Nine fails because both earlier accepted requests remain in the window. At ten, zero expires, leaving one counted request, so ten passes. The active accepted timestamps are now one and ten. Repeating the state aloud makes the boundary rule concrete.

**Round 3.** What changes if you append rejected requests too? Describe how repeated rejected attempts could affect the user's ability to recover.

#### Worked answer — Explanation and next challenge

You would be implementing an attempt-based policy rather than the stated accepted-request policy. Rejected attempts could keep the remembered count high and extend the effective lockout as traffic continues. That may be intended in another security rule, but it is not a harmless coding choice.

**Round 4.** Why does chronological order let you remove expired requests only from the queue front? What breaks if an older timestamp is inserted behind a fresh one?

#### Worked answer — Explanation and next challenge

With ordered timestamps, all expired entries form a prefix. If the front is fresh, every later timestamp is at least as fresh. An out-of-order old entry can hide behind it, defeating that inference. Server admission time can provide the monotone ordering for this limiter; event-time analysis needs a different lateness contract.

**Round 5.** A user sends two requests just before a fixed-minute boundary and two just after it. With a fixed limit of two per minute, what happens? Why does this differ from the rolling rule?

#### Worked answer — Explanation and next challenge

All four can pass because each fixed bucket contains only two. A short rolling interval spanning the boundary contains four and would reject some under the exact rolling guarantee. The bucket counter forgot within-minute placement. Use this example to explain the policy trade-off rather than simply calling one implementation wrong.

**Round 6.** A token bucket has capacity three and refill rate one per second. It starts full. Describe four immediate requests, then one request a second later.

#### Worked answer — Explanation and next challenge

The first three consume the three tokens and pass. The fourth fails. One second later one token has refilled, so another request passes. The bucket allows burst capacity followed by a sustained refill rate. It does not need to retain the individual timestamps of the first three requests.

**Round 7.** The same bucket sits idle for an hour. How many immediate requests can it then allow, and why doesn't the long idle period permit thousands?

#### Worked answer — Explanation and next challenge

At most three, because token balance is capped at capacity. Without that cap the limiter would bank an unbounded amount of permission and permit a huge later burst. Capacity and refill rate control different parts of behavior: burst size and recovery or sustained throughput.

**Round 8.** A request triggers refill but is rejected because the balance is still below one token. Why should the last-refill timestamp still be updated?

#### Worked answer — Explanation and next challenge

The elapsed interval has already contributed tokens. Leaving the old timestamp causes the same interval to be credited again on the next request. Even a rejected operation may legitimately update internal accounting. The balance and its reference time must describe the same instant.

**Round 9.** Two concurrent callers see one token. Describe the bad interleaving and the smallest state transition that must be atomic.

#### Worked answer — Explanation and next challenge

Both check the balance before either decrements, then both accept. Protect refill, check, and debit together. A mutex around only the map lookup does not protect the policy decision. Inside one process a simple lock is often a sufficient baseline; optimization can come after correctness.

**Round 10.** Two servers each maintain a separate limit of one hundred for the same user. Can you promise a global limit of one hundred? Give the counterexample.

#### Worked answer — Explanation and next challenge

The user can consume one hundred on each server, reaching two hundred globally. Global enforcement needs shared atomic state, stable ownership, partitioned quotas, or another explicit coordination policy. Naming a data store does not itself prove the read-modify-write decision is atomic.

**Round 11.** A limiter stores at most one hundred timestamps per user but never deletes user records. Why can memory still grow indefinitely, and when can exact-log state be safely discarded?

#### Worked answer — Explanation and next challenge

New user IDs create new queues, so a per-user bound does not bound the number of users. Once a user's counted timestamps are all outside the horizon, an empty log is equivalent to no record for future admission. Cleanup must find idle users without relying only on their next request.

**Round 12.** Choose a policy for a strict rolling abuse cap and another for a burst-tolerant backend throughput limit. Explain the guarantee, retained state, and a boundary test for each.

#### Worked answer — Explanation and mastery check

An exact sliding log fits the strict rolling cap, with per-user accepted timestamps and a test exactly at expiry. A token bucket fits burst tolerance with tokens and last-refill time, tested with a burst, short recovery, and long idle saturation. Correctness means matching the promised behavior, not choosing the algorithm with the most sophisticated name.

Have the tutor describe a request sequence and ask which policies accept it. Then reverse roles: invent a short sequence that distinguishes two policies. This tests whether you understand behavior rather than parameter names.

## Worked program design: exact admission and a token bucket

For the exact limiter, the public operation Allow(user, now) returns a boolean. State consists of the limit, window length, and per-user timestamp queues. Each queue has a timestamp slice and a head index. The contract requires nondecreasing admission time and counts only accepted requests in the interval after now minus window through now.

Allow finds or creates the user's queue. Advance its head while the oldest retained timestamp is at or before the cutoff. Count the entries from head to the end. If already at the limit, return false without appending. Otherwise append now and return true. Periodically compact the slice when discarded slots become a significant fraction of its contents. The pruning, count check, and append belong to one critical section if callers overlap.

With limit two and window ten, Allow at zero and one returns true; at nine false; at ten true. After that last call the active queue contains one and ten. A rejected call must not alter that accepted-event sequence. Another user has independent state. Idle user cleanup is a separate bounded-memory concern; no new request arrives to trigger self-cleanup for a disappeared user.

For a token bucket, the fields are capacity, refill rate, current tokens, and last-accounted time. Allow first brings tokens up to now: add elapsed time times rate, capped at capacity, and update the accounted time. If tokens are below one, return false. Otherwise subtract one and return true. Protect that full sequence atomically. Define the clock policy and whether fractional tokens are represented by floating point or exact integer units.

A capacity-three bucket starts with three immediate acceptances and then rejects. After one second at one token per second it accepts one more. Long idle time returns it only to capacity. These tests validate a different guarantee from the exact rolling log. The two structs share an admission API but should not be presented as interchangeable implementations of the same rule.

**Optional spoken walkthrough:** dictate the branch order for each Allow method. Explain what state changes even on a rejected token-bucket request, then show a boundary burst that distinguishes a fixed-window counter from the exact queue.
