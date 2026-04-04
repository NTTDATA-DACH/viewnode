# Contract: Single-namespace output without redundant pod-name prefixes

## Purpose

Define the corrected user-visible behavior for `viewnode` when exactly one namespace is selected.

## Command Surface

### `viewnode --namespace <ns>`

- Uses the existing root command and flag names unchanged
- Loads node and pod data through the current root workflow
- Prints the existing `namespace(s): <ns>` header for explicitly scoped runs
- Prints existing summary lines for:
  - total pod count
  - unscheduled pod count
  - running node count with scheduled pod count
- Renders each runnable node as a tree section
- Renders pod rows beneath each node without repeating the selected namespace inline in each pod row
- Preserves pod ordering, status text, and optional container, timing, and metrics details

## Output Shape

```text
namespace(s): <selected-namespace>
<summary lines>
<node-prefix> <node-name> running <pod-count> pod(s) (<os>/<arch>/<runtime>[ | metrics])
<node-indent><pod-prefix> <pod-name> (<phase>[/<time>][ | metrics]) ...
```

## Behavioral Guarantees

- The namespace header remains visible for the selected namespace.
- Pod rows in single-namespace output omit the redundant `<namespace>: ` prefix.
- Node grouping, counts, pod ordering, and unscheduled pod behavior remain unchanged.
- Multi-namespace and all-namespaces output behavior remains unchanged.

## Example

```text
$ viewnode -n kube-system
namespace(s): kube-system
10 pod(s) in total
0 unscheduled pod(s)
2 running node(s) with 10 scheduled pod(s):
├── camunda-platform-local-c8.8-control-plane running 8 pod(s) (linux/arm64/containerd://2.2.0)
│   ├── coredns-7d764666f9-2dwxj (running)
│   ├── coredns-7d764666f9-m9lbq (running)
│   ├── etcd-camunda-platform-local-c8.8-control-plane (running)
│   ├── kindnet-c4xnx (running)
│   ├── kube-apiserver-camunda-platform-local-c8.8-control-plane (running)
│   ├── kube-controller-manager-camunda-platform-local-c8.8-control-plane (running)
│   ├── kube-proxy-br7hr (running)
│   └── kube-scheduler-camunda-platform-local-c8.8-control-plane (running)
└── camunda-platform-local-c8.8-worker running 2 pod(s) (linux/arm64/containerd://2.2.0)
    ├── kindnet-w2mdf (running)
    └── kube-proxy-6f8jb (running)
```

## Non-Goals

- No new CLI flags, aliases, or command names
- No change to namespace header behavior
- No change to multi-namespace grouped output
- No change to all-namespaces grouped output
- No change to pod loading or namespace filtering semantics

## Verification Expectations

- Tests verify single-namespace pod rows no longer include `<namespace>: `
- Tests verify the namespace header remains visible
- Tests verify container and auxiliary pod-row details still render correctly
- Tests verify grouped namespace modes remain unchanged
- `README.md` documents the corrected single-namespace output
