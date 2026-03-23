# viewnode
The `viewnode` displays Kubernetes cluster nodes with their pods and containers.
It is very useful when you need to monitor multiple resources such as nodes, pods or containers in a dynamic environment like a CI/CD platform.

Source code: [GitHub](https://github.com/NTTDATA-DACH/viewnode)

Copyright: © 2026 NTT DATA Deutschland SE, Adam Boczek

## Usage

```
Usage:
  viewnode [flags]
  viewnode [command]

Available Commands:
  ctx         Manage Kubernetes context
  help        Help about any command
  ns          Manage Kubernetes namespaces
  version     Plugin Version

Flags:
  -A, --all-namespaces             use all namespaces
  -b, --container-tree-view        format containers in tree view, otherwise inline
  -h, --help                       help for viewnode
      --kubeconfig string          kubectl configuration file (default: ~/.kube/config or env: $KUBECONFIG)
  -n, --namespace string           namespace to use; accepts comma-separated values
  -f, --node-filter string         show only nodes according to filter
  -p, --pod-filter string          show only pods according to filter
  -c, --show-containers            show containers in pod
  -m, --show-metrics               show memory footprint of nodes, pods and containers
  -t, --show-pod-start-times       show start times of pods
  -r, --show-requests-and-limits   show requests and limits for containers' cpu and memory (requires -c flag)
      --show-running-only          show running pods only
  -v, --verbosity string           defines log level (debug, info, warn, error, fatal, panic) (default "warning")
  -w, --watch                      executes the command every second so that changes can be observed

Use "viewnode [command] --help" for more information about a command.
```

## Installation
### As a Krew Plugin
Follow the instructions to [install](https://krew.sigs.k8s.io/docs/user-guide/setup/install/) krew and then run the following command:
```
kubectl krew install viewnode
```
The plugin will be available as `kubectl viewnode`.

### Standalone
The `viewnode` is written in _go_, so just download the correct executable suitable for your platform from the [releases](https://github.com/NTTDATA-DACH/viewnode/releases) and run it.

### Build It
You can also download the source code from [GitHub](https://github.com/NTTDATA-DACH/viewnode) and compile it by running the following command:
```
make build
```

## Usage examples
Showing nodes and pods:
```
$ viewnode
namespace(s): jenkins-onprem
8 pod(s) in total
0 unscheduled pod(s)
3 running node(s) with 8 scheduled pod(s):
├── gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-8fws running 2 pod(s) (linux/amd64)
│   ├── docker-in-the-cloud-341-ffkt5-2k64t (running)
│   └── docker-in-the-cloud-86-822pd-d6p3d (running)
├── gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-b0np running 2 pod(s) (linux/amd64)
│   ├── docker-in-the-cloud-338-3wc8r-n1t7z (running)
│   └── docker-in-the-cloud-340-cms5r-pxxq7 (running)
└── gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-v1vr running 4 pod(s) (linux/amd64)
    ├── docker-in-the-cloud-337-4c4lm-q3mtp (running)
    ├── docker-in-the-cloud-339-4x9dq-0xd3q (running)
    ├── docker-in-the-cloud-342-r5khq-9cx2w (running)
    └── liveness-test-4-boom (failed)
```
Showing nodes, pods and containers:
```
$ viewnode --show-containers
namespace(s): jenkins-onprem
8 pod(s) in total
0 unscheduled pod(s)
3 running node(s) with 8 scheduled pod(s):
├── gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-8fws running 2 pod(s) (linux/amd64)
│   ├── docker-in-the-cloud-341-ffkt5-2k64t (running) (3: default/running docker-daemon/running jnlp/running)
│   └── docker-in-the-cloud-86-822pd-d6p3d (running) (3: default/running docker-daemon/running jnlp/terminated)
├── gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-b0np running 2 pod(s) (linux/amd64)
│   ├── docker-in-the-cloud-338-3wc8r-n1t7z (running) (3: default/running docker-daemon/running jnlp/running)
│   └── docker-in-the-cloud-340-cms5r-pxxq7 (running) (3: default/running docker-daemon/running jnlp/running)
└── gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-v1vr running 4 pod(s) (linux/amd64)
    ├── docker-in-the-cloud-337-4c4lm-q3mtp (running) (3: default/running docker-daemon/running jnlp/running)
    ├── docker-in-the-cloud-339-4x9dq-0xd3q (running) (3: default/running docker-daemon/running jnlp/running)
    ├── docker-in-the-cloud-342-r5khq-9cx2w (running) (3: default/running docker-daemon/running jnlp/running)
    └── liveness-test-4-boom (failed) (3: dead-container/terminated liveness1/running liveness2/running)
```
Showing nodes and pods with their start times:
```
$ viewnode --show-pod-start-times
namespace(s): jenkins-onprem
8 pod(s) in total
0 unscheduled pod(s)
3 running node(s) with 8 scheduled pod(s):
├── gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-8fws running 2 pod(s) (linux/amd64)
│   ├── docker-in-the-cloud-341-ffkt5-2k64t (running/Thu Aug 26 09:36:04 CEST 2021)
│   └── docker-in-the-cloud-86-822pd-d6p3d (running/Thu Aug 19 09:09:20 CEST 2021)
├── gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-b0np running 2 pod(s) (linux/amd64)
│   ├── docker-in-the-cloud-338-3wc8r-n1t7z (running/Thu Aug 26 09:36:02 CEST 2021)
│   └── docker-in-the-cloud-340-cms5r-pxxq7 (running/Thu Aug 26 09:36:04 CEST 2021)
└── gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-v1vr running 4 pod(s) (linux/amd64)
    ├── docker-in-the-cloud-337-4c4lm-q3mtp (running/Thu Aug 26 09:35:36 CEST 2021)
    ├── docker-in-the-cloud-339-4x9dq-0xd3q (running/Thu Aug 26 09:36:03 CEST 2021)
    ├── docker-in-the-cloud-342-r5khq-9cx2w (running/Thu Aug 26 09:36:05 CEST 2021)
    └── liveness-test-4-boom (failed/Wed Aug 25 15:07:52 CEST 2021)
```
You can also combine show options:
```
$ viewnode --show-pod-start-times --show-containers
namespace(s): jenkins-onprem
8 pod(s) in total
0 unscheduled pod(s)
3 running node(s) with 8 scheduled pod(s):
├── gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-8fws running 2 pod(s) (linux/amd64)
│   ├── docker-in-the-cloud-341-ffkt5-2k64t (running/Thu Aug 26 09:36:04 CEST 2021) (3: default/running docker-daemon/running jnlp/running)
│   └── docker-in-the-cloud-86-822pd-d6p3d (running/Thu Aug 19 09:09:20 CEST 2021) (3: default/running docker-daemon/running jnlp/terminated)
├── gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-b0np running 2 pod(s) (linux/amd64)
│   ├── docker-in-the-cloud-338-3wc8r-n1t7z (running/Thu Aug 26 09:36:02 CEST 2021) (3: default/running docker-daemon/running jnlp/running)
│   └── docker-in-the-cloud-340-cms5r-pxxq7 (running/Thu Aug 26 09:36:04 CEST 2021) (3: default/running docker-daemon/running jnlp/running)
└── gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-v1vr running 4 pod(s) (linux/amd64)
    ├── docker-in-the-cloud-337-4c4lm-q3mtp (running/Thu Aug 26 09:35:36 CEST 2021) (3: default/running docker-daemon/running jnlp/running)
    ├── docker-in-the-cloud-339-4x9dq-0xd3q (running/Thu Aug 26 09:36:03 CEST 2021) (3: default/running docker-daemon/running jnlp/running)
    ├── docker-in-the-cloud-342-r5khq-9cx2w (running/Thu Aug 26 09:36:05 CEST 2021) (3: default/running docker-daemon/running jnlp/running)
    └── liveness-test-4-boom (failed/Wed Aug 25 15:07:52 CEST 2021) (3: dead-container/terminated liveness1/running liveness2/running)
```
As well as filter nodes and pods:
```
$ viewnode --node-filter v1vr
namespace(s): jenkins-onprem
4 pod(s) in total
0 unscheduled pod(s)
1 running node(s) with 4 scheduled pod(s):
└── gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-v1vr running 4 pod(s) (linux/amd64)
    ├── docker-in-the-cloud-337-4c4lm-q3mtp (running)
    ├── docker-in-the-cloud-339-4x9dq-0xd3q (running)
    ├── docker-in-the-cloud-342-r5khq-9cx2w (running)
    └── liveness-test-4-boom (failed)
```
You can also aggregate pods from selected namespaces only:
```
$ viewnode --namespace jenkins-onprem,default
8 pod(s) in total
0 unscheduled pod(s)
3 running node(s) with 8 scheduled pod(s):
├── gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-8fws running 2 pod(s) (linux/amd64)
│   ├── default: docker-in-the-cloud-86-822pd-d6p3d (running)
│   └── jenkins-onprem: docker-in-the-cloud-341-ffkt5-2k64t (running)
├── gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-b0np running 2 pod(s) (linux/amd64)
│   ├── default: docker-in-the-cloud-340-cms5r-pxxq7 (running)
│   └── jenkins-onprem: docker-in-the-cloud-338-3wc8r-n1t7z (running)
└── gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-v1vr running 4 pod(s) (linux/amd64)
    ├── default: docker-in-the-cloud-339-4x9dq-0xd3q (running)
    ├── default: liveness-test-4-boom (failed)
    ├── jenkins-onprem: docker-in-the-cloud-337-4c4lm-q3mtp (running)
    └── jenkins-onprem: docker-in-the-cloud-342-r5khq-9cx2w (running)
```
Very popular is combining `viewnode` with `watch` command e.g. watching all nodes, pods and containers every second can be configured as follows:
```
watch -n1 viewnode --show-pod-start-times --show-containers
```

### Context command
Manage kubeconfig contexts directly from `viewnode`:

```
$ viewnode ctx list
[*] dev-cluster
[ ] staging-cluster

$ viewnode ctx get-current
dev-cluster

$ viewnode ctx set-current staging-cluster
current context set to staging-cluster
```

Aliases:
- `viewnode ctx getcurrent`
- `viewnode ctx setcurrent staging-cluster`

### Namespace command
Manage namespaces in the current Kubernetes context directly from `viewnode`:

```
$ viewnode ns list
[ ] default
[*] jenkins-onprem
[ ] kube-system

$ viewnode ns list --filter kube
[ ] kube-public
[ ] kube-system

$ viewnode ns list -f missing
no namespaces matched filter "missing"

$ viewnode ns get-current
jenkins-onprem

$ viewnode ns set-current kube-system
current namespace set to kube-system
```

Aliases:
- `viewnode ns ls`
- `viewnode ns getcurrent`
- `viewnode ns setcurrent kube-system`

## Compatibility
`viewnode` uses Kubernetes `client-go` (currently `v0.35.2` in this repository).
It should work with any cloud platform and Kubernetes distribution supported by [client-go](https://github.com/kubernetes/client-go).

Compatibility generally follows the Kubernetes `client-go` policy for nearby cluster versions (for example, same minor version and close minor versions).
