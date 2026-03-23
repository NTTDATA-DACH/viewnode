# Data Model: Namespace list filtering

## Namespace Filter Input

- **Purpose**: Represents the user-provided filter text for `viewnode ns list`.
- **Fields**:
  - `rawValue`: the flag value supplied through `--filter` or `-f`
  - `normalizedValue`: the comparison value used for case-insensitive matching
- **Validation rules**:
  - empty input is treated the same as no filter
  - normalization must preserve substring intent while ignoring letter case
- **Lifecycle**:
  - read from the `ns list` command flags
  - normalized before matching against discovered namespace names
  - reused for no-match messaging when zero namespaces remain

## Namespace Result

- **Purpose**: Represents one namespace row eligible for `viewnode ns list` output.
- **Fields**:
  - `name`: namespace name shown to the user
  - `isActive`: whether the namespace matches the active namespace from kubeconfig
  - `matchesFilter`: whether the namespace satisfies the case-insensitive substring filter
- **Validation rules**:
  - `name` must be non-empty
  - `isActive` must continue to treat an empty kubeconfig namespace as `default`
  - only namespaces with `matchesFilter=true` are rendered when a filter is present
- **Lifecycle**:
  - loaded from one Kubernetes namespace list call
  - compared against the normalized filter input
  - sorted alphabetically before display

## Empty Result Context

- **Purpose**: Captures why a command produced no visible items so output wording can distinguish an empty filter result from other empty states.
- **Fields**:
  - `resourceType`: `namespace` or `node`
  - `filterValue`: the filter text that was applied, if any
  - `isFiltered`: whether filtering caused the empty result path
- **Validation rules**:
  - filtered empty-state messages must still represent a successful command outcome
  - generic empty-state messaging should remain available when no filter is involved
- **Lifecycle**:
  - created after filtering but before output
  - used only to select the correct user-facing message

## Documentation Example

- **Purpose**: Records the CLI example content needed for `README.md`.
- **Fields**:
  - `command`: example command line
  - `expectedBehavior`: short explanation of what the command demonstrates
- **Validation rules**:
  - examples must use the documented flag names and aliases accurately
  - examples must reflect the clarified case-insensitive substring behavior
