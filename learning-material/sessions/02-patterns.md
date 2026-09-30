# Chapter 2 — Hashing arrays and sliding windows

A faster solution often remembers precisely the information that the next step needs.

## Two Sum as a lesson in sufficient state

Suppose playback chunks have sizes two, seven, eleven, and fifteen. Find two distinct chunks totaling nine. Trying every pair is correct and quadratic. When considering seven, the only relevant question about earlier values is whether two occurred. Store value to earlier index. Before inserting the current value, look for target minus current value. This order prevents reusing the same element.

The invariant is that before processing position i, the map contains only positions before i. If the complement exists, the indices are distinct and their values sum to the target. If it does not, storing the current index prepares exactly what future positions need. Duplicates such as three and three with target six work because the first three is stored before the second is examined. The expected time is linear and space is linear. Reference: `TwoSum` in algorithms.go.

Transfer this idea to event matching, reconciling IDs, or counting pairs. Ask whether you need existence, a count, or actual positions. A set cannot answer how many matching records exist; a frequency map can. An index chosen for one query may discard information required by another.

## Grouping means inventing a stable identity

To group anagrams, convert each word to a signature. Sorting its letters gives the same signature to words with the same letter multiset. For lowercase English letters, an array of twenty six counts is a comparable Go map key and avoids sorting each word. For a word of length m, sorting costs O(m log m); counting costs O(m) under the fixed alphabet assumption. A signature that stores only which letters appear would wrongly group “abb” and “ab.”

Real-world analogues include normalizing event identities and deduplicating equivalent requests. Decide what equivalence means first. Case folding, Unicode normalization, and ignoring punctuation are requirements, not automatic improvements. Making a key “more normalized” can merge records that should remain distinct.

## Sliding windows maintain a local invariant

Find the longest contiguous sequence of titles with no repeated ID. A brute-force approach starts at every position and extends until a duplicate appears. The overlap between these searches is wasted work. Instead maintain a left boundary and a map of last seen positions while advancing the right boundary once.

Walk through A, B, B, A. After A and B, the window has length two. The second B was at position one, so move left to position two. At the final A, its previous position is zero, already outside the window. Do not move left backward. The update is left equals the larger of current left and previous position plus one. The window now contains B, A, again length two.

The invariant is that the current window contains no repeated ID. Adding the rightmost element can introduce only a repetition of that element. Moving past its last occurrence repairs the invariant with the smallest necessary movement. Both boundaries only move forward, giving linear time. A nested while loop that shrinks the left side can still be linear because each element leaves once. Reference: `LongestUnique`, using rune positions.

**Think aloud.** Why is “move left to last seen plus one” incorrect without checking the current left boundary? Use A, B, B, A to give a counterexample.

#### Worked answer — Coach notes

The final A would move left back to one and reintroduce the duplicate B. A correct boundary only advances. Explain the invariant before quoting a max expression.

## Know when a window is not enough

For positive numbers, shrinking a sum window reduces its sum, so you can use directional reasoning for certain threshold tasks. With negative numbers, that assumption disappears. For example, the array three, minus two, two contains a length-three subarray summing to three even though an early sum may suggest shrinking. Do not use a positive-only argument on signed data.

To count subarrays with sum k for arbitrary integers, maintain prefix sum p. A subarray ending here sums to k whenever an earlier prefix equals p minus k. Store counts of earlier prefixes, initialized with zero occurring once to represent a subarray starting at the beginning. Add the number of matches before incrementing the current prefix count. For one, minus one, one with target one, the answer is three. This is hashing of accumulated state, not a window. Reference: `CountSubarrays`.

## Prefix and suffix work without division

Product Except Self asks for the product of all values except the current one. Dividing a total product fails at zeros and may violate the prompt. First fill each output position with the product strictly to its left. Walk backward with a running product strictly to its right, multiplying it into the output. On two, three, four, left products are one, two, six; right products supply twelve, four, one; the result is twelve, eight, six.

The empty product is one. That identity makes the edges work without special cases. One zero means only its position may have a nonzero output; two zeros make every output zero. The algorithm is linear with constant auxiliary space excluding output, assuming products fit the chosen numeric type. The deeper pattern is combining summaries from both sides of a position.

## Deriving an algorithm by deciding what the past must remember

