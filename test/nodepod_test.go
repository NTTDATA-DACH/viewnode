package test

import (
	v1 "k8s.io/api/core/v1"
	"kubectl-viewnode/srv"
	"testing"
)

type MockApi struct{}

const NodesCount = 2
const NodeName1 = "Node 1"
const NodeName2 = "Node 2"
const PodsCount = 3
const PodName1 = "Pod 1"
const PodName2 = "Pod 2"
const PodName3 = "Pod 3"

func (m MockApi) RetrieveNodeList() (*v1.NodeList, error) {
	var nl v1.NodeList
	nl.Items = make([]v1.Node, NodesCount)
	nl.Items[0].Name = NodeName1
	nl.Items[1].Name = NodeName2
	return &nl, nil
}

func (m MockApi) RetrievePodList(namespace string) (*v1.PodList, error) {
	var pl v1.PodList
	pl.Items = make([]v1.Pod, PodsCount)
	pl.Items[0].Name = PodName1
	pl.Items[0].Spec.NodeName = NodeName1
	pl.Items[1].Name = PodName2
	pl.Items[1].Spec.NodeName = NodeName1
	pl.Items[2].Name = PodName3
	pl.Items[2].Spec.NodeName = NodeName1
	return &pl, nil
}

func TestNodeLoadAndFilter(t *testing.T) {
	var api MockApi
	nf := srv.NodeFilter{
		SearchText: NodeName1,
		Api: api,
	}
	vns, err := nf.LoadAndFilter(nil)
	if err != nil {
		t.Error(err.Error())
	}
	const expectedNoNodes = 1
	if len(vns) != expectedNoNodes {
		t.Fatalf("Loading and filtering of nodes was not correct. Got: %d, expected: %d nodes.", len(vns), expectedNoNodes)
	}
	if vns[0].Name != NodeName1 {
		t.Errorf("Loading and filtering of nodes was not correct for the node name. Got: '%s', expected: '%s'", vns[0].Name, NodeName1)
	}
}

func TestPodLoadAndFilter(t *testing.T) {
	var api MockApi
	nf := srv.NodeFilter{
		Api: api,
	}
	vns, err := nf.LoadAndFilter(nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	pf := srv.PodFilter{
		Namespace: "",
		SearchText: PodName2,
		Api: api,
	}
	vns, err = pf.LoadAndFilter(vns)
	if err != nil {
		t.Errorf(err.Error())
	}
	const expectedNoPods = 1
	l := len(vns[0].Pods)
	if l != expectedNoPods {
		t.Errorf("Loading and filtering of pods was not correct. Got: %d, expected %d pod.", l, expectedNoPods)
	}
}