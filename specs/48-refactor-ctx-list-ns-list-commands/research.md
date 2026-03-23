# Research: Refactor `ctx list` and `ns list` commands

## Decision 1: Use explicit list-entry structs for display preparation

- **Decision**: Represent context and namespace output rows with dedicated entry structs containing the display name and active-state boolean.
- **Rationale**: The issue explicitly asks for a struct per list entry, and explicit entries make shared ordering and rendering behavior testable without relying on inline loops.
- **Alternatives considered**:
  - Keep anonymous loops and extract only printing: rejected because the duplicated “collect + sort + mark active” behavior would remain split across commands.
  - Use a single generic cross-domain struct: rejected because the repository does not otherwise use generics for Cobra command output, and separate domain types stay clearer for a small refactor.

## Decision 2: Share only list-preparation behavior, not data retrieval

- **Decision**: Extract the duplicated preparation logic that turns raw names plus an active selection into ordered display entries, while leaving context lookup and namespace API retrieval in their current command-specific paths.
- **Rationale**: The commands share formatting semantics, but they do not share the same data source. This keeps the refactor small and avoids forcing unrelated kubeconfig and API access concerns into one abstraction.
- **Alternatives considered**:
  - Introduce one shared service that both loads data and prints output: rejected because it would mix kubeconfig parsing and Kubernetes API access in a wider abstraction than the feature requires.
  - Duplicate new entry-building logic in each command: rejected because the issue specifically asks to remove duplicate loops.

## Decision 3: Treat empty active namespace as `default`

- **Decision**: Namespace entry preparation must resolve an empty active namespace to `default` before active-state comparison.
- **Rationale**: Repository guidance already treats an empty namespace on the active kubeconfig context as `default`, and this behavior matches Kubernetes kubeconfig semantics.
- **Alternatives considered**:
  - Compare directly against the raw setup namespace: rejected because an empty namespace would fail to mark `default` correctly.
  - Mark no namespace as active when the kubeconfig namespace is empty: rejected because it changes user-visible behavior and breaks current namespace semantics.

## Decision 4: Preserve command output exactly

- **Decision**: Keep the existing marker format (`[*]` / `[ ]`), alphabetical ordering, aliases, and error flow for both list commands.
- **Rationale**: The issue is scoped to refactoring, and the constitution requires consistent CLI behavior for user-facing commands.
- **Alternatives considered**:
  - Update wording or add metadata during refactor: rejected because it would turn an internal cleanup into a behavior change.
