# Architecture Synthesis: viewnode

**Input Views**:
- Scenario: `.specify/memory/architecture-scenario-view.md`
- Logical: `.specify/memory/architecture-logical-view.md`
- Process: `.specify/memory/architecture-process-view.md`
- Development: `.specify/memory/architecture-development-view.md`
- Physical: `.specify/memory/architecture-physical-view.md`

**Note**: This synthesis was produced by the architecture-reverse workflow after the five 4+1 view files were updated from observable repository facts in `architecture-repo-facts.md`.

## View Index

| View | File | Purpose | Current Status |
|------|------|---------|----------------|
| Scenario | `.specify/memory/architecture-scenario-view.md` | UC-producing actor, use case, path, branch, and acceptance semantics | Current |
| Logical | `.specify/memory/architecture-logical-view.md` | Capability boundaries, domain objects, states, and invariants | Current |
| Process | `.specify/memory/architecture-process-view.md` | Runtime links, handoffs, approvals, receipts, failure closure | Current |
| Development | `.specify/memory/architecture-development-view.md` | Architecture-level components, package boundaries, contracts, dependencies | Current |
| Physical | `.specify/memory/architecture-physical-view.md` | Deployment, external systems, fact sources, observability, operations | Current |

## Architecture Intent

`viewnode` is a single-binary, read-only Kubernetes inspection CLI that an operator runs against their existing kubeconfig and chosen cluster. Three logical capabilities — CLI Surface (incl. Watch Orchestration), Cluster Access, and View Assembly & Rendering — collaborate inside one short-lived process to produce a single rendered "refresh" or a sustained sequence of refreshes. Everything else (no server, no persistence of its own, no cluster-resource mutation) is intentionally absent and must stay absent.

## Central Design Forces

