# Sources and scope

Original scope research checked on 30 September 2026. Reading edition revised on 4 October 2026; Go specification, heap, sync, and memory-model references rechecked for the revised examples. The exercises and learning sequence are original preparation material; they are not a leaked or verified Netflix question bank.

## What informed the preparation

Your plan.md is the primary scope input. It supplies the core algorithm sequence and practical topics, particularly LRU, TTL, rate limiting, event processing, and iterative requirements. Your clarification establishes an L5 live interview with another senior engineer. No interview date or team-specific preparation guide was supplied, and no chapter requires a particular study duration.

Internet searches covered Netflix technical-screen descriptions, CodeSignal's own documentation, recommendation-system sources, and Go references. Search results contained conflicting commercial guides and secondhand reports. Those reports were not treated as proof of what you will be asked. The live format and level come from you, not from those reports.

## Primary sources used

- [CodeSignal live interviewing](https://support.codesignal.com/hc/en-us/articles/360045985453-Conducting-an-interview-in-CodeSignal) supports the distinction between an interviewer-led session and an assessment. It does not specify Netflix's content or evaluation rubric.

- [CodeSignal practice overview](https://support.codesignal.com/hc/en-us/articles/12984563824279-Practice-Content-Overview) supports practising the IDE and question formats. Your invitation remains the source for session-specific constraints.

- [CodeSignal language support](https://support.codesignal.com/hc/en-us/articles/360061267734-What-languages-are-supported-in-CodeSignal-Interview) lists platform language support visually. We do not infer a particular Go runtime version from it; verify that in your actual environment.

- [Netflix recommendation overview](https://help.netflix.com/en/node/100639) supplies high-level product context, not a coding interview syllabus.

- [Google candidate generation overview](https://developers.google.com/machine-learning/recommendation/overview/candidate-generation) explains content-based and collaborative retrieval. The toy ranking formulas and tasks here are our own.

- [Go container/heap](https://pkg.go.dev/container/heap) and [sort.Search](https://pkg.go.dev/sort#Search) support heap interface and binary-search API details. Current source documentation was also fetched through Context7 using /golang/go.

- [Go sync](https://pkg.go.dev/sync), [memory model](https://go.dev/ref/mem), and [time](https://pkg.go.dev/time#hdr-Monotonic_Clocks) support synchronization and clock discussion.

- [Go specification](https://go.dev/ref/spec) and [strings, bytes, runes and characters](https://go.dev/blog/strings) support language behavior.

- [Go test command](https://pkg.go.dev/cmd/go#hdr-Test_packages) and [race detector](https://go.dev/doc/articles/race_detector) explain local validation tools. Race detection checks executed paths and does not prove absence of all races.

## Coverage of the original plan

Two Sum, anagrams, windows, and array products appear in lesson 2. Three Sum, intervals, binary search, top k, and kth-element selection share lesson 3's ordering tools. Tree traversal, BST validation, islands, course dependencies, and graph cloning appear in lesson 4. LRU, weighted capacity, TTL, combined expiration, and concurrency are lessons 6 and 7. The two rate limiters are lesson 8. Time-based storage, deduplication, counters, and scheduling are lesson 9. Streaming statistics and recommendation tasks are lesson 10. Filesystem, transactions, pub/sub, and playback-service composition are chapter 11 and its worked program design.

The reading edition now includes 65 Go blocks extracted into one compiler-checked
example package. Algorithms and the LRU, weighted cache, TTL map, transaction
stack, and combined metadata-cache baseline are assembled implementations.
Filesystem resolution, pub/sub publishing, and scheduler claiming/draining are
explicitly labelled cores with their surrounding preconditions. They are not
presented as complete implementations of every optional service method.
The final ranking case composes the earlier comparator with filtering,
deduplication, and aggregation. The conversation edition remains separate.

## A further connection

For kth largest in a stream, use a size-k min heap over all seen values; duplicates count as separate observations unless specified otherwise. For a static array, sorting is a simple baseline and quickselect offers expected linear time, but its worst case and pivot handling require care. For this first-round route, reliable heap reasoning is the higher-priority transferable skill.

## Limits of the material

This course provides focused preparation, not exhaustive coverage of every possible interview problem. Deep distributed systems, advanced ML training, and a full Go language course are outside this first coding-round scope. The book can be read independently without an editor, running tests, or an AI. Spoken implementations develop reasoning but are not evidence that code has compiled.

## Repository and conversation access

The book is independent of a particular ChatGPT interface. For the repository
workflow, an agent must actually retrieve the relevant files. Do not infer that
remote read access includes permission or tools to update progress.
[Official ChatGPT Voice documentation](https://learn.chatgpt.com/docs/features/voice)
describes conversational turn-taking and task permissions.
[Official workspace connection guidance](https://learn.chatgpt.com/docs/enterprise/shared-connections)
distinguishes resource access from write actions and recommends checking them
separately. These pages do not verify this reader's account configuration.

