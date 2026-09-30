# Chapter 3 — Sorting binary search intervals and heaps

Ordering can turn a global search into a local decision. Decide whether you need full order, a boundary, or just the next best item.

## Two pointers after sorting

For Three Sum, find distinct value triples whose sum is zero. The cubic baseline tries all triples. Sort first, then fix one position and search the remaining suffix with a left and right pointer. If the sum is too small, moving right inward cannot increase it; move left forward. If it is too large, move right backward. Each pointer move rules out a family of impossible pairs.

For minus four, minus one, minus one, zero, one, two, fixing minus one allows the pairs minus one with two and zero with one. Skip repeated fixed values and repeated pair values after emitting a triple. The result is unique by values, not by indices. Sorting loses original positions unless you retain them. The full complexity is quadratic time after sorting, with output space proportional to the number of triples; copying to avoid input mutation adds linear space.

**Think aloud.** Why can you move the left pointer when the sum is too small? What property makes that argument fail on an unsorted input?

#### Worked answer — Coach notes

With the same left value, every smaller right index gives a value no greater than the current right value, so none can fix a sum that is already too small. Sorted order supplies the comparison. “Because that is the pattern” is not a proof.

## Merge intervals by remembering one frontier

Suppose sessions occupy half-open intervals: start included, end excluded. Sort by start. Keep a current interval; if the next start is less than its end, the intervals overlap, so extend the end to the larger end. Otherwise emit the current interval and begin another. Once a gap is found, no later interval can bridge backward across it, because later starts are no smaller.

For two to five, one to three, and seven to nine, sorting gives one to three, two to five, seven to nine. The first two merge into one to five. The gap before seven finalizes that answer. If touching intervals should merge as continuous coverage, use less than or equal instead of strictly less. For concurrent viewer counts, touching half-open intervals do not overlap: a departure at time five is processed before an arrival at time five.

A sweep line generalizes merging to “maximum simultaneous sessions.” Turn starts into plus one events and ends into minus one events, sort by time with the agreed tie order, and track a running count and maximum. Merging answers union coverage; sweeping answers overlap counts. Reusing one algorithm without noticing this distinction is a common modeling error.

## Binary search asks for the first true position

Binary search requires a monotone predicate: false, false, then true for the rest. To find the latest metadata version at or before time fifteen in timestamps ten, twenty, thirty, search for the first timestamp greater than fifteen. That is index one; the answer is the previous entry, index zero. If the first true index is zero, there is no suitable version. If no index is true, all versions qualify and the final one is the answer.

```go
i := sort.Search(len(versions), func(i int) bool {
    return versions[i].Time > queryTime
})
if i == 0 { return "", false }
return versions[i-1].Value, true
```

Go's [sort.Search](https://pkg.go.dev/sort#Search) returns the first true index or the length, not minus one. The important skill is defining the predicate and handling both edges. For a hand-written version, use a half-open search range and compute middle as low plus half of high minus low. Every iteration must shrink the unresolved range.

In a rotated sorted array with distinct values, at least one half around the midpoint is sorted. Determine which half, check whether the target belongs in its value range, and discard the other. With duplicates, values at both ends and midpoint may fail to identify a sorted side; discarding equal endpoints can degrade to linear time. Do not promise logarithmic worst-case behavior under an assumption you have removed.

## A heap keeps only the ordering you need

A binary min heap is usually stored in a slice representing a complete binary tree. The root is index zero. For a node at index i, its children are at two times i plus one and two times i plus two; its parent is at integer division of i minus one by two. The invariant is that each parent is no greater than either child. This guarantees the root is smallest, but siblings and distant branches need not be ordered.

Insert at the end and move upward while the parent comparison is violated. To remove the minimum, move the last item to the root and repeatedly swap downward with the smaller child until the invariant holds. The height is logarithmic, so insertion and removal cost logarithmic time. Inspecting the root is constant time. Building a heap from a whole slice by repairing internal nodes bottom-up takes linear time; inserting every element separately costs O(n log n). Try drawing the slice three, eight, five, ten, nine: it is a valid heap even though eight appears before five.

To return the top k frequent titles, first count occurrences with a map. Sorting all u unique titles costs O(u log u) and is a fine baseline. If k is small, maintain a heap of at most k winners whose root is the worst retained candidate. A new candidate worse than the root cannot belong in the top k; a better one replaces it. The invariant is that after any processed prefix, the heap contains that prefix's best k candidates.

