# Research: Fix empty namespace groups in namespace-filtered node output

## Decision 1: Keep the fix in the printout layer

- **Decision**: Correct namespace-group visibility inside `srv/view.go` using each node's already loaded `[]ViewPod` data.
- **Rationale**: The bug is about incorrect rendering, not incorrect pod discovery. Fixing it in the printout layer keeps the change small and consistent with the repository's existing view-layer structure.
- **Alternatives considered**:
  - Change pod loading or filtering to drop namespace selections earlier: rejected because it broadens a renderer defect into the data pipeline.
  - Add node-level flags during loading to mark empty namespaces: rejected because the renderer already has the data it needs.

## Decision 2: Preserve selected namespace input without forcing empty groups

- **Decision**: Continue passing `SelectedNamespaces` from `cmd/root.go`, but use it only to scope relevant namespace groups instead of rendering a heading for every selected namespace on every node.
- **Rationale**: The root command still owns canonical parsing, trimming, and deduplication of namespace selections. The renderer should consume that selection without turning it into empty visual rows that misrepresent pod placement.
- **Alternatives considered**:
  - Remove `SelectedNamespaces` from the view config entirely: rejected because scoped grouped output still benefits from explicit selection context.
  - Keep always rendering empty selected namespace groups: rejected because it is the behavior being fixed in issue `#61`.

## Decision 3: Keep alphabetical headings and pod-order parity

- **Decision**: Continue ordering rendered namespace headings alphabetically and preserve the existing pod order within each namespace group.
- **Rationale**: The issue only changes whether an empty namespace heading appears, not how matching namespaces or pods are ordered. Preserving these rules reduces regression risk and keeps the behavior aligned with the established grouped renderer.
- **Alternatives considered**:
  - Switch to user-supplied namespace order: rejected because ordering is not part of the defect and would introduce an unrelated contract change.
  - Re-sort pods within a namespace: rejected because the grouped path already has an established order.

## Decision 4: Keep nodes visible even when no namespace groups render

- **Decision**: Render nodes selected by the existing node-selection workflow even when none of their pods match the selected namespaces, but omit namespace headings beneath those nodes.
- **Rationale**: The issue explicitly requires that namespace filtering affect only rendered namespace sections and pods, not whether a node appears in the output.
- **Alternatives considered**:
  - Hide nodes without matching namespace groups: rejected because it changes a separate, user-visible contract for node visibility.
  - Render a placeholder namespace row or explanatory text under the node: rejected because the accepted behavior is to omit empty namespace sections entirely.

## Decision 5: Update documentation to remove the obsolete empty-group example

- **Decision**: Replace README text that currently states empty selected namespaces are rendered as grouping rows.
- **Rationale**: The current README documents the exact behavior reported as a bug. Leaving it unchanged would violate the repository constitution's documentation rule.
- **Alternatives considered**:
  - Leave README updates to implementation time without planning for them: rejected because documentation impact is already known and user-visible.
  - Keep the example and add a caveat: rejected because the example itself becomes incorrect after the fix.
