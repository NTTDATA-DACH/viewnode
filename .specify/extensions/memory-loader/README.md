# spec-kit-memory-loader

A [Spec Kit](https://github.com/github/spec-kit) extension that loads `.specify/memory/` files before spec-kit lifecycle commands so LLM agents have project governance context (constitution, glossary, conventions, resource standards).

## Problem

Spec-kit lifecycle commands (`/speckit.specify`, `/speckit.plan`, `/speckit.implement`, etc.) execute without awareness of project-specific governance documents stored in `.specify/memory/`. This means:

- Constitution principles are not consulted during specification
- Glossary terms are not available during planning
- Coding conventions are missed during implementation
- Resource standards are ignored during task generation

## Solution

The Memory Loader extension registers `before_*` hooks on all major spec-kit lifecycle commands. Before each command runs, it reads every `.md` file from `.specify/memory/` and outputs their contents, giving the LLM agent full governance context.

## Installation

```bash
# From release
specify extension add memory-loader --from https://github.com/KevinBrown5280/spec-kit-memory-loader/archive/refs/tags/v1.0.0.zip

# From main branch
specify extension add memory-loader --from https://github.com/KevinBrown5280/spec-kit-memory-loader/archive/refs/heads/main.zip

# Development mode (local clone)
specify extension add --dev /path/to/spec-kit-memory-loader
```

## Commands

| Command | Description | Modifies Files? |
|---------|-------------|-----------------|
| `speckit.memory-loader.load` | Read all project memory files and output their contents for context | No — read-only |

## How It Works

1. **Gather**: Reads every `.md` file from `.specify/memory/`
   - If the directory does not exist, skips silently
   - If a file cannot be read, skips it and continues

2. **Output**: For each file, prints a headed section:
   ```
   ## Memory: {filename}

   {file contents}
   ```

3. **Summarize**: After all files, outputs:
   ```
   Context loaded: {memory_count} memory files
   ```

## Hooks

The extension fires automatically before these lifecycle commands:

| Hook | Command | Description |
|------|---------|-------------|
| `before_specify` | `speckit.memory-loader.load` | Load context before specification |
| `before_plan` | `speckit.memory-loader.load` | Load context before planning |
| `before_tasks` | `speckit.memory-loader.load` | Load context before task generation |
| `before_implement` | `speckit.memory-loader.load` | Load context before implementation |
| `before_clarify` | `speckit.memory-loader.load` | Load context before clarification |
| `before_checklist` | `speckit.memory-loader.load` | Load context before checklist generation |
| `before_analyze` | `speckit.memory-loader.load` | Load context before analysis |

## Design Decisions

- **Read-only** — never modifies any files
- **Graceful degradation** — missing directory or unreadable files are skipped silently
- **Governance only** — loads project-level memory; feature-specific reference docs are handled separately by a companion extension

## Requirements

- Spec Kit >= 0.6.0

## License

MIT
