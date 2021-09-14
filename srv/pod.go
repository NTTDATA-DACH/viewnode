package srv

import (
	"errors"
	"strings"
	"time"
)

type PodFilter struct {
	Namespace   string
	SearchText  string
	RunningOnly bool
	Api         Api
}

func (pf PodFilter) LoadAndFilter(vns []ViewNode) (result []ViewNode, err error) {
	if vns == nil {
		return nil, errors.New("list of view nodes must not be nil")
	}
	list, err := pf.Api.RetrievePodList(pf.Namespace)
	if err != nil {
		return nil, err
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
	return vns, nil
}
