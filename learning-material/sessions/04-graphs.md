# Chapter 4 — Trees graphs and dependencies

A graph problem asks how information moves through relationships. Make vertices, edges, and visitation rules explicit.

## Recognizing the graph hiding in a prompt

Movies linked by shared actors form a graph. Encoding jobs that depend on source assets form a directed graph. Adjacent land cells form an implicit graph. A filesystem directory structure is a tree if every node has one parent and links cannot create cycles. Start by defining the vertices and whether edges are directed, weighted, or generated from neighboring coordinates.

An adjacency list stores each vertex's neighbors. It uses O(V plus E) space and supports traversals in O(V plus E) time. An adjacency matrix uses quadratic space but gives constant time edge existence queries. Sparse interview graphs usually favor lists. A grid already represents its vertices; you do not need to allocate an explicit object for every edge.

## Breadth first search and shortest hops

To find the fewest recommendation links from one movie to another, enqueue the start and mark it discovered. Pop in first-in-first-out order. Enqueue unvisited neighbors with distance one greater. The queue processes all distance-zero nodes, then distance-one nodes, then distance-two nodes. Therefore the first discovery of a vertex has the shortest number of edges.

Mark on enqueue, not on dequeue, so two parents do not add the same vertex repeatedly. Keep a parent map if you must reconstruct the actual path. An empty queue before reaching the destination means unreachable. Breadth first search finds shortest paths when each edge has equal cost. If edges represent different download durations, shortest hops may not mean shortest time.

For nonnegative weighted edges, Dijkstra's algorithm uses a min heap of tentative total distances. When removing a heap record, skip it if its distance is no longer the best recorded one. Relax each outgoing edge by checking whether the current distance plus its cost improves the neighbor. With stale heap entries, a typical bound is O((V plus E) log E), often written O((V plus E) log V) for simple graphs. Negative edges break the greedy finalization argument. Dijkstra is a useful optional extension, not the first implementation to master.

**Think aloud.** A has an edge to B costing ten, an edge to C costing one, and C has an edge to B costing one. Why does shortest-hop BFS give the wrong cheapest route?

#### Worked answer — Coach notes

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

## Questions and worked explanations on relationships and reachability

**Round 1.** A movie can lead to another through a shared actor. What are the vertices and edges if the question asks for the fewest movie-to-movie hops? What detail about the connection should you clarify?

#### Worked answer — Explanation and next challenge

Movies are vertices and a qualifying shared-actor relationship creates an edge. Clarify whether the relation is symmetric and whether each hop has equal cost. You could also build an actor-movie bipartite graph, but then path lengths count different edge units. The representation must preserve the distance definition you intend to answer.

**Round 2.** A connects to B and C. B connects to D. Starting from A, describe the breadth first frontier after exploring A, and then after exploring B.

#### Worked answer — Explanation and next challenge

After A, the queue is B, C. After removing B and discovering D, it is C, D. D cannot jump ahead of C because it was discovered later. This is the layer-preserving behavior behind shortest unweighted paths. If you instead choose D next, you have changed the frontier policy toward depth first exploration.

**Round 3.** Both B and C point to D. Why should D be marked discovered when it is enqueued rather than only when removed?

#### Worked answer — Explanation and next challenge

When B first enqueues D, the mark lets C recognize that D already has scheduled work. Without it, both can add D. Repeated duplicates can inflate the queue and obscure parent selection. Registration at discovery makes “unvisited” mean not yet placed on the frontier, not merely not yet processed.

**Round 4.** A direct edge from A to B costs ten. A to C costs one and C to B costs one. Why does shortest-hop reasoning fail for minimum total cost?

#### Worked answer — Explanation and next challenge

The one-hop route costs ten, while two hops cost two. The layer order only measures number of edges, not their weight. Nonnegative weighted shortest paths need a frontier ordered by tentative total cost, such as Dijkstra's min heap. Before naming that algorithm, explain which quantity the frontier must prioritize.

**Round 5.** A grid has two land cells touching only at a corner. Is that one island or two? What does this reveal about converting a story into a graph?

#### Worked answer — Explanation and next challenge

It depends on whether adjacency includes diagonals. Four-way adjacency yields two islands; eight-way adjacency connects them. The graph's edges are part of the contract, not an implementation afterthought. A correct traversal over the wrong edges solves the wrong problem perfectly.

**Round 6.** A BST has root ten, left child five, and five's right child twelve. Every parent-child comparison near five looks sensible. Why is the whole tree invalid?

#### Worked answer — Explanation and next challenge

