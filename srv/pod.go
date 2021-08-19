package srv

import (
	"errors"
	"strings"
)

type PodFilter struct {
	Namespace  string
	SearchText string
	Api        Api
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
			if strings.Contains(p.Name, pf.SearchText) && vns[i].Name == p.Spec.NodeName {
				if vns[i].Pods == nil {
					vns[i].Pods = make([]ViewPod, 0)
				}
				vp := ViewPod{
					Name:      p.Name,
					Phase:     string(p.Status.Phase),
					Namespace: p.Namespace,
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
				vns[i].Pods = append(vns[i].Pods, vp)
			}
		}
	}
	return vns, nil
}
