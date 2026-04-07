# Feature Specification: Context List Filtering

**Feature Branch**: `069-add-ctx-filter`  
**Created**: 2026-04-07  
**Status**: Draft  
**Input**: User description: "GitHub issue #69: feat(ctx): add `--filter` flag to viewnode ctx list command. Add a `--filter` / `-f` flag to the `viewnode ctx list` command so that users can filter the displayed contexts by a substring match on the context name, mirroring the existing `viewnode ns list --filter` behavior."

## GitHub Issue Traceability

- **Issue Number**: 69
- **Issue URL**: https://github.com/NTTDATA-DACH/viewnode/issues/69
- **Issue Title**: feat(ctx): add `--filter` flag to viewnode ctx list command

## Clarifications

### Session 2026-04-07

- Q: What exact no-match message should `viewnode ctx list` use when a filter returns no contexts? → A: `no contexts matched filter "<value>"`
- Q: How should filtered output behave when the current context does not match the filter? → A: Show only matching contexts; if the current context does not match, no active marker appears.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Filter context results (Priority: P1)

As a `viewnode` user working with many kubeconfig contexts, I want to narrow `viewnode ctx list` results by a partial context name so I can find the relevant context faster.

**Why this priority**: This is the core capability requested in the issue and delivers immediate value without depending on any other workflow changes.

**Independent Test**: Can be fully tested by running `viewnode ctx list` with a filter value against a kubeconfig containing matching and non-matching context names, then verifying that only matching contexts are shown.

**Acceptance Scenarios**:

1. **Given** a kubeconfig with context names that include `prod`, **When** the user runs `viewnode ctx list --filter prod`, **Then** the command returns only contexts whose names contain `prod`.
2. **Given** a kubeconfig with context names that include `prod`, **When** the user runs `viewnode ctx list -f prod`, **Then** the command returns the same filtered results as the long flag.
3. **Given** a kubeconfig with context names containing `PROD` in varying letter case, **When** the user runs `viewnode ctx list --filter prod`, **Then** the command returns all contexts whose names match `prod` regardless of letter case.
4. **Given** a kubeconfig with available contexts, **When** the user provides an empty filter value, **Then** the command behaves the same as if no filter was supplied and shows the full context list.
5. **Given** a kubeconfig where the current context name does not match the provided filter, **When** the user runs `viewnode ctx list --filter <value>`, **Then** the command shows only matching contexts and none of the displayed rows are marked as active.

---

### User Story 2 - Understand empty filter results (Priority: P2)

As a `viewnode` user, I want a clear no-match message when my context filter returns nothing so I can tell the difference between an empty result set and a tool failure.

**Why this priority**: The issue asks for ergonomics matching the existing namespace filtering behavior, and clear no-match feedback is part of that user experience.

**Independent Test**: Can be fully tested by running the context list command with a filter that matches no contexts and verifying that the output explains that no contexts matched the filter.

**Acceptance Scenarios**:

1. **Given** a kubeconfig where no context names contain `does-not-exist`, **When** the user runs `viewnode ctx list --filter does-not-exist`, **Then** the command shows a clear no-match message instead of ambiguous empty output.
2. **Given** a kubeconfig where no context names match the provided filter, **When** the user runs the filtered command, **Then** the command completes successfully and shows a clear no-match message.

---

### User Story 3 - Learn the new option from documentation (Priority: P3)

As a `viewnode` user, I want the command reference and examples to explain context filtering so I can discover and use the new option without trial and error.

**Why this priority**: Documentation is secondary to command behavior, but it keeps the CLI discoverable and aligned with the repository's existing user-facing guidance.

**Independent Test**: Can be fully tested by reviewing the command reference and examples to confirm that the filter option, shorthand, and expected behavior are described clearly.

**Acceptance Scenarios**:

1. **Given** the updated documentation, **When** a user reads the `ctx list` command reference, **Then** they can see that `--filter` and `-f` are available and understand that filtering is based on substring matches.
2. **Given** the updated documentation, **When** a user looks for an example of context filtering, **Then** they can find an example showing a filtered context list command and the expected outcome.

### Edge Cases

