# Data Model: Clearer EOF error guidance

## Decorated Retrieval Error

- **Purpose**: Represents a Kubernetes retrieval failure that has been classified into a more specific domain error before the root command renders a user-facing message.
- **Fields**:
  - `originalError`: the underlying Go error returned from the Kubernetes client
  - `resourceScope`: the retrieval path involved, such as nodes or pods
  - `isScopedEOF`: whether the error matches the cluster retrieval API `Get` EOF pattern covered by this feature
- **Validation rules**:
  - `originalError` must remain available for `Error()` output and `errors.As` / `errors.Is` behavior
  - `isScopedEOF` must be false for non-EOF failures and for EOFs outside the targeted retrieval path
  - classification must not remove existing unauthorized or forbidden decoration behavior
- **Lifecycle**:
  - created in `srv/api_errors.go` after a Kubernetes client call fails
  - passed upward through the current service and root command flow
  - consumed by the root fatal message selection logic

## Proxy Guidance Message

- **Purpose**: Captures the user-facing fatal message for scoped EOF retrieval failures.
- **Fields**:
  - `resourceName`: the resource currently being loaded and filtered
  - `proxyHint`: the guidance text explaining that proxy configuration may be the cause
  - `originalDetail`: the original low-level error detail, including the EOF context
- **Validation rules**:
  - `proxyHint` must be phrased as a possible cause, not a certainty
  - `originalDetail` must remain visible alongside the guidance
  - the combined message must stay deterministic across repeated occurrences
- **Lifecycle**:
  - assembled in `cmd/root.go` when the decorated scoped EOF error is encountered
  - emitted through the existing fatal logging path
  - repeated for every matching failure occurrence

## Non-Scoped Failure Outcome

- **Purpose**: Represents failures that should continue through the existing CLI behavior without the new proxy guidance.
- **Fields**:
  - `originalError`: the underlying or already-decorated error
  - `reasonCategory`: unauthorized, forbidden, metrics missing, generic failure, or out-of-scope EOF
  - `messagePath`: the existing root command branch used to render the result
- **Validation rules**:
  - non-EOF retrieval failures must keep their current message semantics
  - out-of-scope EOF failures must not accidentally receive the new proxy hint
  - existing warning and fatal branches must remain reachable
- **Lifecycle**:
  - classified before the new scoped EOF branch is checked
  - rendered via the current warning or fatal message paths
  - preserved by regression tests as the compatibility baseline
