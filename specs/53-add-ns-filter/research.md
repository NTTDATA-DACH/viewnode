# Research: Namespace list filtering

## Decision 1: Add `--filter` / `-f` as a local `ns list` flag

- **Decision**: Define the namespace filter flag on `cmd/ns/list.go` so `viewnode ns list --filter <value>` and `viewnode ns list -f <value>` work without changing the root command flag set.
- **Rationale**: The issue requires the shorthand `-f` for namespace listing, and a local subcommand flag keeps that shorthand scoped to `ns list` while preserving the root command's existing node-filter behavior.
- **Alternatives considered**:
  - Reuse the root `--node-filter` / `-f` flag: rejected because it filters a different resource and would blur two separate command contracts.
  - Add only `--filter` without `-f`: rejected because the issue explicitly requires the shorthand alias.

## Decision 2: Filter namespaces client-side after one namespace list call

- **Decision**: Retrieve namespaces once through the existing Kubernetes namespace list path, then apply case-insensitive substring filtering in command logic before rendering.
- **Rationale**: Kubernetes namespace name substring matching is not naturally expressed as a server-side list selector here, and client-side filtering preserves the current single API request pattern.
- **Alternatives considered**:
  - Introduce additional API queries per filter attempt: rejected because it would add avoidable latency and complexity for an interactive CLI command.
  - Push filtering into a new shared service layer: rejected because the feature is small and the repository favors changes that stay close to the current Cobra command implementation.

## Decision 3: Treat empty namespace filters as no filter

- **Decision**: Normalize an empty namespace filter value to the same behavior as omitting the flag and show the full namespace list.
- **Rationale**: This matches the clarification outcome and prevents a harmless empty input from being treated as an error or a misleading zero-match result.
- **Alternatives considered**:
  - Error on empty filter values: rejected because it would make the command less forgiving without adding meaningful user value.
  - Treat empty values as matching nothing: rejected because it would create surprising behavior for a blank filter.

## Decision 4: Make no-match messaging filter-aware while keeping success status

- **Decision**: When filtering returns no results, `viewnode ns list` and the root node-filter flow should print a clear no-match message and still complete successfully.
- **Rationale**: The issue explicitly asks for a clearer result than the current generic `no nodes to display...`, and the clarification confirmed that no-match is a valid command outcome rather than an execution failure.
- **Alternatives considered**:
  - Keep the current generic empty-state text: rejected because it does not explain that filtering succeeded but matched nothing.
  - Convert no-match results into command errors: rejected because it would make valid filtered searches look like operational failures.

## Decision 5: Cover the behavior through command-level and printout tests

- **Decision**: Add namespace list tests through a small Cobra root tree for inherited config behavior, plus focused tests around root printout/no-match messaging and README examples.
- **Rationale**: Repository guidance prefers realistic Cobra execution for inherited flags, and the user-visible empty-state message needs direct regression coverage in the output layer.
- **Alternatives considered**:
  - Test only helper functions: rejected because flag wiring and output wording are key parts of the requested behavior.
  - Skip README verification because the feature is small: rejected because the constitution and repository instructions require README updates for user-facing CLI changes.
