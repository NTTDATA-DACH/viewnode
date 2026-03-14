# Spec: Namespace current-state subcommands

## Title

Add `ns get-current` and `ns set-current` subcommands

## Problem Statement

`viewnode` already exposes `ctx list`, `ctx get-current`, and `ctx set-current` for managing kubeconfig contexts, but the `ns` command only supports `list`. Users cannot inspect or change the current namespace directly through `viewnode`, even though namespace selection is a core part of the active kubeconfig context.

## Goal

Extend `viewnode ns` so users can read and update the current namespace in a way that matches the existing `ctx` command experience and respects the selected kubeconfig source.

## User Story

As a `viewnode` user working across Kubernetes contexts and namespaces, I want to get the current namespace and set a new current namespace from the CLI so I can manage my active kubeconfig state without switching to another tool.

## Scope

- Add `viewnode ns get-current` to print the current namespace.
- Add `viewnode ns set-current <namespace>` to persist a new current namespace for the active kubeconfig context.
- Align command shape, help text, and aliases with the existing `ctx get-current` and `ctx set-current` pattern where applicable.
- Respect the root-level `--kubeconfig` flag when resolving and mutating kubeconfig data.
- Validate the requested namespace against the current cluster before persisting it.
- Update the manually maintained `README.md` command reference and examples for the new user-facing namespace commands.

## Non-Goals

- Adding new namespace discovery or filtering behavior beyond the existing `ns list` command.
- Changing how the root `--namespace` flag filters `viewnode` output.
- Introducing parallel namespace-management workflows outside the current Cobra command structure.

## Constraints And Compatibility

- The new `ns` subcommands must follow the repository's Cobra conventions for inherited persistent flags and initialize config using the root command so `--kubeconfig` resolves correctly.
- Behavior should remain aligned with kubeconfig-backed operations already used by the `ctx` command family.
- Tests for inherited flag behavior should execute through a small root command tree with `Execute()` rather than calling leaf `RunE` handlers directly.
- README command listings and examples must be updated in the same change because they are maintained manually.

## Acceptance Criteria

- AC-001: `viewnode ns get-current` is available under the `ns` command and prints the current namespace resolved from the active kubeconfig context followed by a trailing newline.
- AC-002: If the active kubeconfig context does not explicitly set a namespace, `viewnode ns get-current` prints `default`.
- AC-003: `viewnode ns get-current` uses the same kubeconfig source selection as the rest of the CLI, including support for the root-level `--kubeconfig` flag.
- AC-004: `viewnode ns set-current <namespace>` is available under the `ns` command and updates only the current namespace stored for the active kubeconfig context.
- AC-005: Before persisting a new namespace, `viewnode ns set-current <namespace>` validates that the namespace exists in the current cluster.
- AC-006: After `viewnode ns set-current <namespace>` succeeds, the command prints `current namespace set to <namespace>`, and a subsequent read of the kubeconfig reflects the updated namespace for the active context.
- AC-007: If validation fails because the requested namespace does not exist, `viewnode ns set-current <namespace>` returns a descriptive error and does not mutate kubeconfig.
- AC-008: If `viewnode ns set-current <namespace>` cannot complete the kubeconfig update for another reason, the command returns a descriptive error instead of silently succeeding.
- AC-009: The new namespace subcommands follow the existing `ctx` naming pattern, including `getcurrent` and `setcurrent` aliases in addition to `get-current` and `set-current`.
- AC-010: Automated tests cover successful `get-current` and `set-current` flows, `default` fallback behavior, namespace validation, kubeconfig mutation for `set-current`, error handling, and inherited `--kubeconfig` flag propagation through a root command tree.
- AC-011: `README.md` documents `viewnode ns get-current` and `viewnode ns set-current` alongside the existing namespace command examples and aliases.

## Clarified Decisions

- `viewnode ns set-current <namespace>` must validate that the namespace exists before writing kubeconfig.
- If no namespace is set on the active kubeconfig context, `viewnode ns get-current` prints `default`.
- `viewnode ns set-current <namespace>` updates only the active kubeconfig context and does not interact with the root `--namespace` runtime filter flag.
- Alias parity with `ctx` is required, including `getcurrent` and `setcurrent`.
- The success message must be `current namespace set to <namespace>`.
