# Quickstart: Context list filtering

## Goal

Implement context-name filtering for `viewnode ctx list`, keep filtered output aligned with the existing context list format, and document the new behavior.

## Steps

1. Add a local `--filter` / `-f` flag in [cmd/ctx/list.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/list.go).
2. Keep the existing raw kubeconfig discovery path, then apply case-insensitive substring filtering before rendering context entries.
3. Preserve alphabetical ordering and active-marker behavior for displayed rows by continuing to use [cmd/internal/listing/entries.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/internal/listing/entries.go).
4. Add the exact success-path no-match message `no contexts matched filter "<value>"` in [cmd/ctx/list.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/list.go).
5. Ensure filtered output contains only matching contexts, even when that means the current context is omitted and no active marker is shown.
6. Add or update tests in [cmd/ctx/ctx_test.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/ctx_test.go) to cover long flag, short flag, empty filter, no-match behavior, filtered-out current context behavior, and kubeconfig error handling.
7. Update [README.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/README.md) with the new context filter flag, shorthand, and examples.
8. Run `make test`.

## Expected Verification

- `viewnode ctx list --filter prod` returns only contexts whose names contain `prod`, regardless of letter case.
- `viewnode ctx list -f prod` matches the long-flag behavior exactly.
- `viewnode ctx list --filter ""` behaves the same as `viewnode ctx list`.
- `viewnode ctx list --filter missing` prints `no contexts matched filter "missing"` and succeeds.
- When the current context does not match the filter, the filtered output contains only matching contexts and no row is marked active.
- README context command examples reflect the new filtering behavior.
