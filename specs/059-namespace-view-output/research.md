# Research: Namespace tree view for namespace-filtered node output

## Decision 1: Keep grouping in the printout layer

- **Decision**: Build scoped namespace groupings inside `srv/view.go` from each node's already loaded `[]ViewPod` data plus the selected namespace set passed in by `cmd/root.go`.
- **Rationale**: The feature changes presentation, not retrieval. Keeping grouping in the printout layer avoids altering node/pod loading behavior and reuses the same renderer seam already established for issue `#57`.
- **Alternatives considered**:
  - Precompute grouped structures during load/filtering: rejected because it broadens a rendering-only change into the data pipeline.
  - Reconstruct selected namespaces by reparsing `ViewNodeData.Namespace`: rejected because display text should not become the renderer's source of truth.

## Decision 2: Carry selected namespaces explicitly in the view context

- **Decision**: Extend the printout input with the explicit selected namespace slice already parsed by `cmd/root.go`.
- **Rationale**: Empty grouping rows for selected namespaces without pods on a node cannot be derived from the node's pod slice alone. Passing the parsed selection is the smallest repository-native way to support the clarified behavior without new API calls.
- **Alternatives considered**:
  - Add empty rows only for namespaces present in pod data: rejected because it violates the accepted clarification for issue `#59`.
  - Query namespaces per node or reload pod data: rejected because it adds unnecessary API work to a rendering feature.

## Decision 3: Sort namespace headings alphabetically and preserve grouped pod order parity

- **Decision**: Render one namespace heading per selected or populated namespace under each node, ordered alphabetically by namespace name, while keeping pod rows within each namespace group in the same order already used by grouped `--all-namespaces`.
- **Rationale**: Alphabetical grouping matches the accepted clarification, and pod-order parity keeps `--namespace` aligned with the established `--all-namespaces` contract instead of introducing another ordering rule.
- **Alternatives considered**:
  - Preserve the user-supplied namespace flag order: rejected because the clarification selected alphabetical grouping.
  - Alphabetize pods within each namespace group: rejected because the clarification chose parity with `--all-namespaces` behavior instead.

## Decision 4: Show empty namespace headings for selected namespaces on every displayed node

- **Decision**: When multiple namespaces are selected, render a grouping row for every selected namespace beneath each displayed node even if that node has no displayed pods in one or more of those namespaces.
- **Rationale**: This matches the accepted clarification and makes the scoped tree shape consistent from node to node.
- **Alternatives considered**:
  - Omit empty namespace headings: rejected because it conflicts with the accepted clarification.
  - Show empty headings only when multiple namespaces are present on the node: rejected because it would make the structure data-dependent and inconsistent.

## Decision 5: Keep summary, empty-state, and container behavior intact

- **Decision**: Preserve unscheduled pod rendering, node summary lines, root empty-state behavior, and both inline and tree container layouts while inserting scoped namespace grouping rows for scheduled pods.
- **Rationale**: Issue `#59` asks for a view refactor, not a broader redesign. Keeping adjacent behavior intact keeps the change small, testable, and repository-native.
- **Alternatives considered**:
  - Redesign the full scoped output flow: rejected because it expands scope and raises regression risk.
  - Skip container-mode validation because the issue is namespace-focused: rejected because grouped pod rows must still work in both container layouts.
