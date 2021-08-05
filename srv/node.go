package srv

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"strings"
)

type NodeFilter struct {
	SearchText string
	Api        Api
}

func (nf NodeFilter) LoadAndFilter() (result []ViewNode, err error) {
	list, err := nf.Api.RetrieveNodeList()
	if err != nil {
		return nil, err
	}
	vns := make([]ViewNode, 0, len(list.Items))
	for _, n := range list.Items {
		if strings.Contains(n.Name, nf.SearchText) {
			vn := ViewNode{
				Name: n.Name,
				Os:   n.Status.NodeInfo.OperatingSystem,
				Arch: n.Status.NodeInfo.Architecture,
			}
			vns = append(vns, vn)
		}
	}
	return vns, err
}