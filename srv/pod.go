package srv

import (
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
)

type PodFilter struct {
	Namespace   string
	SearchText  string
	RunningOnly bool
	WithMetrics bool
	Api         Api
}

func (pf PodFilter) LoadAndFilter(vns []ViewNode) (result []ViewNode, err error) {
	list, err := pf.Api.RetrievePodList(pf.Namespace)
	if err != nil {
		return nil, err
	}
	if vns == nil {
		vns = createNodeListFromPodSpec(*list)
	}
	for i := range vns {
		for _, p := range list.Items {
			if !(strings.Contains(p.Name, pf.SearchText) && vns[i].Name == p.Spec.NodeName) {
				continue
			}
			if pf.RunningOnly && string(p.Status.Phase) != "Running" {
				continue
			}
			if vns[i].Pods == nil {
				vns[i].Pods = make([]ViewPod, 0)
			}
			st := time.Now()
			if p.Status.StartTime != nil {
				st = p.Status.StartTime.Time
			}
			cn := "unknown"
			for _, c := range p.Status.Conditions {
				if c.Status == "True" {
					cn = string(c.Type)
				}
			}
			vp := ViewPod{
				Name:      p.Name,
				Phase:     string(p.Status.Phase),
				Condition: cn,
				Namespace: p.Namespace,
				StartTime: st,
			}
			vp.Containers = make([]ViewContainer, len(p.Status.ContainerStatuses))
			for j, cs := range p.Status.ContainerStatuses {
				vp.Containers[j].Name = cs.Name
				if cs.State.Waiting != nil {
					vp.Containers[j].State = "Waiting"
				}
				if cs.State.Running != nil {
					vp.Containers[j].State = "Running"
				}
				if cs.State.Terminated != nil {
					vp.Containers[j].State = "Terminated"
				}
			}
			for j := range vp.Containers {
				for _, spc := range p.Spec.Containers {
					if vp.Containers[j].Name == spc.Name {
						vp.Containers[j].CpuLimit = spc.Resources.Limits.Cpu().String()
						vp.Containers[j].CpuReq = spc.Resources.Requests.Cpu().String()
						vp.Containers[j].MemoryLimit = spc.Resources.Limits.Memory().String()
						vp.Containers[j].MemoryReq = spc.Resources.Requests.Memory().String()
						continue
					}
				}
			}
			vns[i].Pods = append(vns[i].Pods, vp)
		}
	}
	if pf.WithMetrics {
		pml, err := pf.Api.RetrievePodMetricses(pf.Namespace)
		if err != nil {
			return nil, err
		}
		for i := range vns {
			if vns[i].Name == "" {
				continue
			}
			for j := range vns[i].Pods {
				for _, p := range pml.Items {
					if vns[i].Pods[j].Name == p.Name {
						var pm int64
						for k := range vns[i].Pods[j].Containers {
							for _, c := range p.Containers {
								if vns[i].Pods[j].Containers[k].Name == c.Name {
									vns[i].Pods[j].Containers[k].Metrics.Memory = c.Usage.Memory().Value()
									pm = pm + c.Usage.Memory().Value()
								}
							}
						}
						vns[i].Pods[j].Metrics.Memory = pm
						break
					}
				}
			}
		}
	}
	return vns, nil
}

func (pf PodFilter) ResourceName() string {
	return "pod"
}

func createNodeListFromPodSpec(list v1.PodList) []ViewNode {
	vns := make([]ViewNode, 1)
	vns[0].Name = "" // create placeholder for unscheduled pods
outer:
	for _, p := range list.Items {
		for _, n := range vns {
			if n.Name == p.Spec.NodeName {
				continue outer
			}
		}
		vn := ViewNode{
			Name: p.Spec.NodeName,
			Os:   "na",
			Arch: "na",
		}
		vns = append(vns, vn)
	}
	return vns
}
