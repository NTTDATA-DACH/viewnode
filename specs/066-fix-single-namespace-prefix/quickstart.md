# Quickstart: Remove redundant namespace prefixes from single-namespace pod rows

## Goal

Verify that single-namespace `viewnode` output keeps the namespace header but removes the redundant inline namespace prefix from pod rows.

## Implementation Steps

1. Update `srv/view.go` so the flat single-namespace pod-row formatter prints only the pod name and status.
2. Update `srv/view_test.go` to replace the old `namespace: pod-name` expectation with the corrected single-namespace output.
3. Update `README.md` so the single-namespace examples and explanation match the corrected output.

## Validation

1. Run focused tests:

```bash
go test ./srv ./cmd
```

2. Run the repository validation required by the constitution:

```bash
make test
```

3. Manually compare single-namespace behavior:

```bash
viewnode -n kube-system
```

## Expected Outcomes

- The namespace header still shows the selected namespace.
- Pod rows no longer display the redundant `<namespace>: ` prefix in single-namespace output.
- Pod ordering, node grouping, counts, and optional row details remain unchanged.
- Multi-namespace and all-namespaces behavior remains unchanged.
