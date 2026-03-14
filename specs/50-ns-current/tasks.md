# Tasks: Namespace current-state subcommands

## Task 1: Wire new namespace current-state commands

- Add `get-current` and `set-current` command registration to `cmd/ns/ns.go`.
- Match the existing `ctx` command naming pattern, including `getcurrent` and `setcurrent` aliases.

Acceptance coverage: AC-001, AC-004, AC-009

## Task 2: Implement `ns get-current`

- Add a `cmd/ns/getcurrent.go` command that initializes config from the root command.
- Read the current namespace from the active kubeconfig context.
- Print `default` when the active context does not define a namespace.
- Return descriptive errors for kubeconfig initialization or raw-config lookup failures.

Acceptance coverage: AC-001, AC-002, AC-003

## Task 3: Implement `ns set-current`

- Add a `cmd/ns/setcurrent.go` command that accepts exactly one namespace argument.
- Initialize config from the root command and resolve the active kubeconfig context.
- Validate that the requested namespace exists in the current cluster before persisting.
- Update only the active context namespace in kubeconfig.
- Print `current namespace set to <namespace>` after a successful write.
- Return descriptive errors for missing namespaces, kubeconfig read failures, and kubeconfig write failures.

Acceptance coverage: AC-003, AC-004, AC-005, AC-006, AC-007, AC-008, AC-009

## Task 4: Add automated test coverage for `cmd/ns`

- Add tests for `ns get-current` success with an explicit namespace.
- Add tests for `ns get-current` falling back to `default` when the active context has no namespace.
- Add tests for `ns set-current` success, including persisted kubeconfig mutation for the active context only.
- Add tests for namespace validation failure to confirm kubeconfig is not mutated.
- Add tests for kubeconfig read/write error handling.
- Exercise inherited `--kubeconfig` behavior through a small root command tree with `Execute()`.

Acceptance coverage: AC-003, AC-006, AC-007, AC-008, AC-010

## Task 5: Update manual documentation

- Extend the namespace command section in `README.md` with `ns get-current` and `ns set-current` examples.
- Document the `getcurrent` and `setcurrent` aliases alongside the existing namespace aliases.

Acceptance coverage: AC-011

## Task 6: Validate the change

- Run `make test`.
- Confirm the implemented behavior still matches the clarified spec after tests pass.

Acceptance coverage: AC-010, AC-011
