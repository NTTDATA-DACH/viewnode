# Implementation Plan: Remove redundant namespace prefixes from single-namespace pod rows

**Branch**: `066-fix-single-namespace-prefix` | **Date**: 2026-04-04 | **Spec**: [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/066-fix-single-namespace-prefix/spec.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/066-fix-single-namespace-prefix/spec.md)
**Input**: Feature specification from `/specs/066-fix-single-namespace-prefix/spec.md`

## Summary

Correct single-namespace `viewnode` output so pod rows no longer repeat the selected namespace inline when that namespace is already shown in the header. The change stays in the existing flat render path in `srv/view.go`, preserves grouped multi-namespace and all-namespaces behavior, adds focused regression coverage in `srv/view_test.go`, and updates `README.md` so the user-visible examples match the corrected output.

## Technical Context

**Language/Version**: Go 1.25.0 (toolchain go1.25.7)  
**Primary Dependencies**: `github.com/spf13/cobra`, `k8s.io/client-go`, `k8s.io/apimachinery`, `github.com/sirupsen/logrus`, `github.com/stretchr/testify`  
**Storage**: Kubeconfig files plus live Kubernetes API data for node and pod discovery  
**Testing**: `go test`, `make test`  
**Target Platform**: Cross-platform CLI environments supported by Go and Kubernetes client-go  
**Project Type**: CLI  
**Performance Goals**: Preserve current output latency and avoid any additional Kubernetes API calls by keeping the fix inside the existing pod-row formatter  
**Constraints**: Keep the CLI surface unchanged, preserve the namespace header, preserve node grouping and pod ordering, preserve multi-namespace and all-namespaces behavior, preserve container and metrics output, and update `README.md` for the user-visible single-namespace change  
**Scale/Scope**: One flat-rendering path in `srv/view.go`, focused regression coverage in `srv/view_test.go`, optional confirmation in `cmd/root_test.go` if needed, and one README correction for single-namespace output examples

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Operational Proof Over Intent**: Pass. The fix changes directly observable CLI output and can be verified through output-focused tests and a manual single-namespace command run.
- **II. CLI-First, Script-Safe Interfaces**: Pass. No command names or flags change; only redundant text in the output is removed while preserving deterministic CLI behavior.
- **III. Tests and Validation Are Mandatory**: Pass with required work. The plan adds or updates regression coverage for the single-namespace flat path and requires `make test`.
- **IV. Documentation Matches User Behavior**: Pass with required work. `README.md` must be updated in the same change because the single-namespace examples currently imply the old output shape.
- **V. Small, Compatible, Repository-Native Changes**: Pass. The design reuses the current renderer seam and avoids changing command parsing, pod loading, or grouping behavior for other modes.

## Project Structure

### Documentation (this feature)

```text
specs/066-fix-single-namespace-prefix/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── single-namespace-output-contract.md
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

**Structure Decision**: Keep the fix in the flat single-namespace rendering path in `srv/view.go`, keep `cmd/root.go` unchanged unless verification shows a configuration mismatch, cover the behavior with output-level tests in `srv/view_test.go`, and update the single-namespace examples in `README.md`.

## Phase 0: Research

Research findings are recorded in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/066-fix-single-namespace-prefix/research.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/066-fix-single-namespace-prefix/research.md). Key decisions:

- Keep the fix in `srv/view.go` because the bug is in the flat pod-row formatter, not in namespace parsing or pod loading.
- Continue using the existing namespace header for explicit single-namespace runs and remove only the redundant inline prefix from pod rows.
- Preserve grouped multi-namespace and all-namespaces behavior exactly as it is today.
- Preserve container, timing, and metrics rendering on the same pod rows after the prefix removal.
- Update README examples so the documentation reflects the corrected single-namespace output.

## Phase 1: Design & Contracts

### Data Model

Data entities and flat-render relationships are defined in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/066-fix-single-namespace-prefix/data-model.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/066-fix-single-namespace-prefix/data-model.md).

### Contracts

The corrected user-visible CLI output contract is documented in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/066-fix-single-namespace-prefix/contracts/single-namespace-output-contract.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/066-fix-single-namespace-prefix/contracts/single-namespace-output-contract.md).

### Quickstart

Implementation and validation steps are captured in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/066-fix-single-namespace-prefix/quickstart.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/066-fix-single-namespace-prefix/quickstart.md).

### Design Decisions

- Keep `cmd/root.go` namespace parsing and `ViewNodeDataConfig` behavior unchanged unless tests prove otherwise.
- Refine the flat single-namespace pod-row formatter in `srv/view.go` so it renders just the pod name and phase while the namespace header remains the canonical namespace label.
- Preserve the current grouped rendering path for multi-namespace and all-namespaces modes.
- Preserve pod ordering, node grouping, summary lines, and unscheduled pod behavior.
- Preserve inline and tree container output, timing output, and metrics output when the namespace prefix is removed.
- Add regression coverage for:
  - single-namespace pod rows no longer showing `<namespace>: `
  - single-namespace output preserving the namespace header and node/pod counts
  - container and auxiliary pod-row details remaining compatible after the prefix removal
  - multi-namespace and grouped output remaining unchanged where relevant
  - README examples reflecting the corrected single-namespace output

## Phase 2: Implementation Outline

1. Refactor `srv/view.go` so the flat single-namespace pod-row path omits the redundant inline namespace prefix.
2. Update `srv/view_test.go` to replace the current expectation of `team-a: pod-name` for single-namespace runs with the corrected `pod-name` output and preserve coverage for related pod-row details.
3. Review `cmd/root.go` and `cmd/root_test.go` only as needed to confirm the fix does not require command configuration changes.
4. Update `README.md` examples and explanatory text for single-namespace output.
5. Run `make test` and verify the corrected single-namespace output matches the documented contract.

## Post-Design Constitution Check

- **I. Operational Proof Over Intent**: Pass. The design is directly verifiable through output tests and a manual single-namespace command run.
- **II. CLI-First, Script-Safe Interfaces**: Pass. The CLI surface remains stable and the output becomes less redundant.
- **III. Tests and Validation Are Mandatory**: Pass. The design requires regression tests and `make test`.
- **IV. Documentation Matches User Behavior**: Pass. README updates are part of the same work unit.
- **V. Small, Compatible, Repository-Native Changes**: Pass. The fix stays localized to the established renderer, tests, and docs.

## Complexity Tracking

No constitution violations require justification.
