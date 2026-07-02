# Feature Specification: Optional Refresh Interval and Retryable Watch Loop

**Feature Branch**: `073-watch-interval-retry`

**Created**: 2026-06-26

**Status**: Draft

**Input**: GitHub Issue #73 — fix(watch): support optional refresh interval and retryable watch loop

**GitHub Issue**:

- Number: 73
- URL: https://github.com/NTTDATA-DACH/viewnode/issues/73
- Title: fix(watch): support optional refresh interval and retryable watch loop

## Clarifications

### Session 2026-06-26

- Q: When the very first refresh in a `viewnode --watch` session fails, what should happen? → A: Mirror one-shot — exit immediately with the same non-zero exit code that `viewnode` without `--watch` would produce for the same error, and never enter the watch loop. Transient-error tolerance and on-screen error rendering only apply to refreshes that occur after the first refresh succeeds.
- Q: Which classification rule should the watch loop use to decide whether a mid-session refresh error is transient (retry) vs fatal (exit)? → A: All refresh errors that occur after the first successful refresh are treated as transient. The error is rendered in the current refresh frame and the loop retries on the next interval. The only way to leave the loop is Ctrl-C or an equivalent terminal interruption signal.
- Q: When a refresh fails mid-session, what should the watch frame for that interval show? → A: Replace the frame entirely with a clearly-labeled error message that includes the refresh timestamp and the underlying error text. The next successful refresh replaces the error frame with normal output. No prior-output overlay and no rolling error log are required.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Configurable refresh interval for `--watch` (Priority: P1)

An operator monitoring a Kubernetes cluster wants to control how frequently `viewnode` re-renders its output, so that a busy cluster is not polled aggressively while a quiet cluster is still observed at a useful cadence. They want to keep using the existing single `--watch` / `-w` flag, optionally supplying a number of seconds.

**Why this priority**: This is the primary user-visible change requested by the issue and unblocks both calmer polling and faster polling depending on cluster needs. Without it, the watch flag remains a fixed one-second refresh.

**Independent Test**: Run `viewnode` with and without the `--watch` / `-w` flag, with and without an explicit interval value (e.g., `--watch`, `-w 5`), and confirm that one-shot mode, default-interval watch, and custom-interval watch all behave as documented.

**Acceptance Scenarios**:

1. **Given** the operator runs `viewnode` without `--watch`, **When** the command completes, **Then** the output is rendered exactly once and the process exits as it does today.
2. **Given** the operator runs `viewnode --watch` (or `-w`) with no value, **When** the command runs, **Then** the output refreshes every 1 second until interrupted.
3. **Given** the operator runs `viewnode --watch 5` (or `-w 5`), **When** the command runs, **Then** the output refreshes approximately every 5 seconds until interrupted.
4. **Given** the operator runs `viewnode --watch 0`, `viewnode -w -3`, or `viewnode --watch abc`, **When** the command starts, **Then** the command fails fast with a clear usage error and a non-zero exit code, and does not enter watch mode.

---

### User Story 2 - Resilient watch loop that tolerates transient cluster errors (Priority: P1)

An operator running `viewnode --watch` against a real Kubernetes cluster wants the watch session to survive transient API errors (timeouts, throttling, brief connectivity loss, missing metrics endpoint, momentary RBAC blips) so they can leave the watch open without having to restart it after every hiccup.

**Why this priority**: This is the second core requirement of the issue. Today the watch terminates on many recoverable errors, which makes it unsuitable for live monitoring.

**Independent Test**: Simulate a transient Kubernetes failure between refreshes (for example, by temporarily removing access to the cluster or causing a controlled API error) and confirm the watch surfaces the error in the current refresh and recovers on the next interval instead of exiting.

**Acceptance Scenarios**:

1. **Given** a watch session is active and the next refresh encounters a transient Kubernetes/API error, **When** the refresh fails, **Then** the watch frame for that refresh is replaced with a clearly-labeled error message (refresh timestamp and underlying error text), the loop continues, attempts another refresh after the configured interval, and the next successful refresh replaces the error frame with normal output.
2. **Given** a refresh takes longer than the configured interval, **When** the refresh finally completes, **Then** the loop sleeps for the configured interval before the next refresh starts and does not trigger an immediate catch-up storm of refreshes.
3. **Given** a watch session is active, **When** the operator presses Ctrl-C, **Then** the loop stops cleanly, leaving the terminal in a usable state.

---

### User Story 3 - Watch reflects current command and configuration state on each refresh (Priority: P2)

An operator changes external state between refreshes (for example, switches the active kubeconfig context, namespace, or label selection) and expects subsequent refreshes to reflect that state where practical, so that watch mode behaves like repeated command execution rather than a frozen snapshot of process-level state.

**Why this priority**: Without this, the watch behaves inconsistently with a manual re-run of the command. It is important for correctness but secondary to the existence of a configurable interval and a resilient loop.

**Independent Test**: Start a watch session, change command-relevant external configuration that the underlying command would normally re-read on a fresh invocation, wait for the next refresh, and confirm the refreshed output reflects the change.

**Acceptance Scenarios**:

1. **Given** a watch session is running, **When** the next interval starts a new refresh, **Then** the refresh re-evaluates command and configuration state in the same way a fresh single-run invocation would, where practical.
2. **Given** non-watch invocations of the affected commands, **When** the operator runs them, **Then** their output and behavior remain unchanged from today.

---

### Edge Cases

