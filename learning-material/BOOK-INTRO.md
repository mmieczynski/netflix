# How to use this reading edition

This book prepares you to reason through a live Go coding interview with another senior engineer. The exercises are original study scenarios shaped by your preparation plan; they are not a verified list of Netflix questions.

The reading edition explains each problem directly. It gives a concrete input and expected output before introducing the technique, then traces state changes and shows Go code beside the reasoning. The separate conversation sessions remain available for teaching aloud, feedback, and one-question-at-a-time practice.

## Read the problem before the implementation

A technique is an answer to a particular question. A map might store earlier indices for Two Sum, counts of prefix boundaries for exact-sum subarrays, or stable node pointers for LRU. Knowing only "use a map" does not tell you what to store or when to update it.

For each topic, identify the output, assumptions, and slow baseline first. Then read the example table as the state after each operation. Compare successive rows to see exactly which field changed and why. An invariant is a statement that must remain true; the explanation connects that statement to the implementation's branches.

## Notation used throughout

| Notation | Meaning |
| --- | --- |
| `[2, 7, 11, 15]` | Ordered array or sequence of values |
| `{A: 6, B: 7}` | Map from keys to values; display order is not semantic |
| `[(A, 4), (B, 7)]` | Ordered records or tuples for illustration |
| `nums[i:j]` | Elements i through j-1; end position is excluded |
| `[start, end)` | Interval including start but excluding end |
| `(now-W, now]` | Rolling interval excluding its left boundary |
| `[A, B]` in a cache trace | Most recent A, then least recent B |

Illustrative tuples and maps are notation, not always literal Go syntax. Fenced Go blocks use actual Go declarations. Code omits package and import boilerplate for readability; examples using sorting, heap operations, or locks use `sort`, `container/heap`, or `sync`. The repository's example package assembles all blocks for compiler and behavior checks.

## Reading code without running it

Before a loop, name what each variable means. For one row of the trace, follow the lookup, decision, and update in order. For prefix counting, for example, the map stores earlier boundaries before the current one is inserted. Reversing those lines changes the meaning and can introduce an empty-subarray bug.

Some algorithms are complete functions; small service examples explicitly identify their surrounding preconditions or omitted methods. The LRU, weighted cache, TTL map, transaction stack, and final metadata-cache baseline are assembled implementations. No paragraph asks you to discover an unstated solution before the book becomes useful.

Small examples keep the state visible. Diagrams show a specific relationship or transition, while tables give exact values. Shading highlights active state but is not required to understand the figure on a monochrome screen. Code blocks may continue across pages with a continuation label.

## The route through the material

Chapters 1-5 develop representation, hashing, windows, ordering, traversal, stacks, and choice reasoning. Chapters 6-9 apply them to caches, deadlines, concurrency, limits, and event services. Chapter 10 covers ranking and rolling statistics. Chapter 11 models practical programs. Chapter 12 combines the concepts into two complete interview-style cases.

Read unfamiliar topics in order. On a later pass, you can cover the next table row and predict it, but the worked answers remain fully available. A change in assumptions is often the most useful check: negative numbers break a positive-only window proof; varying TTLs break insertion-order cleanup; falling scores break winner-only storage.

With this repository available to ChatGPT, ask to study or continue a topic. The tutor retrieves the separate conversation session and can use these reading examples too. No prompt copying, editor work, or code execution is required for independent reading or spoken learning. Progress is updated only from actual study, not from authoring changes to this book.
