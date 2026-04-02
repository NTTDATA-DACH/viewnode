# Contract: All-namespaces node view output

## Purpose

Define the user-visible behavior for `viewnode` when `--all-namespaces` is used and scheduled pods must be rendered beneath namespace grouping rows.

## Command Surface

### `viewnode --all-namespaces`

- Uses the existing root command and flag names unchanged
- Loads node and pod data through the current root workflow
- Prints existing summary lines for:
  - total pod count
  - unscheduled pod count
  - running node count with scheduled pod count
- Renders each runnable node as a tree section
- Renders namespace grouping rows beneath each node when scheduled pods are shown
- Orders namespace grouping rows alphabetically by namespace name within each node
- Shows the namespace grouping row even when a node contains pods from only one namespace
- Renders pod rows beneath their namespace heading without repeating the namespace inline on each pod line
- Keeps current pod ordering within each namespace group
- Preserves existing container inline/tree rendering beneath each pod row

## Output Shape

### Node section

```text
<node-prefix> <node-name> running <pod-count> pod(s) (<os>/<arch>/<runtime>[ | metrics])
<node-indent><namespace-prefix> <namespace>
<node-indent><namespace-child-indent><pod-prefix> <pod-name> (<phase>[/<time>][ | metrics]) ...
```

### Behavioral guarantees

- Namespace headings are nested directly under the node that owns the scheduled pods.
- Pod rows are nested directly under the namespace heading that matches each pod's namespace.
- Namespace membership remains unambiguous from the tree structure alone.
- The grouped namespace tree is used only for the all-namespaces path; non-grouped scoped output remains unchanged.

## Example

```text
$ viewnode -A
22 pod(s) in total
0 unscheduled pod(s)
2 running node(s) with 22 scheduled pod(s):
├── control-plane running 9 pod(s) (linux/arm64/containerd://2.2.0)
│   ├── kube-system
│   │   ├── coredns-a (running)
│   │   └── kube-proxy-a (running)
│   └── local-path-storage
│       └── local-path-provisioner-a (running)
└── worker running 13 pod(s) (linux/arm64/containerd://2.2.0)
    ├── camunda
    │   ├── connectors-a (running)
    │   └── zeebe-a (running)
    └── kube-system
        ├── kindnet-a (running)
        └── kube-proxy-b (running)
```

## Non-Goals

- No new CLI flags, aliases, or command names
- No change to pod-loading or namespace-loading API behavior
- No change to unscheduled pod rendering semantics
- No alphabetical reordering of pods inside a namespace group unless the existing pipeline already provides that order
- No change to scoped or single-namespace flat output when `--all-namespaces` is not active

## Verification Expectations

- Tests verify alphabetical namespace grouping within each node
- Tests verify a node with one namespace still renders a namespace heading
- Tests verify pod rows no longer repeat `<namespace>: ` inline in grouped all-namespaces output
- Tests verify container inline/tree rendering still works under grouped pods
- `README.md` documents the grouped all-namespaces output example
