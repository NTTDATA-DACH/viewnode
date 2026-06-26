# Scenario View

**Purpose**: Produce the UC semantics for the architecture workflow. This view is the source for the logical, process, development, and physical views.

## Architecture Intent

`viewnode` is a read-only Kubernetes inspection CLI. The scenario layer must stabilize the operator-facing meaning of "inspect a cluster's nodes, pods, and containers — once, or repeatedly — against a chosen kubeconfig context and namespace scope," so later views never reinterpret that as a server, an agent, or a mutating tool.

## Core Tensions

| Tension | Current Tradeoff Direction | Scenario Consequence |
|---------|----------------------------|----------------------|
| One-shot inspection vs continuous monitoring | Single binary, single command, one flag (`--watch`) flips between modes | Both modes must produce the same kind of output; watch mode is "the same scenario, repeated" |
| Operator control vs sensible defaults | Many flags exist but defaults render a usable view | Scenario must remain meaningful without flags |
| Tolerate cluster volatility vs surface real failures | Read-only nature lets transient errors be retried in watch; fatal errors must still halt the process | Watch and one-shot must classify failure outcomes differently |

## Stable Boundaries

| Boundary | Must Remain Stable Because | Explicitly Does Not Cover |
|----------|----------------------------|---------------------------|
| Operator ↔ CLI | Scripts and operators depend on the documented command surface | Web UI, daemon mode, server APIs |
| CLI ↔ kubeconfig | Context/namespace selection is the only state the tool may mutate | Cluster-resource mutation |
| CLI ↔ cluster | Tool reads node/pod/container/metrics state and renders it | Workload changes, RBAC changes, admission decisions |

## Change Axes

| Expected Change | Isolated By | Scenario Impact |
|-----------------|-------------|-----------------|
| New filter/view flags | `cmd` flag surface + `srv` view assembly | Scenario surface grows, refresh model unchanged |
| Watch behavior refinement (e.g., interval, retry) | `--watch` flag semantics and scheduler | Refresh-cadence and failure-closure semantics change; one-shot scenario unchanged |
| New context/namespace operations | `viewnode ctx`/`viewnode ns` subcommand groups | Additive scenarios at clearly separated subcommand boundaries |

## Invariants

| Invariant | Scenario Evidence | Risk If Violated |
|-----------|-------------------|------------------|
| Without `--watch`, the tool produces exactly one rendering and exits | `cmd/root.go` (`if !watchOn { close(stopCh) }`), `README.md` | Operator scripting and exit-code-driven automation break |
| Documented CLI contracts (flags, aliases, output structure) are part of the product | Constitution principle I, README, `cmd/root.go` flag definitions | Krew plugin users and pipelines silently break |
| The tool is read-only against cluster resources; only kubeconfig state may be mutated, and only via explicit `ctx`/`ns` subcommands | `srv/api.go` (list-only), `cmd/ctx/setcurrent.go`, `cmd/ns/setcurrent.go` | Trust boundary collapses; tool becomes an admin client |
| Ctrl-C ends watch mode cleanly | `cmd/root.go` scheduler, OS signal semantics | Terminals are left in a broken state |

## Non-goals / Anti-patterns

| Non-goal / Anti-pattern | Why It Is Out of Scope or Harmful |
|-------------------------|-----------------------------------|
| Mutating cluster resources (apply/delete pods, nodes, namespaces beyond `current-namespace` selection) | Breaks the read-only trust boundary that makes `viewnode` safe to run anywhere `kubectl get` would be safe |
| Running as a long-lived server or daemon | Contradicts the operator-CLI scenario shape and the krew/binary distribution model |
| Hiding errors during watch refreshes silently | Operators must remain aware that "what they see may be stale" |

## Actors and Participants

| Actor / Participant | Goal | Responsibility | Boundary |
|---------------------|------|----------------|----------|
| Operator | Inspect cluster utilization quickly, optionally live | Invoke `viewnode`, choose context/namespace/filters, read terminal output | Terminal user |
| Kubernetes API server | Provide authoritative node/pod state | Respond to read requests | External system |
| Kubernetes metrics API | Provide node and pod resource usage | Respond to metrics list requests | External system; may be unavailable |
| Kubeconfig file(s) | Provide credentials, context list, current selections | Read for setup; written only by `ctx`/`ns set-current` | Operator workstation |
| Cobra / krew runtime | Dispatch the operator's invocation to the right command | Parse args, run command | Foundation |

