# Tasks: Namespace tree view for all-namespaces node output

**Input**: Design documents from `/specs/57-all-namespaces-tree-view/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/all-namespaces-node-view-contract.md, quickstart.md

**Tests**: Tests are required for this feature because the specification explicitly requires automated verification for grouped rendering behavior, output stability, and `make test` validation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g. `US1`, `US2`, `US3`)
- Each task includes an exact file path

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Align implementation and validation work with the approved rendering contract before code changes begin.

- [x] T001 Review implementation scope, rendering constraints, and verification targets in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/57-all-namespaces-tree-view/plan.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Establish the shared namespace-group rendering support that all story work depends on.

**⚠️ CRITICAL**: No user story work should begin until this phase is complete.

- [x] T002 Create helper logic for per-node alphabetical namespace grouping in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go
- [x] T003 [P] Add focused helper-level coverage for alphabetical grouping and preserved pod order in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go

**Checkpoint**: Namespace grouping rules are defined and regression-tested before user-facing output changes are completed.

---

## Phase 3: User Story 1 - Group node pods by namespace (Priority: P1) 🎯 MVP

**Goal**: Render scheduled pods beneath alphabetical namespace headings in `viewnode --all-namespaces` output so namespace ownership is readable without repeating inline namespace prefixes.

**Independent Test**: Run `go test ./srv` and verify a node with pods from multiple namespaces prints namespace headings in alphabetical order, keeps pod rows nested beneath the correct heading, and still shows a namespace heading when only one namespace exists on the node.

### Tests for User Story 1

- [x] T004 [P] [US1] Add grouped all-namespaces printout coverage for multi-namespace nodes in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go
- [x] T005 [P] [US1] Add grouped all-namespaces printout coverage for single-namespace nodes in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go

### Implementation for User Story 1

- [x] T006 [US1] Refactor scheduled pod rendering to print namespace grouping rows instead of repeated inline namespace prefixes in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go
- [x] T007 [US1] Preserve minimal namespace-display gating for grouped rendering in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root.go

**Checkpoint**: `viewnode --all-namespaces` shows grouped namespace tree output and remains independently testable.

---

## Phase 4: User Story 2 - Preserve readable totals and node summaries (Priority: P2)

**Goal**: Keep existing summary lines, unscheduled pod output, scoped output, and container rendering understandable while the grouped namespace tree is enabled.

**Independent Test**: Run `go test ./srv ./cmd` and verify grouped all-namespaces output keeps current totals and unscheduled pod behavior, scoped output remains flat, and container inline/tree modes still render under grouped pod rows.

### Tests for User Story 2

- [x] T008 [P] [US2] Add regression coverage for grouped output with node summaries and unscheduled pod sections in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go
- [x] T009 [P] [US2] Add regression coverage for grouped output with inline and tree container rendering in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go
- [x] T010 [P] [US2] Add regression coverage that scoped namespace output remains on the flat pod-row path in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go

### Implementation for User Story 2

- [x] T011 [US2] Preserve current summary, unscheduled pod, and scoped-output branches while integrating grouped rendering in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go
- [x] T012 [US2] Keep grouped pod rows compatible with existing inline and tree container formatting in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go

**Checkpoint**: Grouped rendering coexists with existing view behavior and does not regress adjacent output modes.

---

## Phase 5: User Story 3 - Learn the grouped output from docs and examples (Priority: P3)

**Goal**: Update the maintained CLI documentation so users can recognize the grouped all-namespaces output structure from the README alone.

**Independent Test**: Review `README.md` and confirm the root command examples and command reference reflect namespace grouping rows for all-namespaces node output.

### Tests for User Story 3

- [x] T013 [P] [US3] Verify README examples and command reference match the grouped output contract in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/README.md

### Implementation for User Story 3

- [x] T014 [US3] Update all-namespaces node output examples and related command reference text in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/README.md

**Checkpoint**: Users can discover the grouped all-namespaces output behavior from the maintained documentation.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final verification across all implemented stories.

- [x] T015 Run full repository validation with make test from /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/Makefile
- [x] T016 Review final behavior against /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/57-all-namespaces-tree-view/contracts/all-namespaces-node-view-contract.md and /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/57-all-namespaces-tree-view/quickstart.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1: Setup**: No dependencies
- **Phase 2: Foundational**: Depends on Phase 1 and blocks all user story work
- **Phase 3: User Story 1**: Depends on Phase 2
- **Phase 4: User Story 2**: Depends on User Story 1 because adjacent output preservation is validated against the new grouped rendering path
- **Phase 5: User Story 3**: Depends on User Story 1 and User Story 2 so documentation reflects the final rendered behavior
- **Phase 6: Polish**: Depends on all desired user stories being complete

### User Story Dependencies

- **US1**: Can begin once the shared namespace-grouping helper logic exists
- **US2**: Depends on US1 because summary and container compatibility are only meaningful after grouped rendering exists
- **US3**: Depends on US1 and US2 so README examples describe the final output shape

### Within Each User Story

- Tests should be updated before or alongside implementation and must fail when the intended regression is reintroduced
- Shared grouping helper logic precedes user-facing printout changes
- Core grouped rendering precedes compatibility refinements for summaries and containers
- README updates follow the final rendered output behavior
- Full-suite validation happens only after code and documentation changes are complete

### Parallel Opportunities

- `T003` can run in parallel with review of `T002` once the helper shape is agreed
- `T004` and `T005` can run in parallel after Phase 2
- `T008`, `T009`, and `T010` can run in parallel once US1 grouped rendering exists
- `T013` can run in parallel with drafting `T014` after output wording is settled

---

## Parallel Example: User Story 1

```bash
# Launch grouped rendering regression tests together:
Task: "Add grouped all-namespaces printout coverage for multi-namespace nodes in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
Task: "Add grouped all-namespaces printout coverage for single-namespace nodes in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
```

---

## Parallel Example: User Story 2

```bash
# Launch adjacent-output regression tests together:
Task: "Add regression coverage for grouped output with node summaries and unscheduled pod sections in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
Task: "Add regression coverage for grouped output with inline and tree container rendering in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
Task: "Add regression coverage that scoped namespace output remains on the flat pod-row path in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational namespace grouping rules
3. Complete Phase 3: User Story 1
4. Validate grouped all-namespaces rendering independently with `go test ./srv`

### Incremental Delivery

1. Deliver grouped namespace rendering for all-namespaces node output
2. Preserve and validate totals, unscheduled pod sections, scoped output, and container modes
3. Update README examples and run full validation

### Parallel Team Strategy

1. One developer finalizes shared grouping helper logic and core rendering in `srv/view.go`
2. A second developer prepares grouped-output regression tests in `srv/view_test.go`
3. Documentation and full-suite verification happen after rendering behavior stabilizes

---

## Notes

- All tasks follow the required checklist format with explicit task IDs and exact file paths
- User Story 1 is the suggested MVP scope
- The tasks keep the implementation centered in `srv/view.go`, with minimal `cmd/root.go` support, focused `srv/view_test.go` coverage, and README updates
