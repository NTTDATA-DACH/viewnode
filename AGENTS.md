# AGENTS.md

## Ralph PRD generation
- When generating or updating `scripts/ralph/prd.json`, split work into multiple small user stories.
- Each story must be feasible in one Ralph iteration.
- Prioritize stories by dependency and execution order.
- Keep each story narrowly scoped and independently verifiable.
- Write precise, minimal, testable acceptance criteria.
- Reuse existing project patterns and avoid introducing parallel command structures.
- If `scripts/ralph/prd.json` already exists and the task is an update, prefer adding follow-up stories instead of rewriting completed stories, unless they are incorrect.

## Git and GitHub branch rules
- For GitHub issues like #47, always check the existence of the branch on GitHib and use it, or create a branch accoring to GitHib algorithm if it does not exist, like:
  - for issue https://github.com/NTTDATA-DACH/viewnode/issues/47, the branch name should be `47-add-ns-command-list-subcommand-for-listing-namespaces-in-a-context`.
  - the branch name should be in the format `<issue_number>-<short-description>`, where:
    - `<issue_number>` is the number of the GitHub issue (e.g., `47`).
    - `<short-description>` is a concise, hyphen-separated description of the issue (e.g., `add-ns-command-list-subcommand-for-listing-namespaces-in-a-context`).
- Do not create or switch to any other feature branch unless the user explicitly asks.

## Commit rules
- Commit messages must follow Conventional Commits format (https://www.conventionalcommits.org/en/v1.0.0/).
- Use formats like, use capital letters only for text elements where it makes sense (e.g., README.md), and avoid unnecessary capitalization.
- Add a scope in parentheses after the type.
- Add GitHub issue reference in commit message like `feat(cli): add ns list command #47` if applicable.
- Example commit messages:
    - `feat(cli): add ns list command #47`
    - `test(cli): add tests for ns list #47`
    - `docs(readme): document ns list #47`
- Prefer small commits grouped by purpose.
- Do not use vague commit messages like `update`, `fix stuff`, or `changes`.

## Validation
- Before committing, run:
    - `make test`

## Command conventions
- Cobra subcommands that need the global `--kubeconfig` value should call `config.Initialize` with the root command (`cmd.Root()` fallback to the current command) so flag lookup resolves the root-level flag correctly.
- Commands that persist kubeconfig changes must set `clientcmd.ClientConfigLoadingRules.ExplicitPath` from `setup.KubeCfgPath` before calling `clientcmd.ModifyConfig`, otherwise `--kubeconfig` updates can be written to the default config file instead of the requested one.
- Cobra command groups should mirror the existing `cmd/ctx` and `cmd/ns` layout: keep the group command in `<group>.go`, put each subcommand in its own file, and register all subcommands from the group package `init()`.
- Namespace current-state commands should inspect `setup.ClientConfig.RawConfig()` for the active context value and treat an empty context namespace as `default`; this preserves kubeconfig semantics better than relying only on `Setup.Namespace`.

## Testing conventions
- For Cobra subcommands that depend on inherited persistent flags, prefer executing through a small root command tree with `Execute()` in tests instead of calling the leaf command's `RunE` directly; this matches Cobra's real flag propagation.

## Documentation conventions
- `README.md` command listings and CLI examples are maintained manually; when adding or changing user-facing commands, update the README command reference and relevant usage examples in the same change.

## Active Technologies
- Go 1.25.0 (toolchain go1.25.7) + `github.com/spf13/cobra`, `k8s.io/client-go`, `k8s.io/apimachinery`, `github.com/stretchr/testify` (48-refactor-ctx-list-ns-list-commands)
- Kubeconfig files plus live Kubernetes API data for namespace discovery (48-refactor-ctx-list-ns-list-commands)

## Recent Changes
- 48-refactor-ctx-list-ns-list-commands: Added Go 1.25.0 (toolchain go1.25.7) + `github.com/spf13/cobra`, `k8s.io/client-go`, `k8s.io/apimachinery`, `github.com/stretchr/testify`
