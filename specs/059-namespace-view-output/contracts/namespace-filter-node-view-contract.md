# Contract: Namespace-filtered node view output

## Purpose

Define the user-visible behavior for `viewnode` when `--namespace` selects multiple namespaces and scheduled pods must be rendered beneath namespace grouping rows.

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
- Renders namespace grouping rows beneath each displayed node when more than one namespace is selected
- Orders namespace grouping rows alphabetically by namespace name within each node
- Shows a namespace grouping row for every selected namespace under each displayed node, even when that node has no displayed pods in that namespace
- Renders pod rows beneath their namespace heading without repeating the namespace inline on each pod line
- Keeps pod-row ordering within each namespace group aligned with grouped `--all-namespaces` behavior
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
- Selected namespaces with no displayed pods on a node still appear as empty grouping rows for that node.
- Namespace membership remains unambiguous from the tree structure alone.
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

### Empty selected namespace row

```text
$ viewnode --namespace default,team-a
3 pod(s) in total
0 unscheduled pod(s)
1 running node(s) with 3 scheduled pod(s):
└── worker-a running 3 pod(s) (linux/amd64)
    ├── default
    │   ├── api-0 (running)
    │   └── web-0 (running)
    └── team-a
```

## Non-Goals

- No new CLI flags, aliases, or command names
- No change to pod-loading or namespace-loading API behavior
- No change to unscheduled pod rendering semantics
- No change to single-namespace scoped output formatting
- No new pod-ordering rule beyond matching grouped `--all-namespaces`

## Verification Expectations

- Tests verify alphabetical namespace grouping within each node
- Tests verify selected namespaces with no pods still render empty grouping rows
- Tests verify pod rows no longer repeat `<namespace>: ` inline in grouped scoped output
- Tests verify grouped scoped pod ordering matches grouped `--all-namespaces`
- Tests verify single-namespace scoped output remains flat
- `README.md` documents the grouped scoped output example
