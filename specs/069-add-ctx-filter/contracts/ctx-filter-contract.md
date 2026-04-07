# CLI Contract: Context list filtering

## Purpose

Define the user-visible behavior for the new context filter flag and the clarified no-match messaging required by issue `#69`.

## Command Contract

### `viewnode ctx list`

- Command name remains `ctx list`
- Alias remains `ctx ls`
- Adds a local `--filter` flag with shorthand `-f`
- Reads contexts from the current kubeconfig using the existing raw kubeconfig discovery path
- Applies case-insensitive substring matching to context names when a filter is provided
- Treats an empty filter value as equivalent to omitting the filter
- Sorts displayed context names alphabetically before printing
- Marks the active displayed context with `[*]`
- Marks non-active displayed contexts with `[ ]`
- Returns an error and no partial list when kubeconfig loading fails
- Returns success with the exact no-match message `no contexts matched filter "<value>"` when filtering finds zero contexts
- Shows only matching contexts when a filter is present
- Does not display an active marker when the current context is filtered out and therefore absent from the displayed rows

Example:

```text
$ viewnode ctx list --filter prod
[ ] prod-admin
[*] prod-cluster
```

No-match example:

```text
$ viewnode ctx list --filter does-not-exist
no contexts matched filter "does-not-exist"
```

Current-context-filtered-out example:

```text
$ viewnode ctx list --filter staging
[ ] staging-cluster
```

## Non-Goals

- No change to `viewnode` root flag names or aliases
- No label-based, regex-based, or cluster/user-field filtering for contexts
- No change to kubeconfig selection, persistence, or `ctx get-current` / `ctx set-current` behavior
- No change to non-filtered context list output format

## Verification Contract

- Tests verify `--filter` and `-f` produce identical context results
- Tests verify case-insensitive context matching
- Tests verify empty filter input behaves like no filter
- Tests verify no-match outcomes print the exact success message
- Tests verify filtered output omits the active marker when the current context does not match the filter
- `README.md` documents the new context filter flag and example usage
- `make test` passes after the change
