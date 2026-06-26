# Ralph Progress Log

Feature: 073-watch-interval-retry
Started: 2026-06-26 11:53:13

## Codebase Patterns

- `config.Initialize` should hand back a fresh `Setup` for each command execution so watch refreshes can pick up kubeconfig and client changes without reusing stale clientsets across iterations.
- Root command execution is easiest to test by stubbing package-level `runOnce` / `runWatch` hooks while still driving `RootCmd.Execute()` through Cobra's real flag parsing.
- The existing CLI keeps root behavior in `cmd/` with small helper files instead of new packages; build-tagged platform helpers are acceptable for signal differences.
- Cobra `NoOptDefVal` on the optional `--watch` int flag leaves space-separated values like `--watch 5` in `args`, so the root command needs a small normalization step before validation to preserve the documented CLI contract.
- When a watch refresh needs to replace the prior frame, emit the same clear-screen escape sequence that `srv.ViewNodeData.Printout` uses before writing the new frame to stdout.
- Production command execution can keep `log.Fatal` semantics while tests capture returned errors by overriding a small package-level error handler seam in `cmd/root.go`.

---

## Iteration 1 - 2026-06-26 11:59 CEST
**User Story**: US1 - Configurable refresh interval for `--watch`
**Tasks Completed**:
- [x] T001: Confirm baseline `make test` is green on `develop` and capture root help output into a scratch file
- [x] T010: Replace the boolean watch flag with `watchInterval` and derived `watchEnabled`
- [x] T011: Add watch interval validation
- [x] T012: Introduce signal-aware root context with Unix `SIGTERM` support
- [x] T013: Extract the single-refresh path into `runOnce`
- [x] T020: Move watch helpers into `cmd/watch.go` and add `runWatch` plus injectable sleep
- [x] T021: Rewire the root command to `runOnce` + sleep-after-refresh watch loop
- [x] T022: Update root help text for the optional seconds watch flag
- [x] T023: Add watch flag parsing and validation tests
- [x] T024: Add sleep-after-refresh loop invariant test
- [x] T025: Update existing root tests for the new watch state
**Tasks Remaining in Story**: None - story complete
**Commit**: Recorded in Git history for this iteration
**Files Changed**:
- cmd/root.go
- cmd/root_test.go
- cmd/watch.go
- cmd/watch_test.go
- cmd/signal_unix.go
- cmd/signal_windows.go
- specs/073-watch-interval-retry/tasks.md
- specs/073-watch-interval-retry/progress.md
**Learnings**:
- `runOnce` can safely own `config.Initialize(cmd)` without changing one-shot behavior, which also sets up the later per-refresh state re-read contract.
- Returning formatted errors from the load/filter path preserves the existing fatal output while making the single-run path reusable by the watch loop.
---

## Iteration 2 - 2026-06-26 12:00 CEST
**User Story**: US2 - Resilient watch loop that tolerates transient cluster errors
**Tasks Completed**:
- [x] T030: Update `runWatch` to keep retrying after post-start refresh failures and exit only on cancellation
- [x] T031: Add `renderRefreshError` to replace the active watch frame with a labeled refresh error block
- [x] T032: Add transient-error tolerance coverage for repeated watch refresh failures
- [x] T033: Add first-refresh-failure coverage for the root command handoff and watch-loop skip path
**Tasks Remaining in Story**: None - story complete
**Commit**: Recorded in Git history for this iteration
**Files Changed**:
- cmd/root.go
- cmd/root_test.go
- cmd/watch.go
- cmd/watch_test.go
- specs/073-watch-interval-retry/tasks.md
- specs/073-watch-interval-retry/progress.md
**Learnings**:
- `runWatch` should treat context cancellation separately from refresh failures so Ctrl-C does not render a spurious error frame on the way out.
- A tiny `handleRootCommandError` seam keeps production fatal behavior intact while letting Cobra-path tests assert the exact first-refresh handoff semantics.
---

## Iteration 3 - 2026-06-26 12:04 CEST
**User Story**: US3 - Watch reflects current command and configuration state on each refresh
**Tasks Completed**:
- [x] T040: Audit `runOnce` and reload a fresh setup per refresh so watch executions do not retain stale setup or API clients across iterations
- [x] T041: Add per-refresh re-evaluation coverage for mutable state observed between watch refreshes
**Tasks Remaining in Story**: None - story complete
**Commit**: Recorded in Git history for this iteration
**Files Changed**:
- cmd/config/config.go
- cmd/watch.go
- cmd/watch_test.go
- specs/073-watch-interval-retry/tasks.md
- specs/073-watch-interval-retry/progress.md
**Learnings**:
- Passing the freshly initialized setup through the single-run boundary makes the per-refresh reload contract explicit and keeps watch refreshes aligned with one-shot execution.
- The injected `runOnce` seam in `runWatch` is sufficient to model operator-visible state changes between refreshes without a live cluster.
---
