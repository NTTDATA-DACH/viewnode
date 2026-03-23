# Quickstart: Namespace list filtering

## Goal

Implement namespace-name filtering for `viewnode ns list`, add clear no-match messaging for filtered namespace and node output, and document the new behavior.

## Steps

1. Add a local `--filter` / `-f` flag in [cmd/ns/list.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list.go).
2. Keep the existing namespace discovery call, then apply case-insensitive substring filtering before rendering namespace entries.
3. Preserve active namespace handling by continuing to derive the active value from kubeconfig and defaulting an empty active namespace to `default`.
4. Add a clear success-path no-match message for filtered namespace results in [cmd/ns/list.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list.go).
5. Update the root filtered empty-state wording in [srv/view.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go) and any needed plumbing from [cmd/root.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root.go).
6. Add or update tests in:
   - [cmd/ns/list_test.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list_test.go)
   - [srv/view_test.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go)
   - any additional focused tests for namespace filtering helpers if extracted
7. Update [README.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/README.md) with the new namespace filter flag, shorthand, and examples.
8. Run `make test`.

## Expected Verification

- `viewnode ns list --filter kube` returns only namespaces whose names contain `kube`, regardless of letter case.
- `viewnode ns list -f kube` matches the long-flag behavior exactly.
- `viewnode ns list --filter ""` behaves the same as `viewnode ns list`.
- Filtered no-match results for both namespace and node flows produce clear messages while returning success.
- README command reference and examples reflect the new namespace filtering behavior.
