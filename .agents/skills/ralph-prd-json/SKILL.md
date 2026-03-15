---
name: ralph-prd-json
description: Convert a repository-aligned PRD in tasks/prd-*.md into a Ralph-ready scripts/ralph/prd.json, or refresh an existing Ralph spec from updated PRD requirements.
---

# Purpose

Convert a Markdown PRD into a Ralph-ready `scripts/ralph/prd.json`.

This skill is intended to run after a PRD already exists, typically after:
- `$spec-to-prd`

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
Do not read back into upstream Spec Kit artifacts such as `spec.md`, `plan.md`, or `tasks.md` unless the user explicitly asks, or unless repository workflow requires verification against them.

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
4. Determine `branchName` from repository guidance and task context.
5. Convert the PRD into a small set of execution-ready user stories.
6. Order stories by dependency, execution priority, and change type.
7. Write explicit, minimal, testable acceptance criteria for each story.
8. Keep the plan aligned with existing repository patterns and conventions.
9. Write the final JSON to `scripts/ralph/prd.json`.

# Branch handling

Set `branchName` using repository guidance from `AGENTS.md`.

When already on a branch that clearly matches the PRD or issue context:
- reuse that branch if repository guidance supports doing so
- prefer the existing branch over deriving a new one

When the PRD or task maps to a GitHub issue:
- follow the repository's issue-based branch rules
- use an existing issue branch if repository guidance requires checking for one
- otherwise derive the branch name according to the repository rule

When the PRD does not map to a GitHub issue:
- derive the branch name from repository guidance and task context

Do not invent a different branch when repository rules already define branch behavior.

Artifact naming does not need to exactly match the branch name:
- a shorter PRD filename may be used when it remains unambiguous and traceable
- `branchName` in `scripts/ralph/prd.json` should still reflect the actual branch Ralph is expected to use

# Story rules

Each story must:
- represent one coherent execution slice
- be feasible in one Ralph iteration
- be independently verifiable
- avoid bundling unrelated work

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

Include validation criteria when relevant, based on `AGENTS.md` and the PRD.

# Repository alignment

For this repository:
- follow existing project conventions and structure
- reuse established patterns
- use relevant hints from `AGENTS.md`

Do not assume the PRD came from a specific upstream workflow; treat the PRD as the source of truth regardless of whether it originated directly from an issue or from Spec Kit artifacts.

If the PRD includes traceability metadata from Spec Kit or issue-based workflows, preserve the intent of that metadata when useful, but do not let it override the actual PRD requirements.

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

When the PRD clearly maps to a GitHub issue, prefer:
- `#<issue>-<nnn>`

where:
- `<issue>` is the issue number
- `<nnn>` is a zero-padded sequential story number such as `001`, `002`

When no issue number is available, derive IDs from repository-consistent task context.

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