---
name: spec-to-prd
description: Create or refresh a PRD from Spec Kit artifacts and save it under tasks/prd-*.md. Use when a Spec Kit feature exists under specs/<feature>/ and a repository-aligned execution-oriented PRD is needed.
---

Read the relevant Spec Kit artifacts and create or refresh a PRD markdown file for this repository.

# Purpose

Convert a Spec Kit feature into an execution-oriented PRD suitable for later Ralph conversion.

This skill is intended to run after Spec Kit has produced one or more upstream artifacts, typically through:
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

Do not treat the outputs of other Spec Kit skills as primary PRD sources unless the user explicitly asks:
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
- When possible, preserve the issue number in the feature name.
- If a matching feature folder already exists under `specs/`, reuse it.
- Otherwise resolve the feature from the issue number, feature id, slug, or repository context.
- Write to `tasks/prd-<feature>.md` unless the user explicitly requests another path.
- If a matching PRD already exists, update it carefully instead of rewriting it blindly.
- Preserve still-valid content unless the Spec Kit artifacts clearly supersede it.

# Requirements

- Apply relevant repository guidance from `AGENTS.md`.
- Keep the PRD implementation-aware but not implementation-heavy.
- Preserve the intent of the spec.
- Pull technical constraints, integration points, validation notes, sequencing, and work shaping from `plan.md` and `tasks.md` when available.
- If important ambiguity remains, include an explicit `Assumptions / Open Questions` section rather than inventing facts.
- When the spec references an existing pattern, command, API, workflow, feature, or document, preserve that reference as part of the acceptance contract.

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

# Traceability metadata

Include when useful:
- feature name
- spec path
- plan path if present
- tasks path if present
- source status such as `derived from Spec Kit artifacts`

# Guidance

- Use `spec.md` to define what and why.
- Use `plan.md` for constraints, architecture touchpoints, validation strategy, migration concerns, and technical risks.
- Use `tasks.md` to shape work packages and later Ralph conversion.
- Prefer stable identifiers for acceptance criteria and major scope items when useful.
- Do not let task wording distort the original intent from the spec.
- If only `spec.md` exists, derive the PRD from the spec alone and keep missing planning detail explicit rather than guessing.
- If `plan.md` exists but `tasks.md` does not, use the plan to enrich the PRD without pretending that task breakdown already exists.

# Output

- Write the PRD to the resolved path, typically `tasks/prd-<feature>.md`.
- Return the resolved feature name, source paths used, and output path.
- Be explicit about which Spec Kit artifacts were available and which were absent.