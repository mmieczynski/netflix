# Netflix L5 interview learning material

A reading edition for learning away from a computer, and a separate conversation
edition for ChatGPT using this repository. It covers the live technical
interview in Go and the subsequent L5 Developer Platform manager conversation.

Read the [PDF](output/pdf/netflix-go-interview-book.pdf), the
[book introduction](learning-material/BOOK-INTRO.md), or the
[reading chapters](learning-material/reading/). The revised book introduces the problem and expected output before explaining
an approach. It uses literal arrays and maps, 43 state tables, 29 vector diagrams,
and 65 Go code blocks across the 12 technical chapters. The code appears beside the
concept it implements, with assumptions, a trace, an invariant, and boundary
cases. Both final case studies include assembled implementations.
The [conversation sessions](learning-material/sessions/) retain their questions
and coaching notes. Neither route requires an editor, code execution, or timed work.

## Use with ChatGPT

With this repository available in the conversation, ask to start the course,
continue a chapter, or work through a named topic. ChatGPT should begin with this
README, then read [TUTOR.md](learning-material/TUTOR.md), the
[chapter index](learning-material/session-index.json), and the relevant chapter.
There is no prompt-copying or HTML-reader workflow.

The tutor should teach before quizzing, ask one question at a time, and adapt to
your answers. Coding exercises can be done entirely aloud: define the API,
structs and fields, helper functions, control flow, mutation order, and test
cases. Writing Go syntax is a separate skill, not a prerequisite for this book.

The repository's [progress file](learning-material/study-progress.json) starts
empty. An agent with write access can update it. If a ChatGPT connection can only
read GitHub, the tutor should keep a concise progress summary in the conversation
and not claim to have saved a file. A README or AGENTS.md is a retrieval entry
point; its presence alone does not prove every remote chat will automatically
load it. The tutor must establish access by actually reading the chapter.

## Contents

1. Go and data structure choices
2. Hashing, arrays, and sliding windows
3. Sorting, binary search, intervals, and heaps
4. Trees, graphs, and dependencies
5. Stacks, greedy choices, backtracking, and dynamic programming
6. LRU and weighted caches
7. Expiration, clocks, and concurrency
8. Rate limiting and rolling windows
9. Deduplication, counters, versions, and scheduling
10. Recommendations and streaming ranking
11. Filesystems, transactions, pub/sub, and service composition
12. Complete cache and viewing-statistics interview case studies
13. Hiring manager interview preparation for Developer Platform

[Sources and scope](learning-material/SOURCES.md) distinguish official technical
documentation from original exercises. These are not confirmed Netflix questions.

## Hiring manager preparation

All manager-interview materials live in [hiring-manager/](hiring-manager/).
The [research guide](hiring-manager/reading/13-hiring-manager.md) is tailored
to the recruiter brief and [Developer Platform position](https://explore.jobs.netflix.net/careers/job/790301977513),
with sources checked on 13 October 2026. It covers 18 question families, strong
answer ingredients, likely follow-ups, interpretation risks, and questions to
ask the manager. It is included as chapter 13 in the PDF and plain-text edition.

The [story selector](hiring-manager/story-bank.md) maps interview questions to
individual story cards. [Culture framing](hiring-manager/culture-framing.md) and
[remaining gaps](hiring-manager/gaps.md) guide preparation. [Manager preparation progress](hiring-manager/study-progress.json)
is stored separately from technical study. Ask to prepare a project story or practise the
manager conversation; ChatGPT should retrieve
[session 13](hiring-manager/sessions/13-hiring-manager.md) and develop it with
you one question at a time. No stories or completed practice have been inferred
from the authoring of these materials.

## Training projects

[Cache with a time limit](training/cache-with-time-limit/README.md) is an
independent Go exercise with a CodeSignal-style filesystem layout: a problem
statement, method contract, complete reference implementation, and deterministic tests.
It supports both optional coding practice and discussing an implementation aloud.

## Files for maintaining the material

- `learning-material/reading/*.md`: canonical reading chapters, used for the PDF.
- `learning-material/reading/figures.json`: editable vector-diagram specifications.
- `learning-material/reading/figures/*.svg`: generated diagrams visible in Markdown/GitHub.
- `learning-material/reading-index.json`: ordered reading chapters.
- `learning-material/sessions/*.md`: conversation lessons and question bank, read by the tutor.
- `learning-material/TUTOR.md`: teaching and spoken-implementation guidance.
- `learning-material/study-progress.json`: learner state, only updated from actual learning.
- `hiring-manager/`: manager research, coaching session, story bank and preparation progress.
- `learning-material/book.txt`: complete plain-text edition.
- `output/pdf/netflix-go-interview-book.pdf`: portable reading edition.
- `learning-material/build_pdf.py`: regenerates the PDF and text edition.
- `learning-material/sync_reading_code.py`: extracts every reading Go block for compiler checks.
- `learning-material/readingcode/`: generated examples and their behavior tests.
- `learning-material/*.go`: optional reference code; the book stands without it.

To regenerate the editions, install the dependencies in
`learning-material/requirements-pdf.txt`, then run
`python learning-material/build_pdf.py`. This is a maintainer operation, not a
learning exercise. `plan.md` remains the original scope input.

The reading and conversation editions are intentionally separate sources. When
correcting a shared concept, inspect both editions; when changing presentation,
tailor it to that edition. Do not regenerate the book by flattening the tutor's Q&A.
