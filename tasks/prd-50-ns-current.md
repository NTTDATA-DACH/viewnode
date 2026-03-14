# PRD: Namespace Current-State Commands

Issue: [#50](https://github.com/NTTDATA-DACH/viewnode/issues/50)

## Traceability

- Feature: `50-ns-current`
- Source status: derived from Spec Kit artifacts
- Spec: [`specs/50-ns-current/spec.md`](../specs/50-ns-current/spec.md)
- Plan: [`specs/50-ns-current/plan.md`](../specs/50-ns-current/plan.md)
- Tasks: [`specs/50-ns-current/tasks.md`](../specs/50-ns-current/tasks.md)

## Problem Statement

`viewnode` already provides `ctx list`, `ctx get-current`, and `ctx set-current` for kubeconfig context management, but the `ns` command group only exposes `list`. Users can see available namespaces, yet they cannot inspect or update the active namespace for the current kubeconfig context directly from `viewnode`.

## Goal

Add `viewnode ns get-current` and `viewnode ns set-current` so users can inspect and update the active namespace in a way that mirrors the existing `ctx` command workflow, respects the selected kubeconfig source, and validates namespace changes against the current cluster before persisting them.

## User Story

As a `viewnode` user working across contexts and namespaces, I want to read and update the current namespace from the CLI so I can keep my kubeconfig state aligned with the namespace I intend to work in without switching to another tool.

## Scope

- Add `viewnode ns get-current` under the existing `ns` command group.
- Add `viewnode ns set-current <namespace>` under the existing `ns` command group.
- Keep command naming, help text, and alias conventions aligned with the existing `ctx` command family, including `getcurrent` and `setcurrent`.
- Resolve kubeconfig using the same root-command flow as other commands, including support for the root `--kubeconfig` flag.
- For `get-current`, print `default` when the active kubeconfig context does not explicitly define a namespace.
- For `set-current`, validate the requested namespace against the current cluster before mutating kubeconfig.
- Update only the active kubeconfig context namespace field when persisting the change.
- Add or update automated tests covering behavior, validation, mutation, and inherited flag handling.
- Update `README.md` command listings and examples in the same change because user-facing command docs are maintained manually.

## Non-Goals

- Adding new namespace discovery, filtering, or formatting behavior beyond the existing `ns list` command.
- Changing how the root `--namespace` flag controls runtime filtering for the main `viewnode` output.
- Introducing namespace create, delete, or edit operations beyond setting the current namespace in kubeconfig.
- Adding a separate application-owned state store for namespace selection.

## Acceptance Criteria

1. `viewnode ns get-current` is available and prints the current namespace from the active kubeconfig context followed by a trailing newline.
2. If the active kubeconfig context has no explicit namespace configured, `viewnode ns get-current` prints `default`.
3. `viewnode ns get-current` and `viewnode ns set-current` respect the root `--kubeconfig` flag and use the selected kubeconfig source.
4. `viewnode ns set-current <namespace>` is available and updates only the active kubeconfig context namespace.
5. Before persisting a namespace change, `viewnode ns set-current <namespace>` validates that the namespace exists in the current cluster.
6. On success, `viewnode ns set-current <namespace>` prints `current namespace set to <namespace>`, and the kubeconfig file reflects the updated namespace for the active context.
7. If namespace validation fails because the namespace does not exist, `viewnode ns set-current <namespace>` returns a descriptive error and does not mutate kubeconfig.
8. If kubeconfig read or write operations fail during `ns set-current`, the command returns a descriptive error instead of silently succeeding.
9. The new commands follow the existing `ctx` naming pattern, including `getcurrent` and `setcurrent` aliases.
10. Automated tests cover `get-current` success, `default` fallback, `set-current` success, persisted mutation, validation failure, error handling, and inherited `--kubeconfig` propagation through a root command tree.
11. `README.md` documents `viewnode ns get-current`, `viewnode ns set-current`, and their aliases alongside the existing namespace command examples.

## Implementation Notes

- Reuse the existing command organization pattern from [`cmd/ctx`](../cmd/ctx) by extending [`cmd/ns`](../cmd/ns) with dedicated `get-current` and `set-current` command files.
- Cobra subcommands that depend on inherited persistent flags must initialize config from the root command (`cmd.Root()` fallback to the current command) so root-level `--kubeconfig` lookup works correctly.
- `get-current` should use kubeconfig state as the source of truth and preserve the Kubernetes default-namespace behavior by printing `default` when no namespace is set on the active context.
- `set-current` must validate against live cluster namespaces before writing kubeconfig; cluster validation and kubeconfig mutation are separate concerns and should be tested accordingly.
- Tests for inherited flag behavior should execute through a small root command tree with `Execute()` rather than calling the leaf command `RunE` directly.
- Documentation changes belong in the same change as the new commands because `README.md` command listings and examples are maintained manually in this repository.

## Suggested Work Packages

1. Wire `get-current` and `set-current` into `cmd/ns` with `ctx`-aligned naming and aliases.
2. Implement `ns get-current` with kubeconfig lookup and `default` fallback behavior.
3. Implement `ns set-current` with active-context lookup, namespace validation, kubeconfig mutation, and success/error messaging.
4. Add `cmd/ns` tests for success cases, fallback behavior, validation failure, mutation checks, error paths, and inherited `--kubeconfig` handling.
5. Update `README.md` namespace command examples and aliases, then run `make test`.
