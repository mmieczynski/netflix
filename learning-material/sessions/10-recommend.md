# Chapter 10 — Recommendation and streaming ranking tasks

For a coding interview, turn “recommend movies” into a small, testable pipeline with explicit inputs, filters, scores, and tie rules.

## Understand the domain without overbuilding it

Netflix describes recommendations as using member interactions, similar tastes, and title information, with personalization affecting rows and ordering. This gives useful domain context, but it is not an interview specification. [Netflix's recommendation overview](https://help.netflix.com/en/node/100639).

A general recommendation pipeline first retrieves plausible candidates, scores them, and applies final constraints. Content-based retrieval uses item attributes; collaborative methods learn from patterns across users and items. This is a conceptual foundation, not a requirement to train a model in Go. [Google's candidate generation overview](https://developers.google.com/machine-learning/recommendation/overview/candidate-generation).

For an interview, ask what is already given. If scores are supplied, your task may be filtering and top k selection. If co-watch histories are given, it may be graph or frequency aggregation. If scores update live, it may be an indexed heap or an ordered-state problem. The title “recommendation engine” does not choose the algorithm.

## A precise ranking exercise

Input is a list of candidates with ID, genre, integer score, and eligibility flag, plus a set of watched IDs and k. Exclude ineligible and watched candidates. Deduplicate remaining IDs by retaining their highest score. For equal scores of the same ID, retain the lexicographically smaller genre so metadata is deterministic. Return at most k titles ordered by descending score, then ascending ID. Nonpositive k returns an empty result.

For candidates A score eight, B score nine, A score ten, C score nine, with B watched and k two, the winners are A then C. Deduplication must happen before top k selection, or two copies of A can occupy both slots. Eligibility is applied before deduplication in this contract; if eligibility instead belongs to the canonical catalog record, resolve canonical metadata first. That is a requirements distinction worth raising.

The simple baseline builds a map of unique eligible candidates and sorts its values. For n input records and u unique eligible IDs, expected work is O(n plus u log u), and space is O(u). When k is much smaller than u, use the bounded heap from lesson 3, then sort the k winners. The deduplication map still uses O(u) space even though the heap uses only O(k). Do not claim the whole pipeline uses O(k) memory.

**Think aloud.** The interviewer changes the rule to “at most one movie per genre.” Can you keep only the original top k and then discard repeated genres?

#### Worked answer — Coach notes

No. The original top k may all belong to one genre, while the best allowed candidates from other genres lie outside it. Apply the constraint during selection over sufficient candidates, or select each genre's best first under the one-per-genre contract.

## Constraints change what top k means

For a single disjoint genre per movie and a cap of one per genre, keeping the best movie of each genre and then selecting the best k genre winners maximizes total score for that count. Here we fill as many slots as possible up to k, even if scores are negative. If the only objective is maximum total score and fewer selections are allowed, skip negative scores. For a cap greater than one, sorted greedy acceptance respects the per-genre quota. If movies have multiple overlapping tags and several simultaneous quotas, a simple greedy strategy need not be globally optimal. State whether you are maximizing score exactly or applying a reasonable reranking heuristic.

Another common follow-up asks to interleave two ranked lists while removing duplicates. Two pointers with a seen set can preserve a chosen alternation policy. This does not necessarily preserve global score order. Define whether fairness between sources or highest overall scores has priority. If each source is sorted and global score order is required, a heap across source heads gives a k-way merge.

## Small collaborative scoring without machine learning

Suppose a target user watched A and B. Neighbor one watched A, C, D; neighbor two watched B, C; neighbor three watched E. Define similarity as number of shared watched titles. Neighbor one and two each have weight one; neighbor three has zero. Add a neighbor's weight to each of their unwatched titles. C receives two and D one. Recommend C before D. Use sets per history so repeat viewing does not accidentally count as multiple distinct overlaps.

This is an original toy scoring rule. It illustrates joins and aggregation, not Netflix's actual ranking formula. An inverted index from title to users can avoid scanning every user to find overlaps, although extremely popular titles can still generate many candidates. Define cold-start behavior when no overlap exists, such as a supplied popularity list filtered by eligibility.

## Streaming statistics add reversibility

For average watch duration per video, store sum and count, then divide only when count is positive. For a one-hour window, retain events in a time queue and subtract each expired event's duration and count from its video's aggregate. A running mean by itself does not retain enough information to remove arbitrary old contributions easily. If ordering is not monotone, a FIFO expiration queue is no longer sufficient.

Dynamic top k is harder than one-time top k. Counts can fall when events expire. A size-k heap of only current winners may have forgotten the best replacement outsider. Keep full aggregates and rebuild on query, or maintain an indexed structure over all candidates; discuss the update-versus-query trade-off. An approximate heavy-hitter sketch is a different contract with error and ranking uncertainty.

## Make the recommendation problem small enough to reason about

A recommendation prompt can sound enormous because the real product contains retrieval models, ranking models, experiments, data pipelines, and policy rules. A coding exercise usually exposes only a small piece of that system. Your first task is to identify the piece. Are candidate scores already supplied? Are you deriving a toy score from viewing histories? Are you maintaining changing scores over time? These versions require different programs even though they all return movie IDs.

Separate the stages in words. Eligibility decides whether a candidate is allowed at all. Identity resolution ensures one logical title is not counted twice. Scoring assigns a comparable value under an agreed formula. Selection chooses a limited set. Ordering determines the sequence in which the chosen set is returned. These stages can interact: discarding duplicates after selecting winners may leave too few distinct results, and filtering only after selection can hide valid candidates just below the original cutoff.

Think of a supplied score as an input contract, not an explanation of user enjoyment. Your coding task might require ranking by that number without evaluating whether the number is a good model. A toy collaborative rule can be educational, but it should be described honestly as a rule we define for the exercise. It is not evidence of Netflix's actual model architecture or training objective.

Determinism is part of the API. Two candidates with equal scores need an agreed order if callers or tests compare results. A map's iteration order should not accidentally decide which result wins. Ties become particularly subtle at a heap boundary because the heap keeps the worst retained winner at its root. You must reverse both the score direction and the tie direction consistently.

Temporal ranking reveals why a static top-k proof has limits. While processing a fixed set, once k better candidates dominate an outsider, discarding the outsider is safe. Later expiration can reduce the winners' scores. The outsider may then deserve promotion, but a winner-only structure has forgotten it. An algorithm can be correct for static selection and insufficient for dynamic maintenance without any contradiction: the future operations changed what information must be retained.

## Questions and worked explanations on ranking pipelines

**Round 1.** An interviewer says “Recommend five movies” and provides candidates with scores. What would you clarify before discussing machine learning?

#### Worked answer — Explanation and next challenge

Clarify eligibility, watched-title filtering, duplicate IDs, score ordering, tie behavior, and whether scores are already authoritative. If the task is to select from supplied scores, building a learning model is outside the contract. Establish the exact input-to-output transformation first.

**Round 2.** Candidates are A with score eight, B with nine, A with ten, and C with nine. B is watched. Return two distinct titles using each ID's best score. Explain the stages.

#### Worked answer — Explanation and next challenge

Remove watched B, keep A's score-ten record instead of its score-eight duplicate, and rank A before C. The result is A, C. Selecting two records before deduplicating could retain two As, leaving a result with only one distinct title. Stage order must match the promised result.

**Round 3.** Two records for the same title disagree on eligibility. Why is “take the record with the highest score” not enough to resolve the ambiguity?

#### Worked answer — Explanation and next challenge

Eligibility may belong to a canonical catalog record rather than to each recommendation-source record. If the task instead defines filtering per input record, filtering before deduplication is a defensible explicit policy. Ask which source is authoritative. A score comparison cannot settle conflicting metadata semantics.

**Round 4.** A and C have equal scores and lower ID wins. What should the output order be, and why is map iteration not an acceptable tie breaker?

#### Worked answer — Explanation and next challenge

A precedes C. Map iteration is not a specified ranking and can vary, making behavior nondeterministic. A complete comparison rule supplies a stable result independent of storage order. The same rule must be used by full sorting and by the final sort of heap-selected winners.

**Round 5.** You use a size-five heap but also keep a map of every distinct candidate to deduplicate. Can you claim total extra space is just five records?

#### Worked answer — Explanation and next challenge

No. The heap is bounded by five, but the map grows with the number of unique eligible IDs. Complexity counts every retained structure. If memory must be bounded independently of unique IDs, the input contract or algorithm needs a different deduplication strategy or an explicitly approximate guarantee.

**Round 6.** The top three scores all belong to one genre, but the result may contain at most one per genre. Why is selecting top three first and then dropping repeated genres insufficient?

#### Worked answer — Explanation and next challenge

You may return only one title even though good candidates from other genres exist below the original cutoff. Selection must consider the constraint while it still has access to enough candidates. With exactly one genre per title and a cap of one, choosing the best representative of each genre before selecting across genres is sufficient for maximum total score at the required result count.

**Round 7.** Now titles can have several overlapping tags, with limits across tags. Why should you be cautious about claiming a simple greedy ranking is globally optimal?

#### Worked answer — Explanation and next challenge

A chosen title can consume capacity in multiple constraints and block a combination whose total score is better. The disjoint-genre argument no longer applies. A heuristic may be acceptable if the prompt permits it, but exact optimality needs a proof or another algorithm. Be precise about the goal rather than promising every desirable property.

**Round 8.** The target watched A and B. One neighbor watched A and C; another watched B and C and D. Weight each neighbor by number of shared distinct titles. What scores do C and D receive?

#### Worked answer — Explanation and next challenge

Each neighbor has weight one. C receives one from each and totals two; D receives one. A and B are excluded as already watched. Repeated occurrences of A in a history should not increase shared-distinct-title overlap unless that is intentionally part of the scoring rule.

**Round 9.** No neighbor overlaps a new user's history. What should the program return, and why is that a contract question rather than a data-structure failure?

#### Worked answer — Explanation and next challenge

The contract might specify an empty result or an eligible popularity fallback. The algorithm can correctly find no collaborative evidence. A fallback supplies a separate product rule. It should still apply watched and eligibility filters and use deterministic ordering, rather than bypassing constraints because personalization is unavailable.

**Round 10.** For a video's recent average watch duration, why keep sum and count instead of only the current average when old events must expire?

#### Worked answer — Explanation and next challenge

To remove an event, subtract its duration from the sum and decrement the count. The average alone does not identify how much total duration it summarizes or how many observations contributed. Retaining the expiring event contributions and aggregate sufficient statistics makes the update reversible.

**Round 11.** A is the only retained top-one winner. B was discarded as second best. A's old events expire and its score falls below B's. What information is missing from a winner-only heap?

#### Worked answer — Explanation and next challenge

It no longer knows B's current score or even that B is the best replacement. Keep full aggregates and sort or select on query, or maintain a suitable structure over all active candidates. Static selection's discard proof does not survive arbitrary score decreases.

**Round 12.** Build a verbal pipeline for “top unwatched videos by recent total watch duration.” Include duplicate events, time boundaries, aggregation, selection, ties, and one late-event question.

#### Worked answer — Explanation and mastery check

Define event identity and duplicate policy; select events in the agreed time interval; aggregate durations by eligible unwatched title; choose top k with deterministic ties. For streaming updates, retain contributions to remove on expiry. Clarify whether late event-time arrivals are accepted and how they enter ordered expiration state. This task integrates several patterns without requiring a new memorized algorithm.

Have the tutor replace movies with ads, search results, or products. Explain which pipeline invariants remain identical and which eligibility or diversity rules are domain-specific.

## Worked program design: ranking eligible distinct titles

The function accepts candidate records, a watched-ID set, and k. Each candidate has ID, genre, integer score, and eligibility. It returns up to k distinct eligible unwatched titles, highest score first and smallest ID first on ties. For duplicate IDs, retain the highest-scoring eligible record; equal-score duplicates use the smaller genre string for determinism. Nonpositive k returns an empty result.

First scan candidates. Skip ineligible or watched records. Look up each remaining ID in a deduplication map. Keep the incoming record only if no record exists or it wins the duplicate-resolution rule. This stage establishes one canonical candidate per eligible ID under the exercise's policy. It must happen before selecting final winners.

For the baseline, collect the map values and sort using the final comparison. Return at most the first k. For a small k among many unique candidates, instead maintain a heap whose root is the worst retained winner. Push until it has k entries. Thereafter replace its root only when a newcomer is better, then repair the heap. Sort the retained winners into final output order before returning.

On A score eight, B nine, A ten, C nine, with B watched and k two, the canonical map contains A ten and C nine, so return A then C. If A and C tie, A still comes first. Test duplicate IDs, no eligible records, k larger than the result, and all equal scores. The comparison rule should give the same result regardless of map iteration order.

Expected scan cost is linear. Sorting all unique candidates costs u log u; a bounded heap costs u log k plus sorting the final k. Both versions still retain the deduplication map, so their total auxiliary memory is not merely k. A one-per-genre rule changes selection: keep each genre's best representative before global selection, provided each title belongs to exactly one genre and the objective is to fill the available slots up to k, then maximize their total score.

**Optional spoken walkthrough:** describe separate helpers for eligibility, duplicate resolution, and final comparison. Then explain why the heap comparison is the reverse of the final ranking, including its tie rule.
