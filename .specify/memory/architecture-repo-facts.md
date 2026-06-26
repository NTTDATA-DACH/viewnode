# Architecture Repo Facts: viewnode

**Purpose**: Record observable repository facts that support reverse generation of the project-level 4+1 architecture artifacts.

**Note**: This evidence file was filled in by the architecture-reverse workflow before the architecture views were updated. It is a fact source for architecture reasoning, not an implementation audit report.

## Repository Identity

| Fact | Evidence Source | Confidence | Architecture Relevance |
|------|-----------------|------------|--------------------------|
| Go module `viewnode`, Go 1.25.0, toolchain go1.25.7 | `go.mod` | High | Defines the single language baseline and toolchain constraint for the whole architecture |
| Single binary CLI named `viewnode`, distributed standalone and as a krew plugin | `main.go`, `.krew.yaml`, `.goreleaser.yaml`, `README.md` | High | Establishes the runtime surface: one process, two invocation modes (standalone, kubectl plugin) |
| Repository owned by NTTDATA-DACH; Apache-2.0 licensed | `README.md`, `LICENSE`, `.krew.yaml` (`homepage`) | High | Public-facing tool with stable URL identity used by krew manifest |

## Entry Points

| Entry Point | Type | Evidence Source | Observed Responsibility | Supported Scenario |
|-------------|------|-----------------|--------------------------|--------------------|
| `main.go` | Process entry | `main.go` | Wires version/commit and delegates to `cmd.Execute` | Process startup |
| `cmd.RootCmd` (`viewnode`) | Cobra root command | `cmd/root.go` | Loads kubeconfig, fetches nodes/pods/metrics, renders combined view, optional periodic refresh | Inspect cluster nodes with pods/containers/metrics |
| `cmd/ctx.CtxCmd` (`viewnode ctx ...`) | Cobra subcommand group | `cmd/ctx/ctx.go`, `list.go`, `getcurrent.go`, `setcurrent.go` | List, get-current, set-current kubeconfig contexts | Manage kubeconfig context |
| `cmd/ns.NsCmd` (`viewnode ns ...`) | Cobra subcommand group | `cmd/ns/ns.go`, `list.go`, `getcurrent.go`, `setcurrent.go` | List, get-current, set-current namespaces in current context | Manage Kubernetes namespace state |

## User-Visible Behaviors

| Behavior | Evidence Source | Actor / Trigger | Observable Outcome | Supported Use Case |
|----------|-----------------|-----------------|--------------------|--------------------|
| Render nodes with their pods, containers, requests/limits, and metrics in a single terminal view | `cmd/root.go` (`executeLoadAndFilter`, `executePrintOut`), `srv/view.go`, `srv/node.go`, `srv/pod.go` | Operator runs `viewnode` | Tabular terminal output of cluster state | Inspect cluster utilization |
| Filter view by node name, pod name, namespace(s), all-namespaces, container detail flags | `cmd/root.go` (flag definitions and `getViewNodeDataConfig`) | Operator passes flags | Subset of nodes/pods/containers rendered | Targeted inspection |
| Periodically refresh the rendered view (`--watch` / `-w`) | `cmd/root.go` (`watchOn`, `schedule`, `time.NewTicker(1 * time.Second)`) | Operator passes `--watch` | View re-renders on a fixed cadence until interrupted | Live monitoring |
| Manage kubeconfig contexts (list/get/set current) | `cmd/ctx/*.go`, `cmd/ctx/list_test.go` | Operator runs `viewnode ctx ...` | Context list printed; current-context mutated in kubeconfig | Context management |
| Manage namespaces (list/get/set current) | `cmd/ns/*.go`, `cmd/ns/ns_test.go`, `cmd/ns/setcurrent_test.go` | Operator runs `viewnode ns ...` | Namespace list printed; current namespace mutated | Namespace management |
| Reject mutually-incompatible flag combinations with a clear fatal message | `cmd/root.go` (`log.Fatalln` on `-r`/`-b` without `-c`) | Operator runs flag combination | Process exits with explanatory message | Input validation |

## System Boundaries

| Boundary | Evidence Source | Inbound Interaction | Outbound Interaction | Not Proven |
|----------|-----------------|---------------------|----------------------|------------|
| Operator terminal | `cmd/root.go`, `srv/view.go` | CLI flags and args | Rendered table output and log messages to stdout/stderr | TTY-vs-pipe behavioral switches |
| Kubeconfig file(s) | `cmd/config/setup.go` (`KubeCfgPath`, `clientcmd`), `cmd/ctx/setcurrent.go`, `cmd/ns/setcurrent.go` | Read for client setup and context/namespace discovery | Write current-context / current-namespace via `ctx`/`ns` set commands | Multi-file merge semantics beyond what `clientcmd` provides by default |
| Kubernetes API server | `srv/api.go` (`Clientset.CoreV1().Nodes()/Pods()`) | List nodes and pods | None (read-only at the cluster API level) | Mutating cluster operations |
| Kubernetes metrics API | `srv/api.go` (`MetricsClientset.MetricsV1beta1()`) | List node and pod metrics | None | Behavior when metrics API is absent (handled but not architecturally formalized) |

