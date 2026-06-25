---
name: gvb-speckit-to-prd
description: Create or refresh a PRD from Speckit artifacts and save it under tasks/prd-*.md when a repository-aligned execution-oriented PRD is needed.
---

Read the relevant Speckit artifacts and create or refresh a PRD markdown file for this repository.

# Purpose

Convert a Speckit feature into an execution-oriented PRD suitable for later Ralph conversion.
Treat the PRD as a stable Ralph-ready intermediate artifact, not as a generic product brief.

This skill is intended to run after Speckit has produced one or more upstream artifacts, typically through:
- `$speckit-specify`
- optionally `$speckit-clarify`
- `$speckit-plan`
- `$speckit-tasks`

# Inputs

Accept any of:
- a feature folder path
- a feature id or short slug
- a direct path to `spec.md`

If the user provides no explicit input, infer the active feature from repository context when possible, such as:
- the current branch
- an existing `SPECIFY_FEATURE` value
- a single obvious matching folder under `specs/`

Treat branch context as a hint for resolving the feature, not as a requirement that the feature name exactly match the branch name.

# Source priority

Use the best available source in this order, when present:
1. `specs/<feature>/tasks.md`
2. `specs/<feature>/plan.md`
3. `specs/<feature>/spec.md`

Do not fail merely because `tasks.md` or `plan.md` is missing if a valid `spec.md` exists.

Primary PRD sources are:
- `spec.md`
- `plan.md`
- `tasks.md`

Do not treat the outputs of other Speckit skills as primary PRD sources unless the user explicitly asks:
- `speckit-analyze`
- `speckit-checklist`
- `speckit-implement`
- `speckit-taskstoissues`

Those may inform validation or follow-up work, but the PRD should still be derived mainly from the feature spec, plan, and tasks.

# Default behavior

- Resolve the feature from the provided input or current feature context.
- Use current branch context to identify the active issue or feature when helpful.
- Do not require the feature name to exactly match the branch name.
- Prefer a shorter normalized feature name when it remains unambiguous and traceable.
- When the feature originated from the GitHub-issue workflow, preserve GitHub issue traceability while following the repository's established feature and filename formatting.
- If a matching feature folder already exists under `specs/`, reuse it.
- Otherwise resolve the feature from the issue number, feature id, slug, or repository context.
- Write to `tasks/prd-<feature>.md` unless the user explicitly requests another path.
- If a matching PRD already exists, update it carefully instead of rewriting it blindly.
- Preserve still-valid content unless the Speckit artifacts clearly supersede it.

# Requirements

- Apply relevant repository guidance from `AGENTS.md`.
- Keep the PRD implementation-aware but not implementation-heavy.
- Preserve the intent of the spec.
- Pull technical constraints, integration points, validation notes, sequencing, and work shaping from `plan.md` and `tasks.md` when available.
- If important ambiguity remains, include an explicit `Assumptions / Open Questions` section rather than inventing facts.
- When the spec references an existing pattern, command, API, workflow, feature, or document, preserve that reference as part of the acceptance contract.
- Make the resulting PRD structurally consistent so `gvb-prd-to-ralph` can map it with minimal interpretation.

# Adapt to the change type

- For features, emphasize behavior, scope, acceptance criteria, and relevant implementation notes.
- For refactors, emphasize current pain points, target structure, invariants, migration constraints, regression risks, and validation.
- For test-related work, emphasize coverage targets, scenarios, failure modes, and verification.
- For documentation work, emphasize affected docs, target audience, missing or outdated content, examples to add or fix, and review expectations.
- For internal tooling or maintenance work, emphasize operator or maintainer outcomes, constraints, risks, and proof of completion.

# Include as applicable

- problem statement
- goal
- user story, maintainer story, or engineering story
- scope
- non-goals
- acceptance criteria
- implementation notes relevant for later Ralph conversion
- assumptions / open questions when needed
- optional traceability metadata for source feature and source files

# Required PRD structure

Unless the user explicitly requests otherwise, always produce the PRD with these sections in this order:

