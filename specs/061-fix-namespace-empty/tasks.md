# Tasks: Fix empty namespace groups in namespace-filtered node output

**Input**: Design documents from `/specs/061-fix-namespace-empty/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/namespace-filter-node-view-contract.md, quickstart.md

**Tests**: Tests are required for this feature because it fixes a user-visible rendering bug and the plan and constitution require regression coverage plus `make test` validation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g. `US1`, `US2`, `US3`)
- Each task includes an exact file path

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Align the implementation with the approved renderer contract and validation targets before editing code.

- [ ] T001 Review implementation scope, corrected namespace-group contract, and validation targets in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/061-fix-namespace-empty/plan.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Establish the shared grouped-rendering expectations that all story work depends on.

**⚠️ CRITICAL**: No user story work should begin until this phase is complete.

- [ ] T002 Add helper-level regression coverage for omitting empty namespace groups in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go
- [ ] T003 [P] Review grouped namespace construction in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go against /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/061-fix-namespace-empty/data-model.md

**Checkpoint**: Shared grouped-rendering rules are explicit and ready for story-specific implementation.

---

## Phase 3: User Story 1 - Show only matching namespaces per node (Priority: P1) 🎯 MVP

**Goal**: Omit namespace headings that have no matching pods on the current node while preserving grouped rendering for namespaces that do match.

**Independent Test**: Run `go test ./srv` and verify multi-namespace scoped output renders only matching namespace headings for a node, keeps those headings alphabetical, and still nests pods beneath the correct namespace.

### Tests for User Story 1

- [ ] T004 [P] [US1] Add grouped scoped printout coverage for omitting empty selected namespace headings in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go
- [ ] T005 [P] [US1] Add grouped scoped printout coverage for rendering multiple matching namespace headings alphabetically in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go

### Implementation for User Story 1

- [ ] T006 [US1] Refine grouped namespace construction so only namespaces with matching pods are returned in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go
- [ ] T007 [US1] Update grouped pod printout to omit zero-pod namespace headings while preserving grouped pod nesting in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go

**Checkpoint**: Multi-namespace `viewnode --namespace ...` no longer shows empty namespace headings and remains independently testable.

---

## Phase 4: User Story 2 - Keep nodes visible without matching pods (Priority: P2)

**Goal**: Preserve node visibility even when a displayed node has no pods in the selected namespaces.

**Independent Test**: Run `go test ./srv` and verify a node with zero matching namespace groups still renders its node line while no namespace headings appear beneath it.

### Tests for User Story 2

- [ ] T008 [P] [US2] Add grouped scoped printout coverage for nodes that remain visible without rendered namespace headings in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go
- [ ] T009 [P] [US2] Add grouped scoped printout coverage that summary counts and unscheduled pods remain unchanged when a node has no matching namespace groups in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go

### Implementation for User Story 2

- [ ] T010 [US2] Preserve node-line rendering when grouped namespace output is empty for a node in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go
- [ ] T011 [US2] Verify grouped render activation and namespace selection plumbing remain compatible with the corrected renderer behavior in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/cmd/root.go

**Checkpoint**: Nodes still appear even when no selected namespaces match, and adjacent output behavior remains stable.

---

## Phase 5: User Story 3 - Preserve filtered view consistency (Priority: P3)

**Goal**: Keep namespace-filtered grouped output consistent with the all-namespaces grouped path and update the documentation to match the corrected behavior.

**Independent Test**: Run `go test ./srv ./cmd`, review `README.md`, and confirm grouped namespace-filtered output follows the same per-node inclusion rule as `viewnode -A` while single-namespace scoped output remains unchanged.

### Tests for User Story 3

- [ ] T012 [P] [US3] Add regression coverage for grouped scoped parity with the all-namespaces inclusion rule in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go
- [ ] T013 [P] [US3] Add regression coverage that single-namespace scoped output remains on the flat path in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go

### Implementation for User Story 3

- [ ] T014 [US3] Preserve grouped all-namespaces rendering behavior while applying the same non-empty-group rule to scoped output in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view.go
- [ ] T015 [US3] Update namespace-filtered grouped output examples and explanatory text in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/README.md

**Checkpoint**: Scoped grouped output matches the intended all-namespaces parity rule and the README reflects shipped behavior.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final verification across the completed stories.

- [ ] T016 Run focused regression validation with `go test ./srv ./cmd` from /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode
- [ ] T017 Run full repository validation with `make test` from /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/Makefile
- [ ] T018 Review final behavior against /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/061-fix-namespace-empty/contracts/namespace-filter-node-view-contract.md and /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/061-fix-namespace-empty/quickstart.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1: Setup**: No dependencies
- **Phase 2: Foundational**: Depends on Phase 1 and blocks all user story work
- **Phase 3: User Story 1**: Depends on Phase 2
- **Phase 4: User Story 2**: Depends on User Story 1 because node visibility without namespace groups is meaningful only after empty-group omission exists
- **Phase 5: User Story 3**: Depends on User Story 1 and User Story 2 so parity and documentation reflect the final renderer behavior
- **Phase 6: Polish**: Depends on all desired user stories being complete

### User Story Dependencies

- **US1**: Can begin once grouped namespace helper expectations are in place
- **US2**: Depends on US1 because the renderer must already omit empty namespace groups before node-only rendering can be validated
- **US3**: Depends on US1 and US2 because parity and documentation should describe the final corrected output

### Within Each User Story

- Regression tests should be written or updated before the implementation they verify
- Grouping helper behavior precedes final grouped printout changes
- Renderer changes precede README updates
- Full-suite validation happens only after code and documentation changes are complete

### Parallel Opportunities

- `T003` can run in parallel with `T002` after the task scope is reviewed
- `T004` and `T005` can run in parallel once foundational expectations are set
- `T008` and `T009` can run in parallel after US1 behavior is established
- `T012` and `T013` can run in parallel once US2 behavior is stable

---

## Parallel Example: User Story 1

```bash
# Launch grouped scoped renderer regressions together:
Task: "Add grouped scoped printout coverage for omitting empty selected namespace headings in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
Task: "Add grouped scoped printout coverage for rendering multiple matching namespace headings alphabetically in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
```

---

## Parallel Example: User Story 2

```bash
# Launch node-visibility regressions together:
Task: "Add grouped scoped printout coverage for nodes that remain visible without rendered namespace headings in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
Task: "Add grouped scoped printout coverage that summary counts and unscheduled pods remain unchanged when a node has no matching namespace groups in /Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/srv/view_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational grouped-rendering expectations
3. Complete Phase 3: User Story 1
4. Validate the renderer independently with `go test ./srv`

### Incremental Delivery

1. Deliver empty-group omission for grouped namespace-filtered output
2. Preserve node visibility and adjacent summary behavior
3. Confirm parity with the grouped all-namespaces path and update the README
4. Run focused and full validation

### Parallel Team Strategy

1. One developer updates renderer behavior in `srv/view.go`
2. Another developer prepares regression coverage in `srv/view_test.go`
3. Documentation and full validation follow once renderer behavior stabilizes

---

## Notes

- All tasks follow the required checklist format with explicit task IDs and exact file paths
- User Story 1 is the suggested MVP scope
- The work remains centered in `srv/view.go`, with minimal compatibility review in `cmd/root.go`, focused regression coverage in `srv/view_test.go`, and README updates in the same change
