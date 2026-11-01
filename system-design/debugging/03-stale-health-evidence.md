# Rollout decisions based on stale health evidence

[Debugging index](README.md) · [Preparation index](../README.md)

Practice scenario for the diagram-based distributed debugging round. These
questions are reconstructed preparation material, not a disclosed interview.

Example:

- A new version is deployed to an initial cohort.
- Its health pipeline is delayed.
- The controller reads the previous version's successful evaluation.
- It expands the rollout.

Prepare these questions:

- Which artifact does the evaluation describe?
- When was the underlying behavior observed?
- When was the event ingested and evaluated?
- Are required sources reporting?
- Is the dashboard healthy while its collectors are unhealthy?
- Did a rule or query change?
- Is there enough representative traffic to evaluate the cohort?

A plausible immediate mitigation is to stop admitting further rollout work
while preserving already running service capacity. Recovery requires
establishing the actual deployed versions and reevaluating them with
appropriate evidence.

The durable fix is not simply "lower the cache TTL." Bind evidence to the
relevant artifact and decision, specify freshness and coverage requirements,
and represent uncertainty explicitly.

Netflix has published an alerting incident in which evaluation delays caused
critical alerts to arrive roughly 45 minutes late. It is a useful reminder that
the observation system can fail separately from the service being observed.
[Improved Alerting with Atlas Streaming Eval](https://netflixtechblog.com/improved-alerting-with-atlas-streaming-eval-e691c60dc61e)
