# Frame your stories against Netflix culture

Use your decisions and behavior as evidence. The connections below are coaching interpretations, not a private Netflix scorecard. The [official culture memo](https://jobs.netflix.com/culture), checked 14 October 2026, describes autonomy with accountability, seeking other perspectives, constructive feedback and ongoing improvement. It supports useful process, rather than treating all process as undesirable.

## Your strongest connections

| Culture theme | Your evidence | How to frame the answer |
| --- | --- | --- |
| Judgment | [Knowledge base](stories/knowledge-base.md): documents versus vectors; [schema-first](stories/schema-first.md): compatibility costs | Explain the requirements, alternative and accepted cost. Keep observations scoped to your project. |
| Selflessness | [Adoption](stories/adoption.md): implemented initial endpoints so the team could continue | Show the work you did to make another team's adoption feasible. |
| Candor | [Quota incident](stories/quota-incident.md): original assumption and personal response; [feedback](stories/feedback.md): acknowledged pressing too hard | State your part plainly and give colleagues their credit. |
| Curiosity | [Feedback](stories/feedback.md): exploring concerns; [agent platform](stories/agent-platform.md): seeking user feedback | Describe what you sought to understand. Add a changed decision only when it actually happened. |
| People Over Process | [Linter](stories/linter.md) and [schema-first](stories/schema-first.md) | Explain the developer benefit of the checks, including adoption cost. More governance is not itself the result. |
| Great and Always Better | [Knowledge base](stories/knowledge-base.md): ingestion and retrieval revisions; [agent platform](stories/agent-platform.md): SDK refactor | Name the specific limitation you observed and how the design changed. |

These themes and their meaning come from the [culture memo](https://jobs.netflix.com/culture). The project evidence comes from your account.

## Culture traps in these particular stories

**API governance:** “Every change now goes through review” leaves the benefit unexplained. Your supported result is that design issues are caught during contract review before implementation. Explain that benefit, then acknowledge adoption effort. Do not claim reviews replaced other approvals or reduced measured rework without evidence.

**Being right:** “I was objectively right” can obscure the failure to understand the other person. Your account supports “I was convinced my technical position was sound, but continuing to press it was not helping.” You do not need to pretend you changed your technical opinion.

**Autonomy:** Starting the linter yourself is useful evidence. Consulting your manager about introduction does not cancel that initiative. Explain the decision you made and the coordination it required.

**Simplicity:** The document store and deferred event mapping are choices with constraints. Explain why each suited the actual use case; avoid implying that simpler designs always win.

**The new consumer:** In the quota story, the surprise explains the traffic change. Your responsibility remains the original design and contribution to recovery. Avoid making the other team the villain.

**AI novelty:** Use working capabilities, decisions and limitations. Neither project needs claims about transformation, production scale or productivity that have not been established.

## Connections we cannot yet substantiate

The memo describes seeking opposing views and supporting a decision after disagreement. The [SDK priority disagreement](stories/priority-disagreement.md) now supplies a concrete example of supporting the selected work after arguing for another priority. Deliberately seeking opposing opinions remains less established. Normal feedback discussions and adoption support should not be relabeled as those examples. We also lack an inclusion-specific event; helping another team alone does not establish that. [Culture memo](https://jobs.netflix.com/culture).

Your preferences are recorded below. You do not need to invent discomfort; use a genuine question about how support and decision context work on this team.

## Talking points for your culture answer

Choose two or three points you actually believe, then connect each to an example. These are suggested angles, not attributed personal preferences.

- **Autonomy with clear context:** the linter shows self-directed work; the priority-feedback episode illustrates why you want clarity about goals and conflicting signals.
- **Regular, useful feedback:** your proactive approach and conversation with a senior peer can support this; you report that stakeholder alignment resolved the issue. Pair it with your own received-feedback story to show both directions.
- **Letting go of owned work:** the knowledge-base handover supplies a real priority change. Explain the dependency rationale and your documentation/PR support during handover.
- **Learning through revision:** source-specific ingestion, fewer questions and SDK refactoring give concrete examples of revising an initial design.

A useful question for the manager: “When senior stakeholders disagree about priorities, how does this team clarify who makes the final decision and communicate it?” The memo emphasizes a responsible decision owner rather than a committee; alignment discussions do not require unanimity. [Official memo](https://jobs.netflix.com/culture), checked 19 October 2026.

No need to claim every principle feels easy. If true for you, the practical adjustment could be learning how much consultation a decision needs before acting. This is an optional reflection, not an established concern of yours.

## Your stated culture preferences

You connect with ownership and autonomy supported by teammates, who you report trust you to work that way. You also value continuous feedback because it complements autonomous work.

Suggested wording based on that preference:

> What resonates most with me is ownership and autonomy with support. I tend to take responsibility for work and drive it independently, and my teammates trust me to do that. I also value continuous feedback. To me, it complements autonomy: it helps me recognize when my judgment or approach needs adjusting while I remain responsible for the work.

Use the linter for independent initiative, priority feedback for clear context, and the [responsiveness feedback story](stories/feedback.md) for an actual adjustment. These are your preferences and examples, not a claim that autonomy means working in isolation.
