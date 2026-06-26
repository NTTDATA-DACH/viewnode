---
description: "Implementation tasks for issue #73 — Optional Refresh Interval and Retryable Watch Loop"
---

# Tasks: Optional Refresh Interval and Retryable Watch Loop

**Input**: Design documents from `/specs/073-watch-interval-retry/`

**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/cli.md](./contracts/cli.md), [quickstart.md](./quickstart.md)

**Issue**: <https://github.com/NTTDATA-DACH/viewnode/issues/73>

**Tests**: Required. This is a behavior change to a global CLI flag and its scheduler — every behavior-bearing change ships with a nearby Go test (constitution III).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story (US1, US2, US3 from spec.md)
- File paths are exact and relative to repo root

---

## Phase 1: Setup (Shared Infrastructure)

- [ ] T001 Confirm baseline `make test` is green on `develop` before starting and capture the current `./viewnode --help` output for the root command into a scratch file (working dir only, not committed) for SC-005 regression diff in T030. Note: no new dependency or tooling changes required by this feature.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Establish the shared seams used by all three user stories — interval type/validation, the single-run path, and the signal-driven context — without changing observable behavior yet.

- [ ] T010 Replace the global flag declaration in `cmd/root.go`: remove `var watchOn bool`; introduce `var watchInterval int` (seconds, validated separately) and a derived `var watchEnabled bool` set after flag parsing. Register the flag with `RootCmd.Flags().IntVarP(&watchInterval, "watch", "w", 0, "...")` and set `RootCmd.Flag("watch").NoOptDefVal = "1"` so `--watch` / `-w` default to 1 when present with no value. Keep the help text aligned with `contracts/cli.md` (mention seconds, default 1, `>= 1`).
- [ ] T011 [P] Add an interval validator function (e.g., `validateWatchInterval(present bool, seconds int) (bool, error)`) in `cmd/root.go` (or a new `cmd/watch.go` if extracted in T020) that:
  - returns `(false, nil)` when the flag is absent (one-shot mode);
  - returns `(true, nil)` when present with `seconds >= 1`;
  - returns a usage error when present with `seconds <= 0`.
  Cobra rejects non-integer input itself; document that in a code comment so reviewers don't add duplicate checks.
- [ ] T012 [P] Introduce `signal.NotifyContext`-based cancellation in `cmd/root.go`: build a `ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)` at the start of `RootCmd.Run`, register `SIGTERM` additionally on non-Windows via a small Unix-only helper, and `defer stop()`. Do NOT yet wire the loop to it — wiring happens in T021. Keep behavior unchanged for one-shot mode.
- [ ] T013 [P] Extract the existing per-refresh work (current `executeLoadAndFilter` + `executePrintOut` in `cmd/root.go`) into a single function `runOnce(ctx context.Context) error` that performs exactly one refresh end-to-end and returns the first error encountered. Preserve current logging and exit semantics for one-shot mode. This is the single-run path that one-shot mode and the watch loop will both call.

**Checkpoint**: After Phase 2, `viewnode` (no flag) and `viewnode --watch` still behave as today (the loop has not yet been replaced); `make test` is green; the flag is registered with its new type but only one-shot mode is fully wired.

---

## Phase 3: User Story 1 — Configurable refresh interval for `--watch` (Priority: P1) 🎯 MVP

**Goal**: `viewnode --watch [N]` / `viewnode -w [N]` runs the refresh loop on the configured interval; invalid values fail fast; non-watch invocations are unchanged.

**Independent Test**: Run `./viewnode`, `./viewnode --watch`, `./viewnode -w`, `./viewnode --watch 5`, `./viewnode -w 5`, and `./viewnode --watch 0` / `./viewnode --watch abc`. Confirm one-shot exits once; default-interval refreshes ~1s; explicit interval refreshes ~Ns; invalid input fails before any refresh.

