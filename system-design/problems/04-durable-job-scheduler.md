# Durable scheduler design

[Preparation index](../README.md) · [Interview evidence](../interview-context.md)

Practice scenario for Netflix L5 Developer Platform. The follow-ups are
preparation questions, not confirmed Netflix interview questions.

This has direct support in a reported Netflix interview, and the underlying
ideas transfer to Change Systems.

> Design a service that accepts long-running jobs, schedules work onto
> available resources, exposes progress, and recovers from worker or scheduler
> failures.

Be able to distinguish:

- A **job definition** from one logical **run**.
- A run from its execution **attempts**.
- Scheduling a step from performing its external effect.
- A worker heartbeat from proof that useful progress is occurring.

The published candidate report explicitly describes probes around job types,
resource allocation, monitoring resources, crash handling, and multi-tenancy.
The following are additional practice questions:

| Follow-up | Required depth |
|---|---|
| **The scheduler crashes after assigning work.** | Explain where the assignment is durable, how a replacement discovers it, and how duplicate assignment is controlled. |
| **The worker disappears but later returns.** | Ownership generation, stale updates, and protection against obsolete work affecting current state. |
| **One tenant submits a million jobs.** | Admission control, quotas, fairness, queue isolation, and protection for existing work. A priority queue alone does not establish fairness. |
| **A high-priority stream never stops.** | Decide whether low-priority work may starve. If not, explain aging, reserved capacity, or a fair scheduling policy. |
| **A resource is reported free but another scheduler also chooses it.** | Make reservation authoritative and atomic within the chosen coordination domain. Resource observations can be stale. |
| **Cancellation races with completion.** | Define which transitions are valid and what the user sees. "Cancellation requested" is distinct from "work stopped." |
| **A failed job is retried from the beginning.** | Identify repeatable steps, checkpoints, and steps needing idempotency or compensation. |
| **The API connection closes.** | Job lifecycle should follow its explicit contract, not accidentally depend on the client keeping a connection open. |

For an ML-specific prompt, additionally discuss resource shapes, long-running
versus shorter jobs, checkpointing, and whether a job needs several resources
simultaneously. For this role, master durable execution and fairness before
specializing in GPU placement.
