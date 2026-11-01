# Hello Interview guided practice priorities

[Preparation index](README.md) · [Interview context](interview-context.md)

Catalog checked on **1 November 2026**. This order is tailored to the Netflix
L5 Developer Platform role covering HIPR and Change Systems and to the stated
preference for **guided practice**. Exercises without a written solution remain
eligible for the highest priority.

**Start with GitHub Actions, Job Scheduler, and Metrics Monitoring.** Continue
with Ad Click Aggregator, Distributed Cache, and Rate Limiter. Together these
cover change execution, durable scheduling, health signals, event correctness,
stale data, and overload protection.

**YouTube Top K comes next**, adding windowed rankings and aggregation across
partitions. Its overlap with Ad Click Aggregator keeps it behind the first six
in the combined preparation order; move it alongside that exercise for a
confirmed HIPR focus.

All exercises below are listed in the
[System Design Guided Practice catalog](https://www.hellointerview.com/practice/system-design).
GitHub Actions appears under **More Practice**: it has guided practice and
feedback, but no written breakdown. That makes it especially relevant for
this preparation rather than a reason to postpone it.

The ranking is an assessment of relevance to the recruiter brief, not evidence
that Netflix asks these exact exercises. Public problem descriptions and
outlines establish the published scope where available. The full prompts and
feedback inside signed-in guided sessions were not inspected. The Netflix
adaptations and follow-up questions below are original preparation extensions.

## Priority order

| Priority | Guided exercise | Why it belongs here |
|---|---|---|
| 1 | **GitHub Actions** | Closest domain match to Change Systems: coordinated changes, execution infrastructure, validation, and interactions with source control and delivery. |
| 2 | **Job Scheduler** | The durable execution foundation: scheduling, ownership, recovery, retries, and resource allocation. |
| 3 | **Metrics Monitoring** | Closest match to HIPR's measurement responsibilities and to debugging the health of the observation system itself. |
| 4 | **Ad Click Aggregator** | Strong event-processing exercise for duplicate events, delayed data, hot partitions, replay, and trustworthy reports. |
| 5 | **Distributed Cache** | High value for debugging stale data, hot keys, node failure, uneven load, and recovery-induced traffic spikes. |
| 6 | **Rate Limiter** | Shared mutable state, overload policy, tenant fairness, and protecting dependencies from fleet automation. |
| 7 | **YouTube Top K** | Windowed rankings, aggregation across partitions, hot items, freshness, and exact versus approximate results. |
| 8 | **Web Crawler** | Distributed workers, discovery, deduplication, progress tracking, and per-source load control. |
| 9 | **Payment System** | Durable intent, uncertain external outcomes, state transitions, and reconciliation across systems. |
| 10 | **LeetCode** | Running code in isolated workers, resource limits, execution results, and bursty workloads. Relevant to CI and code-transformation execution. |
| 11 | **Notification System** | Internal platform contracts, routing findings to owners, priority isolation, and duplicate delivery. |

This is a combined priority order for both teams. If team matching becomes
explicitly HIPR-focused, move Metrics Monitoring, Ad Click Aggregator, and
YouTube Top K ahead of the execution problems. For a Change Systems focus, keep GitHub Actions and
Job Scheduler first and move Payment System and LeetCode ahead of the lower
priority discovery and notification exercises.

## GitHub Actions

This is the strongest domain match by the publicly listed topic and the role
description. The catalog confirms that the exercise exists; it does not expose
its detailed scenario. The following focus is therefore a preparation
adaptation, not a summary of a hidden answer key.

Use the exercise to practise explaining how a source-control event becomes a
durable workflow run, how steps depend on each other, how workers receive work,
and how artifacts and results are associated with the correct revision. Then
extend the design to a migration spanning many repositories.

Follow-ups to be able to defend:

- **A source-control webhook is delivered twice.** Define the event identity,
  workflow-run identity, and the difference between a duplicate event and an
  intentional rerun.
- **A worker completes an external step and dies before recording success.**
  Identify what the controller can reconcile and what must be idempotent at
  the effect boundary.
- **A workflow definition changes while runs are active.** Bind each run to a
  version and explain whether old runs finish, migrate, or stop.
- **A campaign is paused while builds and deployments are running.** Separate
  admission, cancellation requests, and effects already committed elsewhere.
- **Several campaigns touch the same service.** Explain the shared conflict
  domain, rollout budget, and protection against stale plans.

Pair it with the local [fleet change platform](problems/01-fleet-change-platform.md),
[lost acknowledgment](debugging/01-lost-acknowledgment.md), and
[stale ownership](debugging/02-stale-ownership.md) scenarios.

## Job Scheduler

The public outline includes immediate, future, and recurring execution, status
tracking, timely scheduling, scale, and at-least-once execution. Its initial
scope excludes cancellation and rescheduling; those remain valuable extensions
for this Netflix role.
[Published problem scope](https://www.hellointerview.com/learn/system-design/problem-breakdowns/job-scheduler)

Follow-ups:

- **The scheduler restarts after a scheduled occurrence should have run.**
  Define whether to catch up, skip, or coalesce missed occurrences and how to
  prevent creating the same occurrence twice.
- **A heartbeat continues but the job is stuck.** Distinguish worker liveness
  from useful progress and decide what evidence permits replacement.
- **A lease expires and the old worker resumes.** Explain ownership generation
  and protection for both status writes and external effects.
- **One tenant floods the queue.** Define admission, fairness, priority, and
  starvation behavior rather than relying only on a queue's existence.
- **Cancellation races with completion.** Define the accepted state
  transitions and the outcome visible to the user.

Pair it with the [durable job scheduler](problems/04-durable-job-scheduler.md)
and [growing queue](debugging/05-growing-queue.md).

## Metrics Monitoring

The public scope covers metrics ingestion, dashboard queries, alert rules,
and notifications. Its deep-dive outline includes query latency, alert delay,
availability under spikes, and cardinality growth. Logs and distributed tracing
are outside that example's initial scope.
[Published problem scope](https://www.hellointerview.com/learn/system-design/problem-breakdowns/metrics-monitoring)

For HIPR, extend the exercise from runtime measurements to versioned evidence
used in decisions. A responsive dashboard is insufficient if its evidence
describes an obsolete deployment.

Follow-ups:

- **The pipeline stops receiving data but every dashboard stays green.**
  Define missing-data behavior, freshness checks, and how to observe the
  collector and evaluator independently.
- **A new metric label creates millions of series.** Trace resource growth
  through ingestion, storage, queries, and alert evaluation; explain how to
  isolate the offending producer.
- **Old and new application versions run simultaneously.** Explain how to
  attribute evidence and how aggregate health can hide one failing cohort.
- **An alert arrives late.** Separate observation time, ingestion delay,
  evaluation delay, and delivery delay before selecting a mitigation.

Pair it with the [technical health platform](problems/02-technical-health-platform.md)
and [stale health evidence](debugging/03-stale-health-evidence.md).

## Ad Click Aggregator

The public problem concerns collecting events and serving aggregate analytics.
Its discussion includes scaling partitions, recovering data, and reconciling
streamed results with retained raw events. The transferable value is the event
pipeline; detailed advertising-domain preparation is secondary for this role.
[Published problem](https://www.hellointerview.com/learn/system-design/problem-breakdowns/ad-click-aggregator)

Translate the exercise into deployment events, rule evaluations, or developer
productivity measurements.

Follow-ups:

- **A consumer updates an aggregate, then crashes before recording its input
  position.** Explain the duplicate-effect boundary and how recovery avoids
  double-counting.
- **An event arrives after a reporting window closed.** Define whether the
  result is corrected, versioned, or deliberately left unchanged.
- **One service dominates a partition.** Explain how to redistribute work
  while preserving the ordering or aggregation semantics that matter.
- **A historical replay runs alongside live ingestion.** Protect current
  freshness and distinguish replayed evidence from new logical events.
- **The dashboard disagrees with the authoritative raw record.** Identify a
  reproducible reconciliation procedure and the safe correction path.

Pair it with [event ingestion](problems/05-event-ingestion-and-client-libraries.md)
and the [reconciliation pipeline](problems/07-reconciliation-pipeline.md).

## Distributed Cache

The public scope includes basic key-value operations, expiration, and LRU
eviction. Its outline also covers node failures, distribution, and hot reads
and writes. The example permits eventual consistency and excludes durability;
confirm those assumptions separately in a real interview.
[Published problem scope](https://www.hellointerview.com/learn/system-design/problem-breakdowns/distributed-cache)

This repository already contains local cache exercises. Use this guided
practice to concentrate on behavior across machines and on the backing
system's response when the cache degrades.

Follow-ups:

- **A node disappears.** Explain key ownership, replication, and the sudden
  increase in traffic to the source of truth.
- **An old read repopulates a key after invalidation.** Walk through the exact
  ordering and explain the consistency contract of any fix.
- **Hit rate falls from 99% to 80% at constant request volume.** Miss traffic
  becomes twenty times larger; relate that to downstream queues and retries.
- **One key dominates reads or writes.** Explain why adding ordinary shards
  may not remove that bottleneck.

Pair it with [high latency with normal CPU](debugging/04-high-latency-normal-cpu.md)
and the cache-related [failure probes](debugging/failure-probes.md).

## Rate Limiter

The public problem covers client identification, configurable request limits,
rejection behavior, scale, and availability trade-offs. Its example accepts
some enforcement inconsistency. A safety-critical change budget may require a
different contract.
[Published problem](https://www.hellointerview.com/learn/system-design/problem-breakdowns/distributed-rate-limiter)

The repository's local limiter exercises are useful supporting material. The
guided exercise should emphasize coordination across instances.

Follow-ups:

- **Two gateways admit the same final available token.** Identify the atomic
  read, decision, and update boundary.
- **The shared state service is unreachable.** Choose and justify rejecting,
  allowing, or using a bounded local allocation for the specific operation.
- **Every campaign has its own limit, but their combined traffic overloads a
  dependency.** Define limits shared across campaigns and failure domains.
- **A few requests occupy workers for hours.** Distinguish request-rate limits
  from concurrency limits and resource budgets.
- **Clients retry immediately after rejection.** Explain retry amplification,
  feedback to callers, and the server's independent protection.

Pair it with [distributed counting](problems/06-distributed-counting.md)
and [growing queues](debugging/05-growing-queue.md).

## YouTube Top K

The public exercise ranks videos by views over defined periods. Its initial
contract uses exact results and completed tumbling windows; it later explores
approximation. Guided practice is listed in the catalog.
[Published problem](https://www.hellointerview.com/learn/system-design/problem-breakdowns/top-k)

For HIPR, adapt the exercise to ranking services by a clearly defined signal,
such as new policy violations or failed validations. Clarify whether a count,
rate, or severity-weighted measure actually supports the decision. High traffic
alone should not accidentally determine which service looks least healthy.

Ad Click Aggregator remains first for general event-pipeline coverage. Top K
adds useful depth when the question concerns identifying the largest or most
important items across a fleet. For Change Systems, execution correctness and
recovery still take priority over ranking.

Original follow-ups to practise:

- **Does "last hour" mean a completed clock hour or a continuously moving
  interval?** Explain the additional state and expiration work needed when
  counts can decrease as events leave the window.
- **Can the global result be computed by merging each shard's top K?** This
  works when each item's complete score belongs to one shard, all shards use
  the same window, and ties follow a common deterministic order. It can fail
  when an item's contributions are spread across shards.
- **A hot item is split across several workers.** Combine its partial counts
  before relying on local rankings, or explain another method that preserves
  the required accuracy.
- **Can the system retain only each time bucket's winners?** Not generally
  when later queries combine buckets: an item can win overall without winning
  any individual bucket.
- **A retry repeats an event.** Deduplicate the logical event rather than the
  video or service ID; distinct events for the same item must still count.
- **One partition is delayed.** Explain how result freshness and completeness
  are reported, and whether late events revise an already published result.

For a concrete counterexample, suppose K is 1 and the same item can appear in
both partitions:

| Item | Partition 1 | Partition 2 | Total |
|---|---|---|---|
| A | 60 | 60 | 120 |
| B | 100 | 0 | 100 |
| C | 0 | 100 | 100 |

The local winners are B and C, but the global winner is A. Merging only the
local winners loses A. The same numbers demonstrate why retaining only hourly
winners cannot generally produce an exact ranking across multiple hours.

Pair it with [event ingestion](problems/05-event-ingestion-and-client-libraries.md),
[distributed counting](problems/06-distributed-counting.md), and
[stale health evidence](debugging/03-stale-health-evidence.md).

## Web Crawler

The public problem combines distributed fetching, progress recovery,
per-domain load control, and deduplication. Those ideas transfer to repository
scanning and inventory collection; the actual external-system contracts must
be redefined for that adaptation.
[Published problem](https://www.hellointerview.com/learn/system-design/problem-breakdowns/web-crawler)

Follow-ups:

- **A worker fetched an item but failed before publishing discovered work.**
  Explain how to recover both processing progress and newly discovered items.
- **One source is slow or unavailable.** Keep other sources moving while
  preserving a visible record of incomplete coverage.
- **A probabilistic duplicate filter skips an unseen item.** Decide whether
  that error is acceptable for a health inventory and where exact checks are
  required.
- **The same service appears under several identities.** Explain canonical
  identity, ownership changes, and version-aware reconciliation.

Pair it with the [service catalog and dependency graph](problems/03-service-catalog-and-dependency-graph.md).

## Payment System

The public scope includes initiating an external operation and tracking its
status. The deep-dive outline emphasizes durable records, auditability, and
safety across asynchronous external systems.
[Published problem scope](https://www.hellointerview.com/learn/system-design/problem-breakdowns/payment-system)

Focus on the correctness model that transfers to change automation:

- **The external system succeeded but the response was lost.** Explain why
  timeout is an unknown outcome and how to find authoritative evidence.
- **Callbacks are duplicated or arrive out of order.** Protect state
  transitions and separate one operation from its delivery attempts.
- **Local and external records disagree.** Define reconciliation and a repair
  policy that does not repeat irreversible effects.
- **A compensating operation also fails.** Keep durable recovery state and
  distinguish compensation from atomic rollback.

Pair it with [lost acknowledgments](debugging/01-lost-acknowledgment.md)
and [incomplete rollback](debugging/06-incomplete-rollback.md).

## LeetCode

This is the **system design exercise for a coding platform**, not a request
to solve algorithm questions. The public scope includes executing submissions
in multiple languages, returning results, isolation, and competition bursts.
[Published problem](https://www.hellointerview.com/learn/system-design/problem-breakdowns/leetcode)

Adapt the execution portion to CI checks or code-transformation jobs:

- **A job consumes unbounded CPU or memory.** Explain isolation, budgets, and
  how other tenants remain usable.
- **The worker is killed midway through execution.** Define durable results,
  artifact publication, cleanup, and retry semantics.
- **A result is attached to the wrong source revision.** Specify the identities
  that bind inputs, execution configuration, and output.
- **A transformation worker needs to change a repository.** Separate code
  execution from the authority to publish or deploy its result.

The published exercise's initial scope excludes fault tolerance. Add it for
the Netflix preparation because worker failure is central to the recruiter
brief. The [fleet change platform](problems/01-fleet-change-platform.md) provides
the surrounding workflow.

## Notification System

The public exercise is an internal platform with immediate and scheduled
delivery, bulk sends, and recipient preferences. Its outline explicitly
includes durable acceptance, high-priority traffic during a campaign, and
duplicate delivery.
[Published problem scope](https://www.hellointerview.com/learn/system-design/problem-breakdowns/notification-system)

Translate recipients into service owners and notifications into health or
migration findings:

- **A fleet-wide policy change creates thousands of notifications.** Preserve
  urgent delivery without allowing one producer to monopolize the platform.
- **A provider accepts a message but the acknowledgment is lost.** Distinguish
  logical deduplication from guarantees the provider cannot offer.
- **Ownership changes while a notification waits.** Decide when routing is
  resolved and how to avoid sending stale or inappropriate context.
- **The same health failure is observed repeatedly.** Define an incident or
  finding identity so observations do not become endless duplicate messages.

Pair it with the [technical health platform](problems/02-technical-health-platform.md).

## Use guided practice for both technical rounds

Complete an end-to-end design in the guided exercise, then use that architecture
for a failure discussion. For each important write or external call, explain
what happens if the process stops immediately before it, immediately after it,
or after the effect but before the acknowledgment.

Add a second pass with delayed messages, stale evidence, a hot tenant, and an
unavailable dependency. State the predicted observations, a discriminating
check, the safest mitigation, and the evidence required to declare recovery.
The [debugging method](debugging/README.md) provides the local reference.

Use written explanations and the pattern articles as support when a specific
gap appears. The most relevant supporting patterns are
[Multi-step Processes](https://www.hellointerview.com/learn/system-design/patterns/multi-step-processes)
and [Managing Long Running Tasks](https://www.hellointerview.com/learn/system-design/patterns/long-running-tasks).
They connect workflow state and background execution to the problems above.

## Topics that still need the local problems

The reviewed guided catalog has no directly named exercise for HIPR health
primitives, a service catalog with evidence quality, or a fleet modernization
campaign. Extend the guided problems with the first three
[local design problems](README.md#design-problems) to cover those responsibilities.

Consumer-product problems such as feeds, chat, geospatial matching, and video
delivery provide useful breadth, but have lower immediate relevance to this
role than the exercises above. Do not prioritize a video-service design solely
because the employer is Netflix. Likewise, the presence of agents in the role
does not make the ChatGPT exercise a substitute for reliable change execution.

If preparing only three guided exercises, choose **GitHub Actions, Job
Scheduler, and Metrics Monitoring**. For six, add **Ad Click Aggregator,
Distributed Cache, and Rate Limiter**. A useful stopping criterion for each is
being able to defend its state model, failure boundaries, recovery mechanism,
and remaining uncertainty without relying on the answer key.
