# Quickstart: Namespace tree view for all-namespaces node output

## Goal

Implement and validate the grouped namespace tree for scheduled pods in `viewnode --all-namespaces` output without changing the CLI surface or data-loading behavior.

## Steps

1. Update [srv/view.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go) to group scheduled pods by namespace during printout when namespace display is enabled.
2. Keep [cmd/root.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root.go) minimal, only preserving the existing signal that namespace-aware output is active for all-namespaces and multi-namespace views.
3. Extend [srv/view_test.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go) with grouped output coverage:
   - alphabetical namespace grouping
   - single-namespace grouping row
   - unchanged flat output when grouping is not active
   - compatibility with container rendering
4. Update [README.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/README.md) examples and command reference to show the grouped all-namespaces output.
5. Run `make test`.

## Suggested Validation Scenarios

### Grouped all-namespaces output

Run:

```bash
viewnode -A
```

Verify:

- scheduled pods appear beneath namespace headings under each node
- namespace groups are alphabetical within each node
- a node with only one namespace still shows that namespace heading
- pod rows do not repeat the namespace inline

### Container compatibility

Run:

```bash
viewnode -Ac
viewnode -Acb
```

Verify:

- container details still render beneath the correct pod row
- namespace grouping remains intact in both inline and tree container modes

### Scoped output unchanged

Run:

```bash
viewnode -n kube-system
```

Verify:

- scoped output keeps the current flat pod-row format
- the printed namespace label and empty-state behavior remain unchanged
