# Tasks: Refactor `ctx list` and `ns list` commands

**Input**: Design documents from `/specs/48-refactor-ctx-list-ns-list-commands/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/list-command-contract.md, quickstart.md

**Tests**: Tests are required for this feature because the specification explicitly requires automated coverage for both refactored list commands and `make test` validation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g. `US1`, `US2`, `US3`)
- Each task includes an exact file path

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Capture the feature task scaffold in the active spec directory and align implementation targets with the existing command layout.

- [ ] T001 Review implementation targets and verification steps in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/48-refactor-ctx-list-ns-list-commands/plan.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Establish the shared list-entry preparation unit that both list commands will depend on.

**⚠️ CRITICAL**: No user story work should begin until this phase is complete.

- [ ] T002 Create shared list-entry types and ordered active-state helper in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/config/list_entries.go
- [ ] T003 [P] Add focused helper tests for ordering and active-state rendering in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/config/list_entries_test.go

**Checkpoint**: Shared list-preparation logic is available for both commands.

---

## Phase 3: User Story 1 - Keep context listing behavior stable after refactor (Priority: P1) 🎯 MVP

**Goal**: Refactor `viewnode ctx list` to use explicit context entries and the shared helper without changing output.

**Independent Test**: Run `go test ./cmd/ctx ./cmd/config` and verify `viewnode ctx list` still prints alphabetically ordered contexts with exactly one active marker.

### Tests for User Story 1

- [ ] T004 [P] [US1] Update context list command coverage for ordering and active marker behavior in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/ctx_test.go

### Implementation for User Story 1

- [ ] T005 [US1] Refactor context list assembly to use shared context entries in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/list.go

**Checkpoint**: `ctx list` remains behaviorally unchanged and testable on its own.

---

## Phase 4: User Story 2 - Keep namespace listing behavior aligned with context listing (Priority: P1)

**Goal**: Refactor `viewnode ns list` to use explicit namespace entries and the shared helper while preserving marker and ordering behavior.

**Independent Test**: Run `go test ./cmd/ns ./cmd/config` and verify `viewnode ns list` still prints alphabetically ordered namespaces with the active namespace marked.

### Tests for User Story 2

- [ ] T006 [P] [US2] Expand namespace list success coverage for ordering and active marker behavior in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list_test.go
- [ ] T007 [P] [US2] Add namespace default fallback coverage when the active context namespace is empty in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list_test.go

### Implementation for User Story 2

- [ ] T008 [US2] Refactor namespace list assembly to use shared namespace entries and default fallback semantics in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list.go

**Checkpoint**: `ns list` stays aligned with `ctx list` and handles empty namespace state as `default`.

---

## Phase 5: User Story 3 - Maintain predictable failure handling during the refactor (Priority: P2)

**Goal**: Preserve existing command error behavior while routing both commands through shared list-preparation code.

**Independent Test**: Run targeted failing command tests and confirm both commands still return descriptive errors without partial output.

### Tests for User Story 3

- [ ] T009 [P] [US3] Add or update context list failure coverage for kubeconfig errors in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/ctx_test.go
- [ ] T010 [P] [US3] Add or update namespace list failure coverage for namespace lookup errors in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list_test.go

### Implementation for User Story 3

- [ ] T011 [US3] Preserve existing error wrapping and no-partial-output behavior during refactor updates in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ctx/list.go
- [ ] T012 [US3] Preserve existing error wrapping and no-partial-output behavior during refactor updates in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list.go

**Checkpoint**: Failure paths remain stable after the shared logic extraction.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final verification and cleanup across all stories.

- [ ] T013 Run full repository validation with make test from /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/Makefile
- [ ] T014 Review final command behavior against the contract in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/48-refactor-ctx-list-ns-list-commands/contracts/list-command-contract.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1: Setup**: No dependencies
- **Phase 2: Foundational**: Depends on Phase 1 and blocks all user story work
- **Phase 3: User Story 1**: Depends on Phase 2
- **Phase 4: User Story 2**: Depends on Phase 2
- **Phase 5: User Story 3**: Depends on User Story 1 and User Story 2 because failure behavior must be verified against the refactored commands
- **Phase 6: Polish**: Depends on all user stories being complete

### User Story Dependencies

- **US1**: Can begin as soon as the shared helper exists
- **US2**: Can begin as soon as the shared helper exists
- **US3**: Depends on the refactored `ctx list` and `ns list` implementations being in place

### Within Each User Story

- Tests should be updated before or alongside implementation and must fail when the intended regression is reintroduced
- Shared helper changes precede command refactors
- Command refactors precede final full-suite validation

### Parallel Opportunities

- `T003` can run in parallel with planning of `T002` once helper shape is agreed
- `T004`, `T006`, and `T007` can proceed in parallel after Phase 2
- `T009` and `T010` can proceed in parallel after the corresponding command refactors are in place

---

## Parallel Example: User Story 2

```bash
# Launch namespace list test updates together:
Task: "Expand namespace list success coverage in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list_test.go"
Task: "Add namespace default fallback coverage in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/ns/list_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational helper extraction
3. Complete Phase 3: User Story 1
4. Validate `ctx list` behavior independently

### Incremental Delivery

1. Deliver shared helper plus `ctx list`
2. Deliver `ns list` on top of the same helper
3. Lock down failure handling and run full validation

### Parallel Team Strategy

1. One developer extracts and tests the shared helper
2. One developer refactors `ctx list` while another refactors `ns list`
3. Final pass covers failure-path hardening and `make test`

---

## Notes

- All tasks follow the required checklist format with explicit task IDs and file paths
- No README task is included because the plan and spec treat this issue as a behavior-preserving refactor with no user-facing CLI change
- Keep the helper local to existing command support patterns and avoid introducing a new parallel command structure
