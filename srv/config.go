package srv

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	"path/filepath"
)

type Setup struct {
	Namespace        string
	Clientset        *kubernetes.Clientset
	MetricsClientset *metricsv.Clientset
}

func InitSetup() (*Setup, error) {
	ns, err := retrieveNamespace()
	if err != nil {
		return nil, err
	}
	cs, err := InitCoreSetup()
	if err != nil {
		return nil, err
	}
	mcs, err := InitMetricsSetup()
	if err != nil {
		return nil, err
	}
	s := Setup{
		Namespace:        ns,
		Clientset:        cs,
		MetricsClientset: mcs,
	}
	return &s, nil
}

func InitCoreSetup() (*kubernetes.Clientset, error) {
	config, err := retrieveConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func InitMetricsSetup() (*metricsv.Clientset, error) {
	config, err := retrieveConfig()
	if err != nil {
		return nil, err
	}
	return metricsv.NewForConfig(config)
}

func retrieveConfig() (*rest.Config, error) {
	kubeconfig := ""
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

func retrieveNamespace() (string, error) {
	config, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return "", err
	}
	namespace := config.Contexts[config.CurrentContext].Namespace
	if namespace == "" {
		namespace = "default"
	}
	return namespace, nil
}
