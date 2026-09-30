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

**Failure case.** Two paths can reuse a suffix answer only when their remaining constraints and requested output agree.

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

## Worked example: discover the state from the choices

Three advertising slots offer rewards 6, 10, and 6. Adjacent slots cannot both be chosen. Choosing the largest reward first gives ten, but choosing both ends gives twelve. The failed greedy choice suggests comparing complete alternatives rather than committing to the locally largest value.

![At each position, either skip it or take it and jump past its neighbor. The two branches solve smaller suffixes.](figures/choices.svg)

Let best at position i mean the largest reward obtainable from position i onward. Skip gives best at i plus one. Take gives the current reward plus best at i plus two. Keep the larger. Past the end, the reward is zero. This definition includes the assumption that choosing nothing is allowed.

| Position, evaluated backward | Skip | Take | Best |
| --- | --- | --- | --- |
| Last 6 | 0 | 6 + 0 | 6 |
| Middle 10 | 6 | 10 + 0 | 10 |
| First 6 | 10 | 6 + 6 | 12 |

The suffix state is sufficient because the earlier choices impose no further condition once entry into that suffix is legal. If the problem also limits the number of chosen slots, position alone is insufficient: two visits to the same position with different remaining allowances can have different answers. Add the allowance to the state. Memoization works only when the state actually identifies an equivalent remaining problem.

To recover the selected slots, store which branch won or retain the table and compare alternatives again while walking forward. Computing only the best number can use two rolling values; reconstructing a selection needs additional information. The requested output changes the representation.

## Worked example: a stack postpones decisions

For each value in 2, 1, 3, find the next strictly greater value to its right. A stack holds indices still waiting for an answer. After reading two, index zero waits. After reading one, both zero and one wait, with the smaller value at the top. Reading three answers both unresolved positions: first one, then two. Index two then waits until the input ends and receives no answer.

| New value | Values still waiting | Answers established |
| --- | --- | --- |
| 2 | 2 | None |
| 1 | 2, 1 | None |
| 3 | 3 | Next greater for 1 and 2 is 3 |

Store indices rather than just values so answers can be written to the correct positions, including duplicates. For a strictly greater query, an equal arrival does not resolve an equal waiting value. That single word in the contract changes the comparison.

Every index is pushed once and popped at most once. Several pops during the final arrival do not make the complete scan quadratic. This is the same aggregate accounting used for sliding-window removals, applied to a different invariant: the stack contains precisely the unresolved positions in monotonic value order.

## Putting the program together: choosing nonadjacent promotions

The function takes a slice of signed rewards and returns the greatest total from nonadjacent positions. Empty selection is allowed. The result uses a numeric type large enough for allowed sums. An empty input therefore returns zero. We first solve for the value alone; returning selected positions is a separate extension.

Define best-from-i as the best achievable total using positions i onward. At i, skipping gives best-from-next. Taking gives the current reward plus best-from-two-ahead. Choose the larger. Beyond the input, the contribution is zero. This description is already executable recursive pseudocode, but a naive recursion repeats many identical suffix problems.

A bottom-up implementation walks from the last position toward the first. Keep two fields of working state: the best result for the next suffix and for the suffix two positions ahead. Compute the new best using their old values, then shift the working state so it refers to the correct suffixes for the next iteration. Updating a variable too early can accidentally replace a needed old result.

For four, one, one, four, begin beyond the end with zero and zero. The final four produces best four. The preceding one produces best four. The next one produces best five. The first four produces best eight by combining itself with the best suffix beginning at the third position. The returned eight corresponds to selecting the two ends.

The loop uses linear time and constant auxiliary storage. To return the chosen positions, keep the table of suffix answers or equivalent decisions, then walk forward: compare taking with skipping under a defined tie policy. Taking advances two positions; skipping advances one. This reconstruction requires information discarded by the two-number optimization, so the extra output requirement affects space.
