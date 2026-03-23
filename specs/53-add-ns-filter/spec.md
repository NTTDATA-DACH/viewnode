# Feature Specification: Namespace List Filtering

**Feature Branch**: `53-add-ns-filter`  
**Created**: 2026-03-23  
**Status**: Draft  
**Input**: User description: "GitHub issue #53: Add namespace name filtering to `ns list` via `--filter` flag."

## GitHub Issue Traceability

- **Issue Number**: 53
- **Issue URL**: https://github.com/NTTDATA-DACH/viewnode/issues/53
- **Issue Title**: Add namespace name filtering to `ns list` via `--filter` flag

## Clarifications

### Session 2026-03-23

- Q: Should namespace filtering use case-sensitive or case-insensitive substring matching? → A: Case-insensitive substring matching.
- Q: What should happen when the provided filter value is empty? → A: Treat it the same as no filter and show the full namespace list.
- Q: Should the command succeed or fail when no namespaces match the filter? → A: Succeed and show a clear no-match message.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Filter namespace results (Priority: P1)

As a CLI user working in clusters with many namespaces, I want to narrow `viewnode ns list` results by a partial name so I can find the namespaces I need faster.

**Why this priority**: This is the core requested capability and provides immediate value without depending on any other workflow changes.

**Independent Test**: Can be fully tested by running `viewnode ns list` with a filter value against a namespace set that includes matching and non-matching names, then verifying that only matching namespaces are shown.

**Acceptance Scenarios**:

1. **Given** a cluster with namespaces whose names include `kube`, **When** the user runs `viewnode ns list --filter kube`, **Then** the command returns only namespaces whose names contain `kube`.
2. **Given** a cluster with namespaces whose names include `kube`, **When** the user runs `viewnode ns list -f kube`, **Then** the command returns the same filtered results as the long flag.
3. **Given** a cluster with namespace names containing `KUBE` in varying letter case, **When** the user runs `viewnode ns list --filter kube`, **Then** the command returns all namespaces whose names match `kube` regardless of letter case.
4. **Given** a cluster with available namespaces, **When** the user provides an empty filter value, **Then** the command behaves the same as if no filter was supplied and shows the full namespace list.

---

### User Story 2 - Understand empty filter results (Priority: P2)

As a CLI user, I want a clear no-match message when my filter returns nothing so I can tell the difference between an empty result set and a tool failure.

**Why this priority**: The issue explicitly calls out confusing no-match feedback, and improving that message reduces user uncertainty when filters are used.

**Independent Test**: Can be fully tested by running the namespace list command with a filter that matches no namespaces and verifying that the output explains that no namespaces matched the filter.

**Acceptance Scenarios**:

1. **Given** a cluster where no namespace names contain `does-not-exist`, **When** the user runs `viewnode ns list --filter does-not-exist`, **Then** the command shows a clear no-match message instead of a generic empty-state message.
2. **Given** a cluster where no nodes or namespaces match an applied filter, **When** the user runs the relevant filtered command, **Then** the command explains that no resources matched the filter criteria.
3. **Given** a cluster where no namespace names match the provided filter, **When** the user runs the filtered command, **Then** the command completes successfully and shows a clear no-match message.

---

### User Story 3 - Learn the new behavior from documentation (Priority: P3)

As a CLI user, I want the command reference and examples to explain namespace filtering so I can discover and use the new option without trial and error.

**Why this priority**: Documentation is secondary to command behavior, but the issue requires it and it supports adoption of the new capability.

**Independent Test**: Can be fully tested by reviewing the command reference and examples to confirm that the filter option, shorthand, and expected behavior are described clearly.

**Acceptance Scenarios**:

1. **Given** the updated documentation, **When** a user reads the `ns list` command reference, **Then** they can see that `--filter` and `-f` are available and understand that filtering is based on substring matches.
2. **Given** the updated documentation, **When** a user looks for an example of namespace filtering, **Then** they can find an example showing a filtered namespace list command and the expected outcome.

### Edge Cases

- What happens when users expect a narrowed result but provide an empty filter value? The command should treat it as no filter and show the full namespace list.
- How does the system handle namespace names that contain the filter text in a different letter case or repeated occurrences of the filter text?
- What happens when the filter returns no matches in a context that otherwise has accessible namespaces?
- How does the system communicate when namespace discovery fails for reasons unrelated to filtering so users do not confuse an error with a no-match result?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST allow users to provide a namespace-name filter when running `viewnode ns list`.
- **FR-002**: The system MUST support both `--filter` and `-f` as equivalent ways to provide the namespace-name filter.
- **FR-003**: The system MUST return only namespaces whose names contain the provided filter text using case-insensitive substring matching.
- **FR-004**: The system MUST preserve the existing namespace list behavior when no filter is provided.
- **FR-004a**: The system MUST treat an empty filter value the same as no filter and return the full namespace list.
- **FR-005**: The system MUST present a clear, user-understandable message when no namespaces match the provided filter.
- **FR-005a**: The system MUST treat a no-match result as a successful command outcome rather than an execution error.
- **FR-006**: The system MUST apply the same clear no-match messaging rule to other filtered list output called out by the issue so users receive consistent feedback when filters return no results.
- **FR-007**: The system MUST document the namespace filtering option, its shorthand alias, and at least one example of expected filtered output behavior.
- **FR-008**: The system MUST include automated verification covering matching results, equivalent long and short flags, and the no-match outcome.

### Key Entities *(include if feature involves data)*

- **Namespace Result**: A namespace entry that may or may not be included in command output depending on whether its name satisfies the user-provided filter text.
- **Filter Input**: A user-provided text value used to decide which namespace results are shown.
- **No-Match Outcome**: A user-visible result state indicating that the command completed successfully but found no namespaces satisfying the filter.

## Assumptions

- Filtering applies to namespace names only, because the issue describes namespace name substring matching and does not request filtering on labels or other metadata.
- Matching is based on case-insensitive containment of the provided text within the namespace name.
- An empty filter value should be handled as if the filter was omitted.
- The clearer no-match messaging should distinguish a successful command with zero matches from operational errors.
- A no-match result is a valid outcome and should not change the command to an error state.

## Dependencies

- Access to namespace data in the selected context must already exist through the current `viewnode ns list` workflow.
- The command reference and examples in `README.md` remain the authoritative user-facing documentation for this CLI.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can narrow a namespace list to matching names with a single filtered command invocation.
- **SC-002**: In validation runs, the long and short filter flags produce the same result set for the same input, including names that differ only by letter case.
- **SC-003**: When no namespaces match, users receive a message that clearly states no namespaces matched the filter instead of a generic empty-state message.
- **SC-003a**: When no namespaces match, the command still completes successfully so users can distinguish a valid empty result from a command failure.
- **SC-004**: Documentation includes at least one namespace filtering example and clearly describes the available filter options.
