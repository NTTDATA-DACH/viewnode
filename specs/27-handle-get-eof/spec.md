# Feature Specification: Clearer EOF Error Guidance

**Feature Branch**: `27-handle-get-eof`  
**Created**: 2026-03-24  
**Status**: Draft  
**Input**: User description: "GitHub issue #27: Add better error handling when Get call delivers EOF."

## GitHub Issue Traceability

- **Issue Number**: 27
- **Issue URL**: https://github.com/NTTDATA-DACH/viewnode/issues/27
- **Issue Title**: Add better error handling when Get call delivers EOF

## Clarifications

### Session 2026-03-24

- Q: Which EOF failures should receive the new proxy guidance? → A: Only cluster retrieval/API `Get` failures that surface as EOF.
- Q: How strongly should the message attribute EOF to proxy settings? → A: Proxy misconfiguration may be the cause.
- Q: Should the original low-level error details remain visible? → A: Yes, keep the original error details and add the proxy hint alongside them.
- Q: Should repeated scoped EOF failures keep showing the proxy hint? → A: Yes, show the proxy hint every time the scoped EOF failure occurs.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Get actionable guidance for EOF failures (Priority: P1)

As a CLI user, I want unexpected EOF failures during cluster retrieval API `Get` requests to include helpful guidance so I can quickly recognize that my proxy configuration may be wrong.

**Why this priority**: This is the core problem described in the issue and directly reduces confusion during a failing workflow.

**Independent Test**: Can be fully tested by triggering a retrieval failure that currently surfaces as a raw EOF message and verifying that the user-facing error now includes proxy-related guidance.

**Acceptance Scenarios**:

1. **Given** a cluster retrieval API `Get` request fails with an EOF outcome, **When** the CLI reports the failure, **Then** the message explains that proxy configuration may be a likely cause.
2. **Given** a user sees the EOF-related failure, **When** they read the message, **Then** they can understand what likely went wrong without inspecting source code or logs.
3. **Given** a user sees the EOF-related failure, **When** the CLI reports it, **Then** the original low-level error context remains visible alongside the proxy guidance.

---

### User Story 2 - Preserve clarity for other failures (Priority: P2)

As a CLI user, I want only the relevant EOF-related failure to receive proxy guidance so unrelated errors keep their own clear meaning.

**Why this priority**: Expanding the guidance too broadly could make other failures harder to diagnose and would reduce trust in the output.

**Independent Test**: Can be fully tested by comparing EOF-related API `Get` retrieval failures with different failure types and verifying that only the intended EOF case includes the new proxy hint.

**Acceptance Scenarios**:

1. **Given** a cluster retrieval failure caused by something other than EOF, **When** the CLI reports the error, **Then** it does not incorrectly suggest that proxy settings are the likely cause.
2. **Given** an EOF error outside the cluster retrieval API `Get` path, **When** the CLI reports it, **Then** it does not automatically show the proxy guidance intended for this feature.
3. **Given** a retrieval failure that already has a meaningful user-facing message, **When** the CLI reports it, **Then** that message remains understandable and distinct from the EOF guidance.

---

### User Story 3 - Learn the behavior from documentation and tests (Priority: P3)

As a maintainer, I want the expected EOF guidance behavior to be documented in project artifacts so future changes do not quietly remove it.

**Why this priority**: The immediate user value is in the runtime message, but persistence through tests and recorded requirements keeps the fix reliable.

**Independent Test**: Can be fully tested by reviewing the specification and confirming that the expected EOF guidance behavior is captured clearly enough to support implementation and automated verification.

**Acceptance Scenarios**:

1. **Given** the feature artifacts for this issue, **When** a maintainer reviews them, **Then** they can see that EOF-related failures must guide users toward checking proxy configuration.
2. **Given** future work on error handling, **When** maintainers rely on the feature requirements, **Then** they can verify whether the EOF guidance behavior is still preserved.

### Edge Cases

- What happens when the reported EOF is transient and not caused by a proxy configuration problem?
- How does the system avoid implying that proxy settings are the only possible cause of every EOF-related retrieval failure?
- What happens when an EOF appears outside the cluster retrieval API `Get` path?
- What happens when the failing request already contains useful context, such as the target endpoint, and the user still needs that original context preserved?
- How does the system present guidance when repeated scoped EOF failures occur in the same session? The proxy hint should remain visible each time the scoped failure occurs.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST recognize cluster retrieval API `Get` failures surfaced to the user as EOF-related outcomes.
- **FR-002**: The system MUST present a clearer user-facing error message for EOF-related retrieval failures.
- **FR-003**: The EOF-related error message MUST explain that incorrect proxy configuration may be the cause of the failure, without presenting it as the only or certain cause.
- **FR-004**: The EOF-related error message MUST preserve the original low-level error details alongside the proxy guidance so the user can understand which retrieval action failed.
- **FR-005**: The system MUST avoid applying the proxy-configuration guidance to cluster retrieval failures that are not EOF-related.
- **FR-006**: The system MUST avoid applying the proxy-configuration guidance to EOF errors outside the cluster retrieval API `Get` path covered by this feature.
- **FR-007**: The system MUST keep non-EOF failure messages understandable and distinct from the EOF-specific guidance.
- **FR-008**: The system MUST include automated verification for the EOF-related guidance path and for at least one non-EOF or out-of-scope EOF path to confirm the guidance is not over-applied.
- **FR-009**: The system MUST continue showing the proxy guidance every time the scoped EOF-related cluster retrieval API `Get` failure occurs during a command run or session.

### Key Entities *(include if feature involves data)*

- **Retrieval Failure**: A user-visible failure produced while the CLI is requesting cluster data.
- **EOF-Related Failure**: A cluster retrieval API `Get` failure whose surfaced outcome ends unexpectedly and needs additional user guidance.
- **Proxy Guidance Message**: The user-facing explanation that helps users check whether proxy configuration is contributing to the EOF-related failure.

## Assumptions

- The most important user need is to receive a likely explanation for EOF-related retrieval failures rather than a bare low-level failure message.
- Proxy configuration is a common enough cause of this EOF-related outcome that it should be mentioned explicitly as guidance.
- Proxy guidance should be phrased as a likely explanation rather than a definitive diagnosis.
- The original failure context should remain visible so users can still tell what action failed.
- Repeated occurrences of the same scoped EOF failure should continue to show the same guidance rather than being suppressed.
- The change should improve diagnosis for this specific cluster retrieval failure pattern without redefining every network-related error as a proxy problem.

## Dependencies

- The existing CLI workflow must already surface retrieval failures to the user in a way that can be enhanced with clearer guidance.
- Project tests must be able to exercise both EOF-related and non-EOF error reporting paths.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: When an EOF-related cluster retrieval API `Get` failure occurs, users receive a message that points them toward checking proxy configuration instead of only showing a raw EOF outcome.
- **SC-002**: In validation runs, the intended EOF-related cluster retrieval failure, non-EOF retrieval failures, and out-of-scope EOF failures produce distinct messages so proxy guidance is only shown for the intended case.
- **SC-003a**: In validation runs, the EOF-related message still includes the original low-level failure details together with the added proxy guidance.
- **SC-003b**: In validation runs, repeated scoped EOF failures continue to show the same proxy guidance each time they occur.
- **SC-003**: Maintainers can verify the expected EOF guidance behavior from the feature requirements without needing to infer it from the original GitHub issue alone.
