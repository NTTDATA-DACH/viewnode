package srv

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func captureViewOutput(t *testing.T, data ViewNodeData) string {
	t.Helper()

	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	printErr := data.Printout(false)

	require.NoError(t, w.Close())
	os.Stdout = originalStdout
	require.NoError(t, printErr)

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)

	return buf.String()
}

func TestViewNodeDataPrintoutNamespaceLabel(t *testing.T) {
	output := captureViewOutput(t, ViewNodeData{
		Namespace: "kube-system,local-path-storage",
		Nodes: []ViewNode{
			{Name: ""},
		},
	})
	require.True(t, strings.Contains(output, "namespace(s): kube-system,local-path-storage\n"))
}

func TestViewNodeDataPrintoutNodeFilterNoMatchMessage(t *testing.T) {
	output := captureViewOutput(t, ViewNodeData{
		NodeFilter: "worker-a",
		Nodes: []ViewNode{
			{Name: ""},
		},
	})
	require.Equal(t, "no nodes matched filter \"worker-a\"\n", output)
}

func TestViewNodeDataPrintoutNoNodesWithoutFilterUsesGenericMessage(t *testing.T) {
	output := captureViewOutput(t, ViewNodeData{
		Nodes: []ViewNode{
			{Name: ""},
		},
	})
	require.Equal(t, "no nodes to display...\n", output)
}

func TestGroupPodsByNamespaceOrdersNamespacesAndPreservesPodOrder(t *testing.T) {
	groups := groupPodsByNamespace([]ViewPod{
		{Name: "zeta-1", Namespace: "zeta"},
		{Name: "alpha-1", Namespace: "alpha"},
		{Name: "zeta-2", Namespace: "zeta"},
		{Name: "beta-1", Namespace: "beta"},
	}, nil)

	require.Len(t, groups, 3)
	require.Equal(t, "alpha", groups[0].Namespace)
	require.Equal(t, []ViewPod{{Name: "alpha-1", Namespace: "alpha"}}, groups[0].Pods)
	require.Equal(t, "beta", groups[1].Namespace)
	require.Equal(t, []ViewPod{{Name: "beta-1", Namespace: "beta"}}, groups[1].Pods)
	require.Equal(t, "zeta", groups[2].Namespace)
	require.Equal(t, []ViewPod{
		{Name: "zeta-1", Namespace: "zeta"},
		{Name: "zeta-2", Namespace: "zeta"},
	}, groups[2].Pods)
}

func TestGroupPodsByNamespaceIncludesSelectedNamespacesWithoutPods(t *testing.T) {
	groups := groupPodsByNamespace([]ViewPod{
		{Name: "api-0", Namespace: "team-a"},
	}, []string{"team-b", "team-a", "team-c"})

	require.Len(t, groups, 3)
	require.Equal(t, "team-a", groups[0].Namespace)
	require.Equal(t, []ViewPod{{Name: "api-0", Namespace: "team-a"}}, groups[0].Pods)
	require.Equal(t, "team-b", groups[1].Namespace)
	require.Empty(t, groups[1].Pods)
	require.Equal(t, "team-c", groups[2].Namespace)
	require.Empty(t, groups[2].Pods)
}

func TestGroupPodsByNamespaceDeduplicatesRepeatedSelectedNamespaces(t *testing.T) {
	groups := groupPodsByNamespace(nil, []string{"team-b", "team-a", "team-b", "team-a"})

	require.Len(t, groups, 2)
	require.Equal(t, "team-a", groups[0].Namespace)
	require.Empty(t, groups[0].Pods)
	require.Equal(t, "team-b", groups[1].Namespace)
	require.Empty(t, groups[1].Pods)
}

