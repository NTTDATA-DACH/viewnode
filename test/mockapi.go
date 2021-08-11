package test

import v1 "k8s.io/api/core/v1"

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

