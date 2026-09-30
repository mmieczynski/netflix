# Chapter 5 — Stacks greedy choices and dynamic programming

When a problem does not reduce to lookup or ordering, ask what choices remain and whether different paths reach the same remaining problem.

## Stacks remember unfinished work

Balanced delimiters are the simplest example. Push opening brackets. A closing bracket must match the most recent unfinished opening bracket, so inspect the top. If it does not match, reject immediately. Finish with an empty stack. Counts alone fail on “open square, open round, close square, close round” because totals ignore nesting order.

A monotonic stack is useful for “next greater” questions. Suppose daily loads are seventy, seventy three, seventy one, seventy four. Store indices whose next greater load is unresolved, keeping their values decreasing from bottom to top. When seventy four arrives, it resolves seventy one and seventy three; pop each and compute its distance. Each index is pushed once and popped once, so total work is linear even with a nested loop.

The proof is about discarded candidates. A value popped by a later greater value no longer needs to wait: this is the first greater value encountered since its insertion. For “next greater or equal,” change the equality rule. For a sliding-window maximum, use a monotonic deque because old indices must also expire from the front. A regular stack cannot efficiently remove those oldest elements.

## Greedy requires a reason that local choices are safe

Choose the maximum number of nonoverlapping screenings. Sorting by earliest finish and taking the next compatible screening works: replace the first screening in any optimal schedule with the earliest-finishing one. It cannot finish later, so all later choices remain feasible. Repeat the exchange argument on the remaining schedule.

Sorting by earliest start does not work: a long early interval may exclude many shorter intervals. If screenings have different values and you want maximum total value, earliest finish is no longer enough. The objective changed, so the proof no longer applies. That variant naturally introduces dynamic programming.

## Backtracking explores choices deliberately

Generate playlists of exactly three distinct titles from a small set, subject to a duration budget. A state includes chosen titles, next available index, and remaining budget. Choose a title, recurse, then undo the choice. Choosing indices in increasing order avoids generating permutations when order does not matter. If order matters, that restriction would wrongly remove valid playlists.

Prune only when a branch cannot recover. Exceeding a duration budget is safe to prune if all future durations are nonnegative. It is not safe if negative adjustments exist. Copy the current slice when saving a completed answer; otherwise later mutations may overwrite earlier answers through the same backing array. The search can remain exponential, and the output itself may be exponential. Do not hide that cost behind “recursion.”

**Think aloud.** Two recursive paths have selected different earlier titles but have the same next index and remaining budget. When may they reuse an answer? What if the question asks you to list the actual playlists?

#### Worked answer — Coach notes

For counting or best achievable future score, earlier history may be irrelevant once the state summarizes every future constraint. For enumerating full playlists, you can reuse suffix results carefully, but outputs must still include each distinct prefix. If a genre restriction depends on earlier choices, index and budget are not a sufficient state.

## Dynamic programming is reused subproblem reasoning

Consider choosing nonadjacent episode promotions to maximize total value. At index i, either skip it and solve from i plus one, or take it and solve from i plus two. Define best of i as the best total achievable from index i onward. The recurrence is the larger of those two choices. Beyond the end, the best value is zero. If selecting nothing is allowed, negative values naturally get skipped.

For values four, one, one, four, the answer is eight from the first and last positions. Recursive branching repeats the same suffix problems. Memoization reduces the number of solved states to n. A bottom-up pass from right to left uses only the next two answers, giving linear time and constant auxiliary space. If you must return the selected indices, store decisions or a full table and reconstruct. Reference: `MaxNonAdjacent`.

Before writing a DP table, answer four questions in sentences: what does one state mean; what choices leave it; what smaller states do those choices require; and what are the base cases? Then count states times work per state. For a duration budget B and n titles, a table indexed by item and budget may cost O(nB), which is pseudopolynomial because B's numeric magnitude matters.

## How to discover a state instead of guessing a recurrence

Start with a choice you can explain in everyday language. Suppose promotional slots lie in a row and adjacent slots cannot both be selected. At the first slot, you either take it or skip it. Taking it removes the next slot from consideration. Skipping it leaves the next slot available. Both branches lead to smaller problems of the same kind. You have found a recurrence by describing valid choices, not by looking for a formula.

