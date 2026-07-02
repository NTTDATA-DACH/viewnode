---
name: gvb-ghissue-to-ralph
description: Create an issue-backed Speckit feature and drive it through clarify, plan, tasks, and Ralph launch using the repository's standard flow.
compatibility: Requires spec-kit project structure with .specify/ directory
metadata:
  author: local-wrapper
  wraps:
    - gvb-ghissue-to-speckit
    - speckit-clarify
    - speckit-arch-generate
    - speckit-arch-reverse
    - speckit-plan
    - speckit-tasks
    - speckit-agent-context-update
    - speckit-converge
    - speckit-ralph-run
---

# GVB Issue To Ralph Flow

Use this skill when the user wants the repository's standard issue-driven delivery flow executed from a GitHub issue through ready-to-run implementation.

This wrapper standardizes the following sequence:

1. `gvb-ghissue-to-speckit`
2. `speckit-clarify`
3. architecture grounding with existing architecture memory, `speckit-arch-reverse`, or `speckit-arch-generate`
4. `speckit-plan`
5. `speckit-tasks`
6. `speckit-ralph-run` or the equivalent `after_tasks` hook
7. post-Ralph handoff instructions for `speckit-converge`

The intended outcome is an issue-backed feature with traceable spec artifacts, completed planning artifacts, generated tasks, and a Ralph loop launch path that respects the repository's hook configuration.

## Inputs

- A GitHub issue URL such as `https://github.com/grafvonb/c8volt/issues/45`
- Optionally, extra instructions in the same request
- Optionally, Ralph launcher arguments such as `--max-iterations 5` or `--model gpt-5.5`

A bare issue URL is sufficient to begin.

## Required Behavior

1. Parse the GitHub issue URL from the user request.
2. Treat the issue number as the canonical workflow identifier for the entire run.
3. Run the `gvb-ghissue-to-speckit` workflow first and preserve all of its guarantees:
   - fetch issue title and body
   - enforce issue-number-based feature identity
   - verify shell-layer compatibility before feature creation
   - persist issue traceability into `spec.md`
   - ensure both `spec.md` and `checklists/requirements.md` exist before moving on
4. After specification, always continue immediately into `speckit-clarify` unless the user explicitly asked to skip clarification in the same `gvb-ghissue-to-ralph` request.
5. Treat clarification as a required gate before planning:
   - ask at most the allowed clarification questions
   - wait for the user to answer each clarification question
   - integrate accepted answers into `spec.md`
   - do not run `speckit-plan` or `speckit-tasks` until clarification has completed, all accepted answers have been written to `spec.md`, or the user explicitly instructs you to proceed without resolving remaining clarification questions
6. After clarification completes and before planning, perform architecture grounding:
   - if `.specify/extensions/arch/` is unavailable, continue but report that architecture grounding was skipped because the target repository does not have the architecture extension installed
   - before creating or refreshing architecture memory, verify `.specify/memory/architecture*.md` and `.specify/memory/architecture-repo-facts.md` are trackable by git
   - if architecture memory is missing, create it before planning using `speckit-arch-reverse` for existing repositories with observable code/docs/tests/config/deployment evidence, or `speckit-arch-generate` for greenfield or intent-first work
   - if architecture memory exists, refresh it only when the issue or clarification answers affect stable architecture concerns such as module boundaries, runtime responsibilities, deployment assumptions, external systems, ownership boundaries, or forbidden responsibility crossings
   - for narrow chores, docs updates, dependency bumps, tests, local bug fixes, and other work that does not affect architecture boundaries, reuse the existing architecture memory without refreshing it
7. After architecture grounding completes, run `speckit-plan`.
8. After planning completes, run `speckit-tasks`.
9. After tasks are generated, decide how Ralph should start:
   - always stop and ask the user for explicit confirmation before any Ralph launch path is allowed to continue
   - treat this confirmation as a budget check, not as a mere informational prompt
   - on macOS, always phrase the eventual launch as `speckit-ralph-run in Terminal.app`
   - if `.specify/extensions.yml` defines an enabled `after_tasks` hook targeting `speckit.ralph.run`, prefer that hook-driven path after the user confirms and do not manually invoke `speckit-ralph-run` a second time
   - otherwise invoke `speckit-ralph-run` explicitly only after the user confirms
