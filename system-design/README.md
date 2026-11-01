# Netflix L5 system design and distributed debugging

[Repository index](../README.md)

Preparation for Developer Platform in Warsaw, focused on HIPR and Change
Systems. The original interview research was checked on **31 October 2026**;
Hello Interview recommendations were checked on **1 November 2026**.

Start with the [interview format and evidence](interview-context.md), then use
the individual problems below. The [Hello Interview guided practice priorities](hello-interview-priorities.md)
rank exercises with and without written solutions for this role. The [preparation guide](preparation-guide.md)
covers whiteboard expectations and the study sequence.

Publicly reported questions are identified in the evidence file. The local
design problems and debugging follow-ups are original practice scenarios.
These reference files are separate from the book's reading chapters and
conversation sessions; they do not record study completion.

## Preparation priorities

| Priority | Design to prepare | What it exercises |
|---|---|---|
| **First** | Fleet-wide migration or change platform | Durable workflows, external effects, retries, concurrency, validation, partial completion |
| **First** | Technical-health assessment platform | Evidence quality, freshness, versioning, rules, actionable results, safe automation |
| **First** | Service catalog and dependency/blast-radius system | Identity, graph modeling, incomplete information, ingestion, historical queries |
| **Next** | Durable job scheduler | Ownership, recovery, fairness, cancellation, resource accounting |
| **Next** | Event ingestion and shared client platform | API contracts, delivery semantics, backpressure, compatibility, developer experience |
| **Additional coverage** | Distributed counters and reconciliation pipelines | Atomic decisions, eventual consistency, duplicate events, finding silent omissions |

All practice prompts and follow-up questions below are reconstructions for
preparation. They are not claims about undisclosed Netflix interview questions.

## Design problems

| Order | Problem | Main preparation focus |
|---|---|---|
| 1 | [Fleet change platform](problems/01-fleet-change-platform.md) | Durable intent, execution, validation, pause and recovery |
| 2 | [Technical health platform](problems/02-technical-health-platform.md) | Evidence, rules, freshness, coverage and reporting |
| 3 | [Service catalog and dependency graph](problems/03-service-catalog-and-dependency-graph.md) | Identity, graph semantics, incomplete information and blast radius |
| 4 | [Durable job scheduler](problems/04-durable-job-scheduler.md) | Ownership, resource allocation, fairness and cancellation |
| 5 | [Event ingestion and client libraries](problems/05-event-ingestion-and-client-libraries.md) | Contracts, durability, backpressure, ordering and replay |
| 6 | [Distributed counting](problems/06-distributed-counting.md) | Atomic admission decisions and regional consistency |
| 7 | [Reconciliation pipeline](problems/07-reconciliation-pipeline.md) | Expected versus observed state and safe repair |

## Distributed debugging

The [debugging index](debugging/README.md) contains the troubleshooting method
and links to six detailed scenarios:

1. [Lost acknowledgment after a successful operation](debugging/01-lost-acknowledgment.md)
2. [Stale ownership and overlapping workers](debugging/02-stale-ownership.md)
3. [Rollout decisions based on stale health evidence](debugging/03-stale-health-evidence.md)
4. [High latency with apparently normal CPU](debugging/04-high-latency-normal-cpu.md)
5. [Queue growth that does not improve with more workers](debugging/05-growing-queue.md)
6. [Rollback that does not fully repair the incident](debugging/06-incomplete-rollback.md)

The [additional failure probes](debugging/failure-probes.md) cover cache races,
hot partitions, missing traces, failover, clock skew, competing controllers,
and other diagram-based follow-ups.