1. `# PRD: <title>`
2. `## Overview`
3. `## Goals`
4. `## User Stories`
5. `## Functional Requirements`
6. `## Non-Goals`
7. `## Implementation Notes`
8. `## Validation`
9. `## Traceability`
10. `## Assumptions / Open Questions` when needed

The `## Overview` section must be concise and must be suitable for direct reuse as Ralph `overview`.

Within `## User Stories`, use a stable per-story structure whenever possible:

```md
### US-001: <story title>
**Description:** <one focused execution slice stated as an outcome>

**Acceptance Criteria:**
- <binary, observable criterion>
- <binary, observable criterion>
```

This structure is preferred because `gvb-prd-to-ralph` can map it directly with minimal interpretation.

# User story quality bar

Every PRD story must already be shaped for later Ralph conversion:

- each story must represent one coherent execution slice
- each story must be feasible in one Ralph iteration
- each story must be independently verifiable
- each story must be ordered by dependency
- each story must include explicit, binary acceptance criteria where possible
- each story should read as work Ralph can execute without needing to rediscover hidden context from upstream artifacts

For stories derived from broad Speckit user stories or many tasks, split them into smaller execution-oriented PRD stories instead of preserving oversized planning bundles.

When useful, include stable PRD story identifiers such as `US-001`, `US-002`, etc.

Prefer acceptance criteria that already include the repository-relevant validation expectation when the work clearly implies it, such as:
- the required test or verification command from `AGENTS.md`
- explicit regression checks from `plan.md`
- visual verification for UI changes when the workflow or upstream PRD guidance calls for it

# Traceability metadata

Traceability is mandatory when the workflow originated from `gvb-ghissue-to-speckit`.
Do not omit the `## Traceability` section in that case.

Include when useful, and include all of the following when available:
- GitHub issue number
- GitHub issue URL
- GitHub issue title
- feature name
- feature directory
- spec path
- plan path if present
- tasks path if present
- PRD output path
- source status such as `derived from Speckit artifacts`

Prefer a fixed `## Traceability` block format so downstream skills do not need to guess:

```md
## Traceability

- GitHub Issue: #<number>
- GitHub URL: <url>
- GitHub Title: <title>
- Feature Name: <feature>
- Feature Directory: <path>
- Spec Path: <path>
- Plan Path: <path or n/a>
- Tasks Path: <path or n/a>
- PRD Path: <path>
- Source Status: derived from Speckit artifacts
```

If the workflow is issue-backed and any required traceability item is unavailable, keep the `## Traceability` section and mark the missing value explicitly as `unknown` rather than dropping the field.

# Guidance

- Use `spec.md` to define what and why.
- Use `plan.md` for constraints, architecture touchpoints, validation strategy, migration concerns, and technical risks.
- Use `tasks.md` to shape work packages, dependency ordering, and later Ralph conversion.
- Prefer stable identifiers for acceptance criteria and major scope items when useful.
- Do not let task wording distort the original intent from the spec.
- If only `spec.md` exists, derive the PRD from the spec alone and keep missing planning detail explicit rather than guessing.
- If `plan.md` exists but `tasks.md` does not, use the plan to enrich the PRD without pretending that task breakdown already exists.
- If `tasks.md` exists, prefer its dependency ordering and execution slicing when shaping PRD stories, but rewrite tasks into coherent PRD stories rather than copying checklist items verbatim.
- Use `## Validation` to consolidate repository-level proof of completion that applies across stories, and keep story-level acceptance criteria focused on what Ralph must complete and verify within that slice.

# Output

- Write the PRD to the resolved path, typically `tasks/prd-<feature>.md`.
- Return the resolved feature name, source paths used, and output path.
- Be explicit about which Speckit artifacts were available and which were absent.
- Be explicit about whether the PRD is ready for deterministic Ralph conversion or whether any remaining ambiguity was carried into `Assumptions / Open Questions`.
- If the workflow is issue-backed, explicitly confirm that the `## Traceability` block was written and populated.
