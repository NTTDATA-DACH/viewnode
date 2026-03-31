# Research: Clearer EOF error guidance

## Decision 1: Detect scoped EOF failures in `srv/api_errors.go`

- **Decision**: Add a dedicated decorated error type in `srv/api_errors.go` for cluster retrieval API `Get` failures whose surfaced message includes EOF, and keep the original error wrapped inside it.
- **Rationale**: The repository already classifies special Kubernetes failures at this boundary, so adding the new EOF case here keeps error recognition close to the API call sites and makes it easy for the root command to branch on typed errors.
- **Alternatives considered**:
  - Match raw error strings directly in `cmd/root.go`: rejected because it would spread classification logic into the command layer and diverge from the current decoration pattern.
  - Introduce a repository-wide sentinel package for all errors: rejected because the feature is small and does not justify a new abstraction.

## Decision 2: Scope the new decoration narrowly to Kubernetes API `Get` EOF failures

- **Decision**: Only decorate failures that still look like Kubernetes API request errors surfaced from the root retrieval path and end in EOF, rather than matching every EOF anywhere in the CLI.
- **Rationale**: The clarified spec explicitly limits the proxy hint to cluster retrieval API `Get` failures, and narrow matching reduces false positives for unrelated EOFs.
- **Alternatives considered**:
  - Decorate every EOF string in the application: rejected because it would over-apply the guidance and make unrelated failures less trustworthy.
  - Restrict the behavior to an exact single hard-coded endpoint string: rejected because it would be too brittle and would miss equivalent scoped retrieval failures.

## Decision 3: Add proxy guidance in `cmd/root.go` while preserving original details

- **Decision**: Handle the new decorated error in the root command’s failure switch and emit a fatal message that keeps the resource context and original low-level error text while adding a proxy hint phrased as a possible cause.
- **Rationale**: `cmd/root.go` already owns the translation from internal errors to user-facing fatal messages, and preserving the original detail matches the clarification outcome while still helping users diagnose proxy problems faster.
- **Alternatives considered**:
  - Replace the raw error with only a friendly proxy message: rejected because it would hide useful debugging context.
  - Log a second line with guidance and leave the main fatal text unchanged: rejected because a single deterministic fatal message is easier to test and reason about.

## Decision 4: Keep repeated scoped EOF failures deterministic

- **Decision**: Show the same proxy hint every time the scoped EOF failure occurs rather than suppressing repeats per run or session.
- **Rationale**: The CLI currently treats each failing iteration independently, and repeating the message avoids hidden state, keeps watch-mode behavior predictable, and matches the clarification outcome.
- **Alternatives considered**:
  - Show the hint only once per run: rejected because it adds stateful behavior with little value.
  - Show the hint only once per session: rejected because session-level suppression would be harder to test and would make repeated failures less self-explanatory.

## Decision 5: Cover both decoration and user-facing output in tests

- **Decision**: Add regression tests at the decoration boundary and the root command output boundary, then finish validation with `make test`.
- **Rationale**: The user-visible behavior depends on both correct classification and correct message translation, so testing only one layer would leave a meaningful gap.
- **Alternatives considered**:
  - Test only `DecorateError(...)`: rejected because the final CLI message is the real user-facing contract.
  - Test only the root fatal path with string fixtures: rejected because it would not prove the classifier is scoped correctly.

## Decision 6: Leave README command reference unchanged unless implementation adds troubleshooting guidance

- **Decision**: Treat this feature as a user-visible error-message refinement rather than a command-surface change, so no README command reference update is required unless implementation adds a focused troubleshooting note.
- **Rationale**: Commands, flags, and usage examples stay the same, and the most important durable record for this change is the spec, plan, contract, and automated tests.
- **Alternatives considered**:
  - Update the README command reference anyway: rejected because no command syntax or option descriptions change.
  - Add a mandatory troubleshooting section as part of this feature: rejected because it expands scope beyond the issue’s core behavior change.
