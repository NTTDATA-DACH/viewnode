# Spec Kit Architecture Workflow Extension

Generate or refresh project-level 4+1 architecture artifacts for a Spec Kit project.

This repository provides two independent Spec Kit integrations:

- An extension that installs the explicit architecture workflow commands.
- An optional preset that injects architecture SSOT context into core planning.

## Extension

The extension installs two commands:

```text
/speckit.arch.generate
/speckit.arch.reverse
```

It writes architecture memory under `.specify/memory/` and keeps the workflow separate from the Spec Kit core package.

### Install

After `v1.1.0` is published, install that release with:

```bash
specify extension add arch --from https://github.com/bigsmartben/spec-kit-arch/archive/refs/tags/v1.1.0.zip
```

Install from a local development checkout:

```bash
specify extension add --dev /home/administrator/github/spec-kit-arch
```

### Use

Run the registered command in your AI coding agent:

```text
/speckit.arch.generate
```

The command bootstraps the architecture artifacts if they are missing, reads existing architecture files and optional `.specify/memory/uc.md` context, then asks the agent to fill or refresh the 4+1 views and synthesis.

Run the reverse command in an ordinary historical repository when architecture should be inferred from repo facts:

```text
/speckit.arch.reverse
```

The reverse command first writes `.specify/memory/architecture-repo-facts.md`, then derives the 4+1 view files and synthesis from those evidence-backed facts. Use it when the repository does not already have Spec Kit feature context.

If `.specify/memory/repository-first/` exists, `/speckit.arch.reverse` treats those files as optional evidence inputs. It summarizes repository-first dependency matrices and module invocation specs into repo facts, then carries only architecture-level constraints, dependency rules, gaps, and review triggers into the 4+1 views.

### Files Written

The workflow writes only these files:

```text
.specify/memory/architecture-repo-facts.md
.specify/memory/architecture.md
.specify/memory/architecture-scenario-view.md
.specify/memory/architecture-logical-view.md
.specify/memory/architecture-process-view.md
.specify/memory/architecture-development-view.md
.specify/memory/architecture-physical-view.md
```

`/speckit.arch.generate` populates the six architecture SSOT files. Its setup may also bootstrap `.specify/memory/architecture-repo-facts.md` when missing. `/speckit.arch.reverse` populates those six files plus `.specify/memory/architecture-repo-facts.md` as its evidence layer.

It does not edit `.specify/memory/uc.md`, `.specify/memory/constitution.md`, feature specs, plans, tasks, source code, tests, root `docs/`, deployment manifests, or runbooks.

### Installed Layout

After installation, Spec Kit copies this extension under:

```text
.specify/extensions/arch/
```

The command uses the bundled setup scripts from that installed path:

```text
.specify/extensions/arch/scripts/bash/setup-arch.sh
.specify/extensions/arch/scripts/powershell/setup-arch.ps1
```

## Preset

The optional preset wraps only the core `/speckit.plan` command. It reads architecture files from `.specify/memory/` when present and grounds plan reasoning inside the architecture SSOT.

Install from a local development checkout:

```bash
specify preset add --dev /home/administrator/github/spec-kit-arch/presets/arch-governance
```

The preset does not override `/speckit.tasks` or `/speckit.implement`.

## Development

Validate the extension from a fresh project:

```bash
specify init /tmp/spec-kit-arch-test --ai codex --ignore-agent-tools --script sh
cd /tmp/spec-kit-arch-test
specify extension add --dev /home/administrator/github/spec-kit-arch
.specify/extensions/arch/scripts/bash/setup-arch.sh --json
```

Validate the preset from a fresh project:

```bash
specify init /tmp/spec-kit-arch-preset-test --ai codex --ignore-agent-tools --script sh
cd /tmp/spec-kit-arch-preset-test
specify preset add --dev /home/administrator/github/spec-kit-arch/presets/arch-governance
rg -n "Architecture SSOT Injection|Architecture Grounding Summary" .agents/skills/speckit-plan/SKILL.md
```

The extension intentionally does not provide the legacy `/speckit.arch` command.