- **Refresh as the unit of work**: a one-shot invocation and a watch tick are the same logical operation; watch is "the refresh, repeated, with explicit failure-closure rules."
- **Capability split over package neatness**: the three capabilities are the SSOT for boundaries; current package layout (notably `srv/` carrying both API and rendering) is a soft, documented compromise rather than evidence of a single capability.
- **Operator-side, read-mostly trust model**: only kubeconfig state may be mutated, only by `ctx`/`ns` subcommands; cluster-side mutation is permanently out of scope.
- **Stable CLI contract**: the operator-facing command/flag/output surface is the externally observable architecture and is governed by constitution principle I.
- **Single-process, sleep-after-refresh runtime**: bounded, predictable cadence beats ticker-style throughput; failure-closure is fail-fast on the first refresh and operator-interrupt-only after that (issue #73 Q1+Q2).

## Primary Tradeoffs

| Tradeoff | Chosen Direction | Consequence | Revisit When |
|----------|------------------|-------------|--------------|
| Boolean vs interval-typed `--watch` | Becomes integer-seconds with `NoOptDefVal = 1` (issue #73) | Same flag surface, richer semantics | A new monitoring scenario demands sub-second cadence |
| Per-refresh state re-read vs process-scoped state | Re-read "where practical" at the CLI/Cluster Access boundary | Watch approximates repeated invocations | Deep layers gain caches that diverge from operator expectations |
| Classify error categories vs treat post-first-success errors as universally transient | Universal transient (issue #73 Q2) | Predictable, simple loop; only interrupt exits | A failure mode emerges where retrying is actively harmful |
| Co-located API + rendering in `srv/` vs split packages | Co-located, with documented dependency direction (`view → api`) | Smaller diff today; clear intent for future split | `srv/` becomes large enough that the soft boundary stops paying for itself |

## Stable Boundaries

| Boundary | Affected Views | Must Remain Stable Because | Forbidden Crossing |
|----------|----------------|----------------------------|--------------------|
| CLI Surface ↔ Cluster Access | Logical, Process, Development | The CLI contract is the product; Cluster Access is the only path to the cluster | Cluster I/O performed from `cmd/` directly without going through `cmd/config` + `srv` |
| Cluster Access ↔ Cluster | Scenario, Process, Physical | Read-only against cluster resources is a trust boundary | Any mutating call against cluster resources |
| Watch Orchestration ↔ Refresh | Scenario, Logical, Process | Watch must remain "the refresh, repeated" with documented failure rules | Concurrent refreshes; ticker-style burst; loop-exit on any post-first-success refresh error |
| Tool binary ↔ Operator host | Physical | The tool runs operator-side, with no server-side footprint | Long-lived deployment as a cluster workload |

## Change Axes

| Expected Change | Isolated By | Affected Views | Architecture Consequence |
|-----------------|-------------|----------------|--------------------------|
| New flags / filters / view options | CLI Surface (`cmd/`) + View Assembly (`srv/`) | Scenario, Logical, Development | Additive, no cadence or trust change |
| Watch semantics (interval, retry, render of errors) — issue #73 | Watch Orchestration in `cmd/root.go` | Scenario, Logical, Process | Loop contract evolves; CLI contract surface and capability split unchanged |
| Kubernetes client library evolution | Cluster Access (`cmd/config` + `srv`) | Development, Physical | Coordinated `client-go`/`metrics`/`apimachinery` version bump |
| Distribution targets (OS/arch, channels) | `.goreleaser.yaml` + `.krew.yaml` | Physical | Manifests updated together; no code change |

## Anti-patterns

| Anti-pattern | Why It Violates Intent | Affected Views |
|--------------|------------------------|----------------|
| `srv` importing `cobra` or owning flag parsing | Couples Cluster Access / Rendering to the CLI framework | Logical, Development |
| Background goroutines outliving a refresh | Breaks "refresh as unit of work" | Logical, Process |
| Caching cluster state across refreshes | Watch diverges from repeated one-shot semantics | Scenario, Logical, Process |
| Exiting a healthy watch loop on a refresh error | Reintroduces the issue #73 bug class | Scenario, Process |
| Embedding credentials or telemetry into the binary | Breaks the operator-side, kubeconfig-driven trust model | Scenario, Physical |

## Cross-View Architecture Model

| Architecture Concept | Scenario Meaning | Logical Interpretation | Runtime Role | Development Boundary | Physical Constraint | Architecture Constraint |
|----------------------|------------------|------------------------|--------------|----------------------|---------------------|---------------------------|
| Refresh | The unit of operator-visible inspection (one rendering of cluster state) | Composed work of Cluster Access + View Assembly, triggered by CLI Surface | One sequential pass per cycle; never concurrent within a process | Code path spans `cmd/root.go` → `cmd/config` → `srv` | Bounded by Kubernetes API/metrics availability | Watch ticks are refreshes, not a different operation |
| Watch Orchestration | "Keep showing me the latest" until interrupted | Sub-capability of CLI Surface, owning interval and failure-closure | Sleep-after-refresh loop, signal-driven exit | Lives in `cmd/root.go` scheduler | Operator-side only | First refresh fails like one-shot; subsequent refresh errors are always transient (issue #73) |
| Cluster Access | The operator's existing kubeconfig + cluster reachability | Authority over kubeconfig + cluster I/O | Synchronous list calls per refresh | `cmd/config` (setup) + `srv` (API) | One Kubernetes API server, one optional metrics API | Read-only against cluster resources |
| View Assembly & Rendering | The terminal output the operator reads | Capability that owns the operator-visible layout | Pure transform from API data to text | Within `srv/` | Operator terminal only | Never owns flag parsing or I/O |
| CLI Contract | The flag/command/output surface operators script against | Stable product contract | Defined entirely by CLI Surface | `cmd/` packages | Distributed via goreleaser + krew | Breaking changes require explicit spec + migration (constitution I) |
| Kubeconfig State | Operator-local selection (context, current namespace) | Only mutable state the tool may touch | Mutated only by `ctx`/`ns set-current`; read everywhere else | `cmd/ctx`, `cmd/ns`, `cmd/config` | Operator filesystem | Implicit mutation from inspection commands is forbidden |

## Key Architecture Conclusions

| Conclusion | Affected Views | Boundary/Owner | Consequence |
|------------|----------------|----------------|-------------|
| Watch is "refresh, repeated, with documented failure closure," not a separate scenario | Scenario, Logical, Process | Watch Orchestration (sub-capability of CLI Surface) | Issue #73 changes land entirely inside `cmd/root.go` scheduling; CLI Surface contract and Cluster Access stay untouched |
| Cluster I/O goes through `cmd/config` + `srv` only | Logical, Development | Cluster Access capability | Any new direct `client-go` call in `cmd/` outside `cmd/config` is an architectural violation |
| `srv/` is one package today but two capabilities | Development | View Assembly & Rendering, Cluster Access | If the package grows, split with `view → api` direction; do not introduce a parallel render or API package |
| Read-only against cluster, mutation only on kubeconfig via `ctx`/`ns` | Scenario, Logical, Physical | Cluster Access | A PR adding any cluster mutation is an architectural escalation, not a small change |
| Operator-side, single-binary delivery | Physical | Tool binary boundary | Server, daemon, or in-cluster deployment proposals are scope changes |

## Cross-Cutting Constraints

| Constraint | Source | Affected Views | Scope | Architecture Consequence |
|------------|--------|----------------|-------|--------------------------|
| Tests live next to changed code; race-detection coverage gate | Constitution III, `Makefile`, repo facts | Development | All packages | Behavior changes require nearby tests; CI enforces |
| Conventional Commits with `#<issue-number>` suffix on issue-backed work | `AGENTS.md`, gvb skill rules | Development | History | History remains auditable per-issue |
| Documentation and runtime guidance stay in sync with user-facing behavior changes | Constitution V | Scenario, Development, Physical | `README.md`, CLI help | Watch-mode changes (issue #73) require README/help updates in the same change |
| Krew + goreleaser manifest parity | `.krew.yaml`, `.goreleaser.yaml` | Physical | Release | OS/arch additions land in both manifests together |

## Open Risks and Review Triggers

| Risk or Trigger | Missing Evidence / Change Condition | Affected Views | Required Architecture Review |
|-----------------|-------------------------------------|----------------|------------------------------|
| `srv/` continues to grow | API + render diverge in change rates | Logical, Development | Consider splitting into `srv/api` and `srv/view` with `view → api` dependency |
| Watch behavior beyond issue #73 (e.g., sub-second interval, persistent error logs) | New scenario evidence | Scenario, Process | Re-derive Watch Orchestration invariants before implementing |
| TTY vs non-TTY rendering under `--watch` becomes a contract concern | Operator/CI complaint or new use case | Scenario, Physical | Add explicit scenario + physical-view acceptance |
| Demand for in-cluster deployment | New operator request | Scenario, Physical | Treat as a scope change to the deployment shape; full review required |
| Kubernetes client-go major bump | Upstream release | Development, Physical | Coordinated `client-go`/`metrics`/`apimachinery` update; verify `srv.Api` shape |
