# Quota endpoint rate limit incident

[Story selector](../story-bank.md) · [Culture framing](../culture-framing.md) · [Gaps](../gaps.md)

**Use for:** Your default mistake, failed assumption and operational responsibility story.

**Research questions:** Q12, Q16; supporting material for Q18. **State:** Draft ready for rehearsal. No independent rehearsal recorded.

## Opening

> My quota endpoint relied on expected traffic; a new consumer increased demand roughly a hundredfold and exposed insufficient protection against upstream limits.

## Answer structure

1. **Decision:** Calculated expected traffic below the upstream limit and omitted caching.
2. **Failure:** Months later, a new consumer drove roughly 100 times the traffic; rate limiting caused endpoint failures.
3. **My response:** Helped debug, provided a mitigation option and approved the change implemented by a colleague.
4. **Recovery:** Team added in-memory memoization; a longer-term internal cache was designed.
5. **Lesson:** Consider changing demand, freshness, upstream failure behavior and monitoring, beyond an initial traffic estimate.

## Frame for the question

Own the assumption and explain why it seemed reasonable. Keep the new consumer as context, not the culprit. Finish with the difference between the immediate mitigation and your broader lesson.

**Culture connection:** Candor connects to acknowledging your design’s limitation. Judgment is in the revised reasoning; be clear which improvements are future intentions. See the [memo-based framing guide](../culture-framing.md) for the source and interpretation.

## Spoken draft

> I designed an endpoint for tracking individual users' Claude Code quota usage. It queried an external service with rate limits. I estimated that expected traffic was comfortably below those limits and chose not to add caching.
>
> A few months later, another tool started using the endpoint, increasing traffic roughly a hundredfold. We hit the upstream rate limits and our endpoint began failing. My design had relied on the original traffic assumptions, and it did not handle that increase successfully.
>
> I helped debug the issue and provided a mitigation option. A colleague implemented the change, which I then approved. The team introduced in-memory memoization as the immediate workaround and designed a longer-term solution to cache the information internally, reducing dependence on upstream calls as traffic increased.

> My takeaway was that estimating expected traffic is not enough when depending on a rate-limited service. My calculation reflected the consumers we had at the time, but the endpoint could gain new consumers and traffic could change significantly. Next time, I would explicitly consider how to reduce upstream calls, how fresh the quota information needs to be, how the endpoint behaves when the upstream limit is reached, and monitoring that helps spot pressure before requests fail. The lesson was to make changing demand and dependency limits part of the design rather than relying only on the initial traffic estimate.

## Keep the claims precise

- A colleague implemented the fix. The longer-term cache was designed, not confirmed shipped.
- Do not claim a later application of the lesson or that caching guarantees protection from every spike.

## Follow-up preparation

- Optional follow-up: who was affected and for how long? How much quota staleness could the mitigation allow?
