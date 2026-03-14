<!--
Sync Impact Report
- Version change: 0.0.0 -> 1.0.0
- Modified principles:
  - Template principle 1 -> I. Code Quality Is a Release Gate
  - Template principle 2 -> II. Tests Prove Behavior Before Merge
  - Template principle 3 -> III. CLI Experience Must Stay Consistent
  - Template principle 4 -> IV. Performance Regressions Require Evidence
  - Template principle 5 -> V. Keep Changes Small, Observable, and Reversible
- Added sections:
  - Delivery Standards
  - Workflow & Review
- Removed sections:
  - None
- Templates requiring updates:
  - ✅ updated /.specify/templates/plan-template.md
  - ✅ updated /.specify/templates/spec-template.md
  - ✅ updated /.specify/templates/tasks-template.md
- Follow-up TODOs:
  - None
-->
# viewnode Constitution

## Core Principles

### I. Code Quality Is a Release Gate
Every production change MUST keep the codebase readable, bounded in scope, and
consistent with existing Go and Cobra patterns in this repository. New logic
MUST be implemented in the smallest unit that preserves clarity, avoid parallel
command structures when an established pattern exists, and update user-facing
documentation in `README.md` whenever command behavior changes. Reviewers MUST
reject changes that increase complexity without a concrete operational or
maintainability payoff.

### II. Tests Prove Behavior Before Merge
Every change to behavior MUST include automated tests covering the affected
execution path at the appropriate level. Command behavior that depends on Cobra
flag inheritance MUST be tested through a realistic root command tree using
`Execute()`, and global `--kubeconfig` consumers MUST initialize config from the
root command per repository convention. A change is not mergeable until
`make test` passes locally with any new or updated tests included in the same
change.

### III. CLI Experience Must Stay Consistent
The CLI MUST present predictable naming, flag semantics, output wording, and
error handling across commands. New commands and flags MUST follow existing
verbs, aliases, and help text style; output intended for humans MUST remain
scannable and stable enough for documentation examples; and any user-visible
change MUST update the README command reference and examples in the same
iteration. Breaking UX changes require explicit justification in the feature
plan and reviewer approval.

### IV. Performance Regressions Require Evidence
Features MUST preserve fast feedback for interactive CLI use on real cluster
data. Plans MUST state expected impact on API calls, render cost, memory usage,
and watch-mode behavior when applicable. Any change that could materially add
latency, repeated Kubernetes requests, or output rendering overhead MUST include
measurement or reasoned evidence that the impact stays acceptable for routine
operator workflows.

### V. Keep Changes Small, Observable, and Reversible
Work MUST be sliced into independently verifiable increments with clear
acceptance criteria, matching the project’s preference for narrowly scoped user
stories and small conventional commits. Each increment MUST expose enough
signals to debug failures quickly, either through tests, error messages, or
existing logging patterns, and MUST avoid hidden side effects that make rollback
or diagnosis difficult.

## Delivery Standards

Specifications and plans MUST describe the user-facing outcome, the code paths
being changed, the tests that will prove correctness, the documentation impact,
and any performance-sensitive behavior. Requirements are incomplete until they
state UX consistency expectations for commands, flags, help text, and output
format, plus measurable performance expectations when a change affects cluster
queries, watch loops, or rendering.

## Workflow & Review

Implementation plans MUST pass a constitution check before design and again
before coding starts. Tasks MUST include documentation updates for user-visible
changes and explicit verification work for any affected tests or performance
assumptions. Code review MUST confirm:

- the change follows existing repository patterns before introducing new ones
- tests cover the changed behavior and `make test` passes
- README updates accompany user-facing CLI changes
- performance-sensitive changes include evidence or a measurement plan

## Governance

This constitution overrides ad hoc implementation preferences for this
repository. Amendments require the updated principle text, a justification for
the semantic version bump, and synchronization of affected templates before the
change is complete. Every feature spec, plan, task list, and review is expected
to check compliance against these principles, and exceptions MUST be documented
in the relevant plan under a clear complexity or tradeoff note.

**Version**: 1.0.0 | **Ratified**: 2026-03-13 | **Last Amended**: 2026-03-13
