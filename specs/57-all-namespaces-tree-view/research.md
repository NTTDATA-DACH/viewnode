# Research: Namespace tree view for all-namespaces node output

## Decision 1: Group namespaces during printout, not during data loading

- **Decision**: Build namespace groupings inside `srv/view.go` from each node's already loaded `[]ViewPod` data.
- **Rationale**: The feature changes presentation, not retrieval. Keeping the grouping in the printout layer avoids altering node/pod loading behavior, preserves current API usage, and aligns with repository guidance to keep user-facing output branching in `srv/view.go`.
- **Alternatives considered**:
  - Enrich `srv.ViewNode` with pre-grouped namespace structures during loading: rejected because it broadens a rendering-only change into the load/filter pipeline.
  - Precompute namespace groups in `cmd/root.go`: rejected because the root command should keep passing minimal display context rather than owning view-shaping logic.

## Decision 2: Sort namespace groups alphabetically within each node

- **Decision**: Render one namespace heading per namespace under each node, ordered alphabetically by namespace name.
- **Rationale**: Alphabetical grouping makes the output predictable, easy to scan, and consistent with the accepted clarification recorded in the spec.
- **Alternatives considered**:
  - Preserve first-seen namespace order: rejected because it depends on pod ordering and is harder to predict visually.
  - Sort by pod count: rejected because it introduces a new prioritization rule not requested by the issue.

## Decision 3: Preserve current pod order within each namespace group

- **Decision**: Keep pod rows within a namespace group in their existing order from the loaded node data.
- **Rationale**: This limits the change to namespace grouping without silently changing pod ordering semantics at the same time.
- **Alternatives considered**:
  - Alphabetize pods inside each namespace group: rejected because it creates an additional sorting behavior not required by the issue.
  - Re-sort all pods globally after grouping: rejected because it risks changing existing expectations beyond the requested namespace layout.

## Decision 4: Keep the namespace grouping row even for a single namespace on a node

- **Decision**: Always show a namespace heading in `--all-namespaces` output, even when a node contains pods from only one namespace.
- **Rationale**: A consistent tree shape reduces conditional rendering complexity and matches the accepted clarification captured in the spec.
- **Alternatives considered**:
  - Omit namespace headings for single-namespace nodes: rejected because it creates a special case that makes the output less consistent.
  - Show namespace headings only for multi-namespace nodes: rejected because the structure would vary based on current data shape.

## Decision 5: Preserve existing summary and container rendering behavior

- **Decision**: Keep unscheduled pod rendering, node summary lines, and both inline and tree container layouts intact while inserting namespace grouping rows for scheduled pods.
- **Rationale**: The issue asks for a view refactor, not a broader redesign. Preserving adjacent output behavior keeps the change small, testable, and repository-native.
- **Alternatives considered**:
  - Redesign the whole node printout: rejected because it expands scope and increases regression risk.
  - Skip container-mode validation because the issue is namespace-focused: rejected because grouped pod rows must still work in both container layouts.