## Use Cases

| Use Case | Actor | Goal | Preconditions | Scope Boundary |
|----------|-------|------|---------------|----------------|
| UC-Inspect-Once | Operator | See current node/pod/container layout once | Valid kubeconfig, reachable cluster | `viewnode` without `--watch` |
| UC-Inspect-Watch | Operator | Continuously see node/pod/container layout | Valid kubeconfig, reachable cluster | `viewnode --watch [N]` |
| UC-Filter-View | Operator | Restrict view to specific node/pod/namespace/container detail | Same as UC-Inspect-* | Flag combinations on the root command |
| UC-Manage-Context | Operator | List, see, or change the current kubeconfig context | Writable kubeconfig | `viewnode ctx ...` |
| UC-Manage-Namespace | Operator | List, see, or change the current namespace in current context | Writable kubeconfig | `viewnode ns ...` |
| UC-Validate-Flags | Operator | Get a clear error for incompatible flag combinations | Invocation with disallowed combination | Root command pre-execution validation |

## Scenario Paths

| Scenario | Main Path | Successful Outcome | Alternative / Failure Branches |
|----------|-----------|--------------------|--------------------------------|
| UC-Inspect-Once | Load kubeconfig → fetch nodes/pods/metrics → render → exit 0 | Single rendered view, exit 0 | Missing kubeconfig → exit non-zero with sentinel error; unreachable cluster → exit non-zero |
| UC-Inspect-Watch | Initial load+render (per UC-Inspect-Once) → enter loop → sleep N → reload+render → repeat until interrupt | Repeated rendering until Ctrl-C | First refresh failure → exit per UC-Inspect-Once (does not enter loop); mid-session refresh failure → render error frame, retry on next interval (issue #73 clarifications) |
| UC-Filter-View | Same as UC-Inspect-* with flag-derived view config | Filtered rendering | Disallowed flag combination → exit non-zero via UC-Validate-Flags |
| UC-Manage-Context | Subcommand parses target context → mutate kubeconfig current-context → confirm | Updated current-context in kubeconfig | Invalid context name → exit non-zero |
| UC-Manage-Namespace | Subcommand parses target namespace → mutate current-namespace in current context | Updated current-namespace | Invalid namespace or no current context → exit non-zero |
| UC-Validate-Flags | Pre-execution checks flag combination | Continue to UC-Inspect-* | Disallowed combination → fatal exit with explanatory message |

## Acceptance Semantics

| Acceptance Scenario | Observable Result | Must Hold | Not Covered |
|---------------------|-------------------|-----------|-------------|
| One-shot completion | Process exits after a single render | UC-Inspect-Once invariant | Watch-mode interval behavior |
| Watch refresh cadence | Refresh occurs no sooner than the configured interval after the previous refresh completed | UC-Inspect-Watch sleep-after-refresh invariant (issue #73) | Catch-up bursts (explicitly forbidden) |
| Watch resilience | Mid-session refresh errors render as error frames; the loop continues | UC-Inspect-Watch retry invariant (issue #73) | Behavior when the first refresh fails — that falls back to UC-Inspect-Once exit semantics |
| Read-only against cluster | No cluster-resource mutation occurs for any inspection scenario | Read-only invariant | `ctx`/`ns` scenarios, which mutate kubeconfig only |

## Scenario Gaps

| Gap | Affected Scenario | Why It Matters |
|-----|-------------------|----------------|
| TTY vs non-TTY rendering semantics under `--watch` | UC-Inspect-Watch | Whether watch mode in a pipe re-renders or appends is not formalized; spec #073 only forbids regressions |
| Behavior when metrics API is unavailable but core API works | UC-Inspect-Once / UC-Inspect-Watch | Code tolerates it implicitly; scenario contract is not stated explicitly |

## Prohibited Content

Do not write architecture components, class designs, APIs, database tables, implementation tasks, test strategy, deployment scripts, or framework choices here.
