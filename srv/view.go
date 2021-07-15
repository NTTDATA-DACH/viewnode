package srv

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"strings"
)

type ViewNode struct {
	Name string
	Os string
	Arch string
}

func GetViewNodeData(nodes *v1.NodeList, filter string) []ViewNode {
	var vns []ViewNode
	for _, n := range nodes.Items {
		vn := ViewNode{
			Name: n.Name,
			Os: n.Status.NodeInfo.OperatingSystem,
			Arch: n.Status.NodeInfo.Architecture,
		}
		if filter == "" {
			vns = append(vns, vn)
		}
		if (filter != "") && (strings.Contains(n.Name, filter)) {
			vns = append(vns, vn)
		}
	}
	return vns
}

func PrintNodes(viewNodesData []ViewNode)  {
	for _, vnd := range viewNodesData {
		fmt.Printf("* %s (%s/%s)\n", vnd.Name, vnd.Os, vnd.Arch)
	}
}
