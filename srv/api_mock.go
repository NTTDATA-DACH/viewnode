package srv

import (
	"errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type ApiType int

const (
	NodeListError ApiType = iota + 1
	PodListError
	NodeMetricsesError
	PodMetricsesError
)

type MockApi struct {
	ApiTypeValue ApiType
}

const (
	NodesCount = 2
	NodeName1  = "Node 1"
	NodeName2  = "Node 2"
	PodsCount  = 3
	PodName1   = "Pod 1"
	PodName2   = "Pod 2"
	PodName3   = "Pod 3"
)

func (m MockApi) Init(at ApiType) {
	m.ApiTypeValue = at
}

func (m MockApi) RetrieveNodeList() (*v1.NodeList, error) {
	if m.ApiTypeValue == NodeListError {
		return nil, errors.New("node list error")
	}
	var nl v1.NodeList
	nl.Items = make([]v1.Node, NodesCount)
	nl.Items[0].Name = NodeName1
	nl.Items[1].Name = NodeName2
	return &nl, nil
}

func (m MockApi) RetrievePodList(namespace string) (*v1.PodList, error) {
	if m.ApiTypeValue == PodListError {
		return nil, errors.New("pod list error")
	}
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

func (m MockApi) RetrieveNodeMetricses() (*v1beta1.NodeMetricsList, error) {
	if m.ApiTypeValue == NodeMetricsesError {
		return nil, ErrMetricsServerNotInstalled
	}
	var nml v1beta1.NodeMetricsList
	nml.Items = make([]v1beta1.NodeMetrics, NodesCount)
	nml.Items[0].Name = NodeName1
	nml.Items[1].Name = NodeName2
	return &nml, nil
}

func (m MockApi) RetrievePodMetricses(namespace string) (*v1beta1.PodMetricsList, error) {
	if m.ApiTypeValue == PodMetricsesError {
		return nil, ErrMetricsServerNotInstalled
	}
	panic("implement me")
}
