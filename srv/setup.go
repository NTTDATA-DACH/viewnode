package srv

import (
	"errors"
	"fmt"
	"io/fs"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	"os"
	"strings"
)

// ErrKubeCfgFileNotExist is a sentinel (expected) error, when the config file cannot be found at the given path
type ErrKubeCfgFileNotExist struct {
	KubeCfgPath string
}

const errKubeCfgFileNotExistErrorPrefix = "config file not found at the following path:"

func (e ErrKubeCfgFileNotExist) Error() string {
	return fmt.Sprintf("%s %s", errKubeCfgFileNotExistErrorPrefix, e.KubeCfgPath)
}

func (e ErrKubeCfgFileNotExist) Is(target error) bool {
	if _, ok := target.(ErrKubeCfgFileNotExist); ok {
		return strings.HasPrefix(target.Error(), errKubeCfgFileNotExistErrorPrefix)
	}
	return false
}

type Setup struct {
	KubeCfgPath      string
	KubeContext      string
	KubeCluster      string
	ClientConfig     clientcmd.ClientConfig
	Namespace        string
	Clientset        *kubernetes.Clientset
	MetricsClientset *metricsv.Clientset
}

// Initialize initializes setup struct by setting config, and clientsets objects
func (s *Setup) Initialize() error {
	if s.ClientConfig == nil {
		if err := s.DetermineClientConfig(); err != nil {
			return err
		}
	}
	config, err := s.GetRestConfig()
	if err != nil {
		return err
	}
	if s.Clientset, err = kubernetes.NewForConfig(config); err != nil {
		return err
	}
	if s.MetricsClientset, err = metricsv.NewForConfig(config); err != nil {
		return err
	}
	if s.Namespace, err = s.GetCurrentNamespace(); err != nil {
		return err
	}
	return nil
}

// GetRestConfig returns a complete client config
func (s *Setup) GetRestConfig() (*rest.Config, error) {
	if s.ClientConfig == nil {
		if err := s.DetermineClientConfig(); err != nil {
			return nil, err
		}
	}
	config, err := s.ClientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

// GetCurrentNamespace returns the current namespace which is either provided as option or picked up from kubeconfig
func (s *Setup) GetCurrentNamespace() (string, error) {
	var err error
	if s.ClientConfig == nil {
		if err = s.DetermineClientConfig(); err != nil {
			return "", err
		}
	}
	name, _, err := s.ClientConfig.Namespace()
	return name, err
}

// DetermineClientConfig determines clients configuration by looking for the config file using different path queues
func (s *Setup) DetermineClientConfig() error {
	co := &clientcmd.ConfigOverrides{}
	if s.KubeContext != "" {
		co.CurrentContext = s.KubeContext
	}
	if s.KubeCluster != "" {
		co.Context.Cluster = s.KubeCluster
	}

	lr := clientcmd.NewDefaultClientConfigLoadingRules()
	if len(s.KubeCfgPath) == 0 {
		s.ClientConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(lr, co)
		return nil
	}

	_, err := os.Stat(s.KubeCfgPath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		return ErrKubeCfgFileNotExist{KubeCfgPath: s.KubeCfgPath}
	}

	lr.ExplicitPath = s.KubeCfgPath
	s.ClientConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(lr, co)
	return nil
}
