package srv

import (
	"context"
	"errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"strings"
)

func GetNodes(clientset *kubernetes.Clientset) (*v1.NodeList, error) {
	return clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
}

type NodeFilter struct {
	SearchText string
	Source     *v1.NodeList
}

func (nf NodeFilter) Transform() (result []ViewNode, err error) {
	if nf.Source == nil {
		return nil, errors.New("node list must not be nil")
	}
	vns := make([]ViewNode, len(nf.Source.Items))
	for i, n := range nf.Source.Items {
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
