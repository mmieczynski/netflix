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

**A boundary that must not move backward.** In A, B, B, A, assigning left to the previous occurrence plus one without comparing it with the current left boundary is incorrect.

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

## Worked example: watch a sliding window repair itself

Find the longest contiguous session containing at most two distinct titles in A, B, A, C, C. The slow baseline starts at every position and scans until a third title appears. Adjacent starts repeat much of the same counting. A moving window reuses those counts.

![When C introduces a third title, two removals repair the window to A, C. The next C can then extend it without another removal.](figures/window.svg)

Keep a left boundary, a right boundary, and a frequency map for the records between them. Adding the first C gives counts A:2, B:1, C:1. Removing the leftmost A is insufficient: A still has count one, so three distinct titles remain. Remove B next and delete its zero-count entry. Now the window is A, C and is valid again.

| New title | Window after repair | Counts | Best length |
| --- | --- | --- | --- |
| A | A | A:1 | 1 |
| B | A, B | A:1, B:1 | 2 |
| A | A, B, A | A:2, B:1 | 3 |
| C | A, C | A:1, C:1 | 3 |
| C | A, C, C | A:1, C:2 | 3 |

The repair is a loop, not a single conditional. Its stopping condition is that the number of positive counts is at most two. Each record enters once and leaves at most once, so the total number of boundary moves is linear even though one arrival can trigger several removals. Record the best length only after repair; otherwise an invalid window can become the answer.

## Worked example: negative numbers change the approach

Now count subarrays whose sum is two in the array 2, -1, 1. A shrinking window based on whether the sum is too large is unreliable because removing a negative number increases the sum. The needed relationship is between prefix sums.

The running prefix sums, including the empty prefix, are 0, 2, 1, 2. A subarray sums to two exactly when its ending prefix is two larger than its starting prefix. Maintain frequencies of earlier prefix sums. Seed zero with frequency one so a subarray starting at the first element can be counted.

At running sum two, look for zero and count the first element. At running sum one, look for minus one and find none. At the final sum two, look for zero again and count the whole array. There are two answers. Insert each current prefix only after looking up the required earlier value; that order prevents counting a zero-length subarray when the target is zero.

The two examples both process an array from left to right, but their proofs are different. The title window can repair a violation by discarding a prefix. Prefix counting instead remembers all earlier boundary values that could complete a valid difference. Recognizing that distinction is more useful than identifying both as map problems.

## Putting the program together: a session with at most two distinct titles

The function takes an ordered sequence of title IDs and returns the longest contiguous length containing at most two distinct titles. It preserves input order and returns zero for an empty sequence. “Contiguous” rules out sorting. We need counts within the active window, a left index, and the best length seen so far.

For each right index, increment the incoming title's count. While the count map contains more than two distinct keys, decrement the title at the left boundary, delete its map entry if the count becomes zero, and advance left. Once the window is valid, compare its length with the best answer. After the final right index, return the best length.

The shrinking loop is essential. One removal may not eliminate a distinct title if another copy remains. For A, A, B, C, adding C creates three distinct titles. Removing the first A leaves another A, so the window is still invalid. Removing the second A leaves B, C and restores the rule. The best length was three, from A, A, B.

The map invariant is exact counts of titles between left and right, inclusive. The validity invariant is checked after shrinking. Left never moves backward, so each event is added once and removed at most once. Total work is linear; the count map holds at most three keys during the immediate repair step for this fixed distinct-title bound.

To return the actual interval, store the best starting position whenever a new best length is found. Define ties: retaining the first best found returns the earliest longest interval. This extension changes output bookkeeping without changing the window invariant. Useful mental tests are an empty sequence, all one title, A-B-A, and A-B-C-B. They exercise absence, duplicates, a valid repeated title, and necessary shrinking.
