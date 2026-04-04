# Research: Remove redundant namespace prefixes from single-namespace pod rows

## Decision 1: Keep the fix in the flat render path

- **Decision**: Correct the output in `srv/view.go` by changing the flat pod-row formatting used when a single namespace is selected.
- **Rationale**: The bug is purely presentational. The namespace selection is already parsed correctly and the header is already rendered, so the smallest repository-native fix is to adjust how the pod row is printed.
- **Alternatives considered**:
  - Change `cmd/root.go` configuration to disable namespace-aware rendering entirely for single-namespace runs: rejected because the header still needs the namespace context.
  - Alter pod-loading data so namespaces are stripped earlier: rejected because the renderer still needs namespace information for other modes.

## Decision 2: Keep the namespace header as the canonical single-namespace label

- **Decision**: Preserve the existing `namespace(s): <name>` header for explicitly scoped single-namespace runs and remove only the redundant inline `<namespace>: ` pod-row prefix.
- **Rationale**: The header already communicates the selected namespace once at the appropriate scope. Removing the per-row prefix cleans up the output without losing information.
- **Alternatives considered**:
  - Remove the namespace header too: rejected because the issue asks only to remove redundant row prefixes, not to redesign scoped output.
  - Keep the prefix for consistency with the old flat path: rejected because the old behavior is the reported defect.

## Decision 3: Preserve other display modes and pod-row details unchanged

- **Decision**: Leave grouped multi-namespace and all-namespaces rendering unchanged, and preserve container, timing, and metrics details on single-namespace pod rows.
- **Rationale**: The issue is narrowly scoped to the redundant single-namespace prefix. Preserving adjacent behavior keeps the fix small and lowers regression risk.
- **Alternatives considered**:
  - Harmonize all namespace display modes in one refactor: rejected because it broadens the scope beyond the bug report.
  - Drop or restructure pod-row details while touching the formatter: rejected because those details are not part of the defect.

## Decision 4: Update README examples in the same change

- **Decision**: Refresh single-namespace examples and explanation in `README.md` so documentation matches the corrected output.
- **Rationale**: The repository constitution requires user-facing documentation to change alongside user-visible behavior changes.
- **Alternatives considered**:
  - Defer documentation updates to a later cleanup: rejected because it would leave shipped behavior and documentation out of sync.
