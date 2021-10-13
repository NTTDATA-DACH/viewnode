package srv

import (
	"strings"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type NodeFilter struct {
	SearchText string
	Api        Api
}

func (nf NodeFilter) LoadAndFilter(vns []ViewNode) (result []ViewNode, err error) {
	list, err := nf.Api.RetrieveNodeList()
	if err != nil {
		return nil, err
	}
	if vns == nil {
		vns = make([]ViewNode, 1, len(list.Items)+1)
		vns[0].Name = ""
	}
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
	return vns, nil
}

func (nf NodeFilter) ResourceName() string {
	return "node"
}
