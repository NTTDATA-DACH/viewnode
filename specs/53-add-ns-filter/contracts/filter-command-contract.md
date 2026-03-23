# CLI Contract: Namespace filtering and filter-aware empty states

## Purpose

Define the user-visible behavior for the new namespace filter flag and the clearer no-match messaging required by issue `#53`.

## Command Contract

### `viewnode ns list`

- Command name remains `ns list`
- Alias remains `ns ls`
- Adds a local `--filter` flag with shorthand `-f`
- Reads namespaces from the current cluster context using the existing namespace discovery path
- Applies case-insensitive substring matching to namespace names when a filter is provided
- Treats an empty filter value as equivalent to omitting the filter
- Sorts displayed namespace names alphabetically before printing
- Marks the active namespace with `[*]`
- Marks non-active namespaces with `[ ]`
- Treats an empty namespace on the active kubeconfig context as `default`
- Returns an error and no partial list when namespace discovery fails
- Returns success with a clear no-match message when filtering finds zero namespaces

Example:

```text
$ viewnode ns list --filter kube
[ ] default
[ ] KUBE-public
[*] kube-system
```

No-match example:

```text
$ viewnode ns list --filter does-not-exist
no namespaces matched filter "does-not-exist"
```

### `viewnode --node-filter`

- Existing root command flag names remain unchanged
- Existing node filtering behavior remains intact unless needed to support clearer output wording
- When node filtering yields no visible nodes, the command returns success and prints a message that clearly states no nodes matched the filter instead of the generic empty-state wording

Example:

```text
$ viewnode --node-filter missing-node
no nodes matched filter "missing-node"
```

## Non-Goals

- No change to `viewnode` root flag names or aliases
- No label-based or regex-based namespace filtering
- No change to namespace retrieval source or kubeconfig selection behavior
- No change to non-filtered namespace list output format

## Verification Contract

- Tests verify `--filter` and `-f` produce identical namespace results
- Tests verify case-insensitive namespace matching
- Tests verify empty filter input behaves like no filter
- Tests verify no-match outcomes print clear success messages for both namespace and node filtering paths
- `README.md` documents the new namespace filter flag and example usage
- `make test` passes after the change