## Data and State Clues

| Fact / Entity | Evidence Source | Observed Lifecycle Clue | Fact Source | Not Proven |
|---------------|-----------------|-------------------------|-------------|------------|
| `Setup` (kubeconfig path, client config, namespace, clientsets) | `cmd/config/setup.go` | Built once per process invocation in `Initialize`; consumed by `srv` API layer | Built in-process from kubeconfig + cluster | None observed beyond process scope |
| Node / Pod / Container view model | `srv/view.go`, `srv/node.go`, `srv/pod.go`, `srv/transform.go` | Composed per refresh from API list responses + metrics | Kubernetes API + metrics API | Domain persistence (none — the tool is stateless) |
| Selected namespaces, filters, container view modes | `cmd/root.go` (`getViewNodeDataConfig`) | Constructed each refresh from flag state | CLI flag state | Whether downstream layers re-read flag/config state per refresh today (root.go currently does, but boundary is informal) |

## Runtime and Process Clues

| Runtime Fact | Evidence Source | Trigger / Handoff | Failure or Retry Clue | Not Proven |
|--------------|-----------------|-------------------|-----------------------|------------|
| Single-process, single-goroutine refresh loop with a goroutine-fed error channel | `cmd/root.go` (`schedule`, `handleErrors`, `errCh`, `sync.WaitGroup`) | Goroutine receives via channels; main goroutine waits | Errors from `executeLoadAndFilter`/`executePrintOut` are surfaced via `errCh`; current loop is ticker-driven | Whether errors classify as transient vs fatal (issue #73 evidence) |
| One-shot vs watch mode driven by a boolean `--watch` flag | `cmd/root.go` (`watchOn`, `BoolVarP`, `if !watchOn { close(stopCh) }`) | Same scheduler used for both modes | Watch mode currently exits on errors propagated through `errCh`/`log.Fatal` (basis for issue #73) | Sleep-after-refresh vs ticker-catch-up semantics formalized at the architecture level |
| Synchronous reads of kubeconfig and cluster state per command execution | `cmd/config/setup.go`, `srv/api.go` | Caller invokes `Initialize` then API list calls | Sentinel error `ErrKubeCfgFileNotExist` for missing kubeconfig; other errors decorated via `srv.DecorateError` | Cluster-side rate limiting / throttling strategy |

## Development Structure Clues

| Module / Package Area | Evidence Source | Observed Responsibility | Dependency Clue | Boundary Risk |
|-----------------------|-----------------|--------------------------|-----------------|---------------|
| `cmd/` (root + subcommand groups) | `cmd/root.go`, `cmd/ctx/`, `cmd/ns/`, `cmd/config/`, `cmd/internal/listing` | CLI surface, flag parsing, command orchestration, scheduling | Imports `cmd/config`, `cmd/ctx`, `cmd/ns`, `srv` | Root command currently owns scheduling loop in addition to wiring; conflates CLI concerns with runtime control |
| `cmd/config/` | `cmd/config/setup.go` | Kubeconfig discovery, REST config, clientset construction | Imports `client-go`, `clientcmd`, `metrics` clientset | Sole owner of cluster-credentials boundary |
| `cmd/ctx/`, `cmd/ns/` | `cmd/ctx/*.go`, `cmd/ns/*.go` | Context/namespace listing and mutation commands | Imports `cobra`, `clientcmd` indirectly via shared config | Two parallel sets of list/get/set commands (mirrored shape) |
| `srv/` | `srv/api.go`, `srv/view.go`, `srv/node.go`, `srv/pod.go`, `srv/transform.go`, `srv/api_mock.go`, `*_test.go` | Kubernetes API access (`Api` interface + `KubernetesApi` impl), view-model assembly, table rendering, error decoration | Imports `k8s.io/api`, `apimachinery`, `client-go`, `metrics` | Mixed concerns: API access and rendering coexist in one package |
| `utils/` | `utils/byte-count.go`, `utils/byte-count_test.go` | Byte-count formatting helpers | None outward | Small leaf utility; low risk |
| Top-level `main.go` | `main.go` | Process entry, version/commit injection, calls `cmd.Execute` | Imports `cmd` only | None |

## Repository-First Projection

Repository-first artifacts (`.specify/memory/repository-first/`) are not present in this repository. No repository-first projections recorded.

### Build Manifest Detection

| Ecosystem | Manifest Evidence | Detection Status | Runtime Surface Notes |
|-----------|-------------------|------------------|-----------------------|
| Go modules | `go.mod`, `go.sum` | Present | Single Go module rooted at repo root, module name `viewnode` |
| GoReleaser | `.goreleaser.yaml` | Present | Builds `linux`, `darwin`, `windows` binaries; CGO disabled |
| Krew plugin | `.krew.yaml` | Present | Distributes binaries for darwin/amd64, darwin/arm64, linux/amd64, windows/amd64 |
| GitHub Actions | `.github/workflows/go.yml`, `.github/workflows/release.yaml` | Present | CI and release pipelines |

### First-Party Module Edges

| From Module | To Module | Evidence Source | Observed Direction | Architecture Boundary Meaning |
|-------------|-----------|-----------------|--------------------|-------------------------------|
| `cmd` | `cmd/config`, `cmd/ctx`, `cmd/ns`, `srv` | `cmd/root.go` imports | Downward from CLI orchestration to capabilities | CLI surface owns wiring; downstream packages own behavior |
| `srv` | `cmd/config` | `srv/api.go` (`*config.Setup`) | Upward from capability into CLI config | Coupling worth watching: capability package depends on a config type owned under `cmd/` |

### Module Invocation Governance

| Rule Source | Allowed Direction | Forbidden Direction | Architecture Constraint | Risk If Violated |
|-------------|-------------------|---------------------|-------------------------|------------------|
| Implicit (no repository-first spec) | `cmd` → `srv`, `cmd` → `cmd/*` subcommand groups | `srv` should not own CLI orchestration; `utils` must not depend on `cmd`/`srv` | Keep CLI orchestration out of `srv` and `utils`; keep `utils` leaf | If violated, render/API package starts owning flag parsing or scheduling, complicating reuse |

### Dependency Governance Signals

| Signal Source | Dependency / Concern | Signal Type | Affected Boundary | Architecture Review Trigger |
|---------------|----------------------|-------------|-------------------|-----------------------------|
| `go.mod` | `k8s.io/{api,apimachinery,client-go,metrics}` all at v0.35.2 | Version alignment (intentional) | `cmd/config`, `srv` | Bumping client-go major must update all four together |
| `go.mod` | `github.com/spf13/cobra` v1.10.x | Single CLI framework | `cmd/*` | Changing CLI framework re-opens the entire CLI surface boundary |
| `go.mod` | `github.com/sirupsen/logrus` + `nested-logrus-formatter` | Single logging library | `cmd`, `srv` | Switching loggers affects every package that emits errors/messages |

## Physical / Deployment Clues

| Deployment Fact | Evidence Source | Environment / External System | Operational Constraint | Not Proven |
|-----------------|-----------------|-------------------------------|------------------------|------------|
| Statically built binary, `CGO_ENABLED=0`, `-trimpath`, ldflags inject version/commit | `.goreleaser.yaml`, `Makefile` | Any OS/arch listed in goreleaser/krew | No native dependency on host C libraries | Container image build |
| Distributed via GitHub Releases and krew index | `.goreleaser.yaml`, `.krew.yaml` | krew plugin index, GitHub Releases | Krew plugin index acceptance | Other package managers |
| No server-side deployment; tool runs on operator workstation/CI | `README.md`, lack of Dockerfile/Compose/K8s manifests | Operator workstation, CI runners | None at runtime beyond a kubeconfig and cluster reachability | Long-running daemon usage |

## Git History Signals

| Signal | Evidence Source | Architecture Meaning | Confidence | Review Trigger |
|--------|-----------------|----------------------|------------|----------------|
| Active feature work in `specs/` covering namespace and context filtering features (e.g., `059-namespace-view-output`, `069-add-ctx-filter`) | `specs/` directory listing | View/render and CLI flag surface evolve frequently | Medium | Future flag/filter additions should land through the same `cmd/*` + `srv/view` boundary, not new parallel packages |
| Tests present alongside almost every behavior-bearing file (`*_test.go`) | `cmd/*`, `srv/*`, `utils/*` | Test-backed behavior change is the established norm (matches constitution III) | High | Behavior changes without nearby tests are an architectural smell |

## Evidence Gaps

| Gap | Affected View | Why It Blocks Architecture Conclusion |
|-----|---------------|----------------------------------------|
| No formalized contract for transient vs fatal Kubernetes/API error classification | Process view, Logical view | The current `--watch` loop terminates on errors propagated via `errCh`; issue #73 is the first explicit clarification |
| `srv` package mixes API access and view rendering with no documented split | Development view | Cannot yet state whether render is a true sub-boundary or just file-level separation |
| TTY/pipe rendering behavior under `--watch` is undocumented | Scenario, Physical view | Cannot make an architecture conclusion about non-interactive output without source evidence beyond `cmd/root.go` |

## Evidence Rules

- Every non-placeholder fact must name an evidence source such as a file, directory, configuration, test, script, manifest, command output, or commit signal.
- Confidence values are `High`, `Medium`, or `Low`.
- `High` means multiple independent sources agree, or docs/tests and code entry points agree.
- `Medium` means one strong source is present, such as clear configuration, entry point, route declaration, or behavior test.
- `Low` means naming, directory structure, isolated code, or Git history suggests a fact but lacks behavior evidence.
- Git history is an auxiliary signal for change axes and boundary risks. It cannot independently prove architecture conclusions.
- Repository-first dependency matrices are fact inputs. Do not copy full dependency inventories into architecture views; summarize only architecture constraints, governance signals, gaps, or review triggers.
- Repository-first module invocation specs may support development-view dependency rules when each rule maps to a concrete module edge or dependency signal.
- Concrete classes, functions, fields, endpoints, database tables, and implementation data structures may appear only as evidence sources, not architecture conclusions.
