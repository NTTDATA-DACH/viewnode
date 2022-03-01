package srv

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"viewnode/utils"
)

type ViewNode struct {
	Name    string
	Os      string
	Arch    string
	Metrics ViewMetrics
	Pods    []ViewPod
}

type ViewPod struct {
	Name       string
	Phase      string
	Condition  string
	Namespace  string
	StartTime  time.Time
	Containers []ViewContainer
}

type ViewContainer struct {
	Name        string
	State       string
	MemoryLimit string
	MemoryReq   string
	CpuLimit    string
	CpuReq      string
}

type ViewMetrics struct {
	Memory int64
}

type ViewNodeData struct {
	Config ViewNodeDataConfig
	Nodes  []ViewNode
}

type ViewType int

const (
	Inline ViewType = iota
	Block
)

type ViewNodeDataConfig struct {
	ShowNamespaces bool
	ShowContainers bool
	ShowTimes      bool
	ShowReqLimits  bool
	ShowMetrics    bool

	ContainerViewType ViewType
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
			fmt.Printf("- %s running %d pod(s) (%s/%s", n.Name, len(n.Pods), n.Os, n.Arch)
			if vnd.Config.ShowMetrics {
				fmt.Printf(" | mem: %s", utils.ByteCountIEC(n.Metrics.Memory))
			}
			fmt.Println(")")
			for _, p := range n.Pods {
				if vnd.Config.ShowNamespaces {
					fmt.Printf("  * %s: %s (%s", p.Namespace, p.Name, strings.ToLower(p.Phase))
				} else {
					fmt.Printf("  * %s (%s", p.Name, strings.ToLower(p.Phase))
				}
				if vnd.Config.ShowTimes {
					fmt.Printf("/%s", p.StartTime.Format(time.UnixDate))
				}
				fmt.Printf(")")
				if vnd.Config.ShowContainers {
					switch vnd.Config.ContainerViewType {
					case Inline:
						fmt.Printf(" (%d:", len(p.Containers))
						for _, c := range p.Containers {
							fmt.Printf(" %s/%s", c.Name, strings.ToLower(c.State))
							if vnd.Config.ShowReqLimits {
								fmt.Printf(" [C:%s M:%s]", fmtRes(c.CpuReq, c.CpuLimit), fmtRes(c.MemoryReq, c.MemoryLimit))
							}
						}
						fmt.Printf(")")
						break
					case Block:
						fmt.Printf(" %d container/s:", len(p.Containers))
						for i, c := range p.Containers {
							fmt.Printf("\n    %d: %s (%s)", i, c.Name, strings.ToLower(c.State))
							if vnd.Config.ShowReqLimits {
								fmt.Printf(" [Cpu: %s Memory: %s]", fmtRes(c.CpuReq, c.CpuLimit), fmtRes(c.MemoryReq, c.MemoryLimit))
							}
						}
						break
					}
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
