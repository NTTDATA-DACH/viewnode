---
name: gvb-prd-to-ralph
description: Convert a repository-aligned PRD in tasks/prd-*.md into a Ralph-ready scripts/ralph/prd.json, or refresh an existing Ralph spec from updated PRD requirements.
---

# Purpose

Convert a Markdown PRD into a Ralph-ready `scripts/ralph/prd.json`.
Prefer deterministic mapping from a well-structured PRD over creative reinterpretation.

This skill is intended to run after a PRD already exists, typically after:
- `$gvb-speckit-to-prd`

Use repository guidance from `AGENTS.md` for branch handling, validation, commit rules, naming, and project conventions. Do not restate stable repo rules in the JSON unless they materially affect Ralph execution.

# Inputs

Accept any of:
- a source PRD path, typically `tasks/prd-*.md`
- an existing `scripts/ralph/prd.json`, if present
- repository guidance from `AGENTS.md`

If the user does not provide an explicit PRD path, resolve the intended PRD from repository context when possible.

Do not require the PRD filename or feature name to exactly match the current branch name.
Treat branch context as a hint for issue and feature resolution, while allowing shorter artifact names when they remain unambiguous and traceable.

Treat the PRD as the planning source of truth for this skill.
Do not read back into upstream Speckit artifacts such as `spec.md`, `plan.md`, or `tasks.md` unless the user explicitly asks, or unless repository workflow requires verification against them.

When the PRD follows the expected structure from `gvb-speckit-to-prd`, map it directly:
- `## Overview` -> JSON `overview`
- `## User Stories` -> JSON `userStories`
- `## Traceability` -> issue-aware branch and story identifier preservation
- `## Validation` -> shared verification context that should be reflected in story acceptance criteria when relevant

When the PRD appears to come from the issue-backed workflow, treat a missing or malformed `## Traceability` block as suspicious input that must be surfaced explicitly before writing `scripts/ralph/prd.json`.

# Output

Write valid JSON to:

- `scripts/ralph/prd.json`

The result must be immediately usable by Ralph without manual cleanup.

# Mode selection

Choose the correct mode from task context:

## New generation
Use this when creating a new Ralph spec from a PRD.

## Update existing spec
Use this when `scripts/ralph/prd.json` already exists and the task is to revise or extend previously planned work.

In update mode:
- preserve correctly completed stories that remain valid
- preserve unchanged valid stories whenever possible
- prefer adding new follow-up stories for changed or additional requirements
- only reopen, replace, or rewrite completed stories if they are no longer valid under the revised PRD
- preserve existing `passes: true` only for unchanged stories that still satisfy the revised PRD

# Workflow

1. Read the source PRD carefully.
2. Read and apply relevant guidance from `AGENTS.md`.
3. Inspect existing `scripts/ralph/prd.json` if it exists.
4. Check whether the PRD appears to come from the issue-backed workflow, for example by the presence of GitHub issue references, issue-numbered feature naming, or explicit upstream workflow notes.
5. If issue-backed workflow evidence exists, verify that a usable `## Traceability` block is present before continuing.
6. If the traceability block is missing or malformed in an issue-backed PRD, stop and warn clearly instead of silently guessing branch or issue identity.
7. Determine `branchName` from repository guidance and task context.
8. Reuse the PRD's existing execution-oriented story slicing whenever it already satisfies Ralph's one-iteration rule.
9. Map each PRD story heading, description, and acceptance-criteria block into one JSON user story whenever the PRD structure is already compatible.
10. Convert the PRD into a small set of execution-ready user stories only when additional splitting or normalization is still required.
11. Order stories by dependency, execution priority, and change type.
12. Write explicit, minimal, testable acceptance criteria for each story.
13. Keep the plan aligned with existing repository patterns and conventions.
14. Write the final JSON to `scripts/ralph/prd.json`.

# Branch handling

Set `branchName` using repository guidance from `AGENTS.md`.

When already on a branch that clearly matches the PRD or issue context:
- reuse that branch if repository guidance supports doing so
- prefer the existing branch over deriving a new one

When the PRD includes GitHub-issue traceability, preserve that canonical issue identifier while following repository branch rules.
Otherwise derive the branch name from repository guidance and task context.

Do not invent a different branch when repository rules already define branch behavior.

