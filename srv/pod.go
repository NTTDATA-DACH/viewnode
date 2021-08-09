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
	for _, n := range vns {
		for _, p := range list.Items {
			if strings.Contains(p.Name, pf.SearchText) && n.Name == p.Spec.NodeName {
				if n.Pods == nil {
					n.Pods = make([]ViewPod, 0)
				}
				pn := ViewPod{
					Name: p.Name,
				}
				n.Pods = append(n.Pods, pn)
			}
		}
	}
	return vns, nil
}
