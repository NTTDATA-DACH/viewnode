# PRD: Namespace List Filtering

## Overview

Add namespace-name filtering to `viewnode ns list` through a local `--filter` / `-f` flag, using case-insensitive substring matching, while also replacing generic filtered empty-state output with clearer success-path no-match messages for both namespace listing and the existing root node-filter flow.

## Goals

- Let users narrow `viewnode ns list` results by namespace name with a local filter flag.
- Support both `--filter` and `-f` for namespace filtering without changing the existing root node-filter contract.
- Preserve current namespace ordering and active-namespace marker behavior.
- Treat an empty namespace filter as equivalent to omitting the filter.
- Return clear success-path no-match messages for filtered namespace and node output.
- Update manual CLI documentation in `README.md`.
- Keep validation deterministic through automated tests and `make test`.

## User Stories

### US-001: Add namespace filter flag wiring
**Description:** Add a local namespace filter flag to `viewnode ns list` so users can request filtered namespace output without changing root command semantics.

**Acceptance Criteria:**
- `viewnode ns list --filter <value>` is accepted by the `ns list` command.
- `viewnode ns list -f <value>` is accepted by the `ns list` command.
- The new namespace filter flag does not replace or rename the existing root `--node-filter` flag.

### US-002: Apply case-insensitive namespace filtering
**Description:** Filter namespace results by case-insensitive substring matching while preserving existing ordering and active marker behavior.

**Acceptance Criteria:**
- `viewnode ns list --filter kube` returns only namespaces whose names contain `kube`, regardless of letter case.
- Filtered namespace output remains alphabetically ordered.
- The active namespace marker continues to follow kubeconfig semantics, including treating an empty active namespace as `default`.

### US-003: Treat empty namespace filters as no filter
**Description:** Normalize empty filter input so a blank filter does not create surprising or failing behavior.

**Acceptance Criteria:**
- An empty namespace filter value produces the same namespace output as running `viewnode ns list` without a filter.
- Empty filter handling does not turn the command into an error state.
- Existing unfiltered namespace listing behavior remains unchanged.

### US-004: Improve filtered no-match messaging
**Description:** Replace generic filtered empty-state text with clear messages that explain no resources matched the filter while keeping command execution successful.

**Acceptance Criteria:**
- When `viewnode ns list` finds no namespaces matching the filter, it prints a clear no-match message and exits successfully.
- When the root node-filter flow finds no nodes matching the filter, it prints a clear no-match message instead of `no nodes to display...`.
- Filtered no-match outcomes are distinguishable from namespace or cluster access failures.

### US-005: Add regression coverage for filter behavior
**Description:** Add automated tests that lock in the namespace filter contract, filtered empty-state messages, and repository-specific Cobra behavior.

**Acceptance Criteria:**
- Tests cover `--filter` and `-f` equivalence for `viewnode ns list`.
- Tests cover case-insensitive namespace matching and empty-filter fallback.
- Tests cover filtered no-match messaging for namespace and node flows.
- Relevant Cobra command tests execute through a small root command tree where inherited flags matter.

### US-006: Update CLI documentation
**Description:** Document the new namespace filtering behavior and examples in the manually maintained README.

**Acceptance Criteria:**
- `README.md` documents `viewnode ns list --filter` and `viewnode ns list -f`.
- `README.md` includes at least one filtered namespace example.
- `README.md` reflects the clearer filtered no-match behavior where user-facing output examples are relevant.

## Functional Requirements

- **FR-001**: `viewnode ns list` must accept a local `--filter` flag.
- **FR-002**: `viewnode ns list` must accept `-f` as the shorthand alias for the local namespace filter flag.
- **FR-003**: Namespace filtering must use case-insensitive substring matching against namespace names.
- **FR-004**: Filtered namespace output must preserve alphabetical ordering.
- **FR-005**: Filtered namespace output must preserve active-namespace marking and continue to treat an empty active context namespace as `default`.
- **FR-006**: An empty namespace filter value must behave the same as omitting the namespace filter.
- **FR-007**: A namespace no-match result must return success and print a clear message stating that no namespaces matched the filter.
- **FR-008**: A root node-filter no-match result must return success and print a clear message stating that no nodes matched the filter.
- **FR-009**: Filtered no-match messaging must remain distinct from command errors such as namespace discovery failures.
- **FR-010**: Automated tests must cover long-flag and shorthand namespace filtering, case-insensitive matching, empty-filter fallback, and filtered no-match messaging.
- **FR-011**: `README.md` must be updated in the same change to document the new namespace filter behavior and examples.
- **FR-012**: Repository validation must include `make test` before the feature is considered complete.

## Non-Goals

- Changing the name or meaning of the existing root `--node-filter` flag.
- Adding namespace filtering by labels, regex, or metadata other than the namespace name.
- Changing how `viewnode ns list` discovers namespaces from the cluster.
- Changing the unfiltered namespace list output format beyond the requested flag addition and filtered empty-state wording.
- Adding new namespace management commands beyond the filter behavior requested in issue `#53`.

## Implementation Notes

- Keep the implementation inside the existing CLI structure under `cmd/ns`, `cmd/root.go`, and `srv/view.go`; do not introduce a parallel command structure.
- Define namespace filtering as a local concern of `cmd/ns/list.go` so `-f` remains scoped to `ns list`.
- Apply filtering client-side after the existing single namespace list call rather than adding extra Kubernetes API requests.
- Reuse the existing namespace entry preparation path so alphabetical ordering and active-marker rendering stay aligned with current behavior.
- Root-command-dependent config lookups should continue using `config.Initialize` with `cmd.Root()` fallback to preserve `--kubeconfig` inheritance behavior.
- Filter-aware empty-state messaging for the root flow should be implemented with minimal additional context passed into the existing printout path in `srv/view.go`.
- Tests should follow repository guidance for Cobra inherited flags by executing through a small root command tree where appropriate.
- `README.md` is maintained manually in this repository and must be updated alongside the command change.
- Validation should include targeted command/output tests plus `make test`.

## Validation

- Run targeted Go tests covering namespace filter behavior and filtered empty-state messaging.
- Verify `viewnode ns list --filter kube` and `viewnode ns list -f kube` produce equivalent filtered results.
- Verify an empty namespace filter behaves the same as no filter.
- Verify filtered namespace and node no-match outcomes return success and use clear filter-aware wording.
- Review `README.md` to confirm command reference and examples match the implemented CLI behavior.
- Run `make test`.

## Traceability

- GitHub Issue: #53
- GitHub URL: https://github.com/NTTDATA-DACH/viewnode/issues/53
- GitHub Title: Add namespace name filtering to `ns list` via `--filter` flag
- Feature Name: 53-add-ns-filter
- Feature Directory: /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/53-add-ns-filter
- Spec Path: /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/53-add-ns-filter/spec.md
- Plan Path: /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/53-add-ns-filter/plan.md
- Tasks Path: /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/53-add-ns-filter/tasks.md
- PRD Path: /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/tasks/prd-53-add-ns-filter.md
- Source Status: derived from Speckit artifacts

## Assumptions / Open Questions

- The clearer root no-match wording should target the existing node-filter path specifically, because that is the additional improvement explicitly called out in the issue.
- The PRD assumes the existing `viewnode ns list` cluster query path remains the only namespace data source for this feature.
