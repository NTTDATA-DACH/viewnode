# Data Model: Context list filtering

## Context Filter Input

- **Purpose**: Represents the user-provided filter text for `viewnode ctx list`.
- **Fields**:
  - `rawValue`: the flag value supplied through `--filter` or `-f`
  - `normalizedValue`: the comparison value used for case-insensitive matching
- **Validation rules**:
  - empty input is treated the same as no filter
  - normalization must preserve substring intent while ignoring letter case
- **Lifecycle**:
  - read from the `ctx list` command flags
  - normalized before matching against discovered context names
  - reused for no-match messaging when zero contexts remain

## Context Result

- **Purpose**: Represents one context row eligible for `viewnode ctx list` output.
- **Fields**:
  - `name`: kubeconfig context name shown to the user
  - `isActive`: whether the context matches the current kubeconfig context
  - `matchesFilter`: whether the context satisfies the case-insensitive substring filter
- **Validation rules**:
  - `name` must be non-empty
  - `isActive` must remain tied to `rawConfig.CurrentContext`
  - only contexts with `matchesFilter=true` are rendered when a filter is present
- **Lifecycle**:
  - loaded from one raw kubeconfig read
  - compared against the normalized filter input
  - sorted alphabetically before display through the existing shared listing helper

## Filtered Output State

- **Purpose**: Captures how filtered context rows should be rendered after selection is applied.
- **Fields**:
  - `displayedContexts`: the contexts that remain after filtering
  - `activeContextIncluded`: whether the current context is present in `displayedContexts`
  - `renderedActiveMarker`: whether any displayed row should show `[*]`
- **Validation rules**:
  - `displayedContexts` must contain only matching contexts
  - `renderedActiveMarker` is false when `activeContextIncluded` is false
  - deterministic alphabetical ordering must be preserved for displayed rows
- **Lifecycle**:
  - derived after context-name filtering and before printing
  - used to ensure non-matching active contexts are not reintroduced into output

## Empty Result Context

- **Purpose**: Captures why `ctx list` produced no visible rows so output can distinguish a filtered zero-match result from kubeconfig errors or normal rendering.
- **Fields**:
  - `filterValue`: the normalized filter text that was applied
  - `isFiltered`: whether filtering caused the empty result path
  - `message`: the user-visible output string for the successful no-match case
- **Validation rules**:
  - filtered empty-state messaging must still represent a successful command outcome
  - the exact success-path wording is `no contexts matched filter "<value>"`
- **Lifecycle**:
  - created only when a non-empty filter is applied and all contexts are filtered out
  - used instead of row rendering for that execution path

## Documentation Example

- **Purpose**: Records the CLI example content needed for `README.md`.
- **Fields**:
  - `command`: example command line
  - `expectedBehavior`: short explanation of what the command demonstrates
- **Validation rules**:
  - examples must use the documented flag names and aliases accurately
  - examples must reflect the clarified case-insensitive substring behavior and no-match wording
