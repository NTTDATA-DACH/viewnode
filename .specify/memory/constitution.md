<!--
Sync Impact Report
Version change: unversioned template -> 1.0.0
Modified principles:
- Placeholder principle 1 -> I. Stable CLI Contracts
- Placeholder principle 2 -> II. Kubernetes Compatibility Boundaries
- Placeholder principle 3 -> III. Test-Backed Behavior Changes
- Placeholder principle 4 -> IV. Small, Traceable Delivery
- Placeholder principle 5 -> V. Documentation and Operability Parity
Added sections:
- Technology and Compatibility Constraints
- Delivery Workflow and Quality Gates
Removed sections:
- None
Templates requiring updates:
- updated: .specify/templates/plan-template.md
- updated: .specify/templates/spec-template.md
- updated: .specify/templates/tasks-template.md
- not present: .specify/templates/commands/*.md
Follow-up TODOs:
- None
-->
# viewnode Constitution

## Core Principles

### I. Stable CLI Contracts
All user-facing commands, flags, aliases, exit behavior, and text output are part of
the product contract. Changes MUST preserve documented behavior unless a feature
spec explicitly calls for a breaking change, documents the migration path, and
updates README examples in the same change. Output changes MUST be verified with
focused tests that cover the affected command path.

Rationale: `viewnode` is used as a Kubernetes CLI and krew plugin, so scripts and
operators depend on stable command behavior and readable terminal output.

### II. Kubernetes Compatibility Boundaries
Kubernetes integration MUST use the repository's established `client-go`,
`apimachinery`, and related Kubernetes libraries unless a plan justifies a
version or dependency change. Features that read kubeconfig, contexts,
namespaces, nodes, pods, containers, or metrics MUST define the expected behavior
for unavailable APIs, missing permissions, empty results, and multiple namespace
scopes.

Rationale: The tool operates against live and varied Kubernetes clusters, so
compatibility and graceful degradation are core product requirements.

### III. Test-Backed Behavior Changes
Every behavior change, bug fix, or refactor that can affect users MUST include or
update Go tests close to the changed package. Tests MUST verify preserved
behavior as well as new behavior, and command-facing changes MUST cover Cobra
command execution or output formatting where applicable. Documentation-only and
generated metadata changes MAY omit tests only when the change cannot affect
runtime behavior.

Rationale: The repository already runs `make test` with race detection and
coverage, and CLI regressions are easiest to catch near the package that owns the
command or rendering behavior.

### IV. Small, Traceable Delivery
Feature work MUST be split into independently verifiable user stories or logical
change groups. Each story MUST have acceptance criteria, exact file-level task
paths, and validation steps. Implementations MUST reuse existing repository
patterns and MUST NOT introduce parallel command, configuration, rendering, or
service structures without a documented reason in the implementation plan.

Rationale: Small slices fit the Spec Kit and Ralph workflows, reduce review risk,
and keep feature branches understandable.

### V. Documentation and Operability Parity
User-facing command, flag, output, installation, compatibility, and workflow
changes MUST update README or relevant runtime guidance in the same change.
Warnings, fatal errors, and logs MUST remain actionable for Kubernetes operators;
new failure modes MUST include clear messages and be covered by tests when they
are reachable without a live cluster.

Rationale: Operators often learn and troubleshoot `viewnode` directly from CLI
help, README examples, and terminal messages.

## Technology and Compatibility Constraints

The project baseline is Go 1.25.0 with toolchain go1.25.7. Plans MUST use the
existing Go module, Cobra command structure, Kubernetes client libraries,
logrus-based logging, and testify-based tests unless the feature explicitly
requires and justifies a change.

The supported runtime surface is the `viewnode` CLI and `kubectl viewnode` krew
plugin usage. Features MUST keep standalone and plugin invocation behavior
aligned. Cluster-facing behavior MUST account for kubeconfig selection, current
context and namespace state, selected namespaces, all-namespaces mode, missing
metrics APIs, and restricted node permissions when those concerns are in scope.

## Delivery Workflow and Quality Gates

Plans MUST pass the Constitution Check before research and again after design.
Any violation MUST be recorded in Complexity Tracking with the rejected simpler
alternative. Specs MUST define independently testable user stories and measurable
acceptance criteria before implementation tasks are generated.

Tasks MUST be ordered by dependency and user value. Behavior-changing tasks MUST
include nearby test tasks and a validation task that runs the closest relevant Go
test command. Broader validation, such as `make test`, MUST be run before
commit when practical; if it cannot be run, the reason MUST be documented.

Generated or generated-adjacent artifacts MUST be changed from their source and
regenerated when the repository provides a generation path. Commits MUST follow
Conventional Commits format and use a clear scope when one exists.

## Governance

This constitution supersedes conflicting local practice for Spec Kit features,
Ralph PRDs, implementation plans, generated tasks, and agent guidance in this
repository. Amendments require a written rationale, a semantic version bump, an
updated Sync Impact Report, and review of dependent templates and runtime
guidance.

Versioning policy:
- MAJOR: Removes or redefines a principle in a backward-incompatible way.
- MINOR: Adds a new principle or materially expands governance requirements.
- PATCH: Clarifies wording, fixes typos, or tightens non-semantic guidance.

Compliance review is required during planning, task generation, implementation,
and review. Pull requests and agent handoffs MUST call out any intentional
constitution violation, the reason it is necessary, and the validation performed.

**Version**: 1.0.0 | **Ratified**: 2026-06-26 | **Last Amended**: 2026-06-26
