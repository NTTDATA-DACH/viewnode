# kubectl-viewnode
The `kubectl-viewnode` shows nodes with their pods and containers.
It is very useful when you need to monitor multiple resources such as nodes, pods or containers in a dynamic environment like a CI/CD platform.

```
Usage:
  kubectl-viewnode [flags]
  kubectl-viewnode [command]

Available Commands:
  help        Help about any command
  version     Plugin Version

Flags:
  -A, --all-namespaces         use all namespaces
  -d, --debug                  run in debug mode (shows stack trace in case of errors)
  -h, --help                   help for kubectl-viewnode
  -n, --namespace string       namespace to use
  -f, --node-filter string     show only nodes according to filter
  -p, --pod-filter string      show only pods according to filter
  -c, --show-containers        show containers in pod
  -t, --show-pod-start-times   show start times of pods
      --show-running-only      show running pods only

Use "kubectl-viewnode [command] --help" for more information about a command.
```
## Installation
The `kubectl-viewnode` is written in _go_, so just download the correct executable suitable for your platform from the releases and run it.
You can also download the source code and compile it.

If you copy the executable file into your PATH, you can use it like a `kubectl` command, i.e. without the hyphen:
```
kubectl viewnode --show-running-only
```

## Usage examples
Showing nodes and pods:
```
$ kubectl-viewnode
namespace: jenkins-onprem
8 pod(s) in total
0 unscheduled pod(s)
3 running node(s) with 8 scheduled pod(s):
- gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-8fws running 2 pod(s) (linux/amd64)
  * docker-in-the-cloud-341-ffkt5-2k64t (running)
  * docker-in-the-cloud-86-822pd-d6p3d (running)
- gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-b0np running 2 pod(s) (linux/amd64)
  * docker-in-the-cloud-338-3wc8r-n1t7z (running)
  * docker-in-the-cloud-340-cms5r-pxxq7 (running)
- gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-v1vr running 4 pod(s) (linux/amd64)
  * docker-in-the-cloud-337-4c4lm-q3mtp (running)
  * docker-in-the-cloud-339-4x9dq-0xd3q (running)
  * docker-in-the-cloud-342-r5khq-9cx2w (running)
  * liveness-test-4-boom (failed)
```
Showing nodes, pods and containers:
```
$ kubectl-viewnode --show-containers
namespace: jenkins-onprem
8 pod(s) in total
0 unscheduled pod(s)
3 running node(s) with 8 scheduled pod(s):
- gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-8fws running 2 pod(s) (linux/amd64)
  * docker-in-the-cloud-341-ffkt5-2k64t (running) (3: default/running docker-daemon/running jnlp/running)
  * docker-in-the-cloud-86-822pd-d6p3d (running) (3: default/running docker-daemon/running jnlp/terminated)
- gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-b0np running 2 pod(s) (linux/amd64)
  * docker-in-the-cloud-338-3wc8r-n1t7z (running) (3: default/running docker-daemon/running jnlp/running)
  * docker-in-the-cloud-340-cms5r-pxxq7 (running) (3: default/running docker-daemon/running jnlp/running)
- gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-v1vr running 4 pod(s) (linux/amd64)
  * docker-in-the-cloud-337-4c4lm-q3mtp (running) (3: default/running docker-daemon/running jnlp/running)
  * docker-in-the-cloud-339-4x9dq-0xd3q (running) (3: default/running docker-daemon/running jnlp/running)
  * docker-in-the-cloud-342-r5khq-9cx2w (running) (3: default/running docker-daemon/running jnlp/running)
  * liveness-test-4-boom (failed) (3: dead-container/terminated liveness1/running liveness2/running)
```
Showing nodes and pods with their start times:
```
$ kubectl-viewnode --show-pod-start-times
namespace: jenkins-onprem
8 pod(s) in total
0 unscheduled pod(s)
3 running node(s) with 8 scheduled pod(s):
- gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-8fws running 2 pod(s) (linux/amd64)
  * docker-in-the-cloud-341-ffkt5-2k64t (running/Thu Aug 26 09:36:04 CEST 2021)
  * docker-in-the-cloud-86-822pd-d6p3d (running/Thu Aug 19 09:09:20 CEST 2021)
- gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-b0np running 2 pod(s) (linux/amd64)
  * docker-in-the-cloud-338-3wc8r-n1t7z (running/Thu Aug 26 09:36:02 CEST 2021)
  * docker-in-the-cloud-340-cms5r-pxxq7 (running/Thu Aug 26 09:36:04 CEST 2021)
- gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-v1vr running 4 pod(s) (linux/amd64)
  * docker-in-the-cloud-337-4c4lm-q3mtp (running/Thu Aug 26 09:35:36 CEST 2021)
  * docker-in-the-cloud-339-4x9dq-0xd3q (running/Thu Aug 26 09:36:03 CEST 2021)
  * docker-in-the-cloud-342-r5khq-9cx2w (running/Thu Aug 26 09:36:05 CEST 2021)
  * liveness-test-4-boom (failed/Wed Aug 25 15:07:52 CEST 2021)
```
You can also combine show options:
```
$ kubectl-viewnode --show-pod-start-times --show-containers
namespace: jenkins-onprem
8 pod(s) in total
0 unscheduled pod(s)
3 running node(s) with 8 scheduled pod(s):
- gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-8fws running 2 pod(s) (linux/amd64)
  * docker-in-the-cloud-341-ffkt5-2k64t (running/Thu Aug 26 09:36:04 CEST 2021) (3: default/running docker-daemon/running jnlp/running)
  * docker-in-the-cloud-86-822pd-d6p3d (running/Thu Aug 19 09:09:20 CEST 2021) (3: default/running docker-daemon/running jnlp/terminated)
- gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-b0np running 2 pod(s) (linux/amd64)
  * docker-in-the-cloud-338-3wc8r-n1t7z (running/Thu Aug 26 09:36:02 CEST 2021) (3: default/running docker-daemon/running jnlp/running)
  * docker-in-the-cloud-340-cms5r-pxxq7 (running/Thu Aug 26 09:36:04 CEST 2021) (3: default/running docker-daemon/running jnlp/running)
- gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-v1vr running 4 pod(s) (linux/amd64)
  * docker-in-the-cloud-337-4c4lm-q3mtp (running/Thu Aug 26 09:35:36 CEST 2021) (3: default/running docker-daemon/running jnlp/running)
  * docker-in-the-cloud-339-4x9dq-0xd3q (running/Thu Aug 26 09:36:03 CEST 2021) (3: default/running docker-daemon/running jnlp/running)
  * docker-in-the-cloud-342-r5khq-9cx2w (running/Thu Aug 26 09:36:05 CEST 2021) (3: default/running docker-daemon/running jnlp/running)
  * liveness-test-4-boom (failed/Wed Aug 25 15:07:52 CEST 2021) (3: dead-container/terminated liveness1/running liveness2/running)
```
As well as filter nodes and pods:
```
$ kubectl-viewnode --node-filter v1vr
namespace: jenkins-onprem
4 pod(s) in total
0 unscheduled pod(s)
1 running node(s) with 4 scheduled pod(s):
- gke-dcgsecigke001-dcgsecigke001-linux-1cd8c3b9-v1vr running 4 pod(s) (linux/amd64)
  * docker-in-the-cloud-337-4c4lm-q3mtp (running)
  * docker-in-the-cloud-339-4x9dq-0xd3q (running)
  * docker-in-the-cloud-342-r5khq-9cx2w (running)
  * liveness-test-4-boom (failed)
```
Very popular is combining `kubectl-viewnode` with `watch` command e.g. watching all nodes, pods and containers every second can be configured as follows:
```
watch -n1 kubectl-viewnode --show-pod-start-times --show-containers
```
#Compatibility
The `kubectl-viewnode` was tested against _Google Cloud Platform_ with _Kubernetes_ v1.19.
It should however work with all clouds supported by [client-go](https://github.com/kubernetes/client-go).