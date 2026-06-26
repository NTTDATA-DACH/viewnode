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
- Fatal one-shot output includes ANSI color codes and a fresh timestamp prefix, so SC-005 diffs on failing paths need ANSI/timestamp normalization to isolate semantic regressions.
- Cobra treats `-w -1` as a shorthand-flag parse failure (`unknown shorthand flag: '1' in -1`) rather than as a negative integer value for `-w`, so shorthand-negative coverage should expect a parser-level usage error unless the operator uses `-w=-1`.
- A temporary `KUBECONFIG` copy is enough to simulate mid-session config loss and recovery against a real cluster because each watch refresh re-runs setup initialization instead of reusing stale clients.

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

## Iteration 4 - 2026-06-26 12:06 CEST
**User Story**: Partial progress on Phase 6 - Polish & Cross-Cutting Concerns
**Tasks Completed**:
- [x] T050: Update `README.md` for the optional watch interval, validation rules, retry behavior, and operator examples
- [x] T051: Verify `./viewnode --help` matches the watch flag contract
- [x] T052: Run `make test` with race detection
**Tasks Remaining in Story**: 3
**Commit**: No commit - partial progress
**Files Changed**:
- README.md
- specs/073-watch-interval-retry/tasks.md
- specs/073-watch-interval-retry/progress.md
**Learnings**:
- None of the configured kube contexts were reachable in this environment (`kind-*` endpoints refused connections and the EKS context failed DNS lookup), so the real-cluster quickstart scenarios remain blocked for a later iteration.
- Rebuilding the pre-feature binary from Git history (`088023b`) showed the expected `--help` diff is confined to the documented watch text; the failing one-shot path kept the same exit code and fatal payload, with only the existing log timestamp varying between runs.
---

## Iteration 5 - 2026-06-26 12:09 CEST
**User Story**: Partial progress on Phase 6 - Polish & Cross-Cutting Concerns
**Tasks Completed**:
- [x] T054: Diff `./viewnode --help` and a representative one-shot invocation against the pre-feature baseline
**Tasks Remaining in Story**: 2
**Commit**: No commit - partial progress
**Files Changed**:
- specs/073-watch-interval-retry/tasks.md
- specs/073-watch-interval-retry/progress.md
**Learnings**:
- A detached `git worktree` at `088023b` is the safest way to rebuild the pre-feature baseline for this repository; building from `mktemp`-style directories caused Go 1.25 to ignore `go.mod` as a temporary-root module.
- The current `./viewnode --help` diff versus `088023b` is limited to the intended `--watch` documentation, and the default one-shot failure path still matches the baseline exit code and fatal payload after stripping the pre-existing ANSI/timestamp prefix.
---

## Iteration 6 - 2026-06-26 12:10 CEST
**User Story**: Partial progress on Phase 6 - Polish & Cross-Cutting Concerns
**Tasks Completed**:
- None - T053 remains blocked because no configured kube context is reachable in this environment
**Tasks Remaining in Story**: 2
**Commit**: No commit - partial progress
**Files Changed**:
- specs/073-watch-interval-retry/progress.md
**Learnings**:
- All configured contexts still fail before Scenario B/F/H can start: the EKS endpoint fails DNS lookup and the local `kind-*` endpoints refuse connections.
- `kind` is installed locally, but Docker is not running, so a disposable replacement cluster could not be created to finish the real-cluster quickstart pass.
---

## Iteration 7 - 2026-06-26 12:11 CEST
**User Story**: Partial progress on Phase 6 - Polish & Cross-Cutting Concerns
**Tasks Completed**:
- None - T053 still cannot run end-to-end against a real cluster, so T055 remains blocked behind the unfinished work unit
**Tasks Remaining in Story**: 2
**Commit**: No commit - partial progress
**Files Changed**:
- specs/073-watch-interval-retry/progress.md
**Learnings**:
- `kubectl cluster-info --request-timeout=5s` against the current `kind-camunda-platform-local-c8.9` context still fails with `dial tcp 127.0.0.1:62403: connect: connection refused`, confirming the local kind control plane is down rather than just slow.
- `kind get clusters` currently fails because it cannot reach the Docker daemon at `unix:///Users/adam.boczek/.docker/run/docker.sock`, so there is no available local runtime to recreate the quickstart cluster in this environment.
---