Artifact naming does not need to exactly match the branch name:
- a shorter PRD filename may be used when it remains unambiguous and traceable
- `branchName` in `scripts/ralph/prd.json` should still reflect the actual branch Ralph is expected to use

# Overview mapping

When the PRD has a distinct `## Overview` section, use it as the primary source for JSON `overview`.
Only synthesize `overview` from other PRD content when no dedicated overview section exists.

Keep `overview` concise. It should summarize the execution goal, not duplicate all functional requirements.

# Story rules

Each story must:
- represent one coherent execution slice
- be feasible in one Ralph iteration
- be independently verifiable
- avoid bundling unrelated work

Prefer preserving PRD story boundaries when they already meet these rules.
Only split or merge PRD stories when the existing PRD structure is clearly too large, too vague, or dependency-incoherent for Ralph.

When the PRD uses stable story headings such as `### US-001: Title`, preserve both the story boundary and the title unless a change is needed for Ralph-sized execution.

Adapt story slicing to the change type:
- for behavior changes, slice by observable behavior and dependency order
- for refactors, slice by safe structural steps plus verification
- for test-related work, slice by coverage area, harness, or validation goal
- for documentation work, slice by document surface or audience-facing update
- for maintenance or internal tooling work, slice by operator outcome and verification

Choose execution order based on dependencies and change type.

For behavior changes, prefer this order when applicable:
1. wiring or entry-point changes
2. core behavior
3. validation and tests
4. documentation or polish

If the change is too large, split by execution slice rather than by broad technical layers.

# Acceptance criteria rules

Acceptance criteria must be:
- observable
- binary where possible
- minimal
- specific enough for Ralph to verify completion

When the PRD already includes suitable acceptance criteria, preserve them rather than rewriting them stylistically.

Include validation criteria when relevant, based on `AGENTS.md` and the PRD.
If the PRD includes a shared `## Validation` section, reflect the relevant repository-level validation expectations in the affected stories without copying unrelated checks into every story.
Prefer criteria that Ralph can verify directly in one iteration, such as tests, commands, generated artifacts, or specific observable behavior.

# Repository alignment

For this repository:
- follow existing project conventions and structure
- reuse established patterns
- use relevant hints from `AGENTS.md`

Do not assume the PRD came from a specific upstream workflow; treat the PRD as the source of truth regardless of whether it originated directly from an issue or from Speckit artifacts.

If the PRD includes traceability metadata from Speckit or issue-based workflows, preserve the intent of that metadata when useful, but do not let it override the actual PRD requirements.

When a fixed `## Traceability` block is present, prefer reading issue number, issue URL, and feature identity from that block instead of inferring them again from filenames or branch names.
For issue-backed PRDs, this preferred source becomes the required source unless the user explicitly asks to repair a broken PRD first.

# Output schema

Produce valid JSON with exactly these top-level fields:
- `branchName`
- `overview`
- `userStories`

Each user story must include:
- `id`
- `title`
- `description`
- `acceptanceCriteria`
- `priority`
- `passes`
- `notes`

# Story ID rules

The `id` must:
- reuse existing story IDs defined in the PRD when applicable
- otherwise preserve stable IDs already present in an existing `scripts/ralph/prd.json` when the story is unchanged
- otherwise create a new unique ID

When the PRD clearly carries GitHub-issue traceability from the upstream workflow, prefer:
- `#<issue>-<nnn>`

where:
- `<issue>` is the issue number
- `<nnn>` is a zero-padded sequential story number such as `001`, `002`

When no canonical issue number is available, derive IDs from repository-consistent task context.

If the PRD already defines stable story IDs and they remain compatible with the workflow, preserve them rather than replacing them.

# Pass state rules

Set `passes` to `false` for:
- newly generated stories
- newly added stories
- materially changed stories

Only preserve `passes: true` when the story is unchanged and still valid under the revised PRD.

# Completion bar

The final `scripts/ralph/prd.json` must:
- be ready for immediate Ralph execution
- contain multiple small stories when the work is non-trivial
- reflect relevant repository guidance from `AGENTS.md`
- avoid unnecessary duplication of stable repo policy inside the JSON
- be derived from explicit traceability data when the PRD came from the issue-backed workflow, not from silent re-inference
