# Quickstart: Validate `viewnode --watch`

This guide is the post-implementation validation runbook for issue #73. It assumes the implementation has been built and the operator has a reachable Kubernetes cluster and a working `kubeconfig`.

## Prerequisites

- Repository cloned and on branch `073-watch-interval-retry`.
- Go 1.25.x available (matches `go.mod`).
- A reachable Kubernetes cluster the operator can list nodes/pods in.
- A second terminal (for transient-failure simulation).

## Build

```sh
make build
./viewnode --help | sed -n '/--watch/,+2p'
```

Expected: `--watch` shown with `int` type (seconds), short alias `-w`, default `0`/absent meaning one-shot.

## Test (unit, no cluster required)

```sh
make test
```

Expected: all tests pass. New tests cover flag parsing, single-run path, sleep-after-refresh, transient-error tolerance, and clean stop.

## Scenario A — One-shot behavior unchanged (SC-005)

```sh
./viewnode > /tmp/viewnode.before.out 2>&1; echo "exit=$?"
```

Expected: identical output and exit code to the pre-change `viewnode` build on the same cluster. (Diff against a known-good baseline if available.)

## Scenario B — Default interval (FR-002, US1 #2)

```sh
./viewnode --watch
# Wait ~10 seconds, observe ~10 refreshes, then press Ctrl-C
echo "exit=$?"
```

Expected: refresh roughly every 1 second; Ctrl-C ends cleanly with exit 0.

## Scenario C — Explicit interval (FR-002, US1 #3)

```sh
./viewnode -w 3
# Wait ~15 seconds, observe ~5 refreshes, then press Ctrl-C
```

Expected: refresh roughly every 3 seconds.

## Scenario D — Validation (FR-005, US1 #4)

```sh
./viewnode --watch 0; echo "exit=$?"
./viewnode -w -1; echo "exit=$?"
./viewnode --watch abc; echo "exit=$?"
```

Expected for each: usage error on stderr, non-zero exit, no refresh attempted.

## Scenario E — First-refresh failure (clarification Q1)

```sh
KUBECONFIG=/nonexistent ./viewnode --watch; echo "exit=$?"
```

Expected: same non-zero exit code and error message as `KUBECONFIG=/nonexistent ./viewnode` (without `--watch`). No watch frames rendered.

## Scenario F — Mid-session transient failure tolerance (US2, clarifications Q2+Q3)

1. Start a watch session against a reachable cluster: `./viewnode -w 2`.
2. In a second terminal, temporarily make the cluster unreachable (e.g., move/rename your kubeconfig, disable VPN, point to a stale context). Keep `viewnode` running.
3. Observe: the failing refresh frame is replaced with `[<RFC3339 timestamp>] watch refresh failed: <error>`. The loop **does not exit**.
4. Restore connectivity. Observe the next refresh replaces the error frame with normal output.
5. Press Ctrl-C.

Expected: the watch session survives all transient failures; only Ctrl-C exits the loop (exit 0).

## Scenario G — Slow refresh, no catch-up (FR-007, SC-003)

If a controlled way to introduce per-refresh latency exists (e.g., a busy cluster or a fault-injection proxy), set `-w 1` and confirm that when a refresh takes ~3 seconds, the next refresh starts ~1 second after the slow one completes, not immediately. Two refreshes never run back-to-back without an interval sleep between them.

## Scenario H — Clean signal stop (FR-006)

On Unix:

```sh
./viewnode -w 2 &
sleep 3
kill -INT $!     # then repeat with SIGTERM:
./viewnode -w 2 &
sleep 3
kill -TERM $!
```

Expected for each: process exits cleanly with exit 0, terminal returns to a usable state.

## References

- Spec: [`spec.md`](./spec.md)
- Plan: [`plan.md`](./plan.md)
- CLI contract: [`contracts/cli.md`](./contracts/cli.md)
- Data model: [`data-model.md`](./data-model.md)
- Issue: <https://github.com/NTTDATA-DACH/viewnode/issues/73>
