# PRD: Refactor `ctx list` and `ns list` commands

Issue: [#48](https://github.com/NTTDATA-DACH/viewnode/issues/48)

## Traceability

- Feature: `48-refactor-ctx-list-ns-list-commands`
- Source status: derived from Spec Kit artifacts
- Spec: [`specs/48-refactor-ctx-list-ns-list-commands/spec.md`](../specs/48-refactor-ctx-list-ns-list-commands/spec.md)
- Plan: [`specs/48-refactor-ctx-list-ns-list-commands/plan.md`](../specs/48-refactor-ctx-list-ns-list-commands/plan.md)
- Tasks: [`specs/48-refactor-ctx-list-ns-list-commands/tasks.md`](../specs/48-refactor-ctx-list-ns-list-commands/tasks.md)

## Problem Statement

`viewnode ctx list` and `viewnode ns list` currently implement similar list-building behavior separately. Each command gathers names, sorts them, determines which item is active, and renders marker-prefixed lines on its own. That duplication makes the list commands harder to maintain and increases the risk that behavior diverges when one command changes and the other does not.

## Goal

Refactor `viewnode ctx list` and `viewnode ns list` so they share the duplicated list-preparation behavior through explicit entry structures and reusable helper logic, while preserving current CLI behavior, output format, ordering, active-state semantics, and error handling.

## Engineering Story

As a maintainer of `viewnode`, I want the context and namespace list commands to use the same internal list-entry pattern so I can reduce duplicated logic without changing the user-facing contract of either command.

## Scope

- Introduce explicit list-entry structures for context list rows and namespace list rows.
- Extract the duplicated logic that prepares ordered active-marked list entries from raw names.
- Keep `ctx list` data retrieval based on kubeconfig raw context data.
- Keep `ns list` data retrieval based on the existing Kubernetes namespace query path.
- Preserve current marker formatting (`[*]` and `[ ]`) and alphabetical output ordering.
- Preserve current command names, aliases, help behavior, and existing error flow.
- Preserve namespace semantics by treating an empty namespace on the active kubeconfig context as `default` when marking the active namespace.
- Add or update automated tests for successful listing, active marker behavior, `default` fallback for namespace selection, and failure paths.
- Run `make test` before completion.

## Non-Goals

- Adding new commands, flags, aliases, or output columns.
- Changing how `ctx list` resolves kubeconfig state.
- Changing how `ns list` queries namespaces from the cluster.
- Introducing new namespace filtering, formatting variants, or sorting modes.
- Updating `README.md`, unless the refactor unexpectedly changes user-visible behavior and documentation must be corrected.
- Introducing a broader abstraction that combines kubeconfig parsing and Kubernetes API access into one new command service.

## Acceptance Criteria

1. `viewnode ctx list` continues to print all discovered contexts in deterministic alphabetical order and marks only the current context with `[*]`.
2. `viewnode ns list` continues to print all discovered namespaces in deterministic alphabetical order and marks only the active namespace with `[*]`.
3. If the active kubeconfig context does not explicitly set a namespace, `viewnode ns list` treats `default` as the active namespace for marker comparison.
4. Both commands continue to render each list item in the existing `[marker] name` format without adding or removing user-visible fields.
5. The duplicated list-preparation logic used by `ctx list` and `ns list` is replaced by a shared reusable path based on explicit list-entry structures.
6. `viewnode ctx list` continues returning a descriptive error and no partial list output when kubeconfig loading fails.
7. `viewnode ns list` continues returning a descriptive error and no partial list output when namespace discovery fails.
8. Automated tests cover context listing success, namespace listing success, namespace `default` fallback behavior, and at least one failure case for each affected command path.
9. `make test` passes after the refactor is complete.

## Implementation Notes

- Reuse the existing command layout under [`cmd/ctx`](../cmd/ctx), [`cmd/ns`](../cmd/ns), and [`cmd/config`](../cmd/config) rather than introducing a new parallel command structure.
- Keep shared behavior narrowly scoped to list-entry preparation and rendering inputs; data retrieval should remain command-specific because kubeconfig context discovery and cluster namespace discovery have different sources and failure modes.
- Introduce a dedicated context entry type and namespace entry type with fields for display name and active state so ordering and marker behavior can be tested independently of the Cobra command wrapper.
- Namespace active-state comparison should follow kubeconfig semantics and normalize an empty active namespace to `default` before entry preparation.
- Existing error wrapping in [`cmd/ctx/list.go`](../cmd/ctx/list.go) and [`cmd/ns/list.go`](../cmd/ns/list.go) should be preserved to avoid behavioral regressions.
- Tests should extend the existing command coverage in [`cmd/ctx/ctx_test.go`](../cmd/ctx/ctx_test.go) and [`cmd/ns/list_test.go`](../cmd/ns/list_test.go), with focused helper tests added if the shared preparation logic is extracted into a directly testable unit.
- Performance-sensitive behavior should remain unchanged: the refactor must not add extra namespace API calls or additional kubeconfig reads beyond the current command flow.

## Regression Risks

- The active marker could shift if the shared helper changes ordering or active-name comparison semantics.
- `ns list` could fail to mark `default` correctly if empty namespace normalization is lost during extraction.
- Error behavior could regress if helper integration prints partial output before command failures are returned.
- Over-abstracting the helper could couple context and namespace retrieval in a way that makes later maintenance harder instead of easier.

## Suggested Work Packages

1. Add shared list-entry preparation types and helper logic in the existing command support area, with focused helper tests.
2. Refactor `ctx list` to use the shared helper while preserving current output and context-selection behavior.
3. Refactor `ns list` to use the shared helper while preserving namespace query behavior and `default` fallback semantics.
4. Extend command-level tests for ordering, active marker behavior, and failure propagation across both commands.
5. Run `make test` and verify final behavior against the existing list-command contract before merge.

## Assumptions / Open Questions

- This issue remains a behavior-preserving refactor, so no README update is required unless implementation uncovers a mismatch between actual command behavior and current documentation.
- The most appropriate shared location for helper code is whichever existing command support package keeps the refactor small and avoids creating a new parallel structure; the exact file can be finalized during implementation.
