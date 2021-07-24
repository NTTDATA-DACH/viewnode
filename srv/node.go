package srv

import (
	"errors"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"strings"
)

type NodeFilter struct {
	SearchText string
	Api        Api
}

func (nf NodeFilter) Transform() (result []ViewNode, err error) {
	list, err := nf.Api.RetrieveNodeList()
	if err != nil {
		return nil, err
	}
	vns := make([]ViewNode, len(list.Items))
	for i, n := range list.Items {
		vn := ViewNode{
			Name: n.Name,
			Os:   n.Status.NodeInfo.OperatingSystem,
			Arch: n.Status.NodeInfo.Architecture,
		}
		vns[i] = vn
	}
	return vns, err
}

func (nf NodeFilter) Filter(vns []ViewNode) (result []ViewNode, err error) {
	if vns == nil {
		return nil, errors.New("view node array must not be empty")
	}
	if nf.SearchText == "" {
		return vns, nil
	}
	res := make([]ViewNode, 0, len(vns))
	for _, vn := range vns {
		if strings.Contains(vn.Name, nf.SearchText) {
			res = append(res, vn)
		}
	}
	return res, nil
}
