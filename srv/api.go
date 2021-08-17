package srv

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type Api interface {
	RetrieveNodeList() (*v1.NodeList, error)
	RetrievePodList(namespace string) (*v1.PodList, error)
}

type KubernetesApi struct {
	Setup *Setup
}

func (k KubernetesApi) RetrieveNodeList() (*v1.NodeList, error) {
	return k.Setup.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
}

func (k KubernetesApi) RetrievePodList(namespace string) (*v1.PodList, error) {
	return k.Setup.Clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
}
