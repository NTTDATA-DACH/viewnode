# Implementation Plan: Context list filtering

**Branch**: `069-add-ctx-filter` | **Date**: 2026-04-07 | **Spec**: [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/069-add-ctx-filter/spec.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/069-add-ctx-filter/spec.md)
**Input**: Feature specification from `/specs/069-add-ctx-filter/spec.md`

## Summary

Add case-insensitive context-name filtering to `viewnode ctx list` through a local `--filter` / `-f` flag, matching the ergonomics already established for `viewnode ns list`. The design keeps kubeconfig context discovery in the current `cmd/ctx` command path, applies filtering client-side after the existing raw kubeconfig load, preserves sorted output and active-marker semantics for displayed rows, emits the clarified no-match message `no contexts matched filter "<value>"`, and covers the change with focused command tests plus README updates.

## Technical Context

**Language/Version**: Go 1.25.0 (toolchain go1.25.7)  
**Primary Dependencies**: `github.com/spf13/cobra`, `k8s.io/client-go`, `k8s.io/apimachinery`, `github.com/stretchr/testify`, `github.com/sirupsen/logrus`  
**Storage**: Kubeconfig files for raw context discovery and current-context selection state  
**Testing**: `go test`, `make test`  
**Target Platform**: Cross-platform CLI environments supported by Go and Kubernetes client-go  
**Project Type**: CLI  
**Performance Goals**: Preserve current interactive latency and avoid adding extra kubeconfig loads or Kubernetes API calls beyond the existing single raw kubeconfig read for `ctx list`  
**Constraints**: Preserve existing command names and aliases, keep output ordering deterministic, preserve marker formatting for displayed rows, treat no-match results as successful outcomes, update `README.md`, and follow existing `cmd/ctx`, `cmd/internal/listing`, and Cobra command-tree patterns without introducing parallel structures  
**Scale/Scope**: One existing `ctx list` command, one focused filtered output path, targeted tests in the context command package, and README command reference/example updates

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Operational Proof Over Intent**: Pass. The planned change affects a read-only listing command and preserves observable success and failure states, including an explicit no-match success message and unchanged kubeconfig error propagation.
- **II. CLI-First, Script-Safe Interfaces**: Pass. The feature extends the existing Cobra subcommand with a local flag while preserving command names, aliases, output format, and script-friendly text output.
- **III. Tests and Validation Are Mandatory**: Pass with required work. The plan adds focused automated tests around filter behavior and ends with `make test`.
- **IV. Documentation Matches User Behavior**: Pass with required work. The plan includes `README.md` updates for the new flag, shorthand, examples, and no-match wording.
- **V. Small, Compatible, Repository-Native Changes**: Pass. The design stays within the current `cmd/ctx` and shared listing patterns and intentionally mirrors the repository’s earlier namespace-filter implementation instead of introducing a new abstraction layer.

## Project Structure

### Documentation (this feature)

```text
specs/069-add-ctx-filter/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── ctx-filter-contract.md
└── tasks.md
```

### Source Code (repository root)

```text
cmd/
├── config/
│   ├── config.go
│   ├── setup.go
│   └── setup_test.go
├── ctx/
│   ├── ctx.go
│   ├── ctx_test.go
│   ├── getcurrent.go
│   ├── list.go
│   └── setcurrent.go
├── internal/
│   └── listing/
│       ├── entries.go
│       └── entries_test.go
├── ns/
│   ├── list.go
│   └── list_test.go
└── root.go

README.md
Makefile
go.mod
```

**Structure Decision**: Keep context filter flag parsing and filtered context rendering in `cmd/ctx/list.go`, reuse the existing shared ordering and active-marker helper in `cmd/internal/listing/entries.go`, place regression coverage in `cmd/ctx/ctx_test.go` with only minimal helper extraction if needed, and update `README.md` where the context command examples already live.

## Phase 0: Research

Research findings are recorded in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/069-add-ctx-filter/research.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/069-add-ctx-filter/research.md). Key decisions:

- Add `--filter` / `-f` as a local `ctx list` flag instead of reusing root flags or changing the command shape.
- Filter contexts client-side after the existing raw kubeconfig read.
- Treat empty filter input as no filter.
- Use the exact no-match message `no contexts matched filter "<value>"` while keeping success status.
- Preserve filtered output as matching rows only, with no active marker if the current context is filtered out.

## Phase 1: Design & Contracts

### Data Model

Data entities and validation rules are defined in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/069-add-ctx-filter/data-model.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/069-add-ctx-filter/data-model.md).

### Contracts

User-facing CLI behavior is documented in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/069-add-ctx-filter/contracts/ctx-filter-contract.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/069-add-ctx-filter/contracts/ctx-filter-contract.md).

### Quickstart

Implementation and validation steps are captured in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/069-add-ctx-filter/quickstart.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/069-add-ctx-filter/quickstart.md).

### Design Decisions

- Extend `cmd/ctx/list.go` with a local filter flag and normalize the filter for case-insensitive substring matching.
- Keep context retrieval unchanged:
  - one raw kubeconfig read through the existing config initialization path
  - current context still derived from `rawConfig.CurrentContext`
- Filter context names after retrieval and before `listing.PrepareContextEntries(...)` so ordering and active-marker logic continue to flow through the existing shared helper.
- Return the clarified success-path no-match message `no contexts matched filter "<value>"` when filtering removes all rows.
- Preserve filtered output as only the matching contexts, even when the active context does not match and no displayed row is marked active.
- Add or update tests to verify:
  - `--filter` and `-f` produce identical context results
  - case-insensitive context matching
  - empty filter behaves like no filter
  - filtered no-match output is exact and successful
  - filtered output omits the active marker when the current context is not present
  - README examples and command reference stay aligned with the CLI surface

## Phase 2: Implementation Outline

1. Add local context filter flag parsing, filter normalization, and case-insensitive matching in `cmd/ctx/list.go` while preserving the current kubeconfig initialization and raw-config error flow.
2. Reuse or minimally extend shared listing behavior in `cmd/internal/listing/entries.go` only if needed to keep sorted output and active-marker handling repository-native and testable.
3. Extend `cmd/ctx/ctx_test.go` with realistic command execution coverage for long flag, short flag, empty filter, no-match messaging, current-context-filtered-out behavior, and kubeconfig error handling.
4. Update `README.md` command reference and examples to document `viewnode ctx list --filter` / `-f` and the clarified no-match behavior.
5. Run `make test` and confirm the user-visible contract matches the documented examples.

## Post-Design Constitution Check

- **I. Operational Proof Over Intent**: Pass. The design preserves concrete success and failure outcomes and makes the filtered empty state explicitly observable through stable CLI text.
- **II. CLI-First, Script-Safe Interfaces**: Pass. The contract is limited to the requested local flag addition and explicit no-match output while preserving the existing command tree and formatting conventions.
- **III. Tests and Validation Are Mandatory**: Pass. The plan requires targeted command-level tests and `make test`.
- **IV. Documentation Matches User Behavior**: Pass. README updates remain part of the implementation scope and validation contract.
- **V. Small, Compatible, Repository-Native Changes**: Pass. The work remains tightly scoped to the current context command path, shared listing helper, tests, and documentation.

## Complexity Tracking

No constitution violations require justification.
