# Implementation Plan: Refactor `ctx list` and `ns list` commands

**Branch**: `48-refactor-ctx-list-ns-list-commands` | **Date**: 2026-03-15 | **Spec**: [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/48-refactor-ctx-list-ns-list-commands/spec.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/48-refactor-ctx-list-ns-list-commands/spec.md)
**Input**: Feature specification from `/specs/48-refactor-ctx-list-ns-list-commands/spec.md`

**Note**: `.specify/scripts/bash/setup-plan.sh --json` could not be used directly because the upstream Spec Kit helper currently rejects this repository's issue-based branch format (`48-...`) and expects `001-...`. Artifact paths were resolved from the active feature directory while preserving the required plan workflow.

## Summary

Refactor `viewnode ctx list` and `viewnode ns list` so both commands build and render display entries through shared list-preparation logic, while preserving current command shape, output format, active-item semantics, and error behavior. The design keeps the existing Cobra package layout, introduces explicit list-entry types for contexts and namespaces, and extends tests to cover shared behavior and namespace default handling.

## Technical Context

**Language/Version**: Go 1.25.0 (toolchain go1.25.7)  
**Primary Dependencies**: `github.com/spf13/cobra`, `k8s.io/client-go`, `k8s.io/apimachinery`, `github.com/stretchr/testify`  
**Storage**: Kubeconfig files plus live Kubernetes API data for namespace discovery  
**Testing**: `go test`, `make test`  
**Target Platform**: Cross-platform CLI environments supported by Go and Kubernetes client-go  
**Project Type**: CLI  
**Performance Goals**: Preserve current interactive listing latency and avoid adding extra Kubernetes API calls beyond the existing single namespace list request  
**Constraints**: Preserve user-visible output format, keep alphabetical ordering, follow existing `cmd/ctx` and `cmd/ns` structure, and avoid parallel command implementations  
**Scale/Scope**: Two existing list subcommands, one shared list-preparation flow, and focused command-level test updates

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Code Quality Is a Release Gate**: Pass. The plan keeps the change limited to the existing command packages and extracts only duplicated list-preparation behavior.
- **II. Tests Prove Behavior Before Merge**: Pass with required work. The plan adds and updates automated tests for both command flows and finishes with `make test`.
- **III. CLI Experience Must Stay Consistent**: Pass. The feature is behavior-preserving and does not add or rename user-facing commands or flags, so no README change is expected.
- **IV. Performance Regressions Require Evidence**: Pass. The design keeps the same data sources and request counts; verification will confirm no extra namespace API requests are introduced.
- **V. Keep Changes Small, Observable, and Reversible**: Pass. The change is a narrow refactor with direct command tests and no persistent side effects.

## Project Structure

### Documentation (this feature)

```text
specs/48-refactor-ctx-list-ns-list-commands/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── list-command-contract.md
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
├── ns/
│   ├── getcurrent.go
│   ├── list.go
│   ├── list_test.go
│   ├── ns.go
│   ├── ns_test.go
│   └── setcurrent.go
├── root.go
└── root_test.go

README.md
Makefile
go.mod
```

**Structure Decision**: Keep all changes inside the existing Go CLI layout under `cmd/ctx`, `cmd/ns`, and their tests. Shared list-entry preparation should live alongside the current command packages or in the existing config-oriented command support area only if reuse stays clear and limited.

## Phase 0: Research

Research findings are recorded in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/48-refactor-ctx-list-ns-list-commands/research.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/48-refactor-ctx-list-ns-list-commands/research.md). Key decisions:

- Use explicit display-entry structs for contexts and namespaces instead of anonymous loops.
- Share list-preparation behavior only where semantics are truly common: sorting and active-state rendering inputs.
- Preserve namespace activity semantics by deriving the active namespace from kubeconfig context data and defaulting empty values to `default`.

## Phase 1: Design & Contracts

### Data Model

Data entities and validation rules are defined in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/48-refactor-ctx-list-ns-list-commands/data-model.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/48-refactor-ctx-list-ns-list-commands/data-model.md).

### Contracts

User-facing CLI behavior for the affected commands is documented in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/48-refactor-ctx-list-ns-list-commands/contracts/list-command-contract.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/48-refactor-ctx-list-ns-list-commands/contracts/list-command-contract.md).

### Quickstart

Implementation and validation steps are captured in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/48-refactor-ctx-list-ns-list-commands/quickstart.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/48-refactor-ctx-list-ns-list-commands/quickstart.md).

### Design Decisions

- Introduce a small display-entry model for each command domain:
  - context entries contain `name` and `active`
  - namespace entries contain `name` and `active`
- Extract duplicate “build ordered active-marked entries from raw names” behavior into a shared helper with command-specific inputs.
- Keep Kubernetes access boundaries unchanged:
  - `ctx list` still reads from raw kubeconfig
  - `ns list` still queries namespaces from the cluster once and compares against the active namespace
- Add or update tests to verify:
  - alphabetical ordering
  - exactly one active marker when present
  - `default` namespace fallback
  - error propagation for kubeconfig and namespace lookup failures

## Phase 2: Implementation Outline

1. Add explicit list-entry types and a shared helper for ordered active-state preparation in the most local package scope that avoids new parallel structures.
2. Refactor `cmd/ctx/list.go` to build context entries through the shared helper while preserving existing output.
3. Refactor `cmd/ns/list.go` to build namespace entries through the shared helper and use kubeconfig-default namespace semantics for active-state comparison.
4. Expand command tests to cover shared behavior and refactor-sensitive regressions.
5. Run `make test` and confirm no documentation changes are required because the CLI surface remains unchanged.

## Post-Design Constitution Check

- **I. Code Quality Is a Release Gate**: Pass. The design removes duplication without broadening responsibility.
- **II. Tests Prove Behavior Before Merge**: Pass. The plan explicitly requires command test updates and `make test`.
- **III. CLI Experience Must Stay Consistent**: Pass. The contract locks output and command shape to current behavior.
- **IV. Performance Regressions Require Evidence**: Pass. No extra API calls are introduced; verification checks focus on behavioral parity.
- **V. Keep Changes Small, Observable, and Reversible**: Pass. The work remains isolated to two commands and their tests.

## Complexity Tracking

No constitution violations require justification.