func TestViewNodeDataPrintoutAllNamespacesGroupsPodsByNamespace(t *testing.T) {
	output := captureViewOutput(t, ViewNodeData{
		Config: ViewNodeDataConfig{
			ShowNamespaces:       true,
			GroupPodsByNamespace: true,
		},
		Nodes: []ViewNode{
			{Name: ""},
			{
				Name:             "worker-a",
				Os:               "linux",
				Arch:             "amd64",
				ContainerRuntime: "containerd://2.2.0",
				Pods: []ViewPod{
					{Name: "zeta-pod-1", Namespace: "zeta", Phase: "Running"},
					{Name: "alpha-pod-1", Namespace: "alpha", Phase: "Running"},
					{Name: "zeta-pod-2", Namespace: "zeta", Phase: "Pending"},
				},
			},
		},
	})

	require.Contains(t, output, "1 running node(s) with 3 scheduled pod(s):\n")
	require.Contains(t, output, "└── worker-a running 3 pod(s) (linux/amd64/containerd://2.2.0)\n")
	require.Contains(t, output, "    ├── alpha\n")
	require.Contains(t, output, "    │   └── alpha-pod-1 (running)\n")
	require.Contains(t, output, "    └── zeta\n")
	require.Contains(t, output, "        ├── zeta-pod-1 (running)\n")
	require.Contains(t, output, "        └── zeta-pod-2 (pending)\n")
	require.NotContains(t, output, "alpha: alpha-pod-1")
	require.NotContains(t, output, "zeta: zeta-pod-1")
	require.Less(t, strings.Index(output, "    ├── alpha\n"), strings.Index(output, "    └── zeta\n"))
}

func TestViewNodeDataPrintoutAllNamespacesShowsSingleNamespaceGroup(t *testing.T) {
	output := captureViewOutput(t, ViewNodeData{
		Config: ViewNodeDataConfig{
			ShowNamespaces:       true,
			GroupPodsByNamespace: true,
		},
		Nodes: []ViewNode{
			{Name: ""},
			{
				Name:             "worker-a",
				Os:               "linux",
				Arch:             "amd64",
				ContainerRuntime: "containerd://2.2.0",
				Pods: []ViewPod{
					{Name: "api-0", Namespace: "team-a", Phase: "Running"},
					{Name: "worker-0", Namespace: "team-a", Phase: "Running"},
				},
			},
		},
	})

	require.Contains(t, output, "    └── team-a\n")
	require.Contains(t, output, "        ├── api-0 (running)\n")
	require.Contains(t, output, "        └── worker-0 (running)\n")
}

func TestViewNodeDataPrintoutAllNamespacesPreservesSummariesAndUnscheduledPods(t *testing.T) {
	output := captureViewOutput(t, ViewNodeData{
		Config: ViewNodeDataConfig{
			ShowNamespaces:       true,
			GroupPodsByNamespace: true,
		},
		Nodes: []ViewNode{
			{
				Name: "",
				Pods: []ViewPod{
					{Name: "pending-0", Phase: "Pending"},
					{Name: "pending-1", Phase: "Running"},
				},
			},
			{
				Name:             "worker-a",
				Os:               "linux",
				Arch:             "amd64",
				ContainerRuntime: "containerd://2.2.0",
				Pods: []ViewPod{
					{Name: "api-0", Namespace: "team-a", Phase: "Running"},
					{Name: "dns-0", Namespace: "kube-system", Phase: "Running"},
				},
			},
		},
	})

	require.Contains(t, output, "4 pod(s) in total\n")
	require.Contains(t, output, "2 unscheduled pod(s):\n")
	require.Contains(t, output, "  ├── pending-0 (pending)\n")
	require.Contains(t, output, "  └── pending-1 (running)\n")
	require.Contains(t, output, "1 running node(s) with 2 scheduled pod(s):\n")
	require.Contains(t, output, "└── worker-a running 2 pod(s) (linux/amd64/containerd://2.2.0)\n")
	require.Contains(t, output, "    ├── kube-system\n")
	require.Contains(t, output, "    └── team-a\n")
}