Every node in ten's left subtree must be less than ten, so twelve violates an inherited constraint. Checking only immediate children forgets ancestry. Pass bounds down or check strict inorder ordering. This is a general lesson: a local relationship can be valid while a global invariant fails.

**Round 7.** Jobs A and B must both finish before C. Which jobs are initially ready? After only A finishes, why is C still blocked?

#### Worked answer — Explanation and next challenge

A and B are ready if they have no prerequisites. C starts with two outstanding prerequisites and drops to one after A finishes. It becomes ready only at zero. Using “at least one predecessor completed” would implement an OR rule rather than the required AND rule.

**Round 8.** A needs B and B needs A. What does the ready queue look like initially, and what observation lets you report a cycle?

#### Worked answer — Explanation and next challenge

Neither has zero outstanding prerequisites, so the queue is empty while two jobs remain. More generally, if the algorithm processes fewer jobs than exist, the unresolved graph contains a cycle. The useful test is not merely an empty queue; an empty queue after all jobs finish is success.

**Round 9.** A leads to B and C, and both lead to D. Is seeing D through both paths proof of a directed cycle? Explain what is missing.

#### Worked answer — Explanation and next challenge

No. There is no path returning from D to an ancestor. Shared reachability and cycles are different. In DFS, a back edge to a node still on the active recursion path indicates a directed cycle; an edge to a fully processed node need not. Kahn's method avoids that distinction by counting prerequisites.

**Round 10.** You clone two nodes A and B that point at each other. Why must the original-to-clone map record A's clone before recursively cloning A's neighbors?

#### Worked answer — Explanation and next challenge

When cloning B encounters A again, it must find the already-created A clone rather than start cloning A anew. Registering early both terminates recursion and preserves identity: B's neighbor is the same A clone, not a duplicate object. The clone can exist before all its links are populated.

**Round 11.** A graph has ten thousand isolated vertices and no edges. Is traversal cost proportional only to edges? Why must a complete complexity statement include vertices?

#### Worked answer — Explanation and next challenge

You still need to represent or visit those vertices to answer many whole-graph tasks. The usual adjacency-list traversal bound is vertices plus edges. A source-reachable traversal can restrict the count to the reachable subgraph, but a whole-graph scan includes disconnected vertices. Define the scope of the task before simplifying the bound.

**Round 12.** Explain a dependency scheduler to a teammate without saying “topological sort” until the end. Include readiness, state updates, success, and failure.

#### Worked answer — Explanation and mastery check

Count unfinished prerequisites per job. Queue jobs with none. Repeatedly finish a ready job and decrement the counts of its dependents, queuing those that reach zero. If every job is output, the order works; otherwise a cycle blocks the rest. If you can explain those operations and why readiness is safe, the algorithm name becomes a useful label rather than a memory crutch.

Ask the tutor to change one graph assumption: weighted edges, a duplicate edge, a disconnected component, or a deterministic output requirement. Explain exactly which part of your representation or argument changes.

## Worked program design: a prerequisite planner

The function takes a number of jobs and pairs of prerequisite and dependent IDs. It returns a valid execution order and a success flag. IDs lie in a specified range; invalid endpoints fail validation. A cycle means no complete order exists. Independent jobs may appear in any order unless the contract requests deterministic tie-breaking.

Create an adjacency slice for outgoing dependency edges and an integer indegree slice. For each validated pair, append the dependent to the prerequisite's neighbors and increment the dependent's indegree. The counts and edge list must treat duplicate pairs consistently. Initialize a queue with every zero-indegree job, including isolated jobs.

Use a head index to consume that queue. Each removed job becomes the next output. For every dependent, decrement its outstanding count; when the count reaches zero, append that dependent to the queue. No other condition makes a job ready. After the queue is exhausted, compare output size with the total number of jobs. Equal means success; smaller means a cycle prevented completion.

For A before C and B before C, initial counts are zero for A and B and two for C. Finishing A drops C to one. Finishing B drops C to zero and makes it ready. A, B, C is valid, as is B, A, C. Add an isolated D and it must still appear somewhere. Add C before A and a dependency cycle now blocks at least A and C.

Validation, graph construction, readiness processing, and final completeness checking are useful conceptual stages. Their combined work is linear in vertices plus edges. If the smallest available job must always be chosen, replace the ready queue with a min heap, accepting logarithmic ready-set operations. This changes tie policy rather than dependency meaning.

**Optional spoken walkthrough:** describe the fields and loop body precisely enough that another person can execute them on the three-job example. Explain why a job is enqueued only when its outstanding count becomes zero, and why the final size check is necessary.
