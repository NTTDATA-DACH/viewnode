# Logical View

**Input**: `.specify/memory/architecture-scenario-view.md`

**Purpose**: Derive capability boundaries, domain objects, states, relationships, and invariants from the scenario view.

## Architecture Intent

Separate "what the operator wants to see and control" from "how cluster state is obtained" from "how state is rendered for humans," so each can evolve without dragging the others along. The logical view names three capability boundaries — CLI Surface, Cluster Access, and View Assembly & Rendering — plus the kubeconfig-state authority owned by Cluster Access.

## Core Tensions

| Tension | Current Tradeoff Direction | Logical Consequence |
|---------|----------------------------|---------------------|
| Centralize CLI orchestration vs let each subcommand stand alone | Root command owns scheduling and overall flow; subcommands (`ctx`, `ns`) own their own subtrees | Watch/scheduling logic logically belongs to the root capability, not to subcommands |
| Single capability for cluster access vs split read-state vs mutate-state | Single Cluster Access capability owns kubeconfig setup and cluster reads; mutation is restricted to kubeconfig state via `ctx`/`ns` | Logical "mutation" is bounded to operator-local state only |
| Couple rendering tightly to API model vs introduce an intermediate view model | Current code transforms API objects into a view model before rendering | Logical "View Assembly" boundary exists even if today it shares a package with API access |

## Stable Boundaries

| Boundary | Must Remain Stable Because | Explicitly Does Not Own |
|----------|----------------------------|-------------------------|
| CLI Surface | Stable CLI contracts are a constitutional invariant | Cluster I/O, rendering implementation |
| Cluster Access | Sole owner of kubeconfig/cluster-credentials authority and Kubernetes API I/O | Flag parsing, rendering, scheduling |
| View Assembly & Rendering | Sole owner of operator-facing layout and formatting | Cluster I/O, flag parsing |
| Watch Orchestration (logical sub-capability of CLI Surface) | Determines refresh cadence and failure-closure semantics | Cluster I/O details, rendering implementation |

## Change Axes

| Expected Change | Isolated By | Logical Impact |
|-----------------|-------------|----------------|
| New flags, filters, or subcommands | CLI Surface | No change to Cluster Access or Rendering capability contracts |
| New columns/metrics in rendered view | View Assembly & Rendering | Cluster Access may grow new read calls; CLI Surface unchanged |
| Watch loop semantics (interval, retry) | Watch Orchestration sub-capability | CLI Surface and Cluster Access contracts unchanged |
| Kubernetes client-go version bumps | Cluster Access | Other capabilities unchanged as long as the read calls stay shape-stable |

## Invariants

