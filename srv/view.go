package srv

import (
	"errors"
	"fmt"
	"strings"
)

type ViewNode struct {
	Name string
	Os   string
	Arch string
	Pods []ViewPod
}

type ViewPod struct {
	Name       string
	Phase      string
	Namespace  string
	Containers []ViewContainer
}

type ViewContainer struct {
	Name  string
	State string
}

type ViewNodeData struct {
	Config ViewNodeDataConfig
	Nodes  []ViewNode
}

type ViewNodeDataConfig struct {
	CanShowNamespaces bool
	CanShowContainers bool
}

type View interface {
	Printout() error
}

func (vnd ViewNodeData) Printout() error {
	if vnd.Nodes == nil {
		return errors.New("list of view nodes must not be null")
	}
	l := len(vnd.Nodes)
	if l <= 1 {
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
	for _, p := range vnd.Nodes[0].Pods {
		fmt.Printf("  * %s (%s)\n", p.Name, strings.ToLower(p.Phase))
	}
	fmt.Printf("%d running node(s) with %d scheduled pod(s):\n", l-1, nsp)
	for _, n := range vnd.Nodes {
		if n.Name != "" {
			fmt.Printf("- %s running %d pod(s) (%s/%s)\n", n.Name, len(n.Pods), n.Os, n.Arch)
			for _, p := range n.Pods {
				if vnd.Config.CanShowNamespaces {
					fmt.Printf("  * %s: %s (%s)", p.Namespace, p.Name, strings.ToLower(p.Phase))
				} else {
					fmt.Printf("  * %s (%s)", p.Name, strings.ToLower(p.Phase))
				}
				if vnd.Config.CanShowContainers {
					fmt.Printf(" (%d:", len(p.Containers))
					for _, c := range p.Containers {
						fmt.Printf(" %s/%s", c.Name, strings.ToLower(c.State))
					}
					fmt.Printf(")")
				}
				fmt.Println()
			}
		}
	}
	return nil
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
