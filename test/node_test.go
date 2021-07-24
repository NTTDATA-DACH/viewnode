package test

import (
	v1 "k8s.io/api/core/v1"
	"kubectl-viewnode/srv"
	"testing"
)

type MockApi struct {}

const NodesCount = 10
const NodeName1 = "Node 1"

func (m MockApi) RetrieveNodeList() (*v1.NodeList, error) {
	var nl v1.NodeList
	nl.Items = make([]v1.Node, NodesCount)
	nl.Items[0].Name = NodeName1
	return &nl, nil
}

func TestNodeTransform(t *testing.T) {
	var api MockApi
	nf := srv.NodeFilter{
		Api: api,
	}
	vns, err := nf.Transform()
	if err != nil {
		t.Error(err.Error())
	}
	if len(vns) != NodesCount {
		t.Errorf("Transformation was not correct. Got: %d, expected: %d nodes.", len(vns), NodesCount)
	}
	if vns[0].Name != NodeName1 {
		t.Errorf("Transformation was not correct for node name. Got: %s, expected: %s", vns[0].Name, NodeName1)
	}
}
