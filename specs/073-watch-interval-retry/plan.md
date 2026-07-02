# Implementation Plan: Optional Refresh Interval and Retryable Watch Loop

**Branch**: `073-watch-interval-retry` | **Date**: 2026-06-26 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/073-watch-interval-retry/spec.md`

**Issue**: [#73](https://github.com/NTTDATA-DACH/viewnode/issues/73)

## Summary

Replace the boolean global `--watch` / `-w` flag with an optional positive-integer seconds value (Cobra `NoOptDefVal = "1"`), and replace the current ticker-driven `schedule` loop in `cmd/root.go` with a sleep-after-refresh loop that:
- runs exactly one refresh in one-shot mode,
- enters a watch loop only after the first refresh succeeds,
- treats every subsequent refresh error as transient and renders a labeled error frame for that refresh,
- sleeps the configured interval after each refresh completes (no catch-up bursts),
- stops cleanly on Ctrl-C / SIGINT / SIGTERM.

The change lands entirely inside `cmd/` (CLI Surface / Watch Orchestration sub-capability per the architecture views). `cmd/config` and `srv/` capability contracts are unchanged. README and root command help text are updated in the same change.

## Technical Context

**Language/Version**: Go 1.25.0 (toolchain go1.25.7) per `go.mod`.

**Primary Dependencies**: `github.com/spf13/cobra` v1.10.x (flag + `NoOptDefVal`), `github.com/sirupsen/logrus` + `github.com/antonfisher/nested-logrus-formatter` (existing logging), `k8s.io/client-go` + `k8s.io/apimachinery` + `k8s.io/metrics` v0.35.2 (unchanged), `github.com/stretchr/testify` (tests). No new dependencies introduced.

**Storage**: N/A. Tool is stateless beyond per-invocation in-memory state and the operator-owned kubeconfig (not touched by this feature).

**Testing**: `go test -race -covermode=atomic ./...` (existing `make test`). Unit tests next to changed code: `cmd/root_test.go` and a new `cmd/watch_test.go` (or equivalent file co-located with the extracted loop helpers). No live cluster required for new tests — error injection is done through an injected refresh function or a fake error channel.

**Target Platform**: All platforms `viewnode` already supports — `linux/amd64`, `darwin/amd64`, `darwin/arm64`, `linux/arm64` (goreleaser) and `windows/amd64` (krew manifest). Signal handling must work on Unix and Windows (Cobra's existing patterns + `os/signal` for `SIGINT`; `SIGTERM` Unix-only).

**Project Type**: CLI tool (single Go module, single binary, distributed standalone + as a krew plugin).

**Performance Goals**: Per-refresh wall-clock latency is dominated by Kubernetes API + metrics API round trips; the loop adds only the configured interval as sleep. No measurable per-refresh CPU/memory regression vs current implementation.

**Constraints**:
- Preserve existing one-shot behavior byte-for-byte (FR-011, SC-005).
- Keep the flag global on the root command (FR-001).
- Sleep-after-refresh: next refresh starts no sooner than the configured interval after the previous refresh completes (FR-007, SC-003).
- First-refresh failure exits with the same non-zero code as `viewnode` without `--watch` (clarification Q1).
- Mid-session refresh errors must not exit the loop (clarification Q2).
- Error frame replaces the failing refresh's normal frame; includes refresh timestamp + underlying error text (clarification Q3).

**Scale/Scope**: Behavior change is localized to `cmd/root.go` and a new small test file. No package re-org. No public API change to `srv` or `cmd/config`.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Stable CLI Contracts**: The `--watch` / `-w` flag remains global and keeps both forms. The user-visible change is from "boolean, fixed 1s cadence" to "optional integer seconds, default 1, validated". Non-watch behavior is unchanged. README and root command long-help are updated in the same change. ✅ Pass.
- **Kubernetes Compatibility Boundaries**: No changes to `cmd/config` or `srv` capability contracts. Cluster reads, metrics API tolerance, and kubeconfig handling are unchanged. The loop classifies refresh errors operationally (transient after first success) but does not change which Kubernetes errors are surfaced. ✅ Pass.
- **Test-Backed Behavior Changes**: New unit tests cover (1) interval flag parsing (valid default, valid explicit, rejected invalid), (2) single-run path runs exactly one refresh and propagates exit code, (3) sleep-after-refresh semantics under an injectable clock, (4) post-first-success error tolerance and error-frame rendering, (5) clean stop on cancellation signal. Existing `cmd/root_test.go` is updated to reflect the new flag type and loop entry point. ✅ Pass.
- **Small, Traceable Delivery**: Spec defines 3 prioritized user stories (P1, P1, P2), each independently testable. Tasks (next phase) will be ordered by dependency: flag parsing → single-run extraction → loop → rendering → docs. No parallel command/configuration/service structure is introduced. ✅ Pass.
- **Documentation and Operability Parity**: README's `--watch` section, root command long-help, and any usage examples will be updated in the same change. Error-frame format is defined in the plan (timestamp + error text) so help text and tests describe the same thing. ✅ Pass.

No violations recorded → Complexity Tracking section is empty.

## Project Structure

### Documentation (this feature)

```text
specs/073-watch-interval-retry/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/
│   └── cli.md           # CLI contract for the updated --watch flag
└── checklists/
    └── requirements.md  # From /speckit-specify