The complexity is expected O(n) for counting and O(u log k) selection, followed by O(k log k) if output must be ranked. Define deterministic ties, such as lower title ID winning equal counts. That means the heap's worst element has lower count, or a lexicographically larger ID at the same count. A reversed tie comparison silently returns the wrong boundary item. When k is one, treat the heap step as constant work.

Go's heap interface supplies length, comparison, swapping, appending, and removing the final slice item. Your type's Pop method removes the final element; `heap.Pop` rearranges the root there first. Calling the type's Pop directly bypasses heap repair. Changing a priority requires `heap.Fix` with a maintained index, or a fresh versioned record that makes old records stale. [Go heap documentation](https://pkg.go.dev/container/heap).

## Why order lets you stop looking

An efficient search is often built around a justified refusal to inspect some candidates. In an unordered list, discovering that one number is too small says little about the others. In a sorted list, that observation tells you about an entire region. The sorting step creates evidence that later decisions can exploit. When explaining an ordering algorithm, identify which candidates a step rules out and why they cannot contain the answer.

Binary search is not restricted to finding an exact value. Its deeper task is locating a boundary between false and true. Suppose versions exist at times ten, twenty, and thirty and you query time twenty. The predicate “version time is greater than twenty” is false, false, true. The first true version is too new; the previous version is the newest allowed one. The exact-match problem disappears into a boundary definition. This same reasoning supports first capacity that succeeds, first invalid position, and insertion points, provided the predicate really is monotone.

A heap deliberately creates less order. Think of a tournament where you know the current weakest retained winner, but you have not ranked every contestant. To select the best two candidates, compare each newcomer with that weakest winner. If the newcomer loses, it cannot enter the best two. If it wins, replace the weakest and repair the tournament. The complete order is unnecessary until the final two need to be printed.

Sorting and a heap are therefore not rival tricks with one always superior. If you need all results ranked, sorting is direct. If you only need a few winners from a large collection, maintaining a small boundary can save work. If you repeatedly update priorities and ask for the next task, a heap can retain useful state between calls. But if priorities change without heap repair, the stored comparisons become stale: the old tournament result is no longer evidence of today's winner.

## Questions and worked explanations on ruling candidates out

**Round 1.** In sorted values one, three, five, eight, you want a pair summing to six. Starting with one and eight gives nine. Why is moving the right pointer left justified?

#### Worked answer — Explanation and next challenge

With eight fixed, all values at or to the right of the current left position are at least one, so no pair using eight can reduce the sum to six. Discard eight from consideration. Now one and five gives six. The safe move follows from the sorted order; arbitrary pointer movement without that ordering would lose possible answers.

**Round 2.** You sort viewing events to apply a two-pointer method, but the task asks for a consecutive session in the original order. What did sorting change that matters?

#### Worked answer — Explanation and next challenge

It changed adjacency and therefore which ranges qualify as sessions. Sorting may preserve a set of values while destroying the relationship the query is about. Ask whether input order is meaningful before using sorting as an optimization. If original indices are needed only for output, retaining indices may be enough; if contiguity is meaningful, it is not.

**Round 3.** Intervals are one to three, two to five, and seven to nine. After sorting by start, why can you finalize one to five when you encounter seven?

#### Worked answer — Explanation and next challenge

Every later start is at least seven, so no later interval can begin early enough to bridge the gap from five to seven. The current merged component is complete. This is the algorithm's frontier: everything emitted is settled, and only the last interval remains open to extension.

**Round 4.** One session ends at five and another starts at five. Must they overlap? How would the answer affect merging and counting simultaneous sessions?

#### Worked answer — Explanation and next challenge

Under half-open intervals they do not overlap: the first excludes its endpoint. For coverage display you may still merge touching intervals into continuous coverage, but that is a chosen output policy. For simultaneous counts, process departure before arrival at equal time. Boundary semantics come before comparison operators.

**Round 5.** Versions occur at ten, twenty, and thirty. Query twenty using the first version strictly greater than the query. What is the predicate sequence and the answer?

#### Worked answer — Explanation and next challenge

The sequence is false, false, true. The first true position is the time-thirty version, and its predecessor is time twenty. Query before ten and the first true position is zero, meaning there is no predecessor. Query after thirty and the search returns the length, so the final version is the answer.

**Round 6.** An engineer wants to binary-search the predicate “this title is currently popular” over an arbitrary title list. What must be established before binary search is valid?

#### Worked answer — Explanation and next challenge

