# Data Model: Namespace tree view for all-namespaces node output

## Overview

The feature does not introduce persistent storage or new Kubernetes resources. It adds an intermediate rendering shape that groups a node's scheduled pods by namespace for CLI output.

## Entities

### ViewNodeData

- **Purpose**: Root printout input assembled by `cmd/root.go` and rendered by `srv/view.go`.
- **Relevant fields**:
  - `Namespace`: Current namespace label shown for explicitly scoped runs
  - `NodeFilter`: Filter text used for empty-state messaging
  - `Nodes`: Ordered slice containing the unscheduled-pod bucket at index `0` followed by runnable nodes
  - `Config.ShowNamespaces`: Enables namespace-aware rendering, including grouped all-namespaces output
  - `Config.ShowContainers`: Enables container details per pod
  - `Config.ContainerViewType`: Chooses inline vs tree container rendering
- **Constraints**:
  - `Nodes` must be non-nil for printout
  - `Nodes[0]` continues to represent unscheduled pods

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
  - `Namespace` is the canonical grouping key for scheduled pods when `ShowNamespaces` is active
  - Pod row content remains unchanged except for nesting under a namespace heading

### NamespaceGroup

- **Purpose**: Transient rendering structure used only during printout for a single node.
- **Fields**:
  - `Namespace`: Group heading label
  - `Pods`: Ordered pod rows belonging to that namespace
- **Constraints**:
  - Exactly one group per namespace per node
  - Groups are ordered alphabetically by `Namespace`
  - `Pods` preserve the existing order from the node's `Pods` slice
  - Groups are rendered even when a node has only one namespace

## Relationships

- `ViewNodeData` contains many `ViewNode` entries.
- Each runnable `ViewNode` contains many `ViewPod` entries.
- During printout, each runnable `ViewNode` is transformed into many `NamespaceGroup` entries keyed by `ViewPod.Namespace`.
- Each `NamespaceGroup` contains one or more `ViewPod` entries.

## Rendering Rules

- Namespace grouping is applied only for scheduled pods rendered under runnable nodes when `Config.ShowNamespaces` is true.
- Unscheduled pods in `Nodes[0]` are rendered exactly as they are today and are not namespace-grouped in the dedicated unscheduled section.
- Non-`--all-namespaces` output remains on the current flat pod-row path.
- Container output remains attached to each `ViewPod`, regardless of namespace grouping.

## Validation Notes

- Tests should validate alphabetical namespace group ordering.
- Tests should validate that a single-namespace node still renders a namespace heading.
- Tests should validate that grouped pod rows preserve pod identity and container details.
