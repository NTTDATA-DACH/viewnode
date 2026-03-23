# Implementation Plan: Namespace list filtering

**Branch**: `53-add-ns-filter` | **Date**: 2026-03-23 | **Spec**: [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/53-add-ns-filter/spec.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/53-add-ns-filter/spec.md)
**Input**: Feature specification from `/specs/53-add-ns-filter/spec.md`

## Summary

Add case-insensitive namespace-name filtering to `viewnode ns list` through a local `--filter` / `-f` flag, while also replacing generic filtered empty-state output with clearer success-path no-match messaging for both `ns list` and the existing root node-filter flow. The design keeps namespace discovery and root filtering in their current packages, applies filtering client-side after existing API calls, and covers the change with command/output tests plus README updates.

## Technical Context

**Language/Version**: Go 1.25.0 (toolchain go1.25.7)  
**Primary Dependencies**: `github.com/spf13/cobra`, `k8s.io/client-go`, `k8s.io/apimachinery`, `github.com/stretchr/testify`, `github.com/sirupsen/logrus`  
**Storage**: Kubeconfig files plus live Kubernetes API data for namespace and node/pod discovery  
**Testing**: `go test`, `make test`  
**Target Platform**: Cross-platform CLI environments supported by Go and Kubernetes client-go  
**Project Type**: CLI  
**Performance Goals**: Preserve current interactive latency and avoid adding extra Kubernetes API calls beyond the existing single namespace list request and current node/pod loading flow  
**Constraints**: Preserve existing command names and root flag behavior, keep `ns list` output ordering and active markers stable, treat no-match results as successful outcomes, update `README.md`, and follow existing `cmd/ns`, `cmd/config`, and `srv` patterns without introducing parallel command structures  
**Scale/Scope**: One existing namespace list command, one existing root filtered printout path, focused test additions, and README command reference/example updates

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Code Quality Is a Release Gate**: Pass. The plan keeps the feature inside the current Cobra command and printout layers and avoids introducing a broader abstraction than the change requires.
- **II. Tests Prove Behavior Before Merge**: Pass with required work. The plan adds automated coverage for namespace filtering, no-match output, and command flag behavior, then finishes with `make test`.
- **III. CLI Experience Must Stay Consistent**: Pass. Existing command names, aliases, and non-filtered output remain stable, while README updates capture the user-facing namespace flag addition and message change.
- **IV. Performance Regressions Require Evidence**: Pass. Namespace filtering remains a client-side refinement over the current single list call, and the plan explicitly avoids extra Kubernetes requests.
- **V. Keep Changes Small, Observable, and Reversible**: Pass. The work is narrowly scoped to `ns list`, root filtered empty-state messaging, tests, and documentation.

## Project Structure

### Documentation (this feature)

```text
specs/53-add-ns-filter/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── filter-command-contract.md
└── tasks.md
```

### Source Code (repository root)

```text
cmd/
├── config/
│   ├── config.go
│   ├── setup.go
│   └── setup_test.go
├── internal/
│   └── listing/
│       ├── entries.go
│       └── entries_test.go
├── ns/
│   ├── getcurrent.go
│   ├── list.go
│   ├── list_test.go
│   ├── ns.go
│   ├── ns_test.go
│   └── setcurrent.go
├── root.go
└── root_test.go

srv/
├── node.go
├── pod.go
├── view.go
└── view_test.go

README.md
Makefile
go.mod
```

**Structure Decision**: Keep namespace flag parsing and filtered namespace rendering in `cmd/ns/list.go`, keep root filtered empty-state messaging in the existing `srv/view.go` printout path with only the minimal plumbing needed from `cmd/root.go`, and place all regression coverage in the current command and service test files.

## Phase 0: Research

Research findings are recorded in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/53-add-ns-filter/research.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/53-add-ns-filter/research.md). Key decisions:

- Add `--filter` / `-f` as a local `ns list` flag instead of reusing or changing the root node-filter flag.
- Apply case-insensitive substring matching client-side after the existing namespace list call.
- Treat empty namespace filters as no filter.
- Use clearer no-match success messages for filtered namespace and node flows.

## Phase 1: Design & Contracts

### Data Model

Data entities and validation rules are defined in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/53-add-ns-filter/data-model.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/53-add-ns-filter/data-model.md).

### Contracts

User-facing CLI behavior is documented in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/53-add-ns-filter/contracts/filter-command-contract.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/53-add-ns-filter/contracts/filter-command-contract.md).

### Quickstart

Implementation and validation steps are captured in [/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/53-add-ns-filter/quickstart.md](/Users/adam.boczek/Development/Workspace/NTTDATA/viewnode/specs/53-add-ns-filter/quickstart.md).

### Design Decisions

- Extend `cmd/ns/list.go` with a local filter flag and normalize the filter for case-insensitive matching.
- Keep namespace retrieval unchanged:
  - one namespace list API call
  - active namespace still derived from raw kubeconfig with empty meaning `default`
- Filter namespace names after retrieval and before `listing.PrepareNamespaceEntries(...)` so ordering and active markers continue to use the existing helper.
- Add a dedicated filtered no-match success path to `ns list` output instead of rendering an empty list.
- Update `srv/view.go` so root node filtering can emit a filter-aware no-match message instead of the generic `no nodes to display...`, with only enough context added from `cmd/root.go` to choose the correct wording.
- Add or update tests to verify:
  - `--filter` and `-f` produce identical namespace results
  - case-insensitive namespace matching
  - empty filter behaves like no filter
  - filtered no-match output for namespace and node flows remains successful
  - README examples and command reference stay aligned with the CLI surface

## Phase 2: Implementation Outline

1. Add local namespace filter flag parsing and case-insensitive filtering in `cmd/ns/list.go`, preserving existing kubeconfig initialization and active-namespace semantics.
2. Extend `cmd/ns/list_test.go` with command-tree execution coverage for long flag, short flag, empty filter, no-match messaging, and namespace discovery errors.
3. Introduce minimal filter-aware empty-state context in the root printout flow by updating `cmd/root.go` and `srv/view.go`, then cover it in `srv/view_test.go`.
4. Update `README.md` command reference and examples to document `viewnode ns list --filter` / `-f` and the clearer filtered empty-state behavior.
5. Run `make test` and confirm the user-visible contract matches the documented examples.

## Post-Design Constitution Check

- **I. Code Quality Is a Release Gate**: Pass. The design stays within existing packages and relies on current helpers rather than creating parallel filtering paths.
- **II. Tests Prove Behavior Before Merge**: Pass. The plan requires command-level tests, output tests, and `make test`.
- **III. CLI Experience Must Stay Consistent**: Pass. The contract limits user-visible change to the requested namespace flag and clearer filtered empty-state wording, with README updates included.
- **IV. Performance Regressions Require Evidence**: Pass. Namespace filtering is computed locally after the current single list call, and no extra cluster round trips are added.
- **V. Keep Changes Small, Observable, and Reversible**: Pass. The work remains focused and easy to validate or revert by file area.

## Complexity Tracking

No constitution violations require justification.
