package srv

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Setup struct {
	Namespace string
	Clientset *kubernetes.Clientset
}

func GetCurrentNamespaceAndClientset() (*Setup, error) {
	cclr := clientcmd.NewDefaultClientConfigLoadingRules()
	co := clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(cclr, &co)
	namespace, _, err := kubeConfig.Namespace()
	if err != nil {
		return nil, err
	}
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	s := Setup{
		Namespace: namespace,
		Clientset: clientset,
	}
	return &s, nil
}
