# Schema first API initiative

[Story selector](../story-bank.md) · [Culture framing](../culture-framing.md) · [Gaps](../gaps.md)

**Use for:** Your main project for complexity, ownership, API design and developer impact. Use S03 for the adoption interaction and S04 for initiative; they are episodes from this project.

**Research questions:** Q03, Q04, Q06, Q13. **State:** Draft ready for rehearsal. No independent rehearsal recorded.

## Opening

> I made contract-first API development useful to developers by delivering validation, linting and type-generation tooling and helping teams adopt it.

## Answer structure

1. **Problem:** Incomplete contracts forced integrators to inspect code or ask owners; validation differed between teams.
2. **My scope:** Validation library, schema linter, type-generation tooling, adoption support and joint introduction of breaking-change detection.
3. **Decision:** OpenAPI 3.1 aligned with JSON Schema tooling, at the cost of upgrades from 2.0/3.0 and incompatible tools.
4. **Result:** Teams reported easier integration; architecture valued linting; SDK teams could generate SDKs; main-API contract review catches design issues before implementation.
5. **Reflection:** Explore framework wrappers and integration with the standard CI pipeline.

## Frame for the question

Lead with a concrete developer problem, then explain the version tradeoff. Connect the tooling to easier integration and earlier feedback. Do not lead with a catalogue of libraries.

**Culture connection:** Judgment is visible in the compatibility tradeoff. For People Over Process, explain the practical benefit of checks and contract reviews; do not equate adding governance with success. See the [memo-based framing guide](../culture-framing.md) for the source and interpretation.

## Spoken draft

> One project I enjoyed was moving our APIs from code-first to contract-first development. Documentation was inconsistent, and developers often had to read another service's code or ask its owning team to understand the contract. Request and response validation also varied across teams.

> My responsibility was to make the switch practical and valuable for developers. I shipped a generic OpenAPI-based validation library, built a schema linter using an architect's API guidelines, and worked on type-generation tooling. I also helped teams adopt the approach and introduced breaking-change detection together with the main API team. I initiated the linting work myself.

> A difficult choice was targeting OpenAPI 3.1 for its JSON Schema alignment, which allowed us to use existing validators. The cost was migrating APIs from 2.0 and 3.0 and dealing with incompatible tooling. During the transition, my teammate wrote a converter that upgraded the old 2.0 schema to 3.1 and merged in existing schema-first endpoints, supporting coexistence of both approaches.

> Teams told us integrations became easier and quicker, and they valued automatic validation. The architecture team appreciated the standardization enabled by linting. My work also enabled the SDK teams to introduce generated SDKs. All contributions to the main API now go through contract review, where design issues are caught before implementation.

> Looking back, I'd explore framework-specific libraries on top of the generic core to reduce integration boilerplate. I'd also work with the CI team to explore incorporating the tools into our standard delivery pipeline; currently, they run through separate GitHub Actions. Those would be my next improvements to the developer experience.

## Keep the claims precise

- Architect authored the guidelines; teammate authored the converter; SDK teams delivered generated SDKs.
- Contract testing remains prospective. Performance requirements are not measured results. You have not claimed ownership of the governance process.

## Follow-up preparation

- Which validator and framework integration did you actually ship, and how did you assess overhead?
- What part of the version choice did you personally own, and what alternative did you consider?
