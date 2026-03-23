# CLI Contract: `ctx list` and `ns list` refactor

## Purpose

Define the user-visible behavior that must remain stable while refactoring the internal list logic for `viewnode ctx list` and `viewnode ns list`.

## Command Contract

### `viewnode ctx list`

- Command name remains `ctx list`
- Alias remains `ctx ls`
- Reads contexts from the selected kubeconfig source
- Produces one output line per discovered context
- Sorts context names alphabetically before printing
- Marks the current context with `[*]`
- Marks non-current contexts with `[ ]`
- Returns an error and no partial list when kubeconfig loading fails

Example:

```text
[*] dev-cluster
[ ] staging-cluster
```

### `viewnode ns list`

- Command name remains `ns list`
- Alias remains `ns ls`
- Reads namespaces from the current cluster context
- Produces one output line per discovered namespace
- Sorts namespace names alphabetically before printing
- Marks the active namespace with `[*]`
- Marks non-active namespaces with `[ ]`
- Treats an empty namespace on the active kubeconfig context as `default`
- Returns an error and no partial list when namespace discovery fails

Example:

```text
[ ] default
[*] jenkins-onprem
[ ] kube-system
```

## Non-Goals

- No new flags, aliases, or output columns
- No change to command help text beyond internal cleanup if needed
- No change to namespace retrieval source or context retrieval source

## Verification Contract

- Existing command invocations continue to work unchanged
- Tests verify ordering, marker placement, default namespace fallback, and error propagation
- `make test` passes after the refactor
