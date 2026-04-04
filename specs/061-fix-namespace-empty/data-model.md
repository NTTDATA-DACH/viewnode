# Data Model: Fix empty namespace groups in namespace-filtered node output

## Overview

The feature does not add new storage or Kubernetes resources. It adjusts how the existing view-layer grouping model chooses which namespace sections to render for a node during grouped namespace-aware output.

## Entities

### ViewNodeData

- **Purpose**: Root printout input assembled by `cmd/root.go` and rendered by `srv/view.go`.
- **Relevant fields**:
  - `Namespace`: Current namespace label shown for explicitly scoped runs
  - `NodeFilter`: Filter text used for empty-state messaging
  - `Nodes`: Ordered slice containing the unscheduled-pod bucket at index `0` followed by runnable nodes
  - `Config.ShowNamespaces`: Enables namespace-aware rendering
  - `Config.GroupPodsByNamespace`: Enables tree grouping beneath namespace headings
  - `Config.SelectedNamespaces`: Canonical selected namespace set parsed by the root command for scoped grouped output
  - `Config.ShowContainers`: Enables container details per pod
  - `Config.ContainerViewType`: Chooses inline vs tree container rendering
- **Constraints**:
  - `Nodes` must be non-nil for printout
  - `Nodes[0]` continues to represent unscheduled pods
  - `SelectedNamespaces` remains trimmed and deduplicated before rendering

### ViewNode

- **Purpose**: Represents a runnable node section plus its scheduled pods.
- **Relevant fields**:
  - `Name`
  - `Os`
  - `Arch`
  - `ContainerRuntime`
  - `Metrics`
  - `Pods`
- **Constraints**:
  - Nodes with empty `Name` are not rendered as runnable node sections
  - Node summary counts continue to derive from `len(Pods)`
  - A rendered node may have zero pods matching the selected namespaces and must still remain visible

### ViewPod

- **Purpose**: Represents a rendered pod row and provides the namespace key used for grouping.
- **Relevant fields**:
  - `Name`
  - `Namespace`
  - `Phase`
  - `Condition`
  - `StartTime`
  - `Metrics`
  - `Containers`
- **Constraints**:
  - `Namespace` remains the canonical grouping key for scheduled pods
  - Pod content remains unchanged except for whether it appears beneath a rendered namespace heading
  - Pod ordering within a namespace group must remain stable

### NamespaceGroup

- **Purpose**: Transient rendering structure used during printout for a single node.
- **Fields**:
  - `Namespace`: Group heading label
  - `Pods`: Ordered pod rows belonging to that namespace
- **Constraints**:
  - Exactly one group per namespace per node
  - Groups are ordered alphabetically by `Namespace`
  - Each group contains at least one pod
  - Empty namespace groups are not rendered

### NamespaceFilterSelection

- **Purpose**: The canonical, already-parsed set of namespaces requested through `--namespace`.
- **Fields**:
  - `Values`: Ordered unique namespace names after trimming and deduplication
  - `Grouped`: Whether the selection should trigger grouped rendering
- **Constraints**:
  - Grouped rendering remains enabled only when more than one namespace is selected
  - The selection constrains which pod namespaces may appear in scoped grouped output
  - The selection no longer implies that every chosen namespace must render under every displayed node

## Relationships

- `ViewNodeData` contains many `ViewNode` entries.
- Each runnable `ViewNode` contains many `ViewPod` entries.
- During grouped scoped printout, each runnable `ViewNode` is transformed into zero or more `NamespaceGroup` entries derived from the node's matching pods.
- Each `NamespaceGroup` contains one or more `ViewPod` entries.

## Rendering Rules

- Namespace grouping is applied for grouped output only.
- Unscheduled pods in `Nodes[0]` are rendered exactly as they are today and are not namespace-grouped in the unscheduled section.
- Single-namespace scoped runs remain on the current flat pod-row path.
- In grouped scoped output, a namespace heading is rendered only if the current node has one or more pods in that namespace after filtering.
- If a displayed node has zero matching namespace groups, the node line still renders and no namespace headings are shown beneath it.
- Container output remains attached to each `ViewPod`, regardless of namespace grouping.

## Validation Notes

- Tests should validate alphabetical ordering of rendered namespace headings.
- Tests should validate omission of selected namespaces that have zero pods on a node.
- Tests should validate that nodes remain visible even when no namespace groups render beneath them.
- Tests should validate pod identity, pod order, and container output remain unchanged for rendered groups.