func TestViewNodeDataPrintoutScopedNamespacesGroupPodsByNamespace(t *testing.T) {
	output := captureViewOutput(t, ViewNodeData{
		Config: ViewNodeDataConfig{
			ShowNamespaces:       true,
			GroupPodsByNamespace: true,
			SelectedNamespaces:   []string{"zeta", "alpha"},
		},
		Namespace: "zeta,alpha",
		Nodes: []ViewNode{
			{Name: ""},
			{
				Name:             "worker-a",
				Os:               "linux",
				Arch:             "amd64",
				ContainerRuntime: "containerd://2.2.0",
				Pods: []ViewPod{
					{Name: "zeta-pod-1", Namespace: "zeta", Phase: "Running"},
					{Name: "alpha-pod-1", Namespace: "alpha", Phase: "Running"},
					{Name: "zeta-pod-2", Namespace: "zeta", Phase: "Pending"},
				},
			},
		},
	})

	require.Contains(t, output, "namespace(s): zeta,alpha\n")
	require.Contains(t, output, "└── worker-a running 3 pod(s) (linux/amd64/containerd://2.2.0)\n")
	require.Contains(t, output, "    ├── alpha\n")
	require.Contains(t, output, "    │   └── alpha-pod-1 (running)\n")
	require.Contains(t, output, "    └── zeta\n")
	require.Contains(t, output, "        ├── zeta-pod-1 (running)\n")
	require.Contains(t, output, "        └── zeta-pod-2 (pending)\n")
	require.NotContains(t, output, "alpha: alpha-pod-1")
	require.NotContains(t, output, "zeta: zeta-pod-1")
	require.Less(t, strings.Index(output, "    ├── alpha\n"), strings.Index(output, "    └── zeta\n"))
}

func TestViewNodeDataPrintoutScopedNamespacesIncludeSelectedNamespacesWithoutPods(t *testing.T) {
	output := captureViewOutput(t, ViewNodeData{
		Config: ViewNodeDataConfig{
			ShowNamespaces:       true,
			GroupPodsByNamespace: true,
			SelectedNamespaces:   []string{"default", "team-a"},
		},
		Namespace: "default,team-a",
		Nodes: []ViewNode{
			{Name: ""},
			{
				Name:             "worker-a",
				Os:               "linux",
				Arch:             "amd64",
				ContainerRuntime: "containerd://2.2.0",
				Pods: []ViewPod{
					{Name: "api-0", Namespace: "default", Phase: "Running"},
					{Name: "web-0", Namespace: "default", Phase: "Running"},
				},
			},
		},
	})

	require.Contains(t, output, "    ├── default\n")
	require.Contains(t, output, "    │   ├── api-0 (running)\n")
	require.Contains(t, output, "    │   └── web-0 (running)\n")
	require.Contains(t, output, "    └── team-a\n")
}

func TestViewNodeDataPrintoutScopedNamespacesKeepPodOrderParityWithAllNamespaces(t *testing.T) {
	scopedOutput := captureViewOutput(t, ViewNodeData{
		Config: ViewNodeDataConfig{
			ShowNamespaces:       true,
			GroupPodsByNamespace: true,
			SelectedNamespaces:   []string{"team-b", "team-a"},
		},
		Namespace: "team-b,team-a",
		Nodes: []ViewNode{
			{Name: ""},
			{
				Name:             "worker-a",
				Os:               "linux",
				Arch:             "amd64",
				ContainerRuntime: "containerd://2.2.0",
				Pods: []ViewPod{
					{Name: "team-b-1", Namespace: "team-b", Phase: "Running"},
					{Name: "team-a-1", Namespace: "team-a", Phase: "Running"},
					{Name: "team-b-2", Namespace: "team-b", Phase: "Pending"},
					{Name: "team-a-2", Namespace: "team-a", Phase: "Pending"},
				},
			},
		},
	})
	allNamespacesOutput := captureViewOutput(t, ViewNodeData{
		Config: ViewNodeDataConfig{
			ShowNamespaces:       true,
			GroupPodsByNamespace: true,
		},
		Nodes: []ViewNode{
			{Name: ""},
			{
				Name:             "worker-a",
				Os:               "linux",
				Arch:             "amd64",
				ContainerRuntime: "containerd://2.2.0",
				Pods: []ViewPod{
					{Name: "team-b-1", Namespace: "team-b", Phase: "Running"},
					{Name: "team-a-1", Namespace: "team-a", Phase: "Running"},
					{Name: "team-b-2", Namespace: "team-b", Phase: "Pending"},
					{Name: "team-a-2", Namespace: "team-a", Phase: "Pending"},
				},
			},
		},
	})

	require.Less(t, strings.Index(scopedOutput, "team-a-1 (running)"), strings.Index(scopedOutput, "team-a-2 (pending)"))
	require.Less(t, strings.Index(scopedOutput, "team-b-1 (running)"), strings.Index(scopedOutput, "team-b-2 (pending)"))
	require.Less(t, strings.Index(allNamespacesOutput, "team-a-1 (running)"), strings.Index(allNamespacesOutput, "team-a-2 (pending)"))
	require.Less(t, strings.Index(allNamespacesOutput, "team-b-1 (running)"), strings.Index(allNamespacesOutput, "team-b-2 (pending)"))
}