Imagine reading a long sequence from left to right. You cannot keep re-reading every earlier element if you want a fast solution. Instead ask: what question will the next element ask about the past? In Two Sum, the next value asks whether its complement appeared earlier. In frequency counting, it asks how many times the same value appeared. In longest unique substring, it asks where its last occurrence was. These are three different summaries of history, even though all can be stored in a map.

This question helps you invent the state rather than guess a named pattern. A map is a container, not the whole solution. You must choose the key, the meaning of its value, and when an update occurs. Storing only presence loses frequency. Storing only the latest index loses all earlier positions. Either loss is fine when those earlier details cannot affect any future answer. That last sentence is a correctness argument: discarded information must be irrelevant to the remaining task.

A window adds a boundary to this reasoning. Your summary describes a current contiguous region rather than the whole prefix. Every time you move the left edge, you must remove that departing element's effect. That is why a count map works for an at-most-two-distinct window. A single boolean would not tell you whether another copy remains inside after one copy leaves.

Do not confuse “contiguous” with “sorted.” A viewing session is contiguous because it occupies consecutive events in the original order. Sorting those events changes the session. Two pointers can sometimes work without sorting, but only if the problem supplies another directional property. For a positive sum, moving the right edge increases the sum and moving the left edge decreases it. Signed values remove that property. The method fails because its proof fails, not because a pattern lookup says negative numbers belong somewhere else.

Prefix sums are a different way of representing a contiguous region. Picture cumulative watch duration as an odometer. The duration between two points is the later reading minus the earlier reading. If the current reading is seven and you want a segment totaling three, you are looking for an earlier reading of four. Store how many earlier readings had each value. Negative durations or adjustments do not break subtraction, so this approach survives where directional window reasoning does not.

## Questions and worked explanations on remembering enough

**Round 1.** The values are two, seven, eleven and the target is nine. At seven, what exact question must you answer about the earlier values? Explain the baseline and the improvement.

#### Worked answer — Explanation and next challenge

The question is whether two appeared earlier. A scan searches the earlier values each time; a map records value to index so the question is answered directly. You are replacing repeated search, not merely adding a map because the task resembles Two Sum. Before proceeding, explain why you store an index rather than only a boolean when the output needs positions.

**Round 2.** The input contains just one three and the target is six. Why must you check for the complement before inserting the current value?

#### Worked answer — Explanation and next challenge

Inserting first would let the current value find itself. Checking first guarantees that any match comes from an earlier position. Now add a second three: it correctly finds the first. The distinction is between two equal values at different indices and reusing one observation twice.

**Round 3.** You group words by their distinct letters. Would that correctly group anagrams? Use “ab” and “abb” to explain the information that is missing.

#### Worked answer — Explanation and next challenge

Both have the distinct-letter set containing a and b, but the counts differ. Anagram identity depends on the multiset of letters. A sorted-letter signature or a count vector preserves multiplicity. Generalize: before inventing a signature, state the exact equivalence relation it must preserve.

**Round 4.** In a no-repeated-title window, read A, B, B one event at a time. State the valid window after each event and why the left boundary changes.

#### Worked answer — Explanation and next challenge

After A the window is A. After B it is A, B. The next B conflicts with the earlier B, so move past that earlier occurrence; the window is now just the latest B. Removing only A would not repair the violation. Move the smallest distance that restores validity so useful candidates are not discarded unnecessarily.

**Round 5.** Continue that sequence with A. Its last occurrence is before the current window. Should the left boundary move? Explain the bug if it moves backward.

#### Worked answer — Explanation and next challenge

Leave the boundary where it is and extend to B, A. Moving backward would reintroduce an earlier B and violate uniqueness. The invariant must be true for the current region; historical duplicates outside it are irrelevant. If this is unclear, ask the tutor to repeat the current window separately from the full history.

**Round 6.** Change the rule to at most two distinct titles. The current window is A, A, B and C arrives. Which removals are needed, and why is a count useful?

#### Worked answer — Explanation and next challenge

The first A removal still leaves A, B, C, with three distinct titles. Removing the second A leaves B, C, which is valid. Counts distinguish one remaining A from none. Delete a title from the distinct set only when its count reaches zero. The left edge can move more than once during one right-edge step without making total work quadratic.

