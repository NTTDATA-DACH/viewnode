# Feature Specification: Namespace Tree View for All-Namespaces Node Output

**Feature Branch**: `57-all-namespaces-tree-view`  
**Created**: 2026-04-01  
**Status**: Draft  
**Input**: User description: "GitHub issue #57: refactor: change output view for nodes when using `--all-namespaces` flag."

## GitHub Issue Traceability

- **Issue Number**: 57
- **Issue URL**: https://github.com/NTTDATA-DACH/viewnode/issues/57
- **Issue Title**: refactor: change output view for nodes when using `--all-namespaces` flag

## Clarifications

### Session 2026-04-01

- Q: How should namespace groups be ordered within each node when `--all-namespaces` is enabled? → A: Alphabetical within each node.
- Q: Should `--all-namespaces` output still show a namespace grouping row when a node has pods from only one namespace? → A: Yes, always show the namespace grouping row.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Group node pods by namespace (Priority: P1)

As a CLI user running the node view across all namespaces, I want pods grouped under namespace headings within each node so I can scan multi-namespace output without reading the namespace name on every pod row.

**Why this priority**: This is the core change requested by the issue and directly improves readability for the main `--all-namespaces` workflow.

**Independent Test**: Can be fully tested by running the node view with `--all-namespaces` in a cluster where a node hosts pods from multiple namespaces and verifying that each node shows namespace headings with child pod entries beneath them.

**Acceptance Scenarios**:

1. **Given** a node with scheduled pods from more than one namespace, **When** the user runs the node view with `--all-namespaces`, **Then** the output shows namespace headings beneath that node and lists the relevant pods beneath each heading.
2. **Given** a node with multiple pods from the same namespace, **When** the user runs the node view with `--all-namespaces`, **Then** the namespace name appears once as a grouping row for that node instead of being repeated inline on every pod row.
3. **Given** a node with pods from several namespaces, **When** the user inspects the tree output, **Then** the namespace grouping remains visually nested under the node and the pod rows remain visually nested under their namespace heading.
4. **Given** a node with pods from namespaces whose names sort differently from their first pod appearance, **When** the user runs the node view with `--all-namespaces`, **Then** the namespace grouping rows appear in alphabetical order within that node.
5. **Given** a node with pods from only one namespace, **When** the user runs the node view with `--all-namespaces`, **Then** the output still shows that namespace as a grouping row above the pod entries for that node.

---

### User Story 2 - Preserve readable totals and node summaries (Priority: P2)

As a CLI user, I want the grouped output to keep the existing node summaries and pod totals understandable so I gain the new grouping without losing the overall cluster view I already rely on.

**Why this priority**: The grouping change is valuable only if the rest of the node view remains familiar and easy to interpret.

**Independent Test**: Can be fully tested by comparing a run with and without the grouping change and confirming that total pod counts, unscheduled pod counts, node headings, and pod membership remain understandable and consistent.

**Acceptance Scenarios**:

1. **Given** a node view that reports total pods, unscheduled pods, and node-level pod counts, **When** the user runs the command with `--all-namespaces`, **Then** those summary lines still appear and remain consistent with the grouped pod output.
2. **Given** a node that contains pods from exactly one namespace, **When** the user runs the command with `--all-namespaces`, **Then** the output still uses the namespace grouping structure without obscuring the node summary.
3. **Given** a pod list that previously showed the namespace inline, **When** the grouped output is rendered, **Then** the namespace association for every pod remains unambiguous from the surrounding tree structure.

---

### User Story 3 - Learn the grouped output from docs and examples (Priority: P3)

As a CLI user, I want the command documentation and examples to reflect the grouped all-namespaces output so I know what to expect when using the flag.

**Why this priority**: Documentation is secondary to the output change itself, but the repository expects user-facing command changes to be documented in the same workflow.

**Independent Test**: Can be fully tested by reviewing the command reference and examples to confirm they describe the grouped namespace tree behavior for all-namespaces node output.

