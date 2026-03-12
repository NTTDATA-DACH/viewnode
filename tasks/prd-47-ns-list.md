# PRD: Namespace List Command

Issue: [#47](https://github.com/NTTDATA-DACH/viewnode/issues/47)

## Problem Statement

`viewnode` already provides a `ctx` command group for kubeconfig context operations, but it does not provide a matching command group for namespace-focused workflows. Users can target namespaces through flags on the root command, yet they cannot inspect which namespaces are available in the currently selected context or see which namespace is currently active without leaving `viewnode`.

## Goal

Add an `ns` command group with an initial `list` subcommand so users can inspect namespaces available in the current Kubernetes context directly from `viewnode`, with output and command structure aligned to the existing `ctx` commands.

## User Story

As a `viewnode` user working across Kubernetes environments, I want to run a namespace listing command in the current context so that I can see which namespaces are available and which namespace is currently selected before running other namespace-scoped operations.

## Scope

- Add a new top-level Cobra command `ns`.
- Add an `ns list` subcommand for listing namespaces reachable in the currently selected kubeconfig context.
- Mark the current namespace in the output using the same marker format and behavior used by `viewnode ctx list`.
- Keep command wording and help text consistent with the existing `ctx` command structure.
- Support the existing `--kubeconfig` global flag when resolving the current context and namespace list.
- Add or update automated tests for command behavior.
- Document the new command in the usage documentation if command listings or examples are maintained manually.

## Non-Goals

- Changing the current namespace from `viewnode`.
- Adding namespace create, delete, or edit operations.
- Introducing namespace filtering, formatting variants, or extra output metadata beyond the list and current marker.
- Changing how the main `viewnode` root command resolves `--namespace` or `--all-namespaces`.

## Acceptance Criteria

1. Running `viewnode ns list` lists namespaces available in the currently selected Kubernetes context.
2. The currently selected namespace is marked with `[*]` and non-selected namespaces with `[ ]`, matching `viewnode ctx list`.
3. `viewnode ns` acts as the root for namespace-related subcommands and follows the same structural conventions as `viewnode ctx`.
4. The command respects the existing `--kubeconfig` flag and uses that kubeconfig when resolving namespaces.
5. If namespace discovery fails, the command returns an actionable error through the existing command error flow.
6. Automated tests cover successful listing, current namespace marking, and at least one failure case.
7. User-facing documentation reflects the new `ns list` command if the command reference in `README.md` is kept in sync manually.

## Implementation Notes

- Reuse the existing command organization pattern from [`cmd/ctx`](../cmd/ctx) by introducing a parallel `cmd/ns` package with a root command and subcommands.
- The current namespace is already derived from kubeconfig through `cmd/config.Setup.GetCurrentNamespace()`. The new command should stay aligned with that source of truth so the marked namespace matches the namespace `viewnode` would use by default.
- Namespace enumeration should use the initialized Kubernetes client configuration and cluster API rather than parsing namespaces from kubeconfig, because available namespaces are cluster state.
- Output ordering should be deterministic, consistent with `ctx list`, so test assertions and user output are stable.
- Tests can follow the existing `cmd/ctx/ctx_test.go` pattern, but namespace-list coverage will likely require API mocking rather than kubeconfig-only fixtures because namespaces are not stored as a list in kubeconfig.
- Keep the initial scope to `ns list` so later Ralph conversion can split follow-up stories for additional namespace subcommands without rewriting this requirement.
