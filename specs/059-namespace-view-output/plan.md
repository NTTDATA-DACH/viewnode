# Implementation Plan: Namespace tree view for namespace-filtered node output

**Branch**: `059-namespace-view-output` | **Date**: 2026-04-04 | **Spec**: [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/059-namespace-view-output/spec.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/059-namespace-view-output/spec.md)
**Input**: Feature specification from `/specs/059-namespace-view-output/spec.md`

## Summary

Refactor the root `viewnode` printout used with multi-value `--namespace` selection so scheduled pods are rendered beneath alphabetical namespace headings within each node, matching the grouped `--all-namespaces` tree structure while also showing empty grouping rows for selected namespaces that have no pods on a displayed node. The design keeps rendering logic in `srv/view.go`, extends the root-command view context only enough to pass the selected namespace set into the renderer, adds focused regression coverage in `srv/view_test.go` and `cmd/root_test.go`, and updates `README.md` examples to match the new scoped output contract.

## Technical Context

**Language/Version**: Go 1.25.0 (toolchain go1.25.7)  
**Primary Dependencies**: `github.com/spf13/cobra`, `k8s.io/client-go`, `k8s.io/apimachinery`, `github.com/sirupsen/logrus`, `github.com/stretchr/testify`  
**Storage**: Kubeconfig files plus live Kubernetes API data for node and pod discovery  
**Testing**: `go test`, `make test`  
**Target Platform**: Cross-platform CLI environments supported by Go and Kubernetes client-go  
**Project Type**: CLI  
**Performance Goals**: Preserve current interactive output latency and avoid any additional Kubernetes API calls by deriving namespace groups from the already loaded pod list plus the selected namespace set already parsed by the root command  
**Constraints**: Keep the root command surface unchanged, preserve summary and unscheduled-pod output, preserve pod ordering semantics from `--all-namespaces`, render alphabetical namespace groups for scoped multi-namespace output, show empty grouping rows for selected namespaces with no displayed pods on a node, and update `README.md` for the user-visible output change  
**Scale/Scope**: One root command wiring path in `cmd/root.go`, one renderer path in `srv/view.go`, focused regression tests in `srv/view_test.go` and `cmd/root_test.go`, and README command/example updates

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Operational Proof Over Intent**: Pass. The plan changes only observable CLI printout behavior for already loaded data and keeps success directly verifiable through output-focused tests and manual command runs.
- **II. CLI-First, Script-Safe Interfaces**: Pass. No command names or flags change; the output contract stays text-based and deterministic, with the scoped namespace tree documented explicitly.
- **III. Tests and Validation Are Mandatory**: Pass with required work. The plan adds targeted renderer and root-context tests and finishes with `make test`.
- **IV. Documentation Matches User Behavior**: Pass with required work. `README.md` will be updated in the same change to reflect the grouped multi-namespace scoped output.
- **V. Small, Compatible, Repository-Native Changes**: Pass. The design reuses the existing grouped renderer shape from issue `#57`, keeps branching in `srv/view.go`, and avoids introducing a parallel rendering subsystem.

## Project Structure

### Documentation (this feature)

