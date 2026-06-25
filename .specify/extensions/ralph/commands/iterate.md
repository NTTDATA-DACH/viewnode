---
description: "Execute a single Ralph loop iteration - complete one work unit from tasks.md with proper commits and progress tracking"
---

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## Scope Constraint

**CRITICAL**: Complete AT MOST ONE work unit in this iteration.

- A work unit is the first incomplete user story or, if no user story tasks remain, the first incomplete final polish/cross-cutting section
- If you cannot complete an entire work unit, complete as many tasks as you can
- Partial progress is fine -- uncompleted tasks will be handled in subsequent iterations
- DO NOT start a second work unit even if you have time remaining
- This prevents context rot and keeps changes reviewable

## Outline

1. **Setup**: Run the prerequisite check script from repo root and parse FEATURE_DIR and AVAILABLE_DOCS list. All paths must be absolute. For single quotes in args like "I'm Groot", use escape syntax appropriate to your shell.

   ```bash
   .specify/scripts/bash/check-prerequisites.sh --json --require-tasks --include-tasks
   ```

   ```powershell
   .specify/scripts/powershell/check-prerequisites.ps1 -Json -RequireTasks -IncludeTasks
   ```

2. **Read context first**:
   - Read `FEATURE_DIR/progress.md` if it exists -- check the `## Codebase Patterns` section for discovered conventions
   - Read `FEATURE_DIR/tasks.md` -- understand task structure and identify next incomplete user story
   - Read `FEATURE_DIR/plan.md` for tech stack, architecture, and file structure
   - Read `FEATURE_DIR/spec.md` when present and treat any explicit traceability metadata there as authoritative workflow context
   - **IF EXISTS**: Read `FEATURE_DIR/data-model.md` for entities and relationships
   - **IF EXISTS**: Read `FEATURE_DIR/contracts/` for API specifications
   - **IF EXISTS**: Read `FEATURE_DIR/research.md` for technical decisions and constraints

3. **Identify scope**:
   - Find the FIRST work unit with incomplete tasks (`- [ ]`)
   - Prefer the first incomplete user story section
   - If no user story tasks remain, use the first incomplete final polish/cross-cutting section
   - Work ONLY on tasks within that single work unit
   - Example: If "US-001: Initialize Ralph Command" has incomplete tasks, work only on US-001

4. **Implement tasks**:
   - Complete tasks in dependency order (non-[P] before parallel [P] where noted)
   - Follow TDD when appropriate: write tests first, then implementation
   - Run quality checks after each task (typecheck, lint, test as appropriate)
   - Mark each completed task by changing `[ ]` to `[x]` in tasks.md

5. **Commit on work unit completion**:
   - When ALL tasks in the current work unit are complete (`[x]`), create exactly ONE commit for that work unit, even if the final changes are only task tracking, progress logging, documentation, or validation updates.
   - Stage `tasks.md` and `progress.md` as part of that SAME commit. Do not create a follow-up bookkeeping commit just to update iteration metadata.
   - Do not create a second "record iteration" or "record progress" commit after the main work-unit commit.
   - If `progress.md` would otherwise need the exact commit hash, prefer stable wording like `Recorded in Git history for this iteration` instead of creating another commit or amending only to backfill a hash.
   - Use a proper Conventional Commit scope that describes the changed area or subsystem (for example `view`, `cli`, `spec`, `docs`).
   - Do **NOT** use the feature directory name, feature ID, or git branch name as the commit scope.
   - Before creating the commit, inspect the active feature artifacts for explicit issue traceability metadata.
   - Prefer a stable, persisted source such as `FEATURE_DIR/spec.md`, `FEATURE_DIR/plan.md`, or another feature artifact over chat history or branch-name guesses.
   - If the active feature artifacts include an issue number or equivalent external work-item identifier, append `#<identifier>` as the final token of the commit subject line.
   - Preserve the normal Conventional Commit structure and add the traceability suffix only at the end of the subject.
   - If no explicit traceability identifier is persisted in the feature artifacts, use the normal Conventional Commit subject without any suffix.

     ```sh
     git add -A
     git commit -m "<type>(<scope>): <work unit title> #<identifier>"
     ```

   - Example with persisted traceability: `git commit -m "feat(view): group all-namespaces node output #59"`
   - Example with persisted traceability: `git commit -m "chore(spec): complete final polish #59"`
   - Example without persisted traceability: `git commit -m "feat(view): group all-namespaces node output"`
   - If only partial progress, NO commit -- let the next iteration continue
   - Before outputting `<promise>COMPLETE</promise>`, ensure there are no intended tracked changes left uncommitted for the just-finished work unit

6. **Update progress log**:
   - Create or append to `FEATURE_DIR/progress.md`
   - Add any discovered patterns to `## Codebase Patterns` section at TOP of file
   - Update `progress.md` BEFORE the main work-unit commit when the work unit is complete, so it is included in that same commit
   - Do not amend or create an extra commit only to replace placeholder commit text
   - Use the Progress Report Format below

## Progress Report Format

APPEND to FEATURE_DIR/progress.md:

```markdown
---
## Iteration [N] - [Current Date/Time]
**User Story**: [US-XXX title, final polish section title, or "Partial progress on ..."]
**Tasks Completed**: 
- [x] Task ID: description
- [x] Task ID: description
**Tasks Remaining in Story**: [count] or "None - story complete"
**Commit**: [commit hash, "Recorded in Git history for this iteration", or "No commit - partial progress"]
**Files Changed**: 
- path/to/file.ext
**Learnings**:
- [patterns discovered, gotchas, useful context for future iterations]
---
```

## Stop Conditions

**If ALL tasks in tasks.md are complete** (`[x]`), output exactly:

```text
<promise>COMPLETE</promise>
```

This signals the ralph loop orchestrator to terminate successfully.

**If tasks remain**, end your response normally. The next iteration will continue.

## Quality Gates

- ALL changes must pass quality checks before marking tasks complete
- DO NOT commit broken code
- Follow existing code patterns (check Codebase Patterns in progress file)
- Reference plan.md for architecture decisions
- Run tests if they exist before committing

## Code Style

Follow the patterns established in the codebase:

- Check existing files for naming conventions
- Match indentation and formatting styles
- Use the same import/module patterns
- Follow any linting rules configured in the project
- During each iteration, review newly added or substantially changed functions and tests for concise intent comments where they help future readers understand purpose, implementation highlights, or behavioral contract
- Match the existing repository comment style, for example:
  - `// newSearchQueryPageRequest builds the v8.9 page request, preferring cursor pagination when available.`
  - `// SearchTenants lists native Camunda tenants, then applies the local literal tenant name or ID filter.`
- Do not add comments mechanically just because a function or test exists
- Do not comment getters, setters, trivial wrappers, obvious one-line helpers, or simple implementation details
- Prefer comments that explain why the helper or test exists, what behavior it protects, or what non-obvious boundary or contract it preserves
- For tests, describe the behavior or regression being protected, not a restatement of the test name
- Keep comments short, precise, and implementation-aware

## Error Handling

| Condition | Expected Behavior |
| --------- | ----------------- |
| User story unclear | Ask for clarification in progress entry, mark tasks as blocked |
| Tests fail | Report failure, do not mark task complete, no commit |
| Cannot complete work unit | Report partial progress, commit only if all completed tasks form coherent unit |
| All tasks done | Commit the final work unit, then output `<promise>COMPLETE</promise>` |
| Dependencies missing | Note in progress file, skip to next available task |
