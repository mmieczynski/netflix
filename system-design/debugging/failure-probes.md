# Additional diagram based debugging probes

[Debugging index](README.md) · [Preparation index](../README.md)

Practice scenario for the diagram-based distributed debugging round. These
questions are reconstructed preparation material, not a disclosed interview.

| Given symptom or change | Deep dive to prepare |
|---|---|
| **A cache returns old data after invalidation** | A slow read fetched an old value before the write, then repopulated the cache afterward. Discuss version-aware publication and the consistency limits of the chosen cache contract. |
| **Only one service instance is failing** | Compare version, configuration, placement, traffic mix, connections, and local resource behavior before assuming an application-wide defect. |
| **The database has low CPU but writes are slow** | Locks, storage latency, connection limits, replication waits, and a hot key can dominate without high aggregate CPU. |
| **An event is processed twice** | Identify the exact duplicated effect. Broker delivery guarantees do not establish sink-side business correctness. |
| **A trace has a missing span** | Sampling, instrumentation, propagation, and exporter failures are alternatives to "the call never happened." |
| **A regional failover makes the incident worse** | Remaining capacity, cold caches, downstream quotas, and replicated-data readiness. Two equally sized regions each at 60% load cannot necessarily fit into one surviving region. |
| **A service catalog lists the wrong owner** | Determine the authority for ownership, synchronization delay, and how incorrect routing affects mitigation and change approval. |
| **The control database is unavailable** | Separate new change admission and coordination from already deployed services. Explain which operations must stop and which can continue from existing state. |
| **Clock skew changes observed ordering** | Avoid treating wall-clock timestamps from different machines as a universal order. Use appropriate sequence/version information and explicit timing assumptions. |
| **A metric improves while users complain** | Check aggregation, cohort selection, missing data, and whether the measured success condition matches user outcomes. |
| **An automation loop keeps changing the same service** | Inspect oscillating policy decisions, stale observations, competing controllers, and whether desired state converges. |
| **A small canary passes but the larger rollout fails** | Examine workload representativeness, sample size, rare paths, dependency capacity, and nonlinear effects at higher concurrency. |

For this round, the value comes from following a failure through the entire
chain. For example:

> "The dependency slows down, callers hold resources longer, retries increase
> attempts, the connection pool fills, unrelated requests begin waiting, and
> the service's own health checks may then fail."

That is substantially stronger than listing "timeouts, retries, circuit
breakers" without explaining how the failure propagates. Netflix's
load-shedding work and Google's cascading-failure chapter are useful sources
for this reasoning.
[Netflix prioritized load shedding](https://netflixtechblog.com/keeping-netflix-reliable-using-prioritized-load-shedding-6cc827b02f94),
[Addressing Cascading Failures](https://sre.google/sre-book/addressing-cascading-failures/)
