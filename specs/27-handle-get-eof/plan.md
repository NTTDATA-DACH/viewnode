# Implementation Plan: Clearer EOF error guidance

**Branch**: `27-handle-get-eof` | **Date**: 2026-03-24 | **Spec**: [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/spec.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/spec.md)
**Input**: Feature specification from `/specs/27-handle-get-eof/spec.md`

## Summary

Improve the root `viewnode` failure path so scoped Kubernetes API `Get` requests that surface as EOF produce a clearer fatal message that hints proxy configuration may be involved, while preserving the original low-level error details and leaving non-EOF or out-of-scope EOF failures unchanged. The implementation stays inside the existing `srv` error-decoration layer and the root command’s fatal handling, with focused regression tests and no new command surface.

## Technical Context

**Language/Version**: Go 1.25.0 (toolchain go1.25.7)  
**Primary Dependencies**: `github.com/spf13/cobra`, `k8s.io/client-go`, `k8s.io/apimachinery`, `github.com/sirupsen/logrus`, `github.com/stretchr/testify`  
**Storage**: Kubeconfig files plus live Kubernetes API data for node and pod retrieval  
**Testing**: `go test`, `make test`  
**Target Platform**: Cross-platform CLI environments supported by Go and Kubernetes client-go  
**Project Type**: CLI  
**Performance Goals**: Preserve current interactive command latency and Kubernetes API call counts by classifying existing errors locally without retries or additional network requests  
**Constraints**: Limit the new guidance to cluster retrieval API `Get` failures that surface as EOF, preserve original low-level error details, keep repeated failures deterministic, avoid changing non-EOF messaging semantics, and follow the current `srv` and `cmd/root.go` error-handling patterns  
**Scale/Scope**: One error-decoration file, one root command fatal-handling path, targeted tests, a small CLI contract, and no new command/flag additions

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Operational Proof Over Intent**: Pass. The feature changes only the observable failure message for a known runtime condition and keeps the command’s fatal outcome explicit.
- **II. CLI-First, Script-Safe Interfaces**: Pass. The design preserves the existing `viewnode` command surface and keeps behavior in the established Cobra root execution path.
- **III. Tests and Validation Are Mandatory**: Pass with required work. The plan adds automated coverage for EOF decoration and root fatal messaging, then finishes with `make test`.
- **IV. Documentation Matches User Behavior**: Pass with explicit scope note. No command names, flags, or examples change, so README command reference updates are not required for this feature unless implementation adds a troubleshooting note.
- **V. Small, Compatible, Repository-Native Changes**: Pass. The work reuses the current typed error decoration pattern in `srv/api_errors.go` and the existing root fatal switch rather than introducing a parallel error framework.

## Project Structure

### Documentation (this feature)

```text
specs/27-handle-get-eof/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── eof-error-guidance-contract.md
└── tasks.md
```

### Source Code (repository root)

```text
cmd/
├── root.go
└── root_test.go

srv/
├── api.go
├── api_errors.go
├── api_mock.go
├── node.go
├── node_test.go
└── view_test.go

README.md
Makefile
go.mod
```

**Structure Decision**: Keep EOF classification in `srv/api_errors.go`, where Kubernetes API errors are already decorated into typed errors, and keep user-facing fatal message selection in `cmd/root.go`, where root resource-loading failures are already translated into CLI messages. Add regression tests close to those files instead of introducing a new shared error package.

## Phase 0: Research

Research findings are recorded in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/research.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/research.md). Key decisions:

- Detect scoped EOF failures in `srv/api_errors.go` with a dedicated decorated error type that still wraps the original error.
- Handle the new decorated error in `cmd/root.go` before the generic fatal path so the CLI can add a proxy hint without changing unrelated errors.
- Preserve the original low-level error details in the final fatal message and keep the proxy hint phrased as a possible cause.
- Keep README command reference content unchanged because the CLI surface does not change.

## Phase 1: Design & Contracts

### Data Model

Data entities and validation rules are defined in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/data-model.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/data-model.md).

### Contracts

User-facing CLI failure behavior is documented in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/contracts/eof-error-guidance-contract.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/contracts/eof-error-guidance-contract.md).

### Quickstart

Implementation and validation steps are captured in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/quickstart.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/27-handle-get-eof/quickstart.md).

### Design Decisions

- Extend `srv/api_errors.go` with a new decorated error type for scoped cluster retrieval API `Get` EOF failures, following the existing `UnauthorizedError` and `NodesIsForbiddenError` pattern.
- Limit decoration to the root retrieval path described by the spec:
  - Kubernetes API `Get`-style request failures surfaced during node or pod retrieval
  - messages that still carry the original `EOF` detail
  - no decoration for unrelated EOFs or non-EOF retrieval failures
- Add a dedicated branch in `cmd/root.go` so the root command emits a clearer fatal message that:
  - keeps the current resource context
  - states proxy configuration may be the cause
  - preserves the original low-level error text
  - repeats on every occurrence rather than suppressing duplicates
- Add or update tests to verify:
  - decorated EOF errors are recognized only for the intended scoped messages
  - non-EOF and out-of-scope EOF cases continue through existing behavior
  - the root fatal path renders proxy guidance alongside the original error details
  - `make test` remains the final validation step

## Phase 2: Implementation Outline

1. Add scoped EOF error decoration in `srv/api_errors.go`, keeping the original error wrapped and accessible.
2. Extend `cmd/root.go` to recognize the new decorated error and emit the proxy-guidance fatal message before the generic fallback path.
3. Add focused regression coverage in `srv` and `cmd` tests for scoped EOF detection, non-EOF preservation, and user-facing fatal output behavior.
4. Confirm whether a README troubleshooting note is warranted; otherwise leave `README.md` unchanged because the CLI surface is stable.
5. Run `make test` and verify the final output still preserves original error details while adding the new hint.

## Post-Design Constitution Check

- **I. Operational Proof Over Intent**: Pass. The design changes the emitted fatal message for an observable failure condition without masking command failure.
- **II. CLI-First, Script-Safe Interfaces**: Pass. No flags, commands, or invocation patterns change, and the fatal output remains deterministic for scripts and users.
- **III. Tests and Validation Are Mandatory**: Pass. The design requires targeted tests plus `make test`.
- **IV. Documentation Matches User Behavior**: Pass. The documentation decision is explicit: no command reference update is needed unless implementation introduces a troubleshooting note.
- **V. Small, Compatible, Repository-Native Changes**: Pass. The plan stays inside the repository’s current error decoration and root fatal handling patterns.

## Complexity Tracking

No constitution violations require justification.
