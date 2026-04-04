# Data Model: Remove redundant namespace prefixes from single-namespace pod rows

## Overview

The feature does not add new storage or Kubernetes resources. It adjusts how the existing flat single-namespace render path presents pod rows while preserving the current header, grouping, and pod metadata behavior.

## Entities

### ViewNodeData

- **Purpose**: Root printout input assembled by `cmd/root.go` and rendered by `srv/view.go`.
- **Relevant fields**:
  - `Namespace`: Current namespace label shown for explicitly scoped runs
  - `NodeFilter`: Filter text used for empty-state messaging
  - `Nodes`: Ordered slice containing the unscheduled-pod bucket at index `0` followed by runnable nodes
  - `Config.ShowNamespaces`: Enables namespace-aware rendering
  - `Config.GroupPodsByNamespace`: Determines whether grouped namespace headings are used
  - `Config.ShowContainers`: Enables container details per pod
  - `Config.ShowTimes`: Enables pod start-time output
  - `Config.ShowMetrics`: Enables metrics output
  - `Config.ContainerViewType`: Chooses inline vs tree container rendering
- **Constraints**:
  - `Nodes` must be non-nil for printout
  - `Nodes[0]` continues to represent unscheduled pods
  - Single-namespace runs keep `ShowNamespaces` enabled while `GroupPodsByNamespace` remains false

### ViewPod

- **Purpose**: Represents the rendered pod row beneath a node.
- **Relevant fields**:
  - `Name`
  - `Namespace`
  - `Phase`
  - `StartTime`
  - `Metrics`
  - `Containers`
- **Constraints**:
  - `Namespace` remains available in memory even when it is not printed inline for single-namespace rows
  - `Name`, `Phase`, and optional details remain unchanged by the fix
  - Pod ordering remains stable

### Namespace Header

- **Purpose**: Top-level output line that communicates the explicit namespace scope for a run.
- **Fields**:
  - `Namespace value`
- **Constraints**:
  - Remains visible for single-namespace runs
  - Remains the canonical namespace label when the inline pod-row prefix is removed

### Display Mode

- **Purpose**: Determines how namespace information is shown.
- **States**:
  - `Single namespace flat output`
  - `Multi-namespace grouped output`
  - `All-namespaces grouped output`
- **Constraints**:
  - Only the single-namespace flat output changes in this feature
  - Other modes remain behaviorally unchanged

## Relationships

- `ViewNodeData` contains many `ViewNode` entries.
- Each runnable `ViewNode` contains many `ViewPod` entries.
- The single-namespace flat output path renders `ViewPod` rows directly beneath their node.
- The namespace header provides the namespace context for all rendered pod rows in the single-namespace flat path.

## Rendering Rules

- Single-namespace runs keep the namespace header visible.
- Single-namespace pod rows print the pod name and status without the inline `<namespace>: ` prefix.
- Multi-namespace and all-namespaces grouped rendering rules remain unchanged.
- Container, timing, and metrics output remain attached to the same pod rows.

## Validation Notes

- Tests should validate that single-namespace pod rows no longer include the inline namespace prefix.
- Tests should validate that the namespace header remains visible.
- Tests should validate that pod ordering, node grouping, and pod-row details remain unchanged.
- Tests should validate that other namespace display modes are unaffected.
