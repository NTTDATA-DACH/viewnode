# kubectl-viewnode
The kubectl-viewnode shows nodes with their pods and containers.
You can find the source code and usage documentation at GitHub: https://github.com/NTTDATA-EMEA/kubectl-viewnode.

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

Use "kubectl-viewnode [command] --help" for more information about a command.
```
