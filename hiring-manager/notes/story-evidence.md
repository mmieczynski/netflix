# Archived story development notes

Historical evidence and coaching notes preserved on 14 October 2026. Some assessments below are superseded. Use the [current selection guide](../story-bank.md) and story cards for preparation.

# Netflix hiring manager story bank

Working notes for preparing the L5 Developer Platform manager conversation. The first DSA round is passed, as you reported on 13 October 2026. Three coach-authored drafts are recorded below: the main project, an adoption discussion and the self-initiated linter. All come from the same initiative. No story has been rehearsed yet.

Use the [research chapter](../reading/13-hiring-manager.md) for questions and answer criteria, and the [conversation session](../sessions/13-hiring-manager.md) for coaching. Stories can be developed entirely aloud. The coach retrieves these files, asks one question at a time, and maintains the notes when write access is available. No editing or copying prompts is required from the learner.

## Interview context

- Position: Senior Software Engineer 5, Developer Platform, Warsaw, JR32657.
- [Supplied position](https://explore.jobs.netflix.net/careers/job/790301977513).
- Recruiter agenda: experience and projects, stack, motivation, recurring culture discussion, concrete contribution in complex projects, and questions for the manager.
- Specific team, manager, date and duration: not yet supplied.
- Personal career and project facts: schema-first initiative, responsibilities, technical tradeoff, migration and qualitative outcomes recorded below; reflection is recorded, while rollout details and specific adoption interactions remain open. DSA practice is not evidence of production experience.

## Story inventory

Unfilled rows are slots to help us find coverage, not events attributed to you. Several slots may use different episodes from the same project. Replace a slot only with an event you actually describe.

| ID | Event to find | Related guide questions | Current state |
| --- | --- | --- | --- |
| S01 | Schema-first API initiative and validation library | Q03, Q04, Q06; coverage to establish | First complete draft prepared; not rehearsed |
| S02 | Internal AI agent platform | Q03, Q04, Q06, Q13 | Complete first draft; technical follow-up story |
| S03 | Adoption discussion about schema synchronization | Q05, Q13 | First complete draft; initial batch and team completion confirmed |
| S04 | Self-initiated schema linting within the same initiative | Q07, Q14; coverage to establish | First complete draft; not rehearsed |
| S05 | Quota endpoint exceeded upstream rate limits | Q10, Q12, Q16 | Complete first draft; reflection agreed, not rehearsed |
| S06 | Feedback about pressing arguments too hard | Q08, Q10, Q18 | Preliminary feedback draft; original context open |
| S07 | Difficult feedback given | Q08, Q11 | Not collected |
| S08 | Helping a colleague or choosing a wider benefit | Q14, Q15 | Not collected |
| S09 | AI knowledge initiative and internal skills marketplace | Q02, Q05, Q06, Q13 | Complete first draft; recommended second project |

## Coverage plan after the pacing correction

The learner correctly noted that the extended discussion had mostly developed one project. S01, S03 and S04 are three answer angles from that initiative, not three independent projects. S06 is a preliminary feedback account. Treat the current coverage as one developed project plus a partial feedback story.

Prepare four core examples before further detailed polishing:

| Example | Current position |
| --- | --- |
| Main project with ownership, judgment and impact | S01 ready as a first draft; S03/S04 provide adoption and initiative follow-ups |
| A different substantial project | S09 recommended default second project; S02 also drafted for technical follow-ups. |
| A meaningful mistake or setback with personal responsibility | S05 quota endpoint incident captured; first draft and assessment |
| Critical feedback and growth | S06 preliminary draft available; original context remains open |

For each new example, invite a rough overview of the problem, personal contribution and result. Draft from that overview, give a candid assessment, and use at most one or two essential follow-ups during the first pass. Record other gaps for later rehearsal. Preserve one conversational question at a time, but make prompts broad enough to avoid extracting every detail through separate turns. The user prioritizes breadth and useful feedback over exhaustive refinement of one story.

## Coaching assessment after the first story-development pass

The learner requested explicit feedback on progress and suitability for the Netflix L5 role. The coach had spent too long eliciting details without sufficiently explaining their value or judging the resulting answers. The assessments below use the recruiter brief and the [research chapter](../reading/13-hiring-manager.md); they are preparation judgments, not Netflix's internal scorecard or a hiring prediction.

| Story | Current assessment | Strongest supported evidence | Main improvement |
| --- | --- | --- | --- |
| S01 Schema-first initiative | Strong project choice; draft ready for rehearsal | Developer problems, personally delivered tooling, a version choice with migration costs, support for other teams, qualitative outcomes and practical reflection | Lead with the problem and own contribution. Keep the teammate's converter as a follow-up detail. Technical implementation depth and performance outcomes have not yet been demonstrated in conversation. |
| S03 Adoption discussion | Strong collaboration episode; draft ready for rehearsal | Connected adoption to a team's schema-synchronization problem, implemented the first batch to address upfront effort, and the team completed the remaining endpoints | Emphasize the concrete support and completed handover. Keep contract testing as the proposed longer-term benefit until delivery is confirmed. |
| S04 Self-initiated linter | Credible initiative episode; draft ready for rehearsal | Identified a gap between schema adoption and consistent API design, researched an approach, started work, aligned it with architecture's direction and involved the manager in planning introduction | Emphasize the decision to automate guidelines and the delivered linter. The manager discussion supports coordination; it is not the central achievement. Detailed rollout is available as a future follow-up. |
| S06 Critical feedback | Relevant growth material; framing and context need work | Acknowledged pressing arguments too hard, received meaningful feedback, and reports an ongoing change across later discussions | Explain technical confidence without treating correctness as established. Describe the listening behaviour clearly. Original event context and initial reaction are still unspecified; a later example may help but is not a prerequisite to continuing. |

**L5 relevance:** The collected experience fits the advertised senior individual-contributor work: independently building useful developer tooling, reasoning through compatibility costs, coordinating adoption, enabling other teams and considering operational constraints. S03 is currently the clearest compact behavioural example because it contains both a concrete personal action and a completed team outcome. Actual performance measurements and detailed implementation reasoning have not been assessed; requirements about low overhead are not achieved results.

**Remaining preparation needs:** Three drafts come from one initiative, so they offer different angles on a project rather than broad coverage of every behavioural question. The current outcomes are useful qualitative evidence; do not manufacture metrics. Spoken delivery and response to realistic follow-ups remain unassessed because no complete answer has been rehearsed. Build broader coverage next, then practise the selected answers. Avoid further refinement of the schema-first project during this first coverage pass.

**Coaching change requested by the learner:** Give direct feedback on what an answer demonstrates, its single most important weakness and useful wording. Explain the purpose of any follow-up. Avoid extended chains of questions and repeated demands for increasingly granular detail when enough material exists for a draft. Do not treat requests for feedback as another question to answer.

## S01 Schema-first API initiative

**Evidence state:** First complete coach-authored draft prepared from learner-supplied facts, including reflection, on 13 October 2026. Not rehearsed or independently retold. Quantitative impact and several follow-up details remain unestablished.

**Facts supplied by the learner:**

- Selected this as a project they enjoyed and want to develop into an interview story.
- The initiative changed how APIs were written, moving from code-first to contracts-first.
- Both public and internal endpoints were poorly documented.
- Internal developers integrating with another endpoint encountered incomplete schemas expressed through different libraries, or no schema at all. They then had to inspect implementation code or ask the owning team to understand the contract.
- Request and response validation practices were fragmented and varied greatly between teams.
- Public schemas were also in poor shape; raw OpenAPI schemas were not published because of their condition.
- Longer-term goals included breaking-change detection, SDK generation and contract testing. Breaking-change detection was subsequently reported as introduced, and the learner reports that their work enabled SDK teams to introduce generated SDKs. Contract testing delivery remains unconfirmed.

**Personal contribution reported by the learner:**

- Was responsible for tooling intended to make the transition easier for developers and give them value from switching.
- Shipped a lightweight library for automatic request and response validation using OpenAPI definitions.
- Created a schema linting tool from the ground up, with rules intended to standardize OpenAPI schemas across the company. The rules came from API guidance documents authored by an architect; the learner's contribution was the tooling, not authorship of those guidelines.
- Introduced a breaking-change detection tool for the main API services in collaboration with the team owning the service. Whether the tool itself was existing or newly built, and the exact integration mechanism, have not been described.
- Promoted the approach, helped other teams switch and worked to persuade them to adopt it. A specific adoption interaction and its outcome have not yet been described.
- Responsibility also included automatic type-generation tooling from OpenAPI schemas. The generated language, implementation details and rollout scope have not been described. Type generation does not establish delivery of complete SDK generation.

**Validation-library decision and constraints reported by the learner:**

- Identified the validation-library design as the hardest technical area.
- The work involved a decision to move to OpenAPI 3.1; the largest API was using OpenAPI 2.0 at that point. Who made the final version decision and the migration scope remain to be established.
- The reason for choosing 3.1 was its full JSON Schema support at the time, allowing use of existing open-source validators such as Ajv. The learner named Ajv as an available option; the actual validator and version shipped have not yet been confirmed.
- The cost included upgrading the largest API's definitions from OpenAPI 2.0 and moving many smaller APIs from OpenAPI 3.0. Required migration scope is reported; completion, effort and rollout sequence have not yet been established.
- Other tooling supported older OpenAPI versions but did not support 3.1 at the time. The learner said finding alternatives was the main response, then described the homegrown generator and conversion work below as a concrete example.
- The library needed to work with multiple API frameworks used across the company.
- It needed to validate requests and responses against their schemas to build confidence that schemas matched service behaviour.
- Latency and memory overhead were constraints, including concern about increased operating costs. These are stated requirements; no performance measurements or achieved overhead have been supplied.
- Separately, the learner identified introducing and implementing linting without prompting from management as a decision they initiated. That episode is recorded as S04 below.

**Migration example reported by the learner:**

- A homegrown tool generated OpenAPI schemas from implementation code, but only produced OpenAPI 2.0. The learner's earlier answer about JavaScript objects describing schemas concerned this code-first generation path.
- Moving away from this tool to schema-first development required a transition period with both approaches enabled.
- The learner clarified that the teammate wrote the migration tool. It converted the OpenAPI 2.0 schema into 3.1 and merged the already existing schema-first endpoint definitions into it. The learner previously reported working with this teammate; their specific contribution to this tool remains unspecified.
- The change also involved checking the schema into source control. The learner viewed this as beneficial; a specific review or workflow benefit has not yet been described.
- Source-of-truth rules during coexistence, conflict handling and migration completion have not yet been described.

**Correction to earlier coaching notes:** The coach conflated the old generator with the teammate's migration tool. JavaScript schema objects were not the reported input boundary of the migration tool; OpenAPI 2.0 was. The version conversion, merge and teammate's authorship above supersede the earlier record.

**Technical wording check — coach research, 13 October 2026:** OpenAPI 3.1 bases its Schema Object on JSON Schema Draft 2020-12, while OpenAPI 3.0 used a restricted and extended form of an earlier draft. Ajv documents support for Draft 2020-12, earlier drafts and selected OpenAPI keywords. The rationale is standards alignment and access to existing schema validators; it does not establish the exact validator configuration or full HTTP integration used in this project. [OpenAPI 3.1 Schema Object](https://spec.openapis.org/oas/v3.1.0.html#schema-object), [OpenAPI 3.0 Schema Object](https://spec.openapis.org/oas/v3.0.3.html#schema-object), [Ajv schema support](https://ajv.js.org/json-schema.html).

**Outcomes and feedback reported by the learner:**

- Teams gave positive feedback that integrating with APIs became easier and quicker. This is qualitative feedback; no measured time reduction, population or specific integration example has been supplied.
- All contributions to the main API now go through a governance process that reviews the contract. The learner reports that API design issues are caught in this review, before implementation. Who established or owns that process has not been specified; no measured reduction in rework has been supplied.
- Request and response validation was valued as a significant benefit.
- The architecture team was pleased that the linter enabled API standardization. Company-wide compliance or coverage has not been quantified.
- The learner's work enabled SDK teams to introduce SDKs generated from OpenAPI schemas. Credit SDK delivery to those teams; the learner's contribution was enabling it through the schema/tooling work.
- Breaking-change detection was introduced in collaboration with the team owning the large API, and that team appreciated it. No specific prevented incident or detected incompatibility has been described.

**Adoption concerns reported by the learner:**

Concerns varied between teams; these are not all attributed to a single team or interaction.

- Additional toil and process on top of other changes teams already had to adopt.
- Performance and stability risks.
- Doubts about the value of the solution.

**Specific adoption episode:** One team needed to synchronize schemas between its API and the main API. The discussion connected schema adoption to the possibility of contract testing to improve that integration. See S03 below for the developing story; the learner implemented the first batch, and the team completed the rest. Contract-testing delivery remains unconfirmed.

**Reflection supplied by the learner:**

- Would explore framework-specific libraries on top of the generic validation library to reduce boilerplate and improve developer experience. This is a proposed improvement, not reported delivered work.
- Would engage the CI team to explore delivering the tools through the standard CI pipeline.
- Currently the tools run in GitHub Actions separately from the main delivery pipeline. Standard-pipeline integration has not been reported as implemented.

**Open facts:** Ownership of the governance process and any specific design issue caught in contract review; a concrete integration example and the source/context of the feedback; adoption scope and any available measurements; the learner's specific collaboration on the teammate-authored converter; source of truth and consistency during coexistence; alternatives actually considered and the learner's role in the version choice; concrete library architecture; decision authority and ownership boundaries; how original problems were discovered or substantiated; implementation languages, validator choice and versions; migration effort and completion; specific collaborator interactions, rollout and performance evidence. Qualitative benefits are now reported; quantitative impact is not established.

**First complete draft — coach-authored, not rehearsed:**

> One project I enjoyed was moving our APIs from code-first to contract-first development. Documentation was inconsistent, and developers often had to read another service's code or ask its owning team to understand the contract. Request and response validation also varied across teams.

> My responsibility was to make the switch practical and valuable for developers. I shipped a generic OpenAPI-based validation library, built a schema linter using an architect's API guidelines, and worked on type-generation tooling. I also helped teams adopt the approach and introduced breaking-change detection together with the main API team. I initiated the linting work myself.

> A difficult choice was targeting OpenAPI 3.1 for its JSON Schema alignment, which allowed us to use existing validators. The cost was migrating APIs from 2.0 and 3.0 and dealing with incompatible tooling. During the transition, my teammate wrote a converter that upgraded the old 2.0 schema to 3.1 and merged in existing schema-first endpoints, supporting coexistence of both approaches.

> Teams told us integrations became easier and quicker, and they valued automatic validation. The architecture team appreciated the standardization enabled by linting. My work also enabled the SDK teams to introduce generated SDKs. All contributions to the main API now go through contract review, where design issues are caught before implementation.

> Looking back, I'd explore framework-specific libraries on top of the generic core to reduce integration boilerplate. I'd also work with the CI team to explore incorporating the tools into our standard delivery pipeline; currently, they run through separate GitHub Actions. Those would be my next improvements to the developer experience.

**Follow-up preparation:** S03 below develops a specific adoption interaction from this project. The complete S01 draft remains unrehearsed.

## S02 Internal AI agent platform

**Evidence state:** Complete first coach-authored draft from learner facts; not rehearsed. Working agent confirmed, broad production adoption and reliability results unestablished.

**Facts supplied:**

- Internal platform for creating/running agents with access to internal tools. Main ownership: agent SDK, event ingress and trigger logic.
- Validation workflow: a code implementer triggered by Jira tickets, also able to follow up on CI failures and PR comments. A working agent can work on features; no autonomous merge or delivery-success rate claimed.
- Jira events enter the platform; agent authors decide whether to trigger from the payload. Agent logic is ordinary code using platform-supported tools through the SDK.
- Chose to pass plain requests through instead of adding an event-to-agent-request mapping abstraction, deferring complexity until a confirmed use case. Exact request representation is unconfirmed; do not invent it.
- Guided teammates on authentication connecting the human actor to agent tool access across entry points; collaborated with security on the design. Specific controls and authorization semantics remain unspecified.
- Periodically sought agent-team feedback, held weekly status syncs and reviewed their design documents for agent improvements.
- Would start event ingress earlier to enable more trigger sources. SDK proved too restrictive for some flows and was refactored; would design it with more flexibility. Exact restrictions and learner's role in refactoring unestablished.

**First draft — coach-authored:**

> I worked on an internal platform for creating and running AI agents with access to our internal tools. I owned the agent SDK, event ingress and triggering logic. We validated the platform with a code-implementation agent triggered by Jira tickets, with follow-up interactions from CI failures and PR comments.
>
> The platform accepted events and let agent authors decide whether to start a run based on the payload. They wrote agent logic as ordinary code using tools provided through our SDK. One decision I made was to pass requests through without introducing an event-to-agent-request mapping layer. I wanted to keep the design simple until we had a confirmed need for that abstraction.
>
> I also guided teammates on connecting the initiating human's identity with the agent's access to internal tools, collaborating with security on the design. I stayed in contact with the agent-development team through feedback discussions, weekly syncs and reviews of their design documents.
>
> We reached a working agent that could work on features. Looking back, I would introduce event ingress earlier to support more trigger sources. I would also make the SDK less restrictive: we had to refactor it when some agent flows did not fit its original constraints.

**Assessment:** Useful for technical scope, judgment about abstractions, collaboration and reflection. The SDK correction is candid evidence of learning, not automatically evidence of poor judgment. Impact currently means a working capability; do not claim broad adoption, secure-by-proof operation, production maturity or productivity improvements.

**Optional rehearsal probes, not required before moving on:** Why mapping was unnecessary for initial consumers; specific SDK constraint and change; identity/authorization decision personally influenced; observed agent reliability and usage. Credit security and agent authors accurately.

## S03 Adoption discussion about schema synchronization

**Evidence state:** First complete coach-authored draft prepared from learner-supplied facts. The learner implemented the initial set of endpoints and the team completed the remainder. Part of S01; not rehearsed or independently retold.

**Facts supplied by the learner:**

- A team was concerned about the extra work and process involved in adopting the new approach.
- The learner discussed the concerns and explained why the longer-term benefit could justify the upfront cost.
- This team had a problem keeping schemas synchronized between its API and the main API.
- The learner and colleagues explained that having schemas at both endpoints would enable introducing contract testing, with the aim of improving integration with the main API.
- Contract testing was presented as a future capability. No implemented tests or realized improvement from them have been reported.
- They agreed that the learner would implement the first batch of endpoints to introduce the pattern, and the team would then pick up the remaining endpoints. The learner confirms they implemented the initial set and the team completed the rest.

**First complete draft — coach-authored, not rehearsed:**

> One team was concerned about the extra work involved in adopting schema-first development. They had a problem keeping their schemas synchronized with the main API, so we discussed how schemas on both sides would enable contract testing to improve that integration.
>
> To help with the upfront work, we agreed that I would implement the first batch of endpoints to establish the pattern, and they would take over the rest. I implemented that initial set, and the team then completed the remaining endpoints.

**Open facts for follow-up:** Whether contract testing was introduced; the scale and period of this team's migration; any ongoing support needed; how the synchronization problem was discovered; and the learner's reflection. Completion of the team's remaining endpoints is confirmed; it does not establish company-wide migration or contract-testing delivery.

**Coaching observation:** The outcome completes the adoption episode: the learner's argument addressed a team-specific problem, the agreement provided initial implementation support, and the team followed through on its remaining work. This supports discussion of collaboration and adoption. No claim of zero ongoing support, measured savings or delivered contract testing is added.

**Last answered question:** How did that agreement play out in practice? The learner implemented the initial set; the team did the rest.

**Next step:** Keep this draft for later independent practice. Continue with S04's self-initiated linting decision.

## S04 Self-initiated schema linting

**Evidence state:** First complete coach-authored initiative draft prepared from learner-supplied facts. Episode within S01; not rehearsed or independently retold. Detailed rollout execution remains a possible follow-up.

**Facts supplied by the learner:**

- Observed that APIs were fragmented and sometimes structured very differently.
- While researching tools, found a Google repository containing a public API linter, which served as inspiration. The exact repository has not been identified; no reuse of Google's code or rules is claimed.
- Written guidance on API structure was being introduced at the same time, providing a further reason to turn the guidance into lint rules.
- Initiated and implemented schema linting without being prompted by management.
- Built the linting tool from the ground up using rules taken from an architect's API guidance documents, with the aim of standardizing company schemas.
- Initially worked on the linter during available time between PR review cycles and during a hackathon.
- Discussed introducing the tool with the manager, linking it to architecture's broader standardization push through new API documents and the governance process.
- Identified a gap: services could all have schemas and still expose APIs with very different structures. The case for linting was to address that consistency gap by automating the written guidelines.
- This was a rationale based on an observed problem and organizational priorities; no measured pilot result or quantified benefit before the discussion has been reported.
- The learner and manager agreed to plan the tasks, timelines and how to introduce the tooling to repositories. The specific schedule, rollout mechanism and changes to other commitments have not been described.
- Previously reported that the architecture team valued the linter's contribution to standardization; no quantitative impact or company-wide coverage has been supplied.

**First complete draft — coach-authored, not rehearsed:**

> During our schema-first initiative, I noticed that our APIs were structured inconsistently. Giving every API a schema would still leave that design inconsistency. Architecture was introducing written API guidelines, and I saw an opportunity to turn them into automated checks.

> I found inspiration in Google's public API linter and started working on our own tool between PR review cycles and during a hackathon. I then discussed introducing it with my manager, connecting the proposal to the broader standardization effort. We agreed to plan the tasks, timelines and how to introduce it to our repositories.

> I built the linter using the architect's guidelines. The architecture team valued how it helped standardize our APIs. The idea and initial work came from me, and I involved my manager in planning its introduction.

**Open facts for follow-up:** Initial scope and alternatives; specific rules; effort and rollout execution; any changes to other priorities; and reflection. The agreement to plan is confirmed; exact schedules, dedicated capacity and company-wide coverage are not claimed.

**Coaching observation:** The episode now includes an observed gap, external research, self-initiated work, alignment with an organizational objective, a manager discussion leading to an agreement to plan introduction, and previously reported implementation and architecture feedback. The draft supports discussion of initiative while preserving the difference between initiating work and consulting about introduction.

**Last answered question:** What did you and your manager agree to do next? They agreed to plan tasks, timelines and repository introduction.

**Next step:** Keep the draft for later independent practice. Broaden preparation to a critical-feedback episode, since S01, S03 and S04 all come from the same project.

## S06 Feedback about pressing arguments too hard

**Evidence state:** Feedback and self-reported ongoing behaviour change collected. Revised preliminary wording below; original event context and initial reaction remain unspecified. No complete contextualized answer or rehearsal yet.

**Facts supplied by the learner:**

- During a discussion, felt certain that their position was objectively right and pushed their arguments too hard. The other person was not persuaded. The subject and technical correctness of the position have not been established in this coaching conversation.
- Received feedback that even when right, continuing to push is not always productive; they should seek common ground with the other person.
- Described asking about the other person's issues and listening as an approach that could reveal where the proposed solution addressed those issues.
- Reports generally starting with a strong opinion, while trying to recognize when to stop pushing the discussion and try a different approach.
- Clarified that they applied this feedback in all subsequent discussions: they regard it as an ongoing change in how they work, not an isolated application.
- No specific later event, changed recommendation or resulting outcome has yet been described. This does not mean the learner cannot recall one or that the change was only an intention. The feedback giver and initial reaction remain unspecified.

**Proposed wording — feedback and general response only, coach-authored:**

> I received feedback that I was pushing my arguments too hard in a disagreement. I was convinced my technical position was sound, and I kept pressing it without making progress. The feedback was to find common ground and understand the other person's concerns. Since then, I pay attention to when pushing my argument stops being useful, ask about the issues they are facing, and look for common ground. That has become part of how I handle discussions.

**Coaching observation:** The learner clarified that the feedback became a general practice across later discussions. A request for one illustrative event was confusing; acknowledge the reported ongoing behaviour and explain the purpose of examples without making a specific later event a prerequisite to continuing. Do not infer that the learner lacks examples or has not changed. Preserve the distinction between self-reported change and an independently assessed outcome.

**Open facts:** Initial reaction to the feedback and original disagreement context; optional later example and outcome for interview follow-up. The reported ongoing change is established as the learner's account. Do not claim every later discussion succeeded or that the learner always starts by listening.

**Deferred question:** How did you react when you first received that feedback? Not answered; the learner requested a review of progress and explicit coaching feedback instead.

**Next step:** Use the assessment above to sharpen and practise existing answers. Revisit original context or initial reaction only when useful to the learner; no question is currently awaiting an answer. Do not invent an event or imply that the reported behaviour change was merely intended.

## S09 AI knowledge initiative

**Evidence state:** Complete first coach-authored draft; not rehearsed. Concrete workflow use and design evolution supplied; adoption counts and comparative evaluation methodology unestablished.

**Facts supplied:**

- Owned setup/promotion of an internal skills marketplace, plus knowledge-base design, implementation, design improvements and increasing usage.
- Local and remote agents had to discover where information lived before searching across Confluence, GitHub, OpsLevel and other sources.
- Ingestion skill accepted URLs/raw documents, structured content, added metadata, assessed whether information was new and worth adding, selected/created a domain in a folder tree, and regenerated index files.
- Retrieval used indexes to understand the tree and information locations, searched relevant subtrees and explicitly reported information it could not find. This does not establish perfect retrieval or hallucination prevention.
- Considered vector database versus plain document storage. Chose documents for maintainability, ability to manage information quality and easier Claude-plugin distribution. Reported similar retrieval times in their comparison; no benchmark conditions, accuracy comparison or universal superiority established.
- Added source-specific ingestion skills because generic ingestion worked poorly for some sources; reduced user questions because the workflow was overwhelming; limited retrieval turns after observing long-tail latency. No quantified improvement or proof of no quality loss supplied.
- Presented the solution multiple times and personally implemented an API design review workflow to demonstrate use. That workflow could search API docs; internal PR analyzer could retrieve relevant guidelines.
- Would explore team-specific knowledge bases because teams requested narrower relevance than the general knowledge base provided.

**First draft — coach-authored:**

> I owned the design and implementation of a knowledge base for our internal AI workflows, alongside setting up and promoting an internal skills marketplace. Previously, agents had to search across sources such as Confluence, GitHub and OpsLevel. Before finding an answer, they first had to work out where the relevant information lived.
>
> We built an ingestion flow that structured documents, added metadata, checked for new information and organized it into a document tree with indexes. Agents used those indexes to locate relevant areas, search them and report what they could not find.
>
> A key decision was choosing document storage over a vector database. In our comparison, retrieval times were similar, while documents were easier for us to maintain, manage for quality and distribute as a Claude plugin.
>
> I continued improving the solution based on use. We added source-specific ingestion, reduced the number of questions asked during ingestion and bounded retrieval turns after seeing long-tail latency. To encourage adoption, I presented the solution and built an API design review workflow demonstrating it. That workflow could retrieve API documentation, and our internal PR analyzer could retrieve relevant review guidelines.
>
> Next, I would explore team-specific knowledge bases. The general solution was useful, but teams wanted knowledge with a narrower scope of relevance.

**Assessment:** Recommended default second project for the manager conversation: clear user problem, end-to-end ownership, explicit alternative, concrete iteration, use cases and reflection. S02 remains valuable for platform architecture and security collaboration. This is a coaching recommendation, not an official L5 rating.

**Optional rehearsal probes:** Basis of retrieval-time comparison and quality evaluation; freshness/source traceability; tradeoff from limiting turns; evidence of usage beyond supplied workflows. Do not claim document storage is generally superior to vector retrieval or that bounded turns preserved quality without evidence.

## S05 Quota endpoint exceeded upstream rate limits

**Evidence state:** Learner-supplied incident; first coach draft, not rehearsed.

**Facts supplied:**

- Designed and implemented an endpoint tracking individual users' Claude Code quota usage. Its logic queried an external endpoint for user and quota information.
- Knew the upstream had rate limits and calculated expected traffic would be well below them; omitted caching/memoization on that basis.
- Months later, another tool began depending on the endpoint without the learner's knowledge, increasing traffic approximately 100-fold according to the learner.
- Upstream rate limits were hit and the endpoint began failing. Failure scope and duration were not supplied.
- The team quickly introduced in-memory memoization as a workaround.
- The team also designed a solution to cache the information internally. Implementation of that longer-term solution has not been confirmed.
- Helped debug the incident and provided a mitigation option. A colleague implemented the actual change, which the learner then approved. Do not claim the learner implemented the fix.
- Agreed with the coach-proposed reflection: expected-traffic estimates alone do not address changing demand on a rate-limited dependency. Future design would consider upstream-call reduction, quota-data freshness, behavior at the limit and monitoring. These are intended practices; no later application was supplied.

**Coaching assessment:** Relevant mistake/setback with a clear personal design choice, original rationale, changed conditions and real operational consequences. Do not treat lack of caching as inherently mistaken or infer that the original traffic calculation was wrong. A defensible reflection is that the design relied on expected traffic without enough protection when demand changed, if the learner agrees. Mention the unannounced consumer as context, not the main explanation of responsibility. Do not promise that caching eliminates every future overload risk.

**First draft — coach-authored from known facts:**

> I designed an endpoint for tracking individual users' Claude Code quota usage. It queried an external service with rate limits. I estimated that expected traffic was comfortably below those limits and chose not to add caching.
>
> A few months later, another tool started using the endpoint, increasing traffic roughly a hundredfold. We hit the upstream rate limits and our endpoint began failing. My design had relied on the original traffic assumptions, and it did not handle that increase successfully.
>
> I helped debug the issue and provided a mitigation option. A colleague implemented the change, which I then approved. The team introduced in-memory memoization as the immediate workaround and designed a longer-term solution to cache the information internally, reducing dependence on upstream calls as traffic increased.

**Personal-response assessment:** Clear, bounded contribution to diagnosis, mitigation and approval; implementation credited to colleague. This is sufficient ownership detail for the first draft.

**Agreed reflection — coach-proposed, accepted by learner:**

> My takeaway was that estimating expected traffic is not enough when depending on a rate-limited service. My calculation reflected the consumers we had at the time, but the endpoint could gain new consumers and traffic could change significantly. Next time, I would explicitly consider how to reduce upstream calls, how fresh the quota information needs to be, how the endpoint behaves when the upstream limit is reached, and monitoring that helps spot pressure before requests fail. The lesson was to make changing demand and dependency limits part of the design rather than relying only on the initial traffic estimate.

**Readiness:** Complete first draft combining the incident, personal response and agreed reflection above; ready for rehearsal, not independently retold. Do not present intended practices as already applied. Long-term caching was designed; delivery remains unconfirmed.

**Coverage update:** Both AI projects S02 and S09 now have first drafts after the requested return pass.

## Record for each actual event

The coach creates one record per event and keeps facts separate from wording. An outline becomes a rehearsed story only after an actual practice exchange.

- **ID and title:** A short, memorable label.
- **Evidence state:** Facts collected, outline drafted, rehearsed, or needs clarification.
- **User and problem:** Who was affected, the previous workflow, and why it mattered.
- **Context:** Approximate period, relevant constraints and project scope. Use role labels for colleagues unless names are necessary.
- **Personal responsibility:** What you owned, what was already decided, who else contributed, and your authority to act.
- **Decision and alternatives:** Options, available evidence, uncertainty, tradeoff and strongest opposing argument.
- **Actions:** What you implemented, investigated, communicated or coordinated, in order.
- **People:** Their concerns, how you learned about them, and how their input affected the outcome.
- **Result and support:** Baseline, observation period, population, measure or qualitative evidence, other causes and remaining limitations.
- **Reflection:** What you would do differently; a later application only if it actually happened.
- **Follow-up depth:** Architecture, failure modes, tests, rollout, support and decommissioning facts the manager could ask about.
- **Open facts:** Details not yet remembered or verified. Keep unknown numbers unknown.
- **Answer outline:** A short opening and a few points tailored to the question. Preserve the underlying facts across variants.
- **Practice observations:** The question asked, what came through clearly, what needed prompting, and one next improvement.

## Introduction and motivation notes

To develop in conversation:

- Current scope and the recurring kind of problem you solve.
- Two achievements most relevant to this role.
- Accurate reasons for important career moves and for considering this move.
- Why improving other engineers' work is appealing, with an actual example.
- What you could contribute to the confirmed team and what you would need to learn.
- A culture idea that connects with your experience, plus a real question or concern.
- Your strongest production technologies, adjacent exposure and genuine gaps.

## Questions for the manager

Select from the research chapter, then record the actual answer after the meeting if you want to retain it. First clarify the team behind the broad posting. Other useful areas are expected L5 scope, users and their problems, success measures, decision boundaries, feedback, and operational responsibility.

No manager answers are recorded yet.

## Coaching and persistence rules

Preserve uncertainty, personal scope, negative outcomes and other people's credit. Do not manufacture metrics, customers, conflicts, escalations or later improvements. A fictional example in the chapter must never be promoted into a personal story. Distinguish an assistant's proposed wording from a fact supplied by the learner.

Keep necessary preparation notes rather than verbatim transcripts or confidential implementation detail. With repository write access, read back saved changes. With a read-only GitHub connection, keep the working outline in conversation and state that it was not persisted. Updating this template or the research chapter does not count as learning progress.
