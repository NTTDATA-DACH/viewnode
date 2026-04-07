# Tasks: Context list filtering

**Input**: Design documents from `/specs/069-add-ctx-filter/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/ctx-filter-contract.md, quickstart.md

**Tests**: Tests are required for this feature because the specification explicitly requires automated verification for filter behavior, no-match messaging, active-marker handling, and `make test` validation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this belongs to (e.g. `US1`, `US2`, `US3`)
- Each task includes an exact file path

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Align implementation work with the approved spec, plan, and CLI contract before code changes begin.

- [x] T001 Review implementation scope and verification targets in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/069-add-ctx-filter/plan.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Establish the shared context-filter normalization and matching support needed by all context-list filtering work.

**⚠️ CRITICAL**: No user story work should begin until this phase is complete.

- [x] T002 Create context filter normalization and case-insensitive matching helpers in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/list.go
- [x] T003 [P] Add focused helper coverage for empty and case-insensitive context filters in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/ctx_test.go

**Checkpoint**: Context filtering rules are defined and regression-tested before user-facing command changes.

---

## Phase 3: User Story 1 - Filter context results (Priority: P1) 🎯 MVP

**Goal**: Let users narrow `viewnode ctx list` results with a case-insensitive substring filter using either `--filter` or `-f`.

**Independent Test**: Run `go test ./cmd/ctx` and verify `viewnode ctx list --filter prod` and `viewnode ctx list -f prod` return only matching contexts, preserve alphabetical ordering, keep the active marker on displayed active rows, and treat an empty filter like no filter.

### Tests for User Story 1

- [x] T004 [P] [US1] Add command-tree coverage for `viewnode ctx list --filter` long-flag matching in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/ctx_test.go
- [x] T005 [P] [US1] Add command-tree coverage for `viewnode ctx list -f` shorthand matching and empty-filter fallback in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/ctx_test.go

### Implementation for User Story 1

- [x] T006 [US1] Add local `--filter` / `-f` flag wiring and case-insensitive context filtering to /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/list.go
- [x] T007 [US1] Preserve filtered context ordering and active-context rendering through the existing entry preparation path in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/list.go

**Checkpoint**: `viewnode ctx list` supports filter-driven context narrowing and remains independently testable.

---

## Phase 4: User Story 2 - Understand empty filter results (Priority: P2)

**Goal**: Show a clear, successful no-match message for filtered context views and preserve the clarified behavior when the current context is filtered out.

**Independent Test**: Run `go test ./cmd/ctx` and verify filtered context no-match output prints the exact success message, and filtered results omit the active marker when the current context does not match the filter.

### Tests for User Story 2

- [x] T008 [P] [US2] Add no-match success-path coverage for `viewnode ctx list --filter` in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/ctx_test.go
- [x] T009 [P] [US2] Add filtered-out-current-context coverage for active-marker omission in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/ctx_test.go

### Implementation for User Story 2

- [x] T010 [US2] Add the exact success-path no-match message `no contexts matched filter "<value>"` to /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/list.go
- [x] T011 [US2] Preserve matching-only filtered output when the current context is excluded in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/list.go

**Checkpoint**: Filtered no-match and current-context-filtered-out outcomes are clear, successful, and covered by regression tests.

---

## Phase 5: User Story 3 - Learn the new option from documentation (Priority: P3)

**Goal**: Document the new context filter flag and examples so users can discover and understand the feature from the README.

**Independent Test**: Review `README.md` and confirm the context command reference includes `--filter` / `-f`, plus examples for filtered listing and the expected no-match behavior.

### Tests for User Story 3

- [x] T012 [P] [US3] Verify README context command examples align with /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/069-add-ctx-filter/contracts/ctx-filter-contract.md in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/README.md

### Implementation for User Story 3

- [x] T013 [US3] Update context command reference and usage examples for `viewnode ctx list --filter` / `-f` in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/README.md

**Checkpoint**: Users can discover the new context filtering behavior from the maintained CLI documentation alone.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final verification across all implemented stories.

- [ ] T014 Run full repository validation with make test from /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/Makefile
- [ ] T015 Review final behavior against /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/069-add-ctx-filter/contracts/ctx-filter-contract.md and /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/069-add-ctx-filter/quickstart.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1: Setup**: No dependencies
- **Phase 2: Foundational**: Depends on Phase 1 and blocks all user story work
- **Phase 3: User Story 1**: Depends on Phase 2
- **Phase 4: User Story 2**: Depends on User Story 1 because no-match and active-marker edge cases build on the new context filter flow
- **Phase 5: User Story 3**: Depends on User Story 1 and User Story 2 so documentation reflects final command behavior
- **Phase 6: Polish**: Depends on all desired user stories being complete

### User Story Dependencies

- **US1**: Can begin as soon as foundational filter rules are in place
- **US2**: Depends on US1 because filtered empty-state behavior is only meaningful once context filtering exists
- **US3**: Depends on US1 and US2 so the README matches the implemented flags and output wording

### Within Each User Story

- Tests should be updated before or alongside implementation and must fail when the intended regression is reintroduced
- Filter normalization logic precedes user-facing context flag wiring
- Exact no-match messaging and active-marker omission precede README examples that describe them
- Full-suite validation happens only after code and documentation changes are complete

### Parallel Opportunities

- `T003` can run in parallel with review of `T002` once the helper shape is agreed
- `T004` and `T005` can run in parallel after Phase 2
- `T008` and `T009` can run in parallel once US1 filtering behavior is in place
- `T012` can run in parallel with drafting `T013` after final wording is agreed

---

## Parallel Example: User Story 1

```bash
# Launch context filter test updates together:
Task: "Add command-tree coverage for viewnode ctx list --filter in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/ctx_test.go"
Task: "Add command-tree coverage for viewnode ctx list -f and empty-filter fallback in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/ctx_test.go"
```

---

## Parallel Example: User Story 2

```bash
# Launch filtered edge-case regression tests together:
Task: "Add no-match success-path coverage for viewnode ctx list in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/ctx_test.go"
Task: "Add filtered-out-current-context coverage for active-marker omission in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/ctx_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational filter rules
3. Complete Phase 3: User Story 1
4. Validate context filtering independently with `go test ./cmd/ctx`

### Incremental Delivery

1. Deliver context filter mechanics and tests
2. Deliver exact no-match messaging and current-context-filtered-out behavior
3. Deliver README updates and run full validation

### Parallel Team Strategy

1. One developer finishes foundational filter normalization and command wiring
2. A second developer prepares context regression tests once the command behavior is agreed
3. Documentation and final validation happen after the CLI output wording stabilizes

---

## Notes

- All tasks follow the required checklist format with explicit task IDs and exact file paths
- User Story 1 is the suggested MVP scope
- The tasks preserve existing command structure and keep changes localized to `cmd/ctx`, tests, and `README.md`
