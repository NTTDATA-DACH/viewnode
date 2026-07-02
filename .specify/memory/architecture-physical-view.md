# Physical View

**Input**: `.specify/memory/architecture-process-view.md`, `.specify/memory/architecture-development-view.md`

**Purpose**: Derive deployment, hosting, external system, fact-source, observability, and operational boundaries from process and development views.

## Architecture Intent

`viewnode` is a statically built CLI binary, run on operator workstations and CI runners, that talks to a Kubernetes cluster the operator already has access to. The physical view stabilizes the "operator-side tool, no server, no persistence, no networked daemon" deployment shape so future decisions don't drift toward a service architecture.

## Core Tensions

| Tension | Current Tradeoff Direction | Physical Consequence |
|---------|----------------------------|----------------------|
| Single distributable binary vs feature-specific images | Static `CGO_ENABLED=0` binary, multi-OS/arch via goreleaser; krew plugin packaging | Same artifact runs standalone and as a `kubectl` plugin |
| Cluster reachability via operator's kubeconfig vs embedded credentials | Always use the operator's kubeconfig and current context | No credential management owned by the tool |
| Local-only run vs scheduled remote monitoring | Local run only; watch loop is an operator-facing convenience, not a service | No availability/uptime SLO for `viewnode` itself |

## Stable Boundaries

| Boundary | Must Remain Stable Because | Explicitly Does Not Carry |
|----------|----------------------------|---------------------------|
| Operator workstation / CI runner | The only host class `viewnode` runs on | Server workloads, scheduled jobs, daemons |
| Kubeconfig file boundary | The only persistent state the tool may touch | Cluster-resource state |
| Kubernetes API server endpoint | Authoritative source for node/pod state | Tool state, configuration storage |
| Kubernetes metrics API endpoint | Optional authoritative source for resource usage | Required-at-all-times availability |

## Change Axes

| Expected Change | Isolated By | Physical Impact |
|-----------------|-------------|-----------------|
| New OS/arch target | `.goreleaser.yaml` builds + `.krew.yaml` platform selectors | Both manifests updated together; no runtime change |
| New release channel (e.g., another package manager) | Release pipeline only | No code change required |
| Kubernetes API version drift | `client-go` version bump in `go.mod` | Cluster Access tracking, no deployment topology change |

## Invariants

| Invariant | Source Deployment / External / Fact Boundary | Risk If Violated |
|-----------|----------------------------------------------|------------------|
| One binary, statically linked, no native runtime dependencies | `.goreleaser.yaml` (`CGO_ENABLED=0`, `-trimpath`) | Krew installs break on operators' systems |
| Tool stores nothing of its own outside the operator's kubeconfig | Process view (no persistence beyond a refresh; `ctx`/`ns set-current` only) | Hidden state would violate operator trust |
| Read-only against cluster resources at the network boundary | `srv/api.go` (list-only calls) | Tool ceases to be safe for read-only RBAC contexts |
| `viewnode` and `kubectl viewnode` invocation modes produce equivalent behavior | `.krew.yaml`, `main.go` | Krew users get a different product from binary users |

## Non-goals / Anti-patterns

| Non-goal / Anti-pattern | Why It Is Out of Scope or Harmful |
|-------------------------|-----------------------------------|
| Running `viewnode` as a long-lived Kubernetes Deployment / DaemonSet | Tool is operator-side, not cluster-side; would introduce credentials, lifecycle, and availability concerns absent today |
| Embedding kubeconfig or service-account credentials in the binary | Breaks the "operator already owns the credentials" trust model |
| Publishing telemetry/metrics from `viewnode` itself | Conflicts with "local tool, no server-side concerns" |

## Deployment and Hosting Boundaries

