package srv

import (
	"errors"
	"strings"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var ErrMetricsServerNotInstalled = errors.New("could not find resource metrics.k8s.io; is the metrics server installed?")

type NodeFilter struct {
	SearchText  string
	WithMetrics bool
	Api         Api
}

func (nf NodeFilter) LoadAndFilter(vns []ViewNode) (result []ViewNode, err error) {
	list, err := nf.Api.RetrieveNodeList()
	if err != nil {
		return nil, err
	}
	if vns == nil {
		vns = make([]ViewNode, 1, len(list.Items)+1)
		vns[0].Name = "" // create placeholder for unscheduled pods
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
	if nf.WithMetrics {
		nml, err := nf.Api.RetrieveNodeMetricses()
		if err != nil {
			if strings.Contains(err.Error(), "(get nodes.metrics.k8s.io)") {
				return vns, ErrMetricsServerNotInstalled
			}
			return nil, err
		}
		for i := range vns {
			if vns[i].Name == "" {
				continue
			}
			for _, nm := range nml.Items {
				if vns[i].Name == nm.Name {
					vns[i].Metrics.Memory = nm.Usage.Memory().Value()
					break
				}
			}
		}
	}
	return vns, nil
}

func (nf NodeFilter) ResourceName() string {
	return "node"
}
