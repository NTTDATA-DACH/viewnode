# Contract: Namespace-filtered node view output without empty namespace groups

## Purpose

Define the corrected user-visible behavior for `viewnode` when namespace-filtered grouped output is rendered beneath nodes.

## Command Surface

### `viewnode --namespace <ns1,ns2,...>`

- Uses the existing root command and flag names unchanged
- Loads node and pod data through the current root workflow
- Prints the existing `namespace(s): ...` label for explicitly scoped runs
- Prints existing summary lines for:
  - total pod count
  - unscheduled pod count
  - running node count with scheduled pod count
- Renders each runnable node as a tree section
- Renders namespace grouping rows beneath a displayed node only for namespaces that actually have matching pods on that node
- Orders rendered namespace grouping rows alphabetically within each node
- Renders pod rows beneath their namespace heading without repeating the namespace inline on each pod line
- Keeps pod-row ordering within each namespace group aligned with the existing grouped renderer behavior
- Keeps nodes visible even when none of the selected namespaces have matching pods on that node
- Preserves existing container inline/tree rendering beneath each pod row

## Output Shape

### Node section with matching namespaces

```text
<node-prefix> <node-name> running <pod-count> pod(s) (<os>/<arch>/<runtime>[ | metrics])
<node-indent><namespace-prefix> <namespace>
<node-indent><namespace-child-indent><pod-prefix> <pod-name> (<phase>[/<time>][ | metrics]) ...
```

### Node section with no matching namespaces

```text
<node-prefix> <node-name> running <pod-count> pod(s) (<os>/<arch>/<runtime>[ | metrics])
```

## Behavioral Guarantees

- Namespace headings are nested directly under the node that owns the matching scheduled pods.
- Pod rows are nested directly under the namespace heading that matches each pod's namespace.
- Namespace headings with zero pods on the current node are omitted.
- Nodes remain visible even when no namespace headings are rendered beneath them.
- Single-namespace scoped output remains on the current flat `<namespace>: <pod>` path.

## Example

```text
$ viewnode --namespace default,jenkins-onprem
8 pod(s) in total
0 unscheduled pod(s)
3 running node(s) with 8 scheduled pod(s):
├── node-a running 2 pod(s) (linux/amd64)
│   ├── default
│   │   └── docker-in-the-cloud-a (running)
│   └── jenkins-onprem
│       └── docker-in-the-cloud-b (running)
├── node-b running 2 pod(s) (linux/amd64)
│   ├── default
│   │   └── docker-in-the-cloud-c (running)
│   └── jenkins-onprem
│       └── docker-in-the-cloud-d (running)
└── node-c running 4 pod(s) (linux/amd64)
    ├── default
    │   ├── docker-in-the-cloud-e (running)
    │   └── liveness-test-4-boom (failed)
    └── jenkins-onprem
        ├── docker-in-the-cloud-f (running)
        └── docker-in-the-cloud-g (running)
```

### Node remains visible without matching namespace groups

```text
$ viewnode --namespace default,team-a
3 pod(s) in total
0 unscheduled pod(s)
2 running node(s) with 3 scheduled pod(s):
├── worker-a running 3 pod(s) (linux/amd64)
│   └── default
│       ├── api-0 (running)
│       └── web-0 (running)
└── worker-b running 0 pod(s) (linux/amd64)
```

## Non-Goals

- No new CLI flags, aliases, or command names
- No change to pod-loading or namespace-loading API behavior
- No change to unscheduled pod rendering semantics
- No change to single-namespace scoped output formatting
- No new ordering rule beyond the existing grouped renderer's behavior

## Verification Expectations

- Tests verify empty selected namespace groups are not rendered
- Tests verify nodes still render when no namespace groups match
- Tests verify matching namespace groups remain alphabetical
- Tests verify pod rows do not repeat `<namespace>: ` inline in grouped scoped output
- Tests verify single-namespace scoped output remains flat
- `README.md` documents the corrected grouped scoped output behavior