The predicate must switch from false to true at most once in the searched order. Popularity labels over arbitrary IDs do not supply that structure. Binary search saves work only when the unvisited range can be eliminated using a monotonicity argument. Fast random access by itself is insufficient.

**Round 7.** A min heap's slice begins three, eight, five. Is that already a violation because eight is before five? Explain what order the heap actually promises.

#### Worked answer — Explanation and next challenge

No. Three is the parent of both eight and five and is no greater than either. The siblings need not be ordered. A heap promises parent-child relationships sufficient to expose a minimum at the root. It does not promise that iterating the slice produces sorted values.

**Round 8.** Keep the largest two scores from eight, five, nine, six. State the retained set and its weakest member after each step.

#### Worked answer — Explanation and next challenge

After eight, retain eight. After five, retain eight and five, weakest five. Nine replaces five, leaving nine and eight, weakest eight. Six is rejected. The discarded values are not needed for the static top-two query because at least two retained values are no worse. That reasoning will need revision if retained scores can later decrease.

**Round 9.** At an equal score, lexicographically smaller IDs win. The retained candidates include A and Z at that score. Which is worse and belongs nearer the root of a heap of weakest winners?

#### Worked answer — Explanation and next challenge

Z is worse because A wins the tie. The heap comparator reverses the final ranking: lowest score first, and largest ID first at equal scores. State the final ordering in words, then derive the heap ordering from “worst winner.” This avoids memorizing a comparator with unexplained inequality signs.

**Round 10.** A task's priority changes while its record stays in the same heap position. Why can the next popped task be wrong, and what are two ways to repair the design?

#### Worked answer — Explanation and next challenge

The stored parent-child ordering may no longer hold. Maintain an index and repair the modified position, or insert a new versioned record and ignore old versions when popped. The first requires index maintenance on swaps; the second requires stale-record handling and memory accounting. Neither is a free update.

**Round 11.** You have three already-sorted event streams. Why can a heap containing only each stream's next event find the next event globally?

#### Worked answer — Explanation and next challenge

Every later event in a stream is no earlier than its current head. Therefore the smallest among the heads is no later than any unseen event. Remove it and replace it with that stream's next event. The proof relies on each stream's internal ordering; one out-of-order stream breaks it.

**Round 12.** Choose between sorting, binary search, and a heap for three requests: print every title in rank order; find the latest stored version before a time; repeatedly run the next scheduled task. Defend each choice by the query it accelerates.

#### Worked answer — Explanation and mastery check

Sort for a full ranking; binary-search an already ordered version history for a predecessor; use a deadline heap for repeated next-task selection with updates. Mention the cost of creating or maintaining each order. If you propose binary search on unsorted history or a heap for constant-time arbitrary-ID lookup, revisit the missing prerequisite.

Tell the tutor which candidates each algorithm can safely ignore. Being able to explain that discarded region is a better mastery test than being able to recite a loop.

## Worked program design: a time-based value store

Expose Set(key, value, timestamp) and Get(key, queryTime). Each key has its own ordered history. Writes must be nondecreasing for that key; an earlier timestamp is rejected, and a write at the same timestamp replaces the latest value. Get returns a value plus a found flag. An empty stored string is different from no suitable version.

The store has a map from key to a slice of version records. Each record contains a timestamp and value. Set finds that key's history. If it is empty, append the new version. Otherwise compare with the final timestamp: reject a smaller timestamp, replace the final value if equal, or append if larger. Reject before mutating so a failed call preserves the sorted-history invariant.

Get finds the history and searches for the first version whose timestamp is greater than the query. If that position is zero, return missing. Otherwise return the version immediately before it. If the search finds no greater version, its returned position is the history length, and the final version is the answer. A helper called firstAfter can isolate that boundary-search responsibility from the public method.

With quality at time ten equal to HD and at twenty equal to UHD, Get at fifteen returns HD. At nine it returns missing; at twenty it returns UHD; at thirty it still returns UHD. Replacing the time-twenty value with HDR changes the last three relevant answers without adding a duplicate timestamp. Attempting a time-twelve write afterward is rejected under this contract.

Set is amortized constant time for ordered writes. Get is logarithmic in the versions for the requested key, plus expected map lookup. Allowing arbitrary historical insertion changes Set: a sorted slice needs a search and potentially linear shifting. The new write contract should not silently invalidate Get's binary-search precondition.

**Optional spoken walkthrough:** define the version record and all Set branches, then explain the firstAfter helper without Go syntax. Test the first, equal, between, and after-last query boundaries in words.