**Acceptance Scenarios**:

1. **Given** the updated documentation, **When** a user reads the examples for node output across all namespaces, **Then** they can see that namespaces appear as grouping rows in the tree rather than as repeated inline prefixes.
2. **Given** the updated documentation, **When** a user looks up the command behavior for `--all-namespaces`, **Then** they can understand how pods are organized beneath nodes in multi-namespace output.

### Edge Cases

- What happens when a node has pods from only one namespace while `--all-namespaces` is enabled?
- What happens when several namespaces on the same node each contain only one pod?
- How does the output behave when a namespace heading would otherwise have no visible pod children after filtering or other view constraints?
- How does the command keep namespace ownership clear for every pod while preserving existing cluster and node summary lines?
- How are namespace grouping rows ordered when a node contains pods from many namespaces?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST change the node view rendered with `--all-namespaces` so pods are grouped beneath namespace headings within each node.
- **FR-002**: The system MUST show each namespace heading only once per node grouping and list all matching pod rows beneath that heading.
- **FR-003**: The system MUST preserve a clear visual tree hierarchy in which nodes contain namespace groupings and namespace groupings contain pod rows.
- **FR-003a**: The system MUST order namespace grouping rows alphabetically by namespace name within each node.
- **FR-004**: The system MUST keep the namespace membership of every pod unambiguous in the grouped output.
- **FR-004a**: The system MUST show a namespace grouping row even when a node contains pods from only one namespace.
- **FR-005**: The system MUST preserve existing summary information for total pods, unscheduled pods, and node-level pod counts when the grouped output is shown.
- **FR-006**: The system MUST preserve the existing non-grouped node view behavior when `--all-namespaces` is not in use.
- **FR-007**: The system MUST produce grouped output that remains readable when a node contains pods from one namespace or many namespaces.
- **FR-008**: The system MUST document the grouped all-namespaces node output in the README command reference or examples.
- **FR-009**: The system MUST include automated verification that the all-namespaces node view groups pods by namespace without losing node or pod clarity.

### Key Entities *(include if feature involves data)*

- **Node Output Group**: A rendered section representing one node and its associated summary information in the CLI tree output.
- **Namespace Grouping Row**: A rendered heading under a node that represents one namespace containing one or more scheduled pods for that node.
- **Pod Output Row**: A rendered pod entry displayed beneath its node and namespace grouping.

## Assumptions

- The requested output change applies specifically to the node view when `--all-namespaces` is enabled.
- Existing totals, node ordering, and pod membership remain authoritative and should continue to match the underlying cluster state.
- Grouping namespaces under a node is sufficient to convey pod ownership without repeating the namespace name inline on every pod row.
- Namespace grouping rows should be ordered alphabetically within each node to keep the output predictable.
- The all-namespaces view should use the same grouping structure for one namespace or many namespaces on a node.
- Documentation updates belong in the same feature because this changes user-visible CLI output.

## Dependencies

- The current node view already gathers pod and namespace information needed to determine which pods belong under each node.
- README command listings and examples remain the authoritative user-facing documentation for CLI behavior in this repository.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In all-namespaces runs, users can identify the namespace of any displayed pod from the tree structure without relying on a repeated inline namespace prefix on each pod row.
- **SC-002**: In validation scenarios where a node hosts pods from multiple namespaces, the output shows one namespace grouping per namespace under that node.
- **SC-002a**: In validation scenarios with multiple namespaces on a node, the namespace grouping rows appear in alphabetical order by namespace name.
- **SC-002b**: In validation scenarios where a node hosts pods from only one namespace, the output still shows that namespace as its own grouping row.
- **SC-003**: Users reviewing all-namespaces node output can still see total pod counts, unscheduled pod counts, and node-level pod counts without added ambiguity.
- **SC-004**: Documentation includes at least one example or reference that explains the grouped namespace tree behavior for all-namespaces node output.
