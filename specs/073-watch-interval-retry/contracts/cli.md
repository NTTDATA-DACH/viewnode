# CLI Contract: `viewnode --watch`

**Scope**: The global `--watch` / `-w` flag on `viewnode` (root command). This contract supersedes the prior boolean flag behavior.

## Flag definition

| Property | Value |
|----------|-------|
| Long name | `--watch` |
| Short name | `-w` |
| Type | integer (seconds) |
| Default when flag absent | one-shot (no watch mode; tool runs exactly one refresh and exits) |
| Default when flag present without value (`NoOptDefVal`) | `1` |
| Allowed values | positive integers (`>= 1`) |
| Disallowed values | `0`, negative integers, non-integer input (rejected with usage error and non-zero exit, before any refresh runs) |
| Scope | global on the root command (inherited by subcommands as today) |

## Accepted invocation forms

| Invocation | Effect |
|------------|--------|
| `viewnode` | One refresh, then exit. Unchanged from today. |
| `viewnode --watch` | Watch mode, refresh every 1 second. |
| `viewnode -w` | Watch mode, refresh every 1 second. |
| `viewnode --watch 3` | Watch mode, refresh every 3 seconds. |
| `viewnode -w 3` | Watch mode, refresh every 3 seconds. |
| `viewnode --watch=5` | Watch mode, refresh every 5 seconds. |
| `viewnode -w=5` | Watch mode, refresh every 5 seconds. |
| `viewnode --watch 0` / `-w -1` / `--watch abc` | Usage error, non-zero exit. |

## Exit code matrix

| Situation | Exit code |
|-----------|-----------|
| One-shot success | `0` (unchanged from today). |
| One-shot failure | non-zero (unchanged from today; comes from existing `cmd/config` / `srv` error paths). |
| Watch session: invalid interval value | non-zero (validation failure, never enters watch mode). |
| Watch session: first refresh fails | same non-zero exit code that `viewnode` without `--watch` would produce for that same failure (clarification Q1). |
| Watch session: stopped by signal (Ctrl-C / SIGINT / SIGTERM on Unix) | `0`. The operator interrupt is the sanctioned way to leave the loop. |

## Loop behavior

- **Refresh cadence**: sleep starts after a refresh completes; next refresh starts no sooner than the configured interval after the previous refresh completed. Slow refreshes do **not** cause catch-up bursts (FR-007 / SC-003).
- **First-refresh gate**: the first refresh runs through the same single-run path as one-shot mode. The loop enters its Active state only after the first refresh succeeds (clarification Q1).
- **Post-first-success failure**: every refresh error after the first successful refresh is treated as transient. The loop renders an error frame and sleeps the configured interval before retrying. The loop never exits on a refresh error (clarification Q2).
- **Signal handling**: the loop terminates cleanly on `os.Interrupt` (Ctrl-C / SIGINT) on all platforms, and on `SIGTERM` on Unix. The current refresh is allowed to finish or is cancelled via context, and the active sleep returns immediately when the signal arrives (FR-006).

## Error-frame format

When a mid-session refresh fails, the watch frame for that refresh is replaced with:

```
[<RFC3339 timestamp>] watch refresh failed: <decorated error text>
```

- `<RFC3339 timestamp>` is the refresh's `startedAt`.
- `<decorated error text>` is the error returned from the single-run path (decorated by `srv.DecorateError` where applicable).
- The next successful refresh replaces this frame with normal output.

## Non-changes (explicit)

- Non-watch invocations of every existing command produce byte-identical user-facing output and identical exit codes compared to the pre-change behavior on the same inputs (SC-005).
- The flag stays on the root command. It does **not** become per-subcommand.
- No new flags introduced.
