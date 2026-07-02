# Phase 1 Data Model: Watch Interval and Retryable Loop

This feature introduces no persistent or wire-format data. The following are in-process logical entities owned by the CLI Surface / Watch Orchestration sub-capability (`cmd/`).

## Entity: WatchInterval

**Purpose**: Parsed, validated cadence between refreshes in a watch session.

**Fields**:
- `seconds` (int, positive): number of seconds to sleep after each completed refresh.

**Source**: CLI flag `--watch` / `-w` (Cobra integer flag with `NoOptDefVal = "1"`).

**Validation rules** (from spec FR-002, FR-005):
- Absent value → `seconds = 1`.
- `seconds >= 1` required.
- Non-integer input rejected by Cobra's flag parser with a usage error.
- `seconds <= 0` rejected by an explicit pre-execution check that returns a usage error and a non-zero exit code.

**Lifetime**: Created once per process invocation. Never mutated after validation.

## Entity: Refresh

**Purpose**: A single end-to-end "obtain cluster state and render it" operation; identical work whether invoked one-shot or as a watch tick.

**Fields**:
- `startedAt` (`time.Time`): captured immediately before the refresh runs; used in error frames.
- `outcome` (`enum: Succeeded | Failed`): set when the refresh returns.
- `err` (`error`, set only when `outcome == Failed`): decorated by `srv.DecorateError` where applicable; rendered verbatim in error frames.

**Source**: Constructed by the single-run path in `cmd/` for each invocation/tick.

**Validation rules**:
- A refresh runs to completion (success or failure) before the next refresh may be considered (one Refresh at a time within a watch loop — process view invariant).

**Lifetime**: Discarded after the result is rendered (no per-refresh state survives the cycle, satisfying the logical-view non-goal "no cross-refresh caching").

## Entity: WatchLoopState

**Purpose**: Tracks loop lifecycle for the Watch Orchestration sub-capability.

**Fields**:
- `state` (`enum: NotStarted | Active | AbortedBeforeLoop | Stopped`).
- `firstRefresh` (`*Refresh`): the first refresh, retained only long enough to decide whether to enter `Active`.
- `cancelReason` (`enum: Signal | ContextCanceled`, set only when `state == Stopped`).

**Source**: Created when a watch invocation begins; lives for the duration of the process.

**State transitions** (from logical view):
- `NotStarted` → `AbortedBeforeLoop` when `firstRefresh.outcome == Failed`. Process exits with the same non-zero code the one-shot path would produce. **Forbidden**: transitioning to `Active` without a successful first refresh.
- `NotStarted` → `Active` when `firstRefresh.outcome == Succeeded`.
- `Active` → `Stopped` when the cancellation context is done (operator interrupt or equivalent signal). **Forbidden**: `Active` → `Stopped` due to a refresh error (any refresh error inside `Active` is rendered as an error frame and the loop continues).

## Relationships

- One process invocation → one `WatchInterval` (only when `--watch` is given) and exactly one `WatchLoopState`.
- One `WatchLoopState` → many `Refresh` instances (≥1 in `Active`; exactly 1 in `AbortedBeforeLoop`; 1 in one-shot mode, which uses the single-run path without a loop state).
- No `Refresh` instance crosses processes.

## Invariants

- A `Refresh` failure inside an `Active` `WatchLoopState` never advances the state to anything other than another `Refresh` cycle. (Spec FR-008; clarification Q2.)
- The first `Refresh` of any watch invocation goes through the same single-run path that one-shot mode uses. (Spec FR-009; clarification Q1.)
- `WatchInterval.seconds` is read once and used as the sleep duration after every `Refresh` completes — never used as a deadline that could fire mid-refresh. (Spec FR-007; SC-003.)
