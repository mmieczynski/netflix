# Knowledge base for agent workflows

[Story selector](../story-bank.md) · [Culture framing](../culture-framing.md) · [Gaps](../gaps.md)

**Use for:** Your default second project: technical choice, developer needs, iteration and adoption.

**Research questions:** Q03, Q04, Q05, Q06, Q16; evidence to support Q02. **State:** Draft ready for rehearsal. No independent rehearsal recorded.

## Opening

> I designed and built a knowledge base that let agents locate internal guidance, then improved ingestion and retrieval based on how it was used.

## Answer structure

1. **Problem:** Agents had to discover where knowledge lived before searching Confluence, GitHub, OpsLevel and other sources.
2. **Design:** Ingestion structured documents, added metadata, assessed new information, organized a folder tree and regenerated indexes. Retrieval navigated indexes and relevant subtrees.
3. **Choice:** Documents instead of a vector database: similar observed retrieval time, easier maintenance, quality management and Claude-plugin distribution.
4. **Iteration:** Source-specific ingestion, fewer user questions and a retrieval-turn limit after long-tail latency.
5. **Use:** Built an API design review example; that workflow and the PR analyzer could retrieve internal guidance. Promoted it in presentations.
6. **Reflection:** Explore team-specific scopes after requests for narrower relevance.

## Frame for the question

For judgment, focus on the storage comparison and its limits. For user focus, lead with ingestion friction and changes you made. For impact, use the two concrete workflows, not unspecified productivity gains.

**Culture connection:** Judgment connects to the storage decision; Great and Always Better connects to revising ingestion and retrieval after observing friction. The AI label itself adds no evidence. See the [memo-based framing guide](../culture-framing.md) for the source and interpretation.

## Spoken draft

> I owned the design and implementation of a knowledge base for our internal AI workflows, alongside setting up and promoting an internal skills marketplace. Previously, agents had to search across sources such as Confluence, GitHub and OpsLevel. Before finding an answer, they first had to work out where the relevant information lived.
>
> We built an ingestion flow that structured documents, added metadata, checked for new information and organized it into a document tree with indexes. Agents used those indexes to locate relevant areas, search them and report what they could not find.
>
> A key decision was choosing document storage over a vector database. In our comparison, retrieval times were similar, while documents were easier for us to maintain, manage for quality and distribute as a Claude plugin.
>
> I continued improving the solution based on use. We added source-specific ingestion, reduced the number of questions asked during ingestion and bounded retrieval turns after seeing long-tail latency. To encourage adoption, I presented the solution and built an API design review workflow demonstrating it. That workflow could retrieve API documentation, and our internal PR analyzer could retrieve relevant review guidelines.
>
> Next, I would explore team-specific knowledge bases. The general solution was useful, but teams wanted knowledge with a narrower scope of relevance.

## Keep the claims precise

- Similar retrieval time applies to your comparison; no general superiority or equivalent retrieval quality was established.
- Bounding turns can affect completeness. Do not claim that reporting missing information prevents hallucinations.
- Presentations demonstrate promotion, not adoption counts. Marketplace setup is additional scope, not the main outcome of this story.

## Follow-up preparation

- What tasks and observations supported the storage comparison?
- How did you check retrieval usefulness or completeness after limiting turns?
