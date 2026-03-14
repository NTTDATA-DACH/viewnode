# Plan: Namespace current-state subcommands

## Summary

Implement `viewnode ns get-current` and `viewnode ns set-current` by extending the existing `cmd/ns` package to mirror the established `cmd/ctx` command pattern while preserving repository conventions for kubeconfig handling, inherited flags, tests, and README maintenance.

## Implementation Plan

### 1. Add namespace current-state command wiring

- Extend `cmd/ns/ns.go` to register `get-current` and `set-current` alongside `list`.
- Keep command naming and aliases aligned with `cmd/ctx`, including `getcurrent` and `setcurrent`.
- Reuse the root-command-based config initialization pattern so inherited `--kubeconfig` lookup works correctly.

### 2. Implement `ns get-current`

- Add a dedicated command file in `cmd/ns/` for `get-current`.
- Initialize config from the root command, then read the current namespace from the active kubeconfig context.
- Print `default` when the active context has no explicit namespace configured.
- Return kubeconfig-related failures with descriptive wrapped errors.

### 3. Implement `ns set-current`

- Add a dedicated command file in `cmd/ns/` for `set-current`.
- Initialize config from the root command and load the active kubeconfig.
- Resolve the active context and validate the requested namespace against the current cluster before mutating kubeconfig.
- Update only the active context's namespace field and persist the kubeconfig change.
- Print `current namespace set to <namespace>` on success.
- Return explicit errors for missing namespaces, kubeconfig read failures, and kubeconfig write failures.

### 4. Add targeted tests

- Add `cmd/ns` tests for `get-current` success and `default` fallback behavior using kubeconfig fixtures.
- Add `cmd/ns` tests for `set-current` success, persisted kubeconfig mutation, missing-namespace validation failure, and write/error paths.
- Cover inherited `--kubeconfig` propagation by executing through a small root command tree with `Execute()` rather than invoking leaf handlers directly when flag inheritance matters.
- Reuse the existing `ctx` test style where practical to keep fixture handling consistent.

### 5. Update documentation and validate

- Update `README.md` command reference and namespace command examples to document `ns get-current`, `ns set-current`, and their aliases.
- Run `make test` before considering the change complete, per repository guidance.

## Risks And Checks

- Namespace validation must query the current cluster before writing kubeconfig; tests should ensure validation failure does not mutate the file.
- The active context may omit the namespace field entirely; `get-current` must still produce `default` instead of an empty value.
- `set-current` must change only kubeconfig state for the active context and must not affect the root `--namespace` runtime filter behavior.

## Expected Output

- New `cmd/ns` subcommands for `get-current` and `set-current`.
- Automated test coverage for success, fallback, validation, mutation, and inherited flag behavior.
- Updated `README.md` examples and aliases for namespace current-state management.