```text
specs/059-namespace-view-output/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── namespace-filter-node-view-contract.md
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

**Structure Decision**: Keep namespace-grouped rendering in `srv/view.go`, keep `cmd/root.go` responsible only for parsing and passing the selected namespace context needed by the renderer, add or update output tests in `srv/view_test.go` plus root-context tests in `cmd/root_test.go`, and document the changed scoped output in `README.md`.

## Phase 0: Research

Research findings are recorded in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/059-namespace-view-output/research.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/059-namespace-view-output/research.md). Key decisions:

- Keep namespace grouping in `srv/view.go` because the feature changes presentation, not data loading.
- Extend the printout input with the explicit selected-namespace set so scoped rendering can show empty namespace headings without reparsing display strings or making new API calls.
- Reuse the `--all-namespaces` pod ordering behavior inside each namespace group, while sorting namespace headings alphabetically.
- Apply grouped rendering for scoped output only when more than one namespace is selected, leaving single-namespace and non-namespace runs on the current flat path.
- Preserve existing summary lines, unscheduled pod output, and inline/tree container rendering beneath grouped pod rows.

## Phase 1: Design & Contracts

### Data Model

Data entities and rendering relationships are defined in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/059-namespace-view-output/data-model.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/059-namespace-view-output/data-model.md).

### Contracts

The user-visible CLI output contract is documented in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/059-namespace-view-output/contracts/namespace-filter-node-view-contract.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/059-namespace-view-output/contracts/namespace-filter-node-view-contract.md).

### Quickstart

Implementation and validation steps are captured in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/059-namespace-view-output/quickstart.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/059-namespace-view-output/quickstart.md).

### Design Decisions

- Keep namespace-grouped rendering in `srv/view.go` so the feature remains a view-layer refactor consistent with issue `#57`.
- Extend `srv.ViewNodeData` or `srv.ViewNodeDataConfig` with the selected namespace set needed for scoped grouping, instead of inferring empty headings from formatted strings or changing pod-loading behavior.
- Enable `GroupPodsByNamespace` for multi-namespace scoped runs in `cmd/root.go` while preserving the existing flat path for single-namespace and non-namespace runs.
- Build per-node namespace groups from the union of:
  - namespaces represented by the node's scheduled pods
  - namespaces explicitly selected by the user for the current run
- Sort namespace grouping rows alphabetically within each node.
- Preserve the current pod ordering behavior already used by grouped `--all-namespaces` output within each namespace grouping.
- Render empty namespace grouping rows for selected namespaces that have no displayed pods on a node.
- Preserve existing unscheduled pod output, node summaries, and both inline and tree container rendering beneath grouped pod rows.
- Add regression coverage for:
  - grouped multi-namespace scoped output
  - alphabetical namespace grouping for scoped output
  - empty namespace rows for selected namespaces missing on a node
  - pod-order parity with grouped `--all-namespaces`
  - unchanged flat output for single-namespace scoped runs
  - root-command wiring that turns on grouped rendering only for multi-namespace scoped selection and `--all-namespaces`
- Update `README.md` examples and command reference text to show the grouped tree output for `viewnode --namespace a,b`.

## Phase 2: Implementation Outline

1. Refactor `cmd/root.go` to pass the parsed selected-namespace slice into the printout context and enable grouped namespace rendering for multi-namespace scoped runs without changing CLI flags.
2. Refactor `srv/view.go` to build per-node namespace groups from both scheduled pods and the selected namespace set, render alphabetical namespace headings, and keep grouped pod/container output consistent with the `--all-namespaces` path.
3. Expand `srv/view_test.go` with output-capture tests for grouped scoped output, empty namespace headings, pod-order parity, and unchanged single-namespace flat behavior.
4. Expand `cmd/root_test.go` with root-context tests covering grouped-render activation for multi-namespace scoped selection and non-activation for single-namespace scoped selection.
5. Update `README.md` command reference and examples to reflect the grouped namespace-filtered node output structure.
6. Run `make test` and confirm the rendered contract matches the documented scoped examples.

## Post-Design Constitution Check

- **I. Operational Proof Over Intent**: Pass. The design remains directly verifiable through printout tests, root-context tests, and CLI examples using actual rendered output.
- **II. CLI-First, Script-Safe Interfaces**: Pass. No CLI surface changes are introduced, and the output contract stays deterministic and documented.
- **III. Tests and Validation Are Mandatory**: Pass. The design includes focused regression tests and requires `make test`.
- **IV. Documentation Matches User Behavior**: Pass. README updates remain part of the same work unit as the rendering change.
- **V. Small, Compatible, Repository-Native Changes**: Pass. The feature reuses the existing grouped renderer path and only adds the minimal context needed for scoped grouping semantics.

## Complexity Tracking

No constitution violations require justification.
