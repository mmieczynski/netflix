# Internal agent platform

[Story selector](../story-bank.md) · [Culture framing](../culture-framing.md) · [Gaps](../gaps.md)

**Use for:** Platform architecture, SDK design, security collaboration and avoiding premature abstractions. Use as your technical alternative to S09.

**Research questions:** Q03, Q04, Q07, Q16; collaboration support for Q13. **State:** Draft ready for rehearsal. No independent rehearsal recorded.

## Opening

> I owned the SDK, event ingress and triggering logic for a platform validated with a Jira-triggered coding agent.

## Answer structure

1. **Workflow:** Jira ticket triggers a code implementer, with follow-ups on CI failures and PR comments.
2. **My scope:** Agent SDK, event ingress, trigger logic and guidance on connecting the human actor to agent tool access.
3. **Choice:** Passed requests through rather than introduce an event-to-agent-request mapping abstraction before a confirmed need.
4. **Collaboration:** Worked with security on auth design; sought agent-team feedback, held weekly syncs and reviewed their designs.
5. **Outcome:** A working agent could work on features.
6. **Reflection:** Start event ingress earlier; SDK restrictions required refactoring for some flows.

## Frame for the question

Explain one agent run first. Then discuss why the mapping abstraction was deferred and the cost of overly restrictive SDK choices. Make your contribution and the security team’s role clear.

**Culture connection:** Judgment connects to deferring an unneeded abstraction. Curiosity can be supported by seeking agent-team feedback, but do not claim you deliberately sought opposing opinions without a real example. See the [memo-based framing guide](../culture-framing.md) for the source and interpretation.

## Spoken draft

> I worked on an internal platform for creating and running AI agents with access to our internal tools. I owned the agent SDK, event ingress and triggering logic. We validated the platform with a code-implementation agent triggered by Jira tickets, with follow-up interactions from CI failures and PR comments.
>
> The platform accepted events and let agent authors decide whether to start a run based on the payload. They wrote agent logic as ordinary code using tools provided through our SDK. One decision I made was to pass requests through without introducing an event-to-agent-request mapping layer. I wanted to keep the design simple until we had a confirmed need for that abstraction.
>
> I also guided teammates on connecting the initiating human's identity with the agent's access to internal tools, collaborating with security on the design. I stayed in contact with the agent-development team through feedback discussions, weekly syncs and reviews of their design documents.
>
> We reached a working agent that could work on features. Looking back, I would introduce event ingress earlier to support more trigger sources. I would also make the SDK less restrictive: we had to refactor it when some agent flows did not fit its original constraints.

## Keep the claims precise

- Working capability does not establish broad adoption, autonomous merging or production reliability.
- Security collaboration alone does not establish specific authorization controls. Do not infer token propagation or impersonation details.

## Follow-up preparation

- Which SDK restriction blocked a flow, and what changed?
- What identity or authorization decision did you personally influence?
