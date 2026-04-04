# Tasks: Remove redundant namespace prefixes from single-namespace pod rows

**Input**: Design documents from `/specs/066-fix-single-namespace-prefix/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/single-namespace-output-contract.md, quickstart.md

**Tests**: Tests are required for this feature because it fixes a user-visible CLI rendering bug and the plan plus constitution require regression coverage and `make test` validation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g. `US1`, `US2`, `US3`)
- Each task includes an exact file path

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Align implementation with the approved single-namespace output contract before changing code.

- [x] T001 Review implementation scope, single-namespace output contract, and validation targets in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/066-fix-single-namespace-prefix/plan.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Establish the shared understanding of the flat single-namespace render path that all story work depends on.

**⚠️ CRITICAL**: No user story work should begin until this phase is complete.

- [x] T002 Confirm the single-namespace flat render path and formatter seam in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go
- [x] T003 [P] Review the existing single-namespace regression coverage and affected expectations in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go

**Checkpoint**: The exact implementation seam and regression surface are identified before behavior changes begin.

---

## Phase 3: User Story 1 - Remove redundant prefixes in single-namespace output (Priority: P1) 🎯 MVP

**Goal**: Make single-namespace pod rows display only the pod name and status while keeping the namespace header visible.

**Independent Test**: Run `go test ./srv` and verify `viewnode -n <namespace>` style output keeps the namespace header but no longer prints `<namespace>: ` before each pod name.

### Tests for User Story 1

- [x] T004 [P] [US1] Update single-namespace flat-output regression expectations in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go
- [x] T005 [P] [US1] Add regression coverage that the namespace header remains visible while pod rows omit the inline namespace prefix in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go

### Implementation for User Story 1

- [x] T006 [US1] Refine the flat single-namespace pod-row formatter to omit the redundant inline namespace prefix in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go
- [x] T007 [US1] Preserve existing pod-row phase rendering and flat single-namespace structure while applying the formatter change in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go

**Checkpoint**: Single-namespace output no longer repeats the namespace on every pod row and remains independently testable.

---

## Phase 4: User Story 2 - Preserve existing single-namespace layout semantics (Priority: P2)

**Goal**: Keep counts, node grouping, pod ordering, and auxiliary pod-row details unchanged while removing the redundant prefix.

**Independent Test**: Run `go test ./srv` and verify the only single-namespace output difference is the missing inline namespace prefix, while counts, grouping, ordering, and pod-row details remain intact.

### Tests for User Story 2

- [x] T008 [P] [US2] Add regression coverage that single-namespace counts, node sections, and pod ordering remain unchanged in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go
- [x] T009 [P] [US2] Add regression coverage that container, timing, or metrics-capable pod rows remain compatible after prefix removal in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go

### Implementation for User Story 2

- [x] T010 [US2] Preserve auxiliary pod-row detail rendering in the single-namespace flat path after removing the prefix in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go
- [x] T011 [US2] Verify no command configuration changes are needed for the corrected single-namespace flat path in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root.go

**Checkpoint**: The single-namespace cleanup is isolated to the redundant prefix and does not alter adjacent output behavior.

---

## Phase 5: User Story 3 - Keep single-namespace output consistent with other views (Priority: P3)

**Goal**: Keep the corrected single-namespace output aligned with the repository’s other namespace display modes and document the change accurately.

**Independent Test**: Run `go test ./srv ./cmd`, review `README.md`, and confirm single-namespace output omits inline namespace prefixes while grouped multi-namespace and all-namespaces behavior remains unchanged.

### Tests for User Story 3

- [ ] T012 [P] [US3] Add regression coverage that grouped multi-namespace behavior remains unchanged while single-namespace output uses the corrected flat row format in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go
- [ ] T013 [P] [US3] Review and align single-namespace documentation examples with the corrected output in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/README.md

### Implementation for User Story 3

- [ ] T014 [US3] Keep multi-namespace and all-namespaces rendering behavior unchanged while finalizing the single-namespace formatter fix in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go
- [ ] T015 [US3] Update single-namespace examples and explanatory text in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/README.md

**Checkpoint**: The fix is consistent with other namespace display modes and the README matches shipped behavior.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final verification across the completed stories.

- [ ] T016 Run focused regression validation with `go test ./srv ./cmd` from /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode
- [ ] T017 Run full repository validation with `make test` from /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/Makefile
- [ ] T018 Review final behavior against /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/066-fix-single-namespace-prefix/contracts/single-namespace-output-contract.md and /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/066-fix-single-namespace-prefix/quickstart.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1: Setup**: No dependencies
- **Phase 2: Foundational**: Depends on Phase 1 and blocks all user story work
- **Phase 3: User Story 1**: Depends on Phase 2
- **Phase 4: User Story 2**: Depends on User Story 1 because layout-preservation validation only matters after the redundant prefix is removed
- **Phase 5: User Story 3**: Depends on User Story 1 and User Story 2 so consistency and documentation reflect the final output
- **Phase 6: Polish**: Depends on all desired user stories being complete

### User Story Dependencies

- **US1**: Can begin once the flat-render seam and current regression expectations are confirmed
- **US2**: Depends on US1 because it validates that the single-namespace fix remains narrowly scoped
- **US3**: Depends on US1 and US2 because documentation and cross-mode consistency should describe the final corrected output

### Within Each User Story

- Regression tests should be updated before or alongside implementation
- Formatter changes precede documentation updates
- Full validation happens after code and documentation changes are complete

### Parallel Opportunities

- `T003` can run in parallel with `T002` after the implementation scope is reviewed
- `T004` and `T005` can run in parallel once the single-namespace seam is confirmed
- `T008` and `T009` can run in parallel after US1 behavior is established
- `T012` and `T013` can run in parallel once US2 behavior is stable

---

## Parallel Example: User Story 1

```bash
# Launch single-namespace regression updates together:
Task: "Update single-namespace flat-output regression expectations in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
Task: "Add regression coverage that the namespace header remains visible while pod rows omit the inline namespace prefix in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
```

---

## Parallel Example: User Story 2

```bash
# Launch layout-preservation regressions together:
Task: "Add regression coverage that single-namespace counts, node sections, and pod ordering remain unchanged in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
Task: "Add regression coverage that container, timing, or metrics-capable pod rows remain compatible after prefix removal in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational flat-path analysis
3. Complete Phase 3: User Story 1
4. Validate the corrected single-namespace renderer independently with `go test ./srv`

### Incremental Delivery

1. Deliver the single-namespace prefix removal
2. Preserve counts, grouping, ordering, and auxiliary pod-row details
3. Confirm cross-mode consistency and update `README.md`
4. Run focused and full validation

### Parallel Team Strategy

1. One developer updates the formatter in `srv/view.go`
2. Another developer updates regression coverage in `srv/view_test.go`
3. Documentation and full validation follow once the output stabilizes

---

## Notes

- All tasks follow the required checklist format with explicit task IDs and exact file paths
- User Story 1 is the suggested MVP scope
- The work stays centered in `srv/view.go`, with focused regression coverage in `srv/view_test.go`, minimal compatibility review in `cmd/root.go`, and README updates in the same change
