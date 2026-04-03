# AGENTS.md

## Purpose
- This file defines shared agent guidance for repositories that install tooling from this project.
- It is intentionally generic and should be safe to reuse across multiple repositories.
- Consumer repositories may replace or extend this file with stricter local rules when needed.

## Ralph PRD generation
- When generating or updating `scripts/ralph/prd.json`, split work into multiple small user stories.
- Each story must be feasible in one Ralph iteration.
- Prioritize stories by dependency and execution order.
- Keep each story narrowly scoped and independently verifiable.
- Write precise, minimal, testable acceptance criteria.
- Reuse existing project patterns and avoid introducing parallel structures.
- If `scripts/ralph/prd.json` already exists and the task is an update, prefer adding follow-up stories instead of rewriting completed stories, unless they are no longer valid.

## Branch and repository handling
- For issue-based work, first check whether a matching local or remote branch already exists and reuse it when appropriate.
- If the issue already has a linked or existing branch, use that exact branch name.
- Do not invent a different branch name when an issue branch already exists.
- Preserve the repository's existing branch and feature numbering format, including zero padding when the repository uses it.
- When a generator or repository script returns the branch or feature name, treat that returned name as the source of truth instead of reformatting the numeric prefix by hand.
- Do not add extra prefixes unless the user explicitly asks or the repository explicitly requires them.
- Do not create or switch to a different feature branch unless the user explicitly asks.

## Commit rules
- Commit messages must follow Conventional Commits format.
- Add a scope in parentheses after the type when a clear scope exists.
- Reference the issue in the subject when applicable.
- Use lowercase by default, except where capitalization is required for correct names.
- Prefer small commits grouped by purpose.
- Do not use vague commit messages such as `update`, `fix stuff`, or `changes`.

## Validation
- Run the closest relevant automated checks before committing when the repository provides them.
- Prefer targeted validation near the changed area first, then broader repository validation when appropriate.
- If validation cannot be run, call that out explicitly.

## Project conventions
- Prefer existing repository patterns over introducing new structural styles.
- Preserve externally observable behavior unless the task explicitly asks for behavior changes.
- Favor incremental refactors with verification over broad rewrites.
- When changing generated or generated-adjacent artifacts, update the source and regenerate rather than editing derived output by hand when the repository provides a generation path.

## Testing conventions
- Add or update tests alongside refactoring and bug fixes.
- Prefer tests close to the changed package, module, or feature area.
- For refactors, ensure tests verify preserved behavior, not just new internal structure.

## Documentation conventions
- User-facing documentation and examples should stay in sync with behavior changes.
- When changing user-facing commands, APIs, or workflows, update the relevant documentation in the same change.
- If documentation is generated from source metadata, update the source and regenerate rather than hand-editing derived output.

## Technology baseline
- Follow the repository's current toolchain, dependency, and framework conventions.
- Prefer the libraries and frameworks already established in the repository unless the user explicitly asks for a change.

## Issue-specific guidance
- Issue-specific requirements belong in the issue, feature artifacts, and PRD.
- Do not add changing issue-specific details to this file unless they become stable reusable repository rules.
