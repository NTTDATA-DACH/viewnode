# Feature Specification: Namespace Tree View for Namespace-Filtered Node Output

**Feature Branch**: `059-namespace-view-output`  
**Created**: 2026-04-04  
**Status**: Draft  
**Input**: User description: "GitHub issue #59: refactor: change output view for nodes when using `--namespace` flag."

## GitHub Issue Traceability

- **Issue Number**: 59
- **Issue URL**: https://github.com/NTTDATA-DACH/viewnode/issues/59
- **Issue Title**: refactor: change output view for nodes when using `--namespace` flag

## Related Prior Work

- This feature MUST match the namespace-grouping presentation already established for `--all-namespaces` in GitHub issue #57.
- Issue #57 remains the behavioral reference for how namespace headings and pod rows are nested in the tree output.

## Clarifications

### Session 2026-04-04

- Q: How should namespace groups be ordered within each node for `--namespace` output? → A: Alphabetical by namespace name.
- Q: Should selected namespaces be shown under every node even when a node has no pods in that namespace? → A: Yes, show every selected namespace under every displayed node.
- Q: How should pod rows be ordered within each namespace group? → A: Preserve the behavior of `--all-namespaces`.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Group filtered pods under namespaces (Priority: P1)

As a CLI user filtering the node view to specific namespaces, I want pods grouped beneath namespace headings within each node so I can scan the output using the same tree structure already used by the all-namespaces view.

**Why this priority**: This is the core behavior requested by the issue and delivers the main readability improvement for namespace-filtered output.

**Independent Test**: Can be fully tested by running the node view with `--namespace` set to multiple namespaces on a cluster where a node hosts pods from more than one selected namespace, then verifying that each node shows namespace headings with pod entries nested beneath them.

**Acceptance Scenarios**:

1. **Given** a node with scheduled pods from more than one selected namespace, **When** the user runs the node view with `--namespace`, **Then** the output shows one heading per selected namespace beneath that node and lists the matching pods beneath each heading.
2. **Given** a node with multiple pods from the same selected namespace, **When** the user runs the node view with `--namespace`, **Then** the namespace name appears once as a grouping row for that node instead of being repeated inline on every pod row.
3. **Given** a node with pods from several selected namespaces, **When** the user inspects the output, **Then** the tree visually nests namespace headings under the node and pod rows under the correct namespace heading.
4. **Given** a node with selected namespaces whose names sort differently from the user's flag order or first pod appearance, **When** the user runs the node view with `--namespace`, **Then** the namespace grouping rows appear in alphabetical order within that node.
5. **Given** the existing pod ordering behavior in `--all-namespaces`, **When** the user runs the node view with `--namespace`, **Then** pod rows within each namespace group follow that same ordering behavior.

---

### User Story 2 - Keep namespace-filtered output aligned with all-namespaces behavior (Priority: P2)

As a CLI user who switches between namespace-filtered and all-namespaces views, I want the output structures to match so I do not have to relearn how to read the tree when changing flags.

**Why this priority**: Consistency with the already-delivered all-namespaces behavior is central to the requested refactor and reduces user confusion.

**Independent Test**: Can be fully tested by comparing output from `--namespace` and `--all-namespaces` for equivalent node and namespace combinations and confirming that both use the same namespace-grouped tree structure.

**Acceptance Scenarios**:

1. **Given** the previously delivered all-namespaces tree behavior, **When** the user runs the node view with `--namespace`, **Then** the namespace-filtered output uses the same namespace-heading and child-pod presentation model.
2. **Given** a node that contains pods from only one selected namespace, **When** the user runs the node view with `--namespace`, **Then** the output still shows the namespace as its own grouping row above the pod entries for that node.
3. **Given** a user comparing filtered output with all-namespaces output, **When** both outputs include the same node and namespace, **Then** the grouping style is consistent even if the overall set of visible namespaces differs.
4. **Given** a user selects multiple namespaces, **When** a displayed node has no pods for one of the selected namespaces, **Then** that namespace still appears as a grouping row under the node with no pod entries beneath it.
5. **Given** a user compares the pods listed under the same namespace in `--namespace` and `--all-namespaces` output, **When** both commands display the same pod set for a node, **Then** the pod row ordering matches.

---

### User Story 3 - Preserve summary clarity while changing the tree layout (Priority: P3)

As a CLI user, I want the grouped namespace-filtered output to preserve the existing summary lines and node-level counts so I gain better readability without losing the overall cluster view.

**Why this priority**: The refactor is only successful if it improves pod presentation while keeping the rest of the command output familiar and trustworthy.

**Independent Test**: Can be fully tested by running the node view with `--namespace` before and after the change and confirming that total pod counts, unscheduled pod counts, and node-level pod counts remain clear and consistent with the grouped pod output.

**Acceptance Scenarios**:

1. **Given** a namespace-filtered node view that reports total pods, unscheduled pods, and node-level pod counts, **When** the grouped output is shown, **Then** those summary lines still appear and remain consistent with the displayed pods.
2. **Given** a node view previously showing namespace names inline on each pod row, **When** the grouped output is rendered, **Then** the namespace ownership of every pod remains unambiguous from the surrounding tree structure.
3. **Given** updated user-facing documentation or examples for this command behavior, **When** a user reviews them, **Then** they can understand that `--namespace` now uses the same grouped tree style as `--all-namespaces`.

