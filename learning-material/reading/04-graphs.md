# Chapter 4 — Trees graphs and dependencies

A graph problem asks how information moves through relationships. Make vertices, edges, and visitation rules explicit.

## Recognizing the graph hiding in a prompt

Movies linked by shared actors form a graph. Encoding jobs that depend on source assets form a directed graph. Adjacent land cells form an implicit graph. A filesystem directory structure is a tree if every node has one parent and links cannot create cycles. Start by defining the vertices and whether edges are directed, weighted, or generated from neighboring coordinates.

An adjacency list stores each vertex's neighbors. It uses O(V plus E) space and supports traversals in O(V plus E) time. An adjacency matrix uses quadratic space but gives constant time edge existence queries. Sparse interview graphs usually favor lists. A grid already represents its vertices; you do not need to allocate an explicit object for every edge.

## Breadth first search and shortest hops

To find the fewest recommendation links from one movie to another, enqueue the start and mark it discovered. Pop in first-in-first-out order. Enqueue unvisited neighbors with distance one greater. The queue processes all distance-zero nodes, then distance-one nodes, then distance-two nodes. Therefore the first discovery of a vertex has the shortest number of edges.

Mark on enqueue, not on dequeue, so two parents do not add the same vertex repeatedly. Keep a parent map if you must reconstruct the actual path. An empty queue before reaching the destination means unreachable. Breadth first search finds shortest paths when each edge has equal cost. If edges represent different download durations, shortest hops may not mean shortest time.

For nonnegative weighted edges, Dijkstra's algorithm uses a min heap of tentative total distances. When removing a heap record, skip it if its distance is no longer the best recorded one. Relax each outgoing edge by checking whether the current distance plus its cost improves the neighbor. With stale heap entries, a typical bound is O((V plus E) log E), often written O((V plus E) log V) for simple graphs. Negative edges break the greedy finalization argument. Dijkstra is a useful optional extension, not the first implementation to master.

**Failure case.** A one-edge route costing ten is more expensive than two edges costing one each.

The direct route has one edge but cost ten. A to C to B has two edges and cost two. State whether the task minimizes edge count or total weight before choosing a traversal.

## Depth first search and connected components

Number of Islands asks how many connected land components exist. Scan every cell. On an unvisited land cell, increment the component count and explore all reachable land with a stack or recursion. Mark visited before pushing neighbors. The outer scan finds component entry points; the traversal prevents counting the same island twice. Clarify four-way versus diagonal adjacency and whether mutating the grid is allowed.

For a large grid, an explicit stack avoids relying on deep recursion. Each cell is visited once and each has a constant number of neighbors, so time is O(rows times columns). A separate visited structure uses the same order of extra space. Mutating land to water avoids a separate visited map but the traversal stack can still grow linearly.

## Tree invariants can be global

Validating a binary search tree by comparing each node only with its children is insufficient. A value of twelve can be the right child of a five inside the left subtree of a ten: its local relationship to five is fine but its relationship to ten is wrong. Pass allowable lower and upper bounds down the recursion, or perform an inorder traversal and verify strictly increasing values when duplicates are forbidden.

Inorder means visit left subtree, then node, then right subtree. Keep a previous value and a boolean saying whether it exists; avoid using the smallest integer as an uninitialized sentinel. A tree traversal costs O(n) time and O(h) call-stack space, with h potentially n for a skewed tree. Level order traversal is simply BFS with either a captured queue length per level or explicit distances.

## Dependencies and cycles

For job prerequisites, direct each edge from prerequisite to dependent. Count each vertex's incoming edges. Put all zero-indegree vertices in a queue; remove one, output it, and decrement its dependents. A dependent becomes ready when its indegree reaches zero. If you output fewer than V vertices, the remaining subgraph contains a cycle. This is Kahn's topological sort.

With A before C and B before C, A and B may appear in either order, but C must follow both. If deterministic output is required, use a min heap of ready IDs rather than relying on map iteration. Clarify whether duplicate edges are distinct requirements or should be deduplicated. A consistent adjacency list and indegree count can accommodate duplicates, but mismatching their treatment is a bug.

To clone a graph, memoize original pointer to clone pointer. Create and register the clone before recursively cloning neighbors, otherwise a cycle causes infinite recursion. Reuse that clone whenever the original is encountered again. A map keyed only by a displayed label is unsafe if different vertices may share labels.

## Think about the frontier before thinking about recursion

Every traversal has discovered territory and territory still waiting to be explored. The waiting territory is its frontier. Breadth first search chooses the oldest discovered work next; depth first search chooses the newest unfinished work next. That small policy difference changes the order in which facts become known. Neither traversal needs a particular product story: movies, cities, cells, and jobs can all be vertices.

Consider A connected to B and C, with B connected to D. Breadth first search discovers B and C one step from A before it explores D two steps away. The queue preserves layers. Depth first search might follow A, B, D before returning to C. That is useful when exploring an entire branch, but the first route it finds is not necessarily the shortest route in edge count. To choose between them, identify the property you need from the discovery order.