- [ ] T020 [US1] Create `cmd/watch.go` (new file) and move the existing `schedule`/`handleErrors`/`executeLoadAndFilter`/`executePrintOut` helpers out of `cmd/root.go` into it, alongside a new `runWatch(ctx context.Context, interval time.Duration, runOnce func(context.Context) error, sleep func(context.Context, time.Duration) error) error` function. Keep all code in `package cmd`. The default production `sleep` is a small wrapper around `time.NewTimer` that returns immediately when `ctx.Done()` fires; expose it so tests can substitute.
- [ ] T021 [US1] Rewire `RootCmd.Run` in `cmd/root.go` to use the new wiring: call `validateWatchInterval` (T011) → run `runOnce(ctx)` from T013 → if not in watch mode (or first refresh failed), return its result and exit; otherwise call `runWatch(ctx, time.Duration(watchInterval)*time.Second, runOnce, productionSleep)`. Remove the old ticker-driven `schedule` goroutine and `sync.WaitGroup`. Map errors from `runOnce`/`runWatch` to existing fatal/exit behavior so one-shot exit codes remain byte-identical.
- [ ] T022 [US1] Update the root command long help and flag help in `cmd/root.go` to match `contracts/cli.md`: explain the optional integer seconds value, default `1`, `>= 1` requirement, and the fail-fast-on-invalid rule. Keep the help text concise and operator-readable; no implementation details.
- [ ] T023 [P] [US1] Add unit tests in `cmd/watch_test.go` for interval flag parsing and validation:
  - `--watch` absent → `watchEnabled == false`;
  - `--watch` present, no value → `watchInterval == 1`, `watchEnabled == true`;
  - `--watch 5` and `-w 5` → `watchInterval == 5`, `watchEnabled == true`;
  - `--watch 0` → validation error, non-zero exit;
  - `--watch -1` → validation error;
  - `--watch abc` → Cobra parse error (assert error type/message contains "invalid argument").
  Use Cobra's `RootCmd.SetArgs` + `Execute` pattern, mirroring the style already used in `cmd/root_test.go`.
- [ ] T024 [P] [US1] Add a unit test in `cmd/watch_test.go` for the sleep-after-refresh invariant (FR-007, SC-003): inject a `runOnce` that records call timestamps and a `sleep` that records requested durations; assert that for a scripted slow refresh (e.g., `runOnce` sleeps 50ms while interval is 10ms), the next `sleep` call requests the full interval after `runOnce` returned, and that two `runOnce` calls are never back-to-back without a `sleep` in between.
- [ ] T025 [US1] Update `cmd/root_test.go` for the new flag type: any test that previously set `watchOn = false` or asserted on the boolean flag now uses `watchInterval = 0` / `watchEnabled = false`, or resets state via the flag's `Value.Set`. Confirm `make test` stays green.

**Checkpoint**: US1 complete and independently demonstrable.

---

## Phase 4: User Story 2 — Resilient watch loop that tolerates transient cluster errors (Priority: P1)

**Goal**: A mid-session refresh failure renders an error frame for that refresh and the loop continues. Only Ctrl-C / SIGINT / SIGTERM stops the loop. The first refresh still fails fast per clarification Q1.

**Independent Test**: Start `./viewnode -w 2` against a reachable cluster; between refreshes, make the cluster unreachable; observe the error frame and the loop continuing; restore connectivity and observe normal frames; press Ctrl-C and observe exit 0.

- [ ] T030 [US2] In `cmd/watch.go`, implement the loop-state logic in `runWatch`:
  - Treat the **first** `runOnce` call as part of the foundational handoff in `RootCmd.Run` (already done in T021 — first refresh runs through the single-run path and its error returns to the caller before `runWatch` is entered). Document in a code comment that `runWatch` is entered only after the first refresh succeeded.
  - Inside `runWatch`, every `runOnce` error must be rendered as an error frame (via T031) and the loop must continue with `sleep(ctx, interval)`. The loop exits only when `ctx.Done()` fires.
