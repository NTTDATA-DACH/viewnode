# Feature Specification: Single Namespace Pod Row Simplification

**Feature Branch**: `[066-fix-single-namespace-prefix]`  
**Created**: 2026-04-04  
**Status**: Draft  
**Input**: User description: "bug: single-namespace view redundantly prefixes pod names with namespace"

## GitHub Issue Traceability

- **Issue Number**: 66
- **Issue URL**: https://github.com/NTTDATA-DACH/viewnode/issues/66
- **Issue Title**: bug: single-namespace view redundantly prefixes pod names with namespace

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

### User Story 1 - Remove Redundant Prefixes In Single-Namespace Output (Priority: P1)

As a user viewing pods for exactly one selected namespace, I want pod rows to show only the pod name so the output is easier to scan and does not repeat information already present in the namespace header.

**Why this priority**: This is the directly reported bug and it affects the primary output users see when filtering to a single namespace.

**Independent Test**: Can be fully tested by running `viewnode -n kube-system` or an equivalent single-namespace command and verifying the header still shows the namespace while pod rows no longer start with `kube-system: `.

**Acceptance Scenarios**:

1. **Given** the user runs `viewnode -n kube-system`, **When** pod rows are rendered for nodes with matching pods, **Then** each pod row shows only the pod name and status without the `kube-system: ` prefix.
2. **Given** a single namespace is selected and pods exist on multiple nodes, **When** output is rendered, **Then** the namespace header remains visible once at the top and node grouping remains unchanged while pod rows omit the namespace prefix.

---

### User Story 2 - Preserve Existing Single-Namespace Layout Semantics (Priority: P2)

As a user already familiar with single-namespace output, I want the layout, counts, and grouping to stay the same except for the removed redundant namespace prefix so the fix feels like a cleanup rather than a redesign.

**Why this priority**: The issue asks for one targeted correction, so preserving adjacent behavior limits regression risk.

**Independent Test**: Can be fully tested by comparing output before and after the fix for a single-namespace run and confirming that only the inline namespace prefix is removed while counts, node sections, and pod ordering stay the same.

**Acceptance Scenarios**:

1. **Given** a single-namespace run currently shows the namespace header and grouped node output, **When** the fix is applied, **Then** the same header, node sections, counts, and pod ordering remain visible with only the inline namespace prefix removed from pod rows.

---

### User Story 3 - Keep Single-Namespace Output Consistent With Other Views (Priority: P3)

As a user switching between single-namespace, multi-namespace, and all-namespaces views, I want each mode to show only the namespace information needed for that mode so the output stays consistent and avoids redundant labels.

**Why this priority**: Cross-mode consistency helps users build predictable mental models of how namespace information is displayed.

**Independent Test**: Can be fully tested by comparing single-namespace output with multi-namespace and all-namespaces output and confirming that single-namespace rows omit inline namespace prefixes while other modes keep their established namespace display rules.

**Acceptance Scenarios**:

1. **Given** multi-namespace and all-namespaces views already use mode-appropriate namespace labeling, **When** the user runs a single-namespace view, **Then** the output avoids redundant inline namespace prefixes and remains consistent with the repository's existing display conventions.

---

### Edge Cases

- A single namespace is selected but no pods match: the existing empty or no-node behavior remains unchanged.
- A single namespace is selected and container details are enabled: container rows remain attached to the pod row and only the inline namespace prefix is removed.
- Multiple namespaces are selected: existing multi-namespace behavior remains unchanged.
- All namespaces are selected: existing all-namespaces grouped behavior remains unchanged.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: When exactly one namespace is selected, the system MUST render pod rows without prefixing the pod name with that namespace.
- **FR-002**: When exactly one namespace is selected, the system MUST continue to show the selected namespace in the existing namespace header.
- **FR-003**: The system MUST preserve node grouping, pod ordering, counts, and status rendering for single-namespace output, except for removing the redundant inline namespace prefix.
- **FR-004**: The system MUST preserve existing multi-namespace output behavior.
- **FR-005**: The system MUST preserve existing all-namespaces output behavior.
- **FR-006**: The system MUST preserve container, timing, metrics, and other pod-row details in single-namespace output after removing the redundant namespace prefix.
- **FR-007**: The change MUST remain limited to presentation behavior and MUST NOT require additional Kubernetes API calls or altered pod-loading semantics.

### Key Entities *(include if feature involves data)*

- **Single-Namespace Selection**: The command state in which exactly one namespace is requested through the namespace flag.
- **Namespace Header**: The top-level output line that shows the selected namespace value for explicitly scoped runs.
- **Pod Row**: The rendered output line for a pod beneath a node, including pod name, phase, and optional extra details.
- **Display Mode**: The current namespace presentation mode, such as single-namespace, multi-namespace, or all-namespaces.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In acceptance checks for single-namespace output, 100% of pod rows omit the redundant `<namespace>: ` prefix.
- **SC-002**: In acceptance checks for single-namespace output, 100% of runs still display the namespace header and the same node and pod counts as before the fix.
- **SC-003**: Multi-namespace and all-namespaces acceptance checks show no behavior change in namespace labeling.
- **SC-004**: For the issue reproduction scenario, the rendered output matches the expected example where pod names appear without the redundant namespace prefix.

## Assumptions

- The selected namespace is already represented elsewhere in the single-namespace output, so repeating it inline on every pod row adds no required information.
- Existing multi-namespace and all-namespaces behavior already uses distinct namespace display rules and remains the source of truth for those modes.
- The fix is limited to rendering behavior and does not require changes to command flags, API calls, or data retrieval.
- Existing output options such as container details, timing, and metrics remain supported on the same pod rows after the prefix is removed.
