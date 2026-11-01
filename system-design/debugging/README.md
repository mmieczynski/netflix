# Distributed debugging preparation

[Preparation index](../README.md) · [Interview format](../interview-context.md)

The recruiter describes a supplied system diagram and a discussion of failure
modes, troubleshooting, and mitigation. Start with the method below, then
follow the scenarios. All scenario questions are original practice material.

## Troubleshooting method

First, explain the intended system back to the interviewer:

> "An accepted request creates durable intent here. This queue carries work
> to these workers. The external system performs the change. These events
> update our view, and this component decides whether execution can continue.
> I want to clarify which acknowledgment means durable acceptance and which
> store is authoritative for actual completion."

Then identify the important invariant. Examples:

- An accepted change must remain discoverable after a restart.
- A stale worker must not overwrite a newer decision.
- A target must not be marked verified using another artifact's evidence.
- A paused campaign must not keep admitting new work.
- Missing telemetry must remain distinguishable from healthy behavior.

For each suspected failure, use the sequence:

**Observed symptom → plausible mechanism → predicted evidence → discriminating
check → mitigation → verification.**

"Check the logs" is incomplete. State which records are needed, which
identifiers connect them, and what different results would mean. Google's
troubleshooting chapter is useful grounding for this hypothesis-driven
approach.
[Effective Troubleshooting](https://sre.google/sre-book/effective-troubleshooting/)

## Scenarios

- [Lost acknowledgment after a successful operation](01-lost-acknowledgment.md)
- [Stale ownership and overlapping workers](02-stale-ownership.md)
- [Rollout decisions based on stale health evidence](03-stale-health-evidence.md)
- [High latency with apparently normal CPU](04-high-latency-normal-cpu.md)
- [Queue growth that does not improve with more workers](05-growing-queue.md)
- [Rollback that does not fully repair the incident](06-incomplete-rollback.md)
- [Additional failure probes](failure-probes.md)
