# Research: Context list filtering

## Decision 1: Add `--filter` / `-f` as a local `ctx list` flag

- **Decision**: Define the context filter flag on `cmd/ctx/list.go` so `viewnode ctx list --filter <value>` and `viewnode ctx list -f <value>` work without changing root flags or neighboring context subcommands.
- **Rationale**: The issue asks for ergonomics matching `viewnode ns list --filter`, and a local flag keeps the shorthand scoped to the context-listing command while preserving the existing root `-f` meaning for node filtering.
- **Alternatives considered**:
  - Reuse the root `--node-filter` / `-f` flag: rejected because it targets a different resource and would blur two distinct command contracts.
  - Add only `--filter` without `-f`: rejected because the issue explicitly requests both flag forms.

## Decision 2: Filter contexts client-side after one raw kubeconfig read

- **Decision**: Continue loading kubeconfig through the existing `config.Initialize(...)` and `RawConfig()` path, gather all context names once, then apply case-insensitive substring filtering in command logic before rendering.
- **Rationale**: `ctx list` already relies on local kubeconfig data rather than live cluster APIs, so client-side filtering keeps the current one-read command shape and avoids unnecessary architectural change.
- **Alternatives considered**:
  - Push filtering into a new shared service layer: rejected because the feature is small and the repository favors incremental changes close to the Cobra command implementation.
  - Introduce a new kubeconfig abstraction just for filtering: rejected because the existing raw-config path already provides all required data.

## Decision 3: Treat empty context filters as no filter

- **Decision**: Normalize an empty context filter value to the same behavior as omitting the flag and show the full context list.
- **Rationale**: This matches the clarification outcome and keeps the CLI forgiving for harmless blank input rather than turning it into an error or a misleading zero-match state.
- **Alternatives considered**:
  - Error on empty filter values: rejected because it adds friction without improving correctness.
  - Treat empty input as matching nothing: rejected because it would make `--filter ""` unexpectedly hide all contexts.

## Decision 4: Use exact no-match wording while keeping success status

- **Decision**: When filtering returns no contexts, `viewnode ctx list` should print `no contexts matched filter "<value>"` and still return success.
- **Rationale**: The clarified wording keeps CLI behavior consistent with the repository’s existing filter-aware messaging pattern and tells users that filtering succeeded but found no matches.
- **Alternatives considered**:
  - Keep silent empty output: rejected because it is ambiguous and looks like a rendering defect.
  - Return a command error for no matches: rejected because a valid zero-match search is not an operational failure.

## Decision 5: Show only matching contexts even when the current context is filtered out

- **Decision**: Filtered output should contain only contexts whose names match the filter, even when that removes the current context and leaves no active marker among the displayed rows.
- **Rationale**: This matches normal filtering semantics, avoids reintroducing non-matching rows solely for presentation, and keeps the result set easy to reason about in scripts and tests.
- **Alternatives considered**:
  - Always include the current context: rejected because it would violate the expectation that filtered output contains only matches.
  - Fail when the current context is filtered out: rejected because it would turn a legitimate filtered view into an unnecessary error path.

## Decision 6: Cover behavior through realistic command tests and README examples

- **Decision**: Add context-list coverage through Cobra command execution tests in `cmd/ctx/ctx_test.go` and update the README context command examples alongside the code change.
- **Rationale**: Repository guidance prefers tests that exercise actual command wiring, and the constitution requires documentation changes for user-visible CLI behavior.
- **Alternatives considered**:
  - Test only small helper functions: rejected because flag wiring and output wording are part of the requested behavior.
  - Skip README verification because the feature is small: rejected because the change is user-facing and must remain documented.
