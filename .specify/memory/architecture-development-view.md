# Development View

**Input**: `.specify/memory/architecture-logical-view.md`, `.specify/memory/architecture-process-view.md`

**Purpose**: Derive architecture-level components, package boundary intent, contract/artifact semantics, and dependency rules from logical and process views.

## Architecture Intent

Map the three logical capabilities — CLI Surface, Cluster Access, View Assembly & Rendering — onto stable package areas (`cmd/...`, `cmd/config`, `srv`, plus the leaf `utils`) so additions and refactors land in the right place. Treat the current cohabitation of API access and rendering inside `srv/` as a documented soft boundary, not as evidence that the two are one component.

## Core Tensions

| Tension | Current Tradeoff Direction | Development Consequence |
|---------|----------------------------|-------------------------|
| Strict per-capability packages vs current layout | Capability ≠ package one-to-one (Cluster Access spans `cmd/config` and `srv`; View Assembly lives inside `srv`) | Boundary rules are stated at the capability level; refactors should converge code with capability |
| Subcommand grouping in `cmd/*` vs root-command monolith | Subcommands grouped (`ctx`, `ns`); root command owns inspection and scheduling | Root command file is large but coherent; future split should follow capability lines |
| Tests close to code vs centralized test packages | Tests live next to code (`*_test.go`) | Behavioral changes must update nearby tests |

## Stable Boundaries

| Boundary | Must Remain Stable Because | Explicitly Must Not Own |
|----------|----------------------------|-------------------------|
| `cmd/` (root + subcommands) | The CLI Surface capability is the product's contractual surface | Cluster I/O, rendering implementation |
| `cmd/config` | Sole owner of kubeconfig discovery, client-config selection, and clientset construction | Flag parsing, watch scheduling, rendering |
| `srv/` | Houses Cluster Access (API) and View Assembly & Rendering capabilities; sole consumer of `client-go`/`metrics` for read paths | Flag parsing, scheduling, kubeconfig mutation |
| `cmd/ctx`, `cmd/ns` | Sole owners of kubeconfig-state mutation subcommands and their listings | Cluster reads (delegate to `cmd/config`/`srv`) |
| `utils/` | Leaf utility package | Any dependency on `cmd/*` or `srv` |

## Change Axes