### Edge Cases

- What happens when a selected namespace has only one pod on a node?
- What happens when a node contains pods from one selected namespace while other namespaces on that node are excluded by the filter?
- How does the output behave when the filter includes multiple namespaces but a node has pods from only a subset of them?
- What happens when a selected namespace has no pods on a displayed node?
- How does the command keep pod ownership clear when namespace names are no longer repeated inline on each pod row?
- How does the grouped layout remain readable when several selected namespaces appear on the same node?
- How are namespace grouping rows ordered when a node contains several selected namespaces?
- How are pod rows ordered within each namespace grouping?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST change the node view rendered with `--namespace` so pods are grouped beneath namespace headings within each node.
- **FR-002**: The system MUST show each namespace heading only once per node grouping and list all matching pod rows beneath that heading.
- **FR-003**: The system MUST preserve a clear visual tree hierarchy in which nodes contain namespace groupings and namespace groupings contain pod rows.
- **FR-003a**: The system MUST order namespace grouping rows alphabetically by namespace name within each node.
- **FR-004**: The system MUST keep the namespace membership of every displayed pod unambiguous in the grouped output.
- **FR-004a**: The system MUST preserve the same pod-row ordering behavior within each namespace grouping that is already used by `--all-namespaces`.
- **FR-005**: The system MUST use the same namespace-grouping presentation model for `--namespace` that is already used for `--all-namespaces`.
- **FR-006**: The system MUST show a namespace grouping row even when a node contains pods from only one selected namespace.
- **FR-006a**: The system MUST show a namespace grouping row for every selected namespace under each displayed node, even when that node has no displayed pods in that namespace.
- **FR-007**: The system MUST preserve existing summary information for total pods, unscheduled pods, and node-level pod counts when the grouped output is shown.
- **FR-008**: The system MUST limit displayed namespace groupings and pods to those included by the active `--namespace` filter.
- **FR-009**: The system MUST preserve existing node-view behavior when `--namespace` is not in use.
- **FR-010**: The system MUST document the grouped namespace-filtered node output in the repository's user-facing command reference or examples.
- **FR-011**: The system MUST include automated verification that the namespace-filtered node view groups pods by namespace without losing summary clarity or namespace ownership.

### Key Entities *(include if feature involves data)*

- **Node Output Group**: A rendered section representing one node and its summary information in the CLI tree output.
- **Namespace Grouping Row**: A rendered heading under a node that represents one filtered namespace containing one or more displayed pods for that node.
- **Pod Output Row**: A rendered pod entry displayed beneath its node and namespace grouping.
- **Namespace Filter Selection**: The user-supplied set of namespaces that determines which namespace groupings and pod rows appear in the node view.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In namespace-filtered runs that include pods from multiple namespaces on the same node, users can identify the namespace of any displayed pod from the tree structure without relying on a repeated inline namespace prefix.
- **SC-002**: In validation scenarios where a node hosts pods from multiple selected namespaces, the output shows one namespace grouping per displayed namespace under that node.
- **SC-002a**: In validation scenarios with multiple selected namespaces on a node, the namespace grouping rows appear in alphabetical order by namespace name.
- **SC-003**: In validation scenarios where a node hosts pods from only one selected namespace, the output still shows that namespace as its own grouping row above the pod entries.
- **SC-003a**: In validation scenarios where a displayed node has no pods for one or more selected namespaces, the output still shows those selected namespaces as empty grouping rows under that node.
- **SC-003b**: In validation scenarios where `--namespace` and `--all-namespaces` display the same pod set for a node and namespace, the pod rows appear in the same order.
- **SC-004**: Users reviewing namespace-filtered node output can still see total pod counts, unscheduled pod counts, and node-level pod counts without added ambiguity.
- **SC-005**: Repository documentation includes at least one example or reference that explains that `--namespace` uses the same grouped namespace tree style as `--all-namespaces`.

## Assumptions

- The requested output change applies specifically to the node view when `--namespace` is used.
- The namespace-filtered view should mirror the namespace-grouped presentation already delivered for `--all-namespaces`, unless the active namespace filter necessarily hides some namespaces.
- Namespace grouping rows should be ordered alphabetically within each node to keep the output predictable and aligned with the all-namespaces behavior.
- Every displayed node should render all user-selected namespaces as grouping rows, even if some of those namespaces have no displayed pods on that node.
- Pod rows inside each namespace grouping should follow the same ordering rule already established for `--all-namespaces`.
- Existing totals, node ordering, and pod membership remain authoritative and should continue to match the underlying cluster state.
- Grouping namespaces under a node is sufficient to convey pod ownership without repeating the namespace name inline on every pod row.
- Documentation updates belong in the same feature because this changes user-visible CLI output.

## Dependencies

- The behavior implemented for GitHub issue #57 defines the expected grouping model that this feature should reuse for namespace-filtered output.
- The current node view already gathers the pod and namespace information needed to determine which displayed pods belong under each node.
- README command listings and examples remain the authoritative user-facing documentation for CLI behavior in this repository.