func TestViewNodeDataPrintoutScopedNamespacesRemainInline(t *testing.T) {
	output := captureViewOutput(t, ViewNodeData{
		Config: ViewNodeDataConfig{
			ShowNamespaces: true,
		},
		Namespace: "alpha,beta",
		Nodes: []ViewNode{
			{Name: ""},
			{
				Name:             "worker-a",
				Os:               "linux",
				Arch:             "amd64",
				ContainerRuntime: "containerd://2.2.0",
				Pods: []ViewPod{
					{Name: "alpha-pod-1", Namespace: "alpha", Phase: "Running"},
					{Name: "beta-pod-1", Namespace: "beta", Phase: "Running"},
				},
			},
		},
	})

	require.Contains(t, output, "    ├── alpha: alpha-pod-1 (running)\n")
	require.Contains(t, output, "    └── beta: beta-pod-1 (running)\n")
	require.NotContains(t, output, "    ├── alpha\n")
}

func TestViewNodeDataPrintoutGroupedNamespacesKeepTreeContainerOutput(t *testing.T) {
	output := captureViewOutput(t, ViewNodeData{
		Config: ViewNodeDataConfig{
			ShowNamespaces:       true,
			GroupPodsByNamespace: true,
			ShowContainers:       true,
			ContainerViewType:    Tree,
			ShowTimes:            true,
		},
		Nodes: []ViewNode{
			{Name: ""},
			{
				Name:             "worker-a",
				Os:               "linux",
				Arch:             "amd64",
				ContainerRuntime: "containerd://2.2.0",
				Pods: []ViewPod{
					{
						Name:      "api-0",
						Namespace: "team-a",
						Phase:     "Running",
						StartTime: time.Unix(0, 0).UTC(),
						Containers: []ViewContainer{
							{Name: "web", State: "Running"},
							{Name: "sidecar", State: "Waiting"},
						},
					},
				},
			},
		},
	})

	require.Contains(t, output, "        └── api-0 (running/Thu Jan  1 00:00:00 UTC 1970) 2 container/s:\n")
	require.Contains(t, output, "            ├── 0: web (running)\n")
	require.Contains(t, output, "            └── 1: sidecar (waiting)\n")
}

func TestViewNodeDataPrintoutGroupedNamespacesKeepInlineContainerOutput(t *testing.T) {
	output := captureViewOutput(t, ViewNodeData{
		Config: ViewNodeDataConfig{
			ShowNamespaces:       true,
			GroupPodsByNamespace: true,
			ShowContainers:       true,
			ContainerViewType:    Inline,
			ShowReqLimits:        true,
		},
		Nodes: []ViewNode{
			{Name: ""},
			{
				Name:             "worker-a",
				Os:               "linux",
				Arch:             "amd64",
				ContainerRuntime: "containerd://2.2.0",
				Pods: []ViewPod{
					{
						Name:      "api-0",
						Namespace: "team-a",
						Phase:     "Running",
						Containers: []ViewContainer{
							{Name: "web", State: "Running", CpuReq: "100m", CpuLimit: "200m", MemoryReq: "128Mi", MemoryLimit: "256Mi"},
						},
					},
				},
			},
		},
	})

	require.Contains(t, output, "    └── team-a\n")
	require.Contains(t, output, "        └── api-0 (running) (1: web/running [cr:100m<200m mr:128Mi<256Mi])\n")
}