| Expected Change | Isolated By | Development Impact |
|-----------------|-------------|--------------------|
| New flags on the root command | `cmd/root.go` flag block + `getViewNodeDataConfig` | Local to `cmd/` |
| New filter/render columns | `srv/view.go`, `srv/node.go`, `srv/pod.go`, `srv/transform.go` | Local to `srv/` |
| Watch interval/retry semantics (issue #73) | `cmd/root.go` scheduler + flag parsing | Local to `cmd/`; no API/render changes required |
| Kubeconfig discovery changes | `cmd/config/setup.go` | Local to `cmd/config`; consumers unchanged |
| New `ctx` or `ns` subcommands | `cmd/ctx`, `cmd/ns` | Additive to subcommand groups; pattern already established |

## Invariants

| Invariant | Source Boundary / Contract / Dependency Rule | Risk If Violated |
|-----------|----------------------------------------------|------------------|
| `srv` never imports `cobra` or owns CLI flag parsing | CLI Surface vs Cluster Access split | Render/API package becomes coupled to the CLI framework |
| `utils` has no inbound dependencies from `cmd/*` knowledge of orchestration | Leaf utility boundary | Hidden coupling forms; `utils` stops being safely reusable |
| All cluster-read code paths go through `srv.Api` (or the same package) so error decoration and metrics fallbacks stay uniform | Cluster Access invariant | Inconsistent error surfacing across commands |
| Tests live next to the code they cover (`*_test.go`) | Constitution III + observed pattern | Behavior regressions slip through |
| Conventional Commits with issue-suffixed subjects for issue-backed work | Workflow rule from `AGENTS.md` and gvb skills | Lost issue traceability in history |

## Non-goals / Anti-patterns

| Non-goal / Anti-pattern | Why It Is Out of Scope or Harmful |
|-------------------------|-----------------------------------|
| A parallel `internal/api` or `internal/render` package introduced without folding `srv` first | Creates two sources of truth for Cluster Access / View Assembly |
| Introducing a separate config framework (viper-style) alongside the current flag-driven model | Duplicates CLI Surface ownership of operator intent |
| Moving the `--watch` scheduler out of `cmd/` into `srv/` | Pulls runtime orchestration into Cluster Access, violating the capability split |

## Architecture-Level Components

| Component / Capability Package | Responsibility | Input / Output Boundary | Collaborators | Explicitly Must Not Own | Source View Evidence |
|--------------------------------|----------------|-------------------------|---------------|--------------------------|----------------------|
| CLI Surface (`cmd/`) | Parse flags, validate combinations, dispatch use cases, own watch scheduling | Operator invocation in / use-case calls out | `cmd/config`, `cmd/ctx`, `cmd/ns`, `srv` | Cluster I/O, rendering | Logical: CLI Surface; Process: CLI Surface ↔ Cluster Access |
| Cluster Access — Config (`cmd/config`) | Discover kubeconfig, build client and metrics clientsets | Kubeconfig path / `Setup` struct | `client-go`, `clientcmd`, `metrics` clientset | CLI parsing, rendering, scheduling | Logical: Cluster Access; Process: Cluster Access link |
| Cluster Access — API (within `srv/`) | Perform list-only reads against core and metrics APIs; decorate errors | `Setup` in / API list results out | `client-go`, `apimachinery`, `metrics` | Rendering implementation details | Logical: Cluster Access; Process: Cluster Access link |
| View Assembly & Rendering (within `srv/`) | Build view model from API results; render tabular output | API results + view config in / terminal text out | None outside `srv` | Cluster I/O, flag parsing | Logical: View Assembly & Rendering |
| Context Management (`cmd/ctx`) | List/get/set-current kubeconfig context | Operator subcommand in / kubeconfig mutation out | `clientcmd`, `cmd/config` | Cluster reads beyond what context discovery needs | Logical: Cluster Access kubeconfig mutation |
| Namespace Management (`cmd/ns`) | List/get/set-current namespace in current context | Operator subcommand in / kubeconfig mutation out | `clientcmd`, `cmd/config`, `srv` for listings | Cluster mutations | Logical: Cluster Access kubeconfig mutation |
| Shared Utilities (`utils/`) | Cross-cutting helpers (byte-count formatting) | Pure functions | None | Anything from `cmd/*` or `srv` | Development invariant only |

## Package Boundary Intent

| Package / Boundary | Abstraction Level | Owned Concepts | May Depend On | Must Not Depend On | Evolution Rule |
|--------------------|-------------------|----------------|---------------|--------------------|----------------|
| `cmd` (root) | CLI orchestration | Root command, global flags, watch scheduler, use-case dispatch | `cmd/config`, `cmd/ctx`, `cmd/ns`, `srv`, `cobra`, `logrus` | Direct `client-go` calls | New inspection options land here as flags + dispatch wiring |
| `cmd/config` | Cluster credentials and client construction | `Setup`, kubeconfig discovery, sentinel errors | `client-go`, `clientcmd`, `metrics` | Cobra root command internals | Changes to discovery rules go here, not into `cmd/root.go` |
| `cmd/ctx`, `cmd/ns` | Subcommand groups for kubeconfig-state operations | Listing and current-state mutation subcommands | `cobra`, `cmd/config`, `clientcmd` | `srv` for cluster mutation | New `ctx`/`ns` operations follow the existing list/get/set pattern |
| `cmd/internal/listing` | Internal listing helpers for subcommands | Shared listing utilities | None outside `cmd` (internal package) | n/a | Stays internal to `cmd/` |
| `srv` | Cluster Access (API) + View Assembly & Rendering | `Api` interface, `KubernetesApi`, view model assembly, table rendering, error decoration | `client-go`, `apimachinery`, `metrics`, `cmd/config` (for `Setup`) | `cobra` or anything CLI-framework specific | Future refactor may split into `srv/api` and `srv/view`; if so, dependency direction is `view → api`, not the reverse |
| `utils` | Leaf helpers | Formatting helpers | Standard library only | `cmd/*`, `srv` | Stays a leaf utility |

## Contracts and Artifacts

| Contract / Artifact | Semantics | Producer | Consumer | Lifecycle | Architecture Consequence |
|---------------------|-----------|----------|----------|-----------|--------------------------|
| `viewnode` CLI surface (commands, flags, output) | Stable operator-facing contract | `cmd/` packages | Operators, scripts, krew plugin users | Versioned with releases; constitutional invariant | Breaking changes require explicit spec + migration |
| `srv.Api` interface | Cluster-read contract used internally and by tests | `srv/api.go` | CLI Surface, tests (`srv/api_mock.go`) | Co-evolves with cluster-read needs | Extending it triggers updates in both `KubernetesApi` and mocks |
| `config.Setup` struct | Cluster credentials + clientsets bundle | `cmd/config` | `srv`, subcommands | Built per process invocation | Changing it ripples to all consumers; treat as a contract |
| Release artifacts (goreleaser archives, krew manifest) | Distributed binaries and plugin manifest | `.goreleaser.yaml`, `.krew.yaml` | krew users, GitHub Releases consumers | Per tagged release | Adding OS/arch requires updating both manifests together |

## Dependency Rules

| Rule | Allowed Direction | Forbidden Direction | Reason | Risk If Violated |
|------|-------------------|---------------------|--------|------------------|
| CLI Surface depends on Cluster Access and View Assembly | `cmd → cmd/config`, `cmd → srv` | `cmd/config → cmd`, `srv → cmd` (root) | Keeps capabilities reusable from non-root entry points | Cyclic, untestable dependencies |
| View Assembly depends on Cluster Access types | `srv (view) → srv (api)` | `srv (api) → srv (view)` | Renderer consumes data; API does not know about rendering | API code starts owning presentation concerns |
| `utils` is a leaf | nothing → `utils` (one-way) | `utils → cmd/*` or `utils → srv` | Keep helpers reusable and side-effect free | Hidden coupling |
| Kubernetes client libraries are bound to `cmd/config` and `srv` only | `cmd/config`, `srv` → `client-go`/`metrics`/`apimachinery` | Any other first-party package importing them | Concentrates cluster-API knowledge in two places | API changes ripple through unrelated packages |

## Development View Gaps

| Gap | Affected Component / Boundary | Why It Matters |
|-----|-------------------------------|----------------|
| `srv` mixes Cluster Access (API) and View Assembly & Rendering | `srv/` package boundary | The capabilities are clear logically but not packaged separately; future split is non-blocking but should follow the stated dependency rule |
| Watch scheduling lives inline in `cmd/root.go` | CLI Surface (root) | Issue #73 will refine the scheduler in place; if it grows further, extracting a small `cmd/internal/watch` may help, but is not required by this view |
| `cmd/internal/listing` responsibilities not yet inventoried in detail | CLI Surface internals | Listing helpers may be reusable across `ctx`/`ns`; worth confirming during the issue #73 plan only if relevant |

## Prohibited Content

Do not write source file paths, concrete package trees, classes, functions, implementation tasks, framework-specific wiring, or code generation notes here.
