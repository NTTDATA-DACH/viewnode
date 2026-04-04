# Implementation Plan: Fix empty namespace groups in namespace-filtered node output

**Branch**: `061-fix-namespace-empty` | **Date**: 2026-04-04 | **Spec**: [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/061-fix-namespace-empty/spec.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/061-fix-namespace-empty/spec.md)
**Input**: Feature specification from `/specs/061-fix-namespace-empty/spec.md`

## Summary

Correct the grouped namespace-filtered `viewnode` printout so each node renders only namespace headings that actually contain matching pods on that node, while still keeping nodes visible when no selected namespaces match. The implementation stays in the existing `srv/view.go` renderer, preserves current root command flag behavior in `cmd/root.go`, adds focused regression coverage in `srv/view_test.go`, and updates `README.md` to remove the outdated empty-namespace example.

## Technical Context

**Language/Version**: Go 1.25.0 (toolchain go1.25.7)  
**Primary Dependencies**: `github.com/spf13/cobra`, `k8s.io/client-go`, `k8s.io/apimachinery`, `github.com/sirupsen/logrus`, `github.com/stretchr/testify`  
**Storage**: Kubeconfig files plus live Kubernetes API data for node and pod discovery  
**Testing**: `go test`, `make test`  
**Target Platform**: Cross-platform CLI environments supported by Go and Kubernetes client-go  
**Project Type**: CLI  
**Performance Goals**: Preserve current output latency and avoid any additional Kubernetes API calls by deriving namespace visibility from the already loaded per-node pod list during rendering  
**Constraints**: Keep the existing CLI surface unchanged, preserve summary and unscheduled-pod output, preserve grouped namespace ordering and pod ordering behavior for matching groups, keep nodes visible even when they have zero matching pods, and update `README.md` to match the corrected user-visible behavior  
**Scale/Scope**: One renderer path in `srv/view.go`, existing root command wiring in `cmd/root.go`, focused regression tests in `srv/view_test.go` and `cmd/root_test.go` as needed, and one README correction for namespace-filtered grouped output

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Operational Proof Over Intent**: Pass. The change is fully observable in CLI output and can be verified through renderer-focused tests and targeted manual command runs.
- **II. CLI-First, Script-Safe Interfaces**: Pass. No commands or flags change; only incorrect text output is corrected while preserving deterministic CLI behavior.
- **III. Tests and Validation Are Mandatory**: Pass with required work. The implementation must add or update regression coverage around grouped namespace rendering and finish with `make test`.
- **IV. Documentation Matches User Behavior**: Pass with required work. `README.md` must be updated in the same change because it currently documents the incorrect empty-group behavior.
- **V. Small, Compatible, Repository-Native Changes**: Pass. The design reuses the current renderer and view config instead of changing data loading or adding a new rendering pipeline.

## Project Structure

### Documentation (this feature)

```text
specs/061-fix-namespace-empty/
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

**Structure Decision**: Keep namespace visibility decisions in `srv/view.go`, keep `cmd/root.go` limited to passing already parsed namespace-selection context into `srv.ViewNodeDataConfig`, cover the behavior with output-level tests in `srv/view_test.go` and only update `cmd/root_test.go` if configuration expectations change, and document the corrected grouped scoped output in `README.md`.

## Phase 0: Research

Research findings are recorded in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/061-fix-namespace-empty/research.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/061-fix-namespace-empty/research.md). Key decisions:

- Keep namespace-group filtering in `srv/view.go` because the defect is a rendering issue, not a loading issue.
- Continue passing `SelectedNamespaces` from `cmd/root.go`, but use that selection only to constrain eligible namespace groups rather than forcing empty groups to render.
- Preserve alphabetical ordering of rendered namespace headings and existing pod ordering within each namespace group.
- Keep grouped rendering enabled for multi-namespace scoped runs and `--all-namespaces`, while single-namespace scoped runs remain on the flat output path.
- Preserve node visibility, summary lines, unscheduled-pod output, and container rendering while omitting zero-pod namespace headings for a node.

## Phase 1: Design & Contracts

### Data Model

Data entities and renderer relationships are defined in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/061-fix-namespace-empty/data-model.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/061-fix-namespace-empty/data-model.md).

### Contracts

The corrected user-visible CLI output contract is documented in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/061-fix-namespace-empty/contracts/namespace-filter-node-view-contract.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/061-fix-namespace-empty/contracts/namespace-filter-node-view-contract.md).

### Quickstart

Implementation and validation steps are captured in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/061-fix-namespace-empty/quickstart.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/061-fix-namespace-empty/quickstart.md).

### Design Decisions

- Keep `srv.ViewNodeDataConfig.SelectedNamespaces` as the canonical scoped-selection input from `cmd/root.go`.
- Refine grouped rendering so `groupPodsByNamespace` returns namespace groups only for namespaces that have at least one matching pod on the current node.
- Preserve alphabetical namespace heading order among rendered groups.
- Preserve pod ordering within each namespace group exactly as it appears in the already loaded pod slice.
- Keep nodes with zero matching pods visible and render no namespace headings beneath those nodes.
- Leave all-namespaces behavior consistent with the same rule: only namespaces with pods on the current node are shown.
- Preserve unscheduled pod rendering, node summaries, and inline/tree container output.
- Add regression coverage for:
  - grouped scoped output omitting empty selected namespace groups
  - nodes remaining visible when none of the selected namespaces match
  - grouped scoped output still rendering matching namespace groups alphabetically
  - grouped scoped output preserving pod-order parity with the grouped path
  - unchanged single-namespace scoped flat output
  - README examples no longer claiming empty selected namespace groups render

## Phase 2: Implementation Outline

1. Refactor `srv/view.go` so grouped namespace rendering only emits namespace headings that have one or more matching pods on the current node.
2. Keep `cmd/root.go` configuration behavior aligned with the existing grouped-scoped rendering contract and adjust `cmd/root_test.go` only if the config expectations need refinement.
3. Update `srv/view_test.go` to replace the prior empty-group expectation with omission behavior and add coverage that nodes still render without namespace headings when they have zero matching pods.
4. Update `README.md` examples and explanatory text to reflect that empty selected namespace groups are no longer shown.
5. Run `make test` and verify the corrected grouped output matches the documented contract.

## Post-Design Constitution Check

- **I. Operational Proof Over Intent**: Pass. The design remains directly verifiable through output assertions and manual command runs against the documented issue scenario.
- **II. CLI-First, Script-Safe Interfaces**: Pass. The CLI surface stays stable and the output contract becomes more accurate.
- **III. Tests and Validation Are Mandatory**: Pass. The design requires regression tests and `make test`.
- **IV. Documentation Matches User Behavior**: Pass. README corrections are part of the same work unit.
- **V. Small, Compatible, Repository-Native Changes**: Pass. The fix stays localized to the established renderer, tests, and docs.

## Complexity Tracking

No constitution violations require justification.
