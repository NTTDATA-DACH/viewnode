---
name: issue-to-prd
description: Create a feature PRD from a GitHub issue and save it under tasks/prd-*.md. Use when a GitHub issue should be turned into a PRD for this repository.
---

Read the GitHub issue provided by the user and create a PRD markdown file.

Requirements:
- Apply relevant repository guidance from AGENTS.md.
- Keep the PRD implementation-aware but not implementation-heavy.
- Align CLI behavior with existing project command conventions.
- Include:
    - problem statement
    - goal
    - user story
    - scope
    - non-goals
    - acceptance criteria
    - implementation notes relevant for later Ralph conversion

Output:
- Write the PRD to the path requested by the user, typically `tasks/prd-<issue>-<short-name>.md`.