package srv

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
	"viewnode/utils"
)

type ViewNode struct {
	Name             string
	Os               string
	Arch             string
	ContainerRuntime string
	Metrics          ViewMetrics
	Pods             []ViewPod
}

type ViewPod struct {
	Name       string
	Phase      string
	Condition  string
	Namespace  string
	StartTime  time.Time
	Metrics    ViewMetrics
	Containers []ViewContainer
}

type ViewContainer struct {
	Name        string
	State       string
	MemoryLimit string
	MemoryReq   string
	CpuLimit    string
	CpuReq      string
	Metrics     ViewMetrics
}

type ViewMetrics struct {
	Memory int64
}

type ViewNodeData struct {
	Config     ViewNodeDataConfig
	Namespace  string
	NodeFilter string
	Nodes      []ViewNode
}

type ViewType int

const (
	Inline ViewType = iota
	Tree
)

type ViewNodeDataConfig struct {
	ShowNamespaces       bool
	GroupPodsByNamespace bool
	SelectedNamespaces   []string
	ShowContainers       bool
	ShowTimes            bool
	ShowReqLimits        bool
	ShowMetrics          bool

	ContainerViewType ViewType
}

type View interface {
	Printout(cls bool) error
}

type namespaceGroup struct {
	Namespace string
	Pods      []ViewPod
}

func (vnd ViewNodeData) Printout(cls bool) error {
	if cls {
		fmt.Print("\033[2J\033[0;0H")
	}
	if vnd.Nodes == nil {
		return errors.New("list of view nodes must not be null")
	}
	if vnd.Namespace != "" {
		fmt.Printf("namespace(s): %s\n", vnd.Namespace)
	}
	l := len(vnd.Nodes)
	if l <= 1 {
		if vnd.NodeFilter != "" {
			fmt.Printf("no nodes matched filter %q\n", vnd.NodeFilter)
			return nil
		}
		fmt.Println("no nodes to display...")
		return nil
	}
	nup := vnd.getNumberOfUnscheduledPods()
	nsp := vnd.getNumberOfScheduledPods()
	fmt.Printf("%d pod(s) in total\n", nup+nsp)
	fmt.Printf("%d unscheduled pod(s)", nup)
	if nup > 0 {
		fmt.Printf(":")
	}
	fmt.Printf("\n")
	for i, p := range vnd.Nodes[0].Pods {
		prefix := "├──"
		if i == len(vnd.Nodes[0].Pods)-1 {
			prefix = "└──"
		}
		fmt.Printf("  %s %s (%s)\n", prefix, p.Name, strings.ToLower(p.Phase))
	}
	nodes := make([]ViewNode, 0, l-1)
	for _, n := range vnd.Nodes[1:] {
		if n.Name != "" {
			nodes = append(nodes, n)
		}
	}
	fmt.Printf("%d running node(s) with %d scheduled pod(s):\n", len(nodes), nsp)
	for ni, n := range nodes {
		nodePrefix := "├──"
		podIndent := "│   "
		if ni == len(nodes)-1 {
			nodePrefix = "└──"
			podIndent = "    "
		}
		fmt.Printf("%s %s running %d pod(s) (%s/%s/%s", nodePrefix, n.Name, len(n.Pods), n.Os, n.Arch, n.ContainerRuntime)
		if vnd.Config.ShowMetrics {
			fmt.Printf(" | mem: %s", utils.ByteCountIEC(n.Metrics.Memory))
		}
		fmt.Println(")")
		if vnd.Config.GroupPodsByNamespace {
			vnd.printNamespaceGroupedPods(n.Pods, podIndent)
			continue
		}
		showNamespaceInline := vnd.shouldShowNamespaceInline()
		for pi, p := range n.Pods {
			vnd.printPod(p, podIndent, pi == len(n.Pods)-1, showNamespaceInline)
		}
	}
	return nil
}

func (vnd ViewNodeData) printNamespaceGroupedPods(pods []ViewPod, podIndent string) {
	groups := groupPodsByNamespace(pods, vnd.Config.SelectedNamespaces)
	for gi, group := range groups {
		namespacePrefix := "├──"
		namespacePodIndent := podIndent + "│   "
		if gi == len(groups)-1 {
			namespacePrefix = "└──"
			namespacePodIndent = podIndent + "    "
		}
		fmt.Printf("%s%s %s\n", podIndent, namespacePrefix, group.Namespace)
		for pi, pod := range group.Pods {
			vnd.printPod(pod, namespacePodIndent, pi == len(group.Pods)-1, false)
		}
	}
}