Now ask what information a smaller problem needs. If all earlier constraints have already been resolved, the suffix starting at a position is enough. It does not need to remember the exact sequence of earlier decisions. Two branches reaching that same suffix can reuse the same best result. If there is also a remaining money budget, position alone is no longer enough; the state needs budget too. A dynamic-programming state is a summary of all past information that can still affect the future.

There are two different opportunities to avoid work. A greedy proof shows that one choice can always replace another without harming an optimal answer, so you never explore the rejected alternative. Dynamic programming explores the meaningful alternatives but shares repeated subproblems. Backtracking explores a decision tree and may prune branches that cannot lead to valid answers. These techniques can coexist, but they eliminate work for different reasons.

The safest way to test a proposed state is to try to make it lie. Can two histories produce the same state but have different legal futures or different optimal future values? If so, the state has forgotten something important. For example, a playlist state containing only remaining duration cannot enforce “do not choose two films from the same director” unless it also remembers selected directors or uses another representation that preserves the restriction.

Similarly, test a greedy choice by trying to defeat it with a tiny example. Choosing the highest-value promotion first sounds plausible. With values six, ten, six in three adjacent positions, taking the middle ten blocks both sixes, while taking the ends gives twelve. The counterexample shows that this greedy rule lacks the exchange property you would need. It does not mean greedy algorithms are unreliable; it means each greedy rule needs its own proof.

## Questions and worked explanations on choices and reuse

**Round 1.** Why do balanced brackets need a stack rather than just a count of opening and closing symbols? Describe a sequence with matching totals but invalid nesting.

#### Worked answer — Explanation and next challenge

Opening square, opening round, closing square, closing round has balanced totals but crosses the nesting order. The next closer must match the most recent unfinished opener. A stack preserves that relationship. A count preserves quantity while losing order, just as a set preserves membership while losing multiplicity.

**Round 2.** Loads arrive as five, three, six. For each earlier load, find its first later greater load. What does six resolve when it arrives?

#### Worked answer — Explanation and next challenge

Six resolves both three and five. Before six arrives, their next greater values are unknown. A decreasing stack holds these unresolved candidates. Each is removed the first time a sufficiently large future value appears. The order of removal is three then five, but their recorded answer is the same new position.

**Round 3.** A monotonic-stack implementation pops several items in one iteration. Why can its total work still be linear?

#### Worked answer — Explanation and next challenge

Each item is pushed once and popped at most once. Charge each removal to the item removed instead of charging a full scan to every incoming item. This accounting is closely related to sliding-window analysis. Being able to reuse the proof is more useful than memorizing separate complexity claims for each pattern.

**Round 4.** To maximize the number of nonoverlapping screenings, why might finishing earlier help more than starting earlier?

#### Worked answer — Explanation and next challenge

An earlier finish leaves at least as much remaining time for future screenings. Replacing the first selected interval of an optimal schedule with an earliest-finishing compatible one cannot invalidate later choices. A very early start could belong to an extremely long screening. The objective is count, so occupying less future time is what matters.

**Round 5.** Now screenings have different rewards. Is the earliest-finish argument sufficient to maximize total reward? Explain the missing guarantee.

#### Worked answer — Explanation and next challenge

Replacing one interval with an earlier-finishing interval might lose a large reward. Preserving feasibility is not enough; the replacement must also preserve or improve the objective. Weighted interval scheduling needs a different argument, often a recurrence comparing taking an interval with skipping it and using the latest compatible predecessor.

**Round 6.** Nonadjacent values are six, ten, six. What does “always take the largest remaining value” choose, and what better answer disproves it?

#### Worked answer — Explanation and next challenge

The greedy rule takes ten and blocks both neighbors. Taking both sixes gives twelve. This tiny example is a constructive reason to keep alternative decisions. Ask next what smaller problem remains after taking or skipping the first value.

**Round 7.** Define “best from position i onward” in words. What are the two alternatives at i, and why does taking i jump past the next position?

