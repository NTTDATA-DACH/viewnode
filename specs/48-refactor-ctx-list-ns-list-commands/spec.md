# Feature Specification: Refactor `ctx list` and `ns list` commands

**Feature Branch**: `48-refactor-ctx-list-ns-list-commands`  
**Created**: 2026-03-15  
**Status**: Draft  
**Input**: User description: "https://github.com/NTTDATA-DACH/viewnode/issues/48 - Refactor `ctx list` and `ns list` commands by extracting shared list logic, removing duplicated loops, and introducing explicit list entry structs with active-state markers."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Keep context listing behavior stable after refactor (Priority: P1)

As a `viewnode` user, I want `ctx list` to keep showing the same ordered list of contexts and the active marker after the refactor so I can trust that internal cleanup does not change the command result.

**Why this priority**: `ctx list` is existing user-facing behavior that must remain correct while shared logic is extracted.

**Independent Test**: Can be fully tested by running `viewnode ctx list` against a kubeconfig with multiple contexts and verifying the same names, ordering, and single active marker as before.

**Acceptance Scenarios**:

1. **Given** a kubeconfig with multiple contexts and one current context, **When** the user runs `viewnode ctx list`, **Then** the command prints every context in deterministic order and marks only the current context as active.
2. **Given** a kubeconfig with one context, **When** the user runs `viewnode ctx list`, **Then** the command still prints that context in the existing output format without requiring namespace access.

---

### User Story 2 - Keep namespace listing behavior aligned with context listing (Priority: P1)

As a `viewnode` user, I want `ns list` to present namespace names and the active namespace marker through the same list-entry pattern used by `ctx list` so both commands stay consistent and easy to compare.

**Why this priority**: The issue explicitly targets removal of duplicated list logic and preservation of active-item behavior across both commands.

**Independent Test**: Can be fully tested by running `viewnode ns list` with a known namespace set and confirming the command returns each namespace in deterministic order with the active namespace marked.

**Acceptance Scenarios**:

1. **Given** a cluster response containing multiple namespaces and an active namespace, **When** the user runs `viewnode ns list`, **Then** the command prints every namespace in deterministic order and marks only the active namespace.
2. **Given** the active kubeconfig context does not set an explicit namespace, **When** the user runs `viewnode ns list`, **Then** the command treats `default` as the active namespace marker target.

---

### User Story 3 - Maintain predictable failure handling during the refactor (Priority: P2)

As a maintainer, I want list-command errors to continue surfacing through the existing command flow after shared logic is extracted so the refactor reduces duplication without hiding failures.

**Why this priority**: This protects operator trust and keeps the refactor bounded to structure rather than behavior changes.

**Independent Test**: Can be fully tested by forcing kubeconfig or namespace lookup failures and verifying each command still returns a descriptive error instead of partial or misleading output.

**Acceptance Scenarios**:

1. **Given** raw kubeconfig loading fails, **When** the user runs `viewnode ctx list`, **Then** the command returns the existing descriptive error path and produces no partial list output.
2. **Given** namespace discovery fails, **When** the user runs `viewnode ns list`, **Then** the command returns a descriptive error and produces no partial list output.

---

### Edge Cases

- What happens when the current context name exists in kubeconfig but the context map is otherwise empty?
- How does the system handle an active namespace that is not present in the discovered namespace list?
- What happens when duplicate removal changes internal code structure but the output order must remain stable for tests and users?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST preserve the existing `viewnode ctx list` command name, arguments, aliases, and output format while refactoring its internal list-building logic.
- **FR-002**: System MUST preserve the existing `viewnode ns list` command name, arguments, aliases, and output format while refactoring its internal list-building logic.
- **FR-003**: System MUST represent each context list result as a dedicated entry structure containing the context name and whether that context is active.
- **FR-004**: System MUST represent each namespace list result as a dedicated entry structure containing the namespace name and whether that namespace is active.
- **FR-005**: System MUST extract duplicated active-marker and ordered-list preparation logic used by `ctx list` and `ns list` into dedicated reusable functions where the shared behavior is the same.
- **FR-006**: System MUST continue to sort context and namespace names deterministically before rendering output.
- **FR-007**: System MUST determine the active context from the kubeconfig current-context value used by existing context commands.
- **FR-008**: System MUST determine the active namespace using the active kubeconfig context and treat an empty context namespace as `default`.
- **FR-009**: System MUST continue returning descriptive command errors when kubeconfig loading or namespace discovery fails.
- **FR-010**: System MUST keep the refactor within the current command organization patterns and avoid introducing alternative command flows or duplicate implementations.
- **FR-011**: Automated tests MUST cover the shared list-entry behavior for both commands, including active-state marking, deterministic ordering, default-namespace handling for namespace output, and error propagation.

### Key Entities *(include if feature involves data)*

- **Context List Entry**: A display-ready representation of one kubeconfig context with its context name and whether it matches the current context.
- **Namespace List Entry**: A display-ready representation of one namespace with its namespace name and whether it matches the active namespace for the current context.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In validation runs, `ctx list` and `ns list` continue to produce the same user-visible line format and one active marker per list item set before and after the refactor.
- **SC-002**: Automated coverage includes at least one successful listing case and one failure case for each affected command path before the change is considered ready.
- **SC-003**: Reviewers can identify one shared list-preparation flow used by both commands instead of separate duplicated loops for active-state rendering.
- **SC-004**: The refactor introduces no new documented user steps, flags, or command variants for listing contexts and namespaces.

## Assumptions

- Existing command names and aliases for `ctx list` and `ns list` remain unchanged because the issue requests refactoring rather than new CLI behavior.
- The current namespace for namespace-list marking should follow kubeconfig semantics and resolve an empty namespace to `default`, matching repository guidance.
- Deterministic ordering remains alphabetical because both current command implementations sort names before printing.
- Documentation updates are not required for this issue because the intended outcome is internal refactoring without user-facing command changes.