| Invariant | Source Scenario / Object / State | Risk If Violated |
|-----------|----------------------------------|------------------|
| Cluster operations performed by `viewnode` are list-only against cluster resources | Read-only scenario invariant; `srv/api.go` only invokes `List` | The tool stops being a safe inspection tool |
| Only kubeconfig state may be mutated, and only by explicit `ctx`/`ns` subcommands | UC-Manage-Context, UC-Manage-Namespace | Implicit kubeconfig mutation from inspection commands is a trust violation |
| Watch Orchestration enters the loop only after at least one successful refresh; once inside, only an interrupt may exit it (issue #73 clarifications Q1+Q2) | UC-Inspect-Watch acceptance semantics | Operators see unpredictable exits or burst behavior |
| A refresh is the same logical work as a one-shot inspection | UC-Inspect-Watch built on UC-Inspect-Once | Output drift between modes |

## Non-goals / Anti-patterns

| Non-goal / Anti-pattern | Why It Is Out of Scope or Harmful |
|-------------------------|-----------------------------------|
| Caching cluster state across refreshes for performance | Defeats the "refresh = repeat the scenario" invariant |
| Pushing flag parsing into Cluster Access | Couples CLI surface to client-go internals |
| Sharing global mutable state across refreshes beyond what process-level flag/config state requires | Watch behavior diverges from one-shot behavior |

## Capability Boundaries

| Capability / Boundary | Responsibility | Input | Output | Explicitly Does Not Own | Scenario Source |
|-----------------------|----------------|-------|--------|--------------------------|-----------------|
| CLI Surface | Parse args, validate flags, dispatch use cases, own watch scheduling | Operator invocation | Calls into Cluster Access and View Assembly; exit code | Cluster I/O, formatting, kubeconfig mutation semantics beyond delegation | UC-Inspect-*, UC-Filter-View, UC-Validate-Flags |
| Watch Orchestration (sub-capability) | Decide when a refresh happens, classify refresh failures, surface failure to the operator | Configured interval, refresh outcome | Sequence of refresh invocations; lifecycle ends on interrupt | Cluster I/O, rendering | UC-Inspect-Watch |
| Cluster Access | Discover kubeconfig, build clients, read nodes/pods/metrics, mutate kubeconfig context/namespace on demand | Kubeconfig path, target context/namespace, read intent | API list responses, decorated errors; updated kubeconfig | Operator-facing rendering | UC-Inspect-*, UC-Manage-Context, UC-Manage-Namespace |
| View Assembly & Rendering | Assemble node/pod/container view model from API data, render to terminal | API list responses + view config | Tabular terminal output | API I/O, flag parsing | UC-Inspect-*, UC-Filter-View |

## Domain Objects and Relationships

| Object | Meaning | Owning Capability | Key Relationships | Fact Source | Invariants |
|--------|---------|-------------------|-------------------|-------------|------------|
| Operator Inspection Request | The intent expressed by a single CLI invocation | CLI Surface | Drives a Refresh (one or many) | Cluster (read), kubeconfig (selection) | Always parameterized by current context/namespace state |
| Refresh | A single end-to-end "obtain state → render" cycle | CLI Surface (orchestration) / Cluster Access (fetch) / View Assembly (render) | Composed once per invocation in one-shot, repeatedly in watch | Cluster API + metrics API | A refresh's output is the same kind as any other refresh's output |
| Cluster View Model | Operator-facing aggregation of nodes, pods, containers, requests/limits, metrics | View Assembly & Rendering | Derived from cluster API objects per refresh | Cluster API + metrics API | Never persisted beyond a single refresh |
| Kubeconfig Selection State | Current context + current namespace recorded in the operator's kubeconfig | Cluster Access | Read by every inspection; written only by `ctx`/`ns set-current` | Local kubeconfig file | Only mutated by explicit `ctx`/`ns` subcommands |

## State and Lifecycle

| Object / Flow | State | Entered When | Exited When | Forbidden Transition | Responsible Boundary |
|---------------|-------|--------------|-------------|----------------------|----------------------|
| Watch Loop | Not-started | Process begins with `--watch` | First refresh attempted | n/a | Watch Orchestration |
| Watch Loop | Active | First refresh succeeded | Operator interrupt (Ctrl-C / equivalent) | Exit on a refresh error after Active was entered (forbidden per issue #73 Q2) | Watch Orchestration |
| Watch Loop | Aborted-before-loop | First refresh failed | Process exits with the same code one-shot would (issue #73 Q1) | Entering Active without a successful first refresh | Watch Orchestration |
| Refresh | Succeeded | Render completes | Sleep begins (watch) or process exits (one-shot) | n/a | CLI Surface + View Assembly |
| Refresh | Failed-transient | Mid-session refresh error after at least one success | Error frame rendered, loop re-arms sleep | Exiting the loop | Watch Orchestration + View Assembly |

## Logical Decisions

| Decision | Scope | Owner / Boundary | Affected Objects or Flows | Consequence |
|----------|-------|------------------|---------------------------|-------------|
| Watch sleeps after the previous refresh completes (no catch-up) | Watch Orchestration | CLI Surface | Refresh cadence | Sustained refresh interval ≥ configured value; no burst behavior |
| Mid-session refresh errors are universally transient | Watch Orchestration | CLI Surface | Refresh lifecycle | Loop exits only on interrupt; simpler and more predictable failure model |
| First-refresh failure is treated as one-shot failure | Watch Orchestration | CLI Surface | Watch Loop entry | Same exit code as `viewnode` without `--watch` for the same error |
| Read-only against cluster; mutation only via `ctx`/`ns` to kubeconfig | All inspection paths | Cluster Access | All inspection scenarios | Operators can run `viewnode` anywhere `kubectl get` is allowed |

## Logical Gaps

| Gap | Affected Capability / Object | Why It Matters |
|-----|------------------------------|----------------|
| View Assembly and Cluster Access share a development package today | View Assembly & Rendering, Cluster Access | Logical boundary is clear, development boundary is not — captured in the development view |
| No formalized "per-refresh re-evaluation" contract for command/config state | Watch Orchestration | Issue #73 spec FR-010 requires it "where practical" but the logical boundary is not yet enforced anywhere |

## Prohibited Content

Do not write classes, DTOs, database tables, fields, method names, endpoints, schemas, or implementation data structures here.