#### Worked answer — Explanation and next challenge

It is the maximum achievable sum using only positions at or after i under the nonadjacency rule. Skip i and use the best suffix from the next position, or take its value and add the best suffix from two positions later. The jump enforces the constraint directly. The state definition must include whether an empty selection is allowed.

**Round 8.** Two different decision paths reach the same remaining suffix with no other constraints. Why can they share the same best future value?

#### Worked answer — Explanation and next challenge

The legal future choices and their rewards depend only on that suffix. The histories may contribute different already-earned totals, but the additional optimum is the same. Memoization stores the future contribution under the suffix position. If previous choices affect a remaining constraint, include that constraint in the state before sharing results.

**Round 9.** All values are negative. What answer follows if choosing nothing is allowed? What changes if at least one item must be selected?

#### Worked answer — Explanation and next challenge

With empty selection allowed, zero beats any negative total. With at least one required, zero is not a valid answer unless a selected combination actually achieves it. The recurrence or state must distinguish a valid nonempty selection from “nothing chosen.” Base cases encode the contract, not just termination.

**Round 10.** You generate combinations of titles and save the current slice when a combination is complete. Why might later backtracking alter previously saved answers?

#### Worked answer — Explanation and next challenge

If the saved slices share the current path's backing array, later overwrites can modify their elements. Save a copy of the chosen elements. This connects algorithmic correctness to Go ownership semantics: the mathematical answer should be a snapshot, so its representation must preserve that snapshot.

**Round 11.** You prune a playlist branch whenever its duration exceeds the budget. Under what assumption is that safe, and what kind of future value would invalidate it?

#### Worked answer — Explanation and next challenge

It is safe when every additional duration is nonnegative, so extending can never return under budget. Negative adjustments could reduce the total and make an apparently invalid partial branch lead to a valid answer. A pruning rule is a small impossibility proof. State that proof each time you prune.

**Round 12.** A playlist must fit a budget and use each director at most once. Is “next title index plus remaining budget” a sufficient memoization key? Construct the reason verbally.

#### Worked answer — Explanation and mastery check

No, two histories can have equal position and remaining budget but different already-used directors. A future title may be legal in one history and illegal in the other. Add the relevant director information or find a different formulation. This test of state sufficiency is the central skill; the resulting state space might be too large, which itself informs whether exact DP is practical.

Explain the difference between proving a choice safe, sharing a repeated subproblem, and proving a branch impossible. Give one small example of each without writing a recurrence or code.

## Worked program design: choosing nonadjacent promotions

The function takes a slice of signed rewards and returns the greatest total from nonadjacent positions. Empty selection is allowed. The result uses a numeric type large enough for allowed sums. An empty input therefore returns zero. We first solve for the value alone; returning selected positions is a separate extension.

Define best-from-i as the best achievable total using positions i onward. At i, skipping gives best-from-next. Taking gives the current reward plus best-from-two-ahead. Choose the larger. Beyond the input, the contribution is zero. This description is already executable recursive pseudocode, but a naive recursion repeats many identical suffix problems.

A bottom-up implementation walks from the last position toward the first. Keep two fields of working state: the best result for the next suffix and for the suffix two positions ahead. Compute the new best using their old values, then shift the working state so it refers to the correct suffixes for the next iteration. Updating a variable too early can accidentally replace a needed old result.

For four, one, one, four, begin beyond the end with zero and zero. The final four produces best four. The preceding one produces best four. The next one produces best five. The first four produces best eight by combining itself with the best suffix beginning at the third position. The returned eight corresponds to selecting the two ends.

The loop uses linear time and constant auxiliary storage. To return the chosen positions, keep the table of suffix answers or equivalent decisions, then walk forward: compare taking with skipping under a defined tie policy. Taking advances two positions; skipping advances one. This reconstruction requires information discarded by the two-number optimization, so the extra output requirement affects space.

**Optional spoken walkthrough:** narrate the two alternatives at each position of six, ten, six. Explain why the middle value loses to both ends, then state what must change if selecting nothing is forbidden.
