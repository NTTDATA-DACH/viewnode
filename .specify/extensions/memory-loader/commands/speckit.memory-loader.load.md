---
description: "Read all project memory files and output their contents for LLM context"
---

# Load Project Memory

Read ALL `.md` files in `.specify/memory/` and output their contents. This gives you project governance context (constitution, glossary, conventions, resource standards) for the command that follows.

## Steps

1. **Gather**: Read every `.md` file from `.specify/memory/`.
   - If the directory does not exist, skip it silently.
   - If a file cannot be read, skip it and continue.

2. **Output**: For each file, print a headed section:

   ```
   ## Memory: {filename}
   
   {file contents}
   ```

3. **Summarize**: After all files, output:

   ```
   Context loaded: {memory_count} memory files
   ```

## Usage Notes

- Designed as a mandatory `before_*` hook that fires before spec-kit lifecycle commands.
- Loads governance context only. Feature-specific reference docs are loaded by the `spec-reference-loader` extension.
- This is a read-only operation — do NOT modify any files.
