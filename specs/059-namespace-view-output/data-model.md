# Data Model: Namespace tree view for namespace-filtered node output

## Overview

The feature does not introduce persistent storage or new Kubernetes resources. It extends the existing CLI printout context so scoped multi-namespace runs can render per-node namespace groups, including empty grouping rows for selected namespaces with no pods on a displayed node.

## Entities

### ViewNodeData

- **Purpose**: Root printout input assembled by `cmd/root.go` and rendered by `srv/view.go`.
- **Relevant fields**:
  - `Namespace`: Current namespace label shown for explicitly scoped runs
  - `NodeFilter`: Filter text used for empty-state messaging
  - `Nodes`: Ordered slice containing the unscheduled-pod bucket at index `0` followed by runnable nodes
  - `Config.ShowNamespaces`: Enables namespace-aware rendering
  - `Config.GroupPodsByNamespace`: Enables tree grouping beneath namespace headings
  - `Config.ShowContainers`: Enables container details per pod
  - `Config.ContainerViewType`: Chooses inline vs tree container rendering
  - `Config.SelectedNamespaces` or equivalent: Canonical selected namespace set used to render empty grouping rows for scoped output
- **Constraints**:
  - `Nodes` must be non-nil for printout
  - `Nodes[0]` continues to represent unscheduled pods
  - Selected namespaces must already be deduplicated and trimmed by the root command before printout

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
  - Pod row content remains unchanged except for nesting under a namespace heading
  - Pod ordering within a namespace group must match the current grouped `--all-namespaces` behavior

### NamespaceGroup

- **Purpose**: Transient rendering structure used during printout for a single node.
- **Fields**:
  - `Namespace`: Group heading label
  - `Pods`: Ordered pod rows belonging to that namespace
  - `Selected`: Whether the namespace comes from the active scoped selection, even if `Pods` is empty
- **Constraints**:
  - Exactly one group per namespace per node
  - Groups are ordered alphabetically by `Namespace`
  - `Pods` preserve the current grouped renderer's pod ordering behavior
  - A selected namespace group may contain zero pods for a displayed node

### NamespaceFilterSelection

- **Purpose**: The canonical, already-parsed set of namespaces requested through `--namespace`.
- **Fields**:
  - `Values`: Ordered unique namespace names after trimming and deduplication
  - `Grouped`: Whether the selection should trigger grouped rendering
- **Constraints**:
  - Grouped rendering is enabled only when more than one namespace is selected
  - Single-namespace scoped runs continue using the flat pod-row path

## Relationships

- `ViewNodeData` contains many `ViewNode` entries.
- Each runnable `ViewNode` contains many `ViewPod` entries.
- During grouped scoped printout, each runnable `ViewNode` is transformed into many `NamespaceGroup` entries keyed by the union of selected namespaces and pod namespaces.
- Each `NamespaceGroup` contains zero or more `ViewPod` entries.

## Rendering Rules

- Namespace grouping is applied for scoped output only when more than one namespace is selected.
- Unscheduled pods in `Nodes[0]` are rendered exactly as they are today and are not namespace-grouped in the dedicated unscheduled section.
- Single-namespace and non-namespace runs remain on the current flat pod-row path.
- Empty namespace groups are rendered only for explicitly selected namespaces in grouped scoped output.
- Container output remains attached to each `ViewPod`, regardless of namespace grouping.

## Validation Notes

- Tests should validate alphabetical namespace group ordering.
- Tests should validate that selected namespaces with no pods on a node still render empty grouping rows.
- Tests should validate that grouped scoped output preserves pod identity, pod order parity with `--all-namespaces`, and container details.
- Tests should validate that single-namespace scoped output stays flat.
