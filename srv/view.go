package srv

import (
	"errors"
	"fmt"
)

type ViewNode struct {
	Name string
	Os   string
	Arch string
}

func PrintNodes(viewNodesData []ViewNode) error {
	if viewNodesData == nil {
		return errors.New("list of view nodes must not be null")
	}
	l := len(viewNodesData)
	if l == 0 {
		fmt.Println("no nodes to display...")
		return nil
	}
	fmt.Printf("viewing %d node(s):\n", l)
	for _, vnd := range viewNodesData {
		fmt.Printf("* %s (%s/%s)\n", vnd.Name, vnd.Os, vnd.Arch)
	}
	return nil
}
