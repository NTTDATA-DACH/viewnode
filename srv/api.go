package srv

import (
	"context"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type Api interface {
	RetrieveNodeList() (*v1.NodeList, error)
	RetrievePodList(namespace string) (*v1.PodList, error)
	RetrieveNodeMetricses() (*v1beta1.NodeMetricsList, error)
	RetrievePodMetricses(namespace string) (*v1beta1.PodMetricsList, error)
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

func (k KubernetesApi) RetrieveNodeMetricses() (*v1beta1.NodeMetricsList, error) {
	return k.Setup.MetricsClientset.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
}

func (k KubernetesApi) RetrievePodMetricses(namespace string) (*v1beta1.PodMetricsList, error) {
	return k.Setup.MetricsClientset.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{})
}