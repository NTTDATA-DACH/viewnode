# Tasks: Namespace list filtering

**Input**: Design documents from `/specs/53-add-ns-filter/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/filter-command-contract.md, quickstart.md

**Tests**: Tests are required for this feature because the specification explicitly requires automated verification for filtering behavior, no-match messaging, and `make test` validation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g. `US1`, `US2`, `US3`)
- Each task includes an exact file path

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Align implementation work with the approved spec, plan, and contract before code changes begin.

- [ ] T001 Review implementation scope and verification targets in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/53-add-ns-filter/plan.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Establish the shared filter normalization and matching support needed by namespace filtering work.

**⚠️ CRITICAL**: No user story work should begin until this phase is complete.

- [ ] T002 Create namespace filter normalization and case-insensitive matching helpers in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list.go
- [ ] T003 [P] Add focused helper coverage for empty and case-insensitive namespace filters in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list_test.go

**Checkpoint**: Namespace filtering rules are defined and regression-tested before user-facing command changes.

---

## Phase 3: User Story 1 - Filter namespace results (Priority: P1) 🎯 MVP

**Goal**: Let users narrow `viewnode ns list` results with a case-insensitive substring filter using either `--filter` or `-f`.

**Independent Test**: Run `go test ./cmd/ns` and verify `viewnode ns list --filter kube` and `viewnode ns list -f kube` return only matching namespaces, preserving alphabetical ordering and active markers; confirm an empty filter behaves like no filter.

### Tests for User Story 1

- [ ] T004 [P] [US1] Add command-tree coverage for `viewnode ns list --filter` long-flag matching in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list_test.go
- [ ] T005 [P] [US1] Add command-tree coverage for `viewnode ns list -f` shorthand matching and empty-filter fallback in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list_test.go

### Implementation for User Story 1

- [ ] T006 [US1] Add local `--filter` / `-f` flag wiring and case-insensitive namespace filtering to /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list.go
- [ ] T007 [US1] Preserve filtered namespace ordering and active namespace rendering through existing entry preparation in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list.go

**Checkpoint**: `viewnode ns list` supports filter-driven namespace narrowing and remains independently testable.

---

## Phase 4: User Story 2 - Understand empty filter results (Priority: P2)

**Goal**: Show clear, successful no-match messages for filtered namespace and node views so users can distinguish zero matches from command failures.

**Independent Test**: Run `go test ./cmd/ns ./srv` and verify filtered namespace and node flows print explicit no-match messages while still succeeding.

### Tests for User Story 2

- [ ] T008 [P] [US2] Add namespace no-match success-path coverage for filtered results in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list_test.go
- [ ] T009 [P] [US2] Add root filtered empty-state output coverage for node-filter no-match messaging in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go

### Implementation for User Story 2

- [ ] T010 [US2] Add clear success-path namespace no-match messaging to /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list.go
- [ ] T011 [US2] Add filter-aware empty-result context plumbing in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root.go
- [ ] T012 [US2] Replace generic filtered no-node empty-state wording in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go

**Checkpoint**: Filtered no-match outcomes for namespaces and nodes are clear, successful, and covered by regression tests.

---

## Phase 5: User Story 3 - Learn the new behavior from documentation (Priority: P3)

**Goal**: Document the new namespace filter flag and examples so users can discover and understand the feature from the README.

**Independent Test**: Review `README.md` and confirm the namespace command reference includes `--filter` / `-f`, plus examples for filtered listing and expected no-match behavior.

### Tests for User Story 3

- [ ] T013 [P] [US3] Verify README command reference and examples align with the filter contract in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/README.md

### Implementation for User Story 3

- [ ] T014 [US3] Update namespace command reference and usage examples for `viewnode ns list --filter` / `-f` in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/README.md

**Checkpoint**: Users can discover the new namespace filtering behavior from the maintained CLI documentation alone.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final verification across all implemented stories.

- [ ] T015 Run full repository validation with make test from /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/Makefile
- [ ] T016 Review final behavior against /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/53-add-ns-filter/contracts/filter-command-contract.md and /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/53-add-ns-filter/quickstart.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1: Setup**: No dependencies
- **Phase 2: Foundational**: Depends on Phase 1 and blocks all user story work
- **Phase 3: User Story 1**: Depends on Phase 2
- **Phase 4: User Story 2**: Depends on User Story 1 because no-match messaging builds on the new namespace filter flow and shared filtered output expectations
- **Phase 5: User Story 3**: Depends on User Story 1 and User Story 2 so documentation reflects final command behavior
- **Phase 6: Polish**: Depends on all desired user stories being complete

### User Story Dependencies

- **US1**: Can begin as soon as foundational filter rules are in place
- **US2**: Depends on US1 because filtered empty-state messaging is only meaningful once namespace filtering exists
- **US3**: Depends on US1 and US2 so the README matches the implemented flags and output wording

### Within Each User Story

- Tests should be updated before or alongside implementation and must fail when the intended regression is reintroduced
- Filter normalization logic precedes user-facing namespace flag wiring
- Namespace no-match messaging precedes README examples that describe it
- Full-suite validation happens only after code and documentation changes are complete

### Parallel Opportunities

- `T003` can run in parallel with review of `T002` once the helper shape is agreed
- `T004` and `T005` can run in parallel after Phase 2
- `T008` and `T009` can run in parallel once US1 filtering behavior is in place
- `T013` can run in parallel with drafting `T014` after final wording is agreed

---

## Parallel Example: User Story 1

```bash
# Launch namespace filter test updates together:
Task: "Add command-tree coverage for viewnode ns list --filter in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list_test.go"
Task: "Add command-tree coverage for viewnode ns list -f and empty-filter fallback in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list_test.go"
```

---

## Parallel Example: User Story 2

```bash
# Launch filtered no-match regression tests together:
Task: "Add namespace no-match success-path coverage in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list_test.go"
Task: "Add root filtered empty-state output coverage in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational filter rules
3. Complete Phase 3: User Story 1
4. Validate namespace filtering independently with `go test ./cmd/ns`

### Incremental Delivery

1. Deliver namespace filter mechanics and tests
2. Deliver filter-aware no-match messaging across namespace and node flows
3. Deliver README updates and run full validation

### Parallel Team Strategy

1. One developer finishes foundational filter normalization and command wiring
2. A second developer prepares namespace and output regression tests once the interfaces are agreed
3. Documentation and final validation happen after message wording stabilizes

---

## Notes

- All tasks follow the required checklist format with explicit task IDs and exact file paths
- User Story 1 is the suggested MVP scope
- The tasks preserve existing command structure and keep changes localized to `cmd/ns`, `cmd/root.go`, `srv/view.go`, tests, and `README.md`
