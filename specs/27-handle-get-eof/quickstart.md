# Quickstart: Clearer EOF error guidance

## Goal

Implement and validate clearer proxy-related guidance for scoped Kubernetes API `Get` EOF failures in the root `viewnode` command.

## Implementation Steps

1. Update `srv/api_errors.go` to classify scoped retrieval EOF failures into a dedicated decorated error type while preserving the original error details.
2. Update `cmd/root.go` to translate that decorated error into a clearer fatal message that keeps the original low-level detail and adds the proxy hint.
3. Add regression tests for scoped EOF decoration and for the root fatal-message behavior.

## Validation Steps

1. Run targeted tests while iterating:

   ```bash
   go test ./srv ./cmd
   ```

2. Run the full repository validation before considering the change complete:

   ```bash
   make test
   ```

3. Optional manual verification against a cluster environment with a broken proxy path:

   ```bash
   viewnode --kubeconfig /path/to/test-kubeconfig
   ```

   Expected outcome: a fatal message for the scoped EOF path that still shows the original EOF detail and adds a hint that proxy configuration may be the cause.
