---
name: ralph-prd-json
description: Convert a feature PRD in tasks/prd-*.md into a Ralph-ready scripts/ralph/prd.json, or update an existing Ralph spec from revised PRD requirements.
---

# Purpose

Convert a Markdown PRD into a Ralph-ready `scripts/ralph/prd.json`.

Use repository guidance from `AGENTS.md` for branch handling, validation, commit rules, and project conventions. Do not restate stable repo rules in the JSON unless they materially affect Ralph execution.

# Inputs

- Source PRD in `tasks/prd-*.md`
- `AGENTS.md`
- Existing `scripts/ralph/prd.json`, if present

# Output

Write valid JSON to:

- `scripts/ralph/prd.json`

The result must be immediately usable by Ralph without manual cleanup.

# Mode selection

Choose the correct mode from the task context:

## New generation
Use this when creating a new Ralph spec from a PRD.

## Update existing spec
Use this when `scripts/ralph/prd.json` already exists and the task is to revise or extend previously planned work.

In update mode:
- preserve correctly completed stories
- prefer adding new follow-up stories for changed requirements
- only reopen or rewrite completed stories if they are no longer valid

# Workflow

1. Read the source PRD carefully.
2. Read and apply relevant guidance from `AGENTS.md`.
3. Inspect existing `scripts/ralph/prd.json` if it exists.
4. Determine `branchName` from repository guidance and task context.
5. Convert the PRD into a small set of execution-ready user stories.
6. Order stories by dependency and execution priority.
7. Write explicit, minimal, testable acceptance criteria for each story.
8. Keep the plan aligned with existing repository patterns and command conventions.
9. Write the final JSON to `scripts/ralph/prd.json`.

# Branch handling

Set `branchName` using repository guidance from `AGENTS.md`.

When the PRD or task maps to a GitHub issue:
- follow the issue-based branch rules from `AGENTS.md`
- use an existing issue branch if repository guidance requires checking for one
- otherwise derive the branch name according to the repository rule

Do not invent a different branch when the repository rules already define branch behavior.

# Story rules

Each story must:
- represent one coherent execution slice
- be feasible in one Ralph iteration
- be independently verifiable
- avoid bundling unrelated work

Prefer this execution order when applicable:
1. command or API wiring
2. core behavior
3. tests
4. documentation or polish

If the feature is too large, split by execution slice rather than by broad technical layers.

# Acceptance criteria rules

Acceptance criteria must be:
- observable
- binary where possible
- minimal
- specific enough for Ralph to verify completion

Include validation criteria when relevant, based on `AGENTS.md`.

# Repository alignment

For this repository:
- follow existing command conventions and structure
- reuse established patterns
- use relevant hints from `AGENTS.md`

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

The `id` must:
- reuse existing ID defines in tasks/prd-*.md when applicable, or be a new unique ID if the story is new.
- in case of GitHub issue like #47, the `id` should be `#47-<nnn>`, where `<nnn>` is a zero-padded sequential number (e.g., `001`, `002`) representing the story's order within the PRD.

Set `passes` to `false` for newly generated or newly added work unless the user explicitly asks to preserve prior completion state.

# Completion bar

The final `scripts/ralph/prd.json` must:
- be ready for immediate Ralph execution
- contain multiple small stories when the feature is non-trivial
- reflect relevant repository guidance from `AGENTS.md`
- avoid unnecessary duplication of stable repo policy inside the JSON