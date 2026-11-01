# High latency with apparently normal CPU

[Debugging index](README.md) · [Preparation index](../README.md)

Practice scenario for the diagram-based distributed debugging round. These
questions are reconstructed preparation material, not a disclosed interview.

Suppose the interviewer says:

> "Request volume is unchanged. Tail latency has increased sharply. CPU is
> only 35%. What would you investigate?"

Do not infer that capacity is healthy from CPU alone.

| Candidate mechanism | Evidence that would help distinguish it |
|---|---|
| Waiting for database connections | Pool occupancy, pending waiters, connection-acquisition time versus query time |
| Lock contention | Lock waits, blocked transactions, blocked threads or goroutines |
| Slow downstream service | Outbound latency split into waiting, connection establishment, and response time |
| Retry amplification | Original request rate versus total attempt rate, broken down by caller and dependency |
| One hot tenant or partition | Latency and queue depth by key, partition, tenant, and instance |
| Runtime or host contention | Scheduling delay, throttling, pauses, allocation behavior |
| Network setup or connectivity problems | DNS, connection, TLS, timeout, and endpoint-level measurements |

Netflix's noisy-neighbor article illustrates why scheduling delay can be
informative even when application-level CPU measurements do not explain the
symptom.
[Noisy Neighbor Detection with eBPF](https://netflixtechblog.com/noisy-neighbor-detection-with-ebpf-64b1f4b3bbdd)

A useful quantitative explanation is:

> At 100 requests per second and 50 ms **average** time in the system, there
> are about five requests in flight on average. If average time rises to
> 500 ms, that becomes about fifty.

This is Little's Law applied to compatible averages in a stable observation
window. Do not substitute p99 latency into that calculation.

If retries happen at three layers, with each layer allowing three attempts,
one original request can produce up to 27 leaf attempts in the worst case.
That gives a concrete explanation for how unchanged user traffic can create
much more internal work.

Possible mitigations depend on the evidence: reduce admissions, bound
concurrency, isolate a tenant, reduce retry amplification, or remove the
triggering change. Increasing timeouts can increase the amount of work held in
the system, so explain the queueing effect before proposing it.
[Timeouts, retries, and backoff with jitter](https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/)

If the discussion reaches Go, know which diagnostic would answer which
question: CPU profiles for active computation, goroutine stacks for blocked
work, heap profiles for retained allocations, and blocking or mutex profiles
for contention. The choice of Go in the coding interview does not establish
that the supplied distributed system will use Go.
[Go diagnostics](https://go.dev/doc/diagnostics)
