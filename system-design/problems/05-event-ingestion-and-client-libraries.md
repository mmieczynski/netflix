# Event ingestion and client libraries

[Preparation index](../README.md) · [Interview evidence](../interview-context.md)

Practice scenario for Netflix L5 Developer Platform. The follow-ups are
preparation questions, not confirmed Netflix interview questions.

The reported L5 question was about ownership and integration strategy. Prepare
both that product judgment and the distributed-system mechanics beneath it.

> Teams need to emit events reliably from several languages. Design the
> ingestion service and decide what the platform team owns.

A strong answer to the SDK question should account for supported languages,
staffing, compatibility, required reliability behavior, and the cost of
inconsistent implementations. A maintained protocol plus a small number of
supported SDKs and conformance tests is one reasonable option; it is not
automatically the answer for every organization.

| Follow-up | What matters |
|---|---|
| **When is an event accepted?** | Define the acknowledgment's durability meaning. Being received in process memory and being durably stored are different contracts. |
| **The server stores the event but the acknowledgment is lost.** | Stable event identity, duplicate detection where required, and clear retention limits for deduplication. |
| **Only part of a batch succeeds.** | Define per-item outcomes or an atomic batch contract. Avoid ambiguous retries that repeatedly duplicate successful items. |
| **The collector is unavailable for hours.** | Bound client memory and disk use. Choose explicit behavior when buffers fill: blocking, rejecting, dropping selected events, or another policy appropriate to the data. |
| **A client library causes application latency.** | Bound work and queues, isolate background delivery, and avoid allowing telemetry retries to consume application resources indefinitely. |
| **Events arrive out of order.** | Specify whether ordering matters and at what scope. Global ordering is a much stronger requirement than ordering within one entity. |
| **A tenant dominates one partition.** | Show the key distribution and discuss partitioning changes without casually breaking required order. |
| **The event schema changes.** | Compatibility rules, schema/version identification, mixed producer versions, and replay behavior. |
| **You replay historical data.** | Protect live ingestion capacity and define whether sinks overwrite, append, deduplicate, or recompute. |

A useful impossibility to recognize: unlimited outage tolerance, zero event
loss, bounded local storage, and permanently nonblocking producers cannot all
be guaranteed simultaneously. The product must choose how to behave when
capacity runs out.