| Runtime / Hosting Unit | Carries | Boundary | Depends On | Release / Migration Impact |
|------------------------|---------|----------|------------|----------------------------|
| `viewnode` binary | All capabilities (CLI Surface, Cluster Access, View Assembly & Rendering) | Operator workstation / CI runner | OS + arch supported by goreleaser/krew manifests | New OS/arch must be added to both `.goreleaser.yaml` and `.krew.yaml` |
| Krew plugin install | Same binary, installed as `kubectl viewnode` | Operator workstation with `kubectl` + krew | krew plugin index acceptance | Plugin index PR per release |
| Operator's kubeconfig | Cluster credentials and selection state | Operator filesystem | Operator-managed | None at release time; mutated only by explicit `ctx`/`ns` subcommands |

## External System Collaboration

| External System | Purpose | Exchanged Content | Authoritative Fact | Failure Impact | Isolation / Substitute Boundary |
|-----------------|---------|-------------------|--------------------|----------------|---------------------------------|
| Kubernetes API server | Authoritative cluster state (nodes, pods, containers) | List requests; list responses | The cluster, not `viewnode` | First-refresh failure → exit; mid-session failure → transient (issue #73) | `srv.Api` interface enables test substitution via mock |
| Kubernetes metrics API | Resource usage for nodes and pods | List requests; list responses | The cluster | Missing/forbidden → degraded view; mid-session failure → transient | Tolerated absence at runtime; mock available |
| Kubeconfig file(s) | Credentials, context list, current selections | File read (always), file write (only on `ctx`/`ns set-current`) | Operator | Missing/invalid → fatal sentinel error | None — file system is the substrate |
| Krew plugin index | Plugin discovery and install | krew manifest | Krew maintainers | Out-of-date manifest → installs use older `viewnode` version | Direct binary downloads from GitHub Releases remain an alternative |
| GitHub Releases | Binary distribution | Tagged tarballs and checksums | GitHub | Release pipeline failure → no release | Manual goreleaser run is the fallback |

## Fact Sources and Observability

| Fact / Event | Authoritative Source | Observable Location | Consumers | Traceability Requirement |
|--------------|----------------------|---------------------|-----------|--------------------------|
| Cluster node/pod/container state | Kubernetes API server | Rendered terminal frame | Operator | Refresh frame must be self-identifying for the operator (issue #73 Q3 timestamps error frames) |
| Resource usage metrics | Kubernetes metrics API | Rendered terminal frame | Operator | Same as above; degraded gracefully when API absent |
| Tool version / commit | Build-time ldflags (`main.version`, `main.commit`) | `viewnode --version`, `--help` (where exposed) | Operator, support | Always present in released binaries |
| Operator-visible diagnostics | `logrus` output | Standard error stream | Operator | Errors must remain actionable per constitution principle V |

## Operations and Release Boundaries

| Operational Concern | Responsible Boundary | Trigger | Affected Views | Architecture Consequence |
|---------------------|----------------------|---------|----------------|--------------------------|
| Build & test on push/PR | GitHub Actions (`go.yml`) | Push / PR | Development | All packages must remain `go test -race`-clean (constitution III) |
| Tagged release | GitHub Actions (`release.yaml`) + goreleaser | Tag push | Physical | Updates GitHub Releases artifacts |
| Krew plugin manifest sync | `.krew.yaml` (templated by goreleaser) + krew plugin index PR | After release | Physical | Operators get the new version via krew |
| Operator-side install/upgrade | krew or direct binary download | Operator action | Physical | No coordinated rollout — operators upgrade independently |
| Failure of a watch loop in production use | Operator interrupt | Operator decision | Process, Physical | No on-call surface; operator restarts the command |

## Physical View Gaps

| Gap | Affected Deployment / External Boundary | Why It Matters |
|-----|-----------------------------------------|----------------|
| No published container image | Operator workstation / CI runner | If CI patterns shift to "use this in a container step," a thin image build target may be needed; not blocking today |
| No formalized RBAC documentation for "minimum viable kubeconfig" | Kubernetes API server boundary | Operators must currently infer required verbs (get/list across nodes, pods, metrics) |

## Prohibited Content

Do not write Kubernetes YAML, cloud resource manifests, machine sizes, service SKUs, deployment scripts, runbooks, or concrete infrastructure configuration here.
