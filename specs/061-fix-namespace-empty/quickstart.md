# Quickstart: Fix empty namespace groups in namespace-filtered node output

## Goal

Verify that grouped namespace-filtered `viewnode` output omits empty namespace headings per node while continuing to show nodes with no matching pods.

## Implementation Steps

1. Update `srv/view.go` so grouped namespace rendering only emits namespace groups that contain one or more pods for the current node.
2. Keep `cmd/root.go` grouped-render activation unchanged unless tests show a configuration mismatch.
3. Update `srv/view_test.go` to replace empty-group expectations with omission behavior and add node-visibility regression coverage.
4. Update `README.md` to remove text and examples that claim empty selected namespaces render as grouping rows.

## Validation

1. Run focused tests:

```bash
go test ./srv ./cmd
```

2. Run the repository validation required by the constitution:

```bash
make test
```

3. Manually compare expected grouped output behaviors:

```bash
viewnode --namespace default,team-a
viewnode -A
```

## Expected Outcomes

- Matching namespace groups still render alphabetically beneath each node.
- Selected namespaces with zero matching pods on a node are not shown for that node.
- Nodes remain visible even when no namespace groups render beneath them.
- Single-namespace scoped output remains on the flat inline namespace path.
- README examples describe the corrected behavior.
