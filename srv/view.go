package srv

import (
	"fmt"
)

type ViewNode struct {
	Name string
	Os   string
	Arch string
}

func PrintNodes(viewNodesData []ViewNode) {
	for _, vnd := range viewNodesData {
		fmt.Printf("* %s (%s/%s)\n", vnd.Name, vnd.Os, vnd.Arch)
	}
}