```

### Source Code (repository root)

```text
cmd/
├── root.go              # Watch flag changes here (BoolVarP → IntVarP with NoOptDefVal); schedule() is replaced/restructured
├── root_test.go         # Updated for new flag type and validation
├── watch.go             # NEW (optional file split): single-run path and watch loop helpers, kept inside `cmd` package
├── watch_test.go        # NEW: tests for flag parsing, sleep-after-refresh, transient-error tolerance, clean stop
├── config/              # UNCHANGED
├── ctx/                 # UNCHANGED
├── ns/                  # UNCHANGED
└── internal/listing/    # UNCHANGED

srv/                     # UNCHANGED (no API or rendering contract change)
utils/                   # UNCHANGED
main.go                  # UNCHANGED

README.md                # Updated `--watch` section + examples
```

**Structure Decision**: Single Go module, existing package layout. The CLI Surface owns the change per the architecture's logical and development views; no new packages and no movement of capabilities into `srv`. The `cmd/watch.go` + `cmd/watch_test.go` split is preferred to keep `cmd/root.go` from growing further; if the diff stays small, keeping everything in `cmd/root.go` + `cmd/root_test.go` is acceptable. Either choice keeps the dependency rules in `architecture-development-view.md` intact.

## Phase 0: Outline & Research

All technical-context items above resolved without `NEEDS CLARIFICATION`. Behavior was clarified during `/speckit-clarify` (Q1, Q2, Q3). Phase 0 still produces a concise `research.md` covering:
- Cobra `NoOptDefVal` semantics for a single flag that is both bare (`--watch`) and value-carrying (`--watch 5`), and the consequences for `-w`.
- Sleep-after-refresh vs ticker patterns in Go for CLI watch loops (justifies the chosen pattern against current ticker code).
- Signal handling in CLI watch loops with Cobra (`signal.NotifyContext` with `os.Interrupt` and `syscall.SIGTERM` on non-Windows; behavior matrix for Windows).
- Injectable-clock pattern for testing sleep-based loops without making tests slow.

**Output**: `specs/073-watch-interval-retry/research.md`.

## Phase 1: Design & Contracts

**Prerequisites**: `research.md` complete.

1. **Entities → `data-model.md`**: capture three small in-process entities — `WatchInterval`, `Refresh`, `WatchLoopState` — with fields, validation rules from spec FRs, and the state lifecycle from the logical view (Not-started → Active / Aborted-before-loop; Refresh: Succeeded / Failed-transient). No persistence.
2. **Interface contract → `contracts/cli.md`**: document the updated `--watch` / `-w` CLI contract, including: flag type and default value, validation rule, accepted invocation forms, exit-code matrix (one-shot, first-refresh failure, validation failure, clean interrupt), and the error-frame format (`time.RFC3339` timestamp + decorated error string).
3. **Quickstart → `quickstart.md`**: runnable validation steps an operator/reviewer can follow (`make build`, run with each invocation form including invalid inputs, simulate transient failure via a fake/missing cluster between refreshes, verify Ctrl-C stops cleanly).
4. **Agent context update**: refresh the `<!-- SPECKIT START -->` … `<!-- SPECKIT END -->` block in `AGENTS.md` to point at this plan (handled by the `speckit.agent-context.update` hook if it fires; otherwise via the agent-context extension).

**Output**: `data-model.md`, `contracts/cli.md`, `quickstart.md`, updated `AGENTS.md` block.

**Re-check Constitution gates** after Phase 1: still pass — design adds no parallel structure, keeps cluster I/O out of scope, and keeps user-facing changes documented.

## Complexity Tracking

> No Constitution Check violations. Section intentionally empty.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| (none) | (n/a) | (n/a) |