10. When a hook-driven Ralph launch is available, surface that to the user clearly and ask whether they want to spend budget on starting Ralph now.
11. Respect enabled non-Ralph hooks such as `speckit.agent-context.update`; do not suppress hook-managed context refresh while applying the Ralph budget gate.
12. After Ralph is launched externally, provide a post-Ralph handoff: when Ralph finishes, run `speckit-converge`; if it appends tasks to `tasks.md`, ask for budget confirmation before running `speckit-ralph-run` again.
13. Carry the GitHub issue number through all downstream naming and traceability where relevant, including commit messages.

## Commit Rule

For any commit created as part of this workflow or by downstream skills acting on the generated feature, preserve Conventional Commits formatting and append `#<issue-number>` as the final token in the subject.

Examples:

- `feat(spec): add issue traceability #45`
- `chore(plan): record architecture decisions #45`
- `test(view): complete final polish #59`

## Ralph Launch Rule

This repository may already wire Ralph through `.specify/extensions.yml`.

When the `after_tasks` hook already points at `speckit.ralph.run`:

- always ask the user whether Ralph should run now before allowing the hook to proceed
- on macOS, require the launch to use Terminal.app visibly
- prefer the hook-managed flow
- do not manually duplicate the Ralph launch
- report that Ralph is available through the existing hook and that budget confirmation is required

When that hook is absent or disabled:

- always ask the user whether Ralph should run now before launching it
- on macOS, launch it as `speckit-ralph-run in Terminal.app`
- invoke `speckit-ralph-run` directly
- pass through any user-provided Ralph arguments unchanged when they are compatible with the launcher

## Hook Handling Rule

This wrapper gates only Ralph launch because Ralph spends implementation budget and may run externally in Terminal.app.

For other enabled hooks, including `speckit.agent-context.update`, preserve the repository's hook-managed behavior. Surface optional hooks when useful, execute mandatory hooks as required by the wrapped Speckit skill instructions, and do not replace hook-managed context refresh with manual edits to `AGENTS.md`.

## Clarification Gate

When this skill is invoked, `speckit-clarify` is mandatory immediately after `gvb-ghissue-to-speckit` creates or updates the spec.

Only skip clarification when the user explicitly says to skip clarification while invoking `gvb-ghissue-to-ralph`. General urgency, a bare issue URL, Ralph options, or a request to run the full flow are not skip signals.

Clarification is complete only after one of these conditions is true:

- `speckit-clarify` finds no critical ambiguities worth formal clarification
- the user has answered all clarification questions selected by `speckit-clarify`, and those answers have been integrated into `spec.md`
- the user explicitly stops the clarification loop or instructs the agent to proceed despite outstanding clarification items

Do not run `speckit-plan`, `speckit-tasks`, an `after_tasks` hook, or `speckit-ralph-run` before the clarification gate is complete.

## Architecture Grounding Gate

Architecture grounding runs after clarification and before planning.

The purpose is to make sure `speckit-plan` sees stable project-level architecture context instead of deriving architecture boundaries from one issue in isolation.

This gate is conditional:

- If `.specify/extensions/arch/` is unavailable, do not block the whole issue flow solely for that reason. Continue after reporting that architecture grounding was skipped because the target repository needs an ai-tooling refresh with the architecture extension.
- Before creating or refreshing architecture memory, verify the architecture SSOT paths are not ignored by git:
  - `.specify/memory/architecture.md`
  - `.specify/memory/architecture-scenario-view.md`
  - `.specify/memory/architecture-logical-view.md`
  - `.specify/memory/architecture-process-view.md`
  - `.specify/memory/architecture-development-view.md`
  - `.specify/memory/architecture-physical-view.md`
  - `.specify/memory/architecture-repo-facts.md`
- If any architecture memory path is ignored, stop before running `speckit-arch-generate` or `speckit-arch-reverse` and tell the user to fix `.gitignore`. Do not write architecture SSOT into ignored paths.
- If `.specify/memory/architecture.md` or any required 4+1 view file is missing, run an architecture workflow before planning.
- Use `speckit-arch-reverse` when the repository already has observable implementation, docs, tests, entry points, configuration, CI/CD, or deployment evidence.
- Use `speckit-arch-generate` when the work is primarily greenfield or intent-driven and reliable repository evidence does not yet exist.
- If architecture memory already exists, refresh it only when the issue or clarification answers affect stable architecture concerns such as module boundaries, runtime responsibilities, deployment assumptions, external systems, ownership boundaries, or forbidden responsibility crossings.
- For narrow chores, documentation-only changes, dependency bumps, test cleanup, local bug fixes, and other work that does not affect architecture boundaries, do not refresh architecture memory. Preserve the existing architecture SSOT and rely on `memory-loader` to load it before planning.