- [ ] T031 [US2] Add an `renderRefreshError(start time.Time, err error)` helper in `cmd/watch.go` that writes a single block to stderr (or the same stream existing refresh output uses; mirror the current writer) in the exact format from `contracts/cli.md`:

  ```
  [<RFC3339 timestamp>] watch refresh failed: <decorated error text>
  ```

  Where `<decorated error text>` is `err.Error()` (errors from `srv` are already decorated by `srv.DecorateError`). Do not append a rolling log; this single block replaces the failing refresh's normal frame.
- [ ] T032 [P] [US2] Add a unit test in `cmd/watch_test.go` for transient-error tolerance (FR-008, SC-002):
  - Script `runOnce` to return a sequence: success, error("conn refused"), error("api timeout"), success, then cancel `ctx`.
  - Assert: `runWatch` returns `nil` (clean exit from cancellation); both errors were rendered via `renderRefreshError` to the captured writer; the loop made all four `runOnce` calls; `sleep` was called between each pair.
- [ ] T033 [P] [US2] Add a unit test in `cmd/watch_test.go` for the first-refresh-failure path (clarification Q1): drive `RootCmd.Run` with a stubbed `runOnce` that fails on its first call; assert the process path returns the underlying error to Cobra (so `log.Fatal` / Cobra's error reporting produce the same non-zero exit as one-shot), and `runWatch` was never entered. Use a small test seam (export `runOnce` via package-level variable for tests, or refactor `RootCmd.Run` into a thin wrapper around a testable function — pick whichever keeps the diff smallest).

**Checkpoint**: US1 + US2 deliver the complete monitoring story end-to-end.

---

## Phase 5: User Story 3 — Watch reflects current command and configuration state on each refresh (Priority: P2)

**Goal**: Each refresh re-reads command/config state at the same boundary the single-run path uses, so watch behaves like repeated invocations "where practical" (FR-010). Non-watch invocations remain unchanged.

**Independent Test**: Start a watch session, change a flag-relevant kubeconfig setting (e.g., switch current namespace via `viewnode ns set-current` in a second terminal), wait for the next refresh, confirm the refresh reflects the change (subject to anything the single-run path already re-reads).

- [ ] T040 [US3] Audit `runOnce` (T013) and confirm that everything that today happens inside the existing `executeLoadAndFilter` boundary on a fresh invocation runs again on every call. Concretely:
  - `config.Initialize(cmd)` runs once per `runOnce` (the same as today's single-run code path);
  - `getViewNodeDataConfig` is recomputed per call;
  - no module-level cached `Setup`/`Api` clients persist across calls.
  If any caching slipped in during T013/T020, move construction inside `runOnce`. Document the boundary in a short comment on `runOnce`.
- [ ] T041 [P] [US3] Add a unit test in `cmd/watch_test.go` that verifies the per-refresh re-evaluation contract: inject a `runOnce` that captures a per-call view of a mutable struct (e.g., a counter incremented from the test between calls) and assert that successive refreshes observe the latest mutation. This test enforces the architectural contract from `architecture-process-view.md` ("Process Gaps" → per-refresh re-load boundary).

**Checkpoint**: All three user stories complete. Behavior matches spec.md and all three clarification answers.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T050 [P] Update `README.md`: rewrite the `--watch` section to describe the optional integer seconds value, default 1, fail-fast validation, transient-error tolerance, and error-frame format. Update or add examples for `viewnode`, `viewnode --watch`, `viewnode --watch 5`, `viewnode -w 5`, and the invalid-input case. Keep the section operator-readable; no implementation language.
- [ ] T051 [P] Verify `./viewnode --help` reflects the new flag semantics and matches `contracts/cli.md`. If the help still says "executes the command every second", update it (the existing string was the source of operator confusion).
- [ ] T052 Run `make test` and confirm it passes with race detection on. Resolve any flake caused by the new tests (the loop tests must not depend on wall-clock duration — use the injectable `sleep`).
- [ ] T053 Run the quickstart scenarios in `quickstart.md` against a real cluster: Scenarios A (regression baseline), B, C, D, E, F, G (best-effort), H. Record any deviation; if Scenario F or H reveals a bug, return to Phase 4.
- [ ] T054 Diff `./viewnode --help` and `./viewnode <typical inspection invocation>` output against the scratch baseline captured in T001 to confirm SC-005 (one-shot byte-identical output) holds.
- [ ] T055 Prepare commit(s) following Conventional Commits with `#73` suffix per `AGENTS.md` and the gvb skill rule. Suggested grouping: one commit for Phase 2 + 3 (`fix(watch): support optional refresh interval and validate input #73`), one for Phase 4 (`fix(watch): tolerate transient refresh errors with labeled error frame #73`), one for Phase 5 + 6 (`docs(watch): document new --watch behavior and per-refresh state contract #73`). Final commits must include the standard `Co-authored-by: Copilot` trailer.

---

## Dependencies & Execution Order

- **Phase 1 (T001)** → **Phase 2 (T010–T013)** → **Phase 3 (T020–T025)** → **Phase 4 (T030–T033)** → **Phase 5 (T040–T041)** → **Phase 6 (T050–T055)**.
- Within Phase 2: T010 must land before T011 (validator uses the new flag variable); T012 and T013 can land in parallel with each other but both before Phase 3.
- Within Phase 3: T020 → T021 → T022 in sequence (each builds on the previous wiring); T023/T024 (tests) can run in parallel against the T020–T022 code; T025 only after T020/T021.
- Within Phase 4: T030 → T031 in sequence; T032/T033 (tests) in parallel after T031.
- Within Phase 5: T040 → T041.
- Within Phase 6: T050/T051 in parallel; T052 after all preceding code/test tasks; T053 after T052; T054 after T053; T055 last.

## Parallel Opportunities

- **Phase 2**: T011 + T012 + T013 (all touch new code, distinct concerns).
- **Phase 3**: T023 + T024 (independent tests in the same new file but non-overlapping test functions).
- **Phase 4**: T032 + T033 (independent test functions).
- **Phase 6**: T050 + T051 (README vs help string, separate files/strings).

## Implementation Strategy

- **MVP (deliverable on its own)**: Phase 1 + Phase 2 + Phase 3 (US1). At this point, `--watch [N]` with validation is shipped; the loop still inherits today's failure-exit behavior. Operators get the interval they asked for; the resilience win lands in the next increment.
- **Incremental delivery**: Land Phase 4 (US2) next to close the issue's resilience requirement. Phase 5 (US3) closes the "where practical" re-read contract. Phase 6 polishes docs and verifies SC-005.
- **Commit cadence**: One commit per phase boundary (or per logical group within a phase) keeps history reviewable per `AGENTS.md`'s small-commits rule; every subject ends with `#73`.

## Independent Test Criteria (per story)

- **US1**: `./viewnode`, `./viewnode --watch`, `./viewnode -w`, `./viewnode --watch 5`, `./viewnode -w 5`, `./viewnode --watch 0` / `--watch abc` all behave per the `contracts/cli.md` "Accepted invocation forms" table.
- **US2**: A watch session survives a simulated mid-session cluster failure (Scenario F in `quickstart.md`); only Ctrl-C / SIGINT / SIGTERM exit the loop.
- **US3**: A change to command/config state made between refreshes is reflected on the next refresh, to the same extent a fresh single-run invocation would have reflected it.

## Validation Summary

- ✅ All tasks use checklist format `- [ ] [TaskID] [P?] [Story?] Description with file path`.
- ✅ Each user story is independently testable and has nearby Go test tasks (constitution III).
- ✅ No parallel command/configuration/service structure introduced (constitution IV); change is localized to `cmd/`.
- ✅ Documentation parity tasks present (T050, T051) for the CLI surface change (constitution V).
- ✅ Issue traceability preserved: every commit subject ends with `#73` per T055 and `AGENTS.md`.
