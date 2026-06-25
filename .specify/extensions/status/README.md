# Spec Kit Status

A Spec Kit extension that shows project status and SDD workflow progress.

## What it does

Shows a unified overview of your Spec Kit project:

- **Project info** — name, detected AI agent(s), script type
- **Current feature** — from git branch, SPECIFY_FEATURE env var, or specs/ scan
- **SDD artifacts** — which files exist (spec.md, plan.md, tasks.md, etc.)
- **Task progress** — completed/total with percentage from tasks.md
- **Workflow phase** — which SDD step you're on (Plan → Tasks → Implement → Complete)
- **Extensions** — how many are installed

## Installation

```bash
specify extension add status
```

Or install from GitHub:

```bash
specify extension add --from https://github.com/KhawarHabibKhan/spec-kit-status/archive/refs/tags/v1.0.0.zip
```

## Usage

```
/speckit.status.show
```

Or using the alias:

```
/speckit.status
```

## Requirements

- Spec Kit >= 0.1.0