- The operator passes a non-integer or out-of-range interval (`0`, negative, decimal, alphabetic, empty string after the flag separator): the command MUST fail before entering watch mode with a clear usage error.
- The first refresh fails immediately on startup (for example, the kubeconfig is invalid or the cluster is unreachable): the command MUST exit immediately with the same non-zero exit code that a non-watch `viewnode` invocation would produce for that same failure, and MUST NOT enter the watch loop. Transient-error retry behavior applies only to refreshes that occur after at least one refresh has succeeded.
- The cluster, kubeconfig, current context, or selected namespace becomes unavailable mid-session: the watch MUST render the error and continue retrying on the configured interval rather than terminating.
- A refresh takes longer than the configured interval: the loop MUST NOT queue or burst additional refreshes; it MUST sleep the configured interval after the slow refresh completes.
- Ctrl-C is pressed while a refresh is mid-flight: the loop MUST exit promptly and cleanly without leaving the terminal in a broken state.
- Output is redirected to a non-TTY (pipe, file): existing one-shot behavior MUST remain unchanged; watch behavior in non-TTY targets MUST follow whatever the current implementation already does and MUST NOT regress.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The CLI MUST keep a single, global `--watch` flag with the short alias `-w`. The flag MUST remain global so that all subcommands continue to expose it as they do today.
- **FR-002**: The `--watch` / `-w` flag MUST accept an optional positive integer number of seconds. Without an explicit value, the flag MUST default to 1 second.
- **FR-003**: When `--watch` / `-w` is not provided, the CLI MUST run exactly one refresh and exit, preserving existing one-shot output and exit behavior.
- **FR-004**: When `--watch` / `-w` is provided with a value, the CLI MUST refresh approximately every N seconds, where N is the user-supplied positive integer.
- **FR-005**: The CLI MUST reject invalid interval values (zero, negative numbers, non-integer values, non-numeric input) before entering watch mode, exit with a non-zero code, and show a clear usage error.
- **FR-006**: The CLI MUST stop watch mode cleanly on Ctrl-C (SIGINT) and on equivalent terminal interruption signals, leaving the terminal usable.
- **FR-007**: The watch loop MUST sleep for the configured interval **after** each refresh completes; it MUST NOT use ticker catch-up behavior that bursts queued refreshes when a single refresh takes longer than the interval.
- **FR-008**: Once at least one refresh has succeeded, the watch loop MUST treat every subsequent refresh error as transient: the loop MUST replace the current refresh frame with a clearly-labeled error message (refresh timestamp and the underlying error text) and continue, retrying on the next interval. The next successful refresh MUST replace the error frame with normal output. The loop MUST NOT exit on a refresh error after it has entered the loop. The only ways to leave the loop are Ctrl-C / equivalent terminal interruption signals (per FR-006). If the very first refresh fails, the command MUST exit immediately with the same non-zero exit code that a non-watch `viewnode` invocation would produce for that failure, and MUST NOT enter the loop.
- **FR-009**: Each refresh MUST be driven by a single, reusable single-run path that is also used by non-watch invocations, so that one-shot output and per-refresh output remain consistent.
- **FR-010**: Each refresh SHOULD re-read command/config state in the same way a fresh single-run invocation would, where practical, so that watch behavior approximates repeated command execution.
- **FR-011**: Non-watch invocations of all existing commands MUST keep their current output, exit codes, and side effects.
- **FR-012**: User-facing help text, README, and any documented examples for `--watch` / `-w` MUST be updated to describe the optional seconds value, default of 1, validation behavior, and retry behavior in the same change.
- **FR-013**: Unit tests MUST cover watch interval parsing (valid default, valid explicit, and rejected invalid inputs) and watch loop behavior (sleep-after-refresh, transient-error retry, clean exit) without requiring a live cluster.

### Key Entities *(include if feature involves data)*

- **Watch interval**: A positive integer number of seconds that controls how long the watch loop sleeps between refresh attempts; absent from one-shot invocations, defaults to 1 when `--watch` is given without a value.
- **Refresh**: A single execution of the reusable single-run path that produces the same kind of output as a non-watch invocation of the command.
- **Watch loop**: The orchestration around repeated refreshes, responsible for interval sleeping, transient error tolerance, and clean termination on interrupt.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of documented `--watch` / `-w` invocation forms in the acceptance scenarios behave as specified across one-shot, default-interval, and custom-interval modes.
- **SC-002**: Once the watch loop is entered, no refresh error causes the loop to exit; only Ctrl-C or an equivalent terminal interruption signal stops it.
- **SC-003**: When a refresh takes longer than the configured interval, the next refresh starts no sooner than the configured interval after the previous refresh completed (no catch-up bursts).
- **SC-004**: Invalid interval inputs fail before any refresh runs and produce a usage error the operator can act on without reading source code.
- **SC-005**: Non-watch invocations of every affected command produce byte-identical user-facing output and identical exit codes compared to the pre-change behavior on the same inputs.

## Assumptions

- The current `--watch` / `-w` flag is a global flag on the root command and continues to be exposed by all subcommands; this feature preserves that placement.
- "Transient" Kubernetes/API errors are defined by Clarification Q2 (session 2026-06-26): once the watch loop is running (i.e., the first refresh succeeded), every subsequent refresh error is treated as transient and the loop continues. Errors at the very first refresh follow the non-watch failure behavior and exit the command (per FR-008 and Edge Cases).
- "Where practical" in FR-010 means re-reading state at the same boundary the single-run path already uses; this feature does not require pushing per-refresh re-reads into deeper layers that today are intentionally process-scoped.
- Watch behavior against non-TTY outputs (pipes, file redirection) is out of scope for behavior changes; the implementation must not regress whatever the current code does there.
- The implementation will use the repository's established Cobra command structure and `NoOptDefVal = "1"` to model the optional value, consistent with the issue's implementation notes; this is captured here as an assumption rather than a user-visible requirement.
