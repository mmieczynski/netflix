# Learning to solve coding problems

This is the reading edition of the Netflix L5 Go interview course. It prepares you for the first live technical interview with another senior engineer. You can read it independently on an e-reader, without an editor, a timer, or an AI conversation. The companion conversation edition teaches the same topics through short explanations, questions, feedback, and verbal program design.

The book teaches through explanations, worked examples, diagrams, state tables, and Go code. Each chapter has an explained Go example connecting its invariant to concrete fields, branches, and mutations. Some are complete algorithms; cache, expiration, and scheduling examples isolate a core operation and explicitly state what the surrounding program supplies. You can trace the examples on the page without executing them. The PDF is not a transcript of the tutor's question bank.

## What learning an approach means

A memorized solution tells you what to do when the prompt looks familiar. An approach tells you what information a problem needs, which operations are expensive, and what must remain true as the state changes. That lets you recognize when a familiar method fits a new story and when a changed assumption makes it fail.

For example, a queue can support expiration when timestamps arrive in order. It is not enough to remember “time window means queue.” You should understand why expired timestamps form a prefix, why you can stop when the first fresh timestamp appears, and why out-of-order arrival breaks that reasoning. Once you understand those facts, you can apply the idea to rate limits, recent-hit counts, and expiring viewing events.

The recurring method is to define behavior, construct a correct baseline, identify repeated work, choose state that removes that repeated work, and justify what the state means. Then test the argument against a tiny example and a changed requirement. You will see this process across arrays, graphs, caches, and small services.

## A useful kind of mental exercise

Imagine a cache with room for two titles. After writing A and then B, B is most recent. Reading A moves A ahead of B. Writing C now evicts B. You can execute this sequence in your head and explain which map entries and list links must change. You are practising implementation reasoning even though you are not typing pointer assignments.

A good verbal implementation is more specific than “use a map and a list.” It says what each field represents, what a hit and a miss do, which helper performs a removal, and how every return path preserves the contract. The worked program designs model that precision. An AI tutor can play the caller or the computer while you supply the next operation.

Reading and spoken reasoning do not prove that a program compiles or that you can express it fluently under an interview timer. Those are separate skills. This book is meant to make time away from your computer useful, without pretending that you are practising syntax while walking, commuting, or reading on an e-reader.

## How to read the chapters

Read the explanations in order the first time a topic is unfamiliar. The examples supply their reasoning and results directly. You do not need to answer a question to reach the explanation. On a later pass, cover the next table row and predict it if that helps you check understanding. This is optional; the book remains complete when read straight through.

Read a diagram as a picture of a specific state, not as a complete implementation. Arrows show relationships or direction of processing. Shaded boxes emphasize active or selected state. Labels and the surrounding explanation carry the meaning even on a monochrome screen.

State tables show what changed after each operation. Unless explicitly stated otherwise, a row is the state after the operation and its required cleanup have completed. Compare neighboring rows: identify the new input, the field that changes, and the reason another field stays unchanged. This makes bookkeeping visible without requiring code execution.

In the cache example, a Get changes recency but not the number of residents. In a weighted overwrite, both recency and the weight total may change. In expiration cleanup, a stale deadline record can disappear while the current cached value survives. These differences reveal the invariants more clearly than a picture of the final result alone.

The program design at each teaching chapter's end assembles the concepts into a complete small program in plain language. Trace the included example if you wish, but the text already states what happens and why. No section depends on your discovering an unstated solution before the book becomes useful.

Small examples are intentional. Holding a long array in working memory is not the skill being taught. When discussing a chapter, ask the tutor to repeat the current state or shorten the example. You should spend your attention on the invariant and the decision, not remembering ten arbitrary IDs.

## How to use the chapters in conversation

With this repository available to ChatGPT, ask it to start the course, continue a chapter, or work through a named topic. The tutor uses the separate conversation sessions with their question progression and coaching notes. If you refer to an example from the PDF, it can retrieve the reading chapter too. Diagrams should become short descriptions when you cannot see them. You do not need to copy lesson prompts or use a browser reader.

The tutor should teach a missing idea before testing it, give specific feedback, and vary an assumption after you understand the baseline. You can say “walk through the function structure,” “give me a counterexample,” “explain why the pointer moves,” or “turn this into a spoken coding exercise.” The written answers remain available if you are reading without an AI.

## What the book prioritizes

The first five chapters develop representation, hashing, windows, ordering, graph traversal, and reasoning about choices. The next four apply those tools to caches, expiration, concurrency, rate limits, and event services. Chapter ten develops recommendation and ranking exercises. Chapter eleven models larger programs. Chapter twelve contains two complete interview-style case studies with worked answers, also usable as oral mocks.

Your plan.md shaped the practical emphasis. The exercises are original study scenarios, not a verified list of Netflix questions. Sources establish language behavior, platform context, and recommendation concepts; they do not establish what your interviewer will ask. The aim is transferable reasoning across a wide variety of prompts.

There is no required completion time for a chapter. Independent reading is shorter than a conversation with explanations, attempts, corrections, and follow-ups. Progress means that a new example makes sense and that you can defend your decisions, not that a timer has expired.
