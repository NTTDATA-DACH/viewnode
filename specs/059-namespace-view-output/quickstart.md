# Quickstart: Namespace tree view for namespace-filtered node output

## Goal

Implement and validate the grouped namespace tree for multi-namespace `viewnode --namespace ...` output without changing the CLI surface or data-loading behavior.

## Steps

1. Update [cmd/root.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root.go) to pass the parsed selected namespace set into the view context and enable grouped rendering only for multi-namespace scoped runs.
2. Update [srv/view.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go) to build namespace groups from the selected namespace set plus the node's scheduled pods.
3. Extend [srv/view_test.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go) with grouped scoped output coverage:
   - alphabetical namespace grouping
   - empty selected namespace rows
   - pod-order parity with grouped `--all-namespaces`
   - unchanged single-namespace flat output
   - compatibility with container rendering
4. Extend [cmd/root_test.go](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root_test.go) with root-context coverage for grouped scoped activation.
5. Update [README.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/README.md) examples and command reference to show the grouped scoped output.
6. Run `make test`.

## Suggested Validation Scenarios

### Grouped multi-namespace scoped output

Run:

```bash
viewnode --namespace default,jenkins-onprem
```

Verify:

- scheduled pods appear beneath namespace headings under each node
- namespace groups are alphabetical within each node
- pod rows do not repeat the namespace inline
- pod ordering matches the grouped `--all-namespaces` behavior

### Empty selected namespace row

Run:

```bash
viewnode --namespace default,team-a
```

Verify:

- every displayed node shows both selected namespaces as headings
- a selected namespace with no displayed pods still appears as an empty grouping row

### Container compatibility

Run:

```bash
viewnode --namespace default,team-a -c
viewnode --namespace default,team-a -cb
```

Verify:

- container details still render beneath the correct pod row
- namespace grouping remains intact in both inline and tree container modes

### Single-namespace scoped output unchanged

Run:

```bash
viewnode -n kube-system
```

Verify:

- scoped output keeps the current flat `<namespace>: <pod>` format
- the printed namespace label and empty-state behavior remain unchanged
