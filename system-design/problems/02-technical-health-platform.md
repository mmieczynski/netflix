# Technical health platform design

[Preparation index](../README.md) · [Interview evidence](../interview-context.md)

Practice scenario for Netflix L5 Developer Platform. The follow-ups are
preparation questions, not confirmed Netflix interview questions.

> Design a system that assesses the technical health of services and
> repositories, explains the findings to owners, and exposes trustworthy
> context to developer tools and agents.

The first challenge is defining health. For this practice problem, possible
primitives include supported runtime versions, dependency status, ownership,
validation coverage, deployment readiness, and runtime reliability. These are
illustrative requirements, not a disclosed list of HIPR's rules.

Do not collapse everything immediately into a score out of 100. An average can
hide a critical failure, and a missing measurement can accidentally make a
service look healthier.

A useful conceptual separation is:

- **Evidence:** a build result, deployed artifact, dependency inventory, test
  result, or observation.
- **Evaluation:** a versioned rule applied to specific evidence.
- **Decision:** whether a particular action is appropriate given those
  evaluations.

For example, "this service has a healthy runtime" is different from "we have
sufficient evidence to automatically apply this dependency upgrade."

A plausible design ingests events and periodically reconciles source systems,
normalizes identities, evaluates rules, retains history, and serves current
results through an API and a user interface. Each result should identify its
subject, artifact or revision, rule version, evidence source, observation time,
and validity conditions.

| Follow-up | What to explain |
|---|---|
| **A source stops reporting. Does the service remain healthy?** | Evidence has a freshness contract. Expired or absent evidence becomes unknown where appropriate; it should not silently retain a green interpretation forever. |
| **An old successful build event arrives after a newer failed build.** | Associate results with the correct build or artifact. Arrival time is insufficient for deciding which result describes the current deployment. |
| **A rule changes. Must you reevaluate everything?** | Version the rule; determine which subjects are affected; backfill incrementally; distinguish old-policy from new-policy results during transition. |
| **A defective rule flags the whole fleet.** | Shadow evaluation, limited rollout of rule changes, comparison against previous behavior, and bounded downstream automation. Avoid allowing one bad evaluation to trigger unrestricted changes. |
| **Two sources disagree.** | Preserve provenance and define source authority by fact type. A repository declaration and runtime observation may describe different realities rather than one being universally wrong. |
| **A team disputes a finding.** | Show the evidence and rule explanation. Support accountable exceptions with an owner, reason, scope, and expiration where the product requires them. |
| **What does "not applicable" mean?** | Separate it from healthy, failed, unknown, and waived. Otherwise denominators and adoption reports become misleading. |
| **How do you prevent noisy notifications?** | Deduplicate findings, route to the current owner, prioritize actionable changes, and distinguish new failures from repeated observations of the same failure. |
| **How does an agent consume the result safely?** | Return structured evidence, identity, freshness, and uncertainty. An unexplained boolean encourages decisions that exceed what the evidence supports. |
| **What is the platform's SLO?** | Include evidence freshness, coverage, and evaluation delay, alongside API availability. A responsive API serving obsolete results can still fail its purpose. |

A concrete example to discuss:

> Version A passed validation yesterday. Version B deployed today. The
> collector is delayed, and the dashboard still shows A's green result.

The problem is not just "the cache is stale." The system has attached evidence
for the wrong artifact to a current decision. Trace identity and version
through collection, evaluation, storage, and display.

Also prepare the analytics side, because the role explicitly mentions
reporting. If asked, "What percentage of the fleet is compliant?", define:

- The population being measured.
- Whether retired, unowned, or unreachable services count.
- Whether unknown results are visible in the denominator.
- Whether the measurement concerns repository state or deployed state.
- Whether historical reports use the policy and inventory that existed at the
  time.

These decisions often matter more than the dashboard technology.
