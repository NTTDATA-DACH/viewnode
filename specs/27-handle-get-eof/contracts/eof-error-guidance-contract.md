# CLI Contract: Scoped EOF error guidance

## Purpose

Define the user-visible behavior for GitHub issue `#27`: clearer guidance when a cluster retrieval API `Get` failure surfaces as EOF.

## Command Contract

### `viewnode`

- Command name and flags remain unchanged
- Existing node and pod retrieval flow remains unchanged apart from the scoped EOF message refinement
- When node or pod retrieval fails with the targeted Kubernetes API `Get` EOF pattern:
  - the command exits as a failure through the existing fatal path
  - the emitted message keeps the current resource context such as `loading and filtering of nodes failed`
  - the emitted message states that proxy configuration may be the cause
  - the emitted message still includes the original low-level error details containing EOF
- When the failure is not EOF:
  - existing unauthorized, forbidden, metrics, and generic failure behavior remains unchanged
- When an EOF occurs outside the scoped retrieval pattern:
  - the command does not add the proxy guidance from this feature
- Repeated scoped EOF failures during the same run or session continue to emit the same guidance each time

## Example Message Fragments

Scoped EOF failure:

```text
loading and filtering of nodes failed
proxy configuration may be the cause
Get "...": EOF
```

Compatibility expectation for non-scoped failures:

```text
loading and filtering of nodes failed due to: <original error>
```

## Non-Goals

- No new command, flag, or alias
- No retry, backoff, or connectivity probing workflow
- No proxy auto-detection or proxy configuration editing
- No change to non-EOF failure wording unless required to preserve compatibility

## Verification Contract

- Tests verify only the scoped cluster retrieval API `Get` EOF pattern receives the proxy hint
- Tests verify non-EOF retrieval failures keep their existing behavior
- Tests verify out-of-scope EOF failures do not receive the proxy hint
- Tests verify the final user-facing message still includes the original low-level error details
- `make test` passes after the change
