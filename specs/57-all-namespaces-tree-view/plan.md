# Implementation Plan: Namespace tree view for all-namespaces node output

**Branch**: `57-all-namespaces-tree-view` | **Date**: 2026-04-01 | **Spec**: [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/57-all-namespaces-tree-view/spec.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/57-all-namespaces-tree-view/spec.md)
**Input**: Feature specification from `/specs/57-all-namespaces-tree-view/spec.md`

## Summary

Refactor the root `viewnode` printout used with `--all-namespaces` so scheduled pods are grouped beneath alphabetical namespace headings within each node, while preserving existing node summaries, unscheduled pod handling, non-all-namespaces output, and container rendering modes. The design keeps the change centered in `srv/view.go`, passes only the existing namespace-display context from `cmd/root.go`, adds focused printout tests in `srv/view_test.go`, and updates `README.md` examples to match the new user-visible tree structure.

## Technical Context

**Language/Version**: Go 1.25.0 (toolchain go1.25.7)  
**Primary Dependencies**: `github.com/spf13/cobra`, `k8s.io/client-go`, `k8s.io/apimachinery`, `github.com/sirupsen/logrus`, `github.com/stretchr/testify`  
**Storage**: Kubeconfig files plus live Kubernetes API data for node and pod discovery  
**Testing**: `go test`, `make test`  
**Target Platform**: Cross-platform CLI environments supported by Go and Kubernetes client-go  
**Project Type**: CLI  
**Performance Goals**: Preserve current interactive output latency and avoid any additional Kubernetes API calls by deriving namespace groups from the already loaded pod list during printout  
**Constraints**: Keep the root command surface unchanged, preserve existing non-`--all-namespaces` output, keep filter-aware empty-state messaging in `srv/view.go`, preserve container inline/tree rendering, keep namespace groups alphabetically ordered within each node, and update `README.md` for the user-visible output change  
**Scale/Scope**: One root printout path in `srv/view.go`, one root command wiring path in `cmd/root.go`, focused output regression tests in `srv/view_test.go`, and README command/example updates

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Operational Proof Over Intent**: Pass. The plan changes only the rendered output structure for already loaded data and keeps user-visible results directly observable through CLI output and regression tests.
- **II. CLI-First, Script-Safe Interfaces**: Pass. No command names or flags change; the output contract remains text-based and deterministic, with the new tree layout documented explicitly.
- **III. Tests and Validation Are Mandatory**: Pass with required work. The plan adds targeted printout tests for grouped namespace rendering and finishes with `make test`.
- **IV. Documentation Matches User Behavior**: Pass with required work. `README.md` will be updated in the same change to reflect the grouped all-namespaces output.
- **V. Small, Compatible, Repository-Native Changes**: Pass. The design keeps the refactor within existing `cmd/root.go` and `srv/view.go` flows and avoids new abstractions unless needed for a small rendering helper.

## Project Structure

### Documentation (this feature)

```text
specs/57-all-namespaces-tree-view/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── all-namespaces-node-view-contract.md
└── tasks.md
```

### Source Code (repository root)

```text
cmd/
├── root.go
└── root_test.go

srv/
├── view.go
└── view_test.go

README.md
Makefile
go.mod
```

**Structure Decision**: Keep all rendering logic in `srv/view.go`, keep the root command responsible only for passing existing namespace-display context into `srv.ViewNodeData`, add or update tests in `srv/view_test.go`, and document the changed output in `README.md`.

## Phase 0: Research

Research findings are recorded in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/57-all-namespaces-tree-view/research.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/57-all-namespaces-tree-view/research.md). Key decisions:

- Build namespace groups during printout from each node's existing `[]ViewPod` slice instead of changing load/filter logic or API calls.
- Sort namespace grouping rows alphabetically within each node while preserving each namespace's current pod order beneath its heading.
- Apply grouping only when `ShowNamespaces` is active so single-namespace and explicitly scoped views keep their existing output.
- Preserve current unscheduled pod output and existing container inline/tree rendering beneath grouped pod rows.

## Phase 1: Design & Contracts

### Data Model

Data entities and rendering relationships are defined in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/57-all-namespaces-tree-view/data-model.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/57-all-namespaces-tree-view/data-model.md).

### Contracts

The user-visible CLI output contract is documented in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/57-all-namespaces-tree-view/contracts/all-namespaces-node-view-contract.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/57-all-namespaces-tree-view/contracts/all-namespaces-node-view-contract.md).

### Quickstart

Implementation and validation steps are captured in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/57-all-namespaces-tree-view/quickstart.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/57-all-namespaces-tree-view/quickstart.md).

### Design Decisions

- Keep namespace-grouped rendering in `srv/view.go` because the change is presentational and the repository guidance explicitly keeps filter-aware and similar empty-state logic in that layer.
- Extend `srv.ViewNodeData` only if a minimal rendering flag or helper input is needed; avoid moving branching logic into `cmd/root.go`.
- Derive namespace grouping rows from each node's existing `[]ViewPod` collection so no pod-loading or namespace-loading API behavior changes.
- Sort namespace grouping rows alphabetically by namespace name within each node.
- Preserve the existing order of pods within each namespace group unless the current data pipeline already changes it elsewhere.
- Keep namespace grouping rows visible even when a node has pods from only one namespace, so `--all-namespaces` output remains structurally consistent.
- Preserve existing unscheduled pod output at the top of the printout and leave non-`--all-namespaces` pod lines unchanged.
- Reuse current container rendering branches so inline and tree container layouts continue to work beneath grouped pod rows.
- Add printout-level regression tests for:
  - alphabetical namespace grouping within a node
  - single-namespace grouping row behavior under `--all-namespaces`
  - unchanged output when namespace grouping is disabled
  - grouped output coexisting with node summaries and container rendering
  - unchanged no-node and filter-aware empty-state messaging
- Update `README.md` examples and command reference text to show the grouped namespace tree for `viewnode -A`.

## Phase 2: Implementation Outline

1. Refactor `srv/view.go` to compute per-node namespace groupings when `ShowNamespaces` is active and render grouped pod rows under namespace headings while preserving existing summary lines and container output.
2. Keep `cmd/root.go` wiring minimal so the root command still only sets the namespace-display intent that the view layer uses to choose grouped rendering.
3. Expand `srv/view_test.go` with output-capture tests covering alphabetical namespace groups, single-namespace grouping rows, unchanged scoped output, and compatibility with container rendering.
4. Update `README.md` command reference and examples to reflect the grouped all-namespaces node output structure.
5. Run `make test` and confirm the rendered contract matches the documented examples.

## Post-Design Constitution Check

- **I. Operational Proof Over Intent**: Pass. The design remains directly verifiable through printout tests and CLI examples using the actual rendered output.
- **II. CLI-First, Script-Safe Interfaces**: Pass. No CLI surface changes are introduced, and the output contract stays deterministic and documented.
- **III. Tests and Validation Are Mandatory**: Pass. The design includes focused output regression tests and requires `make test`.
- **IV. Documentation Matches User Behavior**: Pass. README updates remain part of the same work unit as the rendering change.
- **V. Small, Compatible, Repository-Native Changes**: Pass. The feature stays inside the existing view rendering path with no new command hierarchy or data-loading abstraction.

## Complexity Tracking

No constitution violations require justification.
