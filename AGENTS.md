# Repository instructions

This repository contains the user's Netflix L5 Go interview preparation book.
Reading chapters live in `learning-material/reading/`; conversation lessons and
coaching questions live in `learning-material/sessions/`. Keep their presentation
distinct. The PDF build uses `reading-index.json`, not the conversation index.
It must work independently on an e-reader and in ChatGPT conversation mode with
access to the remote GitHub repository. Do not ask the user to copy prompts,
upload repository files, open an HTML reader, use an editor, or execute code.

## Learning requests

When the user asks to start, continue, study, practise verbally, review a topic,
or conduct a mock interview:

1. Read `learning-material/TUTOR.md`.
2. Read `learning-material/study-progress.json`.
3. Read the relevant chapter in `learning-material/sessions/`, using
   `learning-material/session-index.json` to locate it. The complete text and
   coach explanations are local; retrieve them yourself.
4. Teach in conversation. Ask one question, wait for the user's answer, then
   give feedback or a smaller explanation before moving on. Follow the tutor
   guide's adaptive sequence and progress rules.

Default to spoken-friendly explanations and exercises requiring no editor.
Treat typed answers the same way. Coding exercises belong in conversation:
discuss structs, helpers, branches, mutations, pseudocode, and mental tests.
Only writing syntax and executing it are optional. Do not add timed-practice
sections or count editor work toward reading/listening time.
The user's live interview is with another senior engineer, at Netflix L5,
using Go through CodeSignal. Do not assume an automated assessment format.

Do not interrupt unrelated repository work with a lesson. Course authoring or
maintenance requests are requests to update the materials, not evidence that
the user has completed a lesson. Never mark authored material as learned.

If the remote connection cannot write files, keep progress in the conversation
and do not claim persistence to GitHub. Do not assume access or automatic
instruction loading; retrieve the README and lesson through available tools.

## Current technical documentation

Use Context7 MCP for current library, framework, SDK, API, CLI-tool, or cloud
service documentation when a user question involves their syntax, configuration,
migration, setup, or library-specific debugging. Begin with resolve-library-id
unless an exact library ID is supplied, then query-docs for each distinct concept.
Prefer exact relevant matches and reputable sources. Do not use Context7 merely
for general algorithm concepts, business logic, refactoring, or code review.
These documentation instructions preserve the user's provided AGENTS.md rules.
