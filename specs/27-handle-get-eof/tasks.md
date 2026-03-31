# Tasks: Clearer EOF error guidance

**Input**: Design documents from `/specs/27-handle-get-eof/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/eof-error-guidance-contract.md, quickstart.md

**Tests**: Tests are required for this feature because the specification explicitly requires automated verification for scoped EOF guidance, non-EOF compatibility, and final `make test` validation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this belongs to (e.g. `US1`, `US2`, `US3`)
- Each task includes an exact file path

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Align implementation work with the approved spec, plan, and CLI contract before code changes begin.

- [ ] T001 Review implementation scope, quickstart validation, and CLI contract in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/plan.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Establish the shared EOF classification support used by all scoped guidance behavior.

**⚠️ CRITICAL**: No user story work should begin until this phase is complete.

- [ ] T002 Add a reusable scoped EOF classification helper and decorated error type in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/api_errors.go
- [ ] T003 [P] Add focused classification coverage for scoped and out-of-scope EOF patterns in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/node_test.go

**Checkpoint**: Scoped EOF classification is defined and regression-tested before user-facing fatal-message changes begin.

---

## Phase 3: User Story 1 - Get actionable guidance for EOF failures (Priority: P1) 🎯 MVP

**Goal**: Show a clearer fatal message for scoped Kubernetes API `Get` EOF failures that hints proxy configuration may be involved while keeping the original low-level error details visible.

**Independent Test**: Run `go test ./srv ./cmd` and verify a scoped retrieval EOF produces a fatal path that preserves the original error detail and adds proxy guidance for the affected resource.

### Tests for User Story 1

- [ ] T004 [P] [US1] Add decoration-specific regression coverage for wrapped scoped EOF errors in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/node_test.go
- [ ] T005 [P] [US1] Add root fatal-output coverage for scoped EOF guidance messaging in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root_test.go

### Implementation for User Story 1

- [ ] T006 [US1] Apply scoped EOF decoration to Kubernetes retrieval failures in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/api.go and /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/api_errors.go
- [ ] T007 [US1] Handle the scoped EOF decorated error with proxy guidance in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root.go
- [ ] T008 [US1] Preserve original low-level EOF details in the user-facing fatal output from /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root.go

**Checkpoint**: Scoped EOF failures now produce an actionable fatal message and remain independently testable.

---

## Phase 4: User Story 2 - Preserve clarity for other failures (Priority: P2)

**Goal**: Keep non-EOF and out-of-scope EOF failures on their existing paths so proxy guidance is only shown for the intended retrieval failures.

**Independent Test**: Run `go test ./srv ./cmd` and verify unauthorized, forbidden, non-EOF, and out-of-scope EOF failures do not receive the new proxy hint.

### Tests for User Story 2

- [ ] T009 [P] [US2] Add non-EOF compatibility coverage for existing decorated error branches in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root_test.go
- [ ] T010 [P] [US2] Add out-of-scope EOF regression coverage in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/node_test.go

### Implementation for User Story 2

- [ ] T011 [US2] Refine scoped EOF matching so only targeted Kubernetes API `Get` retrieval failures are decorated in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/api_errors.go
- [ ] T012 [US2] Keep existing non-EOF and special-case fatal or warning branches intact in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root.go
- [ ] T013 [US2] Ensure repeated scoped EOF failures continue to emit the same deterministic guidance in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root.go

**Checkpoint**: Only the intended scoped EOF failures receive the proxy hint, and compatibility paths remain intact.

---

## Phase 5: User Story 3 - Learn the behavior from documentation and tests (Priority: P3)

**Goal**: Make the EOF guidance behavior durable for maintainers through explicit tests and alignment with the documented feature contract and quickstart.

**Independent Test**: Review the contract and quickstart documents and run the targeted test suite to confirm maintainers can verify the intended behavior without relying on the GitHub issue text alone.

### Tests for User Story 3

- [ ] T014 [P] [US3] Add contract-aligned regression assertions for final EOF guidance wording in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root_test.go

### Implementation for User Story 3

- [ ] T015 [US3] Align implementation-focused verification steps with the final behavior in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/quickstart.md
- [ ] T016 [US3] Review and update the scoped EOF behavior contract if implementation wording changes in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/contracts/eof-error-guidance-contract.md

**Checkpoint**: Maintainers can rely on the contract, quickstart, and tests to preserve the EOF guidance behavior.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final verification across all implemented stories.

- [ ] T017 Run full repository validation with make test from /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/Makefile
- [ ] T018 Review final implementation against /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/spec.md and /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/contracts/eof-error-guidance-contract.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1: Setup**: No dependencies
- **Phase 2: Foundational**: Depends on Phase 1 and blocks all user story work
- **Phase 3: User Story 1**: Depends on Phase 2
- **Phase 4: User Story 2**: Depends on User Story 1 because compatibility preservation builds on the new scoped EOF guidance path
- **Phase 5: User Story 3**: Depends on User Story 1 and User Story 2 so the contract and quickstart reflect final behavior
- **Phase 6: Polish**: Depends on all desired user stories being complete

### User Story Dependencies

- **US1**: Can begin as soon as scoped EOF classification is in place
- **US2**: Depends on US1 because compatibility checks are meaningful only after the new guidance path exists
- **US3**: Depends on US1 and US2 so maintainers document and verify the settled final behavior

### Within Each User Story

- Tests should be updated before or alongside implementation and must fail when the intended regression is reintroduced
- Error classification must precede root-command fatal-message translation
- Compatibility preservation must be checked before final quickstart and contract review
- Full-suite validation happens only after code and artifact updates are complete

### Parallel Opportunities

- `T003` can run in parallel with review of `T002` once the helper shape is agreed
- `T004` and `T005` can run in parallel after Phase 2
- `T009` and `T010` can run in parallel once US1 guidance behavior is in place
- `T015` and `T016` can run in parallel after implementation wording stabilizes

---

## Parallel Example: User Story 1

```bash
# Launch scoped EOF regression tests together:
Task: "Add decoration-specific regression coverage for wrapped scoped EOF errors in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/node_test.go"
Task: "Add root fatal-output coverage for scoped EOF guidance messaging in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root_test.go"
```

---

## Parallel Example: User Story 2

```bash
# Launch compatibility regression tests together:
Task: "Add non-EOF compatibility coverage in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root_test.go"
Task: "Add out-of-scope EOF regression coverage in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/node_test.go"
```

---

## Parallel Example: User Story 3

```bash
# Launch implementation-artifact alignment tasks together:
Task: "Align implementation-focused verification steps in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/quickstart.md"
Task: "Review and update the scoped EOF behavior contract in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/contracts/eof-error-guidance-contract.md"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational EOF classification support
3. Complete Phase 3: User Story 1
4. Validate scoped EOF guidance independently with `go test ./srv ./cmd`

### Incremental Delivery

1. Deliver scoped EOF classification and proxy-guidance fatal output
2. Deliver compatibility protections for non-EOF and out-of-scope EOF paths
3. Align the contract, quickstart, and final verification steps

### Parallel Team Strategy

1. One developer finishes the shared EOF classification and root fatal-message implementation
2. A second developer prepares scoped and compatibility regression tests once the interfaces are agreed
3. Contract and quickstart alignment happen after the final message wording stabilizes

---

## Notes

- All tasks follow the required checklist format with explicit task IDs and exact file paths
- User Story 1 is the suggested MVP scope
- The tasks keep the implementation localized to `srv/api*.go`, `cmd/root.go`, targeted tests, and the feature artifacts under `specs/27-handle-get-eof/`
