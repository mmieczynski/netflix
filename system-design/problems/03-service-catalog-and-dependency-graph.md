# Service catalog and dependency graph design

[Preparation index](../README.md) · [Interview evidence](../interview-context.md)

Practice scenario for Netflix L5 Developer Platform. The follow-ups are
preparation questions, not confirmed Netflix interview questions.

> Given a proposed library or service change, identify potentially affected
> consumers, their owners, and the evidence supporting that assessment.

There are several different graphs hiding inside this requirement:

| Relationship | What it tells you | What it does not establish |
|---|---|---|
| Declared package dependency | A project includes a library/version | That it invokes the changed API |
| Semantic code reference | Code refers to a particular symbol | That the path executes in production |
| Runtime service call | A call was observed during a period | That every possible dependency was observed |
| Deployment relationship | An artifact/version is running somewhere | That its repository's current main branch matches it |
| Ownership relationship | A team is associated with an entity | That ownership information is current |

Netflix's service-topology publications describe combining network
observations, application metrics, and tracing because their coverage and
context differ. That is useful preparation for explaining why one source
cannot answer every dependency question.
[Why Netflix built Service Topology](https://netflixtechblog.com/from-silos-to-service-topology-why-netflix-built-a-real-time-service-map-0165ba13a7bc)

Netflix's Quasar article is also highly relevant: it connects library API
changes to observed consumer references and explicitly discusses the limits of
incomplete or stale indexing. Read it to understand the distinction between
compatibility changes and evidence of downstream use.
[Measuring the blast radius of library changes](https://netflixtechblog.com/beyond-api-compatibility-building-a-tool-to-measure-the-blast-radius-of-library-changes-9474f2fd4f90)

For the practice design, model identities and edges before picking a graph
database. An edge may need its type, version scope, source, observation
interval, and confidence or coverage information. A relational edge table can
be adequate for bounded queries at modest scale; a specialized graph store
becomes easier to justify when traversal patterns and volume demand it.

Prepare this follow-up chain:

1. **"No consumers were found. Can we remove the API?"** Explain what was
   searched and what remains unobserved. A result can be useful without proving
   absence.
2. **"A service has not called its disaster-recovery dependency recently."** A
   short runtime observation window misses rare paths. Combine evidence and
   expose uncertainty.
3. **"The library is wrapped by another library."** Distinguish direct
   consumers from transitive exposure. Explain how to traverse relationships
   without claiming every reachable service will necessarily fail.
4. **"The current graph differs from the graph during yesterday's incident."**
   Preserve temporal information. Reconstruct topology and deployed versions
   for the relevant period.
5. **"An IP address was reused."** Resolve observations against time-scoped
   instance identity. Joining old traffic to today's IP owner can invent
   dependencies.
6. **"A shared component has an enormous number of consumers."** Bound
   traversals, paginate, report truncation, and consider precomputation. Explain
   how to avoid a high-degree node overwhelming ingestion or queries.
7. **"Events arrive after a service was deleted."** Use identity generations,
   authoritative lifecycle information, or suitable tombstones so late
   observations do not resurrect an obsolete entity.
8. **"The ingestion pipeline falls behind."** Expose freshness degradation.
   Discuss buffering and backpressure, retention limits, and how downstream
   decisions change when the graph is old.

The implementation article on Service Topology is particularly valuable for
debugging preparation: it discusses consumer lag, uneven load, memory pressure,
and separating processing stages. Read it as a study of diagnosing bottlenecks
and revising an architecture under real load.
[Building Service Topology at Scale](https://netflixtechblog.com/building-service-topology-at-scale-architecture-challenges-and-lessons-learned-f4b792f3f0d8)
