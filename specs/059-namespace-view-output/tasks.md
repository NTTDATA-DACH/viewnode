# Tasks: Namespace tree view for namespace-filtered node output

**Input**: Design documents from `/specs/059-namespace-view-output/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/namespace-filter-node-view-contract.md, quickstart.md

**Tests**: Tests are required for this feature because the specification explicitly requires automated verification for grouped rendering behavior, ordering rules, empty namespace rows, and `make test` validation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g. `US1`, `US2`, `US3`)
- Each task includes an exact file path

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Align implementation and validation work with the approved scoped-rendering contract before code changes begin.

- [x] T001 Review implementation scope, scoped grouping constraints, and verification targets in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/059-namespace-view-output/plan.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Establish the shared scoped namespace-group rendering support that all story work depends on.

**⚠️ CRITICAL**: No user story work should begin until this phase is complete.

- [x] T002 Extend the view context to carry selected namespace values for grouped scoped rendering in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go
- [x] T003 [P] Add focused helper-level coverage for alphabetical namespace grouping, empty selected namespaces, and preserved pod order in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go

**Checkpoint**: Scoped namespace grouping rules are defined and regression-tested before user-facing output changes are completed.

---

## Phase 3: User Story 1 - Group filtered pods under namespaces (Priority: P1) 🎯 MVP

**Goal**: Render scheduled pods beneath alphabetical namespace headings in multi-namespace `viewnode --namespace ...` output so namespace ownership is readable without repeating inline namespace prefixes.

**Independent Test**: Run `go test ./srv` and verify a node with pods from multiple selected namespaces prints alphabetical namespace headings, nests pods beneath the correct heading, preserves grouped pod order parity with `--all-namespaces`, and shows empty grouping rows for selected namespaces with no displayed pods.

### Tests for User Story 1

- [x] T004 [P] [US1] Add grouped scoped printout coverage for multi-namespace nodes in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go
- [x] T005 [P] [US1] Add grouped scoped printout coverage for selected namespaces with no pods on a displayed node in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go
- [x] T006 [P] [US1] Add grouped scoped printout coverage for pod-order parity with `--all-namespaces` in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go

### Implementation for User Story 1

- [x] T007 [US1] Refactor selected-namespace parsing and grouped-render activation for multi-namespace scoped runs in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root.go
- [x] T008 [US1] Refactor scheduled pod rendering to print alphabetical scoped namespace grouping rows, including empty selected namespaces, in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go

**Checkpoint**: Multi-namespace `viewnode --namespace ...` shows grouped namespace tree output and remains independently testable.

---

## Phase 4: User Story 2 - Keep namespace-filtered output aligned with all-namespaces behavior (Priority: P2)

**Goal**: Keep scoped grouped output aligned with the established `--all-namespaces` structure while preserving existing flat output for single-namespace and non-grouped runs.

**Independent Test**: Run `go test ./srv ./cmd` and verify grouped scoped output matches the `--all-namespaces` structure for shared node/namespace cases, while single-namespace scoped output remains on the flat `<namespace>: <pod>` path.

### Tests for User Story 2

- [x] T009 [P] [US2] Add regression coverage that single-namespace scoped output remains on the flat pod-row path in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go
- [x] T010 [P] [US2] Add regression coverage for grouped scoped container rendering compatibility in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go
- [x] T011 [P] [US2] Add root-context coverage for grouped rendering activation only on multi-namespace scoped selection in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root_test.go

### Implementation for User Story 2

- [x] T012 [US2] Preserve `--all-namespaces` grouping parity and current flat scoped behavior while integrating scoped grouped rendering in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go
- [x] T013 [US2] Keep root-command namespace display and grouped-render gating minimal and compatible with existing flag behavior in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root.go

**Checkpoint**: Scoped grouped rendering matches the existing all-namespaces tree contract without regressing single-namespace or container-output behavior.

---

## Phase 5: User Story 3 - Preserve summary clarity while changing the tree layout (Priority: P3)

**Goal**: Keep summary lines, unscheduled pod output, node counts, and user-facing documentation aligned with the new grouped scoped output.

**Independent Test**: Run `go test ./srv ./cmd`, review `README.md`, and confirm grouped scoped output keeps existing totals and unscheduled pod behavior while the documentation reflects the new output shape.

### Tests for User Story 3

- [x] T014 [P] [US3] Add regression coverage for grouped scoped output with node summaries and unscheduled pod sections in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go
- [x] T015 [P] [US3] Verify README command reference and grouped scoped example match the output contract in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/README.md

### Implementation for User Story 3

- [x] T016 [US3] Preserve summary-line, unscheduled-pod, and node-count behavior while integrating grouped scoped rendering in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go
- [x] T017 [US3] Update scoped namespace command reference and grouped output examples in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/README.md

**Checkpoint**: Grouped scoped rendering preserves adjacent output behavior and the documentation reflects the shipped contract.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final verification across all implemented stories.

- [x] T018 Run full repository validation with make test from /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/Makefile
- [x] T019 Review final behavior against /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/059-namespace-view-output/contracts/namespace-filter-node-view-contract.md and /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/059-namespace-view-output/quickstart.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1: Setup**: No dependencies
- **Phase 2: Foundational**: Depends on Phase 1 and blocks all user story work
- **Phase 3: User Story 1**: Depends on Phase 2
- **Phase 4: User Story 2**: Depends on User Story 1 because grouped scoped rendering must exist before parity and gating regressions are meaningful
- **Phase 5: User Story 3**: Depends on User Story 1 and User Story 2 so summaries and documentation reflect the final rendered behavior
- **Phase 6: Polish**: Depends on all desired user stories being complete

### User Story Dependencies

- **US1**: Can begin once the selected-namespace view context and grouping helper logic exist
- **US2**: Depends on US1 because parity with `--all-namespaces` and single-namespace flat behavior are validated against the new grouped scoped rendering path
- **US3**: Depends on US1 and US2 so summary validation and README examples describe the final output shape

### Within Each User Story

- Tests should be updated before or alongside implementation and must fail when the intended regression is reintroduced
- Shared selected-namespace context precedes user-facing scoped printout changes
- Core grouped scoped rendering precedes parity and compatibility refinements
- README updates follow the final rendered output behavior
- Full-suite validation happens only after code and documentation changes are complete

### Parallel Opportunities

- `T003` can run in parallel with review of `T002` once the selected-namespace context shape is agreed
- `T004`, `T005`, and `T006` can run in parallel after Phase 2
- `T009`, `T010`, and `T011` can run in parallel once US1 grouped scoped rendering exists
- `T014` and `T015` can run in parallel after grouped scoped behavior stabilizes

---

## Parallel Example: User Story 1

```bash
# Launch grouped scoped rendering regression tests together:
Task: "Add grouped scoped printout coverage for multi-namespace nodes in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
Task: "Add grouped scoped printout coverage for selected namespaces with no pods on a displayed node in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
Task: "Add grouped scoped printout coverage for pod-order parity with --all-namespaces in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
```

---

## Parallel Example: User Story 2

```bash
# Launch scoped-parity and gating regression tests together:
Task: "Add regression coverage that single-namespace scoped output remains on the flat pod-row path in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
Task: "Add regression coverage for grouped scoped container rendering compatibility in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
Task: "Add root-context coverage for grouped rendering activation only on multi-namespace scoped selection in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational selected-namespace grouping support
3. Complete Phase 3: User Story 1
4. Validate grouped multi-namespace scoped rendering independently with `go test ./srv`

### Incremental Delivery

1. Deliver grouped namespace rendering for multi-namespace scoped node output
2. Preserve and validate parity with `--all-namespaces`, single-namespace flat output, and container modes
3. Preserve summary behavior, update README examples, and run full validation

### Parallel Team Strategy

1. One developer finalizes selected-namespace context plumbing and core rendering in `cmd/root.go` and `srv/view.go`
2. A second developer prepares grouped scoped regression tests in `srv/view_test.go` and `cmd/root_test.go`
3. Documentation and full-suite verification happen after rendering behavior stabilizes

---

## Notes

- All tasks follow the required checklist format with explicit task IDs and exact file paths
- User Story 1 is the suggested MVP scope
- The tasks keep the implementation centered in `srv/view.go`, with minimal `cmd/root.go` support, focused `srv/view_test.go` and `cmd/root_test.go` coverage, and README updates