## Iteration 8 - 2026-06-26 12:12 CEST
**User Story**: Partial progress on Phase 6 - Polish & Cross-Cutting Concerns
**Tasks Completed**:
- None - T053 still cannot run against a reachable cluster, so T055 remains blocked behind the unfinished real-cluster validation
**Tasks Remaining in Story**: 2
**Commit**: No commit - partial progress
**Files Changed**:
- specs/073-watch-interval-retry/progress.md
**Learnings**:
- Every configured kube context is still unreachable in this environment: the EKS endpoint fails DNS lookup and all three `kind-*` contexts fail with local `connect: connection refused` errors, so Scenarios B/C/F/H cannot be executed end-to-end.
- Git history already contains the earlier issue-linked Conventional Commits for Phases 2-5 (`61f08a3`, `bf90a2e`, `0984e5d`); the remaining Phase 6 commit should wait until T053 can complete so the final polish work can be recorded in one coherent commit.
---

## Iteration 9 - 2026-06-26 12:13 CEST
**User Story**: Partial progress on Phase 6 - Polish & Cross-Cutting Concerns
**Tasks Completed**:
- None - T053 remains incomplete because no configured kube context is reachable for the real-cluster scenarios, and T055 stays blocked behind that unfinished validation
**Tasks Remaining in Story**: 2
**Commit**: No commit - partial progress
**Files Changed**:
- specs/073-watch-interval-retry/progress.md
**Learnings**:
- `make build` still succeeds, and the cluster-free quickstart checks behave as expected: `./viewnode --watch 0` and `./viewnode --watch abc` fail fast with usage errors, while `KUBECONFIG=/nonexistent ./viewnode --watch` matches the one-shot fatal message and exit code for Scenario E.
- Scenario D has one Cobra-specific nuance: `./viewnode -w -1` exits non-zero before any refresh, but the error comes from Cobra's shorthand parser (`unknown shorthand flag: '1' in -1`) rather than the long-flag validator path.
- The real-cluster scenarios remain blocked for the same environmental reason as before: the EKS context still fails DNS lookup, all `kind-*` contexts still refuse local connections, and Docker is still unavailable to recreate a disposable kind cluster.
---

## Iteration 10 - 2026-06-26 12:14 CEST
**User Story**: Partial progress on Phase 6 - Polish & Cross-Cutting Concerns
**Tasks Completed**:
- None - T053 is still blocked by environment reachability, so T055 remains blocked behind the unfinished real-cluster validation
**Tasks Remaining in Story**: 2
**Commit**: No commit - partial progress
**Files Changed**:
- specs/073-watch-interval-retry/progress.md
**Learnings**:
- `make build` and `make test` still pass on the current tree, so the remaining blocker is the environment rather than the implementation.
- Rechecking all configured contexts shows the same hard failures as before: the EKS API hostname still fails DNS resolution, `kind-camunda-platform-local-c8.{7,8,9}` still refuse local TCP connections, and `kind get clusters` still cannot reach the Docker daemon, so Scenarios B/C/F/G/H cannot be executed against a real cluster in this iteration.
---

## Iteration 11 - 2026-06-26 12:17 CEST
**User Story**: Phase 6 - Polish & Cross-Cutting Concerns
**Tasks Completed**:
- [x] T053: Run the quickstart scenarios in `quickstart.md` against a real cluster
- [x] T055: Prepare the final Conventional Commit for the completed polish work unit
**Tasks Remaining in Story**: None - story complete
**Commit**: Recorded in Git history for this iteration
**Files Changed**:
- README.md
- specs/073-watch-interval-retry/tasks.md
- specs/073-watch-interval-retry/progress.md
**Learnings**:
- Starting Docker Desktop brought the existing `kind-camunda-platform-local-c8.9` cluster back, which was enough to finish Scenarios A/B/C/F/H without provisioning a new cluster.
- Real-cluster quickstart validation matched the contract: one-shot mode still exited once, `--watch` refreshed at roughly 1s, `-w 3` refreshed at roughly 3s, first-refresh failure matched the one-shot exit/error after timestamp normalization, and both `SIGINT` and `SIGTERM` exited cleanly with status 0.
- Scenario F confirmed the per-refresh reload contract end-to-end by removing and restoring a temporary `KUBECONFIG` copy mid-session; the watch rendered transient error frames, stayed alive, and recovered on the next successful refresh.
- Scenario G remained best-effort only: there was no safe latency injector available in this environment to force a slow-but-successful refresh, so the contract remains covered by the automated loop test added earlier.
---
