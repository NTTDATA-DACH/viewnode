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