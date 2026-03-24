# PRD: Clearer EOF Error Guidance

## Overview

Improve the root `viewnode` failure path so scoped Kubernetes API `Get` requests that surface as EOF produce a clearer fatal message that hints proxy configuration may be involved, while preserving the original low-level error details and leaving non-EOF or out-of-scope EOF failures unchanged.

## Goals

- Help CLI users diagnose scoped EOF retrieval failures faster.
- Limit the new proxy hint to the intended cluster retrieval API `Get` failure pattern.
- Preserve original low-level error details in the final fatal output.
- Keep non-EOF and out-of-scope EOF behavior compatible with existing CLI messaging.
- Lock the behavior in with deterministic automated tests and repository validation.

## User Stories

### US-001: Classify scoped EOF retrieval failures
**Description:** Add a dedicated decorated error for scoped Kubernetes API `Get` EOF failures so the root command can distinguish them from other retrieval errors.

**Acceptance Criteria:**
- `srv/api_errors.go` defines a dedicated decorated error type for the scoped EOF case.
- The decorated error preserves the original wrapped error details.
- Non-EOF retrieval failures are not classified as the new scoped EOF error.

### US-002: Add proxy guidance to the root fatal path
**Description:** Update the root command so the scoped EOF error produces a clearer fatal message that mentions proxy configuration as a possible cause.

**Acceptance Criteria:**
- The root command handles the scoped EOF decorated error before the generic fatal fallback path.
- The user-facing fatal message states that proxy configuration may be the cause.
- The user-facing fatal message still includes the original low-level EOF detail.

### US-003: Keep unrelated failure paths unchanged
**Description:** Preserve existing unauthorized, forbidden, non-EOF, and out-of-scope EOF behavior so the new guidance appears only for the intended failure slice.

**Acceptance Criteria:**
- Unauthorized and forbidden retrieval failures continue to use their current message paths.
- An EOF outside the scoped retrieval pattern does not receive the proxy hint from this feature.
- Generic non-EOF retrieval failures continue to use the existing fatal wording path.

### US-004: Keep repeated scoped EOF failures deterministic
**Description:** Ensure repeated occurrences of the same scoped EOF failure continue to show the same guidance each time without introducing session-state suppression.

**Acceptance Criteria:**
- Repeated scoped EOF failures in the same run or session continue to emit the proxy hint each time.
- The repeated message remains stable and deterministic across occurrences.
- No new per-run or per-session suppression state is introduced.

### US-005: Add regression coverage for scoped and compatibility behavior
**Description:** Add automated tests that cover scoped EOF classification, root fatal output, and compatibility for non-scoped failures.

**Acceptance Criteria:**
- Tests cover successful recognition of the scoped EOF retrieval pattern.
- Tests cover root fatal output for the proxy-guidance message and retained low-level details.
- Tests cover non-EOF and out-of-scope EOF compatibility behavior.

### US-006: Validate the feature with repository-level proof
**Description:** Complete the feature with repository-aligned validation so the new behavior is proven rather than only described.

**Acceptance Criteria:**
- Quickstart and contract artifacts remain aligned with the implemented behavior.
- Repository validation includes `make test`.
- The final implementation matches the traceable behavior defined by the Speckit artifacts.

## Functional Requirements

- **FR-001**: The system must recognize cluster retrieval API `Get` failures surfaced to the user as EOF-related outcomes.
- **FR-002**: The system must classify the scoped EOF case with a dedicated decorated error that preserves the original wrapped error.
- **FR-003**: The root `viewnode` command must emit a clearer user-facing fatal message for the scoped EOF case.
- **FR-004**: The scoped EOF fatal message must explain that incorrect proxy configuration may be the cause without presenting it as certain or exclusive.
- **FR-005**: The scoped EOF fatal message must preserve the original low-level error details alongside the proxy guidance.
- **FR-006**: The system must avoid applying the proxy guidance to cluster retrieval failures that are not EOF-related.
- **FR-007**: The system must avoid applying the proxy guidance to EOF errors outside the targeted cluster retrieval API `Get` path.
- **FR-008**: Existing unauthorized, forbidden, metrics, and generic fatal or warning behavior must remain intact unless directly required by the scoped EOF feature.
- **FR-009**: Repeated scoped EOF failures must continue to emit the same guidance on every occurrence.
- **FR-010**: Automated tests must cover the scoped EOF path, preserved low-level details, and at least one non-EOF or out-of-scope EOF compatibility path.
- **FR-011**: Feature completion must include repository validation through `make test`.

## Non-Goals

- Adding new commands, flags, or aliases.
- Adding retry, backoff, or connectivity probing behavior.
- Detecting, editing, or auto-correcting proxy configuration.
- Changing unrelated non-EOF user-facing messages beyond preserving compatibility.
- Introducing a new shared error-handling framework outside the existing `srv` and root command flow.

## Implementation Notes

- Keep the implementation inside the existing repository structure, primarily in `srv/api_errors.go`, `srv/api.go`, and `cmd/root.go`.
- Reuse the current typed error-decoration pattern already used for unauthorized and forbidden retrieval failures.
- Keep the root command as the place where internal errors are translated into user-facing fatal messages.
- Scope the new guidance to Kubernetes API `Get` retrieval failures that surface as EOF; do not broaden it to every EOF string in the CLI.
- Preserve the original resource context in the fatal output, such as `loading and filtering of nodes failed`.
- Preserve the original low-level error detail rather than replacing it with only a friendly explanation.
- Keep repeated failures deterministic and stateless so watch-mode and repeated invocations remain predictable.
- Prefer regression coverage in the existing `cmd/root_test.go` and `srv` test files rather than introducing a new test structure.
- Treat README command-reference updates as unnecessary unless implementation adds a user-facing troubleshooting note, because the CLI surface does not change.

## Validation

- Run targeted Go tests covering scoped EOF classification and root fatal output behavior.
- Verify the scoped EOF fatal message includes both the proxy hint and the original low-level error detail.
- Verify non-EOF and out-of-scope EOF paths do not receive the new proxy guidance.
- Review the feature quickstart and contract artifacts to ensure they still match implementation behavior.
- Run `make test`.

## Traceability

- GitHub Issue: #27
- GitHub URL: https://github.com/NTTDATA-DACH/viewnode/issues/27
- GitHub Title: Add better error handling when Get call delivers EOF
- Feature Name: 27-handle-get-eof
- Feature Directory: /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof
- Spec Path: /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/spec.md
- Plan Path: /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/plan.md
- Tasks Path: /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/tasks.md
- PRD Path: /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/tasks/prd-27-handle-get-eof.md
- Source Status: derived from Speckit artifacts

## Assumptions / Open Questions

- The scoped EOF pattern can be identified reliably enough from the existing retrieval error shape without adding new Kubernetes client instrumentation.
- This PRD assumes the implementation will stay within the existing root retrieval flow and will not expand into broader networking diagnostics or proxy remediation.
