# Whiteboard preparation and study sequence

[Preparation index](README.md) · [Hello Interview priorities](hello-interview-priorities.md)

## Whiteboard content and answer depth

For the system-design whiteboard, make the following information visible as the
conversation develops:

| Put on the whiteboard | Why it helps |
|---|---|
| Main user operation and definition of success | Keeps the design connected to its purpose |
| Critical invariants | Makes correctness discussable |
| Synchronous versus asynchronous interactions | Exposes timeout, delay, and recovery behavior |
| Source of truth for each important fact | Prevents conflicting authorities |
| Key entities and identities | Makes retries, joins, and reconciliation concrete |
| Transaction and acknowledgment boundaries | Shows exactly what is guaranteed |
| One state machine for the central operation | Makes partial completion and cancellation understandable |
| Load assumptions and bottlenecks | Gives scaling decisions a basis |
| One failure-and-recovery walkthrough | Demonstrates that the architecture survives beyond its happy path |

Complete one end-to-end path before spending the rest of the discussion on an
isolated component. A technically detailed queue design does not compensate
for never explaining how the user learns whether their change succeeded.

Treat these statements as incomplete until the associated follow-up can be
answered:

| Initial statement | Follow-up to be able to answer |
|---|---|
| "We use retries." | Which failures are retryable, what is the budget, and what makes repetition safe? |
| "We use a distributed lock." | What prevents the previous owner from acting after losing ownership? |
| "The queue gives exactly-once processing." | Within which boundary, and what happens to an external effect after a crash? |
| "We use eventual consistency." | How stale may the information be, and which decisions are unsafe during that interval? |
| "We add a cache." | What happens during a miss storm, invalidation race, or source outage? |
| "We add more workers." | Is the bottleneck parallelizable, and can downstream systems accept the extra load? |
| "We roll back." | Which state is reversible, and how do we prove restoration? |
| "The dashboard is green." | Is the evidence current, complete, and attached to the correct version? |
| "We will monitor it." | Which metric detects this failure, and can the monitoring path fail in the same incident? |

## Preparation sequence and source reading

1. **Design the change platform end to end.** Explain campaign state, target
   identity, external operations, validation, and recovery.
2. **Apply the acknowledgment-loss and stale-worker scenarios to it.** Follow
   every write and external call through the failure.
3. **Design the health platform.** Work through missing sources, stale evidence,
   incorrect rule releases, and misleading compliance reports.
4. **Design the catalog and dependency system.** Explain evidence types,
   temporal identity, incomplete coverage, and bounded queries.
5. **Practice scheduler and ingestion variations.** Concentrate on fairness,
   cancellation, backpressure, schema evolution, and replay.
6. **Take an unfamiliar diagram and explain three distinct failure chains.**
   For each, identify a discriminating observation, a mitigation, and proof of
   recovery.

The highest-value reading order for this role is:

- [Fleet-wide change automation presentation](https://www.infoq.com/presentations/automate-fleetwide-changes/) — understand the domain and the boundaries between orchestration, validation, and delivery.
- [Maestro workflow-engine redesign](https://netflixtechblog.com/100x-faster-how-we-supercharged-netflix-maestros-workflow-engine-028e9637f041) — examine durable execution and recovery guarantees.
- [Service Topology architecture](https://netflixtechblog.com/building-service-topology-at-scale-architecture-challenges-and-lessons-learned-f4b792f3f0d8) — study concrete production bottlenecks and the evidence behind architectural changes.
- [Quasar blast-radius analysis](https://netflixtechblog.com/beyond-api-compatibility-building-a-tool-to-measure-the-blast-radius-of-library-changes-9474f2fd4f90) — examine how uncertainty must accompany change-impact information.
- [Effective Troubleshooting](https://sre.google/sre-book/effective-troubleshooting/) — practice turning a suspicion into a testable hypothesis.
- [Making retries safe](https://builder.aws.com/content/3Ev0BENTyBr0XxzRk5FDZzgNYos/making-retries-safe-with-idempotent-apis) — work through the contract needed when a caller cannot tell whether an operation happened.

For every practice answer, finish with: **what is authoritative, what remains
uncertain, what action is safe now, and what evidence would prove recovery.**