Visited state is not just a performance trick. It gives a meaning to discovery. Without it, two nodes pointing at one another can generate work forever. For a queue traversal, marking a node when it enters the queue means it has already claimed a place on the frontier. Waiting until removal allows multiple parents to enqueue it. In a graph-cloning problem, a similar early registration says that an original node already has a clone identity, even if that clone's neighbors are not fully built yet.

Dependency processing uses a different notion of readiness. A discovered job is not necessarily runnable. It becomes ready only after all its prerequisite edges have been accounted for. Indegree records the number still blocking it. Removing a ready job reduces that count for its dependents. If progress stops while unresolved jobs remain, you have a structural obstruction: every remaining job still needs something inside the unresolved group. In a finite directed graph, following those dependencies must eventually revisit a node, exposing a cycle.

Do not confuse shared descendants with cycles. If A leads to B and C, and both lead to D, D is shared but the graph can still be acyclic. A traversal needs visitation to avoid duplicate work; cycle detection needs evidence of a path that returns to an active ancestor, or equivalent indegree reasoning. “I saw this node before” alone is not enough to diagnose a cycle in a directed graph.

## Worked example: two routes to the same vertex

Consider edges A to B, A to C, B to D, and C to D. Starting at A, find the minimum number of edges to D. Breadth-first search processes the frontier in layers. The queue is not just storage for pending work; its order represents increasing distance from the start.

![A diamond-shaped graph. B and C both discover D, but only the first discovery enqueues it.](figures/graph.svg)

| Removed from queue | Newly discovered | Queue afterward |
| --- | --- | --- |
| A, distance 0 | B and C, distance 1 | B, C |
| B, distance 1 | D, distance 2 | C, D |
| C, distance 1 | None; D already discovered | D |
| D, distance 2 | None | Empty |

Mark a vertex when it is enqueued. When C examines its edge to D, D is already discovered even though it has not yet been removed from the queue. Marking only on removal allows duplicate queue entries. Save D's predecessor as B on first discovery; following predecessors backward reconstructs A, B, D. A, C, D is equally short, so a requirement for deterministic paths may also specify neighbor order.

The distance argument assumes every edge has equal cost. If A to B costs ten while A to C and C to B each cost one, the one-edge route is not cheapest. The frontier must then be ordered by accumulated cost rather than by hop layer, as in Dijkstra's algorithm for nonnegative weights.

## Worked example: a dependency queue means something else

Interpret the same diamond as build dependencies: A must finish before B and C; both B and C must finish before D. The initial indegrees are A:0, B:1, C:1, D:2. The ready queue starts with A. Completing A makes B and C ready. Completing B reduces D to one unfinished prerequisite; it does not release D. Completing C reduces D to zero and releases it.

The queue now means all prerequisites are satisfied, not minimum distance. This is why naming the container is only half a solution. Its membership rule carries the correctness argument.

Add a dependency D to A. Every vertex now has a positive indegree, so no work is initially ready. The process finishes with fewer completed vertices than exist in the graph, proving the dependency graph contains a cycle. In a graph with an independent vertex E, E could finish while the cyclic component remains stuck. An empty final queue by itself therefore says nothing; compare completed count with total count.

Across both uses, store every declared vertex, including isolated ones. Building the vertex set only from outgoing edges can silently lose a leaf or a standalone job.

## Putting the program together: a prerequisite planner

The function takes a number of jobs and pairs of prerequisite and dependent IDs. It returns a valid execution order and a success flag. IDs lie in a specified range; invalid endpoints fail validation. A cycle means no complete order exists. Independent jobs may appear in any order unless the contract requests deterministic tie-breaking.

Create an adjacency slice for outgoing dependency edges and an integer indegree slice. For each validated pair, append the dependent to the prerequisite's neighbors and increment the dependent's indegree. The counts and edge list must treat duplicate pairs consistently. Initialize a queue with every zero-indegree job, including isolated jobs.

Use a head index to consume that queue. Each removed job becomes the next output. For every dependent, decrement its outstanding count; when the count reaches zero, append that dependent to the queue. No other condition makes a job ready. After the queue is exhausted, compare output size with the total number of jobs. Equal means success; smaller means a cycle prevented completion.

For A before C and B before C, initial counts are zero for A and B and two for C. Finishing A drops C to one. Finishing B drops C to zero and makes it ready. A, B, C is valid, as is B, A, C. Add an isolated D and it must still appear somewhere. Add C before A and a dependency cycle now blocks at least A and C.

Validation, graph construction, readiness processing, and final completeness checking are useful conceptual stages. Their combined work is linear in vertices plus edges. If the smallest available job must always be chosen, replace the ready queue with a min heap, accepting logarithmic ready-set operations. This changes tie policy rather than dependency meaning.
