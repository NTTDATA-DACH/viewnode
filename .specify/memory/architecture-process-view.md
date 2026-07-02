# Process View

**Input**: `.specify/memory/architecture-scenario-view.md`, `.specify/memory/architecture-logical-view.md`

**Purpose**: Derive runtime collaboration, handoffs, approvals, receipts, state advancement, and failure closure from scenario paths and logical boundaries.

## Architecture Intent

`viewnode` runs as a single short-lived OS process per invocation. The process view stabilizes how a Refresh moves through the three capability boundaries (CLI Surface → Cluster Access → View Assembly & Rendering), how the Watch Orchestration sub-capability times and protects repeated refreshes, and how failures close out — without exiting a healthy loop.

## Core Tensions

| Tension | Current Tradeoff Direction | Process Consequence |
|---------|----------------------------|---------------------|
| Concurrency for responsiveness vs simplicity | One scheduling goroutine plus an error channel; refresh work is synchronous within the goroutine | Easy to reason about; no shared cluster-state mutex needed |
| Ticker-based cadence vs sleep-after-completion | Sleep after each completed refresh (issue #73 FR-007) | No catch-up bursts on slow refreshes; bounded resource usage |
| Fail fast vs keep watch alive | Fail fast on first refresh; keep alive on every subsequent refresh failure | Predictable startup, resilient steady state |

## Stable Boundaries

| Boundary | Must Remain Stable Because | Explicitly Does Not Control |
|----------|----------------------------|-----------------------------|
| Operator interrupt boundary (Ctrl-C / SIGINT and equivalents) | The only sanctioned way to leave an Active watch loop | Cluster API liveness, render details |
| CLI Surface → Cluster Access call boundary | Defines what one refresh "asks" the cluster | Refresh cadence, render details |
| Cluster Access → View Assembly handoff | Defines the per-refresh data the renderer receives | Cadence, interrupt handling |

## Change Axes

| Expected Change | Isolated By | Process Impact |
|-----------------|-------------|----------------|
| Watch interval becomes configurable (issue #73) | Watch Orchestration sub-capability | Sleep duration value changes; cadence shape unchanged (still sleep-after-refresh) |
| Per-refresh re-evaluation of command/config state | CLI Surface → refresh entry boundary | A single, well-defined "start of refresh" point becomes the place to reload state |
| Tolerance for additional Kubernetes error categories | Cluster Access error decoration | No process-level change as long as errors are returned to the loop, not raised as fatals |

## Invariants

| Invariant | Source Scenario / Runtime Link | Risk If Violated |
|-----------|--------------------------------|------------------|
| Exactly one OS process per invocation | Single CLI binary | Operators lose the "scripts can rely on this command" property |
| One Refresh at a time within a watch loop | Sleep-after-refresh model | Concurrent refreshes would race on rendering and on per-refresh state |
| Active watch loop exits only on operator interrupt | Issue #73 Q2 | Operators experience unpredictable terminations |
| First-refresh failure exits with one-shot exit code | Issue #73 Q1 | Hides startup misconfiguration behind a never-rendering loop |
| Error frames replace the current refresh's normal frame (not overlay) | Issue #73 Q3 | Operator cannot tell which refresh actually failed |

## Non-goals / Anti-patterns

| Non-goal / Anti-pattern | Why It Is Out of Scope or Harmful |
|-------------------------|-----------------------------------|
| Background goroutines that outlive a single refresh | Defeats the "refresh is a complete inspection scenario" invariant |
| Ticker-style scheduling that can fire a queued refresh immediately after a slow refresh | Reintroduces the catch-up storm forbidden by issue #73 FR-007 |
| Treating any post-first-success refresh error as fatal | Reintroduces the bug class issue #73 set out to fix |

## Main Runtime Links

| Runtime Link | Trigger | Source | Target | Transferred Content / Fact | Completion Condition |
|--------------|---------|--------|--------|----------------------------|----------------------|
| Invocation → CLI Surface | Operator runs `viewnode [...]` | OS / Cobra | CLI Surface | Parsed flags and selected use case | Use case dispatched |
| CLI Surface → Cluster Access (per refresh) | Refresh start | CLI Surface | Cluster Access | View config (filters, namespace scope), kubeconfig selection | API responses or decorated error returned |
| Cluster Access → View Assembly (per refresh) | API responses received | Cluster Access | View Assembly | Aggregated node/pod/metrics state | Rendered output emitted |
| View Assembly → Operator terminal | Render complete | View Assembly | Operator terminal | Tabular text (or error frame on failure) | Output flushed |
| Watch Orchestration → CLI Surface (per cycle) | Sleep elapsed | Watch Orchestration | CLI Surface | Signal to start next refresh | Refresh returns success or transient failure |
| Operator interrupt → Watch Orchestration | Ctrl-C / SIGINT / equivalent | Operator | Watch Orchestration | Stop signal | Loop exits cleanly |

## Handoffs and Approvals

| Handoff / Approval | From | To | Meaning | Accepted Path | Rejected / Returned Path |
|--------------------|------|----|---------|---------------|--------------------------|
| Flag-combination validation | CLI Surface (pre-execution) | CLI Surface (execution) | Operator's flag combination is internally consistent | Continue to refresh | Fatal exit with explanatory message |
| First-refresh outcome | First refresh | Watch Orchestration | "May the watch loop enter Active state?" | Enter Active loop | Exit process with one-shot exit code (issue #73 Q1) |
| Per-refresh outcome (after first success) | Refresh | Watch Orchestration | "Is the loop still healthy?" | Render normal frame and sleep | Render error frame and sleep (loop stays Active — issue #73 Q2) |
| Kubeconfig mutation (`ctx`/`ns set-current`) | Subcommand | Kubeconfig file | Operator-approved selection change | Confirmed write, process exits 0 | Invalid target → non-zero exit |

## Receipts and User Participation

| Receipt / Participation Point | Sender | Receiver | Content | User Action | Architecture Consequence |
|-------------------------------|--------|----------|---------|-------------|--------------------------|
| Rendered refresh frame | View Assembly | Operator terminal | Normal cluster view | Read, optionally interrupt | Operator sees the current state |
| Rendered error frame | View Assembly | Operator terminal | Labeled error with refresh timestamp + underlying error text (issue #73 Q3) | Read, optionally interrupt | Operator sees which refresh failed and why; loop continues |
| Process exit code | CLI Surface | Shell / pipeline | 0 on success, non-zero on validation/first-refresh/fatal errors | Branch in calling script | Scripts can rely on exit-code semantics |

## Failure, Degradation, and Closure

| Failure / Branch | Detection Boundary | Responsible Boundary | Degradation or Compensation | User-Visible Result | Closure Condition |
|------------------|--------------------|----------------------|-----------------------------|---------------------|-------------------|
| Invalid flag combination | CLI Surface pre-execution | CLI Surface | None — fatal | Explanatory message, non-zero exit | Process exits |
| Missing/invalid kubeconfig | Cluster Access | Cluster Access | None — fatal | Sentinel error text, non-zero exit | Process exits |
| First-refresh cluster/API failure | Cluster Access → CLI Surface | CLI Surface | None — fatal (issue #73 Q1) | Same as one-shot failure | Process exits with one-shot exit code |
| Mid-session refresh failure (any kind) | Cluster Access → Watch Orchestration | Watch Orchestration | Render error frame; sleep; retry next interval | Error frame in place of normal frame, then later normal frames if cluster recovers | Operator interrupt |
| Slow refresh (longer than configured interval) | Watch Orchestration | Watch Orchestration | Sleep-after-completion; no catch-up | Refresh cadence stretches with cluster latency | Normal refresh resumes |
| Operator interrupt | OS signal | Watch Orchestration | None — cooperative stop | Terminal returned to usable state | Process exits cleanly |

## Process Gaps

| Gap | Affected Runtime Link / Scenario | Why It Matters |
|-----|----------------------------------|----------------|
| No formalized signal-handling boundary in `cmd/root.go` beyond ticker/stop channels | Operator interrupt → Watch Orchestration | Issue #73 FR-006 requires clean cooperative stop; the architectural contract is clear but the runtime implementation needs explicit signal handling |
| No formalized "per-refresh re-load" boundary | CLI Surface → Cluster Access | FR-010 ("where practical") implies a boundary the runtime should respect but has not yet codified |

## Prohibited Content

Do not write call stacks, queue names, retry counts, thread/process details, endpoint sequences, workflow engine configuration, or orchestration code here.
