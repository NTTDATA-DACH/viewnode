---
name: your-skill-name
description: Describe what the skill does and when Codex should use it. Include concrete trigger phrases, task types, file formats, systems, or workflows that should activate the skill.
---

# Purpose

State the skill's job in 1-3 short paragraphs.
Explain the outcome it should produce, the context it assumes, and the general operating style it should follow.

# Inputs

Accept any of:
- the user request
- relevant repository guidance such as `AGENTS.md`
- key files, directories, URLs, commands, or artifacts the skill should inspect
- optional user-provided constraints, preferences, or examples

If important inputs are omitted, say how the skill should resolve or infer them.

# Output

Write the expected result clearly:
- files to create or update
- decisions to report
- commands to run
- artifacts to validate

If the skill should avoid implementation and only plan, say so explicitly.

# Workflow

1. Read the user request and identify the concrete goal.
2. Read only the minimum repository or artifact context needed.
3. Apply any required conventions, constraints, or domain rules.
4. Use bundled resources when they provide a more reliable path than freeform reasoning.
5. Produce or update the target output.
6. Run the closest relevant validation.
7. Report assumptions, results, and any remaining risks.

Replace this sequence with the task-specific steps the skill should follow.

# Resource Usage

Use this section only if the skill includes bundled resources.

## Scripts

Use `scripts/...` when the work is repetitive, deterministic, or error-prone by hand.
List each script and when to run it.

Example:
- `scripts/example.sh`: Generate a normalized artifact from user input.

## References

Read `references/...` only when the request needs deeper domain guidance.
Link each reference directly from this file and explain when it matters.

Example:
- `references/schema.md`: Use when mapping output to the canonical schema.

## Assets

Use `assets/...` for templates, boilerplate, or files that should be copied or adapted rather than loaded into context.

# Decision Rules

Document the choices another Codex instance should make consistently.

Examples:
- prefer existing repository patterns over introducing new ones
- preserve externally observable behavior unless the user asks for changes
- stop and surface ambiguity instead of silently guessing in high-risk cases
- choose targeted validation before broad validation

# Validation

Run the closest relevant checks for the changed area first.
Escalate to broader validation only when needed.

Examples:
- run a specific test package
- validate generated JSON or YAML
- run a formatter or linter if the repository expects it
- inspect the produced file for required fields

If validation cannot run, say that explicitly in the final response.

# Example Triggers

Use short, concrete examples that help Codex recognize when this skill should activate.

- "Create a PRD for this feature"
- "Turn this issue into a Ralph PRD"
- "Review unresolved GitHub PR comments and address the actionable ones"

# Notes For Author

Before finalizing a real skill:
- replace every placeholder with task-specific guidance
- keep `description` explicit, because it is the main trigger surface
- keep this file concise and move heavy details into `references/` when needed
- add scripts only when they improve repeatability or reliability
- validate the finished skill with the repository's skill-validation path if available
