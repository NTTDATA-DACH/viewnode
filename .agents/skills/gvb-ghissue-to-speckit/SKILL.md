---
name: gvb-ghissue-to-speckit
description: Create or update a Speckit feature specification from a GitHub issue URL or issue-backed feature request, while forcing the feature number to match the issue number.
compatibility: Requires spec-kit project structure with .specify/ directory
metadata:
  author: local-wrapper
  wraps:
    - speckit-specify
    - speckit-agent-context-update
---

# GVB Issue To Speckit Spec

Use this skill when the feature is tied to a GitHub issue and the branch number must match the issue number.
This workflow also requires every commit created for the work to append the GitHub issue number at the end of the subject line, for example `test(view): complete final polish #59`.
A bare GitHub issue URL is sufficient input for this skill.

## Inputs

- A GitHub issue URL such as `https://github.com/grafvonb/c8volt/issues/45`
- Optionally, extra feature description text in the same request

If the user provides only the issue URL, treat that as enough information to begin.
The skill must fetch and use the GitHub issue title and body as the primary feature-description context for Speckit.
For Codex, prefer the official GitHub API issue representation when accessible:
- input URL: `https://github.com/<owner>/<repo>/issues/<number>`
- preferred machine-readable source: `https://api.github.com/repos/<owner>/<repo>/issues/<number>`
- fallback source: the original GitHub issue page

## Required Behavior

1. Parse the GitHub issue URL from the user request.
2. Extract the numeric issue id from the URL path segment `/issues/<number>`.
3. Fetch the GitHub issue title and body from that URL before invoking Speckit.
   - Prefer the GitHub API issue endpoint derived from the provided issue URL.
   - If the API is unavailable, unauthorized, rate-limited, or otherwise unsuitable, fall back to the original issue page content.
4. Treat the fetched issue title and body as the primary feature-description context.
5. If extra user text is present in the same request, merge it with the fetched issue context instead of discarding either source.
6. Verify the installed shared Speckit shell layer behavior before creating the feature:
   - `.specify/scripts/bash/create-new-feature.sh --json --dry-run --number "<issue-number>" --short-name "<generated-short-name>" "<issue-title>"` must return a numeric `FEATURE_NUM`
   - the returned `BRANCH_NAME` must start with that numeric prefix and must not be guessed from the raw issue URL alone
   - `.specify/scripts/bash/common.sh` must accept repository-generated numeric feature branch prefixes through `check_feature_branch`
   - treat this as a behavioral check against the shipped `.specify` shell layer; do not assume patch helper scripts exist in the consumer repository
7. If that compatibility check fails or looks suspicious, stop immediately and tell the user that the installed Speckit tooling does not currently guarantee correct issue-backed branch naming and needs an ai-tooling sync for that repository.
8. Generate a concise 2-4 word short name using the same conventions as `speckit-specify`.
9. Run `.specify/scripts/bash/create-new-feature.sh` exactly once with:
   - `--json`
   - `--short-name "<generated-short-name>"`
   - `--number "<issue-number>"`
   - a feature description assembled from the GitHub issue title, GitHub issue body, and any extra user text
10. Treat the GitHub issue number as the canonical feature identity, while preserving repository naming conventions for the final branch and feature directory names.
11. Do not use the script's auto-number detection when an issue URL is present.
12. Do not synthesize or announce a speculative branch name from the raw issue URL alone. A branch like `issue-42` is a workflow bug unless the repository generator explicitly returned it.
13. If the issue URL is missing or malformed, stop and tell the user what was missing instead of guessing the issue number.
14. After the spec is generated, ensure the resulting `spec.md` includes explicit GitHub issue traceability:
   - issue number
   - issue URL
   - issue title
15. Before reporting success, verify project-owned Spec Kit artifacts are trackable by git:
   - `SPEC_FILE`
   - `FEATURE_DIR/checklists/requirements.md`
   - `.specify/extensions.yml`
   - `.specify/memory/constitution.md`
   - `.specify/memory/architecture.md`
   - `.specify/memory/architecture-scenario-view.md`
   - `.specify/memory/architecture-logical-view.md`
   - `.specify/memory/architecture-process-view.md`
   - `.specify/memory/architecture-development-view.md`
   - `.specify/memory/architecture-physical-view.md`
   - `.specify/memory/architecture-repo-facts.md`
16. If any required project-owned artifact path is ignored, stop and tell the user to fix the target repository `.gitignore` before continuing. Do not silently create issue-backed specs or memory artifacts that cannot be committed.
17. Respect enabled extension hooks such as `speckit.agent-context.update`; do not manually edit `AGENTS.md` or other agent context files when the bundled `agent-context` extension is installed.
18. For any commit created as part of this workflow or by downstream skills acting on the generated feature, preserve Conventional Commits formatting and append `#<issue-number>` at the end of the commit subject line.

## Branch Number Override

This skill intentionally overrides the default `speckit-specify` instruction that says not to pass `--number`.
Treat the GitHub issue number as the canonical branch identity for this workflow.

For this workflow, issue identifiers are assumed to come from GitHub issue URLs.
Treat the GitHub issue number as the canonical workflow identifier that downstream skills should preserve when deriving feature names, PRD names, Ralph branch names, or story IDs.
Treat the GitHub issue number as mandatory commit-message traceability as well: every generated commit subject must end with `#<issue-number>`.
The GitHub issue title and body are the canonical source of functional context unless the user explicitly overrides them with additional instructions.
Do not derive the short name or approval prompt from the raw issue URL token before the issue title and body have been fetched.

