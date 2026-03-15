# Data Model: Refactor `ctx list` and `ns list` commands

## Context List Entry

- **Purpose**: Represents one display row for `viewnode ctx list`.
- **Fields**:
  - `name`: context name shown to the user
  - `active`: whether the entry matches the kubeconfig current context
- **Validation rules**:
  - `name` must be non-empty
  - exactly zero or one entries may be active in a rendered list
- **Lifecycle**:
  - created from raw kubeconfig context names
  - sorted alphabetically by `name`
  - rendered as `[marker] name`

## Namespace List Entry

- **Purpose**: Represents one display row for `viewnode ns list`.
- **Fields**:
  - `name`: namespace name shown to the user
  - `active`: whether the entry matches the active namespace for the current context
- **Validation rules**:
  - `name` must be non-empty
  - empty active namespace input is normalized to `default` before comparison
  - exactly zero or one entries may be active in a rendered list
- **Lifecycle**:
  - created from namespace names returned by the Kubernetes API
  - sorted alphabetically by `name`
  - rendered as `[marker] name`

## Shared Preparation Input

- **Purpose**: Provides the minimum information needed to build ordered list entries without coupling retrieval logic across commands.
- **Fields**:
  - `names`: unsorted list of discovered names
  - `activeName`: selected active item after any command-specific normalization
- **Rules**:
  - names are sorted before entry creation
  - active comparison is exact string equality
  - rendering helpers must not mutate retrieval state or perform new I/O
