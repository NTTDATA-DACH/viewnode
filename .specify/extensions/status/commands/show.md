---
description: "Show project status and SDD workflow progress — active feature, artifacts, task completion, phase, and extensions."
scripts:
  sh: scripts/bash/status.sh
  ps: scripts/powershell/status.ps1
---

# Project Status

Show a unified overview of the current Spec Kit project state and SDD workflow progress.

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## Outline

1. **Run status script**: Execute `{SCRIPT}` from the project root and parse the JSON output.

2. **Present project information**:
   - **Project name** — current directory name
   - **AI agent(s)** — which agent folders are detected, whether commands are present
   - **Script type** — bash, PowerShell, or both
   - **Git branch** — current branch name (or SPECIFY_FEATURE env var fallback)
   - **Feature directory** — resolved specs/NNN-name/ path

3. **Show SDD artifact status** for the current feature:
   - Check existence of: `spec.md`, `plan.md`, `tasks.md`, `research.md`, `data-model.md`, `quickstart.md`
   - Check `contracts/` directory (show file count if present)
   - Check `checklists/` directory (show count and pass/fail status)
   - Mark each with ✓ (exists) or ✗ (missing)

4. **Parse task progress** (if `tasks.md` exists):
   - Count lines matching `- [X]` or `- [x]` as completed
   - Count lines matching `- [ ]` as incomplete
   - Show: "12/18 completed (67%)"

5. **Detect workflow phase** based on which artifacts exist:
   - No spec.md → **Not Started** (run /speckit.specify)
   - spec.md only → **Plan** (ready for /speckit.clarify or /speckit.plan)
   - plan.md exists → **Tasks** (ready for /speckit.tasks)
   - tasks.md exists, not all done → **Implement** (ready for /speckit.implement)
   - All tasks [X] → **Complete**

6. **Show extensions summary**:
   - Count installed extensions from `.specify/extensions/` registry
   - Note: available catalog count is not shown (run `specify extension search` for that)

7. **Handle edge cases**:
   - On main/master branch with no feature → show project info + hint to switch branches
   - No git available → use SPECIFY_FEATURE env var or scan specs/ for latest feature
   - No features exist → show "No features created yet"
   - No .specify/ directory → script exits with error
