# Quickstart: Refactor `ctx list` and `ns list` commands

## Goal

Implement the refactor so `viewnode ctx list` and `viewnode ns list` share list-preparation behavior while keeping output unchanged.

## Steps

1. Create explicit list-entry types for context and namespace display rows.
2. Extract the duplicate list-preparation logic into a shared helper that:
   - accepts raw names plus one active name
   - sorts names deterministically
   - marks active entries consistently
3. Refactor [cmd/ctx/list.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/list.go) to use the shared helper without changing kubeconfig retrieval.
4. Refactor [cmd/ns/list.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list.go) to use the shared helper and normalize an empty active namespace to `default`.
5. Update or add tests in:
   - [cmd/ctx/ctx_test.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/ctx_test.go)
   - [cmd/ns/list_test.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list_test.go)
   - any new focused helper tests if the helper is extracted into a directly testable unit
6. Run `make test`.

## Expected Verification

- `viewnode ctx list` still prints alphabetically ordered contexts with a single active marker.
- `viewnode ns list` still prints alphabetically ordered namespaces with a single active marker.
- `viewnode ns list` marks `default` as active when the current context namespace is empty.
- Failure paths still return errors without partial output.
