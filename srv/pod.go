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
				pn := ViewPod{
					Name:      p.Name,
					Phase:     string(p.Status.Phase),
					Namespace: p.Namespace,
				}
				vns[i].Pods = append(vns[i].Pods, pn)
			}
		}
	}
	return vns, nil
}