- What happens when users provide an empty filter value? The command should treat it as no filter and show the full context list.
- How does the system handle context names that contain the filter text in a different letter case or repeated occurrences of the filter text?
- What happens when the filter returns no matches in a kubeconfig that otherwise contains contexts?
- What happens when the current context is excluded by the filter? The command should still succeed, show only matching contexts, and display no active marker in the filtered output.
- How does the system communicate when kubeconfig loading fails for reasons unrelated to filtering so users do not confuse an error with a no-match result?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST allow users to provide a context-name filter when running `viewnode ctx list`.
- **FR-002**: The system MUST support both `--filter` and `-f` as equivalent ways to provide the context-name filter.
- **FR-003**: The system MUST return only contexts whose names contain the provided filter text using case-insensitive substring matching.
- **FR-004**: The system MUST preserve the existing context list behavior when no filter is provided.
- **FR-004a**: The system MUST treat an empty filter value the same as no filter and return the full context list.
- **FR-005**: The system MUST present a clear, user-understandable message when no contexts match the provided filter.
- **FR-005b**: The system MUST use the exact no-match wording `no contexts matched filter "<value>"` when a context filter returns zero results.
- **FR-005a**: The system MUST treat a no-match result as a successful command outcome rather than an execution error.
- **FR-006**: The system MUST preserve the existing context list ordering and active-context marker behavior for contexts that remain after filtering.
- **FR-006a**: The system MUST NOT add a non-matching current context back into filtered output solely to preserve the active marker.
- **FR-006b**: The system MUST show no active marker when the current context does not satisfy the filter and is therefore absent from the filtered results.
- **FR-007**: The system MUST continue to return the existing descriptive error path and no partial list output when kubeconfig loading fails.
- **FR-008**: The system MUST document the context filtering option, its shorthand alias, and at least one example of expected filtered output behavior.
- **FR-009**: The system MUST include automated verification covering matching results, equivalent long and short flags, active-marker preservation, and the no-match outcome.

### Key Entities *(include if feature involves data)*

- **Context Result**: A kubeconfig context entry that may or may not be included in command output depending on whether its name satisfies the user-provided filter text.
- **Filter Input**: A user-provided text value used to decide which context results are shown.
- **No-Match Outcome**: A user-visible result state indicating that the command completed successfully but found no contexts satisfying the filter.

## Assumptions

- Filtering applies to context names only, because the issue describes substring matching on the context name and does not request filtering on clusters, users, or namespaces.
- Matching should mirror the existing `viewnode ns list --filter` behavior, including case-insensitive substring matching.
- An empty filter value should be handled as if the filter was omitted.
- A no-match result is a valid outcome and should not change the command to an error state.
- Existing `ctx list` behavior for alphabetical ordering, active-context indication, and kubeconfig-backed discovery remains authoritative unless explicitly changed by this feature.
- Filtered output should contain only matching contexts, even when that means the current context is not shown and no active marker appears.

## Dependencies

- Access to kubeconfig context data must already exist through the current `viewnode ctx list` workflow.
- The existing namespace filter behavior in issue `#53` remains the behavioral reference for filter semantics and user messaging.
- The command reference and examples in `README.md` remain the authoritative user-facing documentation for this CLI.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can narrow a context list to matching names with a single filtered command invocation.
- **SC-002**: In validation runs, the long and short filter flags produce the same result set for the same input, including names that differ only by letter case.
- **SC-003**: When no contexts match, users receive a message that clearly states no contexts matched the filter instead of ambiguous empty output.
- **SC-003b**: In validation runs, a no-match context filter prints the exact message `no contexts matched filter "<value>"` with the provided filter value substituted into the output.
- **SC-003a**: When no contexts match, the command still completes successfully so users can distinguish a valid empty result from a command failure.
- **SC-004**: Filtered output continues to preserve the same active-context marker convention and deterministic ordering used by the unfiltered context list.
- **SC-004a**: In validation runs, when the current context does not match the filter, the filtered output contains only matching contexts and no row is marked as active.
- **SC-005**: Documentation includes at least one context filtering example and clearly describes the available filter options.
