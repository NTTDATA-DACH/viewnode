---
name: issue-to-spec
description: Create or update a Spec Kit spec from a GitHub issue and save it under specs/<feature>/spec.md. Use when a GitHub issue should become a structured, repository-aligned specification for later clarification, planning, and task generation.
---

Read the GitHub issue provided by the user and create or update a Spec Kit specification.

# Purpose

Turn a GitHub issue into a clean, structured spec that becomes the source of truth for later clarify, plan, tasks, and PRD generation.

# Inputs

Accept any of:
- a GitHub issue URL
- an issue number
- issue text copied by the user

Prefer the issue URL or issue number when available.

# Default behavior

- Use current branch context to identify the active issue or feature when helpful.
- Do not require the feature name to exactly match the branch name.
- Prefer a shorter normalized feature name when it remains unambiguous and traceable.
- When possible, preserve the issue number in the feature name.
- If a matching feature folder already exists under `specs/`, reuse it.
- Otherwise derive the feature name from the issue number and a short normalized issue-title slug.
- Follow the repository's existing feature naming conventions when they are clear.
- Write to `specs/<feature>/spec.md` unless the user explicitly requests another path.
- If the feature folder does not exist, create it.
- If a matching spec already exists, update it carefully instead of rewriting it blindly.

# Requirements

- Apply relevant repository guidance from `AGENTS.md`.
- Follow existing project conventions where they affect observable behavior, compatibility, workflow, or user expectations.
- Keep the spec outcome-oriented and avoid unnecessary implementation detail.
- Prefer explicit statements from the issue; do not invent requirements.
- If the issue is ambiguous or incomplete, capture uncertainty in an `Open Questions` or `Assumptions` section instead of guessing.
- Keep issue-specific details in the spec, not in `AGENTS.md` and not in this skill.
- When the issue references an existing pattern, command, API, workflow, feature, or document, capture that reference explicitly as part of the requirements.

# Adapt to the issue type

- For user-facing features, focus on behavior, scope, and acceptance criteria.
- For refactors, focus on motivation, invariants, constraints, risks, and proof that behavior remains correct.
- For test-related issues, focus on coverage goals, failure modes, validation scope, and expected confidence improvements.
- For documentation issues, focus on audience, affected docs surfaces, missing or incorrect information, and completion criteria.
- For internal tooling or maintenance issues, focus on operator or maintainer outcomes, constraints, and verification.

# Include as applicable

- title
- problem statement
- goal
- user story, maintainer story, or engineering story
- scope
- non-goals
- acceptance criteria
- constraints or compatibility expectations
- assumptions or open questions when needed

# Guidance

- Focus on what should change and why.
- Avoid detailed implementation design unless the issue explicitly requires it.
- Write acceptance criteria so they can support later planning, task generation, and PRD conversion.
- Prefer stable identifiers for major acceptance criteria when useful, for example `AC-001`, `AC-002`.

# Output

- Write the spec to the resolved path, typically `specs/<feature>/spec.md`.
- Return the resolved feature name and output path.
- Keep the resolved feature name short, traceable, and repository-aligned; it does not need to equal the current branch name.