**Round 7.** Explain why a loop inside another loop can still perform linear work in this window algorithm. Avoid relying on the phrase “two pointers” as the explanation.

#### Worked answer — Explanation and next challenge

Each element enters once as the right boundary passes it and leaves at most once as the left boundary passes it. Across the whole input there are at most a linear number of additions and removals. Individual iterations vary, but the total movement is bounded. This is aggregate accounting, not an assumption about how many nested loops appear in the code.

**Round 8.** A positive-number window algorithm discards its leftmost number whenever the sum is too large. Why might that discard lose a future solution if a negative number can arrive next?

#### Worked answer — Explanation and next challenge

A negative arrival can bring the sum down without removing anything. For three followed by minus two, a target sum of one is achieved by both elements together. Discarding three solely because it exceeds one loses that possibility. The algorithm needs a monotonicity assumption about how extension changes the sum.

**Round 9.** Your current prefix sum is seven and the desired subarray sum is three. Which earlier prefix sum would complete a match, and why?

#### Worked answer — Explanation and next challenge

Four, because seven minus four is three. Think of the earlier prefix as the portion removed from the total to leave the desired suffix. If four appeared three times as a prior prefix, there are three matching starting positions. This explains why counting subarrays needs prefix frequencies rather than a set of prefix values.

**Round 10.** Why does the prefix-frequency map begin with zero occurring once? What happens on the single-element array containing three when the target is three?

#### Worked answer — Explanation and next challenge

The initial zero represents the position before any elements. At cumulative sum three, the required earlier prefix is zero, so the whole array is counted. Without it, ranges beginning at the first element are missed. This is a real mathematical boundary, not a patch added only to pass one test.

**Round 11.** For Product Except Self on two, zero, four, explain the left-product and right-product idea aloud. Why is dividing a total product an unreliable shortcut?

#### Worked answer — Explanation and next challenge

For each position, multiply everything strictly to its left by everything strictly to its right. The zero position gets two times four, which is eight; the other positions include a zero and get zero. A total product of zero cannot recover those answers through division. Prefix and suffix summaries keep the needed information without depending on invertibility.

**Round 12.** A service needs the longest consecutive stretch of events involving at most two regions. Then it changes to counting stretches whose signed adjustments total a target. Explain why these similar-looking stories lead to different state.

#### Worked answer — Explanation and mastery check

The region task has a maintainable window-validity condition: remove from the left until at most two region counts remain. The signed-adjustment task lacks the needed directional sum property but supports prefix differences. A strong answer names the invariant or equation before naming the algorithm. If you cannot explain the distinction, revisit the counterexample in round eight.

Give the tutor a new example of “what the next element needs to know about the past.” Let it change existence to counting or contiguous to noncontiguous, then explain which stored information must change.

## Worked program design: a session with at most two distinct titles

The function takes an ordered sequence of title IDs and returns the longest contiguous length containing at most two distinct titles. It preserves input order and returns zero for an empty sequence. “Contiguous” rules out sorting. We need counts within the active window, a left index, and the best length seen so far.

For each right index, increment the incoming title's count. While the count map contains more than two distinct keys, decrement the title at the left boundary, delete its map entry if the count becomes zero, and advance left. Once the window is valid, compare its length with the best answer. After the final right index, return the best length.

The shrinking loop is essential. One removal may not eliminate a distinct title if another copy remains. For A, A, B, C, adding C creates three distinct titles. Removing the first A leaves another A, so the window is still invalid. Removing the second A leaves B, C and restores the rule. The best length was three, from A, A, B.

The map invariant is exact counts of titles between left and right, inclusive. The validity invariant is checked after shrinking. Left never moves backward, so each event is added once and removed at most once. Total work is linear; the count map holds at most three keys during the immediate repair step for this fixed distinct-title bound.

To return the actual interval, store the best starting position whenever a new best length is found. Define ties: retaining the first best found returns the earliest longest interval. This extension changes output bookkeeping without changing the window invariant. Useful mental tests are an empty sequence, all one title, A-B-A, and A-B-C-B. They exercise absence, duplicates, a valid repeated title, and necessary shrinking.

**Optional spoken walkthrough:** act as the function on A, B, C, B. After each arrival, state left, the count map, current length, and best length. Then describe the return-value change needed to return both boundaries.
