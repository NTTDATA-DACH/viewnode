# Feature Specification: Namespace Filter Node Rendering

**Feature Branch**: `[061-fix-namespace-empty]`  
**Created**: 2026-04-04  
**Status**: Draft  
**Input**: User description: "bug: viewnode --namespace shows empty namespaces on nodes without matching pods"

## GitHub Issue Traceability

- **Issue Number**: 61
- **Issue URL**: https://github.com/NTTDATA-DACH/viewnode/issues/61
- **Issue Title**: bug: viewnode `--namespace` shows empty namespaces on nodes without matching pods

## User Scenarios & Testing *(mandatory)*

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.
  
  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - Show Only Matching Namespaces Per Node (Priority: P1)

As a user filtering `viewnode` output by one or more namespaces, I want each node to show only namespaces that actually contain matching pods on that node so the output reflects the real workload placement instead of empty sections.

**Why this priority**: The bug directly misrepresents filtered results, which makes the command output harder to trust and scan during cluster inspection.

**Independent Test**: Can be fully tested by running `viewnode` with namespace filters against nodes where some namespaces have no matching pods and verifying that empty namespace headings are not rendered while matching namespaces still appear.

**Acceptance Scenarios**:

1. **Given** a node has pods in namespace `kube-system` but no pods in namespace `camunda`, **When** the user runs `viewnode --namespace kube-system --namespace camunda`, **Then** that node shows only the `kube-system` section and does not render an empty `camunda` section.
2. **Given** a node has pods in multiple selected namespaces, **When** the user runs `viewnode` with those namespace filters, **Then** each matching namespace is rendered under that node with its matching pods.

---

### User Story 2 - Keep Nodes Visible Without Matching Pods (Priority: P2)

As a user filtering by namespace, I want nodes to remain visible even when they have no pods in the selected namespaces so I can still understand cluster node coverage and placement context.

**Why this priority**: The issue explicitly states that node visibility must remain unchanged by namespace filtering, which preserves a core behavior of the command.

**Independent Test**: Can be fully tested by running `viewnode` with namespace filters against a cluster where at least one node has zero matching pods and verifying that the node still appears without namespace sections beneath it.

**Acceptance Scenarios**:

1. **Given** a node has no pods in any selected namespace, **When** the user runs `viewnode --namespace <selected-namespace>`, **Then** the node is still shown in the output and no empty namespace section is rendered beneath it.

---

### User Story 3 - Preserve Filtered View Consistency (Priority: P3)

As a user comparing namespace-filtered output with other `viewnode` views, I want filtered namespace rendering to follow the same per-node inclusion rules as the all-namespaces view so the command behaves consistently across modes.

**Why this priority**: Consistency reduces confusion and makes it easier to understand filtered output by relying on already-familiar rendering rules.

**Independent Test**: Can be fully tested by comparing `viewnode -A` and namespace-filtered runs for the same cluster and confirming that only namespaces containing pods on a given node are rendered in both modes.

**Acceptance Scenarios**:

1. **Given** the all-namespaces view omits namespace headings with no pods for a node, **When** the user runs a namespace-filtered view for the same cluster state, **Then** the filtered view follows the same rule and omits empty namespace headings per node.

---

### Edge Cases

- A selected namespace exists in the cluster but has no pods on a specific node: that namespace is omitted for that node only.
- A node has no pods in any selected namespace: the node remains visible without any namespace sections beneath it.
- Multiple namespaces are selected and only a subset match on a node: only the matching subset is rendered for that node.
- None of the selected namespaces have pods on any displayed node: all relevant nodes remain visible, but no empty namespace sections are rendered.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST continue to show the same set of nodes that would be shown without namespace-based pod suppression, subject only to existing node-selection options such as node filters.
- **FR-002**: When namespace filters are provided, the system MUST render a namespace section under a node only if that node has at least one pod that matches both the node selection and the requested namespace filters.
- **FR-003**: The system MUST omit namespace headings that would otherwise contain no rendered pods for a given node.
- **FR-004**: When a node has pods from more than one selected namespace, the system MUST render each matching namespace section for that node.
- **FR-005**: When a node has no pods in any selected namespace, the system MUST still render the node entry and MUST NOT render empty namespace sections beneath it.
- **FR-006**: The namespace-filtered view MUST follow the same per-node namespace inclusion rule as the all-namespaces view: namespaces are shown only when that node contains pods in that namespace.
- **FR-007**: The change MUST preserve existing externally visible command behavior for node visibility, except for removing incorrect empty namespace sections.

### Key Entities *(include if feature involves data)*

- **Displayed Node**: A node that is included in the command output after applying existing node-selection rules.
- **Selected Namespace**: A namespace requested by the user through namespace-filtering options.
- **Matching Pod**: A pod that belongs to a selected namespace and is scheduled onto a displayed node.
- **Rendered Namespace Section**: The output grouping shown under a displayed node for a namespace that has one or more matching pods on that node.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In acceptance checks covering namespace-filtered output, 100% of rendered namespace sections contain at least one pod for the node under which they are shown.
- **SC-002**: In acceptance checks where one or more displayed nodes have zero matching pods for the selected namespaces, 100% of those nodes remain visible in the output.
- **SC-003**: For the issue reproduction scenario, the namespace-filtered output no longer shows any empty namespace headings on nodes without matching pods.
- **SC-004**: Comparing all-namespaces and namespace-filtered views for the same cluster state shows consistent per-node namespace inclusion behavior in all defined acceptance scenarios.

## Assumptions

- Namespace filtering already determines which pods are eligible to be shown, and this change only corrects how those filtered pods are grouped beneath nodes.
- The command's existing node-selection behavior, including optional node filtering, remains the authoritative source for which nodes appear.
- The scope of this feature is limited to correcting rendered output for namespace-filtered node views and does not include redesigning other output formats.
- Existing command semantics for all-namespaces output remain correct and can be used as the behavioral reference for this fix.