## Post-Ralph Handoff

Ralph may run in an external Terminal.app session, so this wrapper does not control the post-implementation loop automatically.

After Ralph is launched, end with this handoff:

1. Wait for the external Ralph run to finish.
2. Run `speckit-converge` to assess the current codebase against `spec.md`, `plan.md`, and `tasks.md`.
3. If `speckit-converge` appends a new convergence phase to `tasks.md`, ask the user for another explicit Ralph budget confirmation.
4. Run `speckit-ralph-run` again only after that confirmation.
5. Repeat the converge/Ralph cycle only with explicit user confirmation for each Ralph run.

Do not run `speckit-converge` before the first Ralph implementation run. It is a post-implementation gap-to-task step, not a planning gate.

## Execution Flow

1. Read the user request and locate the GitHub issue URL.
2. Extract any extra feature instructions and any Ralph runtime options.
3. Run the complete `gvb-ghissue-to-speckit` flow.
4. Confirm the created feature artifacts exist and that issue traceability is present in the spec.
5. Unless the current `gvb-ghissue-to-ralph` request explicitly says to skip clarification, run `speckit-clarify` against the active feature immediately after spec creation.
6. Complete the clarification loop before planning:
   - ask one clarification question at a time
   - wait for the user's answer
   - integrate each accepted answer into `spec.md`
   - stop only when `speckit-clarify` reports no critical ambiguities, all selected questions are answered, the question limit is reached, or the user explicitly stops or proceeds with remaining ambiguity
7. Perform the architecture grounding gate:
   - inspect whether `.specify/extensions/arch/` is installed
   - inspect whether `.specify/memory/architecture.md` and the five 4+1 view files exist
   - verify architecture memory paths are not ignored before creating or refreshing them
   - create architecture memory with `speckit-arch-reverse` or `speckit-arch-generate` when it is missing
   - refresh architecture memory only when the issue affects stable architecture concerns
   - otherwise reuse the existing architecture memory for `memory-loader`
8. Run `speckit-plan` only after the clarification and architecture grounding gates are complete.
9. Run `speckit-tasks` after planning.
10. Inspect `.specify/extensions.yml` for an enabled `after_tasks` hook targeting `speckit.ralph.run`.
11. Stop and ask the user whether they want to launch Ralph now, explicitly framing the question as a budget check.
12. If the user declines, stop after confirming tasks generation is complete and report the next command to run later.
13. If the user approves and the hook exists, present the hook-driven Ralph handoff and stop after confirming the launch path.
14. If the user approves and that hook does not exist, invoke `speckit-ralph-run`.
15. In the final handoff, restate:
   - the feature directory
   - the spec, plan, and tasks artifact paths
   - whether architecture memory was created, refreshed, reused, or skipped
   - whether Ralph was hook-driven or manually launched
   - after external Ralph completion, run `speckit-converge`; if it appends tasks, ask before rerunning Ralph
   - the issue-suffixed commit-message rule

## Notes

- This skill is intentionally opinionated for repositories that use GitHub issues as the primary feature entry point.
- `speckit-clarify` remains interactive by design. That is a feature, not a flaw. Do not try to bypass it unless the user explicitly chooses to skip it in the current `gvb-ghissue-to-ralph` request.
- Ralph launch is also an explicit pause point in this repository because the user wants to check budget before implementation starts.
- Prefer the repository's existing hook configuration over inventing a second orchestration path.
- Do not run `speckit-ralph-run` twice.
- Do not automatically run `speckit-converge` or a second Ralph loop after launching Ralph externally. Report the post-Ralph handoff and wait for the user to return after Ralph finishes.
- If `speckit-converge` later appends tasks, treat the next Ralph run as a new budget decision.
- If the issue-backed specification step fails, stop there and surface the exact blocker instead of attempting later stages.
- If planning fails, do not generate tasks.
- If task generation fails, do not attempt Ralph launch.
- If the user asks for a dry run or to stop before implementation, still run clarification unless they explicitly ask to skip it; then finish at the latest completed stage and report the next command to run.