When a valid GitHub issue URL is provided, the issue number determines the feature identity.
Do not strip or rewrite repository-specific numbering conventions such as zero padding when the generator or repository already defines them.

Examples:

- `https://github.com/grafvonb/c8volt/issues/45` -> generated feature must remain tied to issue `45`
- `https://github.com/grafvonb/c8volt/issues/61` -> generated feature must remain tied to issue `61`
- `https://github.com/NTTDATA-DACH/viewnode/issues/47` -> generated feature must remain tied to issue `47`

## Short Name Rules

- Keep the short name concise and descriptive.
- Use hyphen-separated words.
- Prefer action-noun phrasing when it reads naturally.
- Preserve important technical terms and acronyms.

Examples:

- `Add user authentication` -> `user-auth`
- `Implement OAuth2 integration for the API` -> `oauth2-api-integration`
- `Add ns command list subcommand for listing namespaces in a context` -> `add-ns-list`

## Project Artifact Tracking Gate

Spec Kit project-owned artifacts are source-of-truth material and must remain reviewable in git. This wrapper keeps upstream Spec Kit paths intact and enforces tracking policy instead of moving artifacts to a local-only layout.

Before reporting completion, check whether any of these paths are ignored with `git check-ignore`:

- the generated `spec.md`
- the generated `checklists/requirements.md`
- `.specify/extensions.yml`
- `.specify/memory/constitution.md`
- `.specify/memory/architecture.md`
- `.specify/memory/architecture-scenario-view.md`
- `.specify/memory/architecture-logical-view.md`
- `.specify/memory/architecture-process-view.md`
- `.specify/memory/architecture-development-view.md`
- `.specify/memory/architecture-physical-view.md`
- `.specify/memory/architecture-repo-facts.md`

If any path is ignored, stop and report the ignored paths. Recommend replacing blanket `.specify/` ignores with selective rules that keep `.specify/memory/**`, `.specify/extensions.yml`, and `specs/**` trackable.

## Execution Flow

1. Read the user request and find the issue URL.
2. Fetch the issue title and body from GitHub.
3. Extract the issue number.
4. Determine the effective feature description from:
   - fetched issue title
   - fetched issue body
   - any explicit user text in the same request
5. Generate the short name from the fetched issue title and body, not from the raw issue URL.
6. Run the compatibility check first.
7. If the check is clean, run the script with `--number <issue-number>`.
8. Read the JSON output and inspect `BRANCH_NAME`, `SPEC_FILE`, and `FEATURE_NUM`.
9. Trust the script output directly for final naming and do not rename branches or spec directories after creation.
10. If the script output or any planned approval prompt would use a speculative URL-derived name such as `issue-42` instead of the repository-generated branch name, stop and surface that mismatch explicitly.
11. Continue with the normal `speckit-specify` workflow to write and validate the spec using the created names.
12. Ensure the generated spec records GitHub issue traceability explicitly so downstream skills do not depend on the original chat message.
13. Treat the wrapped `speckit-specify` checklist-and-validation step as mandatory completion criteria for this wrapper:
   - create `FEATURE_DIR/checklists/requirements.md`
   - update it with the current pass/fail validation status
   - do not report completion until both `spec.md` and `checklists/requirements.md` exist
14. Run the project artifact tracking gate. Stop if the feature spec, requirements checklist, `.specify/extensions.yml`, or project memory paths are ignored.
15. Respect any enabled `agent-context` hooks. If `speckit-agent-context-update` is available through `.specify/extensions.yml`, treat it as hook-managed context refresh rather than editing agent context files directly.
16. When reporting completion or handing off to downstream work, restate the commit rule explicitly: commit subjects must keep the Conventional Commit prefix and append `#<issue-number>` as the final token.

## Notes

- Only override feature identity selection by passing the issue number via `--number`. Keep the rest of the `speckit-specify` workflow the same, including repository numbering conventions in the script output.
- The shell-layer compatibility check is mandatory because downstream `speckit-plan` and `speckit-tasks` depend on the shared `.specify/scripts/bash` layer.
- Consumer repositories should receive ready-to-use `.specify` shell files from ai-tooling; do not depend on local patch helper scripts in the consumer repo.
- Do not run the feature creation script more than once per feature.
- Prefer behavioral checks over exact source-code signature checks because upstream Speckit may refactor scripts while preserving compatible behavior.
- If the installed `.specify/scripts/bash` state looks suspicious, stop and surface the exact warning instead of guessing whether downstream Speckit commands will work.
- Do not rely on the issue URL alone after spec creation; persist issue traceability inside `spec.md`.
- Do not create or accept project-owned Spec Kit artifacts that are hidden by target repository ignore rules. `specs/**`, `.specify/memory/**`, and `.specify/extensions.yml` must be trackable.
- `speckit-agent-context-update` is an upstream hook-managed context refresh command. Use enabled hooks when present; do not add a parallel manual `AGENTS.md` update path in this wrapper.
- If both a GitHub issue URL and an explicit conflicting `--number` are supplied by the user, prefer the issue number and mention that choice in the final report.
- Even though this wrapper focuses on issue-number preservation, it is not complete until `spec.md` and `checklists/requirements.md` are both present for the created feature directory.
- Commit subject examples:
  - `feat(spec): add issue traceability #45`
  - `test(view): complete final polish #59`