func groupPodsByNamespace(pods []ViewPod, selectedNamespaces []string) []namespaceGroup {
	selectedNamespaceSet := make(map[string]struct{}, len(selectedNamespaces))
	for _, namespace := range selectedNamespaces {
		selectedNamespaceSet[namespace] = struct{}{}
	}
	groupsByNamespace := make(map[string][]ViewPod, len(pods))
	namespaces := make([]string, 0, len(pods))
	seenNamespaces := make(map[string]struct{}, len(pods))
	for _, pod := range pods {
		if len(selectedNamespaceSet) > 0 {
			if _, ok := selectedNamespaceSet[pod.Namespace]; !ok {
				continue
			}
		}
		if _, ok := seenNamespaces[pod.Namespace]; !ok {
			namespaces = append(namespaces, pod.Namespace)
			seenNamespaces[pod.Namespace] = struct{}{}
		}
		groupsByNamespace[pod.Namespace] = append(groupsByNamespace[pod.Namespace], pod)
	}
	sort.Strings(namespaces)
	groups := make([]namespaceGroup, 0, len(namespaces))
	for _, namespace := range namespaces {
		groups = append(groups, namespaceGroup{
			Namespace: namespace,
			Pods:      groupsByNamespace[namespace],
		})
	}
	return groups
}

func (vnd ViewNodeData) shouldShowNamespaceInline() bool {
	return vnd.Config.ShowNamespaces && !vnd.Config.GroupPodsByNamespace
}

func (vnd ViewNodeData) printPod(p ViewPod, podIndent string, isLast bool, showNamespaceInline bool) {
	podPrefix := "├──"
	containerIndent := podIndent + "│   "
	if isLast {
		podPrefix = "└──"
		containerIndent = podIndent + "    "
	}
	fmt.Printf("%s%s ", podIndent, podPrefix)
	if showNamespaceInline {
		fmt.Printf("%s: %s (%s", p.Namespace, p.Name, strings.ToLower(p.Phase))
	} else {
		fmt.Printf("%s (%s", p.Name, strings.ToLower(p.Phase))
	}
	if vnd.Config.ShowTimes {
		fmt.Printf("/%s", p.StartTime.Format(time.UnixDate))
	}
	if vnd.Config.ShowMetrics {
		fmt.Printf(" | mem usage: %s", utils.ByteCountIEC(p.Metrics.Memory))
	}
	fmt.Printf(")")
	if vnd.Config.ShowContainers {
		switch vnd.Config.ContainerViewType {
		case Inline:
			fmt.Printf(" (%d:", len(p.Containers))
			for _, c := range p.Containers {
				fmt.Printf(" %s/%s", c.Name, strings.ToLower(c.State))
				if vnd.Config.ShowReqLimits {
					fmt.Printf(" [cr:%s mr:%s", fmtRes(c.CpuReq, c.CpuLimit), fmtRes(c.MemoryReq, c.MemoryLimit))
				}
				if vnd.Config.ShowMetrics {
					fmt.Printf(" mu:%s", utils.ByteCountIEC(c.Metrics.Memory))
				}
				if vnd.Config.ShowReqLimits || vnd.Config.ShowMetrics {
					fmt.Printf("]")
				}
			}
			fmt.Printf(")")
		case Tree:
			fmt.Printf(" %d container/s:", len(p.Containers))
			for i, c := range p.Containers {
				containerPrefix := "├──"
				if i == len(p.Containers)-1 {
					containerPrefix = "└──"
				}
				fmt.Printf("\n%s%s %d: %s (%s)", containerIndent, containerPrefix, i, c.Name, strings.ToLower(c.State))
				if vnd.Config.ShowReqLimits {
					fmt.Printf(" [cpu: %s | mem: %s", fmtRes(c.CpuReq, c.CpuLimit), fmtRes(c.MemoryReq, c.MemoryLimit))
				}
				if vnd.Config.ShowMetrics {
					if vnd.Config.ShowReqLimits {
						fmt.Printf(" | ")
					} else {
						fmt.Printf(" [")
					}
					fmt.Printf("mem usage: %s", utils.ByteCountIEC(c.Metrics.Memory))
				}
				if vnd.Config.ShowReqLimits || vnd.Config.ShowMetrics {
					fmt.Printf("]")
				}
			}
		}
	}
	fmt.Println()
}

func (vnd ViewNodeData) getNumberOfUnscheduledPods() int {
	if vnd.Nodes != nil {
		return len(vnd.Nodes[0].Pods)
	}
	return 0
}

func (vnd ViewNodeData) getNumberOfScheduledPods() int {
	var c int
	for _, n := range vnd.Nodes[1:] {
		c = c + len(n.Pods)
	}
	return c
}

func fmtRes(req, lim string) string {
	if req == "0" && lim == "0" {
		return "-"
	}
	if req == "0" {
		req = "-"
	}
	if lim == "0" {
		lim = "-"
	}
	return req + "<" + lim
}
