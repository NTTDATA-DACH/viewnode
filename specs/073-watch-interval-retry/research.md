# Phase 0 Research: Watch Interval and Retryable Loop

All `NEEDS CLARIFICATION` items from `plan.md` Technical Context were resolved by spec clarifications Q1/Q2/Q3 and by inspecting the existing `cmd/root.go`. This file consolidates the technical patterns the implementation will rely on.

## Decision 1 — Cobra flag shape: keep `--watch` / `-w`, optional integer value

**Decision**: Change the flag from `BoolVarP(&watchOn, "watch", "w", false, ...)` to an integer flag with `NoOptDefVal = "1"` so that both `--watch` (no value) and `--watch 5` (with value) are accepted, and `-w` / `-w 5` mirror them.

**Rationale**:
- Cobra's `NoOptDefVal` is the canonical, documented way to make an optional-value flag work for both bare and value-carrying forms.
- Preserves the existing one-letter alias (`-w`) per FR-001.
- Keeps the flag global on `RootCmd` (same `Flags()` registration site), matching today's CLI surface.

**Alternatives considered**:
- A separate `--watch-interval N` flag plus `--watch` bool. Rejected — doubles the surface, contradicts FR-001, and the issue explicitly asks for a single flag.
- A string flag parsed manually. Rejected — loses Cobra's built-in type validation and gives a worse error message for non-integer input.

## Decision 2 — Validation: positive integer seconds, default 1

**Decision**: Parse the value as `int`. Treat absent value as `1`. Reject `< 1` (zero, negative) and non-integer input with a clear usage error before entering watch mode.

**Rationale**:
- Matches FR-002 (default 1) and FR-005 (reject zero/negative/non-integer).
- Integer seconds is the simplest model that satisfies all acceptance scenarios; sub-second cadence is explicitly out of scope.

**Alternatives considered**:
- `time.Duration` (`--watch 5s`). Rejected — issue's acceptance scenarios use integer seconds and existing users do not pass duration strings today.
- `uint`. Rejected — Cobra prints `uint` errors that are less friendly than `int` + explicit range check.

## Decision 3 — Loop pattern: sleep-after-refresh, not ticker

**Decision**: Replace the current `time.NewTicker(1 * time.Second)` loop with `for { runRefresh(); if interrupted { break }; sleep(interval) }`. Use an injectable clock interface for tests so sleeps in unit tests are not real wall-clock waits.

**Rationale**:
- FR-007 + SC-003 explicitly require sleep-after-completion to prevent catch-up bursts on slow refreshes.
- Ticker-based code paths queue ticks while a refresh is running and can deliver them back-to-back when the refresh completes — exactly the bug class FR-007 forbids.
- Predictable bounded resource usage; trivially testable with a fake clock.

**Alternatives considered**:
- Keep the ticker, but drain pending ticks after a slow refresh. Rejected — more code, more failure modes, no benefit over straight sleep.
- `time.AfterFunc` chain. Rejected — harder to interrupt cleanly and equivalently testable.

## Decision 4 — Failure-closure: first-refresh fatal, subsequent universally transient

**Decision** (from clarification Q1+Q2): The first refresh runs through the same single-run path as one-shot mode. On error, the process exits with the same non-zero exit code one-shot would produce. Only if the first refresh succeeds does the loop enter the Active state. Once Active, every refresh error is treated as transient; the error is rendered in the current refresh's frame and the loop sleeps and tries again. The only way to leave the Active loop is a cancellation signal.

**Rationale**: Predictable startup (operators don't get a never-rendering loop hiding a bad kubeconfig) plus resilient steady state (operators don't have to babysit the watch). Documented in the architecture process view as `Watch Loop` lifecycle invariants.

**Alternatives considered**:
- Always enter the loop (option B during clarification). Rejected by clarification Q1.
- Classify errors into transient vs fatal categories (option B during Q2). Rejected — more complex, harder to test, and the issue's framing ("exits on many transient errors") points at universal transience post-first-success.

## Decision 5 — Error-frame rendering: replace the frame

**Decision** (from clarification Q3): When a mid-session refresh fails, replace the frame for that refresh with a clearly-labeled block:
- Line 1: refresh timestamp formatted with `time.RFC3339`.
- Line 2: the underlying error text (already decorated by `srv.DecorateError` where applicable).
- The next successful refresh replaces the error frame with normal output.

**Rationale**: Single rendering pattern (no overlay, no rolling log), reuses the existing single-frame model, and is easy to test by capturing stdout/stderr.

**Alternatives considered**:
- Keep last-good frame and overlay error (option B during Q3). Rejected by clarification.
- Append a rolling error log (option C during Q3). Rejected — out of scope and adds state across refreshes.

## Decision 6 — Signal handling: `signal.NotifyContext` with `os.Interrupt` (+ `SIGTERM` on Unix)

**Decision**: Wrap the loop in a `context.Context` produced by `signal.NotifyContext(ctx, os.Interrupt)` (and, on Unix, `syscall.SIGTERM`). The single-run path and the sleep step both check `ctx.Done()` so an in-flight refresh cancels cleanly and a sleep aborts immediately on signal.

**Rationale**:
- `signal.NotifyContext` is the standard Go 1.16+ idiom and is friendlier than manually orchestrated channels (current code uses `chan bool` + goroutine + `WaitGroup`, which is harder to reason about for cancellation).
- Satisfies FR-006 (clean stop on Ctrl-C) without introducing platform-specific code in the hot path.
- `SIGTERM` on Unix is the issue's "equivalent terminal interruption signals" — added behind a build-tag-free guard (`syscall.SIGTERM` is portable via `os/signal` on Unix; not registered on Windows).

**Alternatives considered**:
- Keep current stop-channel pattern, add signal handler. Rejected — duplicates `signal.NotifyContext` and keeps the harder-to-test pattern.
- Use a third-party signal package. Rejected — standard library suffices.

## Decision 7 — Test seam: injectable clock + injectable refresh function

**Decision**: Extract the loop into a small helper (e.g., `runWatch(ctx, interval, refresh func(context.Context) error, sleep func(context.Context, time.Duration) error)`). Tests pass a fake `sleep` that returns immediately and a `refresh` that returns scripted outcomes (success, transient error, sequence). Real production wiring uses `time.Sleep` via a tiny wrapper that respects context cancellation.

**Rationale**:
- Makes loop tests fast and deterministic without `time.Sleep` in CI.
- Keeps production code simple — the indirection is one function pointer, not an interface hierarchy.
- Confines the test seam to `cmd/` per the development view (no production change in `srv/` or `cmd/config`).

**Alternatives considered**:
- Test against real wall-clock with small intervals. Rejected — flaky under load.
- A `Clock` interface. Rejected — overweight for two test seams.